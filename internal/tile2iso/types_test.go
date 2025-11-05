package tile2iso

import "testing"

func TestNewIsometricConverter(t *testing.T) {
	ic := NewIsometricConverter(64, 32)
	if ic == nil {
		t.Fatal("NewIsometricConverter returned nil")
	}
	if ic.defaultWidth != 64 {
		t.Errorf("expected defaultWidth 64, got %d", ic.defaultWidth)
	}
	if ic.defaultHeight != 32 {
		t.Errorf("expected defaultHeight 32, got %d", ic.defaultHeight)
	}
}

func TestDefaultTileOptions(t *testing.T) {
	opts := DefaultTileOptions()
	if opts.Height != 16 {
		t.Errorf("expected Height 16, got %d", opts.Height)
	}
	if opts.LightingMode != LightingGradient {
		t.Errorf("expected LightingMode Gradient, got %s", opts.LightingMode)
	}
	if opts.TileWidth != 64 {
		t.Errorf("expected TileWidth 64, got %d", opts.TileWidth)
	}
	if opts.TileHeight != 32 {
		t.Errorf("expected TileHeight 32, got %d", opts.TileHeight)
	}
}

