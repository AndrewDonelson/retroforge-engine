package imgtool

import (
	"image"
	"image/color"
	"testing"
)


func TestQuantize_SingleColor(t *testing.T) {
	img := createTestImage(16, 16, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	opts := DefaultQuantizeOptions()

	palette, err := Quantize(img, opts)
	if err != nil {
		t.Fatalf("Quantize() error = %v", err)
	}

	if len(palette.Colors) != 50 {
		t.Errorf("Quantize() palette length = %v, want 50", len(palette.Colors))
	}
}

func TestQuantize_WithTransparency(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 16, 16))
	// Fill with transparent pixels
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 0})
		}
	}

	opts := DefaultQuantizeOptions()
	opts.AlphaThreshold = 128

	palette, err := Quantize(img, opts)
	if err != nil {
		t.Fatalf("Quantize() error = %v", err)
	}

	if len(palette.Colors) != 50 {
		t.Errorf("Quantize() palette length = %v, want 50", len(palette.Colors))
	}
}

func TestQuantize_EnforcesBlackWhite(t *testing.T) {
	img := createTestImage(16, 16, color.RGBA{R: 128, G: 128, B: 128, A: 255})
	opts := DefaultQuantizeOptions()
	opts.EnforceBlackWhite = true

	palette, err := Quantize(img, opts)
	if err != nil {
		t.Fatalf("Quantize() error = %v", err)
	}

	if palette.Colors[0] != "#000000" {
		t.Errorf("Quantize() index 0 = %v, want #000000", palette.Colors[0])
	}
	if palette.Colors[1] != "#ffffff" {
		t.Errorf("Quantize() index 1 = %v, want #ffffff", palette.Colors[1])
	}
}

func TestQuantize_Dithering(t *testing.T) {
	img := createTestImage(16, 16, color.RGBA{R: 200, G: 100, B: 50, A: 255})
	
	tests := []struct {
		name string
		dither string
	}{
		{"no dither", "none"},
		{"floyd-steinberg", "floyd-steinberg"},
		{"ordered", "ordered"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := DefaultQuantizeOptions()
			opts.DitherAlgorithm = tt.dither
			
			palette, err := Quantize(img, opts)
			if err != nil {
				t.Fatalf("Quantize() error = %v", err)
			}
			
			if len(palette.Colors) != 50 {
				t.Errorf("Quantize() palette length = %v, want 50", len(palette.Colors))
			}
		})
	}
}

