package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/AndrewDonelson/retroforge-engine/internal/imgtool"
)

func QuantizeCmd() *cobra.Command {
	var (
		outputPath     string
		ditherAlg      string
		alphaThreshold int
		noBlackWhite   bool
	)

	cmd := &cobra.Command{
		Use:   "quantize <input.png>",
		Short: "Reduce image to 50-color palette",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputPath := args[0]

			// Load image using core package
			img, err := imgtool.LoadPNGFile(inputPath)
			if err != nil {
				return err
			}

			// Configure options
			opts := imgtool.QuantizeOptions{
				DitherAlgorithm:   ditherAlg,
				AlphaThreshold:    uint8(alphaThreshold),
				EnforceBlackWhite: !noBlackWhite,
			}

			// Call core package function
			palette, err := imgtool.Quantize(img, opts)
			if err != nil {
				return err
			}

			// Output JSON
			return outputPaletteJSON(outputPath, palette)
		},
	}

	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output file (default: stdout)")
	cmd.Flags().StringVarP(&ditherAlg, "dither", "d", "floyd-steinberg", "Dithering: floyd-steinberg, ordered, none")
	cmd.Flags().IntVarP(&alphaThreshold, "alpha-threshold", "a", 128, "Alpha threshold 0-255")
	cmd.Flags().BoolVar(&noBlackWhite, "no-black-white", false, "Don't enforce black/white at indices 0/1")

	return cmd
}

func outputPaletteJSON(outputPath string, palette *imgtool.Palette) error {
	data, err := json.MarshalIndent(palette, "", "  ")
	if err != nil {
		return err
	}

	if outputPath == "" {
		fmt.Println(string(data))
		return nil
	}

	return os.WriteFile(outputPath, data, 0644)
}

