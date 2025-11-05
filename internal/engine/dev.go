package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
	"github.com/AndrewDonelson/retroforge-engine/internal/gamestate"
	"github.com/AndrewDonelson/retroforge-engine/internal/lua"
	"github.com/AndrewDonelson/retroforge-engine/internal/luabind"
	"github.com/fsnotify/fsnotify"
)

// DevMode tracks development mode state
type DevMode struct {
	enabled        bool
	cartPath       string
	watcher        *fsnotify.Watcher
	lastReload     time.Time
	reloadCooldown time.Duration
	mu             sync.Mutex
	debugLogs      []string
	debugMaxLogs   int
	stats          DevStats
	isReloading    bool // Flag to indicate reload is in progress
}

// DevStats holds debugging statistics
// Note: This must match luabind.DevStats structure
type DevStats struct {
	FPS         float64
	FrameCount  int64
	LuaMemory   int64
	LoadTime    time.Duration
	LastReload  time.Time
	ReloadCount int
}

// NewDevMode creates a new development mode handler
func NewDevMode() *DevMode {
	return &DevMode{
		enabled:        false,
		reloadCooldown: 500 * time.Millisecond, // Cooldown to avoid rapid reloads
		debugMaxLogs:   100,
		debugLogs:      make([]string, 0, 100),
	}
}

// Enable enables development mode and starts file watching
func (dm *DevMode) Enable(cartPath string) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if dm.enabled {
		return nil // Already enabled
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create file watcher: %w", err)
	}

	// Watch the assets directory
	assetsPath := filepath.Join(cartPath, "assets")
	if err := watcher.Add(assetsPath); err != nil {
		watcher.Close()
		return fmt.Errorf("failed to watch assets directory: %w", err)
	}

	// Watch manifest.json
	manifestPath := filepath.Join(cartPath, "manifest.json")
	if err := watcher.Add(manifestPath); err != nil {
		watcher.Close()
		return fmt.Errorf("failed to watch manifest.json: %w", err)
	}

	dm.enabled = true
	dm.cartPath = cartPath
	dm.watcher = watcher

	return nil
}

// Disable disables development mode and stops file watching
func (dm *DevMode) Disable() {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if !dm.enabled {
		return
	}

	if dm.watcher != nil {
		dm.watcher.Close()
		dm.watcher = nil
	}

	dm.enabled = false
}

// IsEnabled returns whether development mode is active
func (dm *DevMode) IsEnabled() bool {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	return dm.enabled
}

// CheckForReload checks if any files have changed and need reloading
func (dm *DevMode) CheckForReload() bool {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if !dm.enabled || dm.watcher == nil {
		return false
	}

	select {
	case event, ok := <-dm.watcher.Events:
		if !ok {
			return false
		}
		// Only reload on write events, ignore chmod
		if event.Op&fsnotify.Write == fsnotify.Write {
			// Cooldown to avoid rapid reloads
			now := time.Now()
			if now.Sub(dm.lastReload) < dm.reloadCooldown {
				return false
			}
			dm.lastReload = now
			dm.stats.ReloadCount++
			dm.stats.LastReload = now
			return true
		}
	case err := <-dm.watcher.Errors:
		if err != nil {
			dm.AddDebugLog(fmt.Sprintf("File watcher error: %v", err))
		}
		return false
	default:
		return false
	}

	return false
}

// AddDebugLog adds a debug log message
func (dm *DevMode) AddDebugLog(msg string) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	timestamp := time.Now().Format("15:04:05")
	logMsg := fmt.Sprintf("[%s] %s", timestamp, msg)
	dm.debugLogs = append(dm.debugLogs, logMsg)

	// Keep only last N logs
	if len(dm.debugLogs) > dm.debugMaxLogs {
		dm.debugLogs = dm.debugLogs[len(dm.debugLogs)-dm.debugMaxLogs:]
	}
}

// GetDebugLogs returns the debug log messages
func (dm *DevMode) GetDebugLogs() []string {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	logs := make([]string, len(dm.debugLogs))
	copy(logs, dm.debugLogs)
	return logs
}

// GetStats returns current debug statistics
func (dm *DevMode) GetStats() DevStats {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	stats := dm.stats
	return stats
}

// UpdateStats updates debug statistics
func (dm *DevMode) UpdateStats(fps float64, frameCount int64, luaMemory int64) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	dm.stats.FPS = fps
	dm.stats.FrameCount = frameCount
	dm.stats.LuaMemory = luaMemory
}

