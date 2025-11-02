package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/AndrewDonelson/retroforge-engine/internal/imgtool"
)

func MapPaletteCmd() *cobra.Command {
	var (
		outputPath     string
		palettePath    string
		ditherAlg      string
		alphaThreshold int
	)

	cmd := &cobra.Command{
		Use:   "mappalette <input.png> <palette.json>",
		Short: "Map image to palette indices",
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
			opts := imgtool.MapPaletteOptions{
				DitherAlgorithm:  ditherAlg,
				AlphaThreshold:   uint8(alphaThreshold),
				TransparentIndex: -1,
			}

			// Map palette
			indices, err := imgtool.MapPalette(img, palette, opts)
			if err != nil {
				return err
			}

			// Output JSON
			return outputIndicesJSON(outputPath, indices)
		},
	}

	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output file (default: stdout)")
	cmd.Flags().StringVarP(&ditherAlg, "dither", "d", "floyd-steinberg", "Dithering: floyd-steinberg, ordered, none")
	cmd.Flags().IntVarP(&alphaThreshold, "alpha-threshold", "a", 128, "Alpha threshold 0-255")

	return cmd
}

func outputIndicesJSON(outputPath string, indices [][]int) error {
	data, err := json.MarshalIndent(indices, "", "  ")
	if err != nil {
		return err
	}

	if outputPath == "" {
		fmt.Println(string(data))
		return nil
	}

	return os.WriteFile(outputPath, data, 0644)
}

