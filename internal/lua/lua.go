package lua

import (
    "github.com/yuin/gopher-lua"
)

// VM is a thin wrapper around gopher-lua exposing init/update hooks.
type VM struct {
    L *lua.LState
}

func New() *VM { return &VM{L: lua.NewState()} }

func (v *VM) Close() { if v.L != nil { v.L.Close(); v.L = nil } }

// LoadString loads a Lua source string.
func (v *VM) LoadString(src string) error {
    return v.L.DoString(src)
}

// CallInit calls global function init() if present.
func (v *VM) CallInit() error {
    if err := v.callIfExists("_INIT", 0); err != nil { return err }
    return v.callIfExists("init", 0)
}

// CallUpdate calls global function update(dt) if present.
func (v *VM) CallUpdate(dtSeconds float64) error {
    if v.L.GetGlobal("_UPDATE") != lua.LNil {
        v.L.Push(v.L.GetGlobal("_UPDATE"))
        v.L.Push(lua.LNumber(dtSeconds))
        return v.L.PCall(1, 0, nil)
    }
    if v.L.GetGlobal("update") != lua.LNil {
        v.L.Push(v.L.GetGlobal("update"))
        v.L.Push(lua.LNumber(dtSeconds))
        return v.L.PCall(1, 0, nil)
    }
    return nil
}

func (v *VM) CallDraw() error { return v.callIfExists("_DRAW", 0) }

func (v *VM) callIfExists(name string, narg int) error {
    if v.L.GetGlobal(name) == lua.LNil { return nil }
    v.L.Push(v.L.GetGlobal(name))
    return v.L.PCall(narg, 0, nil)
}


