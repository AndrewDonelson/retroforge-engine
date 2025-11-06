package pal

import "image/color"

// Default48: default 48-color game palette (for indices 16-63)
// Built-in colors are separate (indices 0-15)
var Default48 = func() []color.RGBA {
	out := make([]color.RGBA, 48)
	// fill with a repeating pattern
	idx := 0
	for i := 0; i < 16 && idx < 48; i++ {
		base := uint8((i * 15) % 256)
		for s := 0; s < 3 && idx < 48; s++ {
			out[idx] = color.RGBA{base, uint8((int(base) + 40) % 256), uint8((int(base) + 80) % 256), 255}
			idx++
		}
	}
	for idx < 48 {
		v := uint8((idx * 5) & 255)
		out[idx] = color.RGBA{v, v, v, 255}
		idx++
	}
	return out
}()

type Manager struct {
	builtin []color.RGBA // Built-in colors (0-15), always available
	game    []color.RGBA // Game palette colors (16-63), 48 colors
}

func NewManager() *Manager {
	return &Manager{
		builtin: append([]color.RGBA{}, BuiltinColors...),
		game:    append([]color.RGBA{}, Default48...),
	}
}

// Color returns the color at index i (0-63)
// Indices 0-15: built-in colors
// Indices 16-63: game palette colors
func (m *Manager) Color(i int) color.RGBA {
	if i < 0 {
		return m.builtin[0] // Invalid index, return black
	}
	if i < 16 {
		return m.builtin[i] // Built-in color
	}
	if i < 64 {
		gameIdx := i - 16
		if gameIdx < len(m.game) {
			return m.game[gameIdx]
		}
	}
	return m.builtin[0] // Out of bounds, return black
}

// Set sets the game palette (48 colors, indices 16-63)
func (m *Manager) Set(name string) {
	// Get the named 48-color game palette (or default if not found)
	pal := GetPalette(name)
	m.game = pal
}
