package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Converts palettes.json format to imgtool palette format
func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: go run convert-palette-for-imgtool.go <palettes.json> <palette-name> <output.json>\n")
		fmt.Fprintf(os.Stderr, "Example: go run convert-palette-for-imgtool.go ../../retroforge-webapp/palettes.json \"RetroForge 50\" palette.json\n")
		os.Exit(1)
	}

	palettesPath := os.Args[1]
	paletteName := os.Args[2]
	outputPath := os.Args[3]

	// Read palettes.json
	data, err := os.ReadFile(palettesPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", palettesPath, err)
		os.Exit(1)
	}

	var root struct {
		Palettes map[string][][]int `json:"palettes"`
	}
	if err := json.Unmarshal(data, &root); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	// Find the palette
	rgbColors, exists := root.Palettes[paletteName]
	if !exists {
		fmt.Fprintf(os.Stderr, "Palette '%s' not found. Available palettes:\n", paletteName)
		for name := range root.Palettes {
			fmt.Fprintf(os.Stderr, "  - %s\n", name)
		}
		os.Exit(1)
	}

	if len(rgbColors) != 50 {
		fmt.Fprintf(os.Stderr, "Error: Palette must have exactly 50 colors, found %d\n", len(rgbColors))
		os.Exit(1)
	}

	// Convert RGB arrays to hex strings
	hexColors := make([]string, 50)
	for i, rgb := range rgbColors {
		if len(rgb) != 3 {
			fmt.Fprintf(os.Stderr, "Error: Color %d must have 3 components [R,G,B]\n", i)
			os.Exit(1)
		}
		r, g, b := rgb[0], rgb[1], rgb[2]
		hexColors[i] = fmt.Sprintf("#%02x%02x%02x", r, g, b)
	}

	// Create output format
	output := map[string]interface{}{
		"colors": hexColors,
	}

	// Write output
	outputData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(outputPath, outputData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", outputPath, err)
		os.Exit(1)
	}

	fmt.Printf("Converted palette '%s' to %s\n", paletteName, outputPath)
}

