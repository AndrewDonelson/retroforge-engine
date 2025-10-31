package luabind

import (
	"github.com/AndrewDonelson/retroforge-engine/internal/graphics"
)

// State holds persistent state for Lua bindings (tilemap, memory, color remapping)
type State struct {
	tileMap   *graphics.TileMap
	memory    []byte   // Memory for poke/peek (default 2MB like PICO-8)
	palRemap  [256]int // Color remapping: palRemap[oldIndex] = newIndex
	palActive bool     // Whether color remapping is active
	cartStore []byte   // Cart storage for cstore/reload (default 64KB like PICO-8)
	cursorX   int      // Text cursor X position
	cursorY   int      // Text cursor Y position
	textColor int      // Text color index (-1 = use default)
	hasCursor bool     // Whether cursor has been set
	hasColor  bool     // Whether color has been set
	rngSeed   uint32   // Random number generator seed (for deterministic rnd())
}

// NewState creates a new state with default tilemap and memory
func NewState() *State {
	s := &State{
		tileMap:   graphics.NewTileMap(256, 256), // Default 256×256 tilemap
		memory:    make([]byte, 2*1024*1024),     // 2MB like PICO-8
		cartStore: make([]byte, 64*1024),         // 64KB cart storage (2x PICO-8's 32KB)
		palRemap:  [256]int{},                    // Will be initialized to identity mapping
		palActive: false,
		cursorX:   0,
		cursorY:   0,
		textColor: 15, // Default white
		hasCursor: false,
		hasColor:  false,
		rngSeed:   1, // Initial seed (PICO-8 compatible)
	}
	// Initialize palRemap to identity mapping
	for i := range s.palRemap {
		s.palRemap[i] = i
	}
	return s
}

// GetTileMap returns the tilemap
func (s *State) GetTileMap() *graphics.TileMap {
	return s.tileMap
}

// GetMemory returns the memory array
func (s *State) GetMemory() []byte {
	return s.memory
}

// SetPalRemap sets color remapping: oldIndex -> newIndex
func (s *State) SetPalRemap(oldIndex, newIndex int, p bool) {
	if oldIndex >= 0 && oldIndex < 256 {
		if p {
			s.palRemap[oldIndex] = newIndex
		} else {
			s.palRemap[oldIndex] = oldIndex // Reset to identity
		}
		s.palActive = true
	}
}

// GetPalRemap gets the remapped color index
func (s *State) GetPalRemap(index int) int {
	if index < 0 || index >= 256 {
		return index // Out of range, return as-is
	}
	if s.palActive {
		return s.palRemap[index]
	}
	return index
}

// ResetPalRemap resets all color remapping
func (s *State) ResetPalRemap() {
	for i := range s.palRemap {
		s.palRemap[i] = i
	}
	s.palActive = false
}

// GetCartStore returns the cart storage array
func (s *State) GetCartStore() []byte {
	return s.cartStore
}

// SetCursor sets the text cursor position
func (s *State) SetCursor(x, y int) {
	s.cursorX = x
	s.cursorY = y
	s.hasCursor = true
}

// GetCursor returns the text cursor position
func (s *State) GetCursor() (x, y int, has bool) {
	return s.cursorX, s.cursorY, s.hasCursor
}

// SetTextColor sets the text color index
func (s *State) SetTextColor(index int) {
	s.textColor = index
	s.hasColor = true
}

// GetTextColor returns the text color index
func (s *State) GetTextColor() (index int, has bool) {
	return s.textColor, s.hasColor
}

// ResetCursor resets cursor to default state
func (s *State) ResetCursor() {
	s.cursorX = 0
	s.cursorY = 0
	s.hasCursor = false
}

// ResetColor resets color to default state
func (s *State) ResetColor() {
	s.textColor = 15 // Default white
	s.hasColor = false
}

// GetRNGSeed returns the current RNG seed
func (s *State) GetRNGSeed() uint32 {
	return s.rngSeed
}

// SetRNGSeed sets the RNG seed (for deterministic random numbers)
func (s *State) SetRNGSeed(seed uint32) {
	s.rngSeed = seed
}

// NextRandom generates next random number and updates seed (PICO-8 compatible LCG)
func (s *State) NextRandom() float64 {
	// PICO-8 LCG: seed = (seed * 1103515245 + 12345) & 0x7FFFFFFF
	s.rngSeed = (s.rngSeed*1103515245 + 12345) & 0x7FFFFFFF
	return float64(s.rngSeed) / 2147483648.0 // Returns 0.0 to ~0.999...
}
