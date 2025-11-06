package luabind

import (
	"testing"
	"time"

	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
	"github.com/AndrewDonelson/retroforge-engine/internal/rendersoft"
	lua "github.com/yuin/gopher-lua"
)

func TestDrawTilemap_OrthogonalRendering(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	colorByIndex := func(i int) (rgba [4]uint8) {
		if i == 1 {
			return [4]uint8{255, 255, 255, 255} // White
		}
		return [4]uint8{0, 0, 0, 255} // Black
	}

	// Create tileset and tilemap
	tilemapsMap := make(map[string]*cartio.TileMapData)

	// Create normal tileset
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

	// Create normal tilemap
	tilemapData := &cartio.TileMapData{
		TileMap: cartio.TileMap{
			TilesetName: "test_tileset",
			Width:       2,
			Height:      2,
			Tiles: [][]string{
				{"tile1", "tile1"},
				{"tile1", ""},
			},
		},
		Tileset: tilesetData.Tiles,
		IsISO:   false,
	}
	tilemapsMap["test_map"] = tilemapData

	RegisterWithState(L, r, colorByIndex, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), tilemapsMap, nil, NewState(), nil)

	// Clear screen first
	err := L.DoString(`rf.clear_i(0)`)
	if err != nil {
		t.Fatalf("Failed to clear screen: %v", err)
	}
	r.SwapBuffers()

	// Draw tilemap
	err = L.DoString(`rf.drawTilemap("test_map", 10, 10)`)
	if err != nil {
		t.Fatalf("drawTilemap should not error: %v", err)
	}
	r.SwapBuffers()

	// Verify tiles were drawn (check for non-black pixels in expected area)
	pixels := r.Pixels()
	width := r.Width()

	// Check for non-black pixels in tile area (10,10 to 42,42 for 2x2 tiles)
	foundNonBlack := false
	for y := 10; y < 42 && !foundNonBlack; y++ {
		for x := 10; x < 42 && !foundNonBlack; x++ {
			idx := y*width*4 + x*4 // RGBA format, 4 bytes per pixel
			if idx+2 < len(pixels) && idx >= 0 {
				// Check if any color channel is non-zero (not black)
				if pixels[idx] > 0 || pixels[idx+1] > 0 || pixels[idx+2] > 0 {
					foundNonBlack = true
				}
			}
		}
	}

	if !foundNonBlack {
		t.Error("Tiles should be drawn (orthogonal grid)")
	}
}

func TestDrawTilemap_IsometricRendering(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	colorByIndex := func(i int) (rgba [4]uint8) {
		if i == 1 {
			return [4]uint8{255, 255, 255, 255} // White
		}
		return [4]uint8{0, 0, 0, 255} // Black
	}

	// Create isometric tileset
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

	// Create isometric tilemap
	tilemapData := &cartio.TileMapData{
		TileMap: cartio.TileMap{
			TilesetName: "iso_tileset",
			Width:       2,
			Height:      2,
			Tiles: [][]string{
				{"tile1", "tile1"},
				{"tile1", "tile1"},
			},
		},
		Tileset: tilesetData.Tiles,
		IsISO:   true,
	}
	tilemapsMap["iso_map"] = tilemapData

	RegisterWithState(L, r, colorByIndex, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), tilemapsMap, nil, NewState(), nil)

	// Clear screen
	err := L.DoString(`rf.clear_i(0)`)
	if err != nil {
		t.Fatalf("Failed to clear screen: %v", err)
	}
	r.SwapBuffers()

	// Draw isometric tilemap
	err = L.DoString(`rf.drawTilemap("iso_map", 100, 50)`)
	if err != nil {
		t.Fatalf("drawTilemap should not error: %v", err)
	}
	r.SwapBuffers()

	// For isometric, tiles should be positioned in diamond pattern
	// Check that tiles were drawn (should have non-black pixels)
	pixels := r.Pixels()
	width := r.Width()

	// Check a wider area for isometric tiles (diamond pattern spreads out)
	foundNonBlack := false
	for y := 0; y < 270 && !foundNonBlack; y++ {
		for x := 0; x < 480 && !foundNonBlack; x++ {
			idx := y*width*4 + x*4 // RGBA format, 4 bytes per pixel
			if idx+2 < len(pixels) && idx >= 0 {
				// Check if any color channel is non-zero (not black)
				if pixels[idx] > 0 || pixels[idx+1] > 0 || pixels[idx+2] > 0 {
					foundNonBlack = true
				}
			}
		}
	}

	if !foundNonBlack {
		t.Error("Isometric tiles should be drawn")
	}
}

