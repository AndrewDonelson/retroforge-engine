package imgtool

import "image/color"

// IsTransparent checks if a pixel should be treated as transparent based on alpha threshold
func IsTransparent(c color.Color, threshold uint8) bool {
	_, _, _, a := c.RGBA()
	alpha := uint8(a >> 8)
	return alpha < threshold
}

// MakeOpaque converts a color to fully opaque if it's above the threshold
func MakeOpaque(c color.Color, threshold uint8) color.RGBA {
	r, g, b, a := c.RGBA()
	alpha := uint8(a >> 8)
	
	if alpha < threshold {
		return color.RGBA{R: 0, G: 0, B: 0, A: 0} // Fully transparent
	}
	
	return color.RGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: 255, // Fully opaque
	}
}

