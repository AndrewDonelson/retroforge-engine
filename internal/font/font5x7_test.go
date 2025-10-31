package font

import "testing"

func TestConstants(t *testing.T) {
	if Width != 5 {
		t.Fatalf("Width should be 5, got %d", Width)
	}
	if Height != 7 {
		t.Fatalf("Height should be 7, got %d", Height)
	}
	if Advance != 6 {
		t.Fatalf("Advance should be 6, got %d", Advance)
	}
}

func TestGet(t *testing.T) {
	// Test existing glyphs
	testCases := []rune{'A', 'B', 'C', '0', '1', ' ', '!', ','}
	for _, r := range testCases {
		g, ok := Get(r)
		if !ok {
			t.Fatalf("Get('%c') should return ok=true", r)
		}
		if g.W != Width {
			t.Fatalf("glyph width should be %d, got %d", Width, g.W)
		}
		if g.H != Height {
			t.Fatalf("glyph height should be %d, got %d", Height, g.H)
		}
	}

	// Test lowercase (should convert to uppercase)
	g, ok := Get('a')
	if !ok {
		t.Fatalf("Get('a') should convert to uppercase and return ok=true")
	}
	g2, _ := Get('A')
	if g.W != g2.W || g.H != g2.H {
		t.Fatalf("lowercase 'a' should return same glyph as 'A'")
	}

	// Test missing glyph
	g, ok = Get('?')
	if ok {
		t.Fatalf("Get('?') should return ok=false for unsupported character")
	}
}

func TestSpaceGlyph(t *testing.T) {
	g, ok := Get(' ')
	if !ok {
		t.Fatalf("space should be supported")
	}
	// Space should be empty (all zeros)
	for i := 0; i < len(g.Rows); i++ {
		if g.Rows[i] != 0 {
			t.Fatalf("space glyph should have all zero rows")
		}
	}
}

func TestNumberGlyphs(t *testing.T) {
	// Test all digits
	for r := '0'; r <= '9'; r++ {
		g, ok := Get(r)
		if !ok {
			t.Fatalf("digit '%c' should be supported", r)
		}
		if g.W != 5 || g.H != 7 {
			t.Fatalf("digit '%c' should have correct dimensions", r)
		}
	}
}

func TestLetterGlyphs(t *testing.T) {
	// Test all uppercase letters
	for r := 'A'; r <= 'Z'; r++ {
		g, ok := Get(r)
		if !ok {
			t.Fatalf("letter '%c' should be supported", r)
		}
		if g.W != 5 || g.H != 7 {
			t.Fatalf("letter '%c' should have correct dimensions", r)
		}
	}
}

func TestCaseInsensitive(t *testing.T) {
	// Lowercase should map to uppercase
	for r := 'a'; r <= 'z'; r++ {
		upper := rune(r - 32)
		gLower, okLower := Get(r)
		gUpper, okUpper := Get(upper)

		if okLower != okUpper {
			t.Fatalf("lowercase '%c' and uppercase '%c' should have same ok value", r, upper)
		}
		if okLower {
			if gLower.W != gUpper.W || gLower.H != gUpper.H {
				t.Fatalf("lowercase '%c' should return same glyph as uppercase", r)
			}
		}
	}
}

func TestGetEdgeCases(t *testing.T) {
	// Test null/zero rune
	g, ok := Get(0)
	if ok {
		t.Fatalf("null rune should return ok=false")
	}
	if g.W != 0 || g.H != 0 {
		t.Fatalf("null rune should return zero glyph")
	}

	// Test unsupported punctuation (excluding ':' which is supported)
	unsupported := []rune{'@', '#', '$', '%', '^', '&', '*', '(', ')', '[', ']', '{', '}', '|', '\\', '/', '?', '<', '>', ';', '"', '\'', '`', '~', '-', '_', '=', '+'}
	for _, r := range unsupported {
		g, ok := Get(r)
		if ok {
			t.Fatalf("unsupported rune '%c' should return ok=false", r)
		}
		if g.W != 0 || g.H != 0 {
			t.Logf("unsupported rune '%c' returned non-zero glyph (may be acceptable)", r)
		}
	}

	// Test unicode characters beyond ASCII
	unicodeStr := "áéñ🚀αβ∞中文"
	for _, r := range unicodeStr {
		g, ok := Get(r)
		if ok {
			// If it converts to uppercase and matches, that's ok
			if r >= 'a' && r <= 'z' {
				continue // handled by case conversion
			}
			t.Logf("unicode rune '%c' returned ok=true (unexpected)", r)
		}
		if g.W == 0 && g.H == 0 {
			// Expected for unsupported characters
		}
	}

	// Test very large rune values
	g, ok = Get(rune(0x10FFFF)) // max unicode
	if ok {
		t.Logf("very large rune returned ok=true")
	}
}

func TestGlyphStructure(t *testing.T) {
	// Test that all returned glyphs have valid structure
	for r := 'A'; r <= 'Z'; r++ {
		g, ok := Get(r)
		if !ok {
			continue
		}

		if g.W <= 0 || g.H <= 0 {
			t.Fatalf("glyph for '%c' has invalid dimensions: %dx%d", r, g.W, g.H)
		}

		if g.W != Width || g.H != Height {
			t.Fatalf("glyph for '%c' has wrong dimensions: got %dx%d, want %dx%d", r, g.W, g.H, Width, Height)
		}

		// Check rows array
		for i := 0; i < len(g.Rows); i++ {
			// Rows should be valid (no specific check, just ensure it doesn't crash)
			_ = g.Rows[i]
		}
	}
}

func TestSpaceGlyphStructure(t *testing.T) {
	g, ok := Get(' ')
	if !ok {
		t.Fatalf("space should be supported")
	}

	if g.W != Width || g.H != Height {
		t.Fatalf("space glyph should have correct dimensions")
	}

	// Space should have all zero rows
	for i := 0; i < len(g.Rows); i++ {
		if g.Rows[i] != 0 {
			t.Fatalf("space glyph row %d should be 0, got %d", i, g.Rows[i])
		}
	}
}
