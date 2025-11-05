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
func DefaultTileOptions() TileOptions {
	return TileOptions{
		Height:       16,
		LightingMode: LightingGradient,
		TileWidth:    64,
		TileHeight:   32,
	}
}

