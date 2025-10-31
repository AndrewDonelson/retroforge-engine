package graphics

import (
	"image/color"
	"testing"

	"github.com/AndrewDonelson/retroforge-engine/internal/rendersoft"
)

// TestRenderer is a concrete implementation for testing the interface
func TestRendererInterface(t *testing.T) {
	var r Renderer = rendersoft.New(100, 100)

	// Test dimensions
	if r.Width() != 100 {
		t.Fatalf("expected width 100, got %d", r.Width())
	}
	if r.Height() != 100 {
		t.Fatalf("expected height 100, got %d", r.Height())
	}

	// Test pixels
	pix := r.Pixels()
	expectedLen := 100 * 100 * 4
	if len(pix) != expectedLen {
		t.Fatalf("expected %d pixels, got %d", expectedLen, len(pix))
	}
}

func TestClear(t *testing.T) {
	r := rendersoft.New(10, 10)
	c := color.RGBA{R: 128, G: 64, B: 32, A: 255}
	r.Clear(c)

	pix := r.Pixels()
	for i := 0; i < len(pix); i += 4 {
		if pix[i+0] != 128 || pix[i+1] != 64 || pix[i+2] != 32 {
			t.Fatalf("Clear failed at pixel %d", i/4)
		}
	}
}

func TestPrint(t *testing.T) {
	r := rendersoft.New(100, 50)
	c := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	r.Print("TEST", 10, 10, c)

	// Verify some pixels were set (text rendering)
	pix := r.Pixels()
	filled := 0
	for i := 0; i < len(pix); i += 4 {
		if pix[i+0] == 255 {
			filled++
		}
	}
	if filled == 0 {
		t.Fatalf("Print should render text")
	}
}

func TestPrintCentered(t *testing.T) {
	r := rendersoft.New(100, 50)
	c := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	r.PrintCentered("HI", 25, c)

	// Verify text was rendered
	pix := r.Pixels()
	filled := 0
	for i := 0; i < len(pix); i += 4 {
		if pix[i+0] == 255 {
			filled++
		}
	}
	if filled == 0 {
		t.Fatalf("PrintCentered should render text")
	}
}

func TestPrimitives(t *testing.T) {
	r := rendersoft.New(20, 20)
	c := color.RGBA{R: 255, G: 0, B: 0, A: 255}

	// Test all primitives exist and don't crash
	r.PSet(5, 5, c)
	r.Line(0, 0, 19, 19, c)
	r.Rect(2, 2, 10, 10, c)
	r.RectFill(5, 5, 15, 15, c)
	r.Circ(10, 10, 5, c)
	r.CircFill(10, 10, 3, c)

	// If we get here, primitives work (basic smoke test)
}

func TestRendererEdgeCases(t *testing.T) {
	// Test zero-size renderer
	r := rendersoft.New(0, 0)
	if r.Width() != 0 || r.Height() != 0 {
		t.Fatalf("zero-size renderer should have zero dimensions")
	}
	pix := r.Pixels()
	if len(pix) != 0 {
		t.Fatalf("zero-size renderer should have empty pixels")
	}

	// Test 1x1 renderer
	r = rendersoft.New(1, 1)
	if len(r.Pixels()) != 4 {
		t.Fatalf("1x1 renderer should have 4 bytes")
	}

	// Test very large renderer
	r = rendersoft.New(10000, 10000)
	if r.Width() != 10000 || r.Height() != 10000 {
		t.Fatalf("large renderer should have correct dimensions")
	}
	expectedSize := 10000 * 10000 * 4
	if len(r.Pixels()) != expectedSize {
		t.Fatalf("large renderer should have correct pixel count")
	}
}

func TestRendererInvalidColors(t *testing.T) {
	r := rendersoft.New(10, 10)

	// Test with colors that have invalid alpha (should still work)
	c := color.RGBA{R: 255, G: 0, B: 0, A: 0}
	r.Clear(c)

	// Test with out-of-range values (uint8 will wrap, but should not crash)
	// Note: Go color.RGBA uses uint8, so values >255 wrap automatically
	// Use explicit wrapping calculation
	c = color.RGBA{R: 300 % 256, G: 300 % 256, B: 300 % 256, A: 300 % 256} // will wrap to valid range
	r.Clear(c)
}

func TestRendererOperationsOnZeroSize(t *testing.T) {
	r := rendersoft.New(0, 0)
	c := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	// All operations should work without crashing on zero-size
	r.Clear(c)
	r.PSet(0, 0, c)
	r.Line(0, 0, 0, 0, c)
	r.Rect(0, 0, 0, 0, c)
	r.RectFill(0, 0, 0, 0, c)
	r.Circ(0, 0, 0, c)
	r.CircFill(0, 0, 0, c)
	r.Print("TEST", 0, 0, c)
	r.PrintCentered("TEST", 0, c)

	// Should not crash
}
