package engine

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/AndrewDonelson/retroforge-engine/internal/audio"
	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
	"github.com/AndrewDonelson/retroforge-engine/internal/eventbus"
	"github.com/AndrewDonelson/retroforge-engine/internal/gamestate"
	"github.com/AndrewDonelson/retroforge-engine/internal/graphics"
	"github.com/AndrewDonelson/retroforge-engine/internal/input"
	"github.com/AndrewDonelson/retroforge-engine/internal/lua"
	"github.com/AndrewDonelson/retroforge-engine/internal/luabind"
	"github.com/AndrewDonelson/retroforge-engine/internal/network"
	"github.com/AndrewDonelson/retroforge-engine/internal/pal"
	"github.com/AndrewDonelson/retroforge-engine/internal/physics"
	"github.com/AndrewDonelson/retroforge-engine/internal/rendersoft"
	"github.com/AndrewDonelson/retroforge-engine/internal/runner"
	"github.com/AndrewDonelson/retroforge-engine/internal/scheduler"
)

// Engine wires together bus, scheduler/runner, and Lua VM for headless runs.
type Engine struct {
	Bus        *eventbus.Bus
	Sched      *scheduler.Scheduler
	Run        *runner.Runner
	VM         *lua.VM
	Ren        graphics.Renderer
	Pal        *pal.Manager
	Physics    *physics.World
	Network    *network.NetworkManager     // Multiplayer networking
	GSM        *gamestate.GameStateMachine // Game state machine
	sfxMap      cartio.SFXMap
	musicMap    cartio.MusicMap
	spritesMap  cartio.SpriteMap
	tilesetsMap map[string]cartio.TilesetMap // Map of tileset name -> tileset data
	tilemapsMap map[string]*cartio.TileMapData // Map of tilemap name -> tilemap data
	devMode     *DevMode // Development mode (only when loading from folder)
}

func New(targetFPS int) *Engine {
	bus := eventbus.New()
	sched := scheduler.New(targetFPS)
	run := runner.New(bus, sched)
	vm := lua.New()
	ren := rendersoft.New(480, 270)
	phys := physics.NewWorld(0, 9.8) // Default gravity: down (Y+) like real physics

	// Create game state machine (will be set to debug mode in dev mode)
	// isDebug=false means splash screen will show in release builds
	// Note: renderer and palette will be set later in registerLuaBindings
	gsm := gamestate.NewGameStateMachine(false, "RetroForge", "v1.0 Alpha", "RetroForge Team", nil, nil)

	e := &Engine{
		Bus:        bus,
		Sched:      sched,
		Run:        run,
		VM:         vm,
		Ren:        ren,
		Pal:        pal.NewManager(),
		Physics:    phys,
		Network:    network.NewNetworkManager(),
		GSM:        gsm,
		tilesetsMap: make(map[string]cartio.TilesetMap),
		tilemapsMap: make(map[string]*cartio.TileMapData),
	}
	// On each tick, call Lua update with dt seconds.
	bus.Subscribe("tick", func(v any) {
		if dt, ok := v.(time.Duration); ok {
			dtSec := dt.Seconds()
			
		// Debug logging disabled - see cmd/wasm/main.go for frame-level logging

		// Check for hot reload (development mode only)
		if e.devMode != nil && e.devMode.CheckForReload() {
			// Reload synchronously (blocks to prevent race condition with VM)
			// This is fast enough that it won't cause noticeable frame drops
			if err := e.ReloadCart(); err != nil {
				e.devMode.AddDebugLog(fmt.Sprintf("Reload failed: %v", err))
			}
			// Skip this frame's update/draw since we just reloaded
			// Next frame will use the new VM
			input.Step()
			return
		}

		// Skip VM calls if reload is in progress (shouldn't happen with sync reload, but safety check)
		if e.devMode != nil && e.devMode.IsReloading() {
			input.Step()
			return
		}

		// Step physics before Lua update
		e.Physics.Step()

		// Update network frame (for multiplayer sync)
		e.Network.UpdateFrame(dt)

		// Step input state FIRST (same as SDL) - copy cur → prev
		// This ensures btnp() works correctly: prev has old state, cur will be updated by input
		// Then HandleInput() checks: prev (old) vs cur (new) → btnp() detects edges correctly
		input.Step()

		// Use state machine if it has active states, otherwise fall back to direct Lua calls
		if e.GSM != nil {
			_, hasActiveState := e.GSM.GetActiveState()
			if hasActiveState {
				// Handle input AFTER Step() - now prev has old state, cur has new state
				// This matches SDL behavior exactly
				e.GSM.HandleInput()

				// Update and draw using state machine
				e.GSM.Update(dtSec)
				e.GSM.Draw()
				// Swap buffers after drawing is complete (double buffering)
				if e.Ren != nil {
					if swapper, ok := e.Ren.(interface{ SwapBuffers() }); ok {
						swapper.SwapBuffers()
					}
				}
				
				// Clear stale justPressed flags at end of frame (after HandleInput/Draw)
				input.ClearStaleJustPressed()
				} else {
					// No active state - this game doesn't use state machine (old-style Lua game)
					// OR the splash was popped intentionally (e.g., game uses direct _UPDATE/_DRAW)
					// Fallback to direct Lua calls - don't force splash again (would cause loop)
					if e.VM != nil && e.VM.L != nil {
						_ = e.VM.CallUpdate(dtSec)
						_ = e.VM.CallDraw()
						// Swap buffers after drawing is complete (double buffering)
						if e.Ren != nil {
							if swapper, ok := e.Ren.(interface{ SwapBuffers() }); ok {
								swapper.SwapBuffers()
							}
						}
						// Clear stale justPressed flags at end of frame
						input.ClearStaleJustPressed()
					}
				}
			} else {
				// No state machine at all - use direct Lua calls
				if e.VM != nil && e.VM.L != nil {
					_ = e.VM.CallUpdate(dtSec)
					_ = e.VM.CallDraw()
					// Swap buffers after drawing is complete (double buffering)
					if e.Ren != nil {
						if swapper, ok := e.Ren.(interface{ SwapBuffers() }); ok {
							swapper.SwapBuffers()
						}
					}
					// Clear stale justPressed flags at end of frame
					input.ClearStaleJustPressed()
				}
			}

			// NOTE: input.Step() is now called at the START of each frame (above)
			// This matches SDL behavior: Step() → Set buttons → HandleInput()
			// This ensures btnp() edge detection works correctly

			// Update debug stats (development mode only)
			if e.devMode != nil && e.devMode.IsEnabled() {
				fps := 1.0 / dtSec
				e.devMode.UpdateStats(fps, 0, 0) // Frame count and Lua memory would need more work
			}
		}
	})
	return e
}

