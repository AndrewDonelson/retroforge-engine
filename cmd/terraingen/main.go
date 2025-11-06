package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/AndrewDonelson/retroforge-engine/internal/imgtool"
	"github.com/AndrewDonelson/retroforge-engine/internal/pal"
)

// Global palette (can be nil for RGBA mode)
var Palette []color.RGBA
var usePalette bool

// ================================
// NOISE FUNCTIONS
// ================================

// PerlinNoise generates seamless (periodic) Perlin noise for tile generation.
// The noise wraps around at tile boundaries (16x16) to ensure seamless tiling.
// When tiles are placed next to each other, the patterns will match at edges.
// This uses toroidal (wrap-around) coordinates to ensure pixel (0,y) matches pixel (16,y) etc.
func PerlinNoise(x, y int, frequency float64, seed int64) float64 {
	tileSize := 16
	
	// Normalize coordinates to tile space (0-15) for seamless wrapping
	// This ensures the pattern repeats every 16 pixels
	normX := x % tileSize
	normY := y % tileSize
	if normX < 0 {
		normX += tileSize
	}
	if normY < 0 {
		normY += tileSize
	}

	fx := float64(normX) * frequency
	fy := float64(normY) * frequency

	// Calculate grid cell boundaries
	x0 := int(math.Floor(fx))
	y0 := int(math.Floor(fy))
	x1 := x0 + 1
	y1 := y0 + 1

	// Wrap grid coordinates for seamless tiling
	// Key: grid points at tile boundaries (0 and 16) must have the same gradient
	wrapX0 := x0 % tileSize
	wrapY0 := y0 % tileSize
	wrapX1 := x1 % tileSize
	wrapY1 := y1 % tileSize
	if wrapX0 < 0 {
		wrapX0 += tileSize
	}
	if wrapY0 < 0 {
		wrapY0 += tileSize
	}
	if wrapX1 < 0 {
		wrapX1 += tileSize
	}
	if wrapY1 < 0 {
		wrapY1 += tileSize
	}

	// Calculate fractional parts for interpolation
	sx := fx - float64(x0)
	sy := fy - float64(y0)

	// Get gradients at the four corners (using wrapped coordinates)
	n00 := dotGridGradientSeamless(wrapX0, wrapY0, fx, fy, tileSize, seed)
	n10 := dotGridGradientSeamless(wrapX1, wrapY0, fx, fy, tileSize, seed)
	n01 := dotGridGradientSeamless(wrapX0, wrapY1, fx, fy, tileSize, seed)
	n11 := dotGridGradientSeamless(wrapX1, wrapY1, fx, fy, tileSize, seed)

	sx = fade(sx)
	sy = fade(sy)

	ix0 := lerp(n00, n10, sx)
	ix1 := lerp(n01, n11, sx)

	return lerp(ix0, ix1, sy)
}

// dotGridGradientSeamless calculates gradient with seamless wrapping.
// For seamless tiling, gradients at wrapped coordinates must be consistent.
func dotGridGradientSeamless(ix, iy int, x, y float64, tileSize int, seed int64) float64 {
	// Normalize gradient coordinates to tile space (0-15)
	// This ensures that gradient at position (0,y) is the same as gradient at (16,y)
	normIx := ix % tileSize
	normIy := iy % tileSize
	if normIx < 0 {
		normIx += tileSize
	}
	if normIy < 0 {
		normIy += tileSize
	}

	// Generate consistent gradient direction based on normalized position
	// Same normalized position = same gradient, regardless of which tile it's in
	random := math.Sin(float64(normIx)*12.9898+float64(normIy)*78.233+float64(seed)) * 43758.5453
	random = random - math.Floor(random)
	angle := random * 2.0 * math.Pi

	gradX := math.Cos(angle)
	gradY := math.Sin(angle)

	// Calculate distance from sample point (x,y) to grid point (normIx, normIy)
	// Use normalized coordinates for distance calculation to ensure seamless wrapping
	dx := x - float64(normIx)
	dy := y - float64(normIy)
	
	// Handle wrap-around for toroidal space (shortest path)
	// This ensures smooth interpolation across tile boundaries
	if dx > float64(tileSize)/2.0 {
		dx -= float64(tileSize)
	} else if dx < -float64(tileSize)/2.0 {
		dx += float64(tileSize)
	}
	if dy > float64(tileSize)/2.0 {
		dy -= float64(tileSize)
	} else if dy < -float64(tileSize)/2.0 {
		dy += float64(tileSize)
	}

	return dx*gradX + dy*gradY
}

func fade(t float64) float64 {
	return t * t * (3.0 - 2.0*t)
}

func lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}

func dotGridGradient(ix, iy int, x, y float64, seed int64) float64 {
	random := math.Sin(float64(ix)*12.9898+float64(iy)*78.233+float64(seed)) * 43758.5453
	random = random - math.Floor(random)
	angle := random * 2.0 * math.Pi

	gradX := math.Cos(angle)
	gradY := math.Sin(angle)

	dx := x - float64(ix)
	dy := y - float64(iy)

	return dx*gradX + dy*gradY
}

// ================================
// BASE TERRAIN GENERATORS
// ================================

