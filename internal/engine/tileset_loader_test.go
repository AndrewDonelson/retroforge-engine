package engine

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
)

func TestLoadTileset_NormalTileset(t *testing.T) {
	// Create temporary tileset JSON
	tmpDir := t.TempDir()
	tilesetPath := filepath.Join(tmpDir, "test_tiles.json")

	tilesetData := map[string]interface{}{
		"isISO": false,
		"tiles": map[string]interface{}{
			"grass": map[string]interface{}{
				"width":  16,
				"height": 16,
				"type":   "static",
				"pixels": createTestPixels(16, 16, 1),
			},
			"dirt": map[string]interface{}{
				"width":  16,
				"height": 16,
				"type":   "static",
				"pixels": createTestPixels(16, 16, 2),
			},
		},
	}

	data, err := json.Marshal(tilesetData)
	if err != nil {
		t.Fatalf("Failed to marshal tileset: %v", err)
	}

	if err := os.WriteFile(tilesetPath, data, 0644); err != nil {
		t.Fatalf("Failed to write tileset file: %v", err)
	}

	// Create engine and load tileset
	engine := &Engine{
		tilesetsMap: make(map[string]*cartio.TilesetData),
		devMode:    nil,
	}

	err = engine.loadTileset("test", tilesetPath)
	if err != nil {
		t.Fatalf("loadTileset() error = %v", err)
	}

	// Verify tileset loaded
	tileset, exists := engine.tilesetsMap["test"]
	if !exists {
		t.Fatal("Tileset not found in map")
	}

	if tileset.IsISO {
		t.Error("Expected isISO=false, got true")
	}

	if len(tileset.Tiles) != 2 {
		t.Errorf("Expected 2 tiles, got %d", len(tileset.Tiles))
	}

	// Verify tiles
	if _, exists := tileset.Tiles["grass"]; !exists {
		t.Error("Tile 'grass' not found")
	}
	if _, exists := tileset.Tiles["dirt"]; !exists {
		t.Error("Tile 'dirt' not found")
	}
}

func TestLoadTileset_IsometricTileset(t *testing.T) {
	// Create temporary isometric tileset JSON
	tmpDir := t.TempDir()
	tilesetPath := filepath.Join(tmpDir, "iso_tiles.json")

	tilesetData := map[string]interface{}{
		"isISO": true,
		"tiles": map[string]interface{}{
			"earth": map[string]interface{}{
				"width":  32,
				"height": 24,
				"type":   "static",
				"pixels": createTestPixels(24, 32, 3), // 32x24 tile (height x width for pixels)
			},
		},
	}

	data, err := json.Marshal(tilesetData)
	if err != nil {
		t.Fatalf("Failed to marshal tileset: %v", err)
	}

	if err := os.WriteFile(tilesetPath, data, 0644); err != nil {
		t.Fatalf("Failed to write tileset file: %v", err)
	}

	// Create engine and load tileset
	engine := &Engine{
		tilesetsMap: make(map[string]*cartio.TilesetData),
		devMode:    nil,
	}

	err = engine.loadTileset("iso", tilesetPath)
	if err != nil {
		t.Fatalf("loadTileset() error = %v", err)
	}

	// Verify tileset loaded
	tileset, exists := engine.tilesetsMap["iso"]
	if !exists {
		t.Fatal("Tileset not found in map")
	}

	if !tileset.IsISO {
		t.Error("Expected isISO=true, got false")
	}
}

func TestLoadTileset_BackwardCompatibility(t *testing.T) {
	// Test old format (tiles at root level, no isISO)
	tmpDir := t.TempDir()
	tilesetPath := filepath.Join(tmpDir, "old_tiles.json")

	tilesetData := map[string]interface{}{
		"grass": map[string]interface{}{
			"width":  16,
			"height": 16,
			"type":   "static",
			"pixels": createTestPixels(16, 16, 1),
		},
	}

	data, err := json.Marshal(tilesetData)
	if err != nil {
		t.Fatalf("Failed to marshal tileset: %v", err)
	}

	if err := os.WriteFile(tilesetPath, data, 0644); err != nil {
		t.Fatalf("Failed to write tileset file: %v", err)
	}

	// Create engine and load tileset
	engine := &Engine{
		tilesetsMap: make(map[string]*cartio.TilesetData),
		devMode:    nil,
	}

	err = engine.loadTileset("old", tilesetPath)
	if err != nil {
		t.Fatalf("loadTileset() error = %v", err)
	}

	// Verify tileset loaded with isISO=false (default)
	tileset, exists := engine.tilesetsMap["old"]
	if !exists {
		t.Fatal("Tileset not found in map")
	}

	if tileset.IsISO {
		t.Error("Expected isISO=false for old format, got true")
	}

	if len(tileset.Tiles) != 1 {
		t.Errorf("Expected 1 tile, got %d", len(tileset.Tiles))
	}
}

