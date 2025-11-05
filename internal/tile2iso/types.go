package tile2iso

// LightingMode indicates the lighting mode for isometric tiles
type LightingMode string

const (
	LightingNormal   LightingMode = "normal"   // No lighting adjustments
	LightingBasic    LightingMode = "basic"    // Uniform 20% adjustment
	LightingFull     LightingMode = "full"     // Enhanced top/bottom regions
	LightingGradient LightingMode = "gradient" // Smooth vertical gradient
)

// TileOptions configures the isometric tile generation
type TileOptions struct {
	Height       int          // Height of the side faces in pixels
	LightingMode LightingMode // Which lighting mode to apply
	TileWidth    int          // Width of isometric tile (default: 64)
	TileHeight   int          // Height of isometric tile (default: 32)
	ShowOutline  bool         // If true, draw outlines around the tile faces (top diamond and side parallelograms)
}

// IsometricConverter converts 2D sprites to isometric tiles
type IsometricConverter struct {
	defaultWidth  int
	defaultHeight int
}

// NewIsometricConverter creates a new isometric converter with default dimensions
func NewIsometricConverter(defaultWidth, defaultHeight int) *IsometricConverter {
	return &IsometricConverter{
		defaultWidth:  defaultWidth,
		defaultHeight: defaultHeight,
	}
}

// DefaultTileOptions returns default options for tile generation
// For 32×24 output tiles (32 wide, 16 top + 8 sides)
func DefaultTileOptions() TileOptions {
	return TileOptions{
		Height:       8,  // Side face height
		LightingMode: LightingGradient,
		TileWidth:    32, // Final tile width
		TileHeight:   16, // Top diamond height
		ShowOutline:  false,
	}
}

