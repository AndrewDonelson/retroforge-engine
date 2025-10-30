package runner

import (
    "time"
    "github.com/AndrewDonelson/retroforge-engine/internal/eventbus"
    "github.com/AndrewDonelson/retroforge-engine/internal/scheduler"
)

// Runner ties the scheduler to the event bus by publishing tick events.
type Runner struct {
    Bus   *eventbus.Bus
    Sched *scheduler.Scheduler
}

func New(bus *eventbus.Bus, sched *scheduler.Scheduler) *Runner { return &Runner{Bus: bus, Sched: sched} }

// Step runs one frame and publishes "tick" with dt.
func (r *Runner) Step() {
    r.Sched.Step(func(dt time.Duration){ r.Bus.Publish("tick", dt) })
}


