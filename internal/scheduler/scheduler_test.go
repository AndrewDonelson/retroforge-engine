package scheduler

import (
    "testing"
    "time"
)

type fakeClock struct { now time.Time; slept time.Duration }
func (f *fakeClock) Now() time.Time { return f.now }
func (f *fakeClock) Sleep(d time.Duration) { f.slept += d }

func TestStepSleepsToFrame(t *testing.T) {
    fc := &fakeClock{ now: time.Unix(0,0) }
    s := New(60).WithClock(fc)
    called := false
    s.Step(func(dt time.Duration){
        called = true
        // simulate work of 5ms
        fc.now = fc.now.Add(5 * time.Millisecond)
    })
    if !called { t.Fatalf("tick not called") }
    frame := time.Second / 60
    if fc.slept != frame - 5*time.Millisecond { t.Fatalf("expected sleep %v, got %v", frame-5*time.Millisecond, fc.slept) }
}


