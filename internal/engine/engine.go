package engine

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
	"github.com/AndrewDonelson/retroforge-engine/internal/eventbus"
	"github.com/AndrewDonelson/retroforge-engine/internal/gamestate"
	"github.com/AndrewDonelson/retroforge-engine/internal/graphics"
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
	sfxMap     cartio.SFXMap
	musicMap   cartio.MusicMap
	spritesMap cartio.SpriteMap
	devMode    *DevMode // Development mode (only when loading from folder)
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
	gsm := gamestate.NewGameStateMachine(false, "RetroForge", "1.0.0", "RetroForge Team", nil, nil)

	e := &Engine{
		Bus:     bus,
		Sched:   sched,
		Run:     run,
		VM:      vm,
		Ren:     ren,
		Pal:     pal.NewManager(),
		Physics: phys,
		Network: network.NewNetworkManager(),
		GSM:     gsm,
	}
	// On each tick, call Lua update with dt seconds.
	bus.Subscribe("tick", func(v any) {
		if dt, ok := v.(time.Duration); ok {
			dtSec := dt.Seconds()

			// Check for hot reload (development mode only)
			if e.devMode != nil && e.devMode.CheckForReload() {
				// Reload in background (don't block)
				go func() {
					if err := e.ReloadCart(); err != nil {
						e.devMode.AddDebugLog(fmt.Sprintf("Reload failed: %v", err))
					}
				}()
			}

			// Step physics before Lua update
			e.Physics.Step()

			// Update network frame (for multiplayer sync)
			e.Network.UpdateFrame(dt)

			// Use state machine if it has active states, otherwise fall back to direct Lua calls
			if e.GSM != nil {
				_, hasActiveState := e.GSM.GetActiveState()
				if hasActiveState {
					// Handle input
					e.GSM.HandleInput()

					// Update and draw using state machine
					e.GSM.Update(dtSec)
					e.GSM.Draw()
				} else {
					// Fallback: direct Lua calls for games not using state machine
					_ = e.VM.CallUpdate(dtSec)
					_ = e.VM.CallDraw()
				}
			} else {
				// No state machine at all - use direct Lua calls
				_ = e.VM.CallUpdate(dtSec)
				_ = e.VM.CallDraw()
			}

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
		luabind.RegisterWithDev(e.VM.L, e.Ren, func(i int) (c [4]uint8) {
			col := e.Pal.Color(i)
			c[0] = col.R
			c[1] = col.G
			c[2] = col.B
			c[3] = col.A
			return
		}, e.Pal.Set, e.sfxMap, e.musicMap, e.spritesMap, e.Physics, devAdapter, e.Network)
	} else {
		luabind.Register(e.VM.L, e.Ren, func(i int) (c [4]uint8) {
			col := e.Pal.Color(i)
			c[0] = col.R
			c[1] = col.G
			c[2] = col.B
			c[3] = col.A
			return
		}, e.Pal.Set, e.sfxMap, e.musicMap, e.spritesMap, e.Physics, e.Network)
	}

	// Register state machine (needed for game.* API)
	luabind.RegisterStateMachine(e.VM.L, e.GSM)
}

// RunFrames advances N frames headlessly.
func (e *Engine) RunFrames(n int) {
	for i := 0; i < n; i++ {
		e.Run.Step()
	}
}

// LoadCartFromReader loads a .rfs from an io.ReaderAt.
func (e *Engine) LoadCartFromReader(r io.ReaderAt, size int64) error {
	result, err := cartio.Read(r, size)
	if err != nil {
		return err
	}

	// Set palette from manifest if specified
	if result.Manifest.Palette != "" {
		e.Pal.Set(result.Manifest.Palette)
	}

	src, ok := result.Files["assets/"+result.Manifest.Entry]
	if !ok {
		return os.ErrNotExist
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

	return e.LoadLuaSource(string(src))
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
