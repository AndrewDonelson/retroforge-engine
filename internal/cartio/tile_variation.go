package cartio

import (
	"math/rand"
)

// TileVariation represents a transformation to apply to a tile
type TileVariation int

const (
	VariationNone TileVariation = iota
	VariationFlipHorizontal
	VariationFlipVertical
	VariationRotateCW
	VariationRotateCCW
)

// GetTileVariation returns a deterministic variation for a tile at (gridX, gridY).
// Uses seed for deterministic randomization - same seed + coordinates = same variation.
// This ensures tile variations are consistent across game sessions.
// Only returns variations for normal (non-isometric) tiles, as isometric tiles
// cannot be rotated/flipped without breaking the 3D visual appearance.
func GetTileVariation(seed int, gridX, gridY int, isIsometric bool) TileVariation {
	// Isometric tiles cannot be rotated/flipped (would break 3D appearance)
	if isIsometric {
		return VariationNone
	}

	// Use seed + grid position to create deterministic random value
	rng := rand.New(rand.NewSource(int64(seed + gridX*1000 + gridY)))
	variation := rng.Intn(5) // 0-4: none, flipH, flipV, rotateCW, rotateCCW

	return TileVariation(variation)
}

// ApplyVariation applies a transformation to pixel data.
// Returns transformed pixels and updated dimensions (width/height may be swapped for rotations).
// The transformation is applied in-place to the pixel array structure.
func ApplyVariation(pixels [][]int, width, height int, variation TileVariation) ([][]int, int, int) {
	if variation == VariationNone {
		return pixels, width, height
	}

	switch variation {
	case VariationFlipHorizontal:
		return flipHorizontal(pixels, width, height), width, height
	case VariationFlipVertical:
		return flipVertical(pixels, width, height), width, height
	case VariationRotateCW:
		return rotateCW(pixels, width, height), height, width
	case VariationRotateCCW:
		return rotateCCW(pixels, width, height), height, width
	default:
		return pixels, width, height
	}
}

// flipHorizontal flips pixels horizontally
func flipHorizontal(pixels [][]int, width, height int) [][]int {
	result := make([][]int, height)
	for y := 0; y < height; y++ {
		result[y] = make([]int, width)
		for x := 0; x < width; x++ {
			result[y][x] = pixels[y][width-1-x]
		}
	}
	return result
}

// flipVertical flips pixels vertically
func flipVertical(pixels [][]int, width, height int) [][]int {
	result := make([][]int, height)
	for y := 0; y < height; y++ {
		result[y] = make([]int, width)
		for x := 0; x < width; x++ {
			result[y][x] = pixels[height-1-y][x]
		}
	}
	return result
}

// rotateCW rotates pixels 90 degrees clockwise (swaps width/height)
func rotateCW(pixels [][]int, width, height int) [][]int {
	result := make([][]int, height)
	for y := 0; y < height; y++ {
		result[y] = make([]int, width)
		for x := 0; x < width; x++ {
			// 90° CW: (x, y) -> (height-1-y, x)
			result[y][x] = pixels[width-1-x][y]
		}
	}
	return result
}

// rotateCCW rotates pixels 90 degrees counter-clockwise (swaps width/height)
func rotateCCW(pixels [][]int, width, height int) [][]int {
	result := make([][]int, height)
	for y := 0; y < height; y++ {
		result[y] = make([]int, width)
		for x := 0; x < width; x++ {
			// 90° CCW: (x, y) -> (y, width-1-x)
			result[y][x] = pixels[x][height-1-y]
		}
	}
	return result
}