// GenerateGrass creates a textured grass tile with seamless patterns
func GenerateGrass() [][]int {
	pixels := make([][]int, 16)
	for i := range pixels {
		pixels[i] = make([]int, 16)
	}

	darkGreen := 27
	medGreen := 25
	lightGreen := 24
	veryDark := 28

	rng := rand.New(rand.NewSource(42))

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			// Use seamless Perlin noise - coordinates are automatically wrapped
			noise := PerlinNoise(x, y, 0.4, 42)
			noise2 := PerlinNoise(x, y, 0.8, 43) * 0.3
			// Use seamless gradient (wraps around)
			gradient := PerlinNoise(0, y, 0.1, 100) * 0.3 // Subtle vertical gradient
			value := noise + noise2 + gradient

			if value < -0.2 {
				pixels[y][x] = veryDark
			} else if value < 0.2 {
				pixels[y][x] = darkGreen
			} else if value < 0.6 {
				pixels[y][x] = medGreen
			} else {
				pixels[y][x] = lightGreen
			}

			// Use seamless random for speckles
			if rng.Float64() < 0.05 {
				pixels[y][x] = darkGreen
			}
		}
	}

	return pixels
}

// GenerateDarkGrass creates a darker forest floor grass
func GenerateDarkGrass() [][]int {
	pixels := make([][]int, 16)
	for i := range pixels {
		pixels[i] = make([]int, 16)
	}

	veryDark := 29  // Dark teal
	darkGreen := 27 // Dark olive
	medGreen := 26  // Teal green
	darkest := 3    // Dark green

	rng := rand.New(rand.NewSource(45))

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			noise := PerlinNoise(x, y, 0.35, 45)
			noise2 := PerlinNoise(x, y, 0.7, 46) * 0.4
			gradient := float64(y) / 18.0
			value := noise + noise2 + gradient

			if value < -0.1 {
				pixels[y][x] = darkest
			} else if value < 0.3 {
				pixels[y][x] = veryDark
			} else if value < 0.7 {
				pixels[y][x] = darkGreen
			} else {
				pixels[y][x] = medGreen
			}

			if rng.Float64() < 0.06 {
				pixels[y][x] = veryDark
			}
		}
	}

	return pixels
}

// GenerateDirt creates a textured dirt tile
func GenerateDirt() [][]int {
	pixels := make([][]int, 16)
	for i := range pixels {
		pixels[i] = make([]int, 16)
	}

	darkBrown := 18
	medBrown := 19
	lightBrown := 20

	rng := rand.New(rand.NewSource(50))

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			noise := PerlinNoise(x, y, 0.45, 50)
			noise2 := PerlinNoise(x, y, 0.9, 51) * 0.4
			value := noise + noise2

			if value < -0.1 {
				pixels[y][x] = darkBrown
			} else if value < 0.4 {
				pixels[y][x] = medBrown
			} else {
				pixels[y][x] = lightBrown
			}

			if rng.Float64() < 0.04 {
				pixels[y][x] = darkBrown
			}
		}
	}

	return pixels
}

// GenerateSand creates a textured sand tile
func GenerateSand() [][]int {
	pixels := make([][]int, 16)
	for i := range pixels {
		pixels[i] = make([]int, 16)
	}

	darkSand := 21
	medSand := 22
	lightSand := 23

	rng := rand.New(rand.NewSource(70))

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			noise := PerlinNoise(x, y, 0.6, 70)
			gradient := float64(y) / 32.0
			value := noise + gradient

			if value < 0.0 {
				pixels[y][x] = darkSand
			} else if value < 0.5 {
				pixels[y][x] = medSand
			} else {
				pixels[y][x] = medSand
			}

			if rng.Float64() < 0.02 {
				pixels[y][x] = lightSand
			}
		}
	}

	return pixels
}

// GenerateGravel creates a rocky gravel texture
func GenerateGravel() [][]int {
	pixels := make([][]int, 16)
	for i := range pixels {
		pixels[i] = make([]int, 16)
	}

	veryDark := 39
	darkGray := 38
	medGray := 37
	tanGray := 21

	rng := rand.New(rand.NewSource(55))

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			noise := PerlinNoise(x, y, 0.5, 55)
			noise2 := PerlinNoise(x, y, 1.0, 56) * 0.3
			value := noise + noise2

			if value < -0.2 {
				pixels[y][x] = veryDark
			} else if value < 0.2 {
				pixels[y][x] = darkGray
			} else if value < 0.6 {
				pixels[y][x] = medGray
			} else {
				pixels[y][x] = tanGray
			}

			// Random pebbles
			if rng.Float64() < 0.08 {
				if rng.Float64() < 0.5 {
					pixels[y][x] = darkGray
				} else {
					pixels[y][x] = tanGray
				}
			}
		}
	}

	return pixels
}

// GenerateStone creates a textured stone tile
func GenerateStone() [][]int {
	pixels := make([][]int, 16)
	for i := range pixels {
		pixels[i] = make([]int, 16)
	}

	veryDark := 39
	darkGray := 38
	medGray := 37
	lightGray := 36

	rng := rand.New(rand.NewSource(80))

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			noise1 := PerlinNoise(x, y, 0.35, 80)
			noise2 := PerlinNoise(x, y, 0.7, 81) * 0.4
			value := noise1 + noise2

			if value < -0.2 {
				pixels[y][x] = veryDark
			} else if value < 0.1 {
				pixels[y][x] = darkGray
			} else if value < 0.5 {
				pixels[y][x] = medGray
			} else {
				pixels[y][x] = lightGray
			}

			if rng.Float64() < 0.05 {
				pixels[y][x] = veryDark
			}
		}
	}

	return pixels
}

