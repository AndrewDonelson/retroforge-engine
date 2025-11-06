package imgtool

import (
	"image"
	"image/color"
	"testing"
)

func TestMapPalette_ExactMatch(t *testing.T) {
	// Create image with red
	img := image.NewRGBA(image.Rect(0, 0, 16, 16))
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}

	// Create palette with red at index 0 (game palette, 48 colors)
	palette := &Palette{Colors: make([]string, 48)}
	palette.Colors[0] = "#ff0000" // Red
	for i := 1; i < 48; i++ {
		palette.Colors[i] = "#808080"
	}

	opts := DefaultMapPaletteOptions()
	indices, err := MapPalette(img, palette, opts)
	if err != nil {
		t.Fatalf("MapPalette() error = %v", err)
	}

	// All pixels should map to index 0 (red in game palette)
	// Note: In full 64-color system, this would be index 16 (game palette starts at 16)
	// But MapPalette just maps to the provided palette indices (0-47 for game palette)
	for y := range indices {
		for x := range indices[y] {
			if indices[y][x] != 0 {
				t.Errorf("MapPalette() [%d][%d] = %v, want 0", y, x, indices[y][x])
			}
		}
	}
}

func TestMapPalette_TransparentPixels(t *testing.T) {
	// Create image with transparent pixels
	img := image.NewRGBA(image.Rect(0, 0, 16, 16))
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 0}) // Fully transparent
		}
	}

	palette := &Palette{Colors: make([]string, 48)}
	for i := 0; i < 48; i++ {
		palette.Colors[i] = "#000000"
	}

	opts := DefaultMapPaletteOptions()
	opts.AlphaThreshold = 128
	opts.TransparentIndex = -1

	indices, err := MapPalette(img, palette, opts)
	if err != nil {
		t.Fatalf("MapPalette() error = %v", err)
	}

	// All pixels should be transparent (-1)
	for y := range indices {
		for x := range indices[y] {
			if indices[y][x] != -1 {
				t.Errorf("MapPalette() [%d][%d] = %v, want -1", y, x, indices[y][x])
			}
		}
	}
}

func TestMapPalette_InvalidPalette(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 16, 16))
	palette := &Palette{Colors: make([]string, 47)} // Wrong size (should be 48)

	opts := DefaultMapPaletteOptions()
	_, err := MapPalette(img, palette, opts)
	if err == nil {
		t.Error("MapPalette() should return error for invalid palette size")
	}
}