func TestLoadTileset_InvalidTile(t *testing.T) {
	tmpDir := t.TempDir()
	tilesetPath := filepath.Join(tmpDir, "invalid_tiles.json")

	tilesetData := map[string]interface{}{
		"isISO": false,
		"tiles": map[string]interface{}{
			"invalid": map[string]interface{}{
				"width":  1, // Too small
				"height": 1,
				"type":   "static",
				"pixels": createTestPixels(1, 1, 1),
			},
		},
	}

	data, err := json.Marshal(tilesetData)
	if err != nil {
		t.Fatalf("Failed to marshal tileset: %v", err)
	}

	if err := os.WriteFile(tilesetPath, data, 0644); err != nil {
		t.Fatalf("Failed to write tileset file: %v", err)
	}

	engine := &Engine{
		tilesetsMap: make(map[string]*cartio.TilesetData),
		devMode:    nil,
	}

	err = engine.loadTileset("invalid", tilesetPath)
	if err == nil {
		t.Error("Expected error for invalid tile, got nil")
	}
}

func TestLoadTilesetFromBytes(t *testing.T) {
	tilesetData := map[string]interface{}{
		"isISO": false,
		"tiles": map[string]interface{}{
			"tile1": map[string]interface{}{
				"width":  16,
				"height": 16,
				"type":   "static",
				"pixels": createTestPixels(16, 16, 1),
			},
		},
	}

	data, err := json.Marshal(tilesetData)
	if err != nil {
		t.Fatalf("Failed to marshal tileset: %v", err)
	}

	engine := &Engine{
		tilesetsMap: make(map[string]*cartio.TilesetData),
		devMode:    nil,
	}

	err = engine.loadTilesetFromBytes("bytes_tileset", data)
	if err != nil {
		t.Fatalf("loadTilesetFromBytes() error = %v", err)
	}

	// Verify tileset loaded
	tileset, exists := engine.tilesetsMap["bytes_tileset"]
	if !exists {
		t.Fatal("Tileset not found in map")
	}

	if len(tileset.Tiles) != 1 {
		t.Errorf("Expected 1 tile, got %d", len(tileset.Tiles))
	}
}

func TestLoadTilemap_NormalTilemap(t *testing.T) {
	tmpDir := t.TempDir()

	// First create a tileset
	tilesetPath := filepath.Join(tmpDir, "test_tiles.json")
	tilesetData := map[string]interface{}{
		"isISO": false,
		"tiles": map[string]interface{}{
			"grass": map[string]interface{}{
				"width":  16,
				"height": 16,
				"type":   "static",
				"pixels": createTestPixels(16, 16, 1),
			},
		},
	}

	data, err := json.Marshal(tilesetData)
	if err != nil {
		t.Fatalf("Failed to marshal tileset: %v", err)
	}

	if err := os.WriteFile(tilesetPath, data, 0644); err != nil {
		t.Fatalf("Failed to write tileset file: %v", err)
	}

	// Create tilemap
	tilemapPath := filepath.Join(tmpDir, "test_map.json")
	tilemapData := map[string]interface{}{
		"tilesetName": "test",
		"width":       2,
		"height":      2,
		"tiles": [][]string{
			{"grass", "grass"},
			{"grass", ""}, // Empty tile
		},
	}

	data, err = json.Marshal(tilemapData)
	if err != nil {
		t.Fatalf("Failed to marshal tilemap: %v", err)
	}

	if err := os.WriteFile(tilemapPath, data, 0644); err != nil {
		t.Fatalf("Failed to write tilemap file: %v", err)
	}

	// Create engine, load tileset first, then tilemap
	engine := &Engine{
		tilesetsMap: make(map[string]*cartio.TilesetData),
		tilemapsMap: make(map[string]*cartio.TileMapData),
		devMode:    nil,
	}

	err = engine.loadTileset("test", tilesetPath)
	if err != nil {
		t.Fatalf("loadTileset() error = %v", err)
	}

	err = engine.loadTilemap("test_map", tilemapPath)
	if err != nil {
		t.Fatalf("loadTilemap() error = %v", err)
	}

	// Verify tilemap loaded
	tilemap, exists := engine.tilemapsMap["test_map"]
	if !exists {
		t.Fatal("Tilemap not found in map")
	}

	if tilemap.Width != 2 || tilemap.Height != 2 {
		t.Errorf("Expected 2x2 tilemap, got %dx%d", tilemap.Width, tilemap.Height)
	}

	if tilemap.IsISO {
		t.Error("Expected isISO=false, got true")
	}

	// Verify tileset linked
	if len(tilemap.Tileset) != 1 {
		t.Errorf("Expected 1 tile in tileset, got %d", len(tilemap.Tileset))
	}
}

