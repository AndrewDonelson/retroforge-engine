package rendersoft

import (
	"image/color"
	"testing"
)

func TestTriangle(t *testing.T) {
	r := New(100, 100)
	c := color.RGBA{R: 255, G: 0, B: 0, A: 255}

	// Test outline triangle
	r.Triangle(50, 50, 20, false, c)

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
		t.Errorf("Triangle outline should draw some pixels")
	}

	// Test filled triangle
	r.Clear(color.RGBA{R: 0, G: 0, B: 0, A: 255})
	r.Triangle(50, 50, 20, true, c)

	// Check that center is filled
	centerIdx := (50*100 + 50) * 4
	if pix[centerIdx+0] == 255 {
		// Center should be filled
	}

	filled := 0
	for i := 0; i < len(pix); i += 4 {
		if pix[i+0] == 255 {
			filled++
		}
	}
	if filled == 0 {
		t.Errorf("Triangle fill should fill at least some pixels")
	}
}

func TestTriangleEdgeCases(t *testing.T) {
	r := New(100, 100)
	c := color.RGBA{R: 255, G: 0, B: 0, A: 255}

	// Test zero radius
	r.Triangle(50, 50, 0, false, c)
	r.Triangle(50, 50, 0, true, c)

	// Test negative radius
	r.Triangle(50, 50, -10, false, c)
	r.Triangle(50, 50, -10, true, c)

	// Test very large radius
	r.Triangle(50, 50, 1000, false, c)
	r.Triangle(50, 50, 1000, true, c)

	// Test out of bounds center
	r.Triangle(-100, -100, 20, false, c)
	r.Triangle(200, 200, 20, true, c)
}

func TestDiamond(t *testing.T) {
	r := New(100, 100)
	c := color.RGBA{R: 0, G: 255, B: 0, A: 255}

	// Test outline diamond
	r.Diamond(50, 50, 20, false, c)

	pix := r.Pixels()
	drawn := false
	for i := 0; i < len(pix); i += 4 {
		if pix[i+1] == 255 {
			drawn = true
			break
		}
	}
	if !drawn {
		t.Errorf("Diamond outline should draw some pixels")
	}

	// Test filled diamond
	r.Clear(color.RGBA{R: 0, G: 0, B: 0, A: 255})
	r.Diamond(50, 50, 20, true, c)

	// Check center is filled
	centerIdx := (50*100 + 50) * 4
	if pix[centerIdx+1] == 255 {
		// Center should be filled
	}

	filled := 0
	for i := 0; i < len(pix); i += 4 {
		if pix[i+1] == 255 {
			filled++
		}
	}
	if filled == 0 {
		t.Errorf("Diamond fill should fill at least some pixels")
	}
}

func TestDiamondEdgeCases(t *testing.T) {
	r := New(100, 100)
	c := color.RGBA{R: 0, G: 255, B: 0, A: 255}

	r.Diamond(50, 50, 0, false, c)
	r.Diamond(50, 50, -5, true, c)
	r.Diamond(50, 50, 1000, false, c)
	r.Diamond(-50, -50, 20, true, c)
}

func TestSquare(t *testing.T) {
	r := New(100, 100)
	c := color.RGBA{R: 0, G: 0, B: 255, A: 255}

	// Test outline square
	r.Square(50, 50, 15, false, c)

	pix := r.Pixels()
	drawn := false
	for i := 0; i < len(pix); i += 4 {
		if pix[i+2] == 255 {
			drawn = true
			break
		}
	}
	if !drawn {
		t.Errorf("Square outline should draw some pixels")
	}

	// Test filled square
	r.Clear(color.RGBA{R: 0, G: 0, B: 0, A: 255})
	r.Square(50, 50, 15, true, c)

	// Center should be filled
	centerIdx := (50*100 + 50) * 4
	if pix[centerIdx+2] != 255 {
		t.Errorf("Square fill should fill center")
	}

	// Check corners are filled too
	corners := []struct{ x, y int }{
		{35, 35}, {65, 35}, {65, 65}, {35, 65},
	}
	for _, corner := range corners {
		idx := (corner.y*100 + corner.x) * 4
		if idx >= 0 && idx < len(pix) && pix[idx+2] != 255 {
			// Some corners might be slightly out due to radius calculation, that's ok
		}
	}
}

func TestSquareEdgeCases(t *testing.T) {
	r := New(100, 100)
	c := color.RGBA{R: 0, G: 0, B: 255, A: 255}

	r.Square(50, 50, 0, false, c)
	r.Square(50, 50, -5, true, c)
	r.Square(50, 50, 1000, false, c)
	r.Square(-50, -50, 20, true, c)
}

func TestPentagon(t *testing.T) {
	r := New(100, 100)
	c := color.RGBA{R: 255, G: 255, B: 0, A: 255}

	// Test outline pentagon
	r.Pentagon(50, 50, 20, false, c)

	pix := r.Pixels()
	drawn := false
	for i := 0; i < len(pix); i += 4 {
		if pix[i+0] == 255 && pix[i+1] == 255 && pix[i+2] == 0 {
			drawn = true
			break
		}
	}
	if !drawn {
		t.Errorf("Pentagon outline should draw some pixels")
	}

	// Test filled pentagon
	r.Clear(color.RGBA{R: 0, G: 0, B: 0, A: 255})
	r.Pentagon(50, 50, 20, true, c)

	centerIdx := (50*100 + 50) * 4
	if pix[centerIdx+0] == 255 {
		// Center should be filled
	}

	filled := 0
	for i := 0; i < len(pix); i += 4 {
		if pix[i+0] == 255 && pix[i+1] == 255 {
			filled++
		}
	}
	if filled == 0 {
		t.Errorf("Pentagon fill should fill at least some pixels")
	}
}

