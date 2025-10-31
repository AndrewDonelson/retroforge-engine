package engine

import (
	"bytes"
	"image/color"
	"testing"
	"time"

	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
)

// use fake clock to avoid sleeping
type itClock struct{ now time.Time }

func (c *itClock) Now() time.Time        { return c.now }
func (c *itClock) Sleep(d time.Duration) { c.now = c.now.Add(d) }

func TestHelloCartRendersCenteredText(t *testing.T) {
	// Build a minimal cart
	m := cartio.Manifest{Title: "Hello", Author: "RF", Entry: "main.lua"}
	lua := `
        function _INIT() 
            rf.palette_set("default")
        end
        function _UPDATE(dt) end
        function _DRAW()
            rf.clear_i(0)
            rf.print_anchored("HELLO FROM RETROFORGE", "middlecenter", 1)
        end
    `
	var buf bytes.Buffer
	if err := cartio.Write(&buf, m, []cartio.Asset{{Name: "main.lua", Data: []byte(lua)}}, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap)); err != nil {
		t.Fatal(err)
	}

	e := New(60)
	t.Cleanup(e.Close)
	e.Sched.WithClock(&itClock{now: time.Unix(0, 0)})

	// Load cart from memory
	if err := e.LoadCartFromReader(bytes.NewReader(buf.Bytes()), int64(buf.Len())); err != nil {
		t.Fatal(err)
	}
	// Run one frame
	e.RunFrames(1)

	// Assert pixels near vertical center are non-black (text drawn)
	pix := e.Ren.Pixels()
	w := e.Ren.Width()
	h := e.Ren.Height()
	y := 135 + 2 // within glyph height
	nonBlack := 0
	for x := w/2 - 40; x < w/2+40; x++ { // scan ~80px across center
		idx := (y*w + x) * 4
		r, g, b := pix[idx+0], pix[idx+1], pix[idx+2]
		if !(r == 0 && g == 0 && b == 0) {
			nonBlack++
		}
	}
	if nonBlack == 0 {
		t.Fatalf("expected some text pixels drawn at y=%d", y)
	}
	_ = h // silence unused if h not used in future
	_ = color.RGBA{}
}
