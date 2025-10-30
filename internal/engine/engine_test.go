package engine

import (
    "testing"
    "time"
)

// fakeClock to avoid real sleeping
type fakeClock struct{ now time.Time }
func (f *fakeClock) Now() time.Time { return f.now }
func (f *fakeClock) Sleep(d time.Duration) { f.now = f.now.Add(d) }

func TestEngineRunsLuaUpdateAcrossFrames(t *testing.T) {
    e := New(60)
    t.Cleanup(e.Close)
    // Inject fake clock
    e.Sched.WithClock(&fakeClock{now: time.Unix(0, 0)})

    src := `
        updates = 0
        function init() end
        function update(dt)
            updates = updates + 1
        end
    `
    if err := e.LoadLuaSource(src); err != nil { t.Fatalf("load: %v", err) }
    e.RunFrames(3)
    if e.VM.L.GetGlobal("updates").String() != "3" {
        t.Fatalf("expected 3 updates")
    }
}

func TestEngineRunFramesRespectsFPS(t *testing.T) {
    e := New(30)
    t.Cleanup(e.Close)
    fc := &fakeClock{now: time.Unix(0, 0)}
    e.Sched.WithClock(fc)
    if err := e.LoadLuaSource(`function update(dt) end`); err != nil { t.Fatal(err) }
    e.RunFrames(2)
    // At 30 FPS, each frame is ~33.333ms; 2 frames advance ~66.666ms
    // We assert time advanced, not the exact float rounding.
    if fc.now.Sub(time.Unix(0, 0)) <= (time.Second/30) {
        t.Fatalf("clock didn't advance as expected")
    }
}


