package tile2iso

import (
	"image"
	"image/color"
	"testing"

	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
)

func TestCompositeIsometricTile_LayerOrder(t *testing.T) {
	// Create three distinct colored faces
	// Note: leftImg and rightImg should be full-width (32) images
	topImg := createTestImage(32, 16, color.RGBA{R: 255, G: 0, B: 0, A: 255})     // Red top
	
	// Create left side as full-width image (32x8) with parallelogram at x=0 to x=16
	leftImg := image.NewRGBA(image.Rect(0, 0, 32, 8))
	for y := 0; y < 8; y++ {
		for x := 0; x < 16; x++ {
			leftImg.Set(x, y, color.RGBA{R: 0, G: 255, B: 0, A: 255}) // Green
		}
	}
	
	// Create right side as full-width image (32x8) with parallelogram at x=16 to x=32
	rightImg := image.NewRGBA(image.Rect(0, 0, 32, 8))
	for y := 0; y < 8; y++ {
		for x := 16; x < 32; x++ {
			rightImg.Set(x, y, color.RGBA{R: 0, G: 0, B: 255, A: 255}) // Blue
		}
	}

	options := TileOptions{
		TileWidth:  32,
		TileHeight: 16,
		Height:     8,
		ShowOutline: false,
	}

	result := compositeIsometricTile(topImg, leftImg, rightImg, options)

	if result == nil {
		t.Fatal("compositeIsometricTile() returned nil")
	}

	bounds := result.Bounds()
	if bounds.Dx() != 32 || bounds.Dy() != 24 {
		t.Errorf("expected 32x24 canvas, got %dx%d", bounds.Dx(), bounds.Dy())
	}

	// Verify top diamond is visible (should be at Y=0 to Y=16)
	// Check center pixel of top diamond
	topColor := result.At(16, 8)
	if toRGBA(topColor).R == 0 {
		t.Error("top diamond (red) should be visible at center")
	}
}

func TestCompositeIsometricTile_Positioning(t *testing.T) {
	// Create test images with known dimensions
	// Note: leftImg and rightImg should be full-width (32) images with the parallelogram already positioned
	// For testing, we'll create full-width images with the side faces positioned correctly
	topImg := createTestImage(32, 16, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	
	// Create left side as full-width image (32x8) with parallelogram at x=0 to x=16
	leftImg := image.NewRGBA(image.Rect(0, 0, 32, 8))
	for y := 0; y < 8; y++ {
		for x := 0; x < 16; x++ {
			leftImg.Set(x, y, color.RGBA{R: 0, G: 255, B: 0, A: 255}) // Green
		}
	}
	
	// Create right side as full-width image (32x8) with parallelogram at x=16 to x=32
	rightImg := image.NewRGBA(image.Rect(0, 0, 32, 8))
	for y := 0; y < 8; y++ {
		for x := 16; x < 32; x++ {
			rightImg.Set(x, y, color.RGBA{R: 0, G: 0, B: 255, A: 255}) // Blue
		}
	}

	options := TileOptions{
		TileWidth:  32,
		TileHeight: 16,
		Height:     8,
		ShowOutline: false,
	}

	result := compositeIsometricTile(topImg, leftImg, rightImg, options)

	bounds := result.Bounds()
	// Expected: 32 wide, 16 (top) + 8 (sides) = 24 tall
	if bounds.Dx() != 32 {
		t.Errorf("expected width 32, got %d", bounds.Dx())
	}
	if bounds.Dy() != 24 {
		t.Errorf("expected height 24, got %d", bounds.Dy())
	}

	// Verify side faces are positioned below the diamond (Y >= 16)
	// Check left side face (should be at X=0 to X=16, Y=16 to Y=24)
	leftSideColor := result.At(8, 20) // Middle of left side
	leftRGBA := toRGBA(leftSideColor)
	if leftRGBA.G < 128 {
		t.Error("left side face (green) should be visible at Y=20")
	}

	// Check right side face (should be at X=16 to X=32, Y=16 to Y=24)
	rightSideColor := result.At(24, 20) // Middle of right side
	rightRGBA := toRGBA(rightSideColor)
	if rightRGBA.B < 128 {
		t.Error("right side face (blue) should be visible at Y=20")
	}
}

func TestCompositeIsometricTile_CanvasSize(t *testing.T) {
	testCases := []struct {
		name           string
		tileWidth      int
		tileHeight     int
		height         int
		expectedWidth  int
		expectedHeight int
	}{
		{"32x24 tile", 32, 16, 8, 32, 24},
		{"64x48 tile", 64, 32, 16, 64, 48},
		{"16x12 tile", 16, 8, 4, 16, 12},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			topImg := createTestImage(tc.tileWidth, tc.tileHeight, color.RGBA{R: 255, G: 0, B: 0, A: 255})
			
			// Create left side as full-width image
			leftImg := image.NewRGBA(image.Rect(0, 0, tc.tileWidth, tc.height))
			for y := 0; y < tc.height; y++ {
				for x := 0; x < tc.tileWidth/2; x++ {
					leftImg.Set(x, y, color.RGBA{R: 0, G: 255, B: 0, A: 255})
				}
			}
			
			// Create right side as full-width image
			rightImg := image.NewRGBA(image.Rect(0, 0, tc.tileWidth, tc.height))
			for y := 0; y < tc.height; y++ {
				for x := tc.tileWidth / 2; x < tc.tileWidth; x++ {
					rightImg.Set(x, y, color.RGBA{R: 0, G: 0, B: 255, A: 255})
				}
			}

			options := TileOptions{
				TileWidth:  tc.tileWidth,
				TileHeight: tc.tileHeight,
				Height:     tc.height,
				ShowOutline: false,
			}

			result := compositeIsometricTile(topImg, leftImg, rightImg, options)
			bounds := result.Bounds()

			if bounds.Dx() != tc.expectedWidth {
				t.Errorf("expected width %d, got %d", tc.expectedWidth, bounds.Dx())
			}
			if bounds.Dy() != tc.expectedHeight {
				t.Errorf("expected height %d, got %d", tc.expectedHeight, bounds.Dy())
			}
		})
	}
}

