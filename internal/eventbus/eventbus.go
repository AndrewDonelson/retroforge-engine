package eventbus

import "sync"

// Bus is a lightweight pub/sub.
type Bus struct {
    mu   sync.RWMutex
    subs map[string][]func(any)
}

func New() *Bus { return &Bus{subs: make(map[string][]func(any))} }

func (b *Bus) Subscribe(topic string, fn func(any)) (unsubscribe func()) {
    b.mu.Lock()
    b.subs[topic] = append(b.subs[topic], fn)
    idx := len(b.subs[topic]) - 1
    b.mu.Unlock()
    return func() {
        b.mu.Lock()
        defer b.mu.Unlock()
        list := b.subs[topic]
        if idx >= 0 && idx < len(list) {
            // remove by index captured at subscribe time
            copy(list[idx:], list[idx+1:])
            b.subs[topic] = list[:len(list)-1]
        }
    }
}

func (b *Bus) Publish(topic string, payload any) {
    b.mu.RLock()
    cbs := append([]func(any){}, b.subs[topic]...)
    b.mu.RUnlock()
    for _, cb := range cbs {
        cb(payload)
    }
}


