package game

import (
    "sync"
    "github.com/AndrewDonelson/retroforge-engine/internal/eventbus"
)

// HelloWorld publishes a single log line on first tick.
type HelloWorld struct { once sync.Once }

func (g *HelloWorld) Init(bus *eventbus.Bus) {
    bus.Subscribe("tick", func(any) {
        g.once.Do(func() { bus.Publish("log", "Hello, RetroForge!") })
    })
}