// GenerateShallowWater creates a light blue shallow water texture
func GenerateShallowWater() [][]int {
	pixels := make([][]int, 16)
	for i := range pixels {
		pixels[i] = make([]int, 16)
	}

	lightBlue := 34  // Light blue
	medBlue := 33    // Sky blue
	brightBlue := 32 // Bright blue

	rng := rand.New(rand.NewSource(61))

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			wave := math.Sin(float64(x)*0.4+float64(y)*0.25) * 0.25
			noise := PerlinNoise(x, y, 0.25, 61)
			value := noise + wave

			if value < -0.05 {
				pixels[y][x] = brightBlue
			} else if value < 0.5 {
				pixels[y][x] = medBlue
			} else {
				pixels[y][x] = lightBlue
			}

			if rng.Float64() < 0.04 {
				pixels[y][x] = lightBlue
			}
		}
	}

	return pixels
}

// GenerateWater creates a textured water tile
func GenerateWater() [][]int {
	pixels := make([][]int, 16)
	for i := range pixels {
		pixels[i] = make([]int, 16)
	}

	darkBlue := 31
	medBlue := 32
	lightBlue := 33

	rng := rand.New(rand.NewSource(60))

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			wave := math.Sin(float64(x)*0.5+float64(y)*0.3) * 0.3
			noise := PerlinNoise(x, y, 0.3, 60)
			value := noise + wave

			if value < -0.1 {
				pixels[y][x] = darkBlue
			} else if value < 0.4 {
				pixels[y][x] = medBlue
			} else {
				pixels[y][x] = lightBlue
			}

			if rng.Float64() < 0.03 {
				pixels[y][x] = lightBlue
			}
		}
	}

	return pixels
}

// GenerateDeepWater creates a very dark blue deep water texture
func GenerateDeepWater() [][]int {
	pixels := make([][]int, 16)
	for i := range pixels {
		pixels[i] = make([]int, 16)
	}

	veryDark := 1   // Dark blue
	darkBlue := 30  // Dark blue
	oceanBlue := 31 // Ocean blue

	rng := rand.New(rand.NewSource(62))

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			wave := math.Sin(float64(x)*0.6+float64(y)*0.35) * 0.2
			noise := PerlinNoise(x, y, 0.35, 62)
			value := noise + wave

			if value < -0.2 {
				pixels[y][x] = veryDark
			} else if value < 0.3 {
				pixels[y][x] = darkBlue
			} else {
				pixels[y][x] = oceanBlue
			}

			if rng.Float64() < 0.02 {
				pixels[y][x] = darkBlue
			}
		}
	}

	return pixels
}

// GenerateMud creates a dark wet mud texture
func GenerateMud() [][]int {
	pixels := make([][]int, 16)
	for i := range pixels {
		pixels[i] = make([]int, 16)
	}

	veryDark := 17 // Dark purple-brown
	darkBrown := 18
	medBrown := 28 // Olive brown

	rng := rand.New(rand.NewSource(65))

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			noise := PerlinNoise(x, y, 0.4, 65)
			noise2 := PerlinNoise(x, y, 0.85, 66) * 0.35
			gradient := float64(y) / 24.0
			value := noise + noise2 + gradient

			if value < 0.0 {
				pixels[y][x] = veryDark
			} else if value < 0.5 {
				pixels[y][x] = darkBrown
			} else {
				pixels[y][x] = medBrown
			}

			// Wet patches
			if rng.Float64() < 0.05 {
				pixels[y][x] = veryDark
			}
		}
	}

	return pixels
}

// GenerateSnow creates a white snow texture with blue shadows
func GenerateSnow() [][]int {
	pixels := make([][]int, 16)
	for i := range pixels {
		pixels[i] = make([]int, 16)
	}

	white := 47
	offWhite := 7
	blueGray := 36
	lightGray := 46

	rng := rand.New(rand.NewSource(85))

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			noise := PerlinNoise(x, y, 0.5, 85)
			noise2 := PerlinNoise(x, y, 1.0, 86) * 0.2
			gradient := float64(y) / 28.0
			value := noise + noise2 + gradient

			if value < 0.1 {
				pixels[y][x] = blueGray
			} else if value < 0.4 {
				pixels[y][x] = lightGray
			} else if value < 0.7 {
				pixels[y][x] = offWhite
			} else {
				pixels[y][x] = white
			}

			// Sparkle highlights
			if rng.Float64() < 0.03 {
				pixels[y][x] = white
			}
		}
	}

	return pixels
}