func (e *Engine) Close() {
	if e.devMode != nil {
		e.devMode.Disable()
	}
	if e.Network != nil {
		e.Network.Close()
	}
	e.VM.Close()
}

// devModeAdapter adapts engine.DevMode to luabind.DevModeHandler interface
type devModeAdapter struct {
	devMode *DevMode
}

func (a *devModeAdapter) IsEnabled() bool {
	return a.devMode.IsEnabled()
}

func (a *devModeAdapter) AddDebugLog(msg string) {
	a.devMode.AddDebugLog(msg)
}

func (a *devModeAdapter) GetStats() interface{} {
	// Convert engine.DevStats to luabind.DevStats-compatible structure
	stats := a.devMode.GetStats()
	// Return struct with matching fields (luabind will type assert)
	return struct {
		FPS         float64
		FrameCount  int64
		LuaMemory   int64
		LoadTime    time.Duration
		LastReload  time.Time
		ReloadCount int
	}{
		FPS:         stats.FPS,
		FrameCount:  stats.FrameCount,
		LuaMemory:   stats.LuaMemory,
		LoadTime:    stats.LoadTime,
		LastReload:  stats.LastReload,
		ReloadCount: stats.ReloadCount,
	}
}

// LoadLuaSource loads script and calls init() if present.
// Note: Lua bindings, RegisterStateMachine and RegisterModuleImport should be called before this.
func (e *Engine) LoadLuaSource(src string) error {
	if err := e.VM.LoadString(src); err != nil {
		return err
	}
	return e.VM.CallInit()
}

