package scheduler

import "time"

// Clock abstracts time for tests.
type Clock interface { Now() time.Time; Sleep(d time.Duration) }

type realClock struct{}
func (realClock) Now() time.Time { return time.Now() }
func (realClock) Sleep(d time.Duration) { time.Sleep(d) }

// Scheduler runs a fixed-interval tick loop.
type Scheduler struct {
    TargetFPS int
    clock     Clock
}

func New(targetFPS int) *Scheduler { return &Scheduler{TargetFPS: targetFPS, clock: realClock{}} }
func (s *Scheduler) WithClock(c Clock) *Scheduler { s.clock = c; return s }

func (s *Scheduler) Step(fn func(dt time.Duration)) {
    if s.TargetFPS <= 0 { s.TargetFPS = 60 }
    frame := time.Second / time.Duration(s.TargetFPS)
    start := s.clock.Now()
    fn(frame)
    elapsed := s.clock.Now().Sub(start)
    if remain := frame - elapsed; remain > 0 { s.clock.Sleep(remain) }
}


