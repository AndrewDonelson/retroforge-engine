package engine

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewDevMode(t *testing.T) {
	dm := NewDevMode()
	if dm == nil {
		t.Fatal("NewDevMode returned nil")
	}
	if dm.IsEnabled() {
		t.Error("New dev mode should not be enabled")
	}
}

func TestDevModeEnableDisable(t *testing.T) {
	dm := NewDevMode()

	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	assetsDir := filepath.Join(tmpDir, "assets")
	os.MkdirAll(assetsDir, 0755)

	// Create a dummy manifest.json
	manifestPath := filepath.Join(tmpDir, "manifest.json")
	os.WriteFile(manifestPath, []byte(`{"title":"test"}`), 0644)

	// Test enabling
	err := dm.Enable(tmpDir)
	if err != nil {
		t.Fatalf("Enable failed: %v", err)
	}

	if !dm.IsEnabled() {
		t.Error("DevMode should be enabled after Enable()")
	}

	// Test enabling again (should be idempotent)
	err = dm.Enable(tmpDir)
	if err != nil {
		t.Fatalf("Enable failed on second call: %v", err)
	}

	// Test disabling
	dm.Disable()
	if dm.IsEnabled() {
		t.Error("DevMode should be disabled after Disable()")
	}

	// Test disabling again (should be safe)
	dm.Disable()
}

func TestDevModeDebugLogs(t *testing.T) {
	dm := NewDevMode()

	// Test adding debug logs
	dm.AddDebugLog("test message 1")
	dm.AddDebugLog("test message 2")

	logs := dm.GetDebugLogs()
	if len(logs) != 2 {
		t.Errorf("Expected 2 debug logs, got %d", len(logs))
	}

	// Test log rotation (max 100 logs)
	for i := 0; i < 150; i++ {
		dm.AddDebugLog(string(rune(i)))
	}
	logs = dm.GetDebugLogs()
	if len(logs) > 100 {
		t.Errorf("Debug logs should be capped at 100, got %d", len(logs))
	}
}

func TestDevModeStats(t *testing.T) {
	dm := NewDevMode()

	// Test initial stats
	stats := dm.GetStats()
	if stats.FPS != 0 {
		t.Error("Initial FPS should be 0")
	}

	// Test updating stats
	dm.UpdateStats(60.0, 1000, 1024)
	stats = dm.GetStats()
	if stats.FPS != 60.0 {
		t.Errorf("Stats after UpdateStats: FPS = %f, expected 60.0", stats.FPS)
	}
	if stats.FrameCount != 1000 {
		t.Errorf("Stats after UpdateStats: FrameCount = %d, expected 1000", stats.FrameCount)
	}
	if stats.LuaMemory != 1024 {
		t.Errorf("Stats after UpdateStats: LuaMemory = %d, expected 1024", stats.LuaMemory)
	}
}

func TestDevModeCheckForReload(t *testing.T) {
	dm := NewDevMode()

	// Should return false when not enabled
	if dm.CheckForReload() {
		t.Error("CheckForReload should return false when disabled")
	}

	// Create temp directory
	tmpDir := t.TempDir()
	assetsDir := filepath.Join(tmpDir, "assets")
	os.MkdirAll(assetsDir, 0755)
	manifestPath := filepath.Join(tmpDir, "manifest.json")
	os.WriteFile(manifestPath, []byte(`{"title":"test"}`), 0644)

	// Enable dev mode
	err := dm.Enable(tmpDir)
	if err != nil {
		t.Fatalf("Enable failed: %v", err)
	}

	// Initially should not reload (no file changes yet)
	if dm.CheckForReload() {
		t.Error("CheckForReload should return false initially")
	}

	// Note: Testing actual file watching is difficult without real file system events
	// This is a smoke test that the function doesn't crash
	dm.CheckForReload()

	dm.Disable()
}

func TestEngineLoadCartFolder(t *testing.T) {
	e := New(60)
	defer e.Close()

	// Create a valid cart directory structure
	tmpDir := t.TempDir()
	assetsDir := filepath.Join(tmpDir, "assets")
	os.MkdirAll(assetsDir, 0755)

	// Create manifest.json
	manifestPath := filepath.Join(tmpDir, "manifest.json")
	os.WriteFile(manifestPath, []byte(`{
		"title": "Test Cart",
		"author": "Test",
		"entry": "main.lua"
	}`), 0644)

	// Create main.lua
	luaPath := filepath.Join(assetsDir, "main.lua")
	os.WriteFile(luaPath, []byte(`
		function _INIT()
		end
		function _UPDATE(dt)
		end
		function _DRAW()
			rf.clear_i(0)
		end
	`), 0644)

	// Create empty sfx.json
	sfxPath := filepath.Join(assetsDir, "sfx.json")
	os.WriteFile(sfxPath, []byte(`{}`), 0644)

	// Create empty music.json
	musicPath := filepath.Join(assetsDir, "music.json")
	os.WriteFile(musicPath, []byte(`{}`), 0644)

	// Create empty sprites.json
	spritesPath := filepath.Join(assetsDir, "sprites.json")
	os.WriteFile(spritesPath, []byte(`{}`), 0644)

	// Test loading
	err := e.LoadCartFolder(tmpDir)
	if err != nil {
		t.Fatalf("LoadCartFolder failed: %v", err)
	}

	// Should enable dev mode
	if e.devMode == nil || !e.devMode.IsEnabled() {
		t.Error("LoadCartFolder should enable dev mode")
	}
}

func TestEngineReloadCart(t *testing.T) {
	e := New(60)
	defer e.Close()

	// Create cart directory
	tmpDir := t.TempDir()
	assetsDir := filepath.Join(tmpDir, "assets")
	os.MkdirAll(assetsDir, 0755)

	manifestPath := filepath.Join(tmpDir, "manifest.json")
	os.WriteFile(manifestPath, []byte(`{
		"title": "Test",
		"entry": "main.lua"
	}`), 0644)

	luaPath := filepath.Join(assetsDir, "main.lua")
	os.WriteFile(luaPath, []byte(`function _INIT() end`), 0644)

	// Load cart first
	err := e.LoadCartFolder(tmpDir)
	if err != nil {
		t.Fatalf("LoadCartFolder failed: %v", err)
	}

	// Test reload
	err = e.ReloadCart()
	if err != nil {
		t.Fatalf("ReloadCart failed: %v", err)
	}

	// Modify file and reload again
	os.WriteFile(luaPath, []byte(`function _INIT() x=1 end`), 0644)
	time.Sleep(100 * time.Millisecond) // Small delay for file system
	err = e.ReloadCart()
	if err != nil {
		t.Fatalf("ReloadCart after file change failed: %v", err)
	}
}

func TestDevModeEdgeCases(t *testing.T) {
	dm := NewDevMode()

	// Test Enable with invalid path
	err := dm.Enable("/nonexistent/path")
	if err == nil {
		t.Error("Enable with invalid path should fail")
	}

	// Test Enable with path that has no manifest
	tmpDir := t.TempDir()
	err = dm.Enable(tmpDir)
	if err == nil {
		t.Error("Enable without manifest.json should fail")
	}

	// Test stats update with various values
	dm.UpdateStats(0.0, 0, 0)
	dm.UpdateStats(999.9, 999999, 9999999)

	// Test GetStats returns valid stats
	stats := dm.GetStats()
	if stats.FPS != 999.9 {
		t.Errorf("GetStats should reflect updated stats, FPS = %f", stats.FPS)
	}

	// Test Disable when not enabled (should be safe)
	dm.Disable()
}
