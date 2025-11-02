package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/AndrewDonelson/retroforge-engine/internal/imgtool"
)

func ScaleCmd() *cobra.Command {
	var (
		outputPath     string
		width          int
		height         int
		algorithm     string
		ensureDiv      bool
		preserveAspect bool
	)

	cmd := &cobra.Command{
		Use:   "scale <input.png>",
		Short: "Scale image to target dimensions",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputPath := args[0]

			// Load image
			img, err := imgtool.LoadPNGFile(inputPath)
			if err != nil {
				return err
			}

			// Configure options
			opts := imgtool.ScaleOptions{
				Width:           width,
				Height:          height,
				Algorithm:       algorithm,
				EnsureDivisible: ensureDiv,
				PreserveAspect:  preserveAspect,
			}

			// Scale
			rgbData, err := imgtool.Scale(img, opts)
			if err != nil {
				return err
			}

			// Output JSON
			return outputRGBJSON(outputPath, rgbData)
		},
	}

	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output file (default: stdout)")
	cmd.Flags().IntVarP(&width, "width", "w", 16, "Target width")
	cmd.Flags().IntVarP(&height, "height", "h", 16, "Target height")
	cmd.Flags().StringVarP(&algorithm, "algorithm", "a", "nearest", "Algorithm: nearest, bilinear, bicubic")
	cmd.Flags().BoolVar(&ensureDiv, "ensure-divisible", true, "Ensure dimensions divisible by 2")
	cmd.Flags().BoolVar(&preserveAspect, "preserve-aspect", false, "Maintain aspect ratio")

	return cmd
}

func outputRGBJSON(outputPath string, rgbData [][][]uint8) error {
	data, err := json.MarshalIndent(rgbData, "", "  ")
	if err != nil {
		return err
	}

	if outputPath == "" {
		fmt.Println(string(data))
		return nil
	}

	return os.WriteFile(outputPath, data, 0644)
}

