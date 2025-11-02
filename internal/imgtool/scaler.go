package imgtool

import (
	"image"
	"image/color"
)

// Scale resizes an image to target dimensions
func Scale(img image.Image, opts ScaleOptions) ([][][]uint8, error) {
	bounds := img.Bounds()
	srcWidth := bounds.Dx()
	srcHeight := bounds.Dy()

	targetWidth := opts.Width
	targetHeight := opts.Height

	// Ensure divisible by 2 if requested
	if opts.EnsureDivisible {
		if targetWidth%2 != 0 {
			targetWidth++
		}
		if targetHeight%2 != 0 {
			targetHeight++
		}
	}

	// Handle preserve aspect ratio
	if opts.PreserveAspect {
		aspectRatio := float64(srcWidth) / float64(srcHeight)
		targetAspectRatio := float64(targetWidth) / float64(targetHeight)

		if aspectRatio > targetAspectRatio {
			// Source is wider - fit to width
			targetHeight = int(float64(targetWidth) / aspectRatio)
			if opts.EnsureDivisible && targetHeight%2 != 0 {
				targetHeight++
			}
		} else {
			// Source is taller - fit to height
			targetWidth = int(float64(targetHeight) * aspectRatio)
			if opts.EnsureDivisible && targetWidth%2 != 0 {
				targetWidth++
			}
		}
	}

	// Scale the image
	scaled := ScaleImage(img, srcWidth, srcHeight, targetWidth, targetHeight, opts.Algorithm)

	// Convert to RGB array
	result := make([][][]uint8, targetHeight)
	for y := 0; y < targetHeight; y++ {
		result[y] = make([][]uint8, targetWidth)
		for x := 0; x < targetWidth; x++ {
			c := scaled.At(x, y)
			r, g, b, a := c.RGBA()
			result[y][x] = []uint8{
				uint8(r >> 8),
				uint8(g >> 8),
				uint8(b >> 8),
				uint8(a >> 8),
			}
		}
	}

	return result, nil
}

// ScaleImage performs the actual scaling operation (exported for internal use)
func ScaleImage(img image.Image, srcWidth, srcHeight, targetWidth, targetHeight int, algorithm string) image.Image {
	switch algorithm {
	case "bilinear":
		return scaleBilinear(img, srcWidth, srcHeight, targetWidth, targetHeight)
	case "bicubic":
		return scaleBicubic(img, srcWidth, srcHeight, targetWidth, targetHeight)
	default: // "nearest"
		return scaleNearest(img, srcWidth, srcHeight, targetWidth, targetHeight)
	}
}

// scaleNearest performs nearest-neighbor scaling
func scaleNearest(img image.Image, srcWidth, srcHeight, targetWidth, targetHeight int) image.Image {
	rgba := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))

	xRatio := float64(srcWidth) / float64(targetWidth)
	yRatio := float64(srcHeight) / float64(targetHeight)

	for y := 0; y < targetHeight; y++ {
		for x := 0; x < targetWidth; x++ {
			srcX := int(float64(x) * xRatio)
			srcY := int(float64(y) * yRatio)
			
			if srcX >= srcWidth {
				srcX = srcWidth - 1
			}
			if srcY >= srcHeight {
				srcY = srcHeight - 1
			}

			rgba.Set(x, y, img.At(srcX, srcY))
		}
	}

	return rgba
}

// scaleBilinear performs bilinear interpolation scaling
func scaleBilinear(img image.Image, srcWidth, srcHeight, targetWidth, targetHeight int) image.Image {
	rgba := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))

	xRatio := float64(srcWidth-1) / float64(targetWidth)
	yRatio := float64(srcHeight-1) / float64(targetHeight)

	for y := 0; y < targetHeight; y++ {
		for x := 0; x < targetWidth; x++ {
			srcX := float64(x) * xRatio
			srcY := float64(y) * yRatio

			x1 := int(srcX)
			y1 := int(srcY)
			x2 := x1 + 1
			y2 := y1 + 1

			if x2 >= srcWidth {
				x2 = srcWidth - 1
			}
			if y2 >= srcHeight {
				y2 = srcHeight - 1
			}

			// Get four corner pixels
			c11 := img.At(x1, y1)
			c12 := img.At(x1, y2)
			c21 := img.At(x2, y1)
			c22 := img.At(x2, y2)

			// Interpolate
			fx := srcX - float64(x1)
			fy := srcY - float64(y1)

			c := bilinearInterpolate(c11, c12, c21, c22, fx, fy)
			rgba.Set(x, y, c)
		}
	}

	return rgba
}

// bilinearInterpolate performs bilinear interpolation between four colors
func bilinearInterpolate(c11, c12, c21, c22 color.Color, fx, fy float64) color.Color {
	r11, g11, b11, a11 := c11.RGBA()
	r12, g12, b12, a12 := c12.RGBA()
	r21, g21, b21, a21 := c21.RGBA()
	r22, g22, b22, a22 := c22.RGBA()

	r := lerp(lerp(float64(r11), float64(r21), fx), lerp(float64(r12), float64(r22), fx), fy)
	g := lerp(lerp(float64(g11), float64(g21), fx), lerp(float64(g12), float64(g22), fx), fy)
	b := lerp(lerp(float64(b11), float64(b21), fx), lerp(float64(b12), float64(b22), fx), fy)
	a := lerp(lerp(float64(a11), float64(a21), fx), lerp(float64(a12), float64(a22), fx), fy)

	return color.RGBA{
		R: uint8(r / 256.0),
		G: uint8(g / 256.0),
		B: uint8(b / 256.0),
		A: uint8(a / 256.0),
	}
}

// scaleBicubic performs bicubic interpolation scaling (simplified - using bilinear for now)
func scaleBicubic(img image.Image, srcWidth, srcHeight, targetWidth, targetHeight int) image.Image {
	// For now, use bilinear as bicubic is complex
	// TODO: Implement proper bicubic interpolation
	return scaleBilinear(img, srcWidth, srcHeight, targetWidth, targetHeight)
}

// lerp performs linear interpolation
func lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

