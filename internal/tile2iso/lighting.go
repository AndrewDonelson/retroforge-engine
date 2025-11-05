package tile2iso

import (
	"image"
	"image/color"
	"math"
)

// applyLighting applies the specified lighting mode to a side texture
func applyLighting(img image.Image, mode LightingMode, side string, height int) (image.Image, error) {
	if img == nil {
		return nil, ErrNilTexture
	}

	switch mode {
	case LightingNormal:
		return img, nil

	case LightingBasic:
		return applyUniformBrightness(img, side == "left", 1.2, 0.8)

	case LightingFull:
		return applyRegionalBrightness(img, side == "left", height)

	case LightingGradient:
		return applyGradientBrightness(img, side == "left", height)

	default:
		return nil, ErrInvalidLightingMode
	}
}

// applyUniformBrightness applies uniform brightness adjustment
// leftSide: true for left side, false for right side
// leftFactor: brightness multiplier for left side (typically 1.2)
// rightFactor: brightness multiplier for right side (typically 0.8)
func applyUniformBrightness(img image.Image, leftSide bool, leftFactor, rightFactor float64) (image.Image, error) {
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	factor := rightFactor
	if leftSide {
		factor = leftFactor
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			rgbaColor, ok := c.(color.RGBA)
			if !ok {
				r, g, b, a := c.RGBA()
				rgbaColor = color.RGBA{
					R: uint8(r >> 8),
					G: uint8(g >> 8),
					B: uint8(b >> 8),
					A: uint8(a >> 8),
				}
			}

			adjusted := adjustPixelBrightness(rgbaColor, factor)
			rgba.Set(x, y, adjusted)
		}
	}

	return rgba, nil
}

// applyRegionalBrightness applies enhanced lighting at top/bottom regions
func applyRegionalBrightness(img image.Image, leftSide bool, height int) (image.Image, error) {
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	imgHeight := bounds.Dy()

	// Calculate region boundaries (top 20%, middle 60%, bottom 20%)
	topBoundary := int(float64(imgHeight) * 0.2)
	bottomBoundary := int(float64(imgHeight) * 0.8)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		var factor float64
		if y < topBoundary {
			// Top 20%
			if leftSide {
				factor = 1.4
			} else {
				factor = 0.6
			}
		} else if y >= bottomBoundary {
			// Bottom 20%
			if leftSide {
				factor = 1.4
			} else {
				factor = 0.6
			}
		} else {
			// Middle 60%
			if leftSide {
				factor = 1.2
			} else {
				factor = 0.8
			}
		}

		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			rgbaColor, ok := c.(color.RGBA)
			if !ok {
				r, g, b, a := c.RGBA()
				rgbaColor = color.RGBA{
					R: uint8(r >> 8),
					G: uint8(g >> 8),
					B: uint8(b >> 8),
					A: uint8(a >> 8),
				}
			}

			adjusted := adjustPixelBrightness(rgbaColor, factor)
			rgba.Set(x, y, adjusted)
		}
	}

	return rgba, nil
}

// applyGradientBrightness applies smooth vertical gradient
func applyGradientBrightness(img image.Image, leftSide bool, height int) (image.Image, error) {
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	imgHeight := float64(bounds.Dy())

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		// Calculate gradient factor: 1.0 at top, 0.0 at bottom
		gradientFactor := 1.0 - (float64(y-bounds.Min.Y) / imgHeight)

		var brightness float64
		if leftSide {
			baseBrightness := 1.2
			brightness = baseBrightness + (gradientFactor*0.4 - 0.2)
			// Top: 1.2 + 0.2 = 1.4, Middle: 1.2, Bottom: 1.2 - 0.2 = 1.0
		} else {
			baseBrightness := 0.8
			brightness = baseBrightness - (gradientFactor*0.4 - 0.2)
			// Top: 0.8 - 0.2 = 0.6, Middle: 0.8, Bottom: 0.8 + 0.2 = 1.0
		}

		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			rgbaColor, ok := c.(color.RGBA)
			if !ok {
				r, g, b, a := c.RGBA()
				rgbaColor = color.RGBA{
					R: uint8(r >> 8),
					G: uint8(g >> 8),
					B: uint8(b >> 8),
					A: uint8(a >> 8),
				}
			}

			adjusted := adjustPixelBrightness(rgbaColor, brightness)
			rgba.Set(x, y, adjusted)
		}
	}

	return rgba, nil
}

// adjustPixelBrightness adjusts a pixel's brightness by a factor
func adjustPixelBrightness(c color.RGBA, factor float64) color.RGBA {
	newR := float64(c.R) * factor
	newG := float64(c.G) * factor
	newB := float64(c.B) * factor

	return color.RGBA{
		R: uint8(math.Min(255, math.Max(0, newR))),
		G: uint8(math.Min(255, math.Max(0, newG))),
		B: uint8(math.Min(255, math.Max(0, newB))),
		A: c.A, // Preserve alpha
	}
}

