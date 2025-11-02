package imgtool

import (
	"encoding/json"
	"fmt"
	"image/color"
)

// Palette represents a 50-color RetroForge palette
type Palette struct {
	Colors []string `json:"colors"` // Exactly 50 hex color strings
}

// Sprite represents a RetroForge sprite (matches sprites.json format)
type Sprite struct {
	Width        int       `json:"width"`
	Height       int       `json:"height"`
	Pixels       [][]int   `json:"pixels"` // 2D array of palette indices
	UseCollision bool      `json:"useCollision"`
	IsUI         bool      `json:"isUI"`
	Lifetime     int       `json:"lifetime"`
	MaxSpawn     int       `json:"maxSpawn"`
	MountPoints  []Point   `json:"mountPoints"`
}

// Point represents a mount point coordinate
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Color represents an RGB color (internal use)
type Color struct {
	R, G, B uint8
}

// Validate validates the palette structure
func (p *Palette) Validate() error {
	if len(p.Colors) != 50 {
		return ErrInvalidPaletteSize
	}
	for i, hex := range p.Colors {
		if !isValidHex(hex) {
			return fmt.Errorf("invalid hex color at index %d: %s", i, hex)
		}
	}
	return nil
}

// Validate validates the sprite structure
func (s *Sprite) Validate() error {
	if s.Width <= 0 || s.Height <= 0 {
		return ErrInvalidDimensions
	}
	if len(s.Pixels) != s.Height {
		return ErrInvalidPixelData
	}
	for i, row := range s.Pixels {
		if len(row) != s.Width {
			return fmt.Errorf("row %d has wrong width: expected %d, got %d", i, s.Width, len(row))
		}
		for j, idx := range row {
			if idx < -1 || idx > 49 {
				return fmt.Errorf("invalid palette index at [%d][%d]: %d", i, j, idx)
			}
		}
	}
	return nil
}

// ToHex converts Color to hex string
func (c Color) ToHex() string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

// HexToColor parses a hex color string to Color
func HexToColor(hex string) (Color, error) {
	if len(hex) != 7 || hex[0] != '#' {
		return Color{}, ErrInvalidHexFormat
	}

	var r, g, b uint8
	_, err := fmt.Sscanf(hex[1:3], "%02x", &r)
	if err != nil {
		return Color{}, ErrInvalidHexFormat
	}
	_, err = fmt.Sscanf(hex[3:5], "%02x", &g)
	if err != nil {
		return Color{}, ErrInvalidHexFormat
	}
	_, err = fmt.Sscanf(hex[5:7], "%02x", &b)
	if err != nil {
		return Color{}, ErrInvalidHexFormat
	}

	return Color{R: r, G: g, B: b}, nil
}

// ToRGBA converts Color to color.RGBA
func (c Color) ToRGBA() color.RGBA {
	return color.RGBA{R: c.R, G: c.G, B: c.B, A: 255}
}

// isValidHex checks if a string is a valid hex color
func isValidHex(hex string) bool {
	if len(hex) != 7 || hex[0] != '#' {
		return false
	}
	for i := 1; i < 7; i++ {
		c := hex[i]
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

// ToJSON converts Sprite to JSON bytes
func (s *Sprite) ToJSON() ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}

