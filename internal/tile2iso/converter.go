package tile2iso

import (
	"fmt"
	"image"
	"image/color"
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
	fmt.Println("=== ISOMETRIC TILE GENERATION DEBUG ===")
	
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
		options.Height = 8  // Default side height for 32×24 tiles
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
	
	// After loading sprites
	fmt.Printf("Top sprite (palette): %dx%d\n", topSprite.Width, topSprite.Height)
	fmt.Printf("Left sprite (palette): %dx%d\n", leftSprite.Width, leftSprite.Height)
	fmt.Printf("Right sprite (palette): %dx%d\n", rightSprite.Width, rightSprite.Height)

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
	
	// After converting to images
	fmt.Printf("Top image (RGBA): %dx%d\n", topImg.Bounds().Dx(), topImg.Bounds().Dy())
	fmt.Printf("Left image (RGBA): %dx%d\n", leftImg.Bounds().Dx(), leftImg.Bounds().Dy())
	fmt.Printf("Right image (RGBA): %dx%d\n", rightImg.Bounds().Dx(), rightImg.Bounds().Dy())

	// Transform top face to isometric
	topIso, err := transformToIsometric(topImg, options.TileWidth, options.TileHeight)
	if err != nil {
		return nil, fmt.Errorf("transform top to isometric: %w", err)
	}
	
	// After isometric transform
	fmt.Printf("Top isometric (RGBA): %dx%d\n", topIso.Bounds().Dx(), topIso.Bounds().Dy())

	// Scale side faces to match dimensions
	// Side faces must be: width = tileWidth/2, height = options.Height
	// This ensures they align with the diamond edges (each edge is tileWidth/2 long)
	targetSideWidth := options.TileWidth / 2
	leftScaled := scaleSideFace(leftImg, targetSideWidth, options.Height)
	rightScaled := scaleSideFace(rightImg, targetSideWidth, options.Height)
	
	// After scaling sides
	fmt.Printf("Left scaled (RGBA): %dx%d\n", leftScaled.Bounds().Dx(), leftScaled.Bounds().Dy())
	fmt.Printf("Right scaled (RGBA): %dx%d\n", rightScaled.Bounds().Dx(), rightScaled.Bounds().Dy())

	// Apply lighting to side faces BEFORE transformation
	// Lighting uses vertical Y-coordinates which only make sense in rectangular space
	// Must be applied before skewing to parallelogram
	leftLit, err := applyLighting(leftScaled, options.LightingMode, "left", options.Height)
	if err != nil {
		return nil, fmt.Errorf("apply lighting to left: %w", err)
	}

	rightLit, err := applyLighting(rightScaled, options.LightingMode, "right", options.Height)
	if err != nil {
		return nil, fmt.Errorf("apply lighting to right: %w", err)
	}
	
	// After lighting
	fmt.Printf("Left lit (RGBA): %dx%d\n", leftLit.Bounds().Dx(), leftLit.Bounds().Dy())
	fmt.Printf("Right lit (RGBA): %dx%d\n", rightLit.Bounds().Dx(), rightLit.Bounds().Dy())

	// Transform side faces to parallelograms (skewed to align with diamond edges)
	// For isometric 2:1 ratio, shear factor is 0.5
	// Transform the already-lit rectangular textures into parallelograms
	leftSkewed := transformSideFaceToParallelogram(leftLit, options.TileWidth, options.Height, true)
	rightSkewed := transformSideFaceToParallelogram(rightLit, options.TileWidth, options.Height, false)
	
	// After parallelogram transform
	if leftSkewed != nil {
		fmt.Printf("Left parallelogram (RGBA): %dx%d\n", leftSkewed.Bounds().Dx(), leftSkewed.Bounds().Dy())
	} else {
		fmt.Printf("Left parallelogram: nil!\n")
	}
	if rightSkewed != nil {
		fmt.Printf("Right parallelogram (RGBA): %dx%d\n", rightSkewed.Bounds().Dx(), rightSkewed.Bounds().Dy())
	} else {
		fmt.Printf("Right parallelogram: nil!\n")
	}

	// Composite the final tile at target size (32×24 for default options)
	// No resizing needed - we create at the correct size directly
	composite := compositeIsometricTile(topIso, leftSkewed, rightSkewed, options)

	// Final dimensions are the same as the composite canvas
	finalWidth := options.TileWidth   // 32 for default
	finalHeight := options.TileHeight + options.Height  // 16 + 8 = 24 for default

	// Convert back to pixel data (no resize - use composite directly)
	resultPixels, _, err := ImageToPixelData(composite, paletteColors)
	if err != nil {
		return nil, fmt.Errorf("convert composite to pixels: %w", err)
	}

	// Create output sprite with final dimensions
	resultSprite := &cartio.SpriteData{
		Width:        finalWidth,
		Height:       finalHeight,
		Type:         cartio.SpriteTypeStatic,
		Pixels:       resultPixels,
		UseCollision: false,
		IsUI:         false, // Isometric tiles are typically gameplay sprites
		Lifetime:     0,
		MaxSpawn:     0,
		MountPoints:  []cartio.MountPoint{},
	}
	
	// Final output
	fmt.Printf("Converting back to palette indices...\n")
	fmt.Printf("Final sprite (palette): %dx%d\n", resultSprite.Width, resultSprite.Height)
	fmt.Println("=== END DEBUG ===")

	return resultSprite, nil
}

