package imgtool

import "testing"

func TestColorDistance(t *testing.T) {
	tests := []struct {
		name string
		c1   Color
		c2   Color
		want float64 // Expected distance (approximately)
	}{
		{
			name: "same color",
			c1:   Color{R: 255, G: 0, B: 0},
			c2:   Color{R: 255, G: 0, B: 0},
			want: 0.0,
		},
		{
			name: "different colors",
			c1:   Color{R: 255, G: 0, B: 0},
			c2:   Color{R: 0, G: 255, B: 0},
			want: 360.624, // Approximate distance
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ColorDistance(tt.c1, tt.c2)
			if tt.name == "same color" && got != 0.0 {
				t.Errorf("ColorDistance() for same color = %v, want 0.0", got)
			}
			if tt.name == "different colors" && got < 100 {
				t.Errorf("ColorDistance() for different colors = %v, want > 100", got)
			}
		})
	}
}

func TestFindClosestPaletteColor(t *testing.T) {
	palette := &Palette{
		Colors: []string{
			"#000000", // 0 - black
			"#ffffff", // 1 - white
			"#ff0000", // 2 - red
			"#00ff00", // 3 - green
			"#0000ff", // 4 - blue
		},
	}
	// Fill to 48 (game palette size)
	for len(palette.Colors) < 48 {
		palette.Colors = append(palette.Colors, "#808080")
	}

	tests := []struct {
		name  string
		color Color
		want  int
	}{
		{"red should match red", Color{R: 255, G: 0, B: 0}, 2},
		{"green should match green", Color{R: 0, G: 255, B: 0}, 3},
		{"blue should match blue", Color{R: 0, G: 0, B: 255}, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewColorCache()
			got, err := FindClosestPaletteColor(tt.color, palette, cache)
			if err != nil {
				t.Errorf("FindClosestPaletteColor() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("FindClosestPaletteColor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColorCache(t *testing.T) {
	cache := NewColorCache()

	// Test Set and Get
	cache.Set(255, 0, 0, 2)
	idx, ok := cache.Get(255, 0, 0)
	if !ok || idx != 2 {
		t.Errorf("Cache.Get() = %v, %v, want 2, true", idx, ok)
	}

	// Test miss
	_, ok = cache.Get(0, 255, 0)
	if ok {
		t.Errorf("Cache.Get() for uncached color should return false")
	}

	// Test stats
	hits, misses := cache.Stats()
	if hits == 0 || misses == 0 {
		t.Logf("Cache stats: hits=%d, misses=%d", hits, misses)
	}
}