// LoadCartFolder loads a cart from a directory (development mode only)
func (e *Engine) LoadCartFolder(cartPath string) error {
	// Enable development mode
	if e.devMode == nil {
		e.devMode = NewDevMode()
	}
	if err := e.devMode.Enable(cartPath); err != nil {
		return fmt.Errorf("failed to enable dev mode: %w", err)
	}

	// Update GSM - use debug mode false to show splash screen even in dev mode
	// (User wants to see splash screen in development too)
	// Note: renderer and palette will be set in registerLuaBindings
	e.GSM = gamestate.NewGameStateMachine(false, "RetroForge", "v1.0 Alpha", "RetroForge Team", nil, nil)

	// Read manifest.json
	manifestPath := filepath.Join(cartPath, "manifest.json")
	mfBytes, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to read manifest.json: %w", err)
	}

	// Handle new manifest structure: extract from fullManifest if present
	var rawManifest map[string]interface{}
	if err := json.Unmarshal(mfBytes, &rawManifest); err != nil {
		return fmt.Errorf("failed to parse manifest.json: %w", err)
	}

	var actualManifest map[string]interface{}
	if fullManifest, ok := rawManifest["fullManifest"].(map[string]interface{}); ok {
		actualManifest = fullManifest
	} else {
		actualManifest = rawManifest
	}

	var m cartio.Manifest
	actualManifestBytes, _ := json.Marshal(actualManifest)
	if err := json.Unmarshal(actualManifestBytes, &m); err != nil {
		return fmt.Errorf("failed to parse manifest structure: %w", err)
	}

	// Set palette from manifest if specified
	if m.Palette != "" {
		e.Pal.Set(m.Palette)
	}

	// Load main.lua
	entryPath := filepath.Join(cartPath, "assets", m.Entry)
	src, err := os.ReadFile(entryPath)
	if err != nil {
		return fmt.Errorf("failed to read entry file %s: %w", entryPath, err)
	}

	// Load SFX
	sfxPath := filepath.Join(cartPath, "assets", "sfx.json")
	e.sfxMap = make(cartio.SFXMap)
	if b, err := os.ReadFile(sfxPath); err == nil {
		json.Unmarshal(b, &e.sfxMap)
	}

	// Load Music
	musicPath := filepath.Join(cartPath, "assets", "music.json")
	e.musicMap = make(cartio.MusicMap)
	if b, err := os.ReadFile(musicPath); err == nil {
		json.Unmarshal(b, &e.musicMap)
	}

	// Load Sprites
	spritesPath := filepath.Join(cartPath, "assets", "sprites.json")
	e.spritesMap = make(cartio.SpriteMap)
	if b, err := os.ReadFile(spritesPath); err == nil {
		if err := json.Unmarshal(b, &e.spritesMap); err == nil {
			// Validate and normalize all loaded sprites
			for spriteName, sprite := range e.spritesMap {
				// Normalize sprite data (set defaults, trim whitespace, etc.)
				cartio.NormalizeSpriteData(&sprite)
				
				// Validate complete sprite structure
				if err := cartio.ValidateSpriteData(&sprite, spriteName); err != nil {
					if e.devMode != nil {
						e.devMode.AddDebugLog(fmt.Sprintf("Sprite '%s' validation error: %v", spriteName, err))
					}
					// Remove invalid sprite from map
					delete(e.spritesMap, spriteName)
				} else {
					// Update map with normalized sprite
					e.spritesMap[spriteName] = sprite
				}
			}
		}
	}
	
	// Load .rpi files (Raw Palette Indexed images)
	assetsPathForRPI := filepath.Join(cartPath, "assets")
	if entries, err := os.ReadDir(assetsPathForRPI); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			if len(name) > 4 && name[len(name)-4:] == ".rpi" {
				spriteName := name[:len(name)-4] // Remove .rpi extension
				rpiPath := filepath.Join(assetsPathForRPI, name)
				if data, err := os.ReadFile(rpiPath); err == nil {
					if rpiSprite, err := cartio.LoadRPI(data); err == nil {
						// Normalize RPI sprite (ensure type is set)
						cartio.NormalizeSpriteData(rpiSprite)
						
						// Validate RPI sprite
						if err := cartio.ValidateSpriteData(rpiSprite, spriteName); err == nil {
							e.spritesMap[spriteName] = *rpiSprite
						} else {
							if e.devMode != nil {
								e.devMode.AddDebugLog(fmt.Sprintf("RPI sprite '%s' validation error: %v", spriteName, err))
							}
						}
					}
				}
			}
		}
	}

	// Load Tilesets (scan for *_tiles.json files)
	e.tilesetsMap = make(map[string]cartio.TilesetMap)
	if entries, err := os.ReadDir(assetsPathForRPI); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			// Check for *_tiles.json pattern
			if len(name) > 11 && name[len(name)-11:] == "_tiles.json" {
				tilesetName := name[:len(name)-11] // Remove _tiles.json extension
				tilesetPath := filepath.Join(assetsPathForRPI, name)
				if err := e.loadTileset(tilesetName, tilesetPath); err != nil {
					if e.devMode != nil {
						e.devMode.AddDebugLog(fmt.Sprintf("Failed to load tileset '%s': %v", tilesetName, err))
					}
				}
			}
		}
	}

	// Load Tilemaps (scan for *_map.json files)
	e.tilemapsMap = make(map[string]*cartio.TileMapData)
	if entries, err := os.ReadDir(assetsPathForRPI); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			// Check for *_map.json pattern
			if len(name) > 9 && name[len(name)-9:] == "_map.json" {
				tilemapName := name[:len(name)-9] // Remove _map.json extension
				tilemapPath := filepath.Join(assetsPathForRPI, name)
				if err := e.loadTilemap(tilemapName, tilemapPath); err != nil {
					if e.devMode != nil {
						e.devMode.AddDebugLog(fmt.Sprintf("Failed to load tilemap '%s': %v", tilemapName, err))
					}
				}
			}
		}
	}

	// Register Lua bindings first (creates rf table)
	e.registerLuaBindings()

	// Register module import with filesystem (dev mode) - rf table now exists
	assetsPath := filepath.Join(cartPath, "assets")
	luabind.RegisterModuleImportWithFilesystem(e.VM.L, e.GSM, assetsPath)

	start := time.Now()
	err = e.LoadLuaSource(string(src))
	if err == nil {
		// Start the state machine (shows splash in release, goes to initial state in debug)
		// Pass "menu" as initial state - engine splash will transition to game splash if exists, then menu
		initialState := "menu" // Default initial state
		if startErr := e.GSM.Start(initialState); startErr != nil {
			// If Start fails (e.g., no initial state), continue anyway
			// State machine will be empty but that's ok for now
			if e.devMode != nil {
				e.devMode.AddDebugLog(fmt.Sprintf("State machine start warning: %v", startErr))
			}
		}
	}
	loadTime := time.Since(start)

	if e.devMode != nil {
		e.devMode.stats.LoadTime = loadTime
		if err != nil {
			e.devMode.AddDebugLog(fmt.Sprintf("Load error: %v", err))
		} else {
			e.devMode.AddDebugLog(fmt.Sprintf("Loaded %s (took %v)", m.Entry, loadTime))
		}
	}

	return err
}