// GenerateIce creates a smooth glossy ice texture
func GenerateIce() [][]int {
	pixels := make([][]int, 16)
	for i := range pixels {
		pixels[i] = make([]int, 16)
	}

	blueGray := 36
	lightBlue := 34
	veryLight := 35
	skyBlue := 33

	rng := rand.New(rand.NewSource(88))

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			noise := PerlinNoise(x, y, 0.3, 88)
			// Diagonal streaks
			streak := math.Sin((float64(x)+float64(y))*0.4) * 0.25
			value := noise + streak

			if value < -0.1 {
				pixels[y][x] = blueGray
			} else if value < 0.3 {
				pixels[y][x] = skyBlue
			} else if value < 0.7 {
				pixels[y][x] = lightBlue
			} else {
				pixels[y][x] = veryLight
			}

			// Glossy highlights
			if rng.Float64() < 0.04 {
				pixels[y][x] = veryLight
			}
		}
	}

	return pixels
}

// ================================
// TRANSITION GENERATORS
// ================================

// GenerateGrassDirtTransition creates grass/dirt transition tile
func GenerateGrassDirtTransition() [][]int {
	pixels := make([][]int, 16)
	for i := range pixels {
		pixels[i] = make([]int, 16)
	}

	grassTile := GenerateGrass()
	dirtTile := GenerateDirt()

	rng := rand.New(rand.NewSource(100))

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			distFromDiagonal := float64(x+y) - 15.0
			transitionNoise := PerlinNoise(x, y, 0.4, 100) * 5.0
			transitionPoint := distFromDiagonal + transitionNoise

			if transitionPoint < -2 {
				pixels[y][x] = grassTile[y][x]
			} else if transitionPoint < 4 {
				if rng.Float64() < 0.6 {
					pixels[y][x] = grassTile[y][x]
				} else {
					pixels[y][x] = dirtTile[y][x]
				}
			} else {
				pixels[y][x] = dirtTile[y][x]
			}
		}
	}

	return pixels
}

// GenerateGrassSandTransition creates grass/sand transition tile
func GenerateGrassSandTransition() [][]int {
	pixels := make([][]int, 16)
	for i := range pixels {
		pixels[i] = make([]int, 16)
	}

	grassTile := GenerateGrass()
	sandTile := GenerateSand()

	rng := rand.New(rand.NewSource(101))

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			distFromDiagonal := float64(x+y) - 15.0
			transitionNoise := PerlinNoise(x, y, 0.4, 101) * 5.0
			transitionPoint := distFromDiagonal + transitionNoise

			if transitionPoint < -1 {
				pixels[y][x] = grassTile[y][x]
			} else if transitionPoint < 3 {
				if rng.Float64() < 0.55 {
					pixels[y][x] = grassTile[y][x]
				} else {
					pixels[y][x] = sandTile[y][x]
				}
			} else {
				pixels[y][x] = sandTile[y][x]
			}
		}
	}

	return pixels
}

// GenerateGrassStoneTransition creates grass/stone transition tile
func GenerateGrassStoneTransition() [][]int {
	pixels := make([][]int, 16)
	for i := range pixels {
		pixels[i] = make([]int, 16)
	}

	grassTile := GenerateGrass()
	stoneTile := GenerateStone()

	rng := rand.New(rand.NewSource(102))

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			distFromDiagonal := float64(x+y) - 15.0
			transitionNoise := PerlinNoise(x, y, 0.35, 102) * 6.0
			transitionPoint := distFromDiagonal + transitionNoise

			if transitionPoint < -2 {
				pixels[y][x] = grassTile[y][x]
			} else if transitionPoint < 5 {
				// Rocky patches with grass
				if rng.Float64() < 0.5 {
					pixels[y][x] = grassTile[y][x]
				} else {
					pixels[y][x] = stoneTile[y][x]
				}
			} else {
				pixels[y][x] = stoneTile[y][x]
			}
		}
	}

	return pixels
}

// GenerateSandShallowWaterTransition creates beach/shallow water transition
func GenerateSandShallowWaterTransition() [][]int {
	pixels := make([][]int, 16)
	for i := range pixels {
		pixels[i] = make([]int, 16)
	}

	sandTile := GenerateSand()
	shallowWaterTile := GenerateShallowWater()

	rng := rand.New(rand.NewSource(90))

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			distFromDiagonal := float64(x+y) - 15.0
			transitionNoise := PerlinNoise(x, y, 0.5, 90) * 4.0
			transitionPoint := distFromDiagonal + transitionNoise

			if transitionPoint < 0 {
				pixels[y][x] = sandTile[y][x]
			} else if transitionPoint < 3 {
				if rng.Float64() < 0.5 {
					pixels[y][x] = sandTile[y][x]
				} else {
					pixels[y][x] = shallowWaterTile[y][x]
				}
				// Foam/wave effect
				if rng.Float64() < 0.1 {
					pixels[y][x] = 35 // Very light blue (foam)
				}
			} else {
				pixels[y][x] = shallowWaterTile[y][x]
			}
		}
	}

	return pixels
}

// GenerateDirtStoneTransition creates dirt/stone transition tile
func GenerateDirtStoneTransition() [][]int {
	pixels := make([][]int, 16)
	for i := range pixels {
		pixels[i] = make([]int, 16)
	}

	dirtTile := GenerateDirt()
	stoneTile := GenerateStone()

	rng := rand.New(rand.NewSource(103))

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			distFromDiagonal := float64(x+y) - 15.0
			transitionNoise := PerlinNoise(x, y, 0.4, 103) * 5.5
			transitionPoint := distFromDiagonal + transitionNoise

			if transitionPoint < -1 {
				pixels[y][x] = dirtTile[y][x]
			} else if transitionPoint < 4 {
				if rng.Float64() < 0.55 {
					pixels[y][x] = dirtTile[y][x]
				} else {
					pixels[y][x] = stoneTile[y][x]
				}
			} else {
				pixels[y][x] = stoneTile[y][x]
			}
		}
	}

	return pixels
}

