package imgtool

import (
	"image"
	"image/color"
	"testing"
)

// createTestImage creates a simple test image
func createTestImage(width, height int, bgColor color.Color) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, bgColor)
		}
	}
	return img
}

func TestScale_ExactDimensions(t *testing.T) {
	img := createTestImage(32, 32, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	opts := DefaultScaleOptions()
	opts.Width = 16
	opts.Height = 16

	rgbData, err := Scale(img, opts)
	if err != nil {
		t.Fatalf("Scale() error = %v", err)
	}

	if len(rgbData) != 16 {
		t.Errorf("Scale() height = %v, want 16", len(rgbData))
	}
	if len(rgbData[0]) != 16 {
		t.Errorf("Scale() width = %v, want 16", len(rgbData[0]))
	}
}

func TestScale_DivisibleBy2(t *testing.T) {
	img := createTestImage(32, 32, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	opts := DefaultScaleOptions()
	opts.Width = 15  // Odd width
	opts.Height = 17 // Odd height
	opts.EnsureDivisible = true

	rgbData, err := Scale(img, opts)
	if err != nil {
		t.Fatalf("Scale() error = %v", err)
	}

	height := len(rgbData)
	width := len(rgbData[0])

	if width%2 != 0 {
		t.Errorf("Scale() width = %v, should be divisible by 2", width)
	}
	if height%2 != 0 {
		t.Errorf("Scale() height = %v, should be divisible by 2", height)
	}
}

func TestScale_Algorithms(t *testing.T) {
	img := createTestImage(32, 32, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	
	algorithms := []string{"nearest", "bilinear", "bicubic"}
	
	for _, alg := range algorithms {
		t.Run(alg, func(t *testing.T) {
			opts := DefaultScaleOptions()
			opts.Width = 16
			opts.Height = 16
			opts.Algorithm = alg

			rgbData, err := Scale(img, opts)
			if err != nil {
				t.Fatalf("Scale() with %s error = %v", alg, err)
			}

			if len(rgbData) != 16 || len(rgbData[0]) != 16 {
				t.Errorf("Scale() with %s wrong dimensions", alg)
			}
		})
	}
}

