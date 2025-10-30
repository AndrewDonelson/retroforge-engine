package rendersoft

import (
    "image/color"
    "github.com/AndrewDonelson/retroforge-engine/internal/font"
)

type Soft struct {
    w, h int
    pix  []uint8 // RGBA
}

func New(w, h int) *Soft { return &Soft{w:w, h:h, pix: make([]uint8, w*h*4)} }

func (s *Soft) Width() int { return s.w }
func (s *Soft) Height() int { return s.h }
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
    if x < 0 || y < 0 || x >= s.w || y >= s.h { return }
    idx := (y*s.w + x) * 4
    s.pix[idx+0] = c.R
    s.pix[idx+1] = c.G
    s.pix[idx+2] = c.B
    s.pix[idx+3] = 0xFF
}

func (s *Soft) PSet(x, y int, c color.RGBA) { s.set(x,y,c) }

func abs(a int) int { if a<0 { return -a }; return a }

func (s *Soft) Line(x0, y0, x1, y1 int, c color.RGBA) {
    dx := abs(x1-x0); sx := -1; if x0 < x1 { sx = 1 }
    dy := -abs(y1-y0); sy := -1; if y0 < y1 { sy = 1 }
    err := dx + dy
    for {
        s.set(x0,y0,c)
        if x0 == x1 && y0 == y1 { break }
        e2 := 2*err
        if e2 >= dy { err += dy; x0 += sx }
        if e2 <= dx { err += dx; y0 += sy }
    }
}

func (s *Soft) Rect(x0, y0, x1, y1 int, c color.RGBA) {
    s.Line(x0,y0,x1,y0,c)
    s.Line(x1,y0,x1,y1,c)
    s.Line(x1,y1,x0,y1,c)
    s.Line(x0,y1,x0,y0,c)
}

func (s *Soft) RectFill(x0, y0, x1, y1 int, c color.RGBA) {
    if x0 > x1 { x0, x1 = x1, x0 }
    if y0 > y1 { y0, y1 = y1, y0 }
    for y := y0; y <= y1; y++ {
        for x := x0; x <= x1; x++ { s.set(x,y,c) }
    }
}

func (s *Soft) Circ(xc, yc, r int, c color.RGBA) {
    x, y, d := r, 0, 1-2*r
    for y <= x {
        s.set(xc+x, yc+y, c); s.set(xc+y, yc+x, c)
        s.set(xc-x, yc+y, c); s.set(xc-y, yc+x, c)
        s.set(xc-x, yc-y, c); s.set(xc-y, yc-x, c)
        s.set(xc+x, yc-y, c); s.set(xc+y, yc-x, c)
        if d < 0 { d += 2*y + 1 } else { d += 2*(y - x) + 1; x-- }
        y++
    }
}

func (s *Soft) CircFill(xc, yc, r int, c color.RGBA) {
    x, y, d := r, 0, 1-2*r
    for y <= x {
        for xi := xc-x; xi <= xc+x; xi++ { s.set(xi, yc+y, c); s.set(xi, yc-y, c) }
        if d < 0 { d += 2*y + 1 } else { d += 2*(y - x) + 1; x-- }
        y++
    }
}

func (s *Soft) Print(text string, x, y int, c color.RGBA) {
    cx := x
    for _, r := range text {
        if r == '\n' { y += font.Height + 1; cx = x; continue }
        g, ok := font.Get(r)
        if !ok { cx += font.Advance; continue }
        for row := 0; row < g.H; row++ {
            bits := g.Rows[row]
            for col := 0; col < g.W; col++ {
                // 5x7 glyphs use the lower 5 bits with LEFTMOST at bit 4
                // so column 0 -> bit 4, column 4 -> bit 0
                bit := uint((font.Width-1) - col)
                if (bits & (1 << bit)) != 0 { s.set(cx+col, y+row, c) }
            }
        }
        cx += font.Advance
    }
}

func (s *Soft) PrintCentered(text string, y int, c color.RGBA) {
    w := len([]rune(text))*font.Advance
    x := (s.w - w)/2
    if x < 0 { x = 0 }
    s.Print(text, x, y, c)
}


