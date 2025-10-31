package graphics

import "image/color"

// Renderer defines minimal 2D drawing for text.
type Renderer interface {
	Width() int
	Height() int
	Clear(c color.RGBA)
	Print(text string, x, y int, c color.RGBA)
	PrintAnchored(text string, anchor string, c color.RGBA)
	// Pixels exposes the backbuffer for tests/snapshots.
	Pixels() []uint8 // RGBA length = width*height*4
	// Primitives
	PSet(x, y int, c color.RGBA)
	PGet(x, y int) color.RGBA // Get pixel color
	Line(x0, y0, x1, y1 int, c color.RGBA)
	Rect(x0, y0, x1, y1 int, c color.RGBA)
	RectFill(x0, y0, x1, y1 int, c color.RGBA)
	Circ(x, y, r int, c color.RGBA)
	CircFill(x, y, r int, c color.RGBA)
	Ellipse(x, y, rx, ry int, c color.RGBA)
	EllipseFill(x, y, rx, ry int, c color.RGBA)
	// Shape primitives
	Triangle(x, y, radius int, filled bool, c color.RGBA)
	Diamond(x, y, radius int, filled bool, c color.RGBA)
	Square(x, y, radius int, filled bool, c color.RGBA)
	Pentagon(x, y, radius int, filled bool, c color.RGBA)
	Hexagon(x, y, radius int, filled bool, c color.RGBA)
	Star(x, y, radius int, filled bool, c color.RGBA)
	// State management
	SetClip(x, y, w, h int)    // Set clipping rectangle (0,0,0,0 to disable)
	GetClip() (x, y, w, h int) // Get current clip rectangle
	SetCamera(x, y int)        // Set camera offset
	GetCamera() (x, y int)     // Get camera offset
}
