package runner

import (
    "testing"
    "time"
    "github.com/AndrewDonelson/retroforge-engine/internal/eventbus"
    "github.com/AndrewDonelson/retroforge-engine/internal/scheduler"
)

type fakeClock struct { now time.Time }
func (f *fakeClock) Now() time.Time { return f.now }
func (f *fakeClock) Sleep(d time.Duration) { f.now = f.now.Add(d) }

func TestRunnerPublishesTick(t *testing.T) {
    bus := eventbus.New()
    fc := &fakeClock{ now: time.Unix(0,0) }
    sched := scheduler.New(60).WithClock(fc)
    r := New(bus, sched)
    got := 0
    bus.Subscribe("tick", func(v any){ if _, ok := v.(time.Duration); ok { got++ } })
    r.Step()
    if got != 1 { t.Fatalf("expected 1 tick, got %d", got) }
}


