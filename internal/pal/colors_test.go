package pal

import "testing"

func TestNewManager(t *testing.T) {
	m := NewManager()
	if m == nil {
		t.Fatalf("NewManager returned nil")
	}

	// Should have default colors
	c := m.Color(0)
	if c.R != 0 || c.G != 0 || c.B != 0 {
		t.Fatalf("index 0 should be black, got %v", c)
	}

	c = m.Color(7)
	if c.R != 255 || c.G != 255 || c.B != 255 {
		t.Fatalf("index 7 should be white (built-in), got %v", c)
	}
}

func TestColor(t *testing.T) {
	m := NewManager()

	// Test valid indices (64 colors: 0-63)
	for i := 0; i < 64; i++ {
		c := m.Color(i)
		if c.A != 255 {
			t.Fatalf("color at index %d should have alpha 255, got %d", i, c.A)
		}
	}

	// Test negative index (should return default)
	c := m.Color(-1)
	if c.R != 0 || c.G != 0 || c.B != 0 {
		t.Fatalf("negative index should return black (index 0)")
	}

	// Test out of bounds (should return default)
	c = m.Color(999)
	if c.R != 0 || c.G != 0 || c.B != 0 {
		t.Fatalf("out of bounds index should return black (index 0)")
	}
}

func TestSet(t *testing.T) {
	m := NewManager()

	// Set to default (should reset to default palette)
	m.Set("default")

	c := m.Color(0)
	if c.R != 0 || c.G != 0 || c.B != 0 {
		t.Fatalf("Set('default') should reset to black at index 0")
	}

	c = m.Color(7)
	if c.R != 255 || c.G != 255 || c.B != 255 {
		t.Fatalf("Set('default') should reset to white at index 7 (built-in)")
	}

	// Test setting different name (currently only supports default, but shouldn't crash)
	m.Set("unknown")
	c = m.Color(0)
	if c.R != 0 || c.G != 0 || c.B != 0 {
		t.Fatalf("Set('unknown') should still give default")
	}
}

func TestDefault48(t *testing.T) {
	// Verify Default48 has correct length (48 colors for game palette)
	if len(Default48) != 48 {
		t.Fatalf("Default48 should have 48 colors, got %d", len(Default48))
	}

	// Verify all have alpha 255
	for i := 0; i < 48; i++ {
		if Default48[i].A != 255 {
			t.Fatalf("Default48[%d] should have alpha 255, got %d", i, Default48[i].A)
		}
	}
}

func TestBuiltinColors(t *testing.T) {
	// Verify BuiltinColors has correct length (16 colors)
	if len(BuiltinColors) != 16 {
		t.Fatalf("BuiltinColors should have 16 colors, got %d", len(BuiltinColors))
	}

	// Verify index 0 is black
	if BuiltinColors[0].R != 0 || BuiltinColors[0].G != 0 || BuiltinColors[0].B != 0 {
		t.Fatalf("BuiltinColors[0] should be black")
	}

	// Verify index 7 is white
	if BuiltinColors[7].R != 255 || BuiltinColors[7].G != 255 || BuiltinColors[7].B != 255 {
		t.Fatalf("BuiltinColors[7] should be white")
	}

	// Verify all have alpha 255
	for i := 0; i < 16; i++ {
		if BuiltinColors[i].A != 255 {
			t.Fatalf("BuiltinColors[%d] should have alpha 255, got %d", i, BuiltinColors[i].A)
		}
	}
}

func TestColorEdgeCases(t *testing.T) {
	m := NewManager()

	// Test exact boundaries
	c := m.Color(0)
	if c.R != 0 || c.G != 0 || c.B != 0 {
		t.Fatalf("Color(0) should return black")
	}

	c = m.Color(63) // last valid index (64 colors: 0-63)
	if c.A != 255 {
		t.Fatalf("Color(63) should have valid alpha")
	}

	// Test at boundary (should return default)
	c = m.Color(64)
	if c.R != 0 || c.G != 0 || c.B != 0 {
		t.Fatalf("Color(64) should return default (black)")
	}

	// Test very large values
	testCases := []int{-1, 64, 100, 999, 2147483647, -2147483648}
	for _, idx := range testCases {
		c := m.Color(idx)
		if c.R != 0 || c.G != 0 || c.B != 0 {
			t.Fatalf("Color(%d) should return default (black), got %v", idx, c)
		}
	}
}

func TestSetEdgeCases(t *testing.T) {
	m := NewManager()

	// Test setting with empty string
	m.Set("")
	c := m.Color(0)
	if c.R != 0 || c.G != 0 || c.B != 0 {
		t.Fatalf("Set(\"\") should still give default")
	}

	// Test setting with very long name
	m.Set(string(make([]byte, 10000)))
	c = m.Color(0)
	if c.R != 0 || c.G != 0 || c.B != 0 {
		t.Fatalf("Set with long string should give default")
	}

	// Test setting multiple times
	m.Set("default")
	m.Set("default")
	m.Set("default")
	c = m.Color(7) // White is now at index 7 (built-in)
	if c.R != 255 || c.G != 255 || c.B != 255 {
		t.Fatalf("Multiple Set calls should work, white should be at index 7")
	}
}

func TestManagerState(t *testing.T) {
	m1 := NewManager()
	m2 := NewManager()

	// Managers should be independent
	m1.Set("default")
	m2.Set("default")

	// Both should have same default state
	c1 := m1.Color(0)
	c2 := m2.Color(0)
	if c1.R != c2.R || c1.G != c2.G || c1.B != c2.B {
		t.Fatalf("Independent managers should have same default")
	}
}

func TestColorAllIndices(t *testing.T) {
	m := NewManager()

	// Test all valid indices (64 colors: 0-63)
	for i := 0; i < 64; i++ {
		c := m.Color(i)
		if c.A != 255 {
			t.Fatalf("Color(%d) should have alpha 255, got %d", i, c.A)
		}
		// Values should be in valid range
		if c.R > 255 || c.G > 255 || c.B > 255 {
			t.Fatalf("Color(%d) has out of range values: %v", i, c)
		}
	}
}
