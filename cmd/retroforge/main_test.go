//go:build !js && !wasm

package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPackDir(t *testing.T) {
	// Create a temporary test directory
	tmpDir := t.TempDir()

	// Create manifest.json
	manifest := `{
  "title": "Test Game",
  "author": "Test",
  "description": "Test game",
  "entry": "main.lua"
}`
	if err := os.WriteFile(filepath.Join(tmpDir, "manifest.json"), []byte(manifest), 0644); err != nil {
		t.Fatalf("failed to create manifest: %v", err)
	}

	// Create assets directory and main.lua
	assetsDir := filepath.Join(tmpDir, "assets")
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		t.Fatalf("failed to create assets dir: %v", err)
	}

	luaCode := `function _UPDATE(dt)
  -- test
end`
	if err := os.WriteFile(filepath.Join(assetsDir, "main.lua"), []byte(luaCode), 0644); err != nil {
		t.Fatalf("failed to create main.lua: %v", err)
	}

	// Test packDir
	outFile := filepath.Join(tmpDir, "test.rf")
	if err := packDir(tmpDir, outFile); err != nil {
		t.Fatalf("packDir failed: %v", err)
	}

	// Verify output file exists
	if _, err := os.Stat(outFile); os.IsNotExist(err) {
		t.Fatalf("packDir did not create output file")
	}
}

func TestSavePNG(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "test.png")

	// Create test image data
	w, h := 10, 10
	rgba := make([]uint8, w*h*4)
	for i := 0; i < len(rgba); i += 4 {
		rgba[i+0] = 255 // R
		rgba[i+1] = 128 // G
		rgba[i+2] = 64  // B
		rgba[i+3] = 255 // A
	}

	// Test savePNG
	if err := savePNG(tmpFile, w, h, rgba); err != nil {
		t.Fatalf("savePNG failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Fatalf("savePNG did not create file")
	}
}

func TestPackDirEdgeCases(t *testing.T) {
	// Test with empty directory
	tmpDir := t.TempDir()
	outFile := filepath.Join(tmpDir, "empty.rf")
	err := packDir(tmpDir, outFile)
	if err == nil {
		// Might succeed but create empty/valid cart
		t.Logf("packDir with empty dir succeeded")
	}

	// Test with missing manifest.json
	tmpDir2 := t.TempDir()
	assetsDir := filepath.Join(tmpDir2, "assets")
	os.MkdirAll(assetsDir, 0755)
	os.WriteFile(filepath.Join(assetsDir, "main.lua"), []byte("-- test"), 0644)
	outFile2 := filepath.Join(tmpDir2, "no-manifest.rf")
	err = packDir(tmpDir2, outFile2)
	if err == nil {
		t.Logf("packDir without manifest succeeded (might be allowed)")
	}

	// Test with invalid manifest JSON
	tmpDir3 := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir3, "manifest.json"), []byte("INVALID JSON"), 0644)
	outFile3 := filepath.Join(tmpDir3, "bad-manifest.rf")
	err = packDir(tmpDir3, outFile3)
	if err == nil {
		t.Logf("packDir with invalid JSON succeeded")
	}

	// Test with very long path
	longPath := tmpDir
	for i := 0; i < 10; i++ {
		longPath = filepath.Join(longPath, "verylongdirname"+string(make([]byte, 100)))
	}
	os.MkdirAll(longPath, 0755)
	outFile4 := filepath.Join(tmpDir, "longpath.rf")
	// This might fail due to path length limits
	err = packDir(longPath, outFile4)
	if err != nil {
		t.Logf("packDir with long path failed (expected): %v", err)
	}
}

func TestSavePNGEdgeCases(t *testing.T) {
	// Test with zero dimensions
	tmpFile := filepath.Join(t.TempDir(), "zero.png")
	rgba := make([]uint8, 0)
	err := savePNG(tmpFile, 0, 0, rgba)
	if err != nil {
		t.Logf("savePNG with zero dimensions failed (might be expected): %v", err)
	}

	// Test with mismatched size
	tmpFile2 := filepath.Join(t.TempDir(), "mismatch.png")
	rgba2 := make([]uint8, 10) // too small for 10x10
	err = savePNG(tmpFile2, 10, 10, rgba2)
	// Note: PNG encoder might pad or handle gracefully, so we just check it doesn't crash
	if err != nil {
		t.Logf("savePNG with mismatched size failed as expected: %v", err)
	}

	// Test with invalid path (parent doesn't exist)
	invalidFile := filepath.Join("/nonexistent/path", "test.png")
	err = savePNG(invalidFile, 10, 10, make([]uint8, 10*10*4))
	if err == nil {
		t.Fatalf("savePNG with invalid path should fail")
	}

	// Test with very large image
	tmpFile3 := filepath.Join(t.TempDir(), "large.png")
	w, h := 10000, 10000
	rgba3 := make([]uint8, w*h*4)
	err = savePNG(tmpFile3, w, h, rgba3)
	if err != nil {
		t.Logf("savePNG with large image failed (might be memory limit): %v", err)
	}

	// Test with 1x1 image
	tmpFile4 := filepath.Join(t.TempDir(), "tiny.png")
	rgba4 := []uint8{255, 255, 255, 255}
	err = savePNG(tmpFile4, 1, 1, rgba4)
	if err != nil {
		t.Fatalf("savePNG with 1x1 should work: %v", err)
	}
}
