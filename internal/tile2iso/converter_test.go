package tile2iso

import (
	"image/color"
	"testing"

	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
)

func TestCreateIsometricTile(t *testing.T) {
	palette := generateTestPalette()
	converter := NewIsometricConverter(64, 32)

	// Create test sprites
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

	options := DefaultTileOptions()

	result, err := converter.CreateIsometricTile(
		"top", "left", "right",
		"", "", "", // No frame names for static sprites
		palette,
		spriteMap,
		options,
	)

	if err != nil {
		t.Fatalf("CreateIsometricTile() error = %v", err)
	}

	if result == nil {
		t.Fatal("CreateIsometricTile() returned nil sprite")
	}

	if result.Width != options.TileWidth {
		t.Errorf("expected width %d, got %d", options.TileWidth, result.Width)
	}

	expectedHeight := options.TileHeight + options.Height
	if result.Height != expectedHeight {
		t.Errorf("expected height %d, got %d", expectedHeight, result.Height)
	}

	if result.Type != cartio.SpriteTypeStatic {
		t.Errorf("expected type static, got %s", result.Type)
	}
}

func TestCreateIsometricTile_AllLightingModes(t *testing.T) {
	palette := generateTestPalette()
	converter := NewIsometricConverter(64, 32)

	spriteMap := cartio.SpriteMap{
		"top":   {Width: 16, Height: 16, Type: cartio.SpriteTypeStatic, Pixels: generateTestPixels(16, 16), IsUI: false},
		"left":  {Width: 16, Height: 16, Type: cartio.SpriteTypeStatic, Pixels: generateTestPixels(16, 16), IsUI: false},
		"right": {Width: 16, Height: 16, Type: cartio.SpriteTypeStatic, Pixels: generateTestPixels(16, 16), IsUI: false},
	}

	modes := []LightingMode{LightingNormal, LightingBasic, LightingFull, LightingGradient}

	for _, mode := range modes {
		t.Run(string(mode), func(t *testing.T) {
			options := DefaultTileOptions()
			options.LightingMode = mode

			result, err := converter.CreateIsometricTile(
				"top", "left", "right",
				"", "", "",
				palette,
				spriteMap,
				options,
			)

			if err != nil {
				t.Fatalf("CreateIsometricTile() error = %v", err)
			}

			if result == nil {
				t.Fatal("CreateIsometricTile() returned nil sprite")
			}
		})
	}
}

func TestCreateIsometricTile_DefaultOptions(t *testing.T) {
	palette := generateTestPalette()
	converter := NewIsometricConverter(64, 32)

	spriteMap := cartio.SpriteMap{
		"top":   {Width: 16, Height: 16, Type: cartio.SpriteTypeStatic, Pixels: generateTestPixels(16, 16), IsUI: false},
		"left":  {Width: 16, Height: 16, Type: cartio.SpriteTypeStatic, Pixels: generateTestPixels(16, 16), IsUI: false},
		"right": {Width: 16, Height: 16, Type: cartio.SpriteTypeStatic, Pixels: generateTestPixels(16, 16), IsUI: false},
	}

	// Test with zero options (should use defaults)
	options := TileOptions{}

	result, err := converter.CreateIsometricTile(
		"top", "left", "right",
		"", "", "",
		palette,
		spriteMap,
		options,
	)

	if err != nil {
		t.Fatalf("CreateIsometricTile() error = %v", err)
	}

	if result == nil {
		t.Fatal("CreateIsometricTile() returned nil sprite")
	}
}

