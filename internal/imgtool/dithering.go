package imgtool

import (
	"image"
	"math"
)

// DitherAlgorithm specifies the dithering algorithm to use
type DitherAlgorithm string

const (
	DitherNone           DitherAlgorithm = "none"
	DitherFloydSteinberg DitherAlgorithm = "floyd-steinberg"
	DitherOrdered        DitherAlgorithm = "ordered"
)

// ApplyDithering applies dithering to an image during palette mapping
// getPaletteIndex is a function that returns the palette index for a given color
// getPaletteColor is a function that returns the hex color string for a given palette index
func ApplyDithering(img image.Image, opts MapPaletteOptions, getPaletteIndex func(Color) int, getPaletteColor func(int) string) [][]int {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	result := make([][]int, height)
	errorBuffer := make([][][3]float64, height) // Error diffusion buffer [R,G,B]

	// Initialize result and error buffer
	for y := 0; y < height; y++ {
		result[y] = make([]int, width)
		errorBuffer[y] = make([][3]float64, width)
	}

	switch DitherAlgorithm(opts.DitherAlgorithm) {
	case DitherFloydSteinberg:
		applyFloydSteinberg(img, opts, width, height, result, errorBuffer, getPaletteIndex, getPaletteColor)
	case DitherOrdered:
		applyOrdered(img, opts, width, height, result, getPaletteIndex)
	default:
		applyNoDither(img, opts, width, height, result, getPaletteIndex)
	}

	return result
}

// applyNoDither maps pixels without dithering
func applyNoDither(img image.Image, opts MapPaletteOptions, width, height int, result [][]int, getPaletteIndex func(Color) int) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := img.At(x, y)
			if IsTransparent(c, opts.AlphaThreshold) {
				result[y][x] = opts.TransparentIndex
			} else {
				targetColor := ExtractColor(c)
				idx := getPaletteIndex(targetColor)
				result[y][x] = idx
			}
		}
	}
}

// applyFloydSteinberg applies Floyd-Steinberg error diffusion dithering
func applyFloydSteinberg(img image.Image, opts MapPaletteOptions, width, height int, result [][]int, errorBuffer [][][3]float64, getPaletteIndex func(Color) int, getPaletteColor func(int) string) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := img.At(x, y)
			
			if IsTransparent(c, opts.AlphaThreshold) {
				result[y][x] = opts.TransparentIndex
				continue
			}

			// Get original color and add error from previous pixels
			r, g, b, _ := c.RGBA()
			originalR := float64(r>>8) + errorBuffer[y][x][0]
			originalG := float64(g>>8) + errorBuffer[y][x][1]
			originalB := float64(b>>8) + errorBuffer[y][x][2]

			// Clamp to valid range
			originalR = math.Max(0, math.Min(255, originalR))
			originalG = math.Max(0, math.Min(255, originalG))
			originalB = math.Max(0, math.Min(255, originalB))

			// Find closest palette color
			targetColor := Color{
				R: uint8(originalR),
				G: uint8(originalG),
				B: uint8(originalB),
			}
			paletteIdx := getPaletteIndex(targetColor)

			// Get palette color
			result[y][x] = paletteIdx
			var paletteColor Color
			if getPaletteColor == nil {
				// Fallback if not provided
				paletteColor = Color{R: 0, G: 0, B: 0}
			} else {
				paletteColor, _ = HexToColor(getPaletteColor(paletteIdx))
			}

			// Calculate error
			errorR := originalR - float64(paletteColor.R)
			errorG := originalG - float64(paletteColor.G)
			errorB := originalB - float64(paletteColor.B)

			// Distribute error to neighboring pixels (Floyd-Steinberg pattern)
			if x+1 < width {
				errorBuffer[y][x+1][0] += errorR * 7.0 / 16.0
				errorBuffer[y][x+1][1] += errorG * 7.0 / 16.0
				errorBuffer[y][x+1][2] += errorB * 7.0 / 16.0
			}
			if y+1 < height {
				if x-1 >= 0 {
					errorBuffer[y+1][x-1][0] += errorR * 3.0 / 16.0
					errorBuffer[y+1][x-1][1] += errorG * 3.0 / 16.0
					errorBuffer[y+1][x-1][2] += errorB * 3.0 / 16.0
				}
				errorBuffer[y+1][x][0] += errorR * 5.0 / 16.0
				errorBuffer[y+1][x][1] += errorG * 5.0 / 16.0
				errorBuffer[y+1][x][2] += errorB * 5.0 / 16.0
				if x+1 < width {
					errorBuffer[y+1][x+1][0] += errorR * 1.0 / 16.0
					errorBuffer[y+1][x+1][1] += errorG * 1.0 / 16.0
					errorBuffer[y+1][x+1][2] += errorB * 1.0 / 16.0
				}
			}
		}
	}
}

// applyOrdered applies ordered (Bayer) dithering
func applyOrdered(img image.Image, opts MapPaletteOptions, width, height int, result [][]int, getPaletteIndex func(Color) int) {
	// 4x4 Bayer matrix
	bayerMatrix := [4][4]float64{
		{0, 8, 2, 10},
		{12, 4, 14, 6},
		{3, 11, 1, 9},
		{15, 7, 13, 5},
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := img.At(x, y)
			
			if IsTransparent(c, opts.AlphaThreshold) {
				result[y][x] = opts.TransparentIndex
				continue
			}

			// Get original color
			r, g, b, _ := c.RGBA()
			threshold := bayerMatrix[y%4][x%4] * 16.0 // Scale to 0-240

			// Apply threshold
			originalR := float64(r >> 8)
			originalG := float64(g >> 8)
			originalB := float64(b >> 8)

			if originalR < threshold {
				originalR = 0
			} else {
				originalR = 255
			}
			if originalG < threshold {
				originalG = 0
			} else {
				originalG = 255
			}
			if originalB < threshold {
				originalB = 0
			} else {
				originalB = 255
			}

			targetColor := Color{
				R: uint8(originalR),
				G: uint8(originalG),
				B: uint8(originalB),
			}
			result[y][x] = getPaletteIndex(targetColor)
		}
	}
}


