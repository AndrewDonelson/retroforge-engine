package imgtool

import (
	"image"
	"image/color"
	"testing"
)

func TestToSprite_FullPipeline(t *testing.T) {
	// Create test image
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}

	// Create palette
	palette := &Palette{Colors: make([]string, 50)}
	palette.Colors[0] = "#000000"
	palette.Colors[1] = "#ffffff"
	palette.Colors[2] = "#ff0000"
	for i := 3; i < 50; i++ {
		palette.Colors[i] = "#808080"
	}

	opts := DefaultToSpriteOptions("test_sprite")
	opts.TargetWidth = 16
	opts.TargetHeight = 16
	opts.UseCollision = true

	sprite, err := ToSprite(img, palette, opts)
	if err != nil {
		t.Fatalf("ToSprite() error = %v", err)
	}

	if sprite.Width != 16 {
		t.Errorf("ToSprite() width = %v, want 16", sprite.Width)
	}
	if sprite.Height != 16 {
		t.Errorf("ToSprite() height = %v, want 16", sprite.Height)
	}
	if sprite.UseCollision != true {
		t.Errorf("ToSprite() UseCollision = %v, want true", sprite.UseCollision)
	}
	if len(sprite.Pixels) != 16 {
		t.Errorf("ToSprite() pixels height = %v, want 16", len(sprite.Pixels))
	}
	if len(sprite.Pixels[0]) != 16 {
		t.Errorf("ToSprite() pixels width = %v, want 16", len(sprite.Pixels[0]))
	}

	// Validate sprite
	if err := sprite.Validate(); err != nil {
		t.Errorf("ToSprite() invalid sprite: %v", err)
	}
}

func TestToSprite_WithAllProperties(t *testing.T) {
	img := createTestImage(32, 32, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	palette := &Palette{Colors: make([]string, 50)}
	for i := 0; i < 50; i++ {
		palette.Colors[i] = "#000000"
	}

	opts := DefaultToSpriteOptions("test")
	opts.TargetWidth = 16
	opts.TargetHeight = 16
	opts.UseCollision = true
	opts.IsUI = false
	opts.Lifetime = 2000
	opts.MaxSpawn = 10

	sprite, err := ToSprite(img, palette, opts)
	if err != nil {
		t.Fatalf("ToSprite() error = %v", err)
	}

	if sprite.Lifetime != 2000 {
		t.Errorf("ToSprite() Lifetime = %v, want 2000", sprite.Lifetime)
	}
	if sprite.MaxSpawn != 10 {
		t.Errorf("ToSprite() MaxSpawn = %v, want 10", sprite.MaxSpawn)
	}
}

