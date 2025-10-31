//go:build !js && !wasm

package sdlrun

import (
	"testing"

	"github.com/AndrewDonelson/retroforge-engine/internal/engine"
)

func TestSaveScreenshot(t *testing.T) {
	// Test that saveScreenshot doesn't crash
	e := engine.New(60)
	defer e.Close()

	// Create a simple render state
	e.RunFrames(1)

	// Call screenshot function (should not crash even if file system fails)
	saveScreenshot(e)

	// If we get here, function completed without panic
}

func TestSaveScreenshotEdgeCases(t *testing.T) {
	// Test with fresh engine (no frames run)
	e1 := engine.New(60)
	defer e1.Close()
	saveScreenshot(e1) // Should handle empty/uninitialized state

	// Test multiple screenshots in rapid succession
	e2 := engine.New(60)
	defer e2.Close()
	e2.RunFrames(1)
	for i := 0; i < 10; i++ {
		saveScreenshot(e2)
	}

	// Test with engine that has run many frames
	e3 := engine.New(60)
	defer e3.Close()
	e3.RunFrames(100)
	saveScreenshot(e3)

	// Test with different frame rates
	e4 := engine.New(30)
	defer e4.Close()
	e4.RunFrames(1)
	saveScreenshot(e4)

	e5 := engine.New(120)
	defer e5.Close()
	e5.RunFrames(1)
	saveScreenshot(e5)
}
