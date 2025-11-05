package tile2iso

import (
	"image"
	"image/color"
	"testing"
)

func TestTransformToIsometric(t *testing.T) {
	// Create a simple test image
	testImg := createTestImage(16, 16, color.RGBA{R: 255, G: 0, B: 0, A: 255})

	result, err := transformToIsometric(testImg, 64, 32)
	if err != nil {
		t.Fatalf("transformToIsometric() error = %v", err)
	}

	if result == nil {
		t.Fatal("transformToIsometric() returned nil image")
	}

	bounds := result.Bounds()
	if bounds.Dx() != 64 {
		t.Errorf("expected width 64, got %d", bounds.Dx())
	}
	if bounds.Dy() != 32 {
		t.Errorf("expected height 32, got %d", bounds.Dy())
	}

	// Test nil image
	_, err = transformToIsometric(nil, 64, 32)
	if err == nil {
		t.Error("expected error for nil image")
	}
}

func TestBilinearSample(t *testing.T) {
	// Create a test image with known colors
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255})   // Red
	img.Set(1, 0, color.RGBA{R: 0, G: 255, B: 0, A: 255})   // Green
	img.Set(0, 1, color.RGBA{R: 0, G: 0, B: 255, A: 255})   // Blue
	img.Set(1, 1, color.RGBA{R: 255, G: 255, B: 255, A: 255}) // White

	bounds := img.Bounds()

	// Test exact pixel
	c := bilinearSample(img, 0.0, 0.0, bounds)
	rgba := c.(color.RGBA)
	if rgba.R != 255 {
		t.Errorf("expected red pixel, got R=%d", rgba.R)
	}

	// Test interpolation (middle of four pixels)
	c = bilinearSample(img, 0.5, 0.5, bounds)
	rgba = c.(color.RGBA)
	// Should be average of all four colors
	if rgba.R == 0 && rgba.G == 0 && rgba.B == 0 {
		t.Error("interpolated color should not be black")
	}

	// Test out of bounds (should clamp)
	c = bilinearSample(img, 10.0, 10.0, bounds)
	if c == nil {
		t.Error("bilinearSample should handle out of bounds")
	}
}

func TestToRGBA(t *testing.T) {
	tests := []struct {
		name string
		input color.Color
	}{
		{"RGBA color", color.RGBA{R: 255, G: 0, B: 0, A: 255}},
		{"NRGBA color", color.NRGBA{R: 255, G: 0, B: 0, A: 255}},
		{"Gray color", color.Gray{Y: 128}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rgba := toRGBA(tt.input)
			if rgba.A == 0 && rgba.R == 0 && rgba.G == 0 && rgba.B == 0 {
				t.Error("toRGBA() returned zero color")
			}
		})
	}
}

