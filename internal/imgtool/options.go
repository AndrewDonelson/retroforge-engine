package imgtool

// QuantizeOptions for color quantization
type QuantizeOptions struct {
	DitherAlgorithm   string // "floyd-steinberg", "ordered", "none" (default: "floyd-steinberg")
	AlphaThreshold    uint8  // 0-255 (default: 128)
	EnforceBlackWhite bool   // Ensure indices 0=black, 1=white (default: true)
}

// MapPaletteOptions for palette mapping
type MapPaletteOptions struct {
	DitherAlgorithm  string // "floyd-steinberg", "ordered", "none" (default: "floyd-steinberg")
	AlphaThreshold   uint8  // 0-255 (default: 128)
	TransparentIndex int    // Usually -1 (default: -1)
}

// ScaleOptions for image scaling
type ScaleOptions struct {
	Width           int    // Target width (default: 16)
	Height          int    // Target height (default: 16)
	Algorithm       string // "nearest", "bilinear", "bicubic" (default: "nearest")
	EnsureDivisible bool   // Ensure dimensions divisible by 2 (default: true)
	PreserveAspect  bool   // Maintain aspect ratio (default: false)
}

// ToSpriteOptions for sprite conversion
type ToSpriteOptions struct {
	Name           string // Sprite name (required)
	TargetWidth    int    // Target width (default: 16)
	TargetHeight   int    // Target height (default: 16)
	UseCollision   bool   // Enable physics collision (default: false)
	IsUI           bool   // UI sprite (no physics) (default: false)
	Lifetime       int    // Auto-destroy ms (default: 0)
	MaxSpawn       int    // Max simultaneous instances (default: 0)
	DitherAlgorithm string // Dithering (default: "floyd-steinberg")
	AlphaThreshold  uint8  // Alpha threshold (default: 128)
}

// DefaultQuantizeOptions returns default quantization options
func DefaultQuantizeOptions() QuantizeOptions {
	return QuantizeOptions{
		DitherAlgorithm:   "floyd-steinberg",
		AlphaThreshold:    128,
		EnforceBlackWhite: true,
	}
}

// DefaultMapPaletteOptions returns default palette mapping options
func DefaultMapPaletteOptions() MapPaletteOptions {
	return MapPaletteOptions{
		DitherAlgorithm:  "floyd-steinberg",
		AlphaThreshold:   128,
		TransparentIndex: -1,
	}
}

// DefaultScaleOptions returns default scaling options
func DefaultScaleOptions() ScaleOptions {
	return ScaleOptions{
		Width:           16,
		Height:          16,
		Algorithm:       "nearest",
		EnsureDivisible: true,
		PreserveAspect:  false,
	}
}

// DefaultToSpriteOptions returns default sprite conversion options
func DefaultToSpriteOptions(name string) ToSpriteOptions {
	return ToSpriteOptions{
		Name:            name,
		TargetWidth:     16,
		TargetHeight:    16,
		UseCollision:    false,
		IsUI:            false,
		Lifetime:        0,
		MaxSpawn:        0,
		DitherAlgorithm: "floyd-steinberg",
		AlphaThreshold:  128,
	}
}

