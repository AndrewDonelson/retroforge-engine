package cartio

import (
	"testing"
)

func TestSpriteDataDefaults(t *testing.T) {
	sprite := SpriteData{
		Width:        16,
		Height:       16,
		Pixels:       make([][]int, 16),
		UseCollision: false,
		IsUI:         true,
		Lifetime:     0,
		MaxSpawn:     0,
	}

	if sprite.IsUI != true {
		t.Error("IsUI should default to true")
	}
	if sprite.Lifetime != 0 {
		t.Error("Lifetime should default to 0")
	}
	if sprite.MaxSpawn != 0 {
		t.Error("MaxSpawn should default to 0")
	}
}

func TestSpriteMap(t *testing.T) {
	sprites := make(SpriteMap)

	sprite := SpriteData{
		Width:        8,
		Height:       8,
		Pixels:       make([][]int, 8),
		UseCollision: true,
		IsUI:         false,
		Lifetime:     1000,
		MaxSpawn:     10,
	}

	sprites["test"] = sprite

	if sprites["test"].Width != 8 {
		t.Error("Sprite width not set correctly")
	}
	if sprites["test"].IsUI != false {
		t.Error("Sprite IsUI not set correctly")
	}
	if sprites["test"].Lifetime != 1000 {
		t.Error("Sprite Lifetime not set correctly")
	}
	if sprites["test"].MaxSpawn != 10 {
		t.Error("Sprite MaxSpawn not set correctly")
	}
}