// scaleSideFace scales a side face to the exact target width and height
// This does NOT maintain aspect ratio - it scales to the exact dimensions needed
// for isometric tiles (width = tileWidth/2, height = block height)
func scaleSideFace(img image.Image, targetWidth, targetHeight int) image.Image {
	bounds := img.Bounds()
	srcWidth := bounds.Dx()
	srcHeight := bounds.Dy()

	if srcWidth == targetWidth && srcHeight == targetHeight {
		return img
	}

	// Create scaled image with exact target dimensions
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

// resizeImage resizes an image to the target dimensions using nearest-neighbor scaling
func resizeImage(img image.Image, targetWidth, targetHeight int) image.Image {
	bounds := img.Bounds()
	srcWidth := bounds.Dx()
	srcHeight := bounds.Dy()

	if srcWidth == targetWidth && srcHeight == targetHeight {
		return img
	}

	// Create resized image
	resized := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))

	// Nearest-neighbor scaling
	for y := 0; y < targetHeight; y++ {
		srcY := (y * srcHeight) / targetHeight
		for x := 0; x < targetWidth; x++ {
			srcX := (x * srcWidth) / targetWidth
			resized.Set(x, y, img.At(srcX, srcY))
		}
	}

	return resized
}

// transformSideFaceToParallelogram transforms a rectangular side face into a parallelogram
// that aligns with the isometric diamond edges.
// For isometric 2:1 ratio, the shear factor is 0.5.
// leftSide: true for left side (positive shear), false for right side (negative shear)
// dstWidth: width of the final isometric tile (for positioning)
// dstHeight: height of the side face
func transformSideFaceToParallelogram(src image.Image, dstWidth, dstHeight int, leftSide bool) image.Image {
	if src == nil {
		return nil
	}

	srcBounds := src.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()

	// Parallelogram width is half the tile width (each side gets half)
	parallelogramWidth := dstWidth / 2

	// Create destination parallelogram image
	// The parallelogram spans the full dstWidth but only draws in its half
	dst := image.NewRGBA(image.Rect(0, 0, dstWidth, dstHeight))

	// Shear factor for isometric 2:1 ratio
	shearFactor := 0.5

	// For each pixel in destination, find corresponding source pixel
	for y := 0; y < dstHeight; y++ {
		for x := 0; x < dstWidth; x++ {
			var srcX, srcY float64
			var withinParallelogram bool

			if leftSide {
				// Left side: parallelogram spans from x=0 to x=dstWidth/2
				// Skip if outside parallelogram bounds
				if x < 0 || x >= parallelogramWidth {
					continue
				}
				withinParallelogram = true

				// Map destination x to source x (proportional within parallelogram)
				relativeX := float64(x) / float64(parallelogramWidth)
				srcX = relativeX * float64(srcWidth)

				// Apply inverse shear transformation
				// The parallelogram's top edge follows: y_top = x * shearFactor
				// The parallelogram's bottom edge follows: y_bottom = y_top + dstHeight
				// For a point at (x, y) in destination, find the corresponding source point
				yTop := float64(x) * shearFactor  // Top edge of parallelogram at this x
				
				// Check if point is within parallelogram bounds
				if float64(y) < yTop || float64(y) >= yTop+float64(dstHeight) {
					continue
				}
				
				// Calculate position relative to top edge of parallelogram
				yRelative := float64(y) - yTop
				
				// Map this relative position to source rectangle (0 to srcHeight)
				srcY = (yRelative / float64(dstHeight)) * float64(srcHeight)
			} else {
				// Right side: parallelogram spans from x=dstWidth/2 to x=dstWidth
				// Skip if outside parallelogram bounds
				if x < parallelogramWidth || x >= dstWidth {
					continue
				}
				withinParallelogram = true

				// Map destination x to source x (proportional within parallelogram)
				xInParallelogram := x - parallelogramWidth
				relativeX := float64(xInParallelogram) / float64(parallelogramWidth)
				srcX = relativeX * float64(srcWidth)

				// Apply inverse shear transformation (negative)
				// The parallelogram's top edge follows: y_top = (parallelogramWidth - xInParallelogram) * shearFactor
				// This means at x=width/2 (xInParallelogram=0), y_top=height/2
				// At x=width (xInParallelogram=width/2), y_top=0
				// So as we move right, the top edge shifts up
				yTop := float64(parallelogramWidth-xInParallelogram) * shearFactor
				
				// Check if point is within parallelogram bounds
				if float64(y) < yTop || float64(y) >= yTop+float64(dstHeight) {
					continue
				}
				
				// Calculate position relative to top edge of parallelogram
				yRelative := float64(y) - yTop
				
				// Map this relative position to source rectangle (0 to srcHeight)
				srcY = (yRelative / float64(dstHeight)) * float64(srcHeight)
			}

			if !withinParallelogram {
				continue
			}

			// Clamp source coordinates to bounds
			if srcX < 0 || srcX >= float64(srcWidth) || srcY < 0 || srcY >= float64(srcHeight) {
				continue
			}

			// Sample from source using bilinear interpolation
			c := bilinearSample(src, srcX, srcY, srcBounds)
			dst.Set(x, y, c)
		}
	}

	return dst
}

