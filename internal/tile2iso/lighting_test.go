package tile2iso

import (
	"image"
	"image/color"
	"testing"
)

func TestApplyLighting(t *testing.T) {
	// Create a simple test image
	testImg := createTestImage(4, 4, color.RGBA{R: 128, G: 128, B: 128, A: 255})

	tests := []struct {
		name      string
		mode      LightingMode
		side      string
		wantError bool
	}{
		{"normal mode", LightingNormal, "left", false},
		{"basic mode left", LightingBasic, "left", false},
		{"basic mode right", LightingBasic, "right", false},
		{"full mode left", LightingFull, "left", false},
		{"full mode right", LightingFull, "right", false},
		{"gradient mode left", LightingGradient, "left", false},
		{"gradient mode right", LightingGradient, "right", false},
		{"invalid mode", LightingMode("invalid"), "left", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := applyLighting(testImg, tt.mode, tt.side, 4)
			if (err != nil) != tt.wantError {
				t.Errorf("applyLighting() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && result == nil {
				t.Error("applyLighting() returned nil image")
			}
		})
	}

	// Test nil image
	_, err := applyLighting(nil, LightingNormal, "left", 4)
	if err == nil {
		t.Error("expected error for nil image")
	}
}

func TestApplyUniformBrightness(t *testing.T) {
	testImg := createTestImage(2, 2, color.RGBA{R: 100, G: 100, B: 100, A: 255})

	// Test left side (brighter)
	result, err := applyUniformBrightness(testImg, true, 1.2, 0.8)
	if err != nil {
		t.Fatalf("applyUniformBrightness() error = %v", err)
	}

	// Check that pixels are brighter
	rgba := result.At(0, 0).(color.RGBA)
	if rgba.R < 100 {
		t.Errorf("expected brighter pixel, got R=%d", rgba.R)
	}

	// Test right side (darker)
	result, err = applyUniformBrightness(testImg, false, 1.2, 0.8)
	if err != nil {
		t.Fatalf("applyUniformBrightness() error = %v", err)
	}

	rgba = result.At(0, 0).(color.RGBA)
	if rgba.R > 100 {
		t.Errorf("expected darker pixel, got R=%d", rgba.R)
	}
}

func TestApplyGradientBrightness(t *testing.T) {
	testImg := createTestImage(2, 10, color.RGBA{R: 100, G: 100, B: 100, A: 255})

	// Test left side
	result, err := applyGradientBrightness(testImg, true, 10)
	if err != nil {
		t.Fatalf("applyGradientBrightness() error = %v", err)
	}

	// Top should be brighter than bottom
	topPixel := result.At(0, 0).(color.RGBA)
	bottomPixel := result.At(0, 9).(color.RGBA)

	if topPixel.R < bottomPixel.R {
		t.Errorf("expected top pixel (R=%d) to be brighter than bottom (R=%d)", topPixel.R, bottomPixel.R)
	}

	// Test right side
	result, err = applyGradientBrightness(testImg, false, 10)
	if err != nil {
		t.Fatalf("applyGradientBrightness() error = %v", err)
	}

	// Top should be darker than bottom
	topPixel = result.At(0, 0).(color.RGBA)
	bottomPixel = result.At(0, 9).(color.RGBA)

	if topPixel.R > bottomPixel.R {
		t.Errorf("expected top pixel (R=%d) to be darker than bottom (R=%d)", topPixel.R, bottomPixel.R)
	}
}

func TestAdjustPixelBrightness(t *testing.T) {
	tests := []struct {
		name   string
		input  color.RGBA
		factor float64
		want   color.RGBA
	}{
		{
			name:   "normal brightness",
			input:  color.RGBA{R: 100, G: 100, B: 100, A: 255},
			factor: 1.0,
			want:   color.RGBA{R: 100, G: 100, B: 100, A: 255},
		},
		{
			name:   "brighter",
			input:  color.RGBA{R: 100, G: 100, B: 100, A: 255},
			factor: 1.5,
			want:   color.RGBA{R: 150, G: 150, B: 150, A: 255},
		},
		{
			name:   "darker",
			input:  color.RGBA{R: 100, G: 100, B: 100, A: 255},
			factor: 0.5,
			want:   color.RGBA{R: 50, G: 50, B: 50, A: 255},
		},
		{
			name:   "clamp to 255",
			input:  color.RGBA{R: 200, G: 200, B: 200, A: 255},
			factor: 2.0,
			want:   color.RGBA{R: 255, G: 255, B: 255, A: 255},
		},
		{
			name:   "preserve alpha",
			input:  color.RGBA{R: 100, G: 100, B: 100, A: 128},
			factor: 1.5,
			want:   color.RGBA{R: 150, G: 150, B: 150, A: 128},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := adjustPixelBrightness(tt.input, tt.factor)
			if got != tt.want {
				t.Errorf("adjustPixelBrightness() = %v, want %v", got, tt.want)
			}
		})
	}
}

// createTestImage creates a simple test image filled with the given color
func createTestImage(width, height int, c color.RGBA) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, c)
		}
	}
	return img
}

