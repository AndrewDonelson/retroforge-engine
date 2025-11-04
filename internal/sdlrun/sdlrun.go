//go:build !js && !wasm

package sdlrun

import (
	"image"
	"image/png"
	"os"
	"time"
	"unsafe"

	"github.com/AndrewDonelson/retroforge-engine/internal/app"
	"github.com/AndrewDonelson/retroforge-engine/internal/audio"
	"github.com/AndrewDonelson/retroforge-engine/internal/engine"
	"github.com/AndrewDonelson/retroforge-engine/internal/input"
	"github.com/veandco/go-sdl2/sdl"
)

// RunWindow opens an SDL window, runs frames, ESC/Close to exit.
func RunWindow(e *engine.Engine, scale int) error {
	if scale <= 0 {
		scale = 2
	}
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return err
	}
	defer sdl.Quit()

	w := int32(e.Ren.Width() * scale)
	h := int32(e.Ren.Height() * scale)
	win, err := sdl.CreateWindow("RetroForge", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, w, h, sdl.WINDOW_SHOWN)
	if err != nil {
		return err
	}
	defer win.Destroy()

	ren, err := sdl.CreateRenderer(win, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		return err
	}
	defer ren.Destroy()
	ren.SetLogicalSize(int32(e.Ren.Width()), int32(e.Ren.Height()))

	tex, err := ren.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(e.Ren.Width()), int32(e.Ren.Height()))
	if err != nil {
		return err
	}
	defer tex.Destroy()

	_ = audio.Init()
	running := true
	for running {
		// Step input state BEFORE polling events (prev = cur, then we update cur)
		// This ensures btnp() works correctly by comparing current vs previous state
		// NOTE: This must happen before polling events, unlike WASM where input is set asynchronously
		input.Step()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch ev := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				// ESC quitting disabled - games handle quitting through their own menus
				// if ev.Type == sdl.KEYDOWN && ev.Keysym.Sym == sdl.K_ESCAPE {
				// 	running = false
				// }
				down := ev.Type == sdl.KEYDOWN
				if ev.Type == sdl.KEYDOWN && ev.Keysym.Sym == sdl.K_PRINTSCREEN {
					// Save screenshot
					saveScreenshot(e)
				}
				
				// Map keyboard keys to 11-button system
				switch ev.Keysym.Sym {
				// SELECT
				case sdl.K_RETURN:
					input.Set(input.BtnSELECT, down)
				// START
				case sdl.K_SPACE:
					input.Set(input.BtnSTART, down)
				// Directions
				case sdl.K_UP:
					input.Set(input.BtnUP, down)
				case sdl.K_DOWN:
					input.Set(input.BtnDOWN, down)
				case sdl.K_LEFT:
					input.Set(input.BtnLEFT, down)
				case sdl.K_RIGHT:
					input.Set(input.BtnRIGHT, down)
				// Action buttons
				case sdl.K_a:
					input.Set(input.BtnA, down)
				case sdl.K_s:
					input.Set(input.BtnB, down)
				case sdl.K_z:
					input.Set(input.BtnX, down)
				case sdl.K_x:
					input.Set(input.BtnY, down)
				// TURBO (modifier)
				case sdl.K_LSHIFT, sdl.K_RSHIFT:
					input.Set(input.BtnTURBO, down)
				}
			}
		}

		// Run one frame (now input state is correct: prev has old state, cur has new state)
		e.RunFrames(1)
		if app.QuitRequested() {
			running = false
		}
		// Upload pixels
		pix := e.Ren.Pixels()
		var ptr unsafe.Pointer
		if len(pix) > 0 {
			ptr = unsafe.Pointer(&pix[0])
		}
		if err := tex.Update(nil, ptr, e.Ren.Width()*4); err != nil {
			return err
		}
		if err := ren.Clear(); err != nil {
			return err
		}
		if err := ren.Copy(tex, nil, nil); err != nil {
			return err
		}
		ren.Present()
	}
	return nil
}

func saveScreenshot(e *engine.Engine) {
	pix := e.Ren.Pixels()
	w := e.Ren.Width()
	h := e.Ren.Height()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	copy(img.Pix, pix)

	// Generate filename with timestamp
	filename := time.Now().Format("screenshot-20060102-150405.png")
	f, err := os.Create(filename)
	if err != nil {
		return // silently fail
	}
	defer f.Close()
	_ = png.Encode(f, img)
}
