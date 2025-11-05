package tile2iso

import (
	"fmt"
	"image"
	"image/color"
)

// PixelDataToImage converts palette-indexed pixel data to an image.Image
// paletteColors should be a slice of 50 hex color strings
func PixelDataToImage(pixels [][]int, paletteColors []string) (image.Image, error) {
	if len(pixels) == 0 {
		return nil, ErrInvalidDimensions
	}

	height := len(pixels)
	width := len(pixels[0])
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Parse palette colors
	palette := make([]color.RGBA, 50)
	for i, hex := range paletteColors {
		if i >= 50 {
			break
		}
		r, g, b, a := parseHexColor(hex)
		palette[i] = color.RGBA{R: r, G: g, B: b, A: a}
	}

	// Convert pixels to image
	for y, row := range pixels {
		for x, idx := range row {
			var c color.RGBA
			if idx == -1 {
				// Transparent
				c = color.RGBA{A: 0}
			} else if idx >= 0 && idx < 50 {
				c = palette[idx]
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
func ImageToPixelData(img image.Image, paletteColors []string) ([][]int, []string, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Parse palette colors
	palette := make([]color.RGBA, 50)
	for i, hex := range paletteColors {
		if i >= 50 {
			break
		}
		r, g, b, a := parseHexColor(hex)
		palette[i] = color.RGBA{R: r, G: g, B: b, A: a}
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

			// Find closest palette color
			bestIdx := 0
			bestDist := colorDistance(rgba, palette[0])
			for i := 1; i < 50; i++ {
				dist := colorDistance(rgba, palette[i])
				if dist < bestDist {
					bestDist = dist
					bestIdx = i
				}
			}
			pixels[y][x] = bestIdx
		}
	}

	return pixels, paletteColors, nil
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
