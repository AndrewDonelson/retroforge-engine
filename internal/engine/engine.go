package engine

import (
    "io"
    "os"
    "time"
    "github.com/AndrewDonelson/retroforge-engine/internal/cartio"
    "github.com/AndrewDonelson/retroforge-engine/internal/eventbus"
    "github.com/AndrewDonelson/retroforge-engine/internal/graphics"
    "github.com/AndrewDonelson/retroforge-engine/internal/luabind"
    "github.com/AndrewDonelson/retroforge-engine/internal/lua"
    "github.com/AndrewDonelson/retroforge-engine/internal/pal"
    "github.com/AndrewDonelson/retroforge-engine/internal/rendersoft"
    "github.com/AndrewDonelson/retroforge-engine/internal/runner"
    "github.com/AndrewDonelson/retroforge-engine/internal/scheduler"
)

// Engine wires together bus, scheduler/runner, and Lua VM for headless runs.
type Engine struct {
    Bus   *eventbus.Bus
    Sched *scheduler.Scheduler
    Run   *runner.Runner
    VM    *lua.VM
    Ren   graphics.Renderer
    Pal   *pal.Manager
}

func New(targetFPS int) *Engine {
    bus := eventbus.New()
    sched := scheduler.New(targetFPS)
    run := runner.New(bus, sched)
    vm := lua.New()
    ren := rendersoft.New(480, 270)
    e := &Engine{Bus: bus, Sched: sched, Run: run, VM: vm, Ren: ren, Pal: pal.NewManager()}
    // On each tick, call Lua update with dt seconds.
    bus.Subscribe("tick", func(v any) {
        if dt, ok := v.(time.Duration); ok {
            _ = e.VM.CallUpdate(dt.Seconds())
            _ = e.VM.CallDraw()
        }
    })
    return e
}

func (e *Engine) Close() { e.VM.Close() }

// LoadLuaSource loads script and calls init() if present.
func (e *Engine) LoadLuaSource(src string) error {
    luabind.Register(e.VM.L, e.Ren, func(i int) (c [4]uint8) { col := e.Pal.Color(i); c[0]=col.R; c[1]=col.G; c[2]=col.B; c[3]=col.A; return }, e.Pal.Set)
    if err := e.VM.LoadString(src); err != nil { return err }
    return e.VM.CallInit()
}

// RunFrames advances N frames headlessly.
func (e *Engine) RunFrames(n int) {
    for i := 0; i < n; i++ { e.Run.Step() }
}

// LoadCartFromReader loads a .rfs from an io.ReaderAt.
func (e *Engine) LoadCartFromReader(r io.ReaderAt, size int64) error {
    m, files, err := cartio.Read(r, size)
    if err != nil { return err }
    _ = m
    src, ok := files["assets/"+m.Entry]
    if !ok { return os.ErrNotExist }
    return e.LoadLuaSource(string(src))
}

// LoadCartFile opens .rfs by path and loads it.
func (e *Engine) LoadCartFile(path string) error {
    f, err := os.Open(path)
    if err != nil { return err }
    defer f.Close()
    st, err := f.Stat()
    if err != nil { return err }
    return e.LoadCartFromReader(f, st.Size())
}