// compositeIsometricTile composites the three faces into a single image
// Layer order (back to front):
// 1. Left side face - drawn BELOW the diamond at Y=tileHeight
// 2. Right side face - drawn BELOW the diamond at Y=tileHeight
// 3. Top face - drawn at Y=0, overlapping the top edge of the sides
func compositeIsometricTile(topImg, leftImg, rightImg image.Image, options TileOptions) image.Image {
	// Final canvas size
	width := options.TileWidth
	height := options.TileHeight + options.Height

	// Create composite image
	composite := image.NewRGBA(image.Rect(0, 0, width, height))

	// Draw left side (back layer) - positioned BELOW the diamond
	// The leftImg is a full-width image (width×height) with parallelogram at x=0 to x=width/2
	// Extract the parallelogram region (x=0 to x=width/2) and draw it below the diamond
	leftBounds := leftImg.Bounds()
	leftStartY := options.TileHeight  // Start BELOW the diamond at Y=tileHeight
	// Source: parallelogram region in leftImg (x=0 to x=width/2)
	leftSrcPoint := image.Point{X: 0, Y: 0}
	// Destination: draw below diamond (x=0 to x=width/2, y=leftStartY to y=leftStartY+height)
	leftDstRect := image.Rect(0, leftStartY, width/2, leftStartY+leftBounds.Dy())
	leftX, leftY := 0, leftStartY
	draw.Draw(composite, leftDstRect, leftImg, leftSrcPoint, draw.Src)

	// Draw right side (middle layer) - positioned BELOW the diamond
	// The rightImg is a full-width image (width×height) with parallelogram at x=width/2 to x=width
	// Extract the parallelogram region (x=width/2 to x=width) and draw it below the diamond
	rightBounds := rightImg.Bounds()
	rightStartY := options.TileHeight  // Start BELOW the diamond at Y=tileHeight
	// Source: parallelogram region in rightImg (x=width/2 to x=width)
	rightSrcPoint := image.Point{X: width / 2, Y: 0}
	// Destination: draw below diamond (x=width/2 to x=width, y=rightStartY to y=rightStartY+height)
	rightDstRect := image.Rect(width/2, rightStartY, width, rightStartY+rightBounds.Dy())
	rightX, rightY := width/2, rightStartY
	draw.Draw(composite, rightDstRect, rightImg, rightSrcPoint, draw.Src)
	
	// Debug: Check if parallelograms have any non-transparent pixels
	leftNonTransparent := 0
	rightNonTransparent := 0
	for y := 0; y < leftBounds.Dy(); y++ {
		for x := 0; x < width/2; x++ {
			if c := leftImg.At(x, y); c != nil {
				_, _, _, a := c.RGBA()
				if a > 0 {
					leftNonTransparent++
				}
			}
		}
	}
	for y := 0; y < rightBounds.Dy(); y++ {
		for x := width/2; x < width; x++ {
			if c := rightImg.At(x, y); c != nil {
				_, _, _, a := c.RGBA()
				if a > 0 {
					rightNonTransparent++
				}
			}
		}
	}
	fmt.Printf("Left parallelogram non-transparent pixels: %d (out of %d)\n", leftNonTransparent, (width/2)*leftBounds.Dy())
	fmt.Printf("Right parallelogram non-transparent pixels: %d (out of %d)\n", rightNonTransparent, (width/2)*rightBounds.Dy())
	
	// Canvas and positioning
	fmt.Printf("Canvas size (RGBA): %dx%d\n", composite.Bounds().Dx(), composite.Bounds().Dy())
	fmt.Printf("Left drawn at: x=%d, y=%d (rect: %v)\n", leftX, leftY, leftDstRect)
	fmt.Printf("Right drawn at: x=%d, y=%d (rect: %v)\n", rightX, rightY, rightDstRect)
	fmt.Printf("Top drawn at: x=0, y=0\n")

	// Draw top face (front layer, with isometric projection)
	// Draw at Y=0, overlapping the top edge of the sides
	topBounds := topImg.Bounds()
	topDraw := image.Rect(0, 0, topBounds.Dx(), topBounds.Dy())
	draw.Draw(composite, topDraw, topImg, topBounds.Min, draw.Over)

	// Draw outlines if requested
	if options.ShowOutline {
		drawOutline(composite, options)
	}

	return composite
}

