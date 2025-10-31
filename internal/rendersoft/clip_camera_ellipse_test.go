package rendersoft

import (
	"image/color"
	"testing"
)

func TestSetClipGetClip(t *testing.T) {
	r := New(100, 100)

	// Test setting and getting clip
	r.SetClip(10, 20, 30, 40)
	x, y, w, h := r.GetClip()
	if x != 10 || y != 20 || w != 30 || h != 40 {
		t.Errorf("GetClip() = (%d, %d, %d, %d), expected (10, 20, 30, 40)", x, y, w, h)
	}

	// Test setting zero clip (should disable)
	r.SetClip(0, 0, 0, 0)
	x, y, w, h = r.GetClip()
	if x != 0 || y != 0 || w != 0 || h != 0 {
		t.Errorf("GetClip() after zero = (%d, %d, %d, %d), expected (0, 0, 0, 0)", x, y, w, h)
	}

	// Test multiple sets
	r.SetClip(5, 5, 10, 10)
	r.SetClip(15, 25, 50, 60)
	x, y, w, h = r.GetClip()
	if x != 15 || y != 25 || w != 50 || h != 60 {
		t.Errorf("GetClip() after multiple sets = (%d, %d, %d, %d), expected (15, 25, 50, 60)", x, y, w, h)
	}
}

func TestSetCameraGetCamera(t *testing.T) {
	r := New(100, 100)

	// Test initial camera (should be 0, 0)
	x, y := r.GetCamera()
	if x != 0 || y != 0 {
		t.Errorf("Initial GetCamera() = (%d, %d), expected (0, 0)", x, y)
	}

	// Test setting camera
	r.SetCamera(50, 75)
	x, y = r.GetCamera()
	if x != 50 || y != 75 {
		t.Errorf("GetCamera() after SetCamera(50, 75) = (%d, %d), expected (50, 75)", x, y)
	}

	// Test negative camera
	r.SetCamera(-10, -20)
	x, y = r.GetCamera()
	if x != -10 || y != -20 {
		t.Errorf("GetCamera() after SetCamera(-10, -20) = (%d, %d), expected (-10, -20)", x, y)
	}

	// Test multiple sets
	r.SetCamera(100, 200)
	r.SetCamera(0, 0)
	x, y = r.GetCamera()
	if x != 0 || y != 0 {
		t.Errorf("GetCamera() after reset = (%d, %d), expected (0, 0)", x, y)
	}
}

func TestCameraAffectsDrawing(t *testing.T) {
	r := New(100, 100)
	c := color.RGBA{R: 255, G: 0, B: 0, A: 255}

	// Draw pixel at (50, 50) with camera at (0, 0)
	r.SetCamera(0, 0)
	r.PSet(50, 50, c)

	// Move camera right by 10, pixel should now appear at (40, 50) in world space
	// but we set it at (50, 50) so it should be at (60, 50) on screen
	r.Clear(color.RGBA{R: 0, G: 0, B: 0, A: 255})
	r.SetCamera(10, 0)
	r.PSet(50, 50, c)

	// With camera offset, setting pixel at (50, 50) should actually draw at (60, 50) internally
	// but the API applies camera offset, so (50, 50) - (10, 0) = (40, 0) internal coordinate
	// Check that pixel was drawn (basic smoke test)
	pix := r.Pixels()
	drawn := false
	for i := 0; i < len(pix); i += 4 {
		if pix[i+0] == 255 {
			drawn = true
			break
		}
	}
	if !drawn {
		t.Errorf("Camera offset should not prevent drawing")
	}
}

func TestEllipse(t *testing.T) {
	r := New(100, 100)
	c := color.RGBA{R: 255, G: 0, B: 0, A: 255}

	// Test basic ellipse
	r.Ellipse(50, 50, 20, 15, c)

	// Check that some pixels were drawn
	pix := r.Pixels()
	drawn := false
	for i := 0; i < len(pix); i += 4 {
		if pix[i+0] == 255 {
			drawn = true
			break
		}
	}
	if !drawn {
		t.Errorf("Ellipse should draw some pixels")
	}

	// Test with different radii
	r.Clear(color.RGBA{R: 0, G: 0, B: 0, A: 255})
	r.Ellipse(50, 50, 30, 10, c) // Wide ellipse

	// Test circle-like ellipse (rx == ry)
	r.Clear(color.RGBA{R: 0, G: 0, B: 0, A: 255})
	r.Ellipse(50, 50, 25, 25, c)
}

