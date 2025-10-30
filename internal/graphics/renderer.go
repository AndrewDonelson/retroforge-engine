package graphics

import "image/color"

// Renderer defines minimal 2D drawing for text.
type Renderer interface {
    Width() int
    Height() int
    Clear(c color.RGBA)
    Print(text string, x, y int, c color.RGBA)
    PrintCentered(text string, y int, c color.RGBA)
    // Pixels exposes the backbuffer for tests/snapshots.
    Pixels() []uint8 // RGBA length = width*height*4
    // Primitives
    PSet(x, y int, c color.RGBA)
    Line(x0, y0, x1, y1 int, c color.RGBA)
    Rect(x0, y0, x1, y1 int, c color.RGBA)
    RectFill(x0, y0, x1, y1 int, c color.RGBA)
    Circ(x, y, r int, c color.RGBA)
    CircFill(x, y, r int, c color.RGBA)
}


