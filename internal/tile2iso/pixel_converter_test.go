package tile2iso

import (
	"image"
	"image/color"
	"testing"
)

func TestPixelDataToImage(t *testing.T) {
	palette := generateTestPalette()

	tests := []struct {
		name      string
		pixels    [][]int
		wantError bool
	}{
		{
			name: "valid 2x2 sprite",
			pixels: [][]int{
				{0, 1},
				{2, 3},
			},
			wantError: false,
		},
		{
			name:      "empty pixels",
			pixels:    [][]int{},
			wantError: true,
		},
		{
			name: "transparent pixels",
			pixels: [][]int{
				{-1, 0},
				{0, -1},
			},
			wantError: false,
		},
		{
			name: "invalid palette index",
			pixels: [][]int{
				{0, 50}, // Index 50 is out of range
				{0, 0},
			},
			wantError: false, // Should handle gracefully
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img, err := PixelDataToImage(tt.pixels, palette)
			if (err != nil) != tt.wantError {
				t.Errorf("PixelDataToImage() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && img == nil {
				t.Error("PixelDataToImage() returned nil image")
			}
		})
	}
}

func TestImageToPixelData(t *testing.T) {
	palette := generateTestPalette()

	// Create a simple test image
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255})   // Red
	img.Set(1, 0, color.RGBA{R: 0, G: 255, B: 0, A: 255})   // Green
	img.Set(0, 1, color.RGBA{R: 0, G: 0, B: 255, A: 255})   // Blue
	img.Set(1, 1, color.RGBA{R: 0, G: 0, B: 0, A: 0})       // Transparent

	pixels, paletteOut, err := ImageToPixelData(img, palette)
	if err != nil {
		t.Fatalf("ImageToPixelData() error = %v", err)
	}

	if len(pixels) != 2 {
		t.Errorf("expected height 2, got %d", len(pixels))
	}
	if len(pixels[0]) != 2 {
		t.Errorf("expected width 2, got %d", len(pixels[0]))
	}

	// Check transparent pixel
	if pixels[1][1] != -1 {
		t.Errorf("expected transparent pixel (-1), got %d", pixels[1][1])
	}

	// Check palette output
	if len(paletteOut) != len(palette) {
		t.Errorf("expected palette length %d, got %d", len(palette), len(paletteOut))
	}
}

func TestParseHexColor(t *testing.T) {
	tests := []struct {
		name string
		hex  string
		want color.RGBA
	}{
		{
			name: "valid hex color",
			hex:  "#ff0000",
			want: color.RGBA{R: 255, G: 0, B: 0, A: 255},
		},
		{
			name: "valid hex color 2",
			hex:  "#00ff00",
			want: color.RGBA{R: 0, G: 255, B: 0, A: 255},
		},
		{
			name: "invalid format - too short",
			hex:  "#ff00",
			want: color.RGBA{A: 255},
		},
		{
			name: "invalid format - no hash",
			hex:  "ff0000",
			want: color.RGBA{A: 255},
		},
		{
			name: "black",
			hex:  "#000000",
			want: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		},
		{
			name: "white",
			hex:  "#ffffff",
			want: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, g, b, a := parseHexColor(tt.hex)
			got := color.RGBA{R: r, G: g, B: b, A: a}
			if got != tt.want {
				t.Errorf("parseHexColor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColorDistance(t *testing.T) {
	c1 := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	c2 := color.RGBA{R: 0, G: 255, B: 0, A: 255}

	dist := colorDistance(c1, c2)
	if dist <= 0 {
		t.Errorf("colorDistance() should be positive, got %f", dist)
	}

	// Same color should have distance 0
	dist2 := colorDistance(c1, c1)
	if dist2 != 0 {
		t.Errorf("colorDistance(same, same) should be 0, got %f", dist2)
	}
}

// generateTestPalette creates a simple test palette with 50 colors
func generateTestPalette() []string {
	palette := make([]string, 50)
	for i := 0; i < 50; i++ {
		// Generate simple colors
		r := uint8(i * 5 % 256)
		g := uint8((i * 7) % 256)
		b := uint8((i * 11) % 256)
		palette[i] = colorToHex(r, g, b)
	}
	return palette
}

func colorToHex(r, g, b uint8) string {
	return "#" + byteToHex(r) + byteToHex(g) + byteToHex(b)
}

func byteToHex(b uint8) string {
	hex := "0123456789abcdef"
	return string(hex[b>>4]) + string(hex[b&0x0f])
}

