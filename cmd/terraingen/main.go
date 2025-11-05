package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"time"
)

// RetroForge 50-color palette
var Palette = []color.RGBA{
	{0, 0, 0, 255},           // 0: Black
	{29, 43, 83, 255},        // 1: Dark blue
	{126, 37, 83, 255},       // 2: Dark purple
	{0, 135, 81, 255},        // 3: Dark green
	{171, 82, 54, 255},       // 4: Brown
	{95, 87, 79, 255},        // 5: Dark gray
	{194, 195, 199, 255},     // 6: Light gray
	{255, 241, 232, 255},     // 7: White
	{255, 0, 77, 255},        // 8: Red
	{255, 163, 0, 255},       // 9: Orange
	{255, 236, 39, 255},      // 10: Yellow
	{0, 228, 54, 255},        // 11: Green
	{41, 173, 255, 255},      // 12: Blue
	{131, 118, 156, 255},     // 13: Lavender
	{255, 119, 168, 255},     // 14: Pink
	{255, 204, 170, 255},     // 15: Peach
	{34, 32, 52, 255},        // 16: Very dark blue
	{69, 40, 60, 255},        // 17: Dark purple-brown
	{102, 57, 49, 255},       // 18: Dark brown
	{143, 86, 59, 255},       // 19: Medium brown
	{223, 113, 38, 255},      // 20: Light brown/orange
	{217, 160, 102, 255},     // 21: Tan
	{238, 195, 154, 255},     // 22: Light tan
	{251, 242, 54, 255},      // 23: Bright yellow
	{153, 229, 80, 255},      // 24: Light green
	{106, 190, 48, 255},      // 25: Medium green
	{55, 148, 110, 255},      // 26: Teal green
	{75, 105, 47, 255},       // 27: Dark olive
	{82, 75, 36, 255},        // 28: Olive brown
	{50, 60, 57, 255},        // 29: Dark teal
	{63, 63, 116, 255},       // 30: Dark blue
	{48, 96, 130, 255},       // 31: Ocean blue
	{91, 110, 225, 255},      // 32: Bright blue
	{99, 155, 255, 255},      // 33: Sky blue
	{95, 205, 228, 255},      // 34: Light blue
	{203, 219, 252, 255},     // 35: Very light blue
	{155, 173, 183, 255},     // 36: Blue-gray
	{132, 126, 135, 255},     // 37: Medium gray
	{105, 106, 106, 255},     // 38: Dark gray
	{89, 86, 82, 255},        // 39: Very dark gray
	{118, 66, 138, 255},      // 40: Purple
	{172, 50, 50, 255},       // 41: Dark red
	{217, 87, 99, 255},       // 42: Pink-red
	{215, 123, 186, 255},     // 43: Pink
	{143, 151, 74, 255},      // 44: Yellow-green
	{138, 111, 48, 255},      // 45: Gold-brown
	{194, 195, 199, 255},     // 46: Light gray
	{255, 255, 255, 255},     // 47: Pure white
	{0, 0, 0, 255},           // 48: Black
	{0, 0, 0, 255},           // 49: Black
}

// ================================
// NOISE FUNCTIONS
// ================================

func PerlinNoise(x, y int, frequency float64, seed int64) float64 {
	fx := float64(x) * frequency
	fy := float64(y) * frequency

	x0 := int(math.Floor(fx))
	y0 := int(math.Floor(fy))
	x1 := x0 + 1
	y1 := y0 + 1

	sx := fx - float64(x0)
	sy := fy - float64(y0)

	n00 := dotGridGradient(x0, y0, fx, fy, seed)
	n10 := dotGridGradient(x1, y0, fx, fy, seed)
	n01 := dotGridGradient(x0, y1, fx, fy, seed)
	n11 := dotGridGradient(x1, y1, fx, fy, seed)

	sx = fade(sx)
	sy = fade(sy)

	ix0 := lerp(n00, n10, sx)
	ix1 := lerp(n01, n11, sx)

	return lerp(ix0, ix1, sy)
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

// GenerateGrass creates a textured grass tile
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
			noise := PerlinNoise(x, y, 0.4, 42)
			noise2 := PerlinNoise(x, y, 0.8, 43) * 0.3
			gradient := float64(y) / 20.0
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

func SaveTilePNG(pixels [][]int, filename string) error {
	img := image.NewRGBA(image.Rect(0, 0, 16, 16))

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			idx := pixels[y][x]
			if idx >= 0 && idx < len(Palette) {
				img.Set(x, y, Palette[idx])
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

func SaveTilePNGScaled(pixels [][]int, filename string, scale int) error {
	size := 16 * scale
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			idx := pixels[y][x]
			var c color.RGBA
			if idx >= 0 && idx < len(Palette) {
				c = Palette[idx]
			}

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

// ================================
// MAIN
// ================================

func main() {
	rand.Seed(time.Now().UnixNano())

	println("==============================================")
	println("  RetroForge Terrain Tile Generator")
	println("  Tier 1 Complete Tileset (18 tiles)")
	println("==============================================")
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

	println("Generating BASE TERRAIN tiles (12)...")
	println()
	for name, genFunc := range baseTiles {
		pixels := genFunc()

		err := SaveTilePNG(pixels, name+"_16x16.png")
		if err != nil {
			println("  ERROR:", name+"_16x16.png:", err.Error())
			continue
		}

		err = SaveTilePNGScaled(pixels, name+"_preview.png", 8)
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
		pixels := genFunc()

		err := SaveTilePNG(pixels, name+"_16x16.png")
		if err != nil {
			println("  ERROR:", name+"_16x16.png:", err.Error())
			continue
		}

		err = SaveTilePNGScaled(pixels, name+"_preview.png", 8)
		if err != nil {
			println("  ERROR:", name+"_preview.png:", err.Error())
			continue
		}

		println("  ✓", name)
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
	println("  • *_16x16.png        - Use these in your isometric converter")
	println("  • *_preview.png      - 128x128 scaled previews")
	println()
}

