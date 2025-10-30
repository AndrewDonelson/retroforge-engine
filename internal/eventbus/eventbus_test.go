package eventbus

import "testing"

func TestPublishSubscribe(t *testing.T) {
    b := New()
    got := 0
    unsub := b.Subscribe("tick", func(v any) {
        if n, ok := v.(int); ok { got += n }
    })
    b.Publish("tick", 1)
    b.Publish("tick", 2)
    if got != 3 { t.Fatalf("expected 3, got %d", got) }
    unsub()
    b.Publish("tick", 10)
    if got != 3 { t.Fatalf("unsubscribe failed") }
}


