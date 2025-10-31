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
	if err := vm.LoadString(src); err != nil {
		t.Fatalf("load: %v", err)
	}
	if err := vm.CallInit(); err != nil {
		t.Fatalf("init: %v", err)
	}
	if err := vm.CallUpdate(1 / 60); err != nil {
		t.Fatalf("update: %v", err)
	}
	if err := vm.CallUpdate(1 / 60); err != nil {
		t.Fatalf("update: %v", err)
	}

	if vm.L.GetGlobal("called_init").String() != "true" {
		t.Fatalf("init not called")
	}
	if vm.L.GetGlobal("updates").String() != "2" {
		t.Fatalf("updates not incremented")
	}
}

func TestCallDraw(t *testing.T) {
	v := New()
	defer v.Close()

	// Test CallDraw when _DRAW doesn't exist (should not error)
	err := v.CallDraw()
	if err != nil {
		t.Errorf("CallDraw when _DRAW doesn't exist should not error, got: %v", err)
	}

	// Test CallDraw with _DRAW function (smoke test - just verify it doesn't crash)
	src := `
		function _DRAW()
			-- Empty function to test that CallDraw can invoke it
		end
	`
	err = v.LoadString(src)
	if err != nil {
		t.Fatalf("LoadString failed: %v", err)
	}

	// Call _DRAW - should not error
	err = v.CallDraw()
	if err != nil {
		t.Errorf("CallDraw failed: %v", err)
	}

	// Test with error in _DRAW
	src2 := `
		function _DRAW()
			error("test error")
		end
	`
	v2 := New()
	defer v2.Close()
	err = v2.LoadString(src2)
	if err != nil {
		t.Fatalf("LoadString failed: %v", err)
	}

	err = v2.CallDraw()
	if err == nil {
		t.Error("CallDraw with error in _DRAW should return error")
	}
}
