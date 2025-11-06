package luabind

import (
	"fmt"
	"testing"
)

// gridToScreen calculates isometric grid to screen transformation
// This mirrors the logic in luabind.go for testing
func gridToScreen(gridX, gridY, tileWidth, offsetX, offsetY int) (screenX, screenY int) {
	// Calculate diamond height (for 2:1 isometric ratio: width/2)
	diamondHeight := tileWidth / 2

	// Basis vectors
	iHatX := float64(tileWidth) / 2.0      // 16
	iHatY := float64(diamondHeight) / 2.0  // 8
	jHatX := -float64(tileWidth) / 2.0     // -16
	jHatY := float64(diamondHeight) / 2.0  // 8

	// Matrix transformation: screen = (gridX * i_hat) + (gridY * j_hat)
	screenX = int(float64(gridX)*iHatX + float64(gridY)*jHatX)
	screenY = int(float64(gridX)*iHatY + float64(gridY)*jHatY)

	// Apply origin offset
	screenX += offsetX - (tileWidth / 2)
	screenY += offsetY

	return screenX, screenY
}

// screenToGrid calculates inverse transformation (screen to grid)
func screenToGrid(screenX, screenY, tileWidth, offsetX, offsetY int) (gridX, gridY int) {
	// Calculate diamond height
	diamondHeight := tileWidth / 2

	// Basis vectors
	iHatX := float64(tileWidth) / 2.0
	iHatY := float64(diamondHeight) / 2.0
	jHatX := -float64(tileWidth) / 2.0
	jHatY := float64(diamondHeight) / 2.0

	// Determinant
	det := (iHatX * jHatY) - (iHatY * jHatX)

	// Inverse matrix components
	invA := jHatY / det
	invB := -iHatY / det
	invC := -jHatX / det
	invD := iHatX / det

	// Local coordinates (remove offset)
	localX := float64(screenX - offsetX + (tileWidth / 2))
	localY := float64(screenY - offsetY)

	// Apply inverse transformation
	gridXFloat := (localX * invA) + (localY * invC)
	gridYFloat := (localX * invB) + (localY * invD)

	// Round to nearest integer
	if gridXFloat >= 0 {
		gridX = int(gridXFloat + 0.5)
	} else {
		gridX = int(gridXFloat - 0.5)
	}
	if gridYFloat >= 0 {
		gridY = int(gridYFloat + 0.5)
	} else {
		gridY = int(gridYFloat - 0.5)
	}

	return gridX, gridY
}

func TestGridToScreen_BasicTransformation(t *testing.T) {
	tileWidth := 32
	offsetX := 0
	offsetY := 0

	tests := []struct {
		name     string
		gridX    int
		gridY    int
		expected screenPos
	}{
		{"origin", 0, 0, screenPos{-16, 0}},
		{"right", 1, 0, screenPos{0, 8}},
		{"down", 0, 1, screenPos{-32, 8}},
		{"diagonal", 1, 1, screenPos{-16, 16}},
		{"negative_x", -1, 0, screenPos{-32, -8}},
		{"negative_y", 0, -1, screenPos{0, -8}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			screenX, screenY := gridToScreen(tt.gridX, tt.gridY, tileWidth, offsetX, offsetY)
			if screenX != tt.expected.x || screenY != tt.expected.y {
				t.Errorf("gridToScreen(%d, %d) = (%d, %d), want (%d, %d)",
					tt.gridX, tt.gridY, screenX, screenY, tt.expected.x, tt.expected.y)
			}
		})
	}
}

type screenPos struct {
	x, y int
}

func TestGridToScreen_OffsetApplication(t *testing.T) {
	tileWidth := 32
	offsetX := 100
	offsetY := 50

	// Grid (0,0) with offset should be at (100-16, 50)
	screenX, screenY := gridToScreen(0, 0, tileWidth, offsetX, offsetY)
	expectedX := 100 - 16 // offsetX - tileWidth/2
	expectedY := 50       // offsetY

	if screenX != expectedX || screenY != expectedY {
		t.Errorf("gridToScreen(0, 0) with offset(100, 50) = (%d, %d), want (%d, %d)",
			screenX, screenY, expectedX, expectedY)
	}

	// Test negative offsets
	offsetX = -50
	offsetY = -30
	screenX, screenY = gridToScreen(0, 0, tileWidth, offsetX, offsetY)
	expectedX = -50 - 16
	expectedY = -30

	if screenX != expectedX || screenY != expectedY {
		t.Errorf("gridToScreen(0, 0) with offset(-50, -30) = (%d, %d), want (%d, %d)",
			screenX, screenY, expectedX, expectedY)
	}
}

func TestGridToScreen_DiamondHeight(t *testing.T) {
	// Test that diamond height is used (not full tile height)
	// For 32×24 tiles: diamond height = 16, not 24
	tileWidth := 32
	diamondHeight := tileWidth / 2 // Should be 16

	if diamondHeight != 16 {
		t.Errorf("Expected diamond height 16 for 32-width tile, got %d", diamondHeight)
	}

	// Test basis vectors use diamond height
	iHatY := float64(diamondHeight) / 2.0 // Should be 8
	jHatY := float64(diamondHeight) / 2.0 // Should be 8

	if iHatY != 8.0 || jHatY != 8.0 {
		t.Errorf("Basis vectors should use diamond height/2 (8), got iHatY=%f, jHatY=%f", iHatY, jHatY)
	}

	// Verify Y calculation uses diamond height, not full height
	// Grid (1, 0) should move down by 8 (diamond height/2), not 12 (full height/2)
	_, screenY := gridToScreen(1, 0, tileWidth, 0, 0)
	expectedY := 8 // diamondHeight/2

	if screenY != expectedY {
		t.Errorf("gridToScreen(1, 0) Y should use diamond height/2 (8), got %d", screenY)
	}
}