func TestCreateIsometricTile_Errors(t *testing.T) {
	palette := generateTestPalette()
	converter := NewIsometricConverter(64, 32)

	tests := []struct {
		name      string
		spriteMap cartio.SpriteMap
		wantError bool
	}{
		{
			name:      "missing top sprite",
			spriteMap: cartio.SpriteMap{},
			wantError: true,
		},
		{
			name: "missing left sprite",
			spriteMap: cartio.SpriteMap{
				"top": {Width: 16, Height: 16, Type: cartio.SpriteTypeStatic, Pixels: generateTestPixels(16, 16)},
			},
			wantError: true,
		},
		{
			name: "missing right sprite",
			spriteMap: cartio.SpriteMap{
				"top":  {Width: 16, Height: 16, Type: cartio.SpriteTypeStatic, Pixels: generateTestPixels(16, 16)},
				"left": {Width: 16, Height: 16, Type: cartio.SpriteTypeStatic, Pixels: generateTestPixels(16, 16)},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := DefaultTileOptions()
			_, err := converter.CreateIsometricTile(
				"top", "left", "right",
				"", "", "",
				palette,
				tt.spriteMap,
				options,
			)

			if (err != nil) != tt.wantError {
				t.Errorf("CreateIsometricTile() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestCreateIsometricTile_WithFrames(t *testing.T) {
	palette := generateTestPalette()
	converter := NewIsometricConverter(64, 32)

	spriteMap := cartio.SpriteMap{
		"top": {
			Width:  16,
			Height: 16,
			Type:   cartio.SpriteTypeFrames,
			Frames: []cartio.SpriteFrame{
				{
					Name:   "default",
					Pixels: generateTestPixels(16, 16),
				},
			},
			IsUI: false,
		},
		"left": {
			Width:  16,
			Height: 16,
			Type:   cartio.SpriteTypeFrames,
			Frames: []cartio.SpriteFrame{
				{
					Name:   "default",
					Pixels: generateTestPixels(16, 16),
				},
			},
			IsUI: false,
		},
		"right": {
			Width:  16,
			Height: 16,
			Type:   cartio.SpriteTypeFrames,
			Frames: []cartio.SpriteFrame{
				{
					Name:   "default",
					Pixels: generateTestPixels(16, 16),
				},
			},
			IsUI: false,
		},
	}

	options := DefaultTileOptions()

	result, err := converter.CreateIsometricTile(
		"top", "left", "right",
		"default", "default", "default",
		palette,
		spriteMap,
		options,
	)

	if err != nil {
		t.Fatalf("CreateIsometricTile() error = %v", err)
	}

	if result == nil {
		t.Fatal("CreateIsometricTile() returned nil sprite")
	}
}

func TestCreateIsometricTile_InvalidOptions(t *testing.T) {
	palette := generateTestPalette()
	converter := NewIsometricConverter(64, 32)

	spriteMap := cartio.SpriteMap{
		"top":   {Width: 16, Height: 16, Type: cartio.SpriteTypeStatic, Pixels: generateTestPixels(16, 16)},
		"left":  {Width: 16, Height: 16, Type: cartio.SpriteTypeStatic, Pixels: generateTestPixels(16, 16)},
		"right": {Width: 16, Height: 16, Type: cartio.SpriteTypeStatic, Pixels: generateTestPixels(16, 16)},
	}

	// Test invalid palette
	_, err := converter.CreateIsometricTile(
		"top", "left", "right",
		"", "", "",
		[]string{"#000000"}, // Too short
		spriteMap,
		DefaultTileOptions(),
	)

	if err == nil {
		t.Error("expected error for invalid palette")
	}

	// Test invalid dimensions (negative values)
	// Note: zero values use defaults, so they're valid
	testCases := []struct {
		name    string
		options TileOptions
	}{
		{"negative width", TileOptions{TileWidth: -1, TileHeight: 32, Height: 16, LightingMode: LightingGradient}},
		{"negative height", TileOptions{TileWidth: 64, TileHeight: -1, Height: 16, LightingMode: LightingGradient}},
		{"negative side height", TileOptions{TileWidth: 64, TileHeight: 32, Height: -1, LightingMode: LightingGradient}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := converter.CreateIsometricTile(
				"top", "left", "right",
				"", "", "",
				palette,
				spriteMap,
				tc.options,
			)

			if err == nil {
				t.Error("expected error for invalid dimensions")
			}
		})
	}
}

func TestScaleSideFace(t *testing.T) {
	// Create a test image
	testImg := createTestImage(8, 16, color.RGBA{R: 255, G: 0, B: 0, A: 255})

	// Test scaling down (width and height)
	scaled := scaleSideFace(testImg, 4, 8)
	bounds := scaled.Bounds()
	if bounds.Dx() != 4 || bounds.Dy() != 8 {
		t.Errorf("expected 4x8, got %dx%d", bounds.Dx(), bounds.Dy())
	}

	// Test scaling up
	scaled = scaleSideFace(testImg, 16, 32)
	bounds = scaled.Bounds()
	if bounds.Dx() != 16 || bounds.Dy() != 32 {
		t.Errorf("expected 16x32, got %dx%d", bounds.Dx(), bounds.Dy())
	}

	// Test same size
	scaled = scaleSideFace(testImg, 8, 16)
	bounds = scaled.Bounds()
	if bounds.Dx() != 8 || bounds.Dy() != 16 {
		t.Errorf("expected 8x16, got %dx%d", bounds.Dx(), bounds.Dy())
	}
	
	// Test isometric tile side face dimensions (32x16 for 64-wide tile)
	scaled = scaleSideFace(testImg, 32, 16)
	bounds = scaled.Bounds()
	if bounds.Dx() != 32 || bounds.Dy() != 16 {
		t.Errorf("expected 32x16 for isometric side face, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

// generateTestPixels creates a simple test pixel array
func generateTestPixels(width, height int) [][]int {
	pixels := make([][]int, height)
	for y := 0; y < height; y++ {
		pixels[y] = make([]int, width)
		for x := 0; x < width; x++ {
			pixels[y][x] = (x + y) % 50 // Cycle through palette indices
		}
	}
	return pixels
}
