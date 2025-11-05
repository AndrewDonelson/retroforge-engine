package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/AndrewDonelson/retroforge-engine/cmd/tile2iso/commands"
)

var version = "0.1.0"

func main() {
	rootCmd := &cobra.Command{
		Use:     "tile2iso",
		Short:   "RetroForge Isometric Tile Converter",
		Long:    "Convert 2D top-down sprites into isometric (2.5D) tiles by combining three textures: a top face, left side face, and right side face.",
		Version: version,
	}

	rootCmd.AddCommand(commands.ConvertCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

