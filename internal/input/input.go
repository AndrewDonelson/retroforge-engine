package input

// Simple frame-based input state with btn/btnp for 6 Pico-like buttons.

const (
    BtnLeft = 0
    BtnRight = 1
    BtnUp = 2
    BtnDown = 3
    BtnO = 4
    BtnX = 5
    num = 6
)

var cur [num]bool
var prev [num]bool

func Step() { prev = cur }
func Set(i int, down bool) { if i>=0 && i<num { cur[i] = down } }
func Btn(i int) bool { if i<0 || i>=num { return false }; return cur[i] }
func Btnp(i int) bool { if i<0 || i>=num { return false }; return cur[i] && !prev[i] }


