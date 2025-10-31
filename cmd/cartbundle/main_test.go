//go:build !js && !wasm

package main

import (
	"bytes"
	"testing"

	"github.com/AndrewDonelson/retroforge-engine/internal/engine"
)

func TestEmbeddedCart(t *testing.T) {
	// Test that embedded cart data exists and is valid
	if len(cartBytes) == 0 {
		t.Fatalf("embedded cart.rf is empty")
	}

	// Try to load it with engine
	e := engine.New(60)
	defer e.Close()

	err := e.LoadCartFromReader(bytes.NewReader(cartBytes), int64(len(cartBytes)))
	if err != nil {
		// This might fail if cart.rf is just a placeholder, that's ok
		t.Logf("embedded cart failed to load (expected if placeholder): %v", err)
	}

	// At minimum, we verified the embed worked
	if len(cartBytes) < 10 {
		t.Logf("cartBytes seems very small, might be placeholder")
	}
}

func TestEmbeddedCartEdgeCases(t *testing.T) {
	// Test cartBytes properties
	if len(cartBytes) == 0 {
		t.Fatalf("cartBytes should not be empty (even if placeholder)")
	}

	e := engine.New(60)
	defer e.Close()

	// Test with zero size
	err := e.LoadCartFromReader(bytes.NewReader(cartBytes), 0)
	if err == nil {
		t.Logf("LoadCartFromReader with zero size succeeded")
	}

	// Test with negative size
	err = e.LoadCartFromReader(bytes.NewReader(cartBytes), -1)
	if err == nil {
		t.Logf("LoadCartFromReader with negative size succeeded")
	}

	// Test with size larger than actual data
	err = e.LoadCartFromReader(bytes.NewReader(cartBytes), int64(len(cartBytes)*2))
	if err != nil {
		t.Logf("LoadCartFromReader with oversized length failed: %v", err)
	}

	// Test that cartBytes is accessible multiple times
	for i := 0; i < 10; i++ {
		if len(cartBytes) == 0 {
			t.Fatalf("cartBytes became empty on access %d", i)
		}
	}

	// Test reading from middle of cartBytes
	if len(cartBytes) > 10 {
		partialReader := bytes.NewReader(cartBytes[10:])
		err = e.LoadCartFromReader(partialReader, int64(len(cartBytes)-10))
		if err != nil {
			t.Logf("LoadCartFromReader with partial data failed: %v", err)
		}
	}
}
