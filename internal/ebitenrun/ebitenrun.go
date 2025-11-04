//go:build !js && !wasm

package ebitenrun

import (
	"image"
	"image/png"
	"os"
	"sync"
	"time"

	"github.com/AndrewDonelson/retroforge-engine/internal/app"
	"github.com/AndrewDonelson/retroforge-engine/internal/audio"
	"github.com/AndrewDonelson/retroforge-engine/internal/engine"
	"github.com/AndrewDonelson/retroforge-engine/internal/input"
	"github.com/hajimehoshi/ebiten/v2"
)

// Game implements ebiten.Game interface
type Game struct {
	engine           *engine.Engine
	scale            int
	screenshotPressed bool
	screenshotMu     sync.Mutex
}

// Update is called every frame (60 FPS)
func (g *Game) Update() error {
	// Step input state BEFORE checking keys (prev = cur, then we update cur)
	// This ensures btnp() works correctly by comparing current vs previous state
	input.Step()

	// Handle keyboard input - map Ebiten keys to 11-button system
	// Ebiten uses ebiten.IsKeyPressed for continuous key state
	// SELECT
	input.Set(input.BtnSELECT, ebiten.IsKeyPressed(ebiten.KeyEnter))

	// START
	input.Set(input.BtnSTART, ebiten.IsKeyPressed(ebiten.KeySpace))

	// Directions
	input.Set(input.BtnUP, ebiten.IsKeyPressed(ebiten.KeyArrowUp))
	input.Set(input.BtnDOWN, ebiten.IsKeyPressed(ebiten.KeyArrowDown))
	input.Set(input.BtnLEFT, ebiten.IsKeyPressed(ebiten.KeyArrowLeft))
	input.Set(input.BtnRIGHT, ebiten.IsKeyPressed(ebiten.KeyArrowRight))

	// Action buttons
	input.Set(input.BtnA, ebiten.IsKeyPressed(ebiten.KeyA))
	input.Set(input.BtnB, ebiten.IsKeyPressed(ebiten.KeyS))
	input.Set(input.BtnX, ebiten.IsKeyPressed(ebiten.KeyZ))
	input.Set(input.BtnY, ebiten.IsKeyPressed(ebiten.KeyX))

	// TURBO (modifier)
	input.Set(input.BtnTURBO, ebiten.IsKeyPressed(ebiten.KeyShiftLeft) || ebiten.IsKeyPressed(ebiten.KeyShiftRight))

	// Screenshot on PrintScreen (detect key press)
	g.screenshotMu.Lock()
	isPressed := ebiten.IsKeyPressed(ebiten.KeyPrintScreen)
	if isPressed && !g.screenshotPressed {
		saveScreenshot(g.engine)
		g.screenshotPressed = true
	} else if !isPressed {
		g.screenshotPressed = false
	}
	g.screenshotMu.Unlock()

	// Run one frame (now input state is correct: prev has old state, cur has new state)
	g.engine.RunFrames(1)
	if app.QuitRequested() {
		return ebiten.Termination
	}

	return nil
}

// Draw is called every frame to render
func (g *Game) Draw(screen *ebiten.Image) {
	// Get pixels from engine renderer
	pix := g.engine.Ren.Pixels()
	w := g.engine.Ren.Width()
	h := g.engine.Ren.Height()

	if len(pix) == 0 {
		return
	}

	// Create RGBA image from pixel buffer (ABGR format from engine)
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	copy(img.Pix, pix)

	// Convert to Ebiten image
	ebitenImg := ebiten.NewImageFromImage(img)

	// Draw with scaling
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(g.scale), float64(g.scale))
	screen.DrawImage(ebitenImg, op)
}

// Layout returns the game's logical screen size
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	w := g.engine.Ren.Width() * g.scale
	h := g.engine.Ren.Height() * g.scale
	return w, h
}

// RunWindow opens an Ebiten window and runs the game loop
func RunWindow(e *engine.Engine, scale int) error {
	if scale <= 0 {
		scale = 2
	}

	// Initialize audio
	_ = audio.Init()

	// Create game instance
	game := &Game{
		engine: e,
		scale:  scale,
	}

	// Configure window
	ebiten.SetWindowSize(e.Ren.Width()*scale, e.Ren.Height()*scale)
	ebiten.SetWindowTitle("RetroForge")
	ebiten.SetWindowResizable(false)

	// Run game loop (blocks until window is closed)
	return ebiten.RunGame(game)
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

