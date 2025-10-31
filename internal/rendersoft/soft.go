package rendersoft

import (
	"image/color"

	"github.com/AndrewDonelson/retroforge-engine/internal/font"
)

type Soft struct {
	w, h                       int
	pix                        []uint8 // RGBA
	clipX, clipY, clipW, clipH int     // Clipping rectangle (0,0,0,0 = disabled)
	cameraX, cameraY           int     // Camera offset
}

func New(w, h int) *Soft { return &Soft{w: w, h: h, pix: make([]uint8, w*h*4)} }

func (s *Soft) Width() int      { return s.w }
func (s *Soft) Height() int     { return s.h }
func (s *Soft) Pixels() []uint8 { return s.pix }

func (s *Soft) Clear(c color.RGBA) {
	for i := 0; i < len(s.pix); i += 4 {
		s.pix[i+0] = c.R
		s.pix[i+1] = c.G
		s.pix[i+2] = c.B
		s.pix[i+3] = 0xFF
	}
}

func (s *Soft) set(x, y int, c color.RGBA) {
	// Camera offset is applied by callers (functions pass world coordinates)
	// Check bounds
	if x < 0 || y < 0 || x >= s.w || y >= s.h {
		return
	}

	// Check clipping
	if s.clipW > 0 && s.clipH > 0 {
		if x < s.clipX || y < s.clipY || x >= s.clipX+s.clipW || y >= s.clipY+s.clipH {
			return
		}
	}

	idx := (y*s.w + x) * 4
	s.pix[idx+0] = c.R
	s.pix[idx+1] = c.G
	s.pix[idx+2] = c.B
	s.pix[idx+3] = 0xFF
}

func (s *Soft) PSet(x, y int, c color.RGBA) {
	// Apply camera offset
	x -= s.cameraX
	y -= s.cameraY
	s.set(x, y, c)
}

func (s *Soft) PGet(x, y int) color.RGBA {
	if x < 0 || y < 0 || x >= s.w || y >= s.h {
		return color.RGBA{0, 0, 0, 0} // Return transparent for out of bounds
	}
	idx := (y*s.w + x) * 4
	return color.RGBA{
		R: s.pix[idx+0],
		G: s.pix[idx+1],
		B: s.pix[idx+2],
		A: s.pix[idx+3],
	}
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func (s *Soft) Line(x0, y0, x1, y1 int, c color.RGBA) {
	// Apply camera offset
	x0 -= s.cameraX
	y0 -= s.cameraY
	x1 -= s.cameraX
	y1 -= s.cameraY

	dx := abs(x1 - x0)
	sx := -1
	if x0 < x1 {
		sx = 1
	}
	dy := -abs(y1 - y0)
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	err := dx + dy
	for {
		s.set(x0, y0, c)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

func (s *Soft) Rect(x0, y0, x1, y1 int, c color.RGBA) {
	// Apply camera offset (Line already applies it, but we need consistent behavior)
	s.Line(x0, y0, x1, y0, c)
	s.Line(x1, y0, x1, y1, c)
	s.Line(x1, y1, x0, y1, c)
	s.Line(x0, y1, x0, y0, c)
}

func (s *Soft) RectFill(x0, y0, x1, y1 int, c color.RGBA) {
	// Apply camera offset
	x0 -= s.cameraX
	y0 -= s.cameraY
	x1 -= s.cameraX
	y1 -= s.cameraY

	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			s.set(x, y, c)
		}
	}
}

func (s *Soft) Circ(xc, yc, r int, c color.RGBA) {
	// Apply camera offset
	xc -= s.cameraX
	yc -= s.cameraY

	x, y, d := r, 0, 1-2*r
	for y <= x {
		s.set(xc+x, yc+y, c)
		s.set(xc+y, yc+x, c)
		s.set(xc-x, yc+y, c)
		s.set(xc-y, yc+x, c)
		s.set(xc-x, yc-y, c)
		s.set(xc-y, yc-x, c)
		s.set(xc+x, yc-y, c)
		s.set(xc+y, yc-x, c)
		if d < 0 {
			d += 2*y + 1
		} else {
			d += 2*(y-x) + 1
			x--
		}
		y++
	}
}

func (s *Soft) CircFill(xc, yc, r int, c color.RGBA) {
	// Apply camera offset
	xc -= s.cameraX
	yc -= s.cameraY

	x, y, d := r, 0, 1-2*r
	for y <= x {
		for xi := xc - x; xi <= xc+x; xi++ {
			s.set(xi, yc+y, c)
			s.set(xi, yc-y, c)
		}
		if d < 0 {
			d += 2*y + 1
		} else {
			d += 2*(y-x) + 1
			x--
		}
		y++
	}
}

func (s *Soft) Print(text string, x, y int, c color.RGBA) {
	// Apply camera offset
	x -= s.cameraX
	y -= s.cameraY
	cx := x
	for _, r := range text {
		if r == '\n' {
			y += font.Height + 1
			cx = x
			continue
		}
		g, ok := font.Get(r)
		if !ok {
			cx += font.Advance
			continue
		}
		for row := 0; row < g.H; row++ {
			bits := g.Rows[row]
			for col := 0; col < g.W; col++ {
				// 5x7 glyphs use the lower 5 bits with LEFTMOST at bit 4
				// so column 0 -> bit 4, column 4 -> bit 0
				bit := uint((font.Width - 1) - col)
				if (bits & (1 << bit)) != 0 {
					s.set(cx+col, y+row, c)
				}
			}
		}
		cx += font.Advance
	}
}

func (s *Soft) PrintCentered(text string, y int, c color.RGBA) {
	w := len([]rune(text)) * font.Advance
	x := (s.w - w) / 2
	if x < 0 {
		x = 0
	}
	s.Print(text, x, y, c)
}

func (s *Soft) PrintAnchored(text string, anchor string, c color.RGBA) {
	textW := len([]rune(text)) * font.Advance
	textH := font.Height

	var x, y int

	// Horizontal alignment
	switch anchor {
	case "topleft", "middleleft", "bottomleft":
		x = 0
	case "topcenter", "middlecenter", "bottomcenter":
		x = (s.w - textW) / 2
		if x < 0 {
			x = 0
		}
	case "topright", "middleright", "bottomright":
		x = s.w - textW
		if x < 0 {
			x = 0
		}
	default:
		// Default to topcenter for invalid anchors
		x = (s.w - textW) / 2
		if x < 0 {
			x = 0
		}
		y = 0
		s.Print(text, x, y, c)
		return
	}

	// Vertical alignment
	switch anchor {
	case "topleft", "topcenter", "topright":
		y = 0
	case "middleleft", "middlecenter", "middleright":
		y = (s.h - textH) / 2
		if y < 0 {
			y = 0
		}
	case "bottomleft", "bottomcenter", "bottomright":
		y = s.h - textH
		if y < 0 {
			y = 0
		}
	}

	s.Print(text, x, y, c)
}
