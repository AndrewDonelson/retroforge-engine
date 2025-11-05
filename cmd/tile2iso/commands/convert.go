package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
	"github.com/AndrewDonelson/retroforge-engine/internal/imgtool"
	"github.com/AndrewDonelson/retroforge-engine/internal/tile2iso"
)

func ConvertCmd() *cobra.Command {
	var (
		spritesPath   string
		palettePath   string
		outputPath    string
		topSprite     string
		leftSprite    string
		rightSprite   string
		topFrame      string
		leftFrame     string
		rightFrame    string
		tileName      string
		height        int
		lightingMode  string
		tileWidth     int
		tileHeight    int
		showOutline   bool
		updateSprites bool
	)

	cmd := &cobra.Command{
		Use:   "convert",
		Short: "Convert three sprites into an isometric tile",
		Long: `Convert three sprites (top, left, right) from sprites.json into a single isometric tile.
The tool supports static, frames, and animation sprite types. For frames/animation sprites,
you can specify which frame to use with --top-frame, --left-frame, --right-frame flags.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate required flags
			if spritesPath == "" {
				return fmt.Errorf("--sprites is required")
			}
			if palettePath == "" {
				return fmt.Errorf("--palette is required")
			}
			if topSprite == "" {
				return fmt.Errorf("--top is required")
			}
			if leftSprite == "" {
				return fmt.Errorf("--left is required")
			}
			if rightSprite == "" {
				return fmt.Errorf("--right is required")
			}

			// Load sprites
			spriteMap, err := tile2iso.LoadSpritesFromFile(spritesPath)
			if err != nil {
				return fmt.Errorf("failed to load sprites: %w", err)
			}

			// Load palette
			palette, err := imgtool.LoadPaletteFile(palettePath)
			if err != nil {
				return fmt.Errorf("failed to load palette: %w", err)
			}

			// Parse lighting mode
			mode := tile2iso.LightingMode(lightingMode)
			if mode == "" {
				mode = tile2iso.LightingGradient
			}
			if mode != tile2iso.LightingNormal && mode != tile2iso.LightingBasic &&
				mode != tile2iso.LightingFull && mode != tile2iso.LightingGradient {
				return fmt.Errorf("invalid lighting mode '%s', must be one of: normal, basic, full, gradient", lightingMode)
			}

			// Configure options
			options := tile2iso.TileOptions{
				Height:       height,
				LightingMode: mode,
				TileWidth:    tileWidth,
				TileHeight:   tileHeight,
				ShowOutline:  showOutline,
			}

			// Create converter with default dimensions (32×24 output tiles)
			converter := tile2iso.NewIsometricConverter(32, 16)

			// Generate isometric tile
			resultSprite, err := converter.CreateIsometricTile(
				topSprite, leftSprite, rightSprite,
				topFrame, leftFrame, rightFrame,
				palette.Colors,
				spriteMap,
				options,
			)
			if err != nil {
				return fmt.Errorf("failed to create isometric tile: %w", err)
			}

			// Determine output name
			if tileName == "" {
				tileName = fmt.Sprintf("%s_%s_%s_iso", topSprite, leftSprite, rightSprite)
			}

			// If updating sprites.json, add to existing map
			if updateSprites {
				spriteMap[tileName] = *resultSprite

				// Write updated sprites.json
				data, err := json.MarshalIndent(spriteMap, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal sprites: %w", err)
				}

				if err := os.WriteFile(spritesPath, data, 0644); err != nil {
					return fmt.Errorf("failed to write sprites.json: %w", err)
				}

				fmt.Printf("✓ Added isometric tile '%s' to %s\n", tileName, spritesPath)
				return nil
			}

			// Otherwise, write to output file or stdout
			if outputPath == "" {
				// Auto-generate output path
				spritesDir := filepath.Dir(spritesPath)
				outputPath = filepath.Join(spritesDir, fmt.Sprintf("tile-%s.json", tileName))
			}

			// Create output with single sprite
			outputMap := cartio.SpriteMap{
				tileName: *resultSprite,
			}

			data, err := json.MarshalIndent(outputMap, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal sprite: %w", err)
			}

			if err := os.WriteFile(outputPath, data, 0644); err != nil {
				return fmt.Errorf("failed to write output file: %w", err)
			}

			fmt.Printf("✓ Created isometric tile '%s' at %s\n", tileName, outputPath)
			return nil
		},
	}

	cmd.Flags().StringVar(&spritesPath, "sprites", "", "Path to sprites.json file (required)")
	cmd.Flags().StringVar(&palettePath, "palette", "", "Path to palette.json file (required)")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output file path (default: auto-generated in same dir as sprites.json)")
	cmd.Flags().StringVar(&topSprite, "top", "", "Name of top sprite (required)")
	cmd.Flags().StringVar(&leftSprite, "left", "", "Name of left side sprite (required)")
	cmd.Flags().StringVar(&rightSprite, "right", "", "Name of right side sprite (required)")
	cmd.Flags().StringVar(&topFrame, "top-frame", "", "Frame name for top sprite (for frames/animation types)")
	cmd.Flags().StringVar(&leftFrame, "left-frame", "", "Frame name for left sprite (for frames/animation types)")
	cmd.Flags().StringVar(&rightFrame, "right-frame", "", "Frame name for right sprite (for frames/animation types)")
	cmd.Flags().StringVar(&tileName, "name", "", "Name for output tile sprite (default: {top}_{left}_{right}_iso)")
	cmd.Flags().IntVar(&height, "height", 8, "Height of side faces in pixels (default: 8 for 32×24 tiles)")
	cmd.Flags().StringVar(&lightingMode, "lighting", "gradient", "Lighting mode: normal, basic, full, gradient")
	cmd.Flags().IntVar(&tileWidth, "tile-width", 32, "Width of isometric tile (default: 32 for 32×24 tiles)")
	cmd.Flags().IntVar(&tileHeight, "tile-height", 16, "Height of isometric tile top face (default: 16 for 32×24 tiles)")
	cmd.Flags().BoolVar(&showOutline, "show-outline", false, "Draw dark outlines around tile faces (top diamond and side parallelograms)")
	cmd.Flags().BoolVar(&updateSprites, "update", false, "Update sprites.json with new tile (instead of creating separate file)")

	return cmd
}

