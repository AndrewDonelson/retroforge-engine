package commands

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/AndrewDonelson/retroforge-engine/internal/imgtool"
)

func ToRPICommand() *cobra.Command {
	var (
		outputPath     string
		palettePath    string
		ditherAlg      string
		alphaThreshold int
		landscape      bool
	)

	cmd := &cobra.Command{
		Use:   "torpi <input.png> <palette.json>",
		Short: "Convert PNG to .rpi (Raw Palette Indexed) format",
		Long: `Convert PNG image to .rpi format (optimized for screens/backgrounds).

.rpi format:
- Always 480x270 (landscape) or 270x480 (portrait)
- Uses 6-bit palette indices (0-49 for colors, 63 for transparent)
- Maximum compression with gzip
- Optimized for large screens/backgrounds

The image will be scaled to fit the target dimensions while preserving aspect ratio.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputPath := args[0]
			palettePath = args[1]

			// Load image
			img, err := imgtool.LoadPNGFile(inputPath)
			if err != nil {
				return fmt.Errorf("failed to load image: %w", err)
			}

			// Load palette
			palette, err := imgtool.LoadPaletteFile(palettePath)
			if err != nil {
				return fmt.Errorf("failed to load palette: %w", err)
			}

			// Determine dimensions based on image aspect ratio
			imgBounds := img.Bounds()
			imgWidth := imgBounds.Dx()
			imgHeight := imgBounds.Dy()
			imgAspect := float64(imgWidth) / float64(imgHeight)

			// Default to landscape (480x270), unless explicitly portrait or image is taller
			targetWidth := 480
			targetHeight := 270
			if landscape || (!landscape && imgAspect < 1.0) {
				// Portrait mode: 270x480
				targetWidth = 270
				targetHeight = 480
			}

			// Configure sprite conversion options
			opts := imgtool.ToSpriteOptions{
				Name:            "splash",
				TargetWidth:     targetWidth,
				TargetHeight:    targetHeight,
				UseCollision:    false,
				IsUI:            true, // Screens are UI elements
				Lifetime:        0,
				MaxSpawn:        0,
				DitherAlgorithm: ditherAlg,
				AlphaThreshold:  uint8(alphaThreshold),
			}

			// Convert to sprite
			sprite, err := imgtool.ToSprite(img, palette, opts)
			if err != nil {
				return fmt.Errorf("failed to convert to sprite: %w", err)
			}

			// Convert sprite to RPI format
			rpiData, err := spriteToRPI(sprite)
			if err != nil {
				return fmt.Errorf("failed to create RPI: %w", err)
			}

			// Auto-generate output path if not specified
			if outputPath == "" {
				inputDir := filepath.Dir(inputPath)
				inputBase := filepath.Base(inputPath)
				inputName := strings.TrimSuffix(inputBase, filepath.Ext(inputBase))
				outputPath = filepath.Join(inputDir, inputName+".rpi")
			}

			// Write RPI file
			if err := os.WriteFile(outputPath, rpiData, 0644); err != nil {
				return fmt.Errorf("failed to write RPI file: %w", err)
			}

			fmt.Printf("✓ Created %s (%dx%d, %d bytes)\n", outputPath, sprite.Width, sprite.Height, len(rpiData))
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output .rpi file path (default: same as input with .rpi extension)")
	cmd.Flags().StringVarP(&ditherAlg, "dither", "d", "floyd-steinberg", "Dithering algorithm: floyd-steinberg, ordered, none")
	cmd.Flags().IntVarP(&alphaThreshold, "alpha-threshold", "a", 128, "Alpha threshold 0-255")
	cmd.Flags().BoolVar(&landscape, "landscape", false, "Force landscape (480x270), default is auto-detect from aspect ratio")

	return cmd
}

// spriteToRPI converts a Sprite to .rpi format (gzip-compressed binary)
func spriteToRPI(sprite *imgtool.Sprite) ([]byte, error) {
	width := sprite.Width
	height := sprite.Height

	// Pack header (8 bytes)
	// flags: bit 0 = landscape (0) / portrait (1)
	flags := uint16(0)
	if height > width {
		flags |= 1 // Portrait mode
	}

	header := make([]byte, 8)
	binary.LittleEndian.PutUint16(header[0:2], uint16(width))
	binary.LittleEndian.PutUint16(header[2:4], uint16(height))
	binary.LittleEndian.PutUint16(header[4:6], flags)
	binary.LittleEndian.PutUint16(header[6:8], 0) // Reserved

	// Pack pixel data as 6-bit values
	// Map: -1 (transparent) -> 63, 0-49 -> 0-49
	totalPixels := width * height
	encodedPixels := make([]byte, totalPixels)

	pixelIdx := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := sprite.Pixels[y][x]
			// Map transparent (-1) to 63, colors (0-49) stay as-is
			if pixel == -1 {
				encodedPixels[pixelIdx] = 63
			} else if pixel >= 0 && pixel <= 49 {
				encodedPixels[pixelIdx] = byte(pixel)
			} else {
				// Invalid palette index, default to transparent
				encodedPixels[pixelIdx] = 63
			}
			pixelIdx++
		}
	}

	// Pack 6-bit values into bytes
	// Pack 4 pixels into 3 bytes: 4 * 6 = 24 bits = 3 bytes
	packedData := make([]byte, (totalPixels+3)/4*3)
	for i := 0; i < len(encodedPixels); i += 4 {
		// Get up to 4 pixels
		var p [4]byte
		for j := 0; j < 4 && i+j < len(encodedPixels); j++ {
			p[j] = encodedPixels[i+j]
		}
		// Pad with 0 if needed
		for j := len(encodedPixels) - i; j < 4; j++ {
			p[j] = 0
		}

		// Pack 4 pixels (24 bits) into 3 bytes:
		// Pixel 0: bits 0-5 -> Byte 0 bits 0-5
		// Pixel 1: bits 0-5 -> Byte 0 bits 6-7, Byte 1 bits 0-3
		// Pixel 2: bits 0-5 -> Byte 1 bits 4-7, Byte 2 bits 0-1
		// Pixel 3: bits 0-5 -> Byte 2 bits 2-7
		byteIdx := (i / 4) * 3
		if byteIdx+2 < len(packedData) {
			packedData[byteIdx] = (p[0] & 0x3F) | ((p[1] & 0x03) << 6)
			packedData[byteIdx+1] = ((p[1] >> 2) & 0x0F) | ((p[2] & 0x0F) << 4)
			packedData[byteIdx+2] = ((p[2] >> 4) & 0x03) | ((p[3] & 0x3F) << 2)
		}
	}

	// Combine header + data
	rpiData := append(header, packedData...)

	// Compress with gzip (best compression)
	compressed := new(bytes.Buffer)
	gzWriter, err := gzip.NewWriterLevel(compressed, gzip.BestCompression)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip writer: %w", err)
	}
	if _, err := gzWriter.Write(rpiData); err != nil {
		return nil, fmt.Errorf("failed to compress: %w", err)
	}
	if err := gzWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close gzip writer: %w", err)
	}

	return compressed.Bytes(), nil
}

