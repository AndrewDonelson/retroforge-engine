package imgtool

import (
	"bytes"
	"encoding/json"
	"image"
	"image/png"
	"os"
)

// LoadPNG loads a PNG image from byte data
func LoadPNG(data []byte) (image.Image, error) {
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, NewInvalidImageError(err)
	}
	return img, nil
}

// LoadPNGFile loads a PNG image from file path
func LoadPNGFile(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, NewInvalidImageError(err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, NewInvalidImageError(err)
	}
	return img, nil
}

// LoadPalette loads a palette from JSON byte data
func LoadPalette(data []byte) (*Palette, error) {
	var palette Palette
	if err := json.Unmarshal(data, &palette); err != nil {
		return nil, err
	}
	if err := palette.Validate(); err != nil {
		return nil, err
	}
	return &palette, nil
}

// LoadPaletteFile loads a palette from JSON file
func LoadPaletteFile(path string) (*Palette, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return LoadPalette(data)
}

// SaveSprite saves a sprite to a JSON file
func SaveSprite(sprite *Sprite, path string) error {
	data, err := json.MarshalIndent(sprite, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// SavePalette saves a palette to a JSON file
func SavePalette(palette *Palette, path string) error {
	data, err := json.MarshalIndent(palette, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// ImageToRGBA converts any image to RGBA
func ImageToRGBA(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}
	return rgba
}