// IsReloading returns whether a reload is currently in progress
func (dm *DevMode) IsReloading() bool {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	return dm.isReloading
}

// ReloadCart reloads the cart (development mode only)
func (e *Engine) ReloadCart() error {
	if e.devMode == nil || !e.devMode.IsEnabled() {
		return fmt.Errorf("reload only available in development mode")
	}

	// Set reloading flag
	e.devMode.mu.Lock()
	if e.devMode.isReloading {
		e.devMode.mu.Unlock()
		return fmt.Errorf("reload already in progress")
	}
	e.devMode.isReloading = true
	e.devMode.mu.Unlock()

	// Ensure flag is cleared even if reload fails
	defer func() {
		e.devMode.mu.Lock()
		e.devMode.isReloading = false
		e.devMode.mu.Unlock()
	}()

	cartPath := e.devMode.cartPath
	if cartPath == "" {
		return fmt.Errorf("no cart path set")
	}

	e.devMode.AddDebugLog("Reloading cart...")

	// Clear all module-based states before reloading (preserve built-in states)
	// This prevents old Lua callbacks from referencing closed VM
	if e.GSM != nil {
		e.clearModuleStates()
	}

	// Close current VM and create new one
	e.VM.Close()
	e.VM = lua.New()

	// Re-read manifest (handle new manifest structure)
	manifestPath := filepath.Join(cartPath, "manifest.json")
	mfBytes, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to read manifest.json: %w", err)
	}

	// Handle new manifest structure: extract from fullManifest if present
	var rawManifest map[string]interface{}
	if err := json.Unmarshal(mfBytes, &rawManifest); err != nil {
		return fmt.Errorf("failed to parse manifest.json: %w", err)
	}

	var actualManifest map[string]interface{}
	if fullManifest, ok := rawManifest["fullManifest"].(map[string]interface{}); ok {
		actualManifest = fullManifest
	} else {
		actualManifest = rawManifest
	}

	var m cartio.Manifest
	actualManifestBytes, _ := json.Marshal(actualManifest)
	if err := json.Unmarshal(actualManifestBytes, &m); err != nil {
		return fmt.Errorf("failed to parse manifest structure: %w", err)
	}

	// Set palette from manifest if specified
	if m.Palette != "" {
		e.Pal.Set(m.Palette)
	}

	// Load main.lua
	entryPath := filepath.Join(cartPath, "assets", m.Entry)
	src, err := os.ReadFile(entryPath)
	if err != nil {
		return fmt.Errorf("failed to read entry file %s: %w", entryPath, err)
	}

	// Reload SFX, Music, Sprites
	sfxPath := filepath.Join(cartPath, "assets", "sfx.json")
	e.sfxMap = make(cartio.SFXMap)
	if b, err := os.ReadFile(sfxPath); err == nil {
		json.Unmarshal(b, &e.sfxMap)
	}

	musicPath := filepath.Join(cartPath, "assets", "music.json")
	e.musicMap = make(cartio.MusicMap)
	if b, err := os.ReadFile(musicPath); err == nil {
		json.Unmarshal(b, &e.musicMap)
	}

	spritesPath := filepath.Join(cartPath, "assets", "sprites.json")
	e.spritesMap = make(cartio.SpriteMap)
	if b, err := os.ReadFile(spritesPath); err == nil {
		if err := json.Unmarshal(b, &e.spritesMap); err == nil {
			// Validate and normalize all loaded sprites
			for spriteName, sprite := range e.spritesMap {
				// Normalize sprite data (set defaults, trim whitespace, etc.)
				cartio.NormalizeSpriteData(&sprite)
				
				// Validate complete sprite structure
				if err := cartio.ValidateSpriteData(&sprite, spriteName); err != nil {
					if e.devMode != nil {
						e.devMode.AddDebugLog(fmt.Sprintf("Sprite '%s' validation error: %v", spriteName, err))
					}
					// Remove invalid sprite from map
					delete(e.spritesMap, spriteName)
				} else {
					// Update map with normalized sprite
					e.spritesMap[spriteName] = sprite
				}
			}
		}
	}
	
	// Load .rpi files (Raw Palette Indexed images) - also on reload
	assetsPathForRPI := filepath.Join(cartPath, "assets")
	if entries, err := os.ReadDir(assetsPathForRPI); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			if len(name) > 4 && name[len(name)-4:] == ".rpi" {
				spriteName := name[:len(name)-4] // Remove .rpi extension
				rpiPath := filepath.Join(assetsPathForRPI, name)
				if data, err := os.ReadFile(rpiPath); err == nil {
					if rpiSprite, err := cartio.LoadRPI(data); err == nil {
						// Normalize RPI sprite (ensure type is set)
						cartio.NormalizeSpriteData(rpiSprite)
						
						// Validate RPI sprite
						if err := cartio.ValidateSpriteData(rpiSprite, spriteName); err == nil {
							e.spritesMap[spriteName] = *rpiSprite
						} else {
							if e.devMode != nil {
								e.devMode.AddDebugLog(fmt.Sprintf("RPI sprite '%s' validation error: %v", spriteName, err))
							}
						}
					}
				}
			}
		}
	}

	// Register Lua bindings first (creates rf table)
	e.registerLuaBindings()

	// Register module import with filesystem (dev mode) - rf table now exists
	assetsPath := filepath.Join(cartPath, "assets")
	luabind.RegisterModuleImportWithFilesystem(e.VM.L, e.GSM, assetsPath)

	start := time.Now()
	err = e.LoadLuaSource(string(src))
	if err == nil {
		// Start the state machine after reload
		initialState := "menu" // Default to menu in debug
		if startErr := e.GSM.Start(initialState); startErr != nil {
			if e.devMode != nil {
				e.devMode.AddDebugLog(fmt.Sprintf("State machine start warning: %v", startErr))
			}
		}
	}
	loadTime := time.Since(start)

	e.devMode.stats.LoadTime = loadTime
	if err != nil {
		e.devMode.AddDebugLog(fmt.Sprintf("Reload error: %v", err))
		return err
	}

	e.devMode.AddDebugLog(fmt.Sprintf("Reloaded successfully (took %v)", loadTime))
	return nil
}

// clearModuleStates removes all module-based states (those imported via rf.import)
// but preserves built-in states (splash, credits)
func (e *Engine) clearModuleStates() {
	if e.GSM == nil || e.VM == nil {
		return
	}

	// First, pop all states from stack so we can unregister them
	e.GSM.PopAllStates()

	// Get module loader for the old VM (before it's closed)
	// Use e.VM.L to get the underlying LState
	loader := luabind.GetModuleLoader(e.VM.L)
	if loader != nil {
		// Get all loaded module names from the loader
		loadedModules := loader.GetLoadedModules()
		for _, name := range loadedModules {
			// Unregister each module state
			// Ignore errors - state might already be gone
			_ = e.GSM.UnregisterState(name)
		}
	} else {
		// Fallback: try common module state names
		moduleNames := []string{"menu", "play", "play_state"}
		for _, name := range moduleNames {
			if e.GSM.IsStateRegistered(name) {
				_ = e.GSM.UnregisterState(name)
			}
		}
	}

	// Also unregister the old module loader from the map
	luabind.UnregisterModuleLoader(e.VM.L)
}
