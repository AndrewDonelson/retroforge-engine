package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
)

func TestConvertCmd_Validation(t *testing.T) {
	cmd := ConvertCmd()

	// Test missing required flags
	tests := []struct {
		name string
		args []string
		wantError bool
	}{
		{"missing sprites", []string{}, true},
		{"missing palette", []string{"--sprites", "test.json"}, true},
		{"missing top", []string{"--sprites", "test.json", "--palette", "palette.json"}, true},
		{"missing left", []string{"--sprites", "test.json", "--palette", "palette.json", "--top", "top"}, true},
		{"missing right", []string{"--sprites", "test.json", "--palette", "palette.json", "--top", "top", "--left", "left"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd.SetArgs(tt.args)
			err := cmd.Execute()
			if (err != nil) != tt.wantError {
				t.Errorf("expected error %v, got %v", tt.wantError, err)
			}
		})
	}
}

func TestConvertCmd_Integration(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test sprites.json
	spritesPath := filepath.Join(tmpDir, "sprites.json")
	spriteMap := cartio.SpriteMap{
		"top": {
			Width:  16,
			Height: 16,
			Type:   cartio.SpriteTypeStatic,
			Pixels: generateTestPixels(16, 16),
			IsUI:   false,
		},
		"left": {
			Width:  16,
			Height: 16,
			Type:   cartio.SpriteTypeStatic,
			Pixels: generateTestPixels(16, 16),
			IsUI:   false,
		},
		"right": {
			Width:  16,
			Height: 16,
			Type:   cartio.SpriteTypeStatic,
			Pixels: generateTestPixels(16, 16),
			IsUI:   false,
		},
	}

	spritesData, _ := json.MarshalIndent(spriteMap, "", "  ")
	os.WriteFile(spritesPath, spritesData, 0644)

	// Create test palette.json
	palettePath := filepath.Join(tmpDir, "palette.json")
	palette := map[string]interface{}{
		"colors": generateTestPalette(),
	}
	paletteData, _ := json.MarshalIndent(palette, "", "  ")
	os.WriteFile(palettePath, paletteData, 0644)

	// Test convert command
	cmd := ConvertCmd()
	outputPath := filepath.Join(tmpDir, "output.json")
	cmd.SetArgs([]string{
		"--sprites", spritesPath,
		"--palette", palettePath,
		"--top", "top",
		"--left", "left",
		"--right", "right",
		"--output", outputPath,
		"--name", "test_tile",
	})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("convert command failed: %v", err)
	}

	// Verify output file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("output file was not created")
	}

	// Verify output contains the sprite
	outputData, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	var outputMap cartio.SpriteMap
	if err := json.Unmarshal(outputData, &outputMap); err != nil {
		t.Fatalf("failed to parse output JSON: %v", err)
	}

	if _, exists := outputMap["test_tile"]; !exists {
		t.Error("output sprite 'test_tile' not found")
	}

	tile := outputMap["test_tile"]
	if tile.Width == 0 || tile.Height == 0 {
		t.Error("output tile has invalid dimensions")
	}
}

func TestConvertCmd_UpdateSprites(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test sprites.json
	spritesPath := filepath.Join(tmpDir, "sprites.json")
	spriteMap := cartio.SpriteMap{
		"top": {
			Width:  16,
			Height: 16,
			Type:   cartio.SpriteTypeStatic,
			Pixels: generateTestPixels(16, 16),
			IsUI:   false,
		},
		"left": {
			Width:  16,
			Height: 16,
			Type:   cartio.SpriteTypeStatic,
			Pixels: generateTestPixels(16, 16),
			IsUI:   false,
		},
		"right": {
			Width:  16,
			Height: 16,
			Type:   cartio.SpriteTypeStatic,
			Pixels: generateTestPixels(16, 16),
			IsUI:   false,
		},
	}

	spritesData, _ := json.MarshalIndent(spriteMap, "", "  ")
	os.WriteFile(spritesPath, spritesData, 0644)

	// Create test palette.json
	palettePath := filepath.Join(tmpDir, "palette.json")
	palette := map[string]interface{}{
		"colors": generateTestPalette(),
	}
	paletteData, _ := json.MarshalIndent(palette, "", "  ")
	os.WriteFile(palettePath, paletteData, 0644)

	// Test convert with --update flag
	cmd := ConvertCmd()
	cmd.SetArgs([]string{
		"--sprites", spritesPath,
		"--palette", palettePath,
		"--top", "top",
		"--left", "left",
		"--right", "right",
		"--update",
		"--name", "iso_tile",
	})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("convert command failed: %v", err)
	}

	// Verify sprites.json was updated
	updatedData, err := os.ReadFile(spritesPath)
	if err != nil {
		t.Fatalf("failed to read updated sprites.json: %v", err)
	}

	var updatedMap cartio.SpriteMap
	if err := json.Unmarshal(updatedData, &updatedMap); err != nil {
		t.Fatalf("failed to parse updated sprites.json: %v", err)
	}

	if _, exists := updatedMap["iso_tile"]; !exists {
		t.Error("updated sprites.json does not contain 'iso_tile'")
	}

	// Verify original sprites still exist
	if _, exists := updatedMap["top"]; !exists {
		t.Error("original sprite 'top' was removed")
	}
}

// Helper functions
func generateTestPixels(width, height int) [][]int {
	pixels := make([][]int, height)
	for y := 0; y < height; y++ {
		pixels[y] = make([]int, width)
		for x := 0; x < width; x++ {
			pixels[y][x] = (x + y) % 50
		}
	}
	return pixels
}

func generateTestPalette() []string {
	palette := make([]string, 50)
	for i := 0; i < 50; i++ {
		r := uint8(i * 5 % 256)
		g := uint8((i * 7) % 256)
		b := uint8((i * 11) % 256)
		palette[i] = fmt.Sprintf("#%02x%02x%02x", r, g, b)
	}
	return palette
}

