package input

// Universal 11-button input system for cross-platform support:
// SELECT, START, UP, DOWN, LEFT, RIGHT, A, B, X, Y, TURBO

const (
	BtnSELECT = 0
	BtnSTART  = 1
	BtnUP     = 2
	BtnDOWN   = 3
	BtnLEFT   = 4
	BtnRIGHT  = 5
	BtnA      = 6
	BtnB      = 7
	BtnX      = 8
	BtnY      = 9
	BtnTURBO  = 10
	num       = 11

	// Legacy button constants for backward compatibility
	BtnLeft  = BtnLEFT  // Deprecated: use BtnLEFT
	BtnRight = BtnRIGHT // Deprecated: use BtnRIGHT
	BtnUp    = BtnUP    // Deprecated: use BtnUP
	BtnDown  = BtnDOWN  // Deprecated: use BtnDOWN
	BtnO     = BtnA     // Deprecated: use BtnA
)

var cur [num]bool
var prev [num]bool
var justPressed [num]bool // Tracks buttons that were just pressed (for WASM async edge detection)

// Step copies current button state to previous for edge detection
// Call this at the START of each frame (before HandleInput)
func Step() { 
	prev = cur 
	// NOTE: justPressed is NOT cleared here - it persists until consumed by Btnp()
	// This ensures async WASM input (rf_set_btn) works correctly even if called
	// just before Step() runs
	// After Step(), prev now has the state from last frame
	// cur will be updated by Set() calls during/after the frame
	// When HandleInput() runs, btnp() compares: cur (new) vs prev (old)
	logStep() // Will be no-op in non-WASM builds (input_noop.go)
}

func Set(i int, down bool) { 
	if i >= 0 && i < num { 
		// Track edge: if button transitions from false to true, mark as justPressed
		if down && !cur[i] {
			justPressed[i] = true
		}
		cur[i] = down 
	} 
}

func Btn(i int) bool { 
	if i < 0 || i >= num { 
		return false 
	}
	return cur[i] 
}

func Btnp(i int) bool {
	if i < 0 || i >= num {
		return false
	}
	// Use justPressed flag for reliable edge detection (works for both SDL and WASM)
	// Fall back to prev comparison if justPressed wasn't set (backward compatibility)
	result := justPressed[i] || (cur[i] && !prev[i])
	
	// Clear justPressed flag after consuming it (one-shot detection)
	// This ensures the press is detected once per frame, even if Btnp() is called multiple times
	if justPressed[i] {
		justPressed[i] = false
	}
	
	// WASM-specific logging is handled in input_wasm.go via build tags
	// Desktop builds use input_noop.go (no logging)
	logBtnp(i, cur[i], prev[i], false) // Will be no-op in non-WASM builds
	if !result && cur[i] {
		logBtnpFalse(i, cur[i], prev[i]) // Log missed edges in WASM
	}
	return result
}

// ClearStaleJustPressed clears any justPressed flags that weren't consumed
// Call this at the END of each frame (after Draw) to prevent stale flags
func ClearStaleJustPressed() {
	for i := range justPressed {
		justPressed[i] = false
	}
}

// Shift returns true if TURBO button is pressed (backward compatibility)
func Shift() bool { return Btn(BtnTURBO) }

// SetShift sets TURBO button (backward compatibility)
func SetShift(down bool) { Set(BtnTURBO, down) }

// SetByName sets button state by name (e.g., "UP", "A", "TURBO")
func SetByName(name string, down bool) bool {
	switch name {
	case "SELECT":
		Set(BtnSELECT, down)
		return true
	case "START":
		Set(BtnSTART, down)
		return true
	case "UP":
		Set(BtnUP, down)
		return true
	case "DOWN":
		Set(BtnDOWN, down)
		return true
	case "LEFT":
		Set(BtnLEFT, down)
		return true
	case "RIGHT":
		Set(BtnRIGHT, down)
		return true
	case "A":
		Set(BtnA, down)
		return true
	case "B":
		Set(BtnB, down)
		return true
	case "X":
		Set(BtnX, down)
		return true
	case "Y":
		Set(BtnY, down)
		return true
	case "TURBO":
		Set(BtnTURBO, down)
		return true
	}
	return false
}