// GenerateStoneGravelTransition creates stone/gravel transition tile
func GenerateStoneGravelTransition() [][]int {
	pixels := make([][]int, 16)
	for i := range pixels {
		pixels[i] = make([]int, 16)
	}

	stoneTile := GenerateStone()
	gravelTile := GenerateGravel()

	rng := rand.New(rand.NewSource(104))

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			distFromDiagonal := float64(x+y) - 15.0
			transitionNoise := PerlinNoise(x, y, 0.45, 104) * 5.0
			transitionPoint := distFromDiagonal + transitionNoise

			if transitionPoint < -1 {
				pixels[y][x] = stoneTile[y][x]
			} else if transitionPoint < 4 {
				// Weathered stone to gravel
				if rng.Float64() < 0.6 {
					pixels[y][x] = stoneTile[y][x]
				} else {
					pixels[y][x] = gravelTile[y][x]
				}
			} else {
				pixels[y][x] = gravelTile[y][x]
			}
		}
	}

	return pixels
}

// ================================
// FILE SAVING
// ================================

// PixelData represents either palette indices or RGBA colors
type PixelData interface {
	GetColor(x, y int) color.RGBA
}

// PalettePixelData stores palette indices
type PalettePixelData struct {
	indices [][]int
	palette []color.RGBA
}

func (p *PalettePixelData) GetColor(x, y int) color.RGBA {
	idx := p.indices[y][x]
	if idx >= 0 && idx < len(p.palette) {
		return p.palette[idx]
	}
	return color.RGBA{0, 0, 0, 255} // Default to black
}

// RGBAPixelData stores RGBA colors directly
type RGBAPixelData struct {
	colors [][]color.RGBA
}

func (p *RGBAPixelData) GetColor(x, y int) color.RGBA {
	return p.colors[y][x]
}

// convertPaletteIndicesToRGBA converts palette indices to natural RGB colors
// This is used when generating tiles in full RGBA mode (not palette-limited)
func convertPaletteIndicesToRGBA(indices [][]int) [][]color.RGBA {
	colors := make([][]color.RGBA, len(indices))
	for y := range indices {
		colors[y] = make([]color.RGBA, len(indices[y]))
		for x := range indices[y] {
			idx := indices[y][x]
			// Map palette indices to natural RGB colors
			// These are approximate color mappings based on common palette usage
			colors[y][x] = paletteIndexToNaturalColor(idx)
		}
	}
	return colors
}

// paletteIndexToNaturalColor converts a palette index to a natural RGB color
// This provides realistic color representation when not using palette mode
func paletteIndexToNaturalColor(idx int) color.RGBA {
	// Map common palette indices to natural RGB colors
	// Based on typical RetroForge 50 palette usage
	// Note: Using the actual RetroForge 50 palette colors for accuracy
	colorMap := map[int]color.RGBA{
		// Grays
		0:  {0, 0, 0, 255},       // Black
		1:  {255, 255, 255, 255}, // White
		5:  {95, 87, 79, 255},    // Dark gray
		6:  {194, 195, 199, 255}, // Light gray
		37: {132, 126, 135, 255}, // Medium gray
		38: {105, 106, 106, 255}, // Dark gray
		39: {89, 86, 82, 255},    // Very dark gray
		46: {194, 195, 199, 255}, // Light gray
		// Browns
		18: {102, 57, 49, 255},   // Dark brown
		19: {143, 86, 59, 255},   // Medium brown
		20: {223, 113, 38, 255},  // Light brown
		21: {217, 160, 102, 255}, // Tan
		22: {238, 195, 154, 255}, // Light tan
		43: {108, 61, 30, 255},   // Dark olive brown
		44: {228, 181, 150, 255}, // Light tan
		45: {160, 177, 90, 255},  // Yellow-green (used for grass)
		// Greens
		3:  {0, 135, 81, 255},    // Dark green
		11: {0, 228, 54, 255},    // Green
		24: {153, 229, 80, 255},  // Light green
		25: {106, 190, 48, 255},  // Medium green
		26: {55, 148, 110, 255},  // Teal green
		27: {75, 105, 47, 255},   // Dark olive
		28: {82, 75, 36, 255},    // Olive brown
		29: {50, 60, 57, 255},    // Dark teal
		// Blues
		12: {41, 173, 255, 255},  // Blue
		30: {63, 63, 116, 255},   // Dark blue
		31: {48, 96, 130, 255},   // Ocean blue
		32: {91, 110, 225, 255},  // Bright blue
		33: {99, 155, 255, 255},  // Sky blue
		34: {95, 205, 228, 255},  // Light blue
		35: {203, 219, 252, 255}, // Very light blue
		36: {155, 173, 183, 255}, // Blue-gray
		// Yellows/Oranges
		8:  {255, 255, 137, 255},  // Yellow
		9:  {255, 163, 0, 255},   // Orange
		10: {255, 236, 39, 255},  // Bright yellow
		23: {251, 242, 54, 255},  // Bright yellow
		// Reds/Pinks
		2:  {255, 137, 137, 255}, // Red highlight
		4:  {195, 17, 17, 255},   // Red shadow
		14: {255, 119, 168, 255}, // Pink
		42: {255, 111, 177, 255}, // Pink-red
	}

	if c, ok := colorMap[idx]; ok {
		return c
	}

	// Default: interpolate based on index for unknown colors
	// Create a gradient from index value
	v := uint8((idx * 255) / 49)
	return color.RGBA{v, v, v, 255}
}

