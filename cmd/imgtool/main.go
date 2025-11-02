package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/AndrewDonelson/retroforge-engine/cmd/imgtool/commands"
)

var version = "0.1.0"

func main() {
	rootCmd := &cobra.Command{
		Use:   "imgtool",
		Short: "RetroForge Image Tool",
		Long:  "Convert PNG images to RetroForge sprites and palettes",
		Version: version,
	}

	rootCmd.AddCommand(commands.QuantizeCmd())
	rootCmd.AddCommand(commands.MapPaletteCmd())
	rootCmd.AddCommand(commands.ScaleCmd())
	rootCmd.AddCommand(commands.ToSpriteCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