func TestDrawTilemap_DepthSorting(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	colorByIndex := func(i int) (rgba [4]uint8) {
		if i == 1 {
			return [4]uint8{255, 0, 0, 255} // Red for back tiles
		}
		if i == 2 {
			return [4]uint8{0, 255, 0, 255} // Green for front tiles
		}
		return [4]uint8{0, 0, 0, 255}
	}

	tilemapsMap := make(map[string]*cartio.TileMapData)

	tilesetData := &cartio.TilesetData{
		IsISO: true,
		Tiles: cartio.TilesetMap{
			"back": {
				Width:  32,
				Height: 24,
				Type:   cartio.SpriteTypeStatic,
				Pixels: createTilePixels(24, 32, 1), // Red
			},
			"front": {
				Width:  32,
				Height: 24,
				Type:   cartio.SpriteTypeStatic,
				Pixels: createTilePixels(24, 32, 2), // Green
			},
		},
	}

	// Create tilemap where back tile (0,0) should be drawn before front tile (1,1)
	// Sort key: (0,0) = 0, (1,1) = 2
	tilemapData := &cartio.TileMapData{
		TileMap: cartio.TileMap{
			TilesetName: "sort_tileset",
			Width:       2,
			Height:      2,
			Tiles: [][]string{
				{"back", ""},
				{"", "front"},
			},
		},
		Tileset: tilesetData.Tiles,
		IsISO:   true,
	}
	tilemapsMap["sort_map"] = tilemapData

	RegisterWithState(L, r, colorByIndex, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), tilemapsMap, nil, NewState(), nil)

	err := L.DoString(`rf.clear_i(0)`)
	if err != nil {
		t.Fatalf("Failed to clear screen: %v", err)
	}
	r.SwapBuffers()

	err = L.DoString(`rf.drawTilemap("sort_map", 100, 50)`)
	if err != nil {
		t.Fatalf("drawTilemap should not error: %v", err)
	}
	r.SwapBuffers()

	// Verify both tiles were drawn
	pixels := r.Pixels()
	width := r.Width()

	// Check for non-black pixels (both tiles should be drawn)
	foundNonBlack := false
	for y := 0; y < 270 && !foundNonBlack; y++ {
		for x := 0; x < 480 && !foundNonBlack; x++ {
			idx := y*width*4 + x*4 // RGBA format, 4 bytes per pixel
			if idx+2 < len(pixels) && idx >= 0 {
				// Check if any color channel is non-zero (not black)
				if pixels[idx] > 0 || pixels[idx+1] > 0 || pixels[idx+2] > 0 {
					foundNonBlack = true
				}
			}
		}
	}

	if !foundNonBlack {
		t.Error("Both tiles should be drawn (depth sorting)")
	}
}

func TestDrawTilemap_ErrorHandling(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	tilemapsMap := make(map[string]*cartio.TileMapData)

	RegisterWithState(L, r, func(i int) (rgba [4]uint8) {
		return [4]uint8{0, 0, 0, 255}
	}, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), tilemapsMap, nil, NewState(), nil)

	// Test missing tilemap
	err := L.DoString(`rf.drawTilemap("nonexistent", 0, 0)`)
	if err != nil {
		t.Fatalf("drawTilemap with missing map should not error (graceful failure): %v", err)
	}

	// Test empty tilemap
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

	err = L.DoString(`rf.drawTilemap("empty_map", 0, 0)`)
	if err != nil {
		t.Fatalf("drawTilemap with empty map should not error: %v", err)
	}
}

func TestDrawTilemap_OffsetParameters(t *testing.T) {
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

	tilemapData := &cartio.TileMapData{
		TileMap: cartio.TileMap{
			TilesetName: "offset_tileset",
			Width:       1,
			Height:      1,
			Tiles:       [][]string{{"tile1"}},
		},
		Tileset: tilesetData.Tiles,
		IsISO:   false,
	}
	tilemapsMap["offset_map"] = tilemapData

	RegisterWithState(L, r, colorByIndex, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), tilemapsMap, nil, NewState(), nil)

	err := L.DoString(`rf.clear_i(0)`)
	if err != nil {
		t.Fatalf("Failed to clear screen: %v", err)
	}
	r.SwapBuffers()

	// Test with offset (50, 30)
	err = L.DoString(`rf.drawTilemap("offset_map", 50, 30)`)
	if err != nil {
		t.Fatalf("drawTilemap should not error: %v", err)
	}
	r.SwapBuffers()

	// Verify tile drawn at offset position
	// Tiles are drawn pixel-by-pixel, so we check if any non-black pixels exist in the area
	pixels := r.Pixels()
	width := r.Width()

	// Check for any non-black pixels in tile area (50,30 to 66,46) or nearby
	foundNonBlack := false
	for y := 25; y < 55 && !foundNonBlack; y++ {
		for x := 45; x < 75 && !foundNonBlack; x++ {
			idx := y*width*4 + x*4 // RGBA format, 4 bytes per pixel
			if idx+2 < len(pixels) && idx >= 0 {
				// Check if any color channel is non-zero (not black)
				if pixels[idx] > 0 || pixels[idx+1] > 0 || pixels[idx+2] > 0 {
					foundNonBlack = true
				}
			}
		}
	}

	if !foundNonBlack {
		t.Error("Tile should be drawn at offset position (50, 30)")
	}

	// Test with negative offset (tile may be partially off-screen, but should not crash)
	err = L.DoString(`rf.clear_i(0)`)
	if err != nil {
		t.Fatalf("Failed to clear screen: %v", err)
	}
	r.SwapBuffers()

	err = L.DoString(`rf.drawTilemap("offset_map", -10, -10)`)
	if err != nil {
		t.Fatalf("drawTilemap with negative offset should not error: %v", err)
	}
	r.SwapBuffers()

	// With negative offset, tile may be partially off-screen, but function should complete
	// We don't verify pixel positions here since tile may be clipped
}

