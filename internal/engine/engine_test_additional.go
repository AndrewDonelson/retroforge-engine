package engine

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
)

func TestDevModeAdapter(t *testing.T) {
	e := New(60)
	defer e.Close()

	// Create a dev mode for testing
	tmpDir := t.TempDir()
	assetsDir := filepath.Join(tmpDir, "assets")
	os.MkdirAll(assetsDir, 0755)
	manifestPath := filepath.Join(tmpDir, "manifest.json")
	os.WriteFile(manifestPath, []byte(`{"title":"test","entry":"main.lua"}`), 0644)

	err := e.LoadCartFolder(tmpDir)
	if err == nil {
		// Dev mode should now exist
		if e.devMode != nil {
			adapter := &devModeAdapter{devMode: e.devMode}

			// These should not crash
			enabled := adapter.IsEnabled()
			if !enabled {
				t.Error("adapter.IsEnabled should return true when dev mode is enabled")
			}

			adapter.AddDebugLog("test")

			stats := adapter.GetStats()
			if stats == nil {
				t.Error("adapter.GetStats should not return nil")
			}
		}
	}
}

func TestLoadCartFile(t *testing.T) {
	e := New(60)
	defer e.Close()

	// Create a temporary cart file
	tmpDir := t.TempDir()
	cartFile := filepath.Join(tmpDir, "test.rf")

	// Build a minimal cart
	m := cartio.Manifest{Title: "Test", Author: "RF", Entry: "main.lua"}
	lua := `
		function _INIT()
		end
		function _UPDATE(dt)
		end
		function _DRAW()
			rf.clear_i(0)
		end
	`

	var buf bytes.Buffer
	if err := cartio.Write(&buf, m, []cartio.Asset{{Name: "main.lua", Data: []byte(lua)}}, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap)); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	os.WriteFile(cartFile, buf.Bytes(), 0644)

	// Test loading
	err := e.LoadCartFile(cartFile)
	if err != nil {
		t.Fatalf("LoadCartFile failed: %v", err)
	}

	// Test loading non-existent file
	err = e.LoadCartFile("/nonexistent/file.rf")
	if err == nil {
		t.Error("LoadCartFile with non-existent file should fail")
	}

	// Test loading invalid cart
	invalidFile := filepath.Join(tmpDir, "invalid.rf")
	os.WriteFile(invalidFile, []byte("not a valid cart"), 0644)
	err = e.LoadCartFile(invalidFile)
	if err == nil {
		t.Error("LoadCartFile with invalid cart should fail")
	}
}
