package tile2iso

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
)

func TestLoadSpritesFromFile(t *testing.T) {
	// Create a temporary test sprites.json file
	tmpDir := t.TempDir()
	spritesPath := filepath.Join(tmpDir, "sprites.json")

	// Create a valid sprites.json
	spritesJSON := `{
  "test_sprite": {
    "width": 16,
    "height": 16,
    "type": "static",
    "pixels": [
      [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15],
      [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15],
      [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15],
      [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15],
      [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15],
      [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15],
      [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15],
      [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15],
      [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15],
      [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15],
      [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15],
      [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15],
      [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15],
      [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15],
      [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15],
      [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15]
    ],
    "isUI": false
  }
}`

	if err := os.WriteFile(spritesPath, []byte(spritesJSON), 0644); err != nil {
		t.Fatalf("failed to write test sprites.json: %v", err)
	}

	// Test loading
	spriteMap, err := LoadSpritesFromFile(spritesPath)
	if err != nil {
		t.Fatalf("LoadSpritesFromFile() error = %v", err)
	}

	if len(spriteMap) != 1 {
		t.Errorf("expected 1 sprite, got %d", len(spriteMap))
	}

	sprite, exists := spriteMap["test_sprite"]
	if !exists {
		t.Fatal("sprite 'test_sprite' not found")
	}

	if sprite.Width != 16 || sprite.Height != 16 {
		t.Errorf("expected 16x16 sprite, got %dx%d", sprite.Width, sprite.Height)
	}
}

func TestGetSpritePixels(t *testing.T) {
	// Test static sprite
	staticSprite := &cartio.SpriteData{
		Width:  2,
		Height: 2,
		Type:   cartio.SpriteTypeStatic,
		Pixels: [][]int{
			{0, 1},
			{2, 3},
		},
	}

	pixels, err := GetSpritePixels(staticSprite, "")
	if err != nil {
		t.Fatalf("GetSpritePixels() error = %v", err)
	}

	if len(pixels) != 2 || len(pixels[0]) != 2 {
		t.Errorf("expected 2x2 pixels, got %dx%d", len(pixels), len(pixels[0]))
	}

	// Test frames sprite
	framesSprite := &cartio.SpriteData{
		Width:  2,
		Height: 2,
		Type:   cartio.SpriteTypeFrames,
		Frames: []cartio.SpriteFrame{
			{
				Name: "frame1",
				Pixels: [][]int{
					{0, 1},
					{2, 3},
				},
			},
			{
				Name: "frame2",
				Pixels: [][]int{
					{4, 5},
					{6, 7},
				},
			},
		},
	}

	pixels, err = GetSpritePixels(framesSprite, "frame1")
	if err != nil {
		t.Fatalf("GetSpritePixels() error = %v", err)
	}

	if pixels[0][0] != 0 {
		t.Errorf("expected pixel value 0, got %d", pixels[0][0])
	}

	// Test invalid frame name
	_, err = GetSpritePixels(framesSprite, "invalid_frame")
	if err == nil {
		t.Error("expected error for invalid frame name")
	}

	// Test nil sprite
	_, err = GetSpritePixels(nil, "")
	if err == nil {
		t.Error("expected error for nil sprite")
	}
}

func TestGetSpriteFromMap(t *testing.T) {
	spriteMap := cartio.SpriteMap{
		"test": {
			Width:  2,
			Height: 2,
			Type:   cartio.SpriteTypeStatic,
			Pixels: [][]int{
				{0, 1},
				{2, 3},
			},
		},
	}

	sprite, err := GetSpriteFromMap(spriteMap, "test")
	if err != nil {
		t.Fatalf("GetSpriteFromMap() error = %v", err)
	}

	if sprite.Width != 2 {
		t.Errorf("expected width 2, got %d", sprite.Width)
	}

	// Test non-existent sprite
	_, err = GetSpriteFromMap(spriteMap, "nonexistent")
	if err == nil {
		t.Error("expected error for non-existent sprite")
	}
}