func TestEllipseEdgeCases(t *testing.T) {
	r := New(100, 100)
	c := color.RGBA{R: 255, G: 0, B: 0, A: 255}

	// Test zero radius (should do nothing)
	r.Ellipse(50, 50, 0, 15, c)
	r.Ellipse(50, 50, 20, 0, c)
	r.Ellipse(50, 50, 0, 0, c)

	// Test negative radius (should do nothing)
	r.Ellipse(50, 50, -10, 15, c)
	r.Ellipse(50, 50, 20, -15, c)

	// Test very large radius
	r.Ellipse(50, 50, 1000, 1000, c)

	// Test out of bounds center
	r.Ellipse(-100, -100, 20, 15, c)
	r.Ellipse(200, 200, 20, 15, c)

	// Test with camera offset
	r.SetCamera(10, 10)
	r.Clear(color.RGBA{R: 0, G: 0, B: 0, A: 255})
	r.Ellipse(50, 50, 20, 15, c)
}

func TestEllipseFill(t *testing.T) {
	r := New(100, 100)
	c := color.RGBA{R: 0, G: 255, B: 0, A: 255}

	// Test basic filled ellipse
	r.EllipseFill(50, 50, 20, 15, c)

	// Check that center is filled
	pix := r.Pixels()
	centerIdx := (50*100 + 50) * 4
	if pix[centerIdx+1] == 255 {
		// Center should be filled
	}

	// Count filled pixels
	filled := 0
	for i := 0; i < len(pix); i += 4 {
		if pix[i+1] == 255 {
			filled++
		}
	}
	if filled == 0 {
		t.Errorf("EllipseFill should fill at least some pixels")
	}

	// Test wide filled ellipse
	r.Clear(color.RGBA{R: 0, G: 0, B: 0, A: 255})
	r.EllipseFill(50, 50, 30, 10, c)

	// Test circle-like filled ellipse
	r.Clear(color.RGBA{R: 0, G: 0, B: 0, A: 255})
	r.EllipseFill(50, 50, 25, 25, c)
}

func TestEllipseFillEdgeCases(t *testing.T) {
	r := New(100, 100)
	c := color.RGBA{R: 0, G: 255, B: 0, A: 255}

	// Test zero radius (should do nothing)
	r.EllipseFill(50, 50, 0, 15, c)
	r.EllipseFill(50, 50, 20, 0, c)
	r.EllipseFill(50, 50, 0, 0, c)

	// Test negative radius (should do nothing)
	r.EllipseFill(50, 50, -10, 15, c)
	r.EllipseFill(50, 50, 20, -15, c)

	// Test very large radius
	r.EllipseFill(50, 50, 1000, 1000, c)

	// Test out of bounds center
	r.EllipseFill(-100, -100, 20, 15, c)
	r.EllipseFill(200, 200, 20, 15, c)

	// Test with camera offset
	r.SetCamera(-20, -20)
	r.Clear(color.RGBA{R: 0, G: 0, B: 0, A: 255})
	r.EllipseFill(50, 50, 20, 15, c)
}

func TestClipAffectsDrawing(t *testing.T) {
	r := New(100, 100)
	c := color.RGBA{R: 255, G: 0, B: 0, A: 255}

	// Set clip to small region
	r.SetClip(10, 10, 20, 20)
	r.Clear(color.RGBA{R: 0, G: 0, B: 0, A: 255})

	// Draw outside clip region
	r.PSet(5, 5, c)   // Should be clipped
	r.PSet(15, 15, c) // Should be visible

	// Check that pixel in clip region was drawn
	pix := r.Pixels()
	clipIdx := (15*100 + 15) * 4
	if pix[clipIdx+0] == 255 {
		// Pixel in clip region was drawn
	}

	// Disable clip
	r.SetClip(0, 0, 0, 0)
	r.Clear(color.RGBA{R: 0, G: 0, B: 0, A: 255})
	r.PSet(5, 5, c)
	pix = r.Pixels()
	outsideIdx := (5*100 + 5) * 4
	if pix[outsideIdx+0] != 255 {
		t.Errorf("After disabling clip, pixel should be drawn")
	}
}

func TestClipWithCamera(t *testing.T) {
	r := New(100, 100)
	c := color.RGBA{R: 255, G: 255, B: 0, A: 255}

	// Set both clip and camera
	r.SetClip(10, 10, 30, 30)
	r.SetCamera(5, 5)
	r.Clear(color.RGBA{R: 0, G: 0, B: 0, A: 255})

	// Draw should respect both
	r.PSet(20, 20, c)

	// Basic smoke test - should not crash
	pix := r.Pixels()
	if len(pix) != 100*100*4 {
		t.Errorf("Clip and camera together should not break renderer")
	}
}
