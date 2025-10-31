package rendersoft

import (
	"image/color"
	"math"
)

// Helper functions for shape drawing
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// fillPolygon uses scanline algorithm to fill a polygon
func (s *Soft) fillPolygon(points [][]int, c color.RGBA) {
	if len(points) < 3 {
		return
	}

	// Find bounding box
	minY := points[0][1]
	maxY := points[0][1]
	for _, p := range points {
		if p[1] < minY {
			minY = p[1]
		}
		if p[1] > maxY {
			maxY = p[1]
		}
	}

	// Clamp to screen bounds
	minY = maxInt(0, minY)
	maxY = minInt(s.h-1, maxY)

	// Scanline fill
	for y := minY; y <= maxY; y++ {
		var intersects []int

		// Check each edge
		for i := 0; i < len(points); i++ {
			next := (i + 1) % len(points)
			px0, py0 := points[i][0], points[i][1]
			px1, py1 := points[next][0], points[next][1]

			if py0 == py1 {
				if py0 == y {
					intersects = append(intersects, px0, px1)
				}
				continue
			}

			if (py0 < y && py1 >= y) || (py1 < y && py0 >= y) {
				denom := py1 - py0
				if abs(denom) > 0 {
					x := px0 + (y-py0)*(px1-px0)/denom
					intersects = append(intersects, x)
				}
			}
		}

		if len(intersects) >= 2 {
			// Sort intersects
			for i := 0; i < len(intersects)-1; i++ {
				for j := i + 1; j < len(intersects); j++ {
					if intersects[i] > intersects[j] {
						intersects[i], intersects[j] = intersects[j], intersects[i]
					}
				}
			}

			minX := maxInt(0, intersects[0])
			maxX := minInt(s.w-1, intersects[len(intersects)-1])
			for x := minX; x <= maxX; x++ {
				s.set(x, y, c)
			}
		}
	}
}

func (s *Soft) Triangle(cx, cy, radius int, filled bool, c color.RGBA) {
	// Apply camera offset
	cx -= s.cameraX
	cy -= s.cameraY

	x0 := cx
	y0 := cy - radius
	x1 := cx - int(float64(radius)*0.866)
	y1 := cy + radius/2
	x2 := cx + int(float64(radius)*0.866)
	y2 := cy + radius/2

	if filled {
		points := [][]int{{x0, y0}, {x1, y1}, {x2, y2}}
		s.fillPolygon(points, c)
	} else {
		// Line already applies camera, but we need to pass world coords
		// So add camera back temporarily
		s.Line(x0+s.cameraX, y0+s.cameraY, x1+s.cameraX, y1+s.cameraY, c)
		s.Line(x1+s.cameraX, y1+s.cameraY, x2+s.cameraX, y2+s.cameraY, c)
		s.Line(x2+s.cameraX, y2+s.cameraY, x0+s.cameraX, y0+s.cameraY, c)
	}
}

func (s *Soft) Diamond(cx, cy, radius int, filled bool, c color.RGBA) {
	// Apply camera offset
	cx -= s.cameraX
	cy -= s.cameraY

	if filled {
		minY := maxInt(0, cy-radius)
		maxY := minInt(s.h-1, cy+radius)
		for y := minY; y <= maxY; y++ {
			dy := abs(y - cy)
			width := radius - dy
			if width >= 0 {
				minX := maxInt(0, cx-width)
				maxX := minInt(s.w-1, cx+width)
				for x := minX; x <= maxX; x++ {
					s.set(x, y, c)
				}
			}
		}
	} else {
		points := [][]int{
			{cx, cy - radius},
			{cx + radius, cy},
			{cx, cy + radius},
			{cx - radius, cy},
		}
		for i := 0; i < len(points); i++ {
			next := (i + 1) % len(points)
			s.Line(points[i][0]+s.cameraX, points[i][1]+s.cameraY, points[next][0]+s.cameraX, points[next][1]+s.cameraY, c)
		}
	}
}

func (s *Soft) Square(cx, cy, radius int, filled bool, c color.RGBA) {
	// RectFill/Rect already apply camera
	x0 := cx - radius
	y0 := cy - radius
	x1 := cx + radius
	y1 := cy + radius
	if filled {
		s.RectFill(x0, y0, x1, y1, c)
	} else {
		s.Rect(x0, y0, x1, y1, c)
	}
}

func (s *Soft) Pentagon(cx, cy, radius int, filled bool, c color.RGBA) {
	// Apply camera offset
	cx -= s.cameraX
	cy -= s.cameraY

	var points [][]int
	for i := 0; i < 5; i++ {
		angle := float64(i)*2*math.Pi/5 - math.Pi/2
		x := cx + int(float64(radius)*math.Cos(angle))
		y := cy + int(float64(radius)*math.Sin(angle))
		points = append(points, []int{x, y})
	}

	if filled {
		s.fillPolygon(points, c)
	} else {
		for i := 0; i < len(points); i++ {
			next := (i + 1) % len(points)
			s.Line(points[i][0]+s.cameraX, points[i][1]+s.cameraY, points[next][0]+s.cameraX, points[next][1]+s.cameraY, c)
		}
	}
}

func (s *Soft) Hexagon(cx, cy, radius int, filled bool, c color.RGBA) {
	// Apply camera offset
	cx -= s.cameraX
	cy -= s.cameraY

	var points [][]int
	for i := 0; i < 6; i++ {
		angle := float64(i)*2*math.Pi/6 - math.Pi/2
		x := cx + int(float64(radius)*math.Cos(angle))
		y := cy + int(float64(radius)*math.Sin(angle))
		points = append(points, []int{x, y})
	}

	if filled {
		s.fillPolygon(points, c)
	} else {
		for i := 0; i < len(points); i++ {
			next := (i + 1) % len(points)
			s.Line(points[i][0]+s.cameraX, points[i][1]+s.cameraY, points[next][0]+s.cameraX, points[next][1]+s.cameraY, c)
		}
	}
}

func (s *Soft) Star(cx, cy, radius int, filled bool, c color.RGBA) {
	// Apply camera offset
	cx -= s.cameraX
	cy -= s.cameraY

	outerRadius := radius
	innerRadius := radius / 2
	var points [][]int
	for i := 0; i < 10; i++ {
		angle := float64(i)*math.Pi/5 - math.Pi/2
		r := outerRadius
		if i%2 == 1 {
			r = innerRadius
		}
		x := cx + int(float64(r)*math.Cos(angle))
		y := cy + int(float64(r)*math.Sin(angle))
		points = append(points, []int{x, y})
	}

	if filled {
		s.fillPolygon(points, c)
	} else {
		for i := 0; i < len(points); i++ {
			next := (i + 1) % len(points)
			s.Line(points[i][0]+s.cameraX, points[i][1]+s.cameraY, points[next][0]+s.cameraX, points[next][1]+s.cameraY, c)
		}
	}
}
