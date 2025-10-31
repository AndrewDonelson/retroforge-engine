package luabind

import (
	"testing"
)

func TestStateGetTileMap(t *testing.T) {
	s := NewState()
	tm := s.GetTileMap()
	if tm == nil {
		t.Fatal("GetTileMap returned nil")
	}
	if tm.Width() != 256 || tm.Height() != 256 {
		t.Errorf("TileMap dimensions wrong: %dx%d, expected 256x256", tm.Width(), tm.Height())
	}
}

func TestStateSetPalRemap(t *testing.T) {
	s := NewState()

	// Test setting remap
	s.SetPalRemap(5, 10, true)
	if s.GetPalRemap(5) != 10 {
		t.Errorf("SetPalRemap(5, 10, true) failed, GetPalRemap(5) = %d, expected 10", s.GetPalRemap(5))
	}

	// Test disabling remap (p=false)
	s.SetPalRemap(5, 10, false)
	if s.GetPalRemap(5) != 5 {
		t.Errorf("SetPalRemap(5, 10, false) should reset, GetPalRemap(5) = %d, expected 5", s.GetPalRemap(5))
	}

	// Test out of bounds
	s.SetPalRemap(-1, 10, true)
	s.SetPalRemap(256, 10, true)
	// Should not crash or affect valid indices
	if s.GetPalRemap(0) != 0 {
		t.Error("Out of bounds SetPalRemap should not affect valid indices")
	}
}

func TestStateResetPalRemap(t *testing.T) {
	s := NewState()

	// Set some remaps
	s.SetPalRemap(1, 5, true)
	s.SetPalRemap(2, 6, true)
	s.SetPalRemap(3, 7, true)

	// Verify they're set
	if s.GetPalRemap(1) != 5 || s.GetPalRemap(2) != 6 || s.GetPalRemap(3) != 7 {
		t.Error("Pal remaps should be set")
	}

	// Reset
	s.ResetPalRemap()

	// Verify all reset
	if s.GetPalRemap(1) != 1 || s.GetPalRemap(2) != 2 || s.GetPalRemap(3) != 3 {
		t.Error("ResetPalRemap should reset all remaps")
	}
}

func TestStateGetRNGSeed(t *testing.T) {
	s := NewState()

	seed := s.GetRNGSeed()
	if seed != 1 {
		t.Errorf("Initial seed should be 1, got %d", seed)
	}

	// Set seed and verify
	s.SetRNGSeed(42)
	if s.GetRNGSeed() != 42 {
		t.Errorf("GetRNGSeed after SetRNGSeed(42) = %d, expected 42", s.GetRNGSeed())
	}
}
