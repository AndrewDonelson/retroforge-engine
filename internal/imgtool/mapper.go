package imgtool

import "image"

// MapPalette converts an image to palette indices using closest color matching
func MapPalette(img image.Image, palette *Palette, opts MapPaletteOptions) ([][]int, error) {
	if err := palette.Validate(); err != nil {
		return nil, err
	}

	cache := NewColorCache()

	// Create closure for getting palette index
	getPaletteIndex := func(target Color) int {
		idx, err := FindClosestPaletteColor(target, palette, cache)
		if err != nil {
			return 0
		}
		return idx
	}

	// Create closure for getting palette color
	getPaletteColor := func(idx int) string {
		if idx >= 0 && idx < len(palette.Colors) {
			return palette.Colors[idx]
		}
		return "#000000"
	}

	// Apply dithering
	result := ApplyDithering(img, opts, getPaletteIndex, getPaletteColor)

	return result, nil
}

