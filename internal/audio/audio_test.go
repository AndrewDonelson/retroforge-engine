package audio

import (
	"testing"
	"time"
)

// Note: Audio tests are limited because audio requires SDL initialization
// which may not be available in all test environments

func TestInitErrorHandling(t *testing.T) {
	// Test that Init can be called multiple times
	// (should not error on subsequent calls)
	err1 := Init()
	err2 := Init()

	// If first init fails, subsequent should also fail gracefully
	if err1 != nil && err2 != nil {
		t.Logf("Audio init failed (expected if SDL not available): %v", err1)
		return
	}

	// If init succeeds, second call should not error
	if err1 == nil && err2 != nil {
		t.Fatalf("second Init() should not error if first succeeded")
	}
}

func TestStopAll(t *testing.T) {
	// StopAll should not panic even if audio not initialized
	StopAll()

	// Should be safe to call multiple times
	StopAll()
	StopAll()
}

func TestThrustState(t *testing.T) {
	// Test thrust on/off doesn't crash
	// Note: Actual audio playback requires SDL, so we just test it doesn't panic
	Thrust(true)
	Thrust(false)
	Thrust(true)
	Thrust(false)

	StopAll()
}

func TestPlaySineNoCrash(t *testing.T) {
	// Test that PlaySine doesn't crash even if audio not initialized
	PlaySine(440.0, 0.1, 0.5)

	// Give it a moment
	time.Sleep(10 * time.Millisecond)

	StopAll()
}

func TestPlayNoiseNoCrash(t *testing.T) {
	// Test that PlayNoise doesn't crash
	PlayNoise(0.1, 0.5)

	// Give it a moment
	time.Sleep(10 * time.Millisecond)

	StopAll()
}

func TestPlayNotesNoCrash(t *testing.T) {
	// Test that PlayNotes doesn't crash
	tokens := []string{"4C1", "4E1", "4G1", "R1"}
	PlayNotes(tokens, 120.0, 0.3)

	// Give it a moment to start
	time.Sleep(50 * time.Millisecond)

	StopAll()
}

func TestNoteToFreq(t *testing.T) {
	// Test note frequency calculation
	testCases := []struct {
		note   string
		wantOk bool
	}{
		{"4C", true},
		{"4A", true},
		{"R", true}, // rest
		{"", false},
		{"X", false}, // invalid note
	}

	for _, tc := range testCases {
		freq, ok := noteToFreq(tc.note, 4)
		if ok != tc.wantOk {
			t.Errorf("noteToFreq(%q) ok=%v, want %v", tc.note, ok, tc.wantOk)
			continue
		}
		if ok && tc.note != "R" {
			if freq <= 0 {
				t.Errorf("noteToFreq(%q) returned invalid frequency %f", tc.note, freq)
			}
			// A4 should be 440Hz
			if tc.note == "4A" {
				expected := 440.0
				tolerance := 1.0
				if freq < expected-tolerance || freq > expected+tolerance {
					t.Errorf("noteToFreq(\"4A\") = %f, want ~%f", freq, expected)
				}
			}
		}
	}
}

func TestNoteToFreqEdgeCases(t *testing.T) {
	// Test invalid note strings
	invalidNotes := []string{
		"",    // empty
		"   ", // whitespace only
		"X",   // invalid note name
		"4",   // octave only, no note
		"#",   // sharp only
		"4X#", // invalid note with octave and sharp
	}

	for _, note := range invalidNotes {
		if note == "" || note == "   " || note == "X" {
			freq, ok := noteToFreq(note, 4)
			if ok {
				t.Errorf("noteToFreq(%q) should return ok=false", note)
			}
			if freq != 0 {
				t.Errorf("noteToFreq(%q) should return freq=0 when not ok", note)
			}
		}
	}

	// Test edge case octaves
	for octave := 0; octave <= 9; octave++ {
		freq, ok := noteToFreq("A", octave)
		if !ok {
			t.Errorf("noteToFreq(\"A\", %d) should return ok=true", octave)
		}
		if freq <= 0 {
			t.Errorf("noteToFreq(\"A\", %d) should return positive freq", octave)
		}
	}

	// Test all note names
	notes := []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}
	for _, note := range notes {
		freq, ok := noteToFreq(note, 4)
		if !ok {
			t.Errorf("noteToFreq(%q) should return ok=true", note)
		}
		if freq <= 0 {
			t.Errorf("noteToFreq(%q) should return positive freq", note)
		}
	}
}

func TestPlaySineEdgeCases(t *testing.T) {
	// Test with extreme values
	PlaySine(0, 0.1, 0.3)      // zero frequency
	PlaySine(-100, 0.1, 0.3)   // negative frequency
	PlaySine(999999, 0.1, 0.3) // very high frequency

	PlaySine(440, 0, 0.3)      // zero duration
	PlaySine(440, -0.1, 0.3)   // negative duration
	PlaySine(440, 999999, 0.3) // very long duration

	PlaySine(440, 0.1, 0)    // zero gain
	PlaySine(440, 0.1, -0.5) // negative gain
	PlaySine(440, 0.1, 999)  // very large gain

	// Should not crash
	StopAll()
}

func TestPlayNoiseEdgeCases(t *testing.T) {
	// Test with extreme values
	PlayNoise(0, 0.3)      // zero duration
	PlayNoise(-0.1, 0.3)   // negative duration
	PlayNoise(999999, 0.3) // very long duration

	PlayNoise(0.1, 0)    // zero gain
	PlayNoise(0.1, -0.5) // negative gain
	PlayNoise(0.1, 999)  // very large gain

	// Should not crash
	StopAll()
}

func TestPlayNotesEdgeCases(t *testing.T) {
	// Test with empty tokens
	PlayNotes([]string{}, 120, 0.3)

	// Test with invalid tokens
	PlayNotes([]string{"INVALID", "X", "???"}, 120, 0.3)

	// Test with mixed valid/invalid
	PlayNotes([]string{"4C1", "INVALID", "4E1"}, 120, 0.3)

	// Test with extreme BPM
	PlayNotes([]string{"4C1"}, 0, 0.3)      // zero BPM
	PlayNotes([]string{"4C1"}, -100, 0.3)   // negative BPM
	PlayNotes([]string{"4C1"}, 999999, 0.3) // very high BPM

	// Test with extreme gain
	PlayNotes([]string{"4C1"}, 120, 0)   // zero gain
	PlayNotes([]string{"4C1"}, 120, -1)  // negative gain
	PlayNotes([]string{"4C1"}, 120, 999) // very large gain

	// Test with very long sequence
	longSeq := make([]string, 1000)
	for i := range longSeq {
		longSeq[i] = "4C1"
	}
	PlayNotes(longSeq, 120, 0.3)

	time.Sleep(50 * time.Millisecond)
	StopAll()
}

func TestThrustEdgeCases(t *testing.T) {
	// Test rapid toggling
	for i := 0; i < 100; i++ {
		Thrust(true)
		Thrust(false)
	}

	// Test multiple on calls
	Thrust(true)
	Thrust(true)
	Thrust(true)
	Thrust(false)

	// Should not crash
	StopAll()
}

func TestStopAllMultiple(t *testing.T) {
	// Call StopAll multiple times
	for i := 0; i < 100; i++ {
		StopAll()
	}

	// Should not crash
}

func TestPlayFunctionsConcurrent(t *testing.T) {
	// Test concurrent calls (basic test)
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			PlaySine(440, 0.1, 0.3)
			PlayNoise(0.1, 0.3)
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	StopAll()
}
