package graphics

// TileMap represents a 2D tilemap grid
type TileMap struct {
	width, height int
	tiles         []int // 1D array: tiles[y*width + x] = tile index
}

// NewTileMap creates a new tilemap of given dimensions
func NewTileMap(w, h int) *TileMap {
	return &TileMap{
		width:  w,
		height: h,
		tiles:  make([]int, w*h),
	}
}

// Get returns tile at (x, y), or 0 if out of bounds
func (tm *TileMap) Get(x, y int) int {
	if x < 0 || y < 0 || x >= tm.width || y >= tm.height {
		return 0
	}
	return tm.tiles[y*tm.width+x]
}

// Set sets tile at (x, y) to value v
func (tm *TileMap) Set(x, y, v int) {
	if x < 0 || y < 0 || x >= tm.width || y >= tm.height {
		return
	}
	tm.tiles[y*tm.width+x] = v
}

// Width returns tilemap width
func (tm *TileMap) Width() int { return tm.width }

// Height returns tilemap height
func (tm *TileMap) Height() int { return tm.height }

// Draw draws a region of the tilemap using a sprite renderer
// celX, celY: tile coordinates of top-left corner to draw
// sx, sy: screen position to draw at
// celW, celH: number of tiles to draw (width and height)
// spriteRenderer: function that draws a sprite at (x, y) with given tile index
func (tm *TileMap) Draw(celX, celY, sx, sy, celW, celH int, spriteRenderer func(x, y, tileIndex int)) {
	for ty := 0; ty < celH; ty++ {
		for tx := 0; tx < celW; tx++ {
			mapX := celX + tx
			mapY := celY + ty
			tileIndex := tm.Get(mapX, mapY)
			if tileIndex != 0 { // 0 = empty/transparent
				screenX := sx + tx*8 // Assuming 8x8 tiles
				screenY := sy + ty*8
				spriteRenderer(screenX, screenY, tileIndex)
			}
		}
	}
}
