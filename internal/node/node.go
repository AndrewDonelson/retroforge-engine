package node

// Node is a lightweight scene-graph node for RetroForge.
// It only includes logic needed for tests and headless behaviors.
type Node struct {
    Name     string
    Parent   *Node
    Children []*Node
    Enabled  bool
}

func New(name string) *Node { return &Node{Name: name, Enabled: true} }

func (n *Node) AddChild(child *Node) {
    if child == nil || child == n { return }
    if child.Parent != nil { child.Parent.RemoveChild(child) }
    child.Parent = n
    n.Children = append(n.Children, child)
}

func (n *Node) RemoveChild(child *Node) bool {
    for i, c := range n.Children {
        if c == child {
            copy(n.Children[i:], n.Children[i+1:])
            n.Children = n.Children[:len(n.Children)-1]
            c.Parent = nil
            return true
        }
    }
    return false
}

func (n *Node) FindByName(name string) *Node {
    if n.Name == name { return n }
    for _, c := range n.Children {
        if r := c.FindByName(name); r != nil { return r }
    }
    return nil
}