// convertToPixelData converts palette indices to appropriate PixelData type
func convertToPixelData(indices [][]int) PixelData {
	if usePalette && len(Palette) > 0 {
		return &PalettePixelData{
			indices: indices,
			palette: Palette,
		}
	}
	// RGBA mode: convert indices to natural colors
	return &RGBAPixelData{
		colors: convertPaletteIndicesToRGBA(indices),
	}
}

func SaveTilePNG(pixels PixelData, filename string) error {
	img := image.NewRGBA(image.Rect(0, 0, 16, 16))

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			img.Set(x, y, pixels.GetColor(x, y))
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

func SaveTilePNGScaled(pixels PixelData, filename string, scale int) error {
	size := 16 * scale
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			c := pixels.GetColor(x, y)

			for dy := 0; dy < scale; dy++ {
				for dx := 0; dx < scale; dx++ {
					img.Set(x*scale+dx, y*scale+dy, c)
				}
			}
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

// TilesetComposer manages the composite tileset image
type TilesetComposer struct {
	tilesetImg *image.RGBA
	tileSize   int
	tilesPerRow int
	currentTile int
	totalTiles  int
}

// NewTilesetComposer creates a new tileset composer
// tilesPerRow determines layout (e.g., 6 = 6 tiles wide)
func NewTilesetComposer(totalTiles, tileSize, tilesPerRow int) *TilesetComposer {
	width := tilesPerRow * tileSize
	height := ((totalTiles + tilesPerRow - 1) / tilesPerRow) * tileSize // Ceiling division
	
	return &TilesetComposer{
		tilesetImg: image.NewRGBA(image.Rect(0, 0, width, height)),
		tileSize:   tileSize,
		tilesPerRow: tilesPerRow,
		currentTile: 0,
		totalTiles:  totalTiles,
	}
}

// AddTile adds a tile to the tileset at the current position
func (tc *TilesetComposer) AddTile(pixels PixelData) {
	if tc.currentTile >= tc.totalTiles {
		return
	}
	
	row := tc.currentTile / tc.tilesPerRow
	col := tc.currentTile % tc.tilesPerRow
	xOffset := col * tc.tileSize
	yOffset := row * tc.tileSize
	
	// Copy tile pixels to tileset
	for y := 0; y < tc.tileSize; y++ {
		for x := 0; x < tc.tileSize; x++ {
			c := pixels.GetColor(x, y)
			tc.tilesetImg.Set(xOffset+x, yOffset+y, c)
		}
	}
	
	tc.currentTile++
}

// Save saves the tileset image to a file
func (tc *TilesetComposer) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, tc.tilesetImg)
}

// GetImage returns the tileset image for color analysis
func (tc *TilesetComposer) GetImage() image.Image {
	return tc.tilesetImg
}

// CountUniqueColors counts the number of unique colors in the tileset
func CountUniqueColors(img image.Image) int {
	colorMap := make(map[color.RGBA]bool)
	bounds := img.Bounds()
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			r, g, b, a := c.RGBA()
			// Convert from 16-bit to 8-bit
			rgba := color.RGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: uint8(a >> 8),
			}
			colorMap[rgba] = true
		}
	}
	
	return len(colorMap)
}

// GeneratePaletteFromImage generates a 48-color game palette from the tileset image
// Note: Built-in colors (0-15) are always available in the engine
// This function generates only the game palette (48 colors, indices 16-63)
func GeneratePaletteFromImage(img image.Image) ([]string, error) {
	// Use imgtool quantizer to get 48 colors (game palette only)
	opts := imgtool.QuantizeOptions{
		EnforceBlackWhite: false, // Built-in colors handle black/white
		AlphaThreshold:    128,
		DitherAlgorithm:   "none", // No dithering for palette extraction
	}
	
	palette, err := imgtool.Quantize(img, opts)
	if err != nil {
		return nil, err
	}
	
	// Ensure exactly 48 colors (quantizer should do this, but double-check)
	hexColors := make([]string, 48)
	for i := 0; i < 48 && i < len(palette.Colors); i++ {
		hexColors[i] = palette.Colors[i]
	}
	
	// Fill remaining slots if needed (shouldn't happen, but safety check)
	for i := len(palette.Colors); i < 48; i++ {
		v := uint8(32 + (i * 223 / 47)) // Range from 32 to 255
		hexColors[i] = fmt.Sprintf("#%02x%02x%02x", v, v, v)
	}
	
	return hexColors, nil
}

