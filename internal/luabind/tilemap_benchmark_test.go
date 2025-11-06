package luabind

import (
	"testing"

	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
	"github.com/AndrewDonelson/retroforge-engine/internal/rendersoft"
	lua "github.com/yuin/gopher-lua"
)

func BenchmarkDrawTilemap_Orthogonal(b *testing.B) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	colorByIndex := func(i int) (rgba [4]uint8) {
		if i == 1 {
			return [4]uint8{255, 255, 255, 255}
		}
		return [4]uint8{0, 0, 0, 255}
	}

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

	// Create 10×10 tilemap for benchmark
	benchMap := make([][]string, 10)
	for y := 0; y < 10; y++ {
		benchMap[y] = make([]string, 10)
		for x := 0; x < 10; x++ {
			benchMap[y][x] = "tile1"
		}
	}

	tilemapData := &cartio.TileMapData{
		TileMap: cartio.TileMap{
			TilesetName: "bench_tileset",
			Width:       10,
			Height:      10,
			Tiles:       benchMap,
		},
		Tileset: tilesetData.Tiles,
		IsISO:   false,
	}
	tilemapsMap["bench_map"] = tilemapData

	RegisterWithState(L, r, colorByIndex, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), tilemapsMap, nil, NewState(), nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		L.DoString(`rf.clear_i(0)`)
		r.SwapBuffers()
		L.DoString(`rf.drawTilemap("bench_map", 0, 0)`)
		r.SwapBuffers()
	}
}

func BenchmarkDrawTilemap_Isometric(b *testing.B) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	colorByIndex := func(i int) (rgba [4]uint8) {
		if i == 1 {
			return [4]uint8{255, 255, 255, 255}
		}
		return [4]uint8{0, 0, 0, 255}
	}

	tilemapsMap := make(map[string]*cartio.TileMapData)

	tilesetData := &cartio.TilesetData{
		IsISO: true,
		Tiles: cartio.TilesetMap{
			"tile1": {
				Width:  32,
				Height: 24,
				Type:   cartio.SpriteTypeStatic,
				Pixels: createTilePixels(24, 32, 1),
			},
		},
	}

	// Create 10×10 isometric tilemap for benchmark
	benchMap := make([][]string, 10)
	for y := 0; y < 10; y++ {
		benchMap[y] = make([]string, 10)
		for x := 0; x < 10; x++ {
			benchMap[y][x] = "tile1"
		}
	}

	tilemapData := &cartio.TileMapData{
		TileMap: cartio.TileMap{
			TilesetName: "bench_iso_tileset",
			Width:       10,
			Height:      10,
			Tiles:       benchMap,
		},
		Tileset: tilesetData.Tiles,
		IsISO:   true,
	}
	tilemapsMap["bench_iso_map"] = tilemapData

	RegisterWithState(L, r, colorByIndex, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), tilemapsMap, nil, NewState(), nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		L.DoString(`rf.clear_i(0)`)
		r.SwapBuffers()
		L.DoString(`rf.drawTilemap("bench_iso_map", 0, 0)`)
		r.SwapBuffers()
	}
}

func BenchmarkDrawTilemap_Large(b *testing.B) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	colorByIndex := func(i int) (rgba [4]uint8) {
		if i == 1 {
			return [4]uint8{255, 255, 255, 255}
		}
		return [4]uint8{0, 0, 0, 255}
	}

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

	// Create 50×50 tilemap for benchmark
	largeMap := make([][]string, 50)
	for y := 0; y < 50; y++ {
		largeMap[y] = make([]string, 50)
		for x := 0; x < 50; x++ {
			largeMap[y][x] = "tile1"
		}
	}

	tilemapData := &cartio.TileMapData{
		TileMap: cartio.TileMap{
			TilesetName: "large_bench_tileset",
			Width:       50,
			Height:      50,
			Tiles:       largeMap,
		},
		Tileset: tilesetData.Tiles,
		IsISO:   false,
	}
	tilemapsMap["large_bench_map"] = tilemapData

	RegisterWithState(L, r, colorByIndex, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), tilemapsMap, nil, NewState(), nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		L.DoString(`rf.clear_i(0)`)
		r.SwapBuffers()
		L.DoString(`rf.drawTilemap("large_bench_map", 0, 0)`)
		r.SwapBuffers()
	}
}

func BenchmarkGridToScreen(b *testing.B) {
	tileWidth := 32
	offsetX := 100
	offsetY := 50

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gridX := i % 100
		gridY := (i * 2) % 100
		_, _ = gridToScreen(gridX, gridY, tileWidth, offsetX, offsetY)
	}
}

func BenchmarkScreenToGrid(b *testing.B) {
	tileWidth := 32
	offsetX := 100
	offsetY := 50

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		screenX := (i * 10) % 1000
		screenY := (i * 15) % 1000
		_, _ = screenToGrid(screenX, screenY, tileWidth, offsetX, offsetY)
	}
}

func BenchmarkCoordinateRoundTrip(b *testing.B) {
	tileWidth := 32
	offsetX := 100
	offsetY := 50

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gridX := i % 100
		gridY := (i * 2) % 100
		screenX, screenY := gridToScreen(gridX, gridY, tileWidth, offsetX, offsetY)
		_, _ = screenToGrid(screenX, screenY, tileWidth, offsetX, offsetY)
	}
}