// drawOutline draws dark outlines around the isometric tile faces
// Top: diamond outline (4 edges)
// Left: parallelogram outline (4 edges)
// Right: parallelogram outline (4 edges)
func drawOutline(img *image.RGBA, options TileOptions) {
	darkColor := color.RGBA{R: 0, G: 0, B: 0, A: 255} // Black outline
	width := options.TileWidth
	tileHeight := options.TileHeight
	sideHeight := options.Height

	// Draw top diamond outline (isometric diamond)
	// Top diamond: from (width/2, 0) to (width, height/2) to (width/2, height) to (0, height/2)
	centerX := width / 2
	centerY := tileHeight / 2
	
	// Top point
	drawLine(img, centerX, 0, width, centerY, darkColor)
	// Right point
	drawLine(img, width, centerY, centerX, tileHeight, darkColor)
	// Bottom point
	drawLine(img, centerX, tileHeight, 0, centerY, darkColor)
	// Left point
	drawLine(img, 0, centerY, centerX, 0, darkColor)

	// Draw left side parallelogram outline
	// Left parallelogram spans x=0 to x=width/2
	// Top-left corner: (0, tileHeight)
	// Top-right corner: (width/2, tileHeight)
	// Bottom-right corner: (width/2, tileHeight + sideHeight)
	// Bottom-left corner: (0, tileHeight + sideHeight + shear)
	// Shear offset = (width/2) * shearFactor = width/4
	shearOffset := width / 4
	leftTopY := tileHeight
	leftBottomY := tileHeight + sideHeight
	
	drawLine(img, 0, leftTopY, width/2, leftTopY, darkColor)                    // Top edge
	drawLine(img, width/2, leftTopY, width/2, leftBottomY, darkColor)          // Right edge
	drawLine(img, width/2, leftBottomY, 0, leftBottomY+shearOffset, darkColor) // Bottom edge (skewed)
	drawLine(img, 0, leftBottomY+shearOffset, 0, leftTopY, darkColor)          // Left edge (skewed)

	// Draw right side parallelogram outline
	// Right parallelogram spans x=width/2 to x=width
	// Top-left corner: (width/2, tileHeight)
	// Top-right corner: (width, tileHeight)
	// Bottom-right corner: (width, tileHeight + sideHeight)
	// Bottom-left corner: (width/2, tileHeight + sideHeight - shear)
	rightTopY := tileHeight
	rightBottomY := tileHeight + sideHeight
	
	drawLine(img, width/2, rightTopY, width, rightTopY, darkColor)                    // Top edge
	drawLine(img, width, rightTopY, width, rightBottomY, darkColor)                 // Right edge
	drawLine(img, width, rightBottomY, width/2, rightBottomY-shearOffset, darkColor) // Bottom edge (skewed)
	drawLine(img, width/2, rightBottomY-shearOffset, width/2, rightTopY, darkColor) // Left edge (skewed)
}

// drawLine draws a line between two points using Bresenham's algorithm
func drawLine(img *image.RGBA, x0, y0, x1, y1 int, c color.Color) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx := 1
	if x0 > x1 {
		sx = -1
	}
	sy := 1
	if y0 > y1 {
		sy = -1
	}
	err := dx - dy

	x, y := x0, y0
	for {
		// Draw pixel if within bounds
		if x >= 0 && x < img.Bounds().Dx() && y >= 0 && y < img.Bounds().Dy() {
			img.Set(x, y, c)
		}

		if x == x1 && y == y1 {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
	}
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

