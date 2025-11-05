package tile2iso

import (
	"image"
	"image/color"
	"math"
)

// TransformToIsometric transforms a square/rectangular texture to isometric diamond projection
// Uses simple rotation and scaling approach:
// 1. Rotate 45 degrees clockwise
// 2. Scale vertically by 0.5
// This is the exported version for use by the engine
func TransformToIsometric(src image.Image, dstWidth, dstHeight int) (image.Image, error) {
	return transformToIsometric(src, dstWidth, dstHeight)
}

// transformToIsometric transforms a square/rectangular texture to isometric diamond projection
// Uses simple rotation and scaling approach:
// 1. Rotate 45 degrees clockwise
// 2. Scale vertically by 0.5
func transformToIsometric(src image.Image, dstWidth, dstHeight int) (image.Image, error) {
	if src == nil {
		return nil, ErrNilTexture
	}

	srcBounds := src.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()

	// Create destination image
	dst := image.NewRGBA(image.Rect(0, 0, dstWidth, dstHeight))

	// Calculate center points
	srcCenterX := float64(srcWidth) / 2.0
	srcCenterY := float64(srcHeight) / 2.0
	dstCenterX := float64(dstWidth) / 2.0
	dstCenterY := float64(dstHeight) / 2.0

	// For each pixel in destination, find corresponding source pixel
	for y := 0; y < dstHeight; y++ {
		for x := 0; x < dstWidth; x++ {
			// Convert destination coordinates to centered coordinates
			dx := float64(x) - dstCenterX
			dy := float64(y) - dstCenterY

			// Apply inverse isometric transformation
			// Inverse of: rotate 45° and scale Y by 0.5
			// This is: scale Y by 2.0 and rotate -45°
			angle := -45.0 * math.Pi / 180.0
			cos := math.Cos(angle)
			sin := math.Sin(angle)

			// Scale Y back
			dyScaled := dy * 2.0

			// Rotate back
			sx := dx*cos - dyScaled*sin
			sy := dx*sin + dyScaled*cos

			// Convert back to source coordinates
			srcX := sx + srcCenterX
			srcY := sy + srcCenterY

			// Sample from source using bilinear interpolation
			c := bilinearSample(src, srcX, srcY, srcBounds)
			dst.Set(x, y, c)
		}
	}

	return dst, nil
}

// bilinearSample performs bilinear interpolation to sample a pixel from an image
func bilinearSample(img image.Image, x, y float64, bounds image.Rectangle) color.Color {
	// Clamp to bounds
	x = math.Max(float64(bounds.Min.X), math.Min(float64(bounds.Max.X-1), x))
	y = math.Max(float64(bounds.Min.Y), math.Min(float64(bounds.Max.Y-1), y))

	// Get integer coordinates
	x0 := int(math.Floor(x))
	y0 := int(math.Floor(y))
	x1 := x0 + 1
	y1 := y0 + 1

	// Clamp to bounds
	if x1 >= bounds.Max.X {
		x1 = bounds.Max.X - 1
	}
	if y1 >= bounds.Max.Y {
		y1 = bounds.Max.Y - 1
	}

	// Calculate fractional parts
	fx := x - float64(x0)
	fy := y - float64(y0)

	// Get four corner pixels
	c00 := img.At(x0, y0)
	c10 := img.At(x1, y0)
	c01 := img.At(x0, y1)
	c11 := img.At(x1, y1)

	// Convert to RGBA
	rgba00 := toRGBA(c00)
	rgba10 := toRGBA(c10)
	rgba01 := toRGBA(c01)
	rgba11 := toRGBA(c11)

	// Interpolate
	interp := func(c00, c10, c01, c11 uint8, fx, fy float64) uint8 {
		top := float64(c00)*(1-fx) + float64(c10)*fx
		bottom := float64(c01)*(1-fx) + float64(c11)*fx
		return uint8(top*(1-fy) + bottom*fy)
	}

	return color.RGBA{
		R: interp(rgba00.R, rgba10.R, rgba01.R, rgba11.R, fx, fy),
		G: interp(rgba00.G, rgba10.G, rgba01.G, rgba11.G, fx, fy),
		B: interp(rgba00.B, rgba10.B, rgba01.B, rgba11.B, fx, fy),
		A: interp(rgba00.A, rgba10.A, rgba01.A, rgba11.A, fx, fy),
	}
}

// toRGBA converts any color to RGBA
func toRGBA(c color.Color) color.RGBA {
	rgba, ok := c.(color.RGBA)
	if ok {
		return rgba
	}
	r, g, b, a := c.RGBA()
	return color.RGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: uint8(a >> 8),
	}
}

