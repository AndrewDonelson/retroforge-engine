package input

import "testing"

func TestInputBasic(t *testing.T) {
	// Reset state
	Step()

	// Test Set
	Set(BtnLeft, true)
	if !Btn(BtnLeft) {
		t.Fatalf("BtnLeft should be true after Set")
	}
	if Btn(BtnRight) {
		t.Fatalf("BtnRight should be false")
	}

	Set(BtnLeft, false)
	if Btn(BtnLeft) {
		t.Fatalf("BtnLeft should be false after Set(false)")
	}
}

func TestBtnp(t *testing.T) {
	// Reset state
	Step()

	// Initially not pressed
	if Btnp(BtnUp) {
		t.Fatalf("BtnUp should not be pressed initially")
	}

	// Set button down
	Set(BtnUp, true)
	if !Btn(BtnUp) {
		t.Fatalf("BtnUp should be true")
	}
	if !Btnp(BtnUp) {
		t.Fatalf("BtnUp should be pressed on first frame")
	}

	// Step to next frame - should not trigger btnp again
	Step()
	if Btnp(BtnUp) {
		t.Fatalf("BtnUp should not trigger btnp on second frame")
	}
	if !Btn(BtnUp) {
		t.Fatalf("BtnUp should still be true")
	}

	// Release and press again
	Set(BtnUp, false)
	Step()
	Set(BtnUp, true)
	if !Btnp(BtnUp) {
		t.Fatalf("BtnUp should trigger btnp after release and press")
	}
}

func TestInvalidInput(t *testing.T) {
	// Test negative index
	Set(-1, true)
	if Btn(-1) {
		t.Fatalf("negative index should return false")
	}
	if Btnp(-1) {
		t.Fatalf("negative index btnp should return false")
	}

	// Test out of bounds
	Set(999, true)
	if Btn(999) {
		t.Fatalf("out of bounds index should return false")
	}
	if Btnp(999) {
		t.Fatalf("out of bounds index btnp should return false")
	}
}

func TestAllButtons(t *testing.T) {
	Step()

	buttons := []int{BtnLeft, BtnRight, BtnUp, BtnDown, BtnO, BtnX}
	for i, btn := range buttons {
		Set(btn, true)
		if !Btn(btn) {
			t.Fatalf("button %d should be true", i)
		}
		Set(btn, false)
		if Btn(btn) {
			t.Fatalf("button %d should be false", i)
		}
	}
}

func TestEdgeCases(t *testing.T) {
	Step()

	// Test boundary values (0 and num-1 are valid)
	Set(0, true)
	if !Btn(0) {
		t.Fatalf("button 0 (first valid) should work")
	}

	Set(num-1, true)
	if !Btn(num - 1) {
		t.Fatalf("button %d (last valid) should work", num-1)
	}

	// Test exactly at boundary (num is invalid)
	Set(num, true)
	if Btn(num) {
		t.Fatalf("button %d (out of bounds) should return false", num)
	}

	// Test large negative values
	Set(-100, true)
	if Btn(-100) {
		t.Fatalf("button -100 should return false")
	}

	// Test maximum int
	Set(2147483647, true)
	if Btn(2147483647) {
		t.Fatalf("button max int should return false")
	}
}

func TestRapidToggle(t *testing.T) {
	Step()

	// Rapidly toggle button multiple times
	for i := 0; i < 100; i++ {
		Set(BtnX, true)
		Set(BtnX, false)
	}

	// Should end in false state
	if Btn(BtnX) {
		t.Fatalf("button should be false after rapid toggles")
	}
}

func TestMultipleButtonsSimultaneously(t *testing.T) {
	Step()

	// Set multiple buttons at once
	Set(BtnLeft, true)
	Set(BtnRight, true)
	Set(BtnUp, true)
	Set(BtnDown, true)
	Set(BtnO, true)
	Set(BtnX, true)

	// All should be true
	if !Btn(BtnLeft) || !Btn(BtnRight) || !Btn(BtnUp) || !Btn(BtnDown) || !Btn(BtnO) || !Btn(BtnX) {
		t.Fatalf("all buttons should be true when set simultaneously")
	}

	// Clear all
	Set(BtnLeft, false)
	Set(BtnRight, false)
	Set(BtnUp, false)
	Set(BtnDown, false)
	Set(BtnO, false)
	Set(BtnX, false)

	// All should be false
	if Btn(BtnLeft) || Btn(BtnRight) || Btn(BtnUp) || Btn(BtnDown) || Btn(BtnO) || Btn(BtnX) {
		t.Fatalf("all buttons should be false after clearing")
	}
}

func TestBtnpEdgeCases(t *testing.T) {
	Step()

	// Test btnp with invalid indices
	if Btnp(-1) {
		t.Fatalf("btnp(-1) should return false")
	}
	if Btnp(num) {
		t.Fatalf("btnp(num) should return false")
	}
	if Btnp(999) {
		t.Fatalf("btnp(999) should return false")
	}

	// Test btnp state machine edge cases
	Set(BtnUp, true)
	Step()
	Set(BtnUp, false)
	Step()
	// Now btnp should be false because button is released
	if Btnp(BtnUp) {
		t.Fatalf("btnp should be false when button is released")
	}

	// Set again - should trigger btnp
	Set(BtnUp, true)
	if !Btnp(BtnUp) {
		t.Fatalf("btnp should be true on press after release")
	}
}