func TestLoadTilemap_IsometricTilemap(t *testing.T) {
	tmpDir := t.TempDir()

	// Create isometric tileset
	tilesetPath := filepath.Join(tmpDir, "iso_tiles.json")
	tilesetData := map[string]interface{}{
		"isISO": true,
		"tiles": map[string]interface{}{
			"earth": map[string]interface{}{
				"width":  32,
				"height": 24,
				"type":   "static",
				"pixels": createTestPixels(24, 32, 1),
			},
		},
	}

	data, err := json.Marshal(tilesetData)
	if err != nil {
		t.Fatalf("Failed to marshal tileset: %v", err)
	}

	if err := os.WriteFile(tilesetPath, data, 0644); err != nil {
		t.Fatalf("Failed to write tileset file: %v", err)
	}

	// Create tilemap
	tilemapPath := filepath.Join(tmpDir, "iso_map.json")
	tilemapData := map[string]interface{}{
		"tilesetName": "iso",
		"width":       2,
		"height":      2,
		"tiles": [][]string{
			{"earth", "earth"},
			{"earth", "earth"},
		},
	}

	data, err = json.Marshal(tilemapData)
	if err != nil {
		t.Fatalf("Failed to marshal tilemap: %v", err)
	}

	if err := os.WriteFile(tilemapPath, data, 0644); err != nil {
		t.Fatalf("Failed to write tilemap file: %v", err)
	}

	engine := &Engine{
		tilesetsMap: make(map[string]*cartio.TilesetData),
		tilemapsMap: make(map[string]*cartio.TileMapData),
		devMode:    nil,
	}

	err = engine.loadTileset("iso", tilesetPath)
	if err != nil {
		t.Fatalf("loadTileset() error = %v", err)
	}

	err = engine.loadTilemap("iso_map", tilemapPath)
	if err != nil {
		t.Fatalf("loadTilemap() error = %v", err)
	}

	// Verify isometric flag propagated
	tilemap, exists := engine.tilemapsMap["iso_map"]
	if !exists {
		t.Fatal("Tilemap not found in map")
	}

	if !tilemap.IsISO {
		t.Error("Expected isISO=true, got false")
	}
}

func TestLoadTilemap_MissingTileset(t *testing.T) {
	tmpDir := t.TempDir()
	tilemapPath := filepath.Join(tmpDir, "missing_map.json")

	tilemapData := map[string]interface{}{
		"tilesetName": "nonexistent",
		"width":       2,
		"height":      2,
		"tiles": [][]string{
			{"tile1", "tile1"},
			{"tile1", "tile1"},
		},
	}

	data, err := json.Marshal(tilemapData)
	if err != nil {
		t.Fatalf("Failed to marshal tilemap: %v", err)
	}

	if err := os.WriteFile(tilemapPath, data, 0644); err != nil {
		t.Fatalf("Failed to write tilemap file: %v", err)
	}

	engine := &Engine{
		tilesetsMap: make(map[string]*cartio.TilesetData),
		tilemapsMap: make(map[string]*cartio.TileMapData),
		devMode:    nil,
	}

	err = engine.loadTilemap("missing", tilemapPath)
	if err == nil {
		t.Error("Expected error for missing tileset, got nil")
	}
}

