package luabind

import (
	"testing"

	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
	"github.com/AndrewDonelson/retroforge-engine/internal/rendersoft"
	lua "github.com/yuin/gopher-lua"
)

func TestCursor(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	state := NewState()

	RegisterWithState(L, r, func(i int) (rgba [4]uint8) {
		return [4]uint8{0, 0, 0, 255}
	}, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil, state)

	// Test setting cursor
	err := L.DoString(`rf.cursor(10, 20)`)
	if err != nil {
		t.Fatalf("Failed to set cursor: %v", err)
	}

	x, y, has := state.GetCursor()
	if !has {
		t.Error("Cursor should be set")
	}
	if x != 10 || y != 20 {
		t.Errorf("Expected cursor at (10, 20), got (%d, %d)", x, y)
	}

	// Test resetting cursor
	err = L.DoString(`rf.cursor()`)
	if err != nil {
		t.Fatalf("Failed to reset cursor: %v", err)
	}

	_, _, has = state.GetCursor()
	if has {
		t.Error("Cursor should be reset (not set)")
	}

	// Test edge cases: negative coordinates
	err = L.DoString(`rf.cursor(-5, -10)`)
	if err != nil {
		t.Fatalf("Failed to set negative cursor: %v", err)
	}

	x, y, has = state.GetCursor()
	if !has || x != -5 || y != -10 {
		t.Errorf("Expected cursor at (-5, -10), got (%d, %d, has=%v)", x, y, has)
	}

	// Test large coordinates
	err = L.DoString(`rf.cursor(1000, 2000)`)
	if err != nil {
		t.Fatalf("Failed to set large cursor: %v", err)
	}

	x, y, has = state.GetCursor()
	if !has || x != 1000 || y != 2000 {
		t.Errorf("Expected cursor at (1000, 2000), got (%d, %d)", x, y)
	}
}

func TestColor(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	state := NewState()

	RegisterWithState(L, r, func(i int) (rgba [4]uint8) {
		return [4]uint8{0, 0, 0, 255}
	}, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil, state)

	// Test setting color
	err := L.DoString(`rf.color(5)`)
	if err != nil {
		t.Fatalf("Failed to set color: %v", err)
	}

	idx, has := state.GetTextColor()
	if !has {
		t.Error("Color should be set")
	}
	if idx != 5 {
		t.Errorf("Expected color 5, got %d", idx)
	}

	// Test resetting color
	err = L.DoString(`rf.color()`)
	if err != nil {
		t.Fatalf("Failed to reset color: %v", err)
	}

	idx, has = state.GetTextColor()
	if has {
		t.Error("Color should be reset (not set)")
	}
	if idx != 15 {
		t.Errorf("Expected default color 15, got %d", idx)
	}

	// Test edge cases: negative color
	err = L.DoString(`rf.color(-1)`)
	if err != nil {
		t.Fatalf("Failed to set negative color: %v", err)
	}

	idx, has = state.GetTextColor()
	if !has || idx != -1 {
		t.Errorf("Expected color -1, got %d (has=%v)", idx, has)
	}

	// Test large color index
	err = L.DoString(`rf.color(99)`)
	if err != nil {
		t.Fatalf("Failed to set large color: %v", err)
	}

	idx, has = state.GetTextColor()
	if !has || idx != 99 {
		t.Errorf("Expected color 99, got %d", idx)
	}
}

func TestPrintWithCursorColor(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	state := NewState()

	RegisterWithState(L, r, func(i int) (rgba [4]uint8) {
		return [4]uint8{0, 0, 0, 255}
	}, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil, state)

	// Test print with cursor and color state
	err := L.DoString(`
		rf.cursor(50, 100)
		rf.color(7)
		rf.print("Hello")
	`)
	if err != nil {
		t.Fatalf("Failed to print with cursor/color: %v", err)
	}

	// Check cursor was advanced
	x, y, has := state.GetCursor()
	if !has {
		t.Error("Cursor should still be set")
	}
	expectedX := 50 + len("Hello")*6 // font.Advance = 6
	if x != expectedX || y != 100 {
		t.Errorf("Expected cursor at (%d, 100), got (%d, %d)", expectedX, x, y)
	}

	// Check color is still set
	idx, has := state.GetTextColor()
	if !has || idx != 7 {
		t.Errorf("Expected color 7, got %d (has=%v)", idx, has)
	}

	// Test print with explicit coordinates (should override cursor)
	err = L.DoString(`rf.print("World", 200, 300, 3)`)
	if err != nil {
		t.Fatalf("Failed to print with explicit coords: %v", err)
	}

	// Cursor should not change when explicit coords are used
	x, y, has = state.GetCursor()
	if !has || x != expectedX || y != 100 {
		t.Errorf("Cursor should not change, expected (%d, 100), got (%d, %d)", expectedX, x, y)
	}
}

func TestPrintWithNewlines(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	state := NewState()

	RegisterWithState(L, r, func(i int) (rgba [4]uint8) {
		return [4]uint8{0, 0, 0, 255}
	}, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil, state)

	// Test print with newlines (cursor should advance properly)
	err := L.DoString(`
		rf.cursor(0, 0)
		rf.print("Line1\nLine2")
	`)
	if err != nil {
		t.Fatalf("Failed to print with newlines: %v", err)
	}

	// Check cursor position (should be at end of last line)
	x, _, has := state.GetCursor()
	if !has {
		t.Error("Cursor should be set")
	}
	// Cursor X should be at end of "Line2", Y should be advanced by line height
	expectedX := len("Line2") * 6
	if x != expectedX {
		t.Errorf("Expected cursor X %d, got %d", expectedX, x)
	}
}
