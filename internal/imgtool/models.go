package imgtool

import (
	"encoding/json"
	"fmt"
	"image/color"
)

// Palette represents a 48-color RetroForge game palette
// Note: Built-in colors (0-15) are always available in the engine
// This structure represents only the game palette (48 colors, indices 16-63)
type Palette struct {
	Colors []string `json:"colors"` // Exactly 48 hex color strings
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

// Validate validates the palette structure (48 colors for game palette)
func (p *Palette) Validate() error {
	if len(p.Colors) != 48 {
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
	// Validate size using new rules (minimum 2x2, gameplay vs UI restrictions)
	if err := s.ValidateSize(); err != nil {
		return err
	}
	
	if len(s.Pixels) != s.Height {
		return ErrInvalidPixelData
	}
	for i, row := range s.Pixels {
		if len(row) != s.Width {
			return fmt.Errorf("row %d has wrong width: expected %d, got %d", i, s.Width, len(row))
		}
		for j, idx := range row {
			if idx < -1 || idx > 63 {
				return fmt.Errorf("invalid palette index at [%d][%d]: %d (valid range: -1 to 63)", i, j, idx)
			}
		}
	}
	return nil
}

// ValidateSize validates sprite dimensions based on new size rules
func (s *Sprite) ValidateSize() error {
	// Minimum size: 2x2
	if s.Width < 2 || s.Height < 2 {
		return fmt.Errorf("sprite dimensions must be at least 2x2, got %dx%d", s.Width, s.Height)
	}

	if s.IsUI {
		// UI sprites: 2-256, both dimensions divisible by 2
		if s.Width > 256 || s.Height > 256 {
			return fmt.Errorf("UI sprite dimensions cannot exceed 256, got %dx%d", s.Width, s.Height)
		}
		if s.Width%2 != 0 || s.Height%2 != 0 {
			return fmt.Errorf("UI sprite dimensions must be divisible by 2, got %dx%d", s.Width, s.Height)
		}
	} else {
		// Gameplay sprites: 2-32
		if s.Width > 32 || s.Height > 32 {
			return fmt.Errorf("gameplay sprite dimensions cannot exceed 32x32, got %dx%d", s.Width, s.Height)
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

