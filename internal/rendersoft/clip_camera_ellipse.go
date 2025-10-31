package rendersoft

import (
	"image/color"
	"math"
)

// Clip and Camera functions
func (s *Soft) SetClip(x, y, w, h int) {
	s.clipX, s.clipY, s.clipW, s.clipH = x, y, w, h
}

func (s *Soft) GetClip() (x, y, w, h int) {
	return s.clipX, s.clipY, s.clipW, s.clipH
}

func (s *Soft) SetCamera(x, y int) {
	s.cameraX, s.cameraY = x, y
}

func (s *Soft) GetCamera() (x, y int) {
	return s.cameraX, s.cameraY
}

// Ellipse drawing functions
func (s *Soft) Ellipse(xc, yc, rx, ry int, c color.RGBA) {
	if rx <= 0 || ry <= 0 {
		return
	}

	// Apply camera offset
	xc -= s.cameraX
	yc -= s.cameraY

	// Midpoint ellipse algorithm for outline
	rx2 := rx * rx
	ry2 := ry * ry
	twoRx2 := 2 * rx2
	twoRy2 := 2 * ry2

	var x, y int
	var px, py int

	// Region 1
	x = 0
	y = ry
	px = 0
	py = twoRx2 * y

	// Draw first set of points
	s.ellipsePlotPoints(xc, yc, x, y, c)

	p := ry2 - (rx2 * ry) + (rx2 / 4)
	for px < py {
		x++
		px += twoRy2
		if p < 0 {
			p += ry2 + px
		} else {
			y--
			py -= twoRx2
			p += ry2 + px - py
		}
		s.ellipsePlotPoints(xc, yc, x, y, c)
	}

	// Region 2
	p = ry2*(x+1)*(x+1) + rx2*(y-1)*(y-1) - rx2*ry2
	for y > 0 {
		y--
		py -= twoRx2
		if p > 0 {
			p += rx2 - py
		} else {
			x++
			px += twoRy2
			p += rx2 - py + px
		}
		s.ellipsePlotPoints(xc, yc, x, y, c)
	}
}

func (s *Soft) EllipseFill(xc, yc, rx, ry int, c color.RGBA) {
	if rx <= 0 || ry <= 0 {
		return
	}

	// Apply camera offset
	xc -= s.cameraX
	yc -= s.cameraY

	// Fill ellipse using horizontal scanlines
	ry2 := ry * ry
	maxY := ry

	for y := 0; y <= maxY; y++ {
		// Calculate width of ellipse at this y coordinate
		if ry2 == 0 {
			continue
		}
		width := int(float64(rx) * math.Sqrt(1.0-float64(y*y)/float64(ry2)))

		// Draw horizontal line from -width to +width
		for x := -width; x <= width; x++ {
			s.set(xc+x, yc+y, c)
			s.set(xc+x, yc-y, c)
		}
	}
}

func (s *Soft) ellipsePlotPoints(xc, yc, x, y int, c color.RGBA) {
	s.set(xc+x, yc+y, c)
	s.set(xc-x, yc+y, c)
	s.set(xc+x, yc-y, c)
	s.set(xc-x, yc-y, c)
}