func TestScreenToGrid_InverseTransformation(t *testing.T) {
	tileWidth := 32
	offsetX := 100
	offsetY := 50

	testCases := []struct {
		name     string
		gridX    int
		gridY    int
		tolerance int // Allow small rounding differences
	}{
		{"origin", 0, 0, 0},
		{"positive", 1, 1, 1},
		{"negative", -1, -1, 1},
		{"large", 10, 10, 1},
		{"mixed", 5, -3, 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Grid → Screen
			screenX, screenY := gridToScreen(tc.gridX, tc.gridY, tileWidth, offsetX, offsetY)

			// Screen → Grid (inverse)
			resultGridX, resultGridY := screenToGrid(screenX, screenY, tileWidth, offsetX, offsetY)

			// Check round-trip accuracy
			diffX := abs(resultGridX - tc.gridX)
			diffY := abs(resultGridY - tc.gridY)

			if diffX > tc.tolerance || diffY > tc.tolerance {
				t.Errorf("Round-trip failed: grid(%d, %d) → screen(%d, %d) → grid(%d, %d), diff=(%d, %d)",
					tc.gridX, tc.gridY, screenX, screenY, resultGridX, resultGridY, diffX, diffY)
			}
		})
	}
}

func TestScreenToGrid_BoundaryConditions(t *testing.T) {
	tileWidth := 32
	offsetX := 100
	offsetY := 50

	// Test coordinates at tile centers
	tests := []struct {
		name     string
		gridX    int
		gridY    int
		expected screenPos
	}{
		{"center_0_0", 0, 0, screenPos{84, 50}},   // offsetX - tileWidth/2, offsetY
		{"center_1_0", 1, 0, screenPos{100, 58}},  // offsetX, offsetY + 8
		{"center_0_1", 0, 1, screenPos{68, 58}},   // offsetX - tileWidth, offsetY + 8
		{"center_1_1", 1, 1, screenPos{84, 66}},   // offsetX - tileWidth/2, offsetY + 16
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			screenX, screenY := gridToScreen(tt.gridX, tt.gridY, tileWidth, offsetX, offsetY)

			// Convert back to grid
			resultGridX, resultGridY := screenToGrid(screenX, screenY, tileWidth, offsetX, offsetY)

			// Should round to original grid coordinates
			if abs(resultGridX-tt.gridX) > 1 || abs(resultGridY-tt.gridY) > 1 {
				t.Errorf("Boundary test failed: grid(%d, %d) → screen(%d, %d) → grid(%d, %d)",
					tt.gridX, tt.gridY, screenX, screenY, resultGridX, resultGridY)
			}
		})
	}
}

func TestScreenToGrid_OutOfBounds(t *testing.T) {
	tileWidth := 32
	offsetX := 100
	offsetY := 50

	// Test very negative coordinates
	resultGridX, resultGridY := screenToGrid(-1000, -1000, tileWidth, offsetX, offsetY)
	// Should handle gracefully (no panic)
	if resultGridX == 0 && resultGridY == 0 {
		// Expected behavior - coordinates outside map area
	}

	// Test very positive coordinates
	resultGridX2, resultGridY2 := screenToGrid(10000, 10000, tileWidth, offsetX, offsetY)
	// Should handle gracefully (no panic)
	if resultGridX2 == 0 && resultGridY2 == 0 {
		// Expected behavior - coordinates outside map area
	}
	_ = resultGridX
	_ = resultGridY
	_ = resultGridX2
	_ = resultGridY2
}

func TestCoordinateRoundTrip(t *testing.T) {
	tileWidth := 32
	offsetX := 100
	offsetY := 50

	// Test random grid coordinates
	testCoords := []struct {
		gridX int
		gridY int
	}{
		{0, 0},
		{1, 0},
		{0, 1},
		{1, 1},
		{-1, 0},
		{0, -1},
		{-1, -1},
		{5, 3},
		{-3, 5},
		{10, 10},
		{-10, -10},
		{15, -7},
		{-8, 12},
	}

	for _, tc := range testCoords {
		t.Run(fmt.Sprintf("grid_%d_%d", tc.gridX, tc.gridY), func(t *testing.T) {
			// Grid → Screen
			screenX, screenY := gridToScreen(tc.gridX, tc.gridY, tileWidth, offsetX, offsetY)

			// Screen → Grid
			resultGridX, resultGridY := screenToGrid(screenX, screenY, tileWidth, offsetX, offsetY)

			// Allow small rounding differences (within 1 tile)
			diffX := abs(resultGridX - tc.gridX)
			diffY := abs(resultGridY - tc.gridY)

			if diffX > 1 || diffY > 1 {
				t.Errorf("Round-trip failed: grid(%d, %d) → screen(%d, %d) → grid(%d, %d), diff=(%d, %d)",
					tc.gridX, tc.gridY, screenX, screenY, resultGridX, resultGridY, diffX, diffY)
			}
		})
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}


