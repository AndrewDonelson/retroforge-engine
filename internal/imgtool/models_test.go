package imgtool

import "testing"

func TestPalette_Validate(t *testing.T) {
	tests := []struct {
		name    string
		palette *Palette
		wantErr bool
	}{
		{
			name: "valid 50 colors",
			palette: &Palette{
				Colors: make([]string, 50),
			},
			wantErr: false,
		},
		{
			name: "invalid - 49 colors",
			palette: &Palette{
				Colors: make([]string, 49),
			},
			wantErr: true,
		},
		{
			name: "invalid - 51 colors",
			palette: &Palette{
				Colors: make([]string, 51),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Fill with valid hex colors
			for i := range tt.palette.Colors {
				tt.palette.Colors[i] = "#000000"
			}

			err := tt.palette.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Palette.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSprite_Validate(t *testing.T) {
	tests := []struct {
		name    string
		sprite  *Sprite
		wantErr bool
	}{
		{
			name: "valid gameplay sprite 16x16",
			sprite: &Sprite{
				Width:  16,
				Height: 16,
				IsUI:   false,
				Pixels: make([][]int, 16),
			},
			wantErr: false,
		},
		{
			name: "valid UI sprite 16x16",
			sprite: &Sprite{
				Width:  16,
				Height: 16,
				IsUI:   true,
				Pixels: make([][]int, 16),
			},
			wantErr: false,
		},
		{
			name: "invalid - width < 2",
			sprite: &Sprite{
				Width:  1,
				Height: 16,
				IsUI:   false,
				Pixels: make([][]int, 16),
			},
			wantErr: true,
		},
		{
			name: "invalid - height < 2",
			sprite: &Sprite{
				Width:  16,
				Height: 1,
				IsUI:   false,
				Pixels: make([][]int, 1),
			},
			wantErr: true,
		},
		{
			name: "invalid - zero width",
			sprite: &Sprite{
				Width:  0,
				Height: 16,
				IsUI:   false,
				Pixels: make([][]int, 16),
			},
			wantErr: true,
		},
		{
			name: "invalid - gameplay sprite too large",
			sprite: &Sprite{
				Width:  33,
				Height: 32,
				IsUI:   false,
				Pixels: make([][]int, 32),
			},
			wantErr: true,
		},
		{
			name: "invalid - UI sprite too large",
			sprite: &Sprite{
				Width:  257,
				Height: 256,
				IsUI:   true,
				Pixels: make([][]int, 256),
			},
			wantErr: true,
		},
		{
			name: "invalid - UI sprite odd width",
			sprite: &Sprite{
				Width:  3,
				Height: 4,
				IsUI:   true,
				Pixels: make([][]int, 4),
			},
			wantErr: true,
		},
		{
			name: "invalid - UI sprite odd height",
			sprite: &Sprite{
				Width:  4,
				Height: 3,
				IsUI:   true,
				Pixels: make([][]int, 3),
			},
			wantErr: true,
		},
		{
			name: "invalid - wrong pixel height",
			sprite: &Sprite{
				Width:  16,
				Height: 16,
				IsUI:   false,
				Pixels: make([][]int, 15),
			},
			wantErr: true,
		},
		{
			name: "valid - gameplay sprite 2x2",
			sprite: &Sprite{
				Width:  2,
				Height: 2,
				IsUI:   false,
				Pixels: make([][]int, 2),
			},
			wantErr: false,
		},
		{
			name: "valid - gameplay sprite 32x32",
			sprite: &Sprite{
				Width:  32,
				Height: 32,
				IsUI:   false,
				Pixels: make([][]int, 32),
			},
			wantErr: false,
		},
		{
			name: "valid - UI sprite 2x2",
			sprite: &Sprite{
				Width:  2,
				Height: 2,
				IsUI:   true,
				Pixels: make([][]int, 2),
			},
			wantErr: false,
		},
		{
			name: "valid - UI sprite 256x256",
			sprite: &Sprite{
				Width:  256,
				Height: 256,
				IsUI:   true,
				Pixels: make([][]int, 256),
			},
			wantErr: false,
		},
		{
			name: "valid - UI sprite 2x256",
			sprite: &Sprite{
				Width:  2,
				Height: 256,
				IsUI:   true,
				Pixels: make([][]int, 256),
			},
			wantErr: false,
		},
		{
			name: "valid - gameplay sprite 3x5 (non-square)",
			sprite: &Sprite{
				Width:  3,
				Height: 5,
				IsUI:   false,
				Pixels: make([][]int, 5),
			},
			wantErr: false,
		},
		{
			name: "invalid - wrong pixel width",
			sprite: &Sprite{
				Width:  16,
				Height: 16,
				IsUI:   false,
				Pixels: make([][]int, 16),
			},
			wantErr: true,
		},
		{
			name: "invalid - palette index too high",
			sprite: &Sprite{
				Width:  8,
				Height: 8,
				IsUI:   false,
				Pixels: make([][]int, 8),
			},
			wantErr: true,
		},
		{
			name: "invalid - palette index too low",
			sprite: &Sprite{
				Width:  8,
				Height: 8,
				IsUI:   false,
				Pixels: make([][]int, 8),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Fill pixels with valid data (or invalid data for error cases)
			for i := range tt.sprite.Pixels {
				if tt.name == "invalid - wrong pixel width" {
					// Make row width wrong
					tt.sprite.Pixels[i] = make([]int, tt.sprite.Width+1) // Wrong width
				} else {
					tt.sprite.Pixels[i] = make([]int, tt.sprite.Width)
				}
				for j := range tt.sprite.Pixels[i] {
					if tt.name == "invalid - palette index too high" {
						tt.sprite.Pixels[i][j] = 50 // Invalid: > 49
					} else if tt.name == "invalid - palette index too low" {
						tt.sprite.Pixels[i][j] = -2 // Invalid: < -1
					} else {
						tt.sprite.Pixels[i][j] = 0 // Valid palette index
					}
				}
			}

			err := tt.sprite.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Sprite.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSprite_ValidateSize(t *testing.T) {
	tests := []struct {
		name    string
		sprite  *Sprite
		wantErr bool
	}{
		{
			name: "valid gameplay 2x2",
			sprite: &Sprite{Width: 2, Height: 2, IsUI: false},
			wantErr: false,
		},
		{
			name: "valid gameplay 32x32",
			sprite: &Sprite{Width: 32, Height: 32, IsUI: false},
			wantErr: false,
		},
		{
			name: "valid gameplay 3x7",
			sprite: &Sprite{Width: 3, Height: 7, IsUI: false},
			wantErr: false,
		},
		{
			name: "valid UI 2x2",
			sprite: &Sprite{Width: 2, Height: 2, IsUI: true},
			wantErr: false,
		},
		{
			name: "valid UI 256x256",
			sprite: &Sprite{Width: 256, Height: 256, IsUI: true},
			wantErr: false,
		},
		{
			name: "valid UI 4x128",
			sprite: &Sprite{Width: 4, Height: 128, IsUI: true},
			wantErr: false,
		},
		{
			name: "invalid - width < 2",
			sprite: &Sprite{Width: 1, Height: 16, IsUI: false},
			wantErr: true,
		},
		{
			name: "invalid - height < 2",
			sprite: &Sprite{Width: 16, Height: 1, IsUI: false},
			wantErr: true,
		},
		{
			name: "invalid - gameplay > 32",
			sprite: &Sprite{Width: 33, Height: 32, IsUI: false},
			wantErr: true,
		},
		{
			name: "invalid - UI > 256",
			sprite: &Sprite{Width: 257, Height: 256, IsUI: true},
			wantErr: true,
		},
		{
			name: "invalid - UI odd width",
			sprite: &Sprite{Width: 3, Height: 4, IsUI: true},
			wantErr: true,
		},
		{
			name: "invalid - UI odd height",
			sprite: &Sprite{Width: 4, Height: 3, IsUI: true},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sprite.ValidateSize()
			if (err != nil) != tt.wantErr {
				t.Errorf("Sprite.ValidateSize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHexToColor(t *testing.T) {
	tests := []struct {
		name    string
		hex     string
		wantErr bool
	}{
		{"valid hex", "#ff0000", false},
		{"valid hex uppercase", "#FF0000", false},
		{"invalid - no hash", "ff0000", true},
		{"invalid - wrong length", "#ff00", true},
		{"invalid - wrong length", "#ff000000", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := HexToColor(tt.hex)
			if (err != nil) != tt.wantErr {
				t.Errorf("HexToColor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestColor_ToHex(t *testing.T) {
	c := Color{R: 255, G: 0, B: 0}
	got := c.ToHex()
	want := "#ff0000"
	if got != want {
		t.Errorf("Color.ToHex() = %v, want %v", got, want)
	}
}

