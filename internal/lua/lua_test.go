package lua

import "testing"

func TestLuaInitUpdate(t *testing.T) {
    vm := New()
    t.Cleanup(vm.Close)
    src := `
        called_init = false
        updates = 0
        function init()
            called_init = true
        end
        function update(dt)
            updates = updates + 1
        end
    `
    if err := vm.LoadString(src); err != nil { t.Fatalf("load: %v", err) }
    if err := vm.CallInit(); err != nil { t.Fatalf("init: %v", err) }
    if err := vm.CallUpdate(1/60); err != nil { t.Fatalf("update: %v", err) }
    if err := vm.CallUpdate(1/60); err != nil { t.Fatalf("update: %v", err) }

    if vm.L.GetGlobal("called_init").String() != "true" {
        t.Fatalf("init not called")
    }
    if vm.L.GetGlobal("updates").String() != "2" {
        t.Fatalf("updates not incremented")
    }
}


