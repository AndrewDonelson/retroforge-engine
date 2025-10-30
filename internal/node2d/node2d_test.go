package node2d

import (
    "math"
    "testing"
)

func almost(a, b float64) bool { return math.Abs(a-b) < 1e-6 }

func TestTransformApply(t *testing.T) {
    t1 := Identity()
    t1.Position = Vec2{10, 5}
    t1.Scale = Vec2{2, 3}
    t1.Rotation = math.Pi / 2 // 90 deg

    out := t1.Apply(Vec2{1, 0})
    // After scale: (2,0), rotate 90deg: (0,2), translate: (10,7)
    if !almost(out.X, 10) || !almost(out.Y, 7) {
        t.Fatalf("unexpected: %#v", out)
    }
}

func TestCombine(t *testing.T) {
    parent := Identity(); parent.Position = Vec2{10,0}; parent.Scale = Vec2{2,2}
    child := Identity(); child.Position = Vec2{1,0}; child.Rotation = math.Pi/2
    world := Combine(parent, child)
    // world apply to origin should equal composed translation (parent acts on child position)
    o := world.Apply(Vec2{0,0})
    if !almost(o.X, 12) || !almost(o.Y, 0) { t.Fatalf("unexpected world origin: %#v", o) }
    // child's local point (1,0): scale(2) -> (2,0), rotate 90 -> (0,2), translate -> (12,2)
    p := world.Apply(Vec2{1,0})
    if !almost(p.X, 12) || !almost(p.Y, 2) { t.Fatalf("unexpected combined apply: %#v", p) }
}


