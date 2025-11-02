package imgtool

import (
	"image/color"
	"math"
)

// ColorDistance calculates the perceptual distance between two colors using Euclidean distance in RGB space
// For better accuracy, consider using CIEDE2000, but this is simpler and faster
func ColorDistance(c1, c2 Color) float64 {
	r := float64(c1.R) - float64(c2.R)
	g := float64(c1.G) - float64(c2.G)
	b := float64(c1.B) - float64(c2.B)
	return math.Sqrt(r*r + g*g + b*b)
}

// FindClosestPaletteColor finds the closest palette color index for a given RGB color
func FindClosestPaletteColor(target Color, palette *Palette, cache *ColorCache) (int, error) {
	// Check cache first
	if cache != nil {
		if idx, ok := cache.Get(target.R, target.G, target.B); ok {
			return idx, nil
		}
	}

	// Validate palette
	if err := palette.Validate(); err != nil {
		return -1, err
	}

	minDist := math.MaxFloat64
	closestIdx := 0

	// Compare against all palette colors
	for i, hex := range palette.Colors {
		paletteColor, err := HexToColor(hex)
		if err != nil {
			continue
		}

		dist := ColorDistance(target, paletteColor)
		if dist < minDist {
			minDist = dist
			closestIdx = i
		}
	}

	// Cache the result
	if cache != nil {
		cache.Set(target.R, target.G, target.B, closestIdx)
	}

	return closestIdx, nil
}

// RGBAToColor converts color.RGBA to our Color type
func RGBAToColor(c color.RGBA) Color {
	return Color{R: c.R, G: c.G, B: c.B}
}

// ExtractColor extracts Color from any image color type
func ExtractColor(c color.Color) Color {
	r, g, b, _ := c.RGBA()
	return Color{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
	}
}

