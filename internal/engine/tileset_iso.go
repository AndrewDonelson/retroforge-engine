package engine

import (
	"image"
	"image/color"

	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
	"github.com/AndrewDonelson/retroforge-engine/internal/tile2iso"
)

// convertTileToIsometric converts a tile to isometric view
// Uses proper isometric transformation: 45° rotation + 50% vertical scale (creates diamond shape)
func (e *Engine) convertTileToIsometric(tile *cartio.TileData) error {
	// Get current palette colors
	paletteColors := e.getPaletteColors()

	// Handle different tile types
	switch tile.Type {
	case cartio.SpriteTypeStatic:
		// Convert static tile
		if err := convertTilePixelsToIsometric(&tile.Pixels, tile.Width, tile.Height, paletteColors); err != nil {
			return err
		}
		// For isometric: width stays same, height becomes width/2 (diamond shape)
		newWidth := tile.Width
		newHeight := tile.Width / 2  // Isometric diamond: height = width/2
		tile.Width = newWidth
		tile.Height = newHeight

	case cartio.SpriteTypeFrames, cartio.SpriteTypeAnimation:
		// Convert all frames
		for i := range tile.Frames {
			if err := convertTilePixelsToIsometric(&tile.Frames[i].Pixels, tile.Width, tile.Height, paletteColors); err != nil {
				return err
			}
		}
		// Update dimensions
		newWidth := tile.Width
		newHeight := tile.Width / 2
		tile.Width = newWidth
		tile.Height = newHeight
	}

	return nil
}

// convertTilePixelsToIsometric converts pixel data to isometric diamond shape
// Uses proper isometric: rotate 45° + scale Y by 0.5
func convertTilePixelsToIsometric(pixels *[][]int, width, height int, paletteColors []string) error {
	if len(*pixels) == 0 {
		return nil
	}

	// Convert pixels to image
	img, err := tile2iso.PixelDataToImage(*pixels, paletteColors)
	if err != nil {
		return err
	}

	// Use proper isometric transformation: 45° rotation + 50% vertical scale
	// This creates a diamond shape
	isometricImg, err := tile2iso.TransformToIsometric(img, width, width/2)
	if err != nil {
		return err
	}

	// Convert back to pixel data
	resultPixels, _, err := tile2iso.ImageToPixelData(isometricImg, paletteColors)
	if err != nil {
		return err
	}

	*pixels = resultPixels
	return nil
}

// isIsometricTileset checks if a tileset contains isometric tiles
// by checking if any tile has height = width/2 (characteristic of isometric tiles)
func isIsometricTileset(tileset cartio.TilesetMap) bool {
	for _, tile := range tileset {
		// Isometric tiles have height = width/2 (after conversion)
		if tile.Height > 0 && tile.Width > 0 {
			// Check if height is approximately width/2 (within 1 pixel tolerance)
			expectedHeight := tile.Width / 2
			if tile.Height == expectedHeight || tile.Height == expectedHeight+1 || tile.Height == expectedHeight-1 {
				return true
			}
		}
	}
	return false
}

// rotate90CW rotates an image 90 degrees clockwise
func rotate90CW(img image.Image) image.Image {
	bounds := img.Bounds()
	srcWidth := bounds.Dx()
	srcHeight := bounds.Dy()

	// After 90° CW rotation: new width = old height, new height = old width
	dst := image.NewRGBA(image.Rect(0, 0, srcHeight, srcWidth))

	for y := 0; y < srcHeight; y++ {
		for x := 0; x < srcWidth; x++ {
			// 90° CW rotation: (x, y) -> (y, width-1-x)
			dstX := srcHeight - 1 - y
			dstY := x
			dst.Set(dstX, dstY, img.At(x, y))
		}
	}

	return dst
}

// scaleVertical scales an image vertically by a factor (maintains width)
func scaleVertical(img image.Image, scaleY float64) image.Image {
	bounds := img.Bounds()
	srcWidth := bounds.Dx()
	srcHeight := bounds.Dy()

	newHeight := int(float64(srcHeight) * scaleY)
	if newHeight < 1 {
		newHeight = 1
	}

	dst := image.NewRGBA(image.Rect(0, 0, srcWidth, newHeight))

	for y := 0; y < newHeight; y++ {
		// Map destination y to source y
		srcY := int(float64(y) / scaleY)
		if srcY >= srcHeight {
			srcY = srcHeight - 1
		}

		for x := 0; x < srcWidth; x++ {
			dst.Set(x, y, img.At(x, srcY))
		}
	}

	return dst
}

// getPaletteColors gets current palette colors as hex strings
func (e *Engine) getPaletteColors() []string {
	colors := make([]string, 50)
	for i := 0; i < 50; i++ {
		c := e.Pal.Color(i)
		colors[i] = colorToHex(c)
	}
	return colors
}

// colorToHex converts color.RGBA to hex string
func colorToHex(c color.RGBA) string {
	return "#" + byteToHex(c.R) + byteToHex(c.G) + byteToHex(c.B)
}

// byteToHex converts a byte to two-character hex string
func byteToHex(b uint8) string {
	hex := "0123456789abcdef"
	return string(hex[b>>4]) + string(hex[b&0x0f])
}

