package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
	"github.com/AndrewDonelson/retroforge-engine/internal/lua"
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

	// Read manifest.json
	manifestPath := filepath.Join(cartPath, "manifest.json")
	mfBytes, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to read manifest.json: %w", err)
	}

	var m cartio.Manifest
	if err := json.Unmarshal(mfBytes, &m); err != nil {
		return fmt.Errorf("failed to parse manifest.json: %w", err)
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
		json.Unmarshal(b, &e.spritesMap)
	}

	start := time.Now()
	err = e.LoadLuaSource(string(src))
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

// ReloadCart reloads the cart (development mode only)
func (e *Engine) ReloadCart() error {
	if e.devMode == nil || !e.devMode.IsEnabled() {
		return fmt.Errorf("reload only available in development mode")
	}

	cartPath := e.devMode.cartPath
	if cartPath == "" {
		return fmt.Errorf("no cart path set")
	}

	e.devMode.AddDebugLog("Reloading cart...")

	// Close current VM and create new one
	e.VM.Close()
	e.VM = lua.New()

	// Re-read manifest
	manifestPath := filepath.Join(cartPath, "manifest.json")
	mfBytes, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to read manifest.json: %w", err)
	}

	var m cartio.Manifest
	if err := json.Unmarshal(mfBytes, &m); err != nil {
		return fmt.Errorf("failed to parse manifest.json: %w", err)
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
		json.Unmarshal(b, &e.spritesMap)
	}

	start := time.Now()
	err = e.LoadLuaSource(string(src))
	loadTime := time.Since(start)

	e.devMode.stats.LoadTime = loadTime
	if err != nil {
		e.devMode.AddDebugLog(fmt.Sprintf("Reload error: %v", err))
		return err
	}

	e.devMode.AddDebugLog(fmt.Sprintf("Reloaded successfully (took %v)", loadTime))
	return nil
}
