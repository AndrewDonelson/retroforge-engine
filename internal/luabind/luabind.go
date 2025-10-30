package luabind

import (
    "image/color"
    "github.com/yuin/gopher-lua"
    "github.com/AndrewDonelson/retroforge-engine/internal/graphics"
    "github.com/AndrewDonelson/retroforge-engine/internal/input"
    "github.com/AndrewDonelson/retroforge-engine/internal/audio"
    "github.com/AndrewDonelson/retroforge-engine/internal/app"
)

// small POD color to avoid importing image/color in caller signature
type ColorByIndex func(i int) (rgba [4]uint8)

// Register attaches rf.* drawing functions to the Lua state.
func Register(L *lua.LState, r graphics.Renderer, colorByIndex ColorByIndex, setPalette func(string)) {
    rf := L.NewTable()
    L.SetGlobal("rf", rf)

    // rf.clear(r,g,b)
    L.SetField(rf, "clear", L.NewFunction(func(L *lua.LState) int {
        r8 := uint8(L.CheckInt(1))
        g8 := uint8(L.CheckInt(2))
        b8 := uint8(L.CheckInt(3))
        r.Clear(color.RGBA{R:r8, G:g8, B:b8, A:0xFF})
        return 0
    }))

    // rf.print_center(text, y, r,g,b)
    L.SetField(rf, "print_center", L.NewFunction(func(L *lua.LState) int {
        txt := L.CheckString(1)
        y := L.CheckInt(2)
        r8 := uint8(L.CheckInt(3))
        g8 := uint8(L.CheckInt(4))
        b8 := uint8(L.CheckInt(5))
        r.PrintCentered(txt, y, color.RGBA{R:r8, G:g8, B:b8, A:0xFF})
        return 0
    }))

    // rf.clear_i(idx)
    L.SetField(rf, "clear_i", L.NewFunction(func(L *lua.LState) int {
        idx := L.CheckInt(1)
        c := colorByIndex(idx)
        r.Clear(color.RGBA{R:c[0], G:c[1], B:c[2], A:c[3]})
        return 0
    }))

    // rf.print_xy(x,y,text, idx)
    L.SetField(rf, "print_xy", L.NewFunction(func(L *lua.LState) int {
        x := L.CheckInt(1)
        y := L.CheckInt(2)
        txt := L.CheckString(3)
        idx := L.CheckInt(4)
        c := colorByIndex(idx)
        r.Print(txt, x, y, color.RGBA{R:c[0], G:c[1], B:c[2], A:c[3]})
        return 0
    }))

    // palette.set(name)
    L.SetField(rf, "palette_set", L.NewFunction(func(L *lua.LState) int {
        name := L.CheckString(1)
        if setPalette != nil { setPalette(name) }
        return 0
    }))

    // Drawing primitives (index-colored)
    L.SetField(rf, "pset", L.NewFunction(func(L *lua.LState) int {
        x := L.CheckInt(1); y := L.CheckInt(2); idx := L.CheckInt(3)
        c := colorByIndex(idx)
        r.PSet(x,y,color.RGBA{c[0],c[1],c[2],c[3]}); return 0
    }))
    L.SetField(rf, "line", L.NewFunction(func(L *lua.LState) int {
        x0 := L.CheckInt(1); y0 := L.CheckInt(2); x1 := L.CheckInt(3); y1 := L.CheckInt(4); idx := L.CheckInt(5)
        c := colorByIndex(idx)
        r.Line(x0,y0,x1,y1,color.RGBA{c[0],c[1],c[2],c[3]}); return 0
    }))
    L.SetField(rf, "rect", L.NewFunction(func(L *lua.LState) int {
        x0 := L.CheckInt(1); y0 := L.CheckInt(2); x1 := L.CheckInt(3); y1 := L.CheckInt(4); idx := L.CheckInt(5)
        c := colorByIndex(idx)
        r.Rect(x0,y0,x1,y1,color.RGBA{c[0],c[1],c[2],c[3]}); return 0
    }))
    L.SetField(rf, "rectfill", L.NewFunction(func(L *lua.LState) int {
        x0 := L.CheckInt(1); y0 := L.CheckInt(2); x1 := L.CheckInt(3); y1 := L.CheckInt(4); idx := L.CheckInt(5)
        c := colorByIndex(idx)
        r.RectFill(x0,y0,x1,y1,color.RGBA{c[0],c[1],c[2],c[3]}); return 0
    }))
    L.SetField(rf, "circ", L.NewFunction(func(L *lua.LState) int {
        x := L.CheckInt(1); y := L.CheckInt(2); rad := L.CheckInt(3); idx := L.CheckInt(4)
        c := colorByIndex(idx)
        r.Circ(x,y,rad,color.RGBA{c[0],c[1],c[2],c[3]}); return 0
    }))
    L.SetField(rf, "circfill", L.NewFunction(func(L *lua.LState) int {
        x := L.CheckInt(1); y := L.CheckInt(2); rad := L.CheckInt(3); idx := L.CheckInt(4)
        c := colorByIndex(idx)
        r.CircFill(x,y,rad,color.RGBA{c[0],c[1],c[2],c[3]}); return 0
    }))

    // RGB variants
    L.SetField(rf, "circfill_rgb", L.NewFunction(func(L *lua.LState) int {
        x := L.CheckInt(1); y := L.CheckInt(2); rad := L.CheckInt(3)
        rr := uint8(L.CheckInt(4)); gg := uint8(L.CheckInt(5)); bb := uint8(L.CheckInt(6))
        r.CircFill(x,y,rad,color.RGBA{rr,gg,bb,0xFF}); return 0
    }))

    // Input
    L.SetField(rf, "btn", L.NewFunction(func(L *lua.LState) int {
        i := L.CheckInt(1); L.Push(lua.LBool(input.Btn(i))); return 1
    }))
    L.SetField(rf, "btnp", L.NewFunction(func(L *lua.LState) int {
        i := L.CheckInt(1); L.Push(lua.LBool(input.Btnp(i))); return 1
    }))

    // Sound effects
    L.SetField(rf, "sfx", L.NewFunction(func(L *lua.LState) int {
        name := L.CheckString(1)
        action := L.OptString(2, "")
        _ = audio.Init()
        switch name {
        case "thrust":
            audio.Thrust(action != "off")
        case "land":
            audio.PlaySine(880, 0.12, 0.3)
        case "crash":
            audio.PlayNoise(0.25, 0.4)
        case "move":
            audio.PlaySine(520, 0.05, 0.25)
        case "select":
            audio.PlaySine(700, 0.08, 0.3)
        case "stopall":
            audio.StopAll()
        }
        return 0
    }))

    // Raw tone/noise
    L.SetField(rf, "tone", L.NewFunction(func(L *lua.LState) int {
        _ = audio.Init()
        f := L.CheckNumber(1)
        d := L.CheckNumber(2)
        g := L.OptNumber(3, 0.3)
        audio.PlaySine(float64(f), float64(d), float64(g))
        return 0
    }))
    L.SetField(rf, "noise", L.NewFunction(func(L *lua.LState) int {
        _ = audio.Init()
        d := L.CheckNumber(1)
        g := L.OptNumber(2, 0.3)
        audio.PlayNoise(float64(d), float64(g))
        return 0
    }))

    // Music: rf.music({"1G#2","R1","A3"}, bpm, gain)
    L.SetField(rf, "music", L.NewFunction(func(L *lua.LState) int {
        _ = audio.Init()
        tbl := L.CheckTable(1)
        bpm := L.OptNumber(2, 120)
        gain := L.OptNumber(3, 0.3)
        var toks []string
        tbl.ForEach(func(k, v lua.LValue) {
            if s, ok := v.(lua.LString); ok { toks = append(toks, string(s)) }
        })
        audio.PlayNotes(toks, float64(bpm), float64(gain))
        return 0
    }))

    // Quit request
    L.SetField(rf, "quit", L.NewFunction(func(L *lua.LState) int {
        app.RequestQuit(); return 0
    }))
}


