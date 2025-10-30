package node

import "testing"

func TestHierarchyAddRemove(t *testing.T) {
    root := New("root")
    a := New("a")
    b := New("b")
    root.AddChild(a)
    a.AddChild(b)

    if b.Parent != a { t.Fatalf("expected parent a") }
    if root.FindByName("b") == nil { t.Fatalf("expected to find b") }

    a.RemoveChild(b)
    if b.Parent != nil { t.Fatalf("expected b to be detached") }
    if root.FindByName("b") != nil { t.Fatalf("b should not be found after removal") }
}


