package graphics

import (
	"testing"
)

func TestNewTileMap(t *testing.T) {
	tm := NewTileMap(10, 20)
	if tm == nil {
		t.Fatal("NewTileMap returned nil")
	}
	if tm.Width() != 10 {
		t.Errorf("Width() = %d, expected 10", tm.Width())
	}
	if tm.Height() != 20 {
		t.Errorf("Height() = %d, expected 20", tm.Height())
	}
}

func TestTileMapGetSet(t *testing.T) {
	tm := NewTileMap(10, 10)

	// Test setting and getting valid coordinates
	tm.Set(5, 5, 42)
	if got := tm.Get(5, 5); got != 42 {
		t.Errorf("Get(5, 5) = %d, expected 42", got)
	}

	// Test out of bounds get
	if got := tm.Get(-1, 5); got != 0 {
		t.Errorf("Get(-1, 5) should return 0, got %d", got)
	}
	if got := tm.Get(5, -1); got != 0 {
		t.Errorf("Get(5, -1) should return 0, got %d", got)
	}
	if got := tm.Get(10, 5); got != 0 {
		t.Errorf("Get(10, 5) should return 0 (out of bounds), got %d", got)
	}
	if got := tm.Get(5, 10); got != 0 {
		t.Errorf("Get(5, 10) should return 0 (out of bounds), got %d", got)
	}

	// Test out of bounds set (should silently fail)
	tm.Set(-1, 5, 99)
	tm.Set(5, -1, 99)
	tm.Set(10, 5, 99)
	tm.Set(5, 10, 99)
	// Verify surrounding cells weren't affected
	if got := tm.Get(0, 5); got != 0 {
		t.Errorf("Get(0, 5) should still be 0 after out-of-bounds Set(-1, 5)")
	}

	// Test multiple sets
	tm.Set(0, 0, 1)
	tm.Set(9, 9, 99)
	if got := tm.Get(0, 0); got != 1 {
		t.Errorf("Get(0, 0) = %d, expected 1", got)
	}
	if got := tm.Get(9, 9); got != 99 {
		t.Errorf("Get(9, 9) = %d, expected 99", got)
	}
}

func TestTileMapDraw(t *testing.T) {
	tm := NewTileMap(8, 8)

	// Set up a pattern
	tm.Set(1, 1, 1)
	tm.Set(2, 1, 2)
	tm.Set(1, 2, 3)
	tm.Set(2, 2, 4)

	// Track what was drawn
	var drawn []struct{ x, y, tile int }
	spriteRenderer := func(x, y, tileIndex int) {
		drawn = append(drawn, struct{ x, y, tile int }{x, y, tileIndex})
	}

	// Draw a 2x2 region starting at tile (1, 1)
	tm.Draw(1, 1, 0, 0, 2, 2, spriteRenderer)

	// Should draw 4 tiles (none are 0)
	if len(drawn) != 4 {
		t.Errorf("Draw should have drawn 4 tiles, got %d", len(drawn))
	}

	// Verify positions
	expected := []struct{ x, y, tile int }{
		{0, 0, 1}, // tile (1,1) -> screen (0,0) with tile 1
		{8, 0, 2}, // tile (2,1) -> screen (8,0) with tile 2 (8px per tile)
		{0, 8, 3}, // tile (1,2) -> screen (0,8) with tile 3
		{8, 8, 4}, // tile (2,2) -> screen (8,8) with tile 4
	}

	if len(drawn) == len(expected) {
		for _, exp := range expected {
			found := false
			for _, d := range drawn {
				if d.x == exp.x && d.y == exp.y && d.tile == exp.tile {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected draw call not found: screen(%d, %d) tile %d", exp.x, exp.y, exp.tile)
			}
		}
	}

	// Test drawing with empty tiles (should skip 0 tiles)
	tm.Set(3, 3, 0) // Empty tile
	drawn = drawn[:0]
	tm.Draw(3, 3, 0, 0, 1, 1, spriteRenderer)
	if len(drawn) != 0 {
		t.Errorf("Draw should skip empty (0) tiles, but drew %d tiles", len(drawn))
	}

	// Test drawing partial region
	drawn = drawn[:0]
	tm.Set(0, 0, 5)
	tm.Draw(0, 0, 100, 200, 1, 1, spriteRenderer)
	if len(drawn) != 1 {
		t.Errorf("Draw should have drawn 1 tile, got %d", len(drawn))
	}
	if len(drawn) > 0 && (drawn[0].x != 100 || drawn[0].y != 200 || drawn[0].tile != 5) {
		t.Errorf("Draw position incorrect, got screen(%d, %d) tile %d, expected (100, 200) tile 5",
			drawn[0].x, drawn[0].y, drawn[0].tile)
	}

	// Test drawing larger region
	tm = NewTileMap(10, 10)
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			tm.Set(x, y, x+y+1) // Non-zero values
		}
	}
	drawn = drawn[:0]
	tm.Draw(0, 0, 0, 0, 5, 5, spriteRenderer)
	if len(drawn) != 25 {
		t.Errorf("Draw should have drawn 25 tiles, got %d", len(drawn))
	}
}

func TestTileMapEdgeCases(t *testing.T) {
	// Test 1x1 tilemap
	tm := NewTileMap(1, 1)
	tm.Set(0, 0, 10)
	if got := tm.Get(0, 0); got != 10 {
		t.Errorf("1x1 tilemap: Get(0, 0) = %d, expected 10", got)
	}

	// Test large tilemap
	tm = NewTileMap(128, 64)
	tm.Set(127, 63, 255)
	if got := tm.Get(127, 63); got != 255 {
		t.Errorf("Large tilemap: Get(127, 63) = %d, expected 255", got)
	}

	// Test Draw with zero size
	var drawn int
	spriteRenderer := func(x, y, tileIndex int) { drawn++ }
	tm.Draw(0, 0, 0, 0, 0, 0, spriteRenderer)
	if drawn != 0 {
		t.Errorf("Draw with zero size should draw nothing, drew %d", drawn)
	}
}
