package cartio

// TileMap represents a tilemap that references a tileset
type TileMap struct {
	TilesetName string     `json:"tilesetName"` // Name of the tileset to use (without _tiles.json extension)
	Width       int        `json:"width"`       // Map width in tiles
	Height      int        `json:"height"`      // Map height in tiles
	Tiles       [][]string `json:"tiles"`       // 2D array of tile names (empty string = no tile)
}

// TileMapData represents a loaded tilemap with its tileset reference
type TileMapData struct {
	TileMap
	Tileset TilesetMap // The loaded tileset for this map
	IsISO   bool       // Whether this tilemap uses isometric rendering (from tileset)
	Seed    int        // Seed for deterministic tile variation (from tileset, only used for normal tiles)
}