func TestLoadTilemap_InvalidTileReference(t *testing.T) {
	tmpDir := t.TempDir()

	// Create tileset
	tilesetPath := filepath.Join(tmpDir, "test_tiles.json")
	tilesetData := map[string]interface{}{
		"isISO": false,
		"tiles": map[string]interface{}{
			"grass": map[string]interface{}{
				"width":  16,
				"height": 16,
				"type":   "static",
				"pixels": createTestPixels(16, 16, 1),
			},
		},
	}

	data, err := json.Marshal(tilesetData)
	if err != nil {
		t.Fatalf("Failed to marshal tileset: %v", err)
	}

	if err := os.WriteFile(tilesetPath, data, 0644); err != nil {
		t.Fatalf("Failed to write tileset file: %v", err)
	}

	// Create tilemap with invalid tile reference
	tilemapPath := filepath.Join(tmpDir, "invalid_map.json")
	tilemapData := map[string]interface{}{
		"tilesetName": "test",
		"width":       2,
		"height":      2,
		"tiles": [][]string{
			{"grass", "nonexistent"},
			{"grass", "grass"},
		},
	}

	data, err = json.Marshal(tilemapData)
	if err != nil {
		t.Fatalf("Failed to marshal tilemap: %v", err)
	}

	if err := os.WriteFile(tilemapPath, data, 0644); err != nil {
		t.Fatalf("Failed to write tilemap file: %v", err)
	}

	engine := &Engine{
		tilesetsMap: make(map[string]*cartio.TilesetData),
		tilemapsMap: make(map[string]*cartio.TileMapData),
		devMode:    nil,
	}

	err = engine.loadTileset("test", tilesetPath)
	if err != nil {
		t.Fatalf("loadTileset() error = %v", err)
	}

	err = engine.loadTilemap("invalid", tilemapPath)
	if err == nil {
		t.Error("Expected error for invalid tile reference, got nil")
	}
}

func TestLoadTilemapFromBytes(t *testing.T) {
	// Create tileset first
	tmpDir := t.TempDir()
	tilesetPath := filepath.Join(tmpDir, "test_tiles.json")
	tilesetData := map[string]interface{}{
		"isISO": false,
		"tiles": map[string]interface{}{
			"tile1": map[string]interface{}{
				"width":  16,
				"height": 16,
				"type":   "static",
				"pixels": createTestPixels(16, 16, 1),
			},
		},
	}

	data, err := json.Marshal(tilesetData)
	if err != nil {
		t.Fatalf("Failed to marshal tileset: %v", err)
	}

	if err := os.WriteFile(tilesetPath, data, 0644); err != nil {
		t.Fatalf("Failed to write tileset file: %v", err)
	}

	// Create tilemap bytes
	tilemapData := map[string]interface{}{
		"tilesetName": "test",
		"width":       2,
		"height":      2,
		"tiles": [][]string{
			{"tile1", "tile1"},
			{"tile1", ""},
		},
	}

	tilemapBytes, err := json.Marshal(tilemapData)
	if err != nil {
		t.Fatalf("Failed to marshal tilemap: %v", err)
	}

	engine := &Engine{
		tilesetsMap: make(map[string]*cartio.TilesetData),
		tilemapsMap: make(map[string]*cartio.TileMapData),
		devMode:    nil,
	}

	err = engine.loadTileset("test", tilesetPath)
	if err != nil {
		t.Fatalf("loadTileset() error = %v", err)
	}

	err = engine.loadTilemapFromBytes("bytes_map", tilemapBytes)
	if err != nil {
		t.Fatalf("loadTilemapFromBytes() error = %v", err)
	}

	// Verify tilemap loaded
	tilemap, exists := engine.tilemapsMap["bytes_map"]
	if !exists {
		t.Fatal("Tilemap not found in map")
	}

	if tilemap.Width != 2 || tilemap.Height != 2 {
		t.Errorf("Expected 2x2 tilemap, got %dx%d", tilemap.Width, tilemap.Height)
	}
}

// Helper function to create test pixel data
func createTestPixels(height, width int, colorIndex int) [][]int {
	pixels := make([][]int, height)
	for y := 0; y < height; y++ {
		pixels[y] = make([]int, width)
		for x := 0; x < width; x++ {
			pixels[y][x] = colorIndex
		}
	}
	return pixels
}