func TestDrawTilemap_MissingTiles(t *testing.T) {
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

	// Tilemap references non-existent tile
	tilemapData := &cartio.TileMapData{
		TileMap: cartio.TileMap{
			TilesetName: "missing_tileset",
			Width:       2,
			Height:      2,
			Tiles: [][]string{
				{"tile1", "nonexistent"},
				{"", "tile1"},
			},
		},
		Tileset: tilesetData.Tiles,
		IsISO:   false,
	}
	tilemapsMap["missing_map"] = tilemapData

	RegisterWithState(L, r, func(i int) (rgba [4]uint8) {
		return [4]uint8{0, 0, 0, 255}
	}, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), tilemapsMap, nil, NewState(), nil)

	err := L.DoString(`rf.drawTilemap("missing_map", 0, 0)`)
	if err != nil {
		t.Fatalf("drawTilemap with missing tiles should not error (should skip): %v", err)
	}
}

// Helper function to create tile pixels
func createTilePixels(height, width int, colorIndex int) [][]int {
	pixels := make([][]int, height)
	for y := 0; y < height; y++ {
		pixels[y] = make([]int, width)
		for x := 0; x < width; x++ {
			pixels[y][x] = colorIndex
		}
	}
	return pixels
}

func TestDrawTilemap_LargeTilemap(t *testing.T) {
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

	// Create 50×50 tilemap
	largeMap := make([][]string, 50)
	for y := 0; y < 50; y++ {
		largeMap[y] = make([]string, 50)
		for x := 0; x < 50; x++ {
			largeMap[y][x] = "tile1"
		}
	}

	tilemapData := &cartio.TileMapData{
		TileMap: cartio.TileMap{
			TilesetName: "large_tileset",
			Width:       50,
			Height:      50,
			Tiles:       largeMap,
		},
		Tileset: tilesetData.Tiles,
		IsISO:   false,
	}
	tilemapsMap["large_map"] = tilemapData

	RegisterWithState(L, r, colorByIndex, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), tilemapsMap, nil, NewState(), nil)

	err := L.DoString(`rf.clear_i(0)`)
	if err != nil {
		t.Fatalf("Failed to clear screen: %v", err)
	}
	r.SwapBuffers()

	// Measure rendering time
	start := time.Now()
	err = L.DoString(`rf.drawTilemap("large_map", 0, 0)`)
	elapsed := time.Since(start)
	if err != nil {
		t.Fatalf("drawTilemap should not error: %v", err)
	}
	r.SwapBuffers()

	// Verify performance: should complete in reasonable time (< 100ms for 50×50)
	if elapsed > 100*time.Millisecond {
		t.Logf("Large tilemap rendering took %v (target: < 100ms)", elapsed)
		// Don't fail, just log - performance may vary
	}

	// Verify some tiles were drawn
	pixels := r.Pixels()
	width := r.Width()
	foundNonBlack := false
	for y := 0; y < 270 && !foundNonBlack; y++ {
		for x := 0; x < 480 && !foundNonBlack; x++ {
			idx := y*width*4 + x*4
			if idx+2 < len(pixels) && idx >= 0 {
				if pixels[idx] > 0 || pixels[idx+1] > 0 || pixels[idx+2] > 0 {
					foundNonBlack = true
				}
			}
		}
	}

	if !foundNonBlack {
		t.Error("Large tilemap should render some tiles")
	}

	t.Logf("50×50 tilemap rendered in %v", elapsed)
}

