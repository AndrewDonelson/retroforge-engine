package engine

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
)

// loadTileset loads a tileset from a JSON file
// Tilesets can be normal (isISO=false) or isometric (isISO=true)
// Isometric tiles should be pre-rendered using the tile2iso tool
func (e *Engine) loadTileset(tilesetName, tilesetPath string) error {
	data, err := os.ReadFile(tilesetPath)
	if err != nil {
		return fmt.Errorf("failed to read tileset file: %w", err)
	}

	var tilesetData cartio.TilesetData
	if err := json.Unmarshal(data, &tilesetData); err != nil {
		return fmt.Errorf("failed to parse tileset JSON: %w", err)
	}

	// Normalize and validate all tiles
	// Check if this is an isometric tileset to apply appropriate size limits
	isIsometric := tilesetData.IsISO
	for tileName, tile := range tilesetData.Tiles {
		cartio.NormalizeTileData(&tile)
		if err := cartio.ValidateTileData(&tile, tileName, isIsometric); err != nil {
			return fmt.Errorf("tile '%s' validation error: %w", tileName, err)
		}

		tilesetData.Tiles[tileName] = tile
	}

	e.tilesetsMap[tilesetName] = &tilesetData
	return nil
}

// loadTilemap loads a tilemap from a JSON file and links it to its tileset
func (e *Engine) loadTilemap(tilemapName, tilemapPath string) error {
	data, err := os.ReadFile(tilemapPath)
	if err != nil {
		return fmt.Errorf("failed to read tilemap file: %w", err)
	}

	var tilemap cartio.TileMap
	if err := json.Unmarshal(data, &tilemap); err != nil {
		return fmt.Errorf("failed to parse tilemap JSON: %w", err)
	}

	// Debug logging
	if e.devMode != nil && e.devMode.IsEnabled() {
		e.devMode.AddDebugLog(fmt.Sprintf("loadTilemap: Loading tilemap '%s' from '%s', references tileset '%s'", tilemapName, tilemapPath, tilemap.TilesetName))
	}

	// Find and link the tileset
	tilesetData, exists := e.tilesetsMap[tilemap.TilesetName]
	if !exists {
		// List available tilesets for debugging
		availableTilesets := make([]string, 0, len(e.tilesetsMap))
		for name := range e.tilesetsMap {
			availableTilesets = append(availableTilesets, name)
		}
		if e.devMode != nil && e.devMode.IsEnabled() {
			e.devMode.AddDebugLog(fmt.Sprintf("loadTilemap: Tileset '%s' not found for tilemap '%s'. Available tilesets: %v", tilemap.TilesetName, tilemapName, availableTilesets))
		}
		return fmt.Errorf("tileset '%s' not found for tilemap '%s'", tilemap.TilesetName, tilemapName)
	}

	// Validate tilemap dimensions
	if len(tilemap.Tiles) != tilemap.Height {
		return fmt.Errorf("tilemap '%s': tile array height (%d) does not match map height (%d)", tilemapName, len(tilemap.Tiles), tilemap.Height)
	}
	for i, row := range tilemap.Tiles {
		if len(row) != tilemap.Width {
			return fmt.Errorf("tilemap '%s': row %d width (%d) does not match map width (%d)", tilemapName, i, len(row), tilemap.Width)
		}
		// Validate tile names exist in tileset
		for j, tileName := range row {
			if tileName != "" {
				if _, exists := tilesetData.Tiles[tileName]; !exists {
					return fmt.Errorf("tilemap '%s': tile '%s' at [%d][%d] not found in tileset", tilemapName, tileName, i, j)
				}
			}
		}
	}

	// Create TileMapData with linked tileset (store the tileset data, not just the map)
	tilemapData := &cartio.TileMapData{
		TileMap: tilemap,
		Tileset: tilesetData.Tiles, // Extract just the tiles map for backward compatibility
		IsISO:   tilesetData.IsISO, // Store the isometric flag
		Seed:    tilesetData.Seed,   // Store the seed for tile variation
	}

	e.tilemapsMap[tilemapName] = tilemapData
	
	// Debug logging
	if e.devMode != nil && e.devMode.IsEnabled() {
		e.devMode.AddDebugLog(fmt.Sprintf("loadTilemap: Successfully loaded tilemap '%s' (IsISO=%v, size=%dx%d)", tilemapName, tilemapData.IsISO, tilemap.Width, tilemap.Height))
	}
	
	return nil
}

// loadTilesetFromBytes loads a tileset from JSON byte data (for cart mode)
func (e *Engine) loadTilesetFromBytes(tilesetName string, data []byte) error {
	var tilesetData cartio.TilesetData
	if err := json.Unmarshal(data, &tilesetData); err != nil {
		return fmt.Errorf("failed to parse tileset JSON: %w", err)
	}

	// Normalize and validate all tiles
	// Check if this is an isometric tileset to apply appropriate size limits
	isIsometric := tilesetData.IsISO
	for tileName, tile := range tilesetData.Tiles {
		cartio.NormalizeTileData(&tile)
		if err := cartio.ValidateTileData(&tile, tileName, isIsometric); err != nil {
			return fmt.Errorf("tile '%s' validation error: %w", tileName, err)
		}

		tilesetData.Tiles[tileName] = tile
	}

	e.tilesetsMap[tilesetName] = &tilesetData
	return nil
}

// loadTilemapFromBytes loads a tilemap from JSON byte data (for cart mode)
func (e *Engine) loadTilemapFromBytes(tilemapName string, data []byte) error {
	var tilemap cartio.TileMap
	if err := json.Unmarshal(data, &tilemap); err != nil {
		return fmt.Errorf("failed to parse tilemap JSON: %w", err)
	}

	// Find and link the tileset
	tilesetData, exists := e.tilesetsMap[tilemap.TilesetName]
	if !exists {
		return fmt.Errorf("tileset '%s' not found for tilemap '%s'", tilemap.TilesetName, tilemapName)
	}

	// Validate tilemap dimensions
	if len(tilemap.Tiles) != tilemap.Height {
		return fmt.Errorf("tilemap '%s': tile array height (%d) does not match map height (%d)", tilemapName, len(tilemap.Tiles), tilemap.Height)
	}
	for i, row := range tilemap.Tiles {
		if len(row) != tilemap.Width {
			return fmt.Errorf("tilemap '%s': row %d width (%d) does not match map width (%d)", tilemapName, i, len(row), tilemap.Width)
		}
		// Validate tile names exist in tileset
		for j, tileName := range row {
			if tileName != "" {
				if _, exists := tilesetData.Tiles[tileName]; !exists {
					return fmt.Errorf("tilemap '%s': tile '%s' at [%d][%d] not found in tileset", tilemapName, tileName, i, j)
				}
			}
		}
	}

	// Create TileMapData with linked tileset
	tilemapData := &cartio.TileMapData{
		TileMap: tilemap,
		Tileset: tilesetData.Tiles, // Extract just the tiles map for backward compatibility
		IsISO:   tilesetData.IsISO, // Store the isometric flag
		Seed:    tilesetData.Seed,   // Store the seed for tile variation
	}

	e.tilemapsMap[tilemapName] = tilemapData
	return nil
}

