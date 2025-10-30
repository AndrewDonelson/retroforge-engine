package palette

import (
    "errors"
    "fmt"
    "strings"
)

// Palette represents a fixed-size RetroForge palette of 50 colors:
// - index 0: black
// - index 1: white
// - indices 2..49: 16 hues × 3 shades (highlight, base, shadow)
type Palette struct {
    Name   string
    Colors [50]string // hex colors like #RRGGBB (lower/upper case accepted)
}

var (
    ErrSize      = errors.New("palette must have exactly 50 colors")
    ErrHexFormat = errors.New("palette colors must be #RRGGBB hex")
)

// New creates a Palette from a slice of 50 hex colors.
func New(name string, colors []string) (Palette, error) {
    var p Palette
    if len(colors) != len(p.Colors) {
        return p, ErrSize
    }
    for i, c := range colors {
        if !isHex(c) {
            return p, fmt.Errorf("%w at index %d: %q", ErrHexFormat, i, c)
        }
        p.Colors[i] = normalizeHex(c)
    }
    p.Name = name
    return p, nil
}

// IsValid returns nil if the palette follows RetroForge constraints.
func (p Palette) IsValid() error {
    if len(p.Colors) != 50 {
        return ErrSize
    }
    for i, c := range p.Colors {
        if !isHex(c) {
            return fmt.Errorf("%w at index %d: %q", ErrHexFormat, i, c)
        }
    }
    return nil
}

func isHex(s string) bool {
    if len(s) != 7 || s[0] != '#' {
        return false
    }
    for i := 1; i < 7; i++ {
        b := s[i]
        if !((b >= '0' && b <= '9') || (b >= 'a' && b <= 'f') || (b >= 'A' && b <= 'F')) {
            return false
        }
    }
    return true
}

func normalizeHex(s string) string { return "#" + strings.ToLower(strings.TrimPrefix(s, "#")) }


