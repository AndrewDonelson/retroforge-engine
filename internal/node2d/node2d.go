package node2d

import "math"

type Vec2 struct{ X, Y float64 }

type Transform struct{
    Position Vec2
    Scale    Vec2
    Rotation float64 // radians
}

func Identity() Transform { return Transform{Scale: Vec2{1,1}} }

func (t Transform) Apply(p Vec2) Vec2 {
    // scale
    x, y := p.X * t.Scale.X, p.Y * t.Scale.Y
    // rotate
    c, s := math.Cos(t.Rotation), math.Sin(t.Rotation)
    rx, ry := x*c - y*s, x*s + y*c
    // translate
    return Vec2{ rx + t.Position.X, ry + t.Position.Y }
}

// Combine composes parent followed by child into a world transform.
func Combine(parent, child Transform) Transform {
    // naive composition: scale, then rotate, then translate
    // scale composition
    out := Identity()
    out.Scale = Vec2{ parent.Scale.X * child.Scale.X, parent.Scale.Y * child.Scale.Y }
    // rotation composition
    out.Rotation = parent.Rotation + child.Rotation
    // position: apply parent to child's position
    out.Position = parent.Apply(child.Position)
    return out
}


