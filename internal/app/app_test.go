package app

import "testing"

func TestQuitRequest(t *testing.T) {
	// Reset state
	Reset()

	if QuitRequested() {
		t.Fatalf("quit should not be requested initially")
	}

	RequestQuit()
	if !QuitRequested() {
		t.Fatalf("quit should be requested after RequestQuit()")
	}

	// Should still be requested
	if !QuitRequested() {
		t.Fatalf("quit should remain requested")
	}

	// Reset and verify
	Reset()
	if QuitRequested() {
		t.Fatalf("quit should not be requested after Reset()")
	}
}

func TestQuitThreadSafety(t *testing.T) {
	Reset()

	// Simulate concurrent access (basic test)
	done := make(chan bool)

	go func() {
		RequestQuit()
		done <- true
	}()

	go func() {
		_ = QuitRequested()
		done <- true
	}()

	<-done
	<-done

	// Should eventually be requested
	if !QuitRequested() {
		t.Fatalf("quit should be requested after concurrent RequestQuit")
	}
}

func TestQuitRapidToggle(t *testing.T) {
	// Test rapid toggling
	for i := 0; i < 100; i++ {
		RequestQuit()
		Reset()
	}

	// Should end in reset state
	if QuitRequested() {
		t.Fatalf("quit should not be requested after Reset")
	}
}

func TestQuitConcurrent(t *testing.T) {
	Reset()

	// Multiple concurrent requests
	done := make(chan bool, 100)
	for i := 0; i < 100; i++ {
		go func() {
			RequestQuit()
			done <- true
		}()
	}

	// Wait for all
	for i := 0; i < 100; i++ {
		<-done
	}

	// Should be requested
	if !QuitRequested() {
		t.Fatalf("quit should be requested after many concurrent requests")
	}

	Reset()

	// Concurrent reads
	readCount := 0
	for i := 0; i < 100; i++ {
		go func() {
			_ = QuitRequested()
			done <- true
		}()
	}
	for i := 0; i < 100; i++ {
		<-done
		readCount++
	}

	if readCount != 100 {
		t.Fatalf("all reads should complete")
	}
}
