package pal

import (
	"fmt"
	"image/color"
	"sort"
)

// Predefined palettes - all palettes share the same 50 color index structure:
// - Index 0: black
// - Index 1: white
// - Indices 2-49: 16 hues × 3 shades each (highlight, base, shadow)

func hexToRGBA(hex string) color.RGBA {
	// Remove # if present
	if len(hex) > 0 && hex[0] == '#' {
		hex = hex[1:]
	}
	if len(hex) != 6 {
		return color.RGBA{0, 0, 0, 255}
	}
	var r, g, b uint8
	_, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return color.RGBA{0, 0, 0, 255}
	}
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

func generate50Palette(baseHues16 []string) []color.RGBA {
	colors := make([]color.RGBA, 50)
	colors[0] = hexToRGBA("#000000") // Black
	colors[1] = hexToRGBA("#ffffff") // White
	
	idx := 2
	for i := 0; i < 16 && idx < 50; i++ {
		base := baseHues16[i%len(baseHues16)]
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

// generateGrayscale50Palette creates a true grayscale palette with smooth gradient from black to white
func generateGrayscale50Palette() []color.RGBA {
	colors := make([]color.RGBA, 50)
	for i := 0; i < 50; i++ {
		// Create smooth gradient from black (0) to white (255)
		// Map index 0-49 to 0-255 with smooth distribution
		v := uint8((i * 255) / 49)
		colors[i] = color.RGBA{v, v, v, 255}
	}
	// Ensure index 0 is pure black and index 1 is pure white (or near-white)
	colors[0] = color.RGBA{0, 0, 0, 255}
	if len(colors) > 1 {
		colors[1] = color.RGBA{255, 255, 255, 255}
	}
	return colors
}

// Base hue sets (16 each) - same as webapp
var huesRetroForge = []string{"#ff4d4d", "#ff914d", "#ffd84d", "#b6ff4d", "#4dd487", "#36d8c7", "#4dd5ff", "#66bfff", "#6f88ff", "#8a75ff", "#b478ff", "#ff6fb1", "#ff7fa0", "#a8795a", "#a0b15a", "#38bdf8"}
var huesPICO8 = []string{"#ff004d", "#ffa300", "#ffec27", "#00e436", "#29adff", "#83769c", "#ff77a8", "#ffccaa", "#1d2b53", "#7e2553", "#008751", "#ab5236", "#5f574f", "#c2c3c7", "#fff1e8", "#000000"}
var huesNeon = []string{"#ff006e", "#ff5400", "#ffbd00", "#a7ff00", "#00f5d4", "#00bbf9", "#4361ee", "#7209b7", "#b5179e", "#f72585", "#3eff99", "#f8961e", "#ffd166", "#06d6a0", "#118ab2", "#ef476f"}
var huesPastel = []string{"#f4a3a3", "#f7c59f", "#f7e8a3", "#cde7b0", "#a3e0dc", "#a3c4f3", "#c9b6e4", "#f7aef8", "#ffd6e0", "#cdb4db", "#ffc8dd", "#ffafcc", "#bde0fe", "#a2d2ff", "#d0f4de", "#f1faee"}
var huesEarth = []string{"#b5651d", "#cb997e", "#ddbea9", "#ffe8d6", "#6b705c", "#a5a58d", "#b7b7a4", "#3d405b", "#81b29a", "#f2cc8f", "#e07a5f", "#8d5524", "#c68642", "#5d4037", "#7d5a50", "#a47148"}
var huesWarcraft = []string{"#7a2e1b", "#c79347", "#d7b26d", "#e4d39b", "#356b2b", "#5a8e3b", "#88a84a", "#2e4057", "#4a6fa1", "#7aa2d6", "#9765a8", "#5a3668", "#a33e3e", "#6b2b2b", "#a38a6b", "#3f2f1b"}
var huesStarCraft = []string{"#2b3a67", "#3b82f6", "#60a5fa", "#93c5fd", "#22c55e", "#16a34a", "#0d9488", "#06b6d4", "#14b8a6", "#a855f7", "#7c3aed", "#f97316", "#ea580c", "#f59e0b", "#eab308", "#64748b"}
var huesSuperMario = []string{"#e52521", "#ff7f27", "#ffbd3a", "#ffe761", "#3cb44b", "#0f7f12", "#00a2e8", "#3f48cc", "#1d2bd7", "#a349a4", "#c51162", "#ff66a1", "#9b7653", "#8b4513", "#ffd1dc", "#ffb347"}
var huesGrayscale = []string{"#0a0a0a", "#1a1a1a", "#2a2a2a", "#3a3a3a", "#4a4a4a", "#5a5a5a", "#6a6a6a", "#7a7a7a", "#8a8a8a", "#9a9a9a", "#aaaaaa", "#bababa", "#cacaca", "#dadada", "#eaeaea", "#f5f5f5"}
var huesNES = []string{"#7c7c7c", "#0000fc", "#0000bc", "#4428bc", "#940084", "#a80020", "#a81000", "#881400", "#503000", "#007800", "#006800", "#005800", "#004058", "#000000", "#bcbcbc", "#0078f8"}
var huesSNES = []string{"#e04048", "#f0a000", "#f8e060", "#60d0a8", "#40a0e0", "#6060e0", "#a860e0", "#e060a8", "#b03030", "#c07030", "#c0a040", "#60a060", "#4080b0", "#6060a0", "#9060a0", "#a06080"}
var huesGenesis = []string{"#e03a3a", "#ff7f00", "#ffd400", "#a4de02", "#00a884", "#00a3ff", "#0051ff", "#6a00ff", "#b100e8", "#ff00a8", "#ff4d6d", "#ff9e00", "#ffd166", "#06d6a0", "#118ab2", "#073b4c"}
var huesAmiga = []string{"#cc0000", "#ff8c00", "#ffd700", "#9acd32", "#2e8b57", "#20b2aa", "#1e90ff", "#4169e1", "#6a5acd", "#9932cc", "#c71585", "#ff69b4", "#cd853f", "#8b4513", "#708090", "#2f4f4f"}
var huesGameBoyColor = []string{"#0b380f", "#306230", "#8bac0f", "#9bbc0f", "#1b2b34", "#343d46", "#4f5b66", "#65737e", "#a7adba", "#c0c5ce", "#6699cc", "#99c794", "#5fb3b3", "#fac863", "#ec5f67", "#ab7967"}
var huesCyberpunk = []string{"#ff006e", "#f9c80e", "#00f5d4", "#00bbf9", "#3a0ca3", "#7209b7", "#4361ee", "#4cc9f0", "#f72585", "#b5179e", "#4895ef", "#560bad", "#2b2d42", "#8d99ae", "#ef233c", "#ffd166"}
var huesMonokai = []string{"#f92672", "#fd971f", "#e6db74", "#a6e22e", "#66d9ef", "#ae81ff", "#f8f8f2", "#75715e", "#272822", "#1e1f29", "#ff6188", "#fc9867", "#ffd866", "#a9dc76", "#78dce8", "#ab9df2"}

// Predefined palettes map
var predefinedPalettes = map[string][]color.RGBA{
	"default":      Default50,
	"RetroForge 50": generate50Palette(huesRetroForge),
	"PICO-8+ 50":   generate50Palette(huesPICO8),
	"Neon 50":      generate50Palette(huesNeon),
	"Pastel 50":    generate50Palette(huesPastel),
	"Earth 50":     generate50Palette(huesEarth),
	"Warcraft 50":  generate50Palette(huesWarcraft),
	"StarCraft 50": generate50Palette(huesStarCraft),
	"Super Mario 50": generate50Palette(huesSuperMario),
	"Grayscale 50": generateGrayscale50Palette(),
	"NES 50":       generate50Palette(huesNES),
	"SNES 50":      generate50Palette(huesSNES),
	"Genesis 50":   generate50Palette(huesGenesis),
	"Amiga 50":     generate50Palette(huesAmiga),
	"Game Boy Color 50": generate50Palette(huesGameBoyColor),
	"Cyberpunk 50": generate50Palette(huesCyberpunk),
	"Monokai 50":   generate50Palette(huesMonokai),
}

// GetPaletteNames returns all available palette names
func GetPaletteNames() []string {
	names := make([]string, 0, len(predefinedPalettes))
	for name := range predefinedPalettes {
		names = append(names, name)
	}
	// Sort for consistent output
	sort.Strings(names)
	return names
}

// GetPalette returns a copy of the named palette, or default if not found
func GetPalette(name string) []color.RGBA {
	if pal, ok := predefinedPalettes[name]; ok {
		// Return a copy
		return append([]color.RGBA{}, pal...)
	}
	// Return default if not found
	return append([]color.RGBA{}, Default50...)
}
