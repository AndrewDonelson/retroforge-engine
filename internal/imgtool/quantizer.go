package imgtool

import (
	"image"
	"sort"
)

// Quantize reduces an image to a 50-color palette
func Quantize(img image.Image, opts QuantizeOptions) (*Palette, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Extract unique colors from image
	colorMap := make(map[Color]int)
	colors := make([]Color, 0)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := img.At(x, y)
			
			// Handle transparency
			if IsTransparent(c, opts.AlphaThreshold) {
				continue
			}

			col := ExtractColor(c)
			if _, exists := colorMap[col]; !exists {
				colorMap[col] = len(colors)
				colors = append(colors, col)
			}
		}
	}

	// If we have 50 or fewer colors, use them directly
	if len(colors) <= 50 {
		palette := &Palette{Colors: make([]string, 50)}
		
		// Ensure black and white at indices 0 and 1

		idx := 0
		if opts.EnforceBlackWhite {
			palette.Colors[0] = "#000000"
			palette.Colors[1] = "#ffffff"
			idx = 2
		}

		// Add existing colors
		for _, c := range colors {
			if idx >= 50 {
				break
			}
			// Skip if it's black/white and we're enforcing
			if opts.EnforceBlackWhite {
				if (c.R == 0 && c.G == 0 && c.B == 0) || (c.R == 255 && c.G == 255 && c.B == 255) {
					continue
				}
			}
			palette.Colors[idx] = c.ToHex()
			idx++
		}

		// Fill remaining with grayscale if needed
		for idx < 50 {
			v := uint8((idx * 255) / 49)
			palette.Colors[idx] = Color{R: v, G: v, B: v}.ToHex()
			idx++
		}

		return palette, nil
	}

	// Too many colors - need quantization
	// Use median cut algorithm (simplified version)
	targetColors := 48 // Default to 48 for remaining slots after black/white
	if !opts.EnforceBlackWhite {
		targetColors = 50
	}
	quantized := medianCutQuantize(colors, targetColors)
	
	palette := &Palette{Colors: make([]string, 50)}
	
	idx := 0
	if opts.EnforceBlackWhite {
		palette.Colors[0] = "#000000"
		palette.Colors[1] = "#ffffff"
		idx = 2
	}

	for _, c := range quantized {
		if idx >= 50 {
			break
		}
		palette.Colors[idx] = c.ToHex()
		idx++
	}

	// Fill remaining slots
	for idx < 50 {
		v := uint8((idx * 255) / 49)
		palette.Colors[idx] = Color{R: v, G: v, B: v}.ToHex()
		idx++
	}

	return palette, nil
}

// medianCutQuantize performs median cut color quantization
func medianCutQuantize(colors []Color, targetColors int) []Color {
	if len(colors) <= targetColors {
		return colors
	}

	// Simple quantization: sort by brightness and group
	type colorWithBrightness struct {
		color     Color
		brightness float64
	}

	sorted := make([]colorWithBrightness, len(colors))
	for i, c := range colors {
		brightness := 0.299*float64(c.R) + 0.587*float64(c.G) + 0.114*float64(c.B)
		sorted[i] = colorWithBrightness{color: c, brightness: brightness}
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].brightness < sorted[j].brightness
	})

	// Group into buckets and average
	bucketSize := len(colors) / targetColors
	if bucketSize < 1 {
		bucketSize = 1
	}

	result := make([]Color, 0, targetColors)
	for i := 0; i < len(sorted); i += bucketSize {
		end := i + bucketSize
		if end > len(sorted) {
			end = len(sorted)
		}

		bucket := sorted[i:end]
		var sumR, sumG, sumB float64
		for _, item := range bucket {
			sumR += float64(item.color.R)
			sumG += float64(item.color.G)
			sumB += float64(item.color.B)
		}
		count := float64(len(bucket))
		result = append(result, Color{
			R: uint8(sumR / count),
			G: uint8(sumG / count),
			B: uint8(sumB / count),
		})

		if len(result) >= targetColors {
			break
		}
	}

	return result
}

