package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/AndrewDonelson/retroforge-engine/internal/imgtool"
)

func ToSpriteCmd() *cobra.Command {
	var (
		outputPath     string
		palettePath    string
		name           string
		width          int
		height         int
		useCollision   bool
		isUI           bool
		lifetime       int
		maxSpawn       int
		ditherAlg      string
		alphaThreshold int
	)

	cmd := &cobra.Command{
		Use:   "tosprite <input.png> <palette.json>",
		Short: "Convert PNG to sprite (complete pipeline)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputPath := args[0]
			palettePath = args[1]

			// Load image
			img, err := imgtool.LoadPNGFile(inputPath)
			if err != nil {
				return err
			}

			// Load palette
			palette, err := imgtool.LoadPaletteFile(palettePath)
			if err != nil {
				return err
			}

			// Configure options
			opts := imgtool.ToSpriteOptions{
				Name:            name,
				TargetWidth:     width,
				TargetHeight:    height,
				UseCollision:    useCollision,
				IsUI:            isUI,
				Lifetime:        lifetime,
				MaxSpawn:        maxSpawn,
				DitherAlgorithm: ditherAlg,
				AlphaThreshold:  uint8(alphaThreshold),
			}

			// Convert to sprite
			sprite, err := imgtool.ToSprite(img, palette, opts)
			if err != nil {
				return err
			}

			// Auto-generate output path if not specified
			if outputPath == "sprite.json" {
				// Generate output filename: sprite-{basename}.json in same folder as input
				inputDir := filepath.Dir(inputPath)
				inputBase := filepath.Base(inputPath)
				inputName := strings.TrimSuffix(inputBase, filepath.Ext(inputBase))
				outputPath = filepath.Join(inputDir, fmt.Sprintf("sprite-%s.json", inputName))
			}

			// Output JSON
			return outputSpriteJSON(outputPath, sprite)
		},
	}

	cmd.Flags().StringVarP(&outputPath, "output", "o", "sprite.json", "Output file")
	cmd.Flags().StringVarP(&name, "name", "n", "sprite", "Sprite name")
	cmd.Flags().IntVarP(&width, "width", "w", 16, "Target width")
	cmd.Flags().IntVar(&height, "height", 16, "Target height (use --height, not -h to avoid conflict with --help)")
	cmd.Flags().BoolVar(&useCollision, "collision", false, "Enable physics collision")
	cmd.Flags().BoolVar(&isUI, "ui", false, "UI sprite (no physics)")
	cmd.Flags().IntVar(&lifetime, "lifetime", 0, "Auto-destroy after ms (0 = no limit)")
	cmd.Flags().IntVar(&maxSpawn, "max-spawn", 0, "Max simultaneous instances (0 = no limit)")
	cmd.Flags().StringVarP(&ditherAlg, "dither", "d", "floyd-steinberg", "Dithering: floyd-steinberg, ordered, none")
	cmd.Flags().IntVarP(&alphaThreshold, "alpha-threshold", "a", 128, "Alpha threshold 0-255")

	return cmd
}

func outputSpriteJSON(outputPath string, sprite *imgtool.Sprite) error {
	data, err := sprite.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal sprite: %w", err)
	}

	if outputPath == "" {
		fmt.Println(string(data))
		return nil
	}

	return os.WriteFile(outputPath, data, 0644)
}