// registerLuaBindings registers all Lua bindings (rf.*, game.*, module import).
// This should be called before LoadLuaSource.
func (e *Engine) registerLuaBindings() {
	// Update GSM with renderer and palette so splash/credits can draw
	if e.GSM != nil {
		e.GSM.SetRenderer(e.Ren)
		e.GSM.SetPalette(e.Pal)
	}

	if e.devMode != nil && e.devMode.IsEnabled() {
		// Create adapter that implements DevModeHandler interface
		devAdapter := &devModeAdapter{devMode: e.devMode}
		luabind.RegisterWithDevMode(e.VM.L, e.Ren, func(i int) (c [4]uint8) {
			col := e.Pal.Color(i)
			c[0] = col.R
			c[1] = col.G
			c[2] = col.B
			c[3] = col.A
			return
		}, e.Pal.Set, e.sfxMap, e.musicMap, e.spritesMap, e.tilemapsMap, e.Physics, luabind.NewState(), devAdapter, e.Network)
	} else {
		luabind.RegisterWithState(e.VM.L, e.Ren, func(i int) (c [4]uint8) {
			col := e.Pal.Color(i)
			c[0] = col.R
			c[1] = col.G
			c[2] = col.B
			c[3] = col.A
			return
		}, e.Pal.Set, e.sfxMap, e.musicMap, e.spritesMap, e.tilemapsMap, e.Physics, luabind.NewState(), e.Network)
	}

	// Register state machine (needed for game.* API)
	luabind.RegisterStateMachine(e.VM.L, e.GSM)

	// Register imgtool API
	luabind.RegisterImgToolAPI(e.VM.L)
}

// RunFrames advances N frames headlessly.
func (e *Engine) RunFrames(n int) {
	for i := 0; i < n; i++ {
		e.Run.Step()
	}
}

// LoadCartFromReader loads a .rfs from an io.ReaderAt.
func (e *Engine) LoadCartFromReader(r io.ReaderAt, size int64) error {
	// Stop all audio from previous cart before loading new one
	// This prevents music/sounds from previous cart continuing to play
	audio.StopAll()

	result, err := cartio.Read(r, size)
	if err != nil {
		return err
	}

	// Set palette from manifest if specified
	if result.Manifest.Palette != "" {
		e.Pal.Set(result.Manifest.Palette)
	}

	entryPath := "assets/" + result.Manifest.Entry
	src, ok := result.Files[entryPath]
	if !ok {
		// Provide more helpful error message
		var availableFiles []string
		for f := range result.Files {
			if len(f) > 7 && f[:7] == "assets/" {
				availableFiles = append(availableFiles, f)
			}
		}
		return fmt.Errorf("entry file not found: %s (manifest entry: %q). Available asset files: %v", entryPath, result.Manifest.Entry, availableFiles)
	}

	// Store SFX, Music, and Sprites for Lua bindings
	e.sfxMap = result.SFX
	e.musicMap = result.Music
	e.spritesMap = result.Sprites

	// Register Lua bindings first (creates rf table)
	e.registerLuaBindings()

	// Register module import with file map for cart mode - rf table now exists
	// Convert result.Files map to the format expected by module import
	fileMap := make(map[string][]byte)
	for path, data := range result.Files {
		// Remove "assets/" prefix for module import
		if len(path) > 7 && path[:7] == "assets/" {
			fileMap[path[7:]] = data
		} else {
			fileMap[path] = data
		}
	}
	luabind.RegisterModuleImportWithMap(e.VM.L, e.GSM, fileMap)

	// Load and execute Lua source - this will register states via rf.import() calls
	if err := e.LoadLuaSource(string(src)); err != nil {
		return err
	}

	// Start the state machine AFTER Lua has loaded and registered states
	// LoadLuaSource executes the Lua code, which may call rf.import() to register states
	// We need to start the state machine after those registrations are complete
	// This will show the splash screen, then transition to the initial state set by the game
	// Default to "menu" if no initial state is set by the game
	if e.GSM != nil {
		initialState := "menu" // Default initial state
		// Note: GSM.Start() will transition to splash screen (if not debug mode),
		// which will then transition to initialState once that state is registered
		if startErr := e.GSM.Start(initialState); startErr != nil {
			// If Start fails, the state machine stack might be empty
			// This means GetActiveState() will return false, and we'll fall back to direct Lua calls
			// That's ok - the game might not use state machine
			if e.devMode != nil {
				e.devMode.AddDebugLog(fmt.Sprintf("State machine start warning: %v", startErr))
			}
		}
		// NOTE: We no longer force the splash screen if no active state is found.
		// The splash screen will be added by Start() if needed (non-debug builds).
		// If the stack becomes empty after the splash pops itself (e.g., for games that
		// don't use the state machine), we fall back to direct Lua calls in the tick handler.
		// Re-adding the splash here would cause an infinite loop.
	}

	return nil
}

// LoadCartFile opens .rfs by path and loads it.
func (e *Engine) LoadCartFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return err
	}
	return e.LoadCartFromReader(f, st.Size())
}