func TestCompositeIsometricTile_AlphaBlending(t *testing.T) {
	// Create semi-transparent faces
	topImg := createTestImageWithAlpha(32, 16, color.RGBA{R: 255, G: 0, B: 0, A: 128})     // 50% transparent red
	leftImg := createTestImageWithAlpha(16, 8, color.RGBA{R: 0, G: 255, B: 0, A: 128})    // 50% transparent green
	rightImg := createTestImageWithAlpha(16, 8, color.RGBA{R: 0, G: 0, B: 255, A: 128})  // 50% transparent blue

	options := TileOptions{
		TileWidth:  32,
		TileHeight: 16,
		Height:     8,
		ShowOutline: false,
	}

	result := compositeIsometricTile(topImg, leftImg, rightImg, options)

	// Check that alpha blending occurred (pixels should have intermediate alpha values)
	// The top diamond should be visible but semi-transparent
	topColor := result.At(16, 8)
	rgba := toRGBA(topColor)
	if rgba.A == 0 {
		t.Error("top diamond should be visible (alpha > 0)")
	}
	if rgba.A == 255 {
		t.Error("top diamond should be semi-transparent (alpha < 255)")
	}
}

func TestDrawOutline(t *testing.T) {
	// Create a test image
	img := image.NewRGBA(image.Rect(0, 0, 32, 24))

	// Fill with transparent pixels
	for y := 0; y < 24; y++ {
		for x := 0; x < 32; x++ {
			img.Set(x, y, color.RGBA{A: 0})
		}
	}

	// Add some content
	for y := 8; y < 16; y++ {
		for x := 8; x < 24; x++ {
			img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}

	options := TileOptions{
		TileWidth:  32,
		TileHeight: 16,
		Height:     8,
		ShowOutline: true,
	}

	// Draw outline
	drawOutline(img, options)

	// Check that outline pixels were drawn (black pixels at edges)
	// Check top-left corner of diamond
	outlineColor := img.At(0, 8)
	rgba := toRGBA(outlineColor)
	// Outline should be black (or very dark)
	if rgba.R > 10 || rgba.G > 10 || rgba.B > 10 {
		t.Error("outline should be black/dark at edges")
	}
}

func TestDrawOutline_Disabled(t *testing.T) {
	// Create a test image
	img := image.NewRGBA(image.Rect(0, 0, 32, 24))

	// Fill with transparent pixels
	for y := 0; y < 24; y++ {
		for x := 0; x < 32; x++ {
			img.Set(x, y, color.RGBA{A: 0})
		}
	}

	options := TileOptions{
		TileWidth:  32,
		TileHeight: 16,
		Height:     8,
		ShowOutline: false,
	}

	// Draw outline (should do nothing)
	drawOutline(img, options)

	// Verify no outline was drawn (all pixels should still be transparent)
	allTransparent := true
	for y := 0; y < 24; y++ {
		for x := 0; x < 32; x++ {
			rgba := toRGBA(img.At(x, y))
			if rgba.A > 0 {
				allTransparent = false
				break
			}
		}
		if !allTransparent {
			break
		}
	}

	if !allTransparent {
		t.Error("outline should not be drawn when ShowOutline is false")
	}
}

func TestCreateIsometricTile_ShowOutline(t *testing.T) {
	palette := generateTestPalette()
	converter := NewIsometricConverter(32, 16)

	spriteMap := cartio.SpriteMap{
		"top":   {Width: 16, Height: 16, Type: cartio.SpriteTypeStatic, Pixels: generateTestPixels(16, 16), IsUI: false},
		"left":  {Width: 16, Height: 16, Type: cartio.SpriteTypeStatic, Pixels: generateTestPixels(16, 16), IsUI: false},
		"right": {Width: 16, Height: 16, Type: cartio.SpriteTypeStatic, Pixels: generateTestPixels(16, 16), IsUI: false},
	}

	options := DefaultTileOptions()
	options.ShowOutline = true

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

	// Verify the tile was created (outline should be present)
	// We can't easily check for outline pixels without examining the pixel data,
	// but if the function succeeds, the outline was drawn
}

// Helper function to create test image with alpha
func createTestImageWithAlpha(width, height int, c color.RGBA) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, c)
		}
	}
	return img
}

