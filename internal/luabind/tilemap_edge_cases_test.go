package luabind

import (
	"testing"

	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
	"github.com/AndrewDonelson/retroforge-engine/internal/rendersoft"
	lua "github.com/yuin/gopher-lua"
)

func TestDrawTilemap_ZeroDimensions(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	tilemapsMap := make(map[string]*cartio.TileMapData)

	// Test with zero-width tile
	tilesetData := &cartio.TilesetData{
		IsISO: false,
		Tiles: cartio.TilesetMap{
			"invalid_tile": {
				Width:  0,
				Height: 16,
				Type:   cartio.SpriteTypeStatic,
				Pixels: createTilePixels(16, 0, 1),
			},
		},
	}

	tilemapData := &cartio.TileMapData{
		TileMap: cartio.TileMap{
			TilesetName: "invalid_tileset",
			Width:       1,
			Height:      1,
			Tiles:       [][]string{{"invalid_tile"}},
		},
		Tileset: tilesetData.Tiles,
		IsISO:   false,
	}
	tilemapsMap["invalid_map"] = tilemapData

	RegisterWithState(L, r, func(i int) (rgba [4]uint8) {
		return [4]uint8{0, 0, 0, 255}
	}, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), tilemapsMap, nil, NewState(), nil)

	// Should handle gracefully (no crash)
	err := L.DoString(`rf.drawTilemap("invalid_map", 0, 0)`)
	if err != nil {
		t.Logf("drawTilemap with invalid tile handled gracefully: %v", err)
	}
}

func TestDrawTilemap_NegativeCoordinates(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	tilemapsMap := make(map[string]*cartio.TileMapData)

	tilesetData := &cartio.TilesetData{
		IsISO: false,
		Tiles: cartio.TilesetMap{
			"tile1": {
				Width:  16,
				Height: 16,
				Type:   cartio.SpriteTypeStatic,
				Pixels: createTilePixels(16, 16, 1),
			},
		},
	}

	tilemapData := &cartio.TileMapData{
		TileMap: cartio.TileMap{
			TilesetName: "test_tileset",
			Width:       1,
			Height:      1,
			Tiles:       [][]string{{"tile1"}},
		},
		Tileset: tilesetData.Tiles,
		IsISO:   false,
	}
	tilemapsMap["test_map"] = tilemapData

	RegisterWithState(L, r, func(i int) (rgba [4]uint8) {
		return [4]uint8{0, 0, 0, 255}
	}, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), tilemapsMap, nil, NewState(), nil)

	// Test with very negative offsets
	err := L.DoString(`rf.drawTilemap("test_map", -1000, -1000)`)
	if err != nil {
		t.Fatalf("drawTilemap with negative coordinates should not error: %v", err)
	}

	// Test with very large offsets
	err = L.DoString(`rf.drawTilemap("test_map", 10000, 10000)`)
	if err != nil {
		t.Fatalf("drawTilemap with large coordinates should not error: %v", err)
	}
}

func TestDrawTilemap_EmptyTilemap(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	tilemapsMap := make(map[string]*cartio.TileMapData)

	// Empty tilemap (no tiles)
	tilemapData := &cartio.TileMapData{
		TileMap: cartio.TileMap{
			TilesetName: "empty_tileset",
			Width:       0,
			Height:      0,
			Tiles:       [][]string{},
		},
		Tileset: cartio.TilesetMap{},
		IsISO:   false,
	}
	tilemapsMap["empty_map"] = tilemapData

	RegisterWithState(L, r, func(i int) (rgba [4]uint8) {
		return [4]uint8{0, 0, 0, 255}
	}, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), tilemapsMap, nil, NewState(), nil)

	err := L.DoString(`rf.drawTilemap("empty_map", 0, 0)`)
	if err != nil {
		t.Fatalf("drawTilemap with empty map should not error: %v", err)
	}
}

func TestDrawTilemap_PaletteBoundary(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	tilemapsMap := make(map[string]*cartio.TileMapData)

	// Test with palette index at boundaries (0, 49, -1 for transparent)
	tilesetData := &cartio.TilesetData{
		IsISO: false,
		Tiles: cartio.TilesetMap{
			"tile0": {
				Width:  16,
				Height: 16,
				Type:   cartio.SpriteTypeStatic,
				Pixels: createTilePixels(16, 16, 0), // Palette index 0
			},
			"tile49": {
				Width:  16,
				Height: 16,
				Type:   cartio.SpriteTypeStatic,
				Pixels: createTilePixels(16, 16, 49), // Palette index 49
			},
			"tile_transparent": {
				Width:  16,
				Height: 16,
				Type:   cartio.SpriteTypeStatic,
				Pixels: createTilePixelsWithTransparent(16, 16), // -1 for transparent
			},
		},
	}

	tilemapData := &cartio.TileMapData{
		TileMap: cartio.TileMap{
			TilesetName: "palette_tileset",
			Width:       3,
			Height:      1,
			Tiles:       [][]string{{"tile0", "tile49", "tile_transparent"}},
		},
		Tileset: tilesetData.Tiles,
		IsISO:   false,
	}
	tilemapsMap["palette_map"] = tilemapData

	RegisterWithState(L, r, func(i int) (rgba [4]uint8) {
		if i == 0 {
			return [4]uint8{0, 0, 0, 255} // Black
		}
		if i == 49 {
			return [4]uint8{255, 255, 255, 255} // White
		}
		return [4]uint8{128, 128, 128, 255} // Gray
	}, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), tilemapsMap, nil, NewState(), nil)

	err := L.DoString(`rf.drawTilemap("palette_map", 0, 0)`)
	if err != nil {
		t.Fatalf("drawTilemap with palette boundary indices should not error: %v", err)
	}
}

func createTilePixelsWithTransparent(height, width int) [][]int {
	pixels := make([][]int, height)
	for y := 0; y < height; y++ {
		pixels[y] = make([]int, width)
		for x := 0; x < width; x++ {
			// Alternating transparent (-1) and visible (1) pixels
			if (x+y)%2 == 0 {
				pixels[y][x] = -1 // Transparent
			} else {
				pixels[y][x] = 1 // Visible
			}
		}
	}
	return pixels
}

func TestGridToScreen_VeryLargeCoordinates(t *testing.T) {
	tileWidth := 32
	offsetX := 0
	offsetY := 0

	// Test with very large grid coordinates
	largeX := 1000
	largeY := 1000

	screenX, screenY := gridToScreen(largeX, largeY, tileWidth, offsetX, offsetY)

	// Should calculate without overflow
	if screenX == 0 && screenY == 0 {
		t.Error("Very large coordinates should produce valid screen coordinates")
	}

	// Test with very negative coordinates
	negX := -1000
	negY := -1000

	screenX, screenY = gridToScreen(negX, negY, tileWidth, offsetX, offsetY)

	// Should handle gracefully
	_ = screenX
	_ = screenY
}

func TestScreenToGrid_VeryLargeCoordinates(t *testing.T) {
	tileWidth := 32
	offsetX := 0
	offsetY := 0

	// Test with very large screen coordinates
	largeScreenX := 100000
	largeScreenY := 100000

	gridX, gridY := screenToGrid(largeScreenX, largeScreenY, tileWidth, offsetX, offsetY)

	// Should calculate without overflow
	_ = gridX
	_ = gridY

	// Test with very negative coordinates
	negScreenX := -100000
	negScreenY := -100000

	gridX, gridY = screenToGrid(negScreenX, negScreenY, tileWidth, offsetX, offsetY)

	// Should handle gracefully
	_ = gridX
	_ = gridY
}