// SavePaletteJSON saves a palette to JSON file in the same format as palette.json
func SavePaletteJSON(colors []string, filename string) error {
	paletteData := map[string]interface{}{
		"colors": colors,
	}
	
	data, err := json.MarshalIndent(paletteData, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(filename, data, 0644)
}

// ================================
// MAIN
// ================================

func main() {
	rand.Seed(time.Now().UnixNano())

	// Parse command-line flags
	paletteFlag := flag.String("palette", "", "Game palette name to use (e.g., 'RetroForge 48'). Leave empty for full RGBA mode (16M colors)")
	outputDirFlag := flag.String("output", "", "Output directory for generated tiles (default: design/tiles relative to engine root)")
	flag.Parse()

	// Set up palette mode
	if *paletteFlag != "" {
		// Use palette mode
		usePalette = true
		gamePalette := pal.GetPalette(*paletteFlag)
		if len(gamePalette) == 0 {
			println("  WARNING: Palette '" + *paletteFlag + "' not found, using default")
			gamePalette = pal.GetPalette("RetroForge 48")
		}
		// Create full 64-color palette (16 built-in + 48 game palette)
		fullPalette := pal.GetFullPalette(*paletteFlag)
		Palette = fullPalette // Use full 64-color palette for rendering
		println("==============================================")
		println("  RetroForge Terrain Tile Generator")
		println("  Tier 1 Complete Tileset (18 tiles)")
		println("  Seamless Tiling Enabled")
		println("  Mode: PALETTE (" + *paletteFlag + ")")
		println("  Colors: 64 total (16 built-in + 48 game palette)")
		println("==============================================")
	} else {
		// Use RGBA mode (16M colors)
		usePalette = false
		Palette = nil
		println("==============================================")
		println("  RetroForge Terrain Tile Generator")
		println("  Tier 1 Complete Tileset (18 tiles)")
		println("  Seamless Tiling Enabled")
		println("  Mode: RGBA (16M colors)")
		println("==============================================")
	}
	println()

	// Determine output directory
	var outputDir string
	if *outputDirFlag != "" {
		// Use provided path (absolute or relative)
		outputDir = *outputDirFlag
	} else {
		// Default: design/tiles relative to engine root
		// terraingen is in cmd/terraingen, so go up 2 levels to retroforge-engine root
		// Then go up 1 level to RetroForge root, then into design/tiles
		if wd, err := os.Getwd(); err == nil {
			// Try to find design/tiles relative to current working directory
			// Check if we're in retroforge-engine directory
			if _, err := os.Stat(filepath.Join(wd, "design", "tiles")); err == nil {
				outputDir = filepath.Join(wd, "design", "tiles")
			} else {
				// Try going up one level (assuming we're in retroforge-engine)
				parentDir := filepath.Dir(wd)
				if _, err := os.Stat(filepath.Join(parentDir, "design", "tiles")); err == nil {
					outputDir = filepath.Join(parentDir, "design", "tiles")
				} else {
					// Fallback: use current directory
					outputDir = wd
				}
			}
		} else {
			outputDir = "."
		}
	}

	// Convert to absolute path
	absOutputDir, err := filepath.Abs(outputDir)
	if err != nil {
		absOutputDir = outputDir
	}
	outputDir = absOutputDir

	// Create output directory if it doesn't exist (overwrite if exists is default behavior)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		println("  ERROR: Could not create output directory:", outputDir)
		println("  Using current directory instead")
		outputDir = "."
	}

	// Create previews subdirectory
	previewsDir := filepath.Join(outputDir, "previews")
	if err := os.MkdirAll(previewsDir, 0755); err != nil {
		println("  ERROR: Could not create previews directory:", previewsDir)
		println("  Previews will be saved to output directory")
		previewsDir = outputDir
	}

	println("Output directory:", outputDir)
	println("Previews directory:", previewsDir)
	println()

	// Base terrains
	baseTiles := map[string]func() [][]int{
		"grass":         GenerateGrass,
		"dark_grass":    GenerateDarkGrass,
		"dirt":          GenerateDirt,
		"sand":          GenerateSand,
		"gravel":        GenerateGravel,
		"stone":         GenerateStone,
		"shallow_water": GenerateShallowWater,
		"water":         GenerateWater,
		"deep_water":    GenerateDeepWater,
		"mud":           GenerateMud,
		"snow":          GenerateSnow,
		"ice":           GenerateIce,
	}

	// Transition tiles
	transitionTiles := map[string]func() [][]int{
		"grass_to_dirt":          GenerateGrassDirtTransition,
		"grass_to_sand":          GenerateGrassSandTransition,
		"grass_to_stone":         GenerateGrassStoneTransition,
		"sand_to_shallow_water":  GenerateSandShallowWaterTransition,
		"dirt_to_stone":          GenerateDirtStoneTransition,
		"stone_to_gravel":        GenerateStoneGravelTransition,
	}

	// Calculate total tiles for tileset composition
	totalTiles := len(baseTiles) + len(transitionTiles)
	
	// Create tileset composer (6 tiles per row for 18 tiles = 3 rows)
	// Layout: 6 tiles wide × 3 rows = 18 tiles (6×16 = 96px wide, 3×16 = 48px tall)
	tilesetComposer := NewTilesetComposer(totalTiles, 16, 6)

	println("Generating BASE TERRAIN tiles (12)...")
	println()
	for name, genFunc := range baseTiles {
		indices := genFunc()
		pixelData := convertToPixelData(indices)

		// Save main tile to output directory
		outputPath := filepath.Join(outputDir, name+"_16x16.png")
		err := SaveTilePNG(pixelData, outputPath)
		if err != nil {
			println("  ERROR:", name+"_16x16.png:", err.Error())
			continue
		}

		// Add tile to tileset composition
		tilesetComposer.AddTile(pixelData)

		// Save preview to previews subdirectory
		previewPath := filepath.Join(previewsDir, name+"_preview.png")
		err = SaveTilePNGScaled(pixelData, previewPath, 8)
		if err != nil {
			println("  ERROR:", name+"_preview.png:", err.Error())
			continue
		}

		println("  ✓", name)
	}

	println()
	println("Generating TRANSITION tiles (6)...")
	println()
	for name, genFunc := range transitionTiles {
		indices := genFunc()
		pixelData := convertToPixelData(indices)

		// Save main tile to output directory
		outputPath := filepath.Join(outputDir, name+"_16x16.png")
		err := SaveTilePNG(pixelData, outputPath)
		if err != nil {
			println("  ERROR:", name+"_16x16.png:", err.Error())
			continue
		}

		// Add tile to tileset composition
		tilesetComposer.AddTile(pixelData)

		// Save preview to previews subdirectory
		previewPath := filepath.Join(previewsDir, name+"_preview.png")
		err = SaveTilePNGScaled(pixelData, previewPath, 8)
		if err != nil {
			println("  ERROR:", name+"_preview.png:", err.Error())
			continue
		}

		println("  ✓", name)
	}

	// Save tileset.png (1:1 sprite sheet)
	println()
	println("Creating tileset.png (composite sprite sheet)...")
	tilesetPath := filepath.Join(outputDir, "tileset.png")
	if err := tilesetComposer.Save(tilesetPath); err != nil {
		println("  ERROR: Failed to save tileset.png:", err.Error())
	} else {
		println("  ✓ tileset.png saved")
	}

	// Count unique colors in tileset
	println()
	println("Analyzing tileset colors...")
	colorCount := CountUniqueColors(tilesetComposer.GetImage())
	println("  Unique colors in tileset:", colorCount)

	// Generate and save custom palette
	println()
	println("Generating custom game palette (48 colors for indices 16-63)...")
	paletteColors, err := GeneratePaletteFromImage(tilesetComposer.GetImage())
	if err != nil {
		println("  ERROR: Failed to generate palette:", err.Error())
	} else {
		palettePath := filepath.Join(outputDir, "tileset_palette.json")
		if err := SavePaletteJSON(paletteColors, palettePath); err != nil {
			println("  ERROR: Failed to save palette:", err.Error())
		} else {
			println("  ✓ tileset_palette.json saved (48 colors, game palette)")
			println("    Note: Built-in colors (0-15) are always available in engine")
		}
	}

	println()
	println("==============================================")
	println("  GENERATION COMPLETE!")
	println("==============================================")
	println()
	println("Generated 18 tiles total:")
	println()
	println("BASE TERRAINS (12):")
	println("  • grass              - Standard grass")
	println("  • dark_grass         - Forest floor")
	println("  • dirt               - Bare earth")
	println("  • sand               - Beach sand")
	println("  • gravel             - Rocky ground")
	println("  • stone              - Solid rock")
	println("  • shallow_water      - Coastal water")
	println("  • water              - Normal water")
	println("  • deep_water         - Ocean depths")
	println("  • mud                - Wetland mud")
	println("  • snow               - Snow covered")
	println("  • ice                - Frozen surface")
	println()
	println("TRANSITIONS (6):")
	println("  • grass_to_dirt      - Grass → Dirt")
	println("  • grass_to_sand      - Grass → Sand")
	println("  • grass_to_stone     - Grass → Stone")
	println("  • sand_to_shallow_water - Beach → Water")
	println("  • dirt_to_stone      - Dirt → Stone")
	println("  • stone_to_gravel    - Stone → Gravel")
	println()
	println("Files created:")
	println("  • *_16x16.png           - Individual tile files (use in isometric converter)")
	println("  • tileset.png          - Composite sprite sheet (1:1, all tiles)")
	println("  • tileset_palette.json  - Custom 48-color game palette (indices 16-63)")
	println("  • previews/*_preview.png - 128x128 scaled previews in previews/ subfolder")
	println()
	println("Workflow Options:")
	println("  MODERN: Use tileset.png directly (full color, 16M colors)")
	println("  RETRO:  Use tileset_palette.json + individual tiles (64 colors: 16 built-in + 48 game)")
	println()
	println("Palette System:")
	println("  • Built-in colors (0-15): Always available (grayscale + primary colors)")
	println("  • Game palette (16-63): From tileset_palette.json or predefined palette")
	println("  • Total: 64 colors available at any time")
	println()
	println("Usage:")
	println("  terraingen                                    - Generate tiles in RGBA mode (16M colors)")
	println("  terraingen -palette \"RetroForge 48\"          - Generate tiles using specified game palette")
	println("  terraingen -output \"/path/to/tiles\"         - Specify output directory")
	println("  terraingen -palette \"RetroForge 48\" -output \"design/tiles\" - Combined options")
	println()
}

