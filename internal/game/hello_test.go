package game

import (
    "testing"
    "time"
    "github.com/AndrewDonelson/retroforge-engine/internal/eventbus"
    "github.com/AndrewDonelson/retroforge-engine/internal/runner"
    "github.com/AndrewDonelson/retroforge-engine/internal/scheduler"
)

func TestHelloWorldPublishesLogOnce(t *testing.T) {
    bus := eventbus.New()
    sched := scheduler.New(60).WithClock(fakeClock{now: time.Unix(0, 0)})
    r := runner.New(bus, sched)

    hw := &HelloWorld{}
    hw.Init(bus)

    got := make(chan string, 2)
    bus.Subscribe("log", func(v any) {
        if s, ok := v.(string); ok { got <- s }
    })

    // First tick should emit
    r.Step()
    select {
    case msg := <-got:
        if msg != "Hello, RetroForge!" { t.Fatalf("unexpected log: %q", msg) }
    case <-time.After(50 * time.Millisecond):
        t.Fatal("expected hello log on first tick")
    }

    // Subsequent ticks should not emit again
    r.Step(); r.Step()
    select {
    case msg := <-got:
        t.Fatalf("expected only one hello, got: %q", msg)
    default:
        // ok
    }
}

// fakeClock is a minimal Clock for deterministic Scheduler.Step timing.
type fakeClock struct{ now time.Time }

func (f fakeClock) Now() time.Time { return f.now }
func (f fakeClock) Sleep(d time.Duration) { f.now = f.now.Add(d) }


