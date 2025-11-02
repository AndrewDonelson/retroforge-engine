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

	// Create palette with red at index 2
	palette := &Palette{Colors: make([]string, 50)}
	palette.Colors[0] = "#000000"
	palette.Colors[1] = "#ffffff"
	palette.Colors[2] = "#ff0000"
	for i := 3; i < 50; i++ {
		palette.Colors[i] = "#808080"
	}

	opts := DefaultMapPaletteOptions()
	indices, err := MapPalette(img, palette, opts)
	if err != nil {
		t.Fatalf("MapPalette() error = %v", err)
	}

	// All pixels should map to index 2 (red)
	for y := range indices {
		for x := range indices[y] {
			if indices[y][x] != 2 {
				t.Errorf("MapPalette() [%d][%d] = %v, want 2", y, x, indices[y][x])
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

	palette := &Palette{Colors: make([]string, 50)}
	for i := 0; i < 50; i++ {
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
	palette := &Palette{Colors: make([]string, 49)} // Wrong size

	opts := DefaultMapPaletteOptions()
	_, err := MapPalette(img, palette, opts)
	if err == nil {
		t.Error("MapPalette() should return error for invalid palette size")
	}
}

