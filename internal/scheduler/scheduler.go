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

// Scheduler runs a variable-interval tick loop using actual elapsed time.
type Scheduler struct {
    TargetFPS int
    clock     Clock
    lastTime  time.Time // Last frame time for delta calculation
    initialized bool    // Whether we've initialized lastTime
}

func New(targetFPS int) *Scheduler { 
    return &Scheduler{
        TargetFPS: targetFPS, 
        clock: realClock{},
        lastTime: time.Time{},
        initialized: false,
    } 
}
func (s *Scheduler) WithClock(c Clock) *Scheduler { s.clock = c; return s }

func (s *Scheduler) Step(fn func(dt time.Duration)) {
    if s.TargetFPS <= 0 { s.TargetFPS = 60 }
    
    now := s.clock.Now()
    
    // Calculate actual delta time from last frame
    var dt time.Duration
    if !s.initialized {
        // First frame - use target frame time (1/60s for 60 FPS)
        dt = time.Second / time.Duration(s.TargetFPS)
        s.initialized = true
    } else {
        // Subsequent frames - use actual elapsed time
        dt = now.Sub(s.lastTime)
        
        // Cap delta time to prevent huge jumps (e.g., if tab was inactive, or system lag)
        // Cap at 10x target frame time (e.g., 166ms for 60 FPS = ~6 FPS minimum)
        maxDt := (time.Second / time.Duration(s.TargetFPS)) * 10
        if dt > maxDt {
            dt = maxDt
        }
        
        // Ensure minimum delta time (prevents division by zero and handles very fast frames)
        minDt := time.Millisecond // 1ms minimum
        if dt < minDt {
            dt = minDt
        }
    }
    
    // Update last frame time
    s.lastTime = now
    
    // Call the frame function with actual delta time
    start := s.clock.Now()
    fn(dt)
    elapsed := s.clock.Now().Sub(start)
    
    // Sleep to maintain target FPS (only on desktop, WASM skips this)
    frame := time.Second / time.Duration(s.TargetFPS)
    if remain := frame - elapsed; remain > 0 {
        s.clock.Sleep(remain)
    }
}


