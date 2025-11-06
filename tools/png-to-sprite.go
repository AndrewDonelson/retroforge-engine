package main

import (
	"fmt"
	"image"
	_ "image/png"
	"os"
)

// Calculate distance between two colors (RGB)
func colorDistance(r1, g1, b1, r2, g2, b2 uint8) float64 {
	dr := float64(r1) - float64(r2)
	dg := float64(g1) - float64(g2)
	db := float64(b1) - float64(b2)
	return dr*dr + dg*dg + db*db
}

// Find closest palette color index
// NOTE: This tool uses the OLD RetroForge 50 palette structure (0-49).
// The new system uses 64 colors: 0-15 built-in, 16-63 game palette (48 colors).
// This tool should be updated to use the new 64-color system.
// Old structure: 0 = black, 1 = white, 2-49 = 16 hues × 3 shades
// New structure: 0-15 = built-in colors (black=0, white=7, etc.), 16-63 = game palette
func findClosestPaletteIndex(r, g, b uint8) int {
	// Define RetroForge 50 palette colors (OLD SYSTEM - matches old palettes.go)
	paletteColors := []struct {
		r, g, b uint8
		index   int
	}{
		{0, 0, 0, 0},     // Black
		{255, 255, 255, 1}, // White
		// Red family (hue 0: #ff4d4d)
		{255, 137, 137, 2}, // highlight (shade +60)
		{255, 77, 77, 3},   // base (#ff4d4d)
		{195, 17, 17, 4},   // shadow (shade -60)
		// Orange family (hue 1: #ff914d)
		{255, 177, 125, 5},
		{255, 145, 77, 6},
		{195, 73, 13, 7},
		// Yellow family (hue 2: #ffd84d)
		{255, 248, 125, 8},
		{255, 216, 77, 9},
		{195, 136, 13, 10},
		// Lime family (hue 3: #b6ff4d)
		{214, 255, 125, 11},
		{182, 255, 77, 12},
		{82, 195, 13, 13},
		// Green family (hue 4: #4dd487)
		{125, 255, 167, 14},
		{77, 212, 135, 15},
		{13, 132, 55, 16},
		// Cyan family (hue 5: #36d8c7)
		{108, 248, 227, 17},
		{54, 216, 199, 18},
		{0, 156, 139, 19},
		// Sky blue family (hue 6: #4dd5ff)
		{125, 255, 255, 20},
		{77, 213, 255, 21},
		{13, 153, 195, 22},
		// Blue family (hue 7: #66bfff)
		{150, 223, 255, 23},
		{102, 191, 255, 24},
		{42, 131, 195, 25},
		// Indigo family (hue 8: #6f88ff)
		{159, 184, 255, 26},
		{111, 136, 255, 27},
		{39, 88, 195, 28},
		// Purple family (hue 9: #8a75ff)
		{186, 181, 255, 29},
		{138, 117, 255, 30},
		{58, 57, 195, 31},
		// Violet family (hue 10: #b478ff)
		{212, 200, 255, 32},
		{180, 120, 255, 33},
		{88, 24, 195, 34},
		// Pink family (hue 11: #ff6fb1)
		{255, 191, 217, 35},
		{255, 111, 177, 36},
		{195, 51, 81, 37},
		// Rose family (hue 12: #ff7fa0)
		{255, 191, 216, 38},
		{255, 127, 160, 39},
		{195, 67, 64, 40},
		// Brown family (hue 13: #a8795a)
		{232, 185, 146, 41},
		{168, 121, 90, 42},
		{72, 57, 34, 43},
		// Olive family (hue 14: #a0b15a)
		{224, 225, 146, 44},
		{160, 177, 90, 45},
		{64, 113, 34, 46},
		// Cyan blue family (hue 15: #38bdf8)
		{152, 253, 255, 47},
		{56, 189, 248, 48},
		{0, 129, 188, 49},
	}

	minDist := colorDistance(r, g, b, paletteColors[0].r, paletteColors[0].g, paletteColors[0].b)
	closest := 0

	for _, pc := range paletteColors {
		dist := colorDistance(r, g, b, pc.r, pc.g, pc.b)
		if dist < minDist {
			minDist = dist
			closest = pc.index
		}
	}

	return closest
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <logo.png> [output_width] [output_height]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Converts a PNG image to RetroForge sprite data (48x48 by default)\n")
		os.Exit(1)
	}

	pngPath := os.Args[1]
	outputWidth := 48
	outputHeight := 48

	if len(os.Args) >= 3 {
		fmt.Sscanf(os.Args[2], "%d", &outputWidth)
	}
	if len(os.Args) >= 4 {
		fmt.Sscanf(os.Args[3], "%d", &outputHeight)
	}

	// Open and decode PNG
	file, err := os.Open(pngPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening PNG: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding PNG: %v\n", err)
		os.Exit(1)
	}

	bounds := img.Bounds()
	origWidth := bounds.Dx()
	origHeight := bounds.Dy()

	// Extract only the left icon portion (first ~56 pixels width for the 48x48 icon)
	// The icon is the leftmost portion of the logo
	// SVG shows icon at x=4, width=48, so we need to account for viewBox scaling
	// If PNG is 140x28, icon is roughly first 56 pixels (140/280 * 56 = 28, but viewBox scaling)
	// Actually, if the PNG is 2x the SVG (common), icon would be 56px wide
	// Use native PNG dimensions - extract square icon from left side
	iconWidth := origHeight // Use height as icon size (should be square)
	if iconWidth > origWidth {
		iconWidth = origWidth // But don't exceed image width
	}
	
	// For RetroForge logo: icon should be roughly 48-56px wide in a 140x28 or 280x56 PNG
	// Check if this looks like a wide logo (width >> height)
	if origWidth > origHeight*2 {
		// Wide logo format - extract left square portion
		iconWidth = origHeight
	}
	
	iconHeight := iconWidth // Icon is square
	cropX := 0
	cropY := 0

	// Use native icon size (no scaling)
	outputWidth = iconWidth
	outputHeight = iconHeight

	fmt.Printf("// PNG image: %dx%d, extracting icon region (%d,%d) size %dx%d (no scaling)\n", 
		origWidth, origHeight, cropX, cropY, iconWidth, iconHeight)
	fmt.Printf("// logoPixels contains the hardcoded RetroForge logo sprite data (%dx%d)\n", outputWidth, outputHeight)
	fmt.Printf("// Colors mapped to RetroForge 50 palette indices (closest match)\n")
	fmt.Printf("// NOTE: This uses the OLD 50-color system. New system uses 64 colors (16 built-in + 48 game)\n")
	fmt.Printf("var logoPixels = [][]int{\n")

	// Create sprite data at native resolution
	pixels := make([][]int, outputHeight)
	for y := 0; y < outputHeight; y++ {
		pixels[y] = make([]int, outputWidth)
		rowStr := "\t{"

		for x := 0; x < outputWidth; x++ {
			// Use direct pixel coordinates (no scaling)
			srcX := cropX + x
			srcY := cropY + y
			
			// Clamp to image bounds
			if srcX >= origWidth {
				srcX = origWidth - 1
			}
			if srcY >= origHeight {
				srcY = origHeight - 1
			}

			// Get pixel color at source coordinates
			r, g, b, a := img.At(srcX, srcY).RGBA()
			// Convert from 16-bit to 8-bit
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			a8 := uint8(a >> 8)

			var colorIndex int
			if a8 < 128 {
				// Transparent pixel
				colorIndex = -1
			} else {
				// Find closest palette color
				colorIndex = findClosestPaletteIndex(r8, g8, b8)
			}

			pixels[y][x] = colorIndex
			if x > 0 {
				rowStr += ","
			}
			rowStr += fmt.Sprintf("%3d", colorIndex)
		}

		rowStr += "},"
		fmt.Println(rowStr)
	}

	fmt.Printf("}\n")
}
