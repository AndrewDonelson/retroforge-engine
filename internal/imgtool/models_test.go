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
			name: "valid sprite",
			sprite: &Sprite{
				Width:  16,
				Height: 16,
				Pixels: make([][]int, 16),
			},
			wantErr: false,
		},
		{
			name: "invalid - zero width",
			sprite: &Sprite{
				Width:  0,
				Height: 16,
				Pixels: make([][]int, 16),
			},
			wantErr: true,
		},
		{
			name: "invalid - wrong pixel height",
			sprite: &Sprite{
				Width:  16,
				Height: 16,
				Pixels: make([][]int, 15),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Fill pixels with valid data
			for i := range tt.sprite.Pixels {
				tt.sprite.Pixels[i] = make([]int, tt.sprite.Width)
				for j := range tt.sprite.Pixels[i] {
					tt.sprite.Pixels[i][j] = 0 // Valid palette index
				}
			}

			err := tt.sprite.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Sprite.Validate() error = %v, wantErr %v", err, tt.wantErr)
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

