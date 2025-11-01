package spritepool

import (
	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
)

// createTestSpriteData creates a test sprite data for testing
func createTestSpriteData(name string, isUI bool, maxSpawn int, lifetime int) cartio.SpriteData {
	pixels := make([][]int, 8)
	for y := range pixels {
		pixels[y] = make([]int, 8)
		for x := range pixels[y] {
			pixels[y][x] = -1
		}
	}

	return cartio.SpriteData{
		Width:        8,
		Height:       8,
		Pixels:       pixels,
		UseCollision: false,
		MountPoints:  []cartio.MountPoint{},
		IsUI:         isUI,
		Lifetime:     lifetime,
		MaxSpawn:     maxSpawn,
	}
}
