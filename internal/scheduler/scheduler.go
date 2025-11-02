package scheduler

import (
	"runtime"
	"time"
)

// Clock abstracts time for tests.
type Clock interface { Now() time.Time; Sleep(d time.Duration) }

type realClock struct{}
func (realClock) Now() time.Time { return time.Now() }
func (realClock) Sleep(d time.Duration) { 
	// In WASM (GOOS=js), time.Sleep can block the event loop
	// Skip sleeping in WASM - the browser's requestAnimationFrame handles frame timing
	if runtime.GOOS != "js" {
		time.Sleep(d)
	}
	// In WASM, we rely on the browser's animation frame timing instead
}

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


