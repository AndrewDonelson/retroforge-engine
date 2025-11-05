package tile2iso

import (
	"fmt"
	"image"
	"image/draw"

	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
)

// CreateIsometricTile creates an isometric tile from three sprites
// topSpriteName, leftSpriteName, rightSpriteName: names of sprites in the sprite map
// topFrameName, leftFrameName, rightFrameName: frame names (empty for static sprites or to use first frame)
// paletteColors: 50-color palette (hex strings)
// spriteMap: map of all available sprites
// options: tile generation options
func (ic *IsometricConverter) CreateIsometricTile(
	topSpriteName, leftSpriteName, rightSpriteName string,
	topFrameName, leftFrameName, rightFrameName string,
	paletteColors []string,
	spriteMap cartio.SpriteMap,
	options TileOptions,
) (*cartio.SpriteData, error) {
	// Validate inputs
	if len(paletteColors) < 50 {
		return nil, fmt.Errorf("palette must have at least 50 colors, got %d", len(paletteColors))
	}

	// Validate dimensions BEFORE setting defaults
	// This allows explicit zero values to be caught as errors
	if options.TileWidth < 0 || options.TileHeight < 0 || options.Height < 0 {
		return nil, ErrInvalidDimensions
	}

	// Set default options (only if not explicitly set)
	if options.TileWidth == 0 {
		options.TileWidth = ic.defaultWidth
	}
	if options.TileHeight == 0 {
		options.TileHeight = ic.defaultHeight
	}
	if options.Height == 0 {
		options.Height = 16
	}
	if options.LightingMode == "" {
		options.LightingMode = LightingGradient
	}

	// Final validation after defaults are set
	if options.TileWidth <= 0 || options.TileHeight <= 0 || options.Height <= 0 {
		return nil, ErrInvalidDimensions
	}

	// Load sprites
	topSprite, err := GetSpriteFromMap(spriteMap, topSpriteName)
	if err != nil {
		return nil, fmt.Errorf("top sprite: %w", err)
	}

	leftSprite, err := GetSpriteFromMap(spriteMap, leftSpriteName)
	if err != nil {
		return nil, fmt.Errorf("left sprite: %w", err)
	}

	rightSprite, err := GetSpriteFromMap(spriteMap, rightSpriteName)
	if err != nil {
		return nil, fmt.Errorf("right sprite: %w", err)
	}

	// Extract pixel data
	topPixels, err := GetSpritePixels(topSprite, topFrameName)
	if err != nil {
		return nil, fmt.Errorf("top sprite pixels: %w", err)
	}

	leftPixels, err := GetSpritePixels(leftSprite, leftFrameName)
	if err != nil {
		return nil, fmt.Errorf("left sprite pixels: %w", err)
	}

	rightPixels, err := GetSpritePixels(rightSprite, rightFrameName)
	if err != nil {
		return nil, fmt.Errorf("right sprite pixels: %w", err)
	}

	// Convert pixel data to images
	topImg, err := PixelDataToImage(topPixels, paletteColors)
	if err != nil {
		return nil, fmt.Errorf("convert top to image: %w", err)
	}

	leftImg, err := PixelDataToImage(leftPixels, paletteColors)
	if err != nil {
		return nil, fmt.Errorf("convert left to image: %w", err)
	}

	rightImg, err := PixelDataToImage(rightPixels, paletteColors)
	if err != nil {
		return nil, fmt.Errorf("convert right to image: %w", err)
	}

	// Transform top face to isometric
	topIso, err := transformToIsometric(topImg, options.TileWidth, options.TileHeight)
	if err != nil {
		return nil, fmt.Errorf("transform top to isometric: %w", err)
	}

	// Scale side faces to match height
	leftScaled := scaleSideFace(leftImg, options.Height)
	rightScaled := scaleSideFace(rightImg, options.Height)

	// Apply lighting to side faces
	leftLit, err := applyLighting(leftScaled, options.LightingMode, "left", options.Height)
	if err != nil {
		return nil, fmt.Errorf("apply lighting to left: %w", err)
	}

	rightLit, err := applyLighting(rightScaled, options.LightingMode, "right", options.Height)
	if err != nil {
		return nil, fmt.Errorf("apply lighting to right: %w", err)
	}

	// Composite the final tile
	composite := compositeIsometricTile(topIso, leftLit, rightLit, options)

	// Convert back to pixel data
	resultPixels, _, err := ImageToPixelData(composite, paletteColors)
	if err != nil {
		return nil, fmt.Errorf("convert composite to pixels: %w", err)
	}

	// Create output sprite
	resultSprite := &cartio.SpriteData{
		Width:        options.TileWidth,
		Height:       options.TileHeight + options.Height,
		Type:         cartio.SpriteTypeStatic,
		Pixels:       resultPixels,
		UseCollision: false,
		IsUI:         false, // Isometric tiles are typically gameplay sprites
		Lifetime:     0,
		MaxSpawn:     0,
		MountPoints:  []cartio.MountPoint{},
	}

	return resultSprite, nil
}

// scaleSideFace scales a side face to the target height (maintaining aspect ratio)
func scaleSideFace(img image.Image, targetHeight int) image.Image {
	bounds := img.Bounds()
	srcWidth := bounds.Dx()
	srcHeight := bounds.Dy()

	if srcHeight == targetHeight {
		return img
	}

	// Calculate target width (maintain aspect ratio)
	targetWidth := (srcWidth * targetHeight) / srcHeight

	// Create scaled image
	scaled := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))

	// Simple nearest-neighbor scaling
	for y := 0; y < targetHeight; y++ {
		srcY := (y * srcHeight) / targetHeight
		for x := 0; x < targetWidth; x++ {
			srcX := (x * srcWidth) / targetWidth
			scaled.Set(x, y, img.At(srcX, srcY))
		}
	}

	return scaled
}

// compositeIsometricTile composites the three faces into a single image
// Layer order (back to front):
// 1. Left side face - drawn at position (0, tileHeight)
// 2. Right side face - drawn at position (tileWidth/2, tileHeight)
// 3. Top face - drawn at position (0, 0) with isometric projection
func compositeIsometricTile(topImg, leftImg, rightImg image.Image, options TileOptions) image.Image {
	// Final canvas size
	width := options.TileWidth
	height := options.TileHeight + options.Height

	// Create composite image
	composite := image.NewRGBA(image.Rect(0, 0, width, height))

	// Draw left side (back layer)
	leftBounds := leftImg.Bounds()
	leftDraw := image.Rect(0, options.TileHeight, leftBounds.Dx(), options.TileHeight+leftBounds.Dy())
	draw.Draw(composite, leftDraw, leftImg, leftBounds.Min, draw.Src)

	// Draw right side (middle layer)
	rightBounds := rightImg.Bounds()
	rightDraw := image.Rect(width/2, options.TileHeight, width/2+rightBounds.Dx(), options.TileHeight+rightBounds.Dy())
	draw.Draw(composite, rightDraw, rightImg, rightBounds.Min, draw.Src)

	// Draw top face (front layer, with isometric projection)
	topBounds := topImg.Bounds()
	topDraw := image.Rect(0, 0, topBounds.Dx(), topBounds.Dy())
	draw.Draw(composite, topDraw, topImg, topBounds.Min, draw.Over)

	return composite
}

