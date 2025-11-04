package cartio

import "testing"

func TestValidateSpriteSize(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		height   int
		isUI     bool
		wantErr  bool
	}{
		// Valid gameplay sprites (2x2 to 32x32)
		{"gameplay_2x2", 2, 2, false, false},
		{"gameplay_8x8", 8, 8, false, false},
		{"gameplay_16x16", 16, 16, false, false},
		{"gameplay_32x32", 32, 32, false, false},
		{"gameplay_3x5", 3, 5, false, false},
		{"gameplay_7x11", 7, 11, false, false},
		{"gameplay_31x2", 31, 2, false, false},
		{"gameplay_2x31", 2, 31, false, false},
		
		// Valid UI sprites (2x2 to 256x256, divisible by 2)
		{"ui_2x2", 2, 2, true, false},
		{"ui_4x4", 4, 4, true, false},
		{"ui_8x8", 8, 8, true, false},
		{"ui_16x16", 16, 16, true, false},
		{"ui_32x32", 32, 32, true, false},
		{"ui_64x64", 64, 64, true, false},
		{"ui_128x128", 128, 128, true, false},
		{"ui_256x256", 256, 256, true, false},
		{"ui_2x256", 2, 256, true, false},
		{"ui_256x2", 256, 2, true, false},
		{"ui_4x128", 4, 128, true, false},
		{"ui_100x200", 100, 200, true, false},
		
		// Invalid - too small
		{"too_small_width", 1, 10, false, true},
		{"too_small_height", 10, 1, false, true},
		{"too_small_both", 1, 1, false, true},
		{"zero_width", 0, 10, false, true},
		{"zero_height", 10, 0, false, true},
		{"zero_both", 0, 0, false, true},
		{"negative_width", -1, 10, false, true},
		{"negative_height", 10, -1, false, true},
		
		// Invalid - gameplay too large
		{"gameplay_too_large_width", 33, 32, false, true},
		{"gameplay_too_large_height", 32, 33, false, true},
		{"gameplay_too_large_both", 33, 33, false, true},
		{"gameplay_100x32", 100, 32, false, true},
		{"gameplay_32x100", 32, 100, false, true},
		
		// Invalid - UI too large
		{"ui_too_large_width", 257, 256, true, true},
		{"ui_too_large_height", 256, 257, true, true},
		{"ui_too_large_both", 257, 257, true, true},
		
		// Invalid - UI odd dimensions
		{"ui_odd_width", 3, 4, true, true},
		{"ui_odd_height", 4, 3, true, true},
		{"ui_odd_both", 3, 3, true, true},
		{"ui_odd_width_5", 5, 8, true, true},
		{"ui_odd_height_7", 8, 7, true, true},
		{"ui_255x256", 255, 256, true, true}, // 255 is odd
		{"ui_256x255", 256, 255, true, true}, // 255 is odd
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSpriteSize(tt.width, tt.height, tt.isUI)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSpriteSize(%d, %d, %v) error = %v, wantErr %v", 
					tt.width, tt.height, tt.isUI, err, tt.wantErr)
			}
		})
	}
}
