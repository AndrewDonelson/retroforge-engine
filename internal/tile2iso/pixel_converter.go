package tile2iso

import (
	"fmt"
	"image"
	"image/color"
)

// PixelDataToImage converts palette-indexed pixel data to an image.Image
// paletteColors should be a slice of 48 hex color strings (game palette)
// Built-in colors (0-15) are handled separately
// For full 64-color support, use GetFullPalette from pal package
func PixelDataToImage(pixels [][]int, paletteColors []string, builtinColors []color.RGBA) (image.Image, error) {
	if len(pixels) == 0 {
		return nil, ErrInvalidDimensions
	}

	height := len(pixels)
	width := len(pixels[0])
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Parse game palette colors (48 colors)
	gamePalette := make([]color.RGBA, 48)
	for i, hex := range paletteColors {
		if i >= 48 {
			break
		}
		r, g, b, a := parseHexColor(hex)
		gamePalette[i] = color.RGBA{R: r, G: g, B: b, A: a}
	}

	// Convert pixels to image
	for y, row := range pixels {
		for x, idx := range row {
			var c color.RGBA
			if idx == -1 {
				// Transparent
				c = color.RGBA{A: 0}
			} else if idx >= 0 && idx < 16 {
				// Built-in color (0-15)
				if idx < len(builtinColors) {
					c = builtinColors[idx]
				} else {
					c = color.RGBA{A: 255} // Black fallback
				}
			} else if idx >= 16 && idx < 64 {
				// Game palette color (16-63)
				gameIdx := idx - 16
				if gameIdx < len(gamePalette) {
					c = gamePalette[gameIdx]
				} else {
					c = color.RGBA{A: 255} // Black fallback
				}
			} else {
				// Invalid index, use black
				c = color.RGBA{A: 255}
			}
			img.Set(x, y, c)
		}
	}

	return img, nil
}

// ImageToPixelData converts an image.Image back to palette-indexed pixel data
// Returns the pixel data and the palette colors used
// fullPalette should be the complete 64-color palette (16 built-in + 48 game)
func ImageToPixelData(img image.Image, fullPalette []color.RGBA) ([][]int, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if len(fullPalette) < 64 {
		return nil, fmt.Errorf("full palette must have 64 colors (got %d)", len(fullPalette))
	}

	// Convert image to pixels
	pixels := make([][]int, height)
	for y := 0; y < height; y++ {
		pixels[y] = make([]int, width)
		for x := 0; x < width; x++ {
			c := img.At(x, y)
			rgba, ok := c.(color.RGBA)
			if !ok {
				// Convert to RGBA
				r, g, b, a := c.RGBA()
				rgba = color.RGBA{
					R: uint8(r >> 8),
					G: uint8(g >> 8),
					B: uint8(b >> 8),
					A: uint8(a >> 8),
				}
			}

			// Check for transparency
			if rgba.A < 128 {
				pixels[y][x] = -1
				continue
			}

			// Find closest palette color in full 64-color palette
			bestIdx := 0
			bestDist := colorDistance(rgba, fullPalette[0])
			for i := 1; i < 64; i++ {
				dist := colorDistance(rgba, fullPalette[i])
				if dist < bestDist {
					bestDist = dist
					bestIdx = i
				}
			}
			pixels[y][x] = bestIdx
		}
	}

	return pixels, nil
}

// parseHexColor parses a hex color string (#RRGGBB) to RGBA
func parseHexColor(hex string) (r, g, b, a uint8) {
	if len(hex) != 7 || hex[0] != '#' {
		return 0, 0, 0, 255
	}

	var rv, gv, bv uint32
	fmt.Sscanf(hex[1:3], "%02x", &rv)
	fmt.Sscanf(hex[3:5], "%02x", &gv)
	fmt.Sscanf(hex[5:7], "%02x", &bv)

	return uint8(rv), uint8(gv), uint8(bv), 255
}

// colorDistance calculates the Euclidean distance between two colors
func colorDistance(c1, c2 color.RGBA) float64 {
	dr := float64(c1.R) - float64(c2.R)
	dg := float64(c1.G) - float64(c2.G)
	db := float64(c1.B) - float64(c2.B)
	return dr*dr + dg*dg + db*db
}
