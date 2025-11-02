package imgtool

import "image"

// ToSprite performs complete PNG-to-sprite conversion pipeline
func ToSprite(img image.Image, palette *Palette, opts ToSpriteOptions) (*Sprite, error) {
	// Step 1: Scale the image
	// Create a temporary scaled image
	scaledImg := ScaleImage(img, img.Bounds().Dx(), img.Bounds().Dy(), opts.TargetWidth, opts.TargetHeight, "nearest")

	// Step 2: Map to palette
	mapOpts := MapPaletteOptions{
		DitherAlgorithm:  opts.DitherAlgorithm,
		AlphaThreshold:   opts.AlphaThreshold,
		TransparentIndex: -1,
	}

	pixels, err := MapPalette(scaledImg, palette, mapOpts)
	if err != nil {
		return nil, err
	}

	// Step 3: Create sprite object
	sprite := &Sprite{
		Width:        opts.TargetWidth,
		Height:       opts.TargetHeight,
		Pixels:       pixels,
		UseCollision: opts.UseCollision,
		IsUI:         opts.IsUI,
		Lifetime:     opts.Lifetime,
		MaxSpawn:     opts.MaxSpawn,
		MountPoints:  []Point{},
	}

	// Validate sprite
	if err := sprite.Validate(); err != nil {
		return nil, err
	}

	return sprite, nil
}

