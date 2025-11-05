package engine

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
)

// loadTileset loads a tileset from a JSON file
func (e *Engine) loadTileset(tilesetName, tilesetPath string) error {
	data, err := os.ReadFile(tilesetPath)
	if err != nil {
		return fmt.Errorf("failed to read tileset file: %w", err)
	}

	var tileset cartio.TilesetMap
	if err := json.Unmarshal(data, &tileset); err != nil {
		return fmt.Errorf("failed to parse tileset JSON: %w", err)
	}

	// Normalize and validate all tiles
	for tileName, tile := range tileset {
		cartio.NormalizeTileData(&tile)
		if err := cartio.ValidateTileData(&tile, tileName); err != nil {
			return fmt.Errorf("tile '%s' validation error: %w", tileName, err)
		}

		// If isISO is true, convert tile to isometric
		if tile.IsISO {
			if err := e.convertTileToIsometric(&tile); err != nil {
				return fmt.Errorf("failed to convert tile '%s' to isometric: %w", tileName, err)
			}
		}

		tileset[tileName] = tile
	}

	e.tilesetsMap[tilesetName] = tileset
	return nil
}

// convertTileToIsometric is implemented in tileset_iso.go

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

	// Find and link the tileset
	tileset, exists := e.tilesetsMap[tilemap.TilesetName]
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
				if _, exists := tileset[tileName]; !exists {
					return fmt.Errorf("tilemap '%s': tile '%s' at [%d][%d] not found in tileset", tilemapName, tileName, i, j)
				}
			}
		}
	}

	// Create TileMapData with linked tileset
	tilemapData := &cartio.TileMapData{
		TileMap: tilemap,
		Tileset: tileset,
	}

	e.tilemapsMap[tilemapName] = tilemapData
	return nil
}

