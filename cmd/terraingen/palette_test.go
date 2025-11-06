package main

import (
	"fmt"
	"image/color"
)

// This file generates the exact RetroForge 50 palette colors
// It matches the palette generation logic in internal/pal/palettes.go

func hexToRGBA(hex string) color.RGBA {
	if len(hex) > 0 && hex[0] == '#' {
		hex = hex[1:]
	}
	if len(hex) != 6 {
		return color.RGBA{0, 0, 0, 255}
	}
	var r, g, b uint8
	fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return color.RGBA{r, g, b, 255}
}

func shade(hex string, amount int) string {
	c := hexToRGBA(hex)
	r := int(c.R) + amount
	g := int(c.G) + amount
	b := int(c.B) + amount
	if r < 0 {
		r = 0
	}
	if r > 255 {
		r = 255
	}
	if g < 0 {
		g = 0
	}
	if g > 255 {
		g = 255
	}
	if b < 0 {
		b = 0
	}
	if b > 255 {
		b = 255
	}
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

func generateRetroForge50Palette() []color.RGBA {
	huesRetroForge := []string{"#ff4d4d", "#ff914d", "#ffd84d", "#b6ff4d", "#4dd487", "#36d8c7", "#4dd5ff", "#66bfff", "#6f88ff", "#8a75ff", "#b478ff", "#ff6fb1", "#ff7fa0", "#a8795a", "#a0b15a", "#38bdf8"}
	
	colors := make([]color.RGBA, 50)
	colors[0] = hexToRGBA("#000000") // Black
	colors[1] = hexToRGBA("#ffffff") // White
	
	idx := 2
	for i := 0; i < 16 && idx < 50; i++ {
		base := huesRetroForge[i%len(huesRetroForge)]
		colors[idx] = hexToRGBA(shade(base, 60))  // highlight
		idx++
		if idx >= 50 {
			break
		}
		colors[idx] = hexToRGBA(base)              // base
		idx++
		if idx >= 50 {
			break
		}
		colors[idx] = hexToRGBA(shade(base, -60))  // shadow
		idx++
	}
	
	// Fill remaining with grayscale if needed
	for idx < 50 {
		v := uint8((idx * 5) & 255)
		colors[idx] = color.RGBA{v, v, v, 255}
		idx++
	}
	
	return colors
}