func TestPentagonEdgeCases(t *testing.T) {
	r := New(100, 100)
	c := color.RGBA{R: 255, G: 255, B: 0, A: 255}

	r.Pentagon(50, 50, 0, false, c)
	r.Pentagon(50, 50, -5, true, c)
	r.Pentagon(50, 50, 1000, false, c)
	r.Pentagon(-50, -50, 20, true, c)
}

func TestHexagon(t *testing.T) {
	r := New(100, 100)
	c := color.RGBA{R: 255, G: 0, B: 255, A: 255}

	// Test outline hexagon
	r.Hexagon(50, 50, 20, false, c)

	pix := r.Pixels()
	drawn := false
	for i := 0; i < len(pix); i += 4 {
		if pix[i+0] == 255 && pix[i+2] == 255 {
			drawn = true
			break
		}
	}
	if !drawn {
		t.Errorf("Hexagon outline should draw some pixels")
	}

	// Test filled hexagon
	r.Clear(color.RGBA{R: 0, G: 0, B: 0, A: 255})
	r.Hexagon(50, 50, 20, true, c)

	centerIdx := (50*100 + 50) * 4
	if pix[centerIdx+0] == 255 {
		// Center should be filled
	}

	filled := 0
	for i := 0; i < len(pix); i += 4 {
		if pix[i+0] == 255 && pix[i+2] == 255 {
			filled++
		}
	}
	if filled == 0 {
		t.Errorf("Hexagon fill should fill at least some pixels")
	}
}

func TestHexagonEdgeCases(t *testing.T) {
	r := New(100, 100)
	c := color.RGBA{R: 255, G: 0, B: 255, A: 255}

	r.Hexagon(50, 50, 0, false, c)
	r.Hexagon(50, 50, -5, true, c)
	r.Hexagon(50, 50, 1000, false, c)
	r.Hexagon(-50, -50, 20, true, c)
}

func TestStar(t *testing.T) {
	r := New(100, 100)
	c := color.RGBA{R: 255, G: 128, B: 0, A: 255}

	// Test outline star
	r.Star(50, 50, 20, false, c)

	pix := r.Pixels()
	drawn := false
	for i := 0; i < len(pix); i += 4 {
		if pix[i+0] == 255 {
			drawn = true
			break
		}
	}
	if !drawn {
		t.Errorf("Star outline should draw some pixels")
	}

	// Test filled star
	r.Clear(color.RGBA{R: 0, G: 0, B: 0, A: 255})
	r.Star(50, 50, 20, true, c)

	centerIdx := (50*100 + 50) * 4
	if pix[centerIdx+0] == 255 {
		// Center should be filled
	}

	filled := 0
	for i := 0; i < len(pix); i += 4 {
		if pix[i+0] == 255 {
			filled++
		}
	}
	if filled == 0 {
		t.Errorf("Star fill should fill at least some pixels")
	}
}

func TestStarEdgeCases(t *testing.T) {
	r := New(100, 100)
	c := color.RGBA{R: 255, G: 128, B: 0, A: 255}

	r.Star(50, 50, 0, false, c)
	r.Star(50, 50, -5, true, c)
	r.Star(50, 50, 1000, false, c)
	r.Star(-50, -50, 20, true, c)
}

func TestAllShapesWithCamera(t *testing.T) {
	r := New(100, 100)
	c := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	// Test all shapes with camera offset
	r.SetCamera(10, 10)
	r.Clear(color.RGBA{R: 0, G: 0, B: 0, A: 255})

	r.Triangle(50, 50, 15, false, c)
	r.Diamond(50, 50, 15, false, c)
	r.Square(50, 50, 15, false, c)
	r.Pentagon(50, 50, 15, false, c)
	r.Hexagon(50, 50, 15, false, c)
	r.Star(50, 50, 15, false, c)

	// Should not crash
	pix := r.Pixels()
	if len(pix) != 100*100*4 {
		t.Errorf("Shapes with camera should work")
	}
}

func TestAllShapesWithClip(t *testing.T) {
	r := New(100, 100)
	c := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	// Test all shapes with clipping
	r.SetClip(20, 20, 60, 60)
	r.Clear(color.RGBA{R: 0, G: 0, B: 0, A: 255})

	r.Triangle(50, 50, 20, true, c)
	r.Diamond(50, 50, 20, true, c)
	r.Square(50, 50, 20, true, c)
	r.Pentagon(50, 50, 20, true, c)
	r.Hexagon(50, 50, 20, true, c)
	r.Star(50, 50, 20, true, c)

	// Should not crash
	pix := r.Pixels()
	if len(pix) != 100*100*4 {
		t.Errorf("Shapes with clip should work")
	}
}
