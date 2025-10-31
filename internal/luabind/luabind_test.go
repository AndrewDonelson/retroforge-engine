package luabind

import (
	"testing"

	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
	"github.com/AndrewDonelson/retroforge-engine/internal/rendersoft"
	lua "github.com/yuin/gopher-lua"
)

func TestRegister(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(100, 100)
	colorByIndex := func(i int) (rgba [4]uint8) {
		if i == 0 {
			return [4]uint8{0, 0, 0, 255}
		}
		if i == 1 {
			return [4]uint8{255, 255, 255, 255}
		}
		return [4]uint8{128, 128, 128, 255}
	}
	setPalette := func(name string) {
		_ = name // not implemented yet
	}

	Register(L, r, colorByIndex, setPalette, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil)

	// Verify rf table exists
	rf := L.GetGlobal("rf")
	if rf == lua.LNil {
		t.Fatalf("rf global should exist")
	}
	if rf.Type() != lua.LTTable {
		t.Fatalf("rf should be a table")
	}

	// Verify some functions exist
	fields := []string{"clear_i", "print", "print_xy", "print_anchored", "cursor", "color",
		"btn", "btnp", "sfx", "tone", "noise", "music", "cstore", "reload",
		"flr", "ceil", "rnd", "mid", "sgn", "chr", "ord",
		"shl", "shr", "band", "bor", "bxor", "bnot", "quit"}
	for _, field := range fields {
		val := L.GetField(rf, field)
		if val == lua.LNil {
			t.Fatalf("rf.%s should exist", field)
		}
		if val.Type() != lua.LTFunction {
			t.Fatalf("rf.%s should be a function", field)
		}
	}
}

func TestClearFunction(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(10, 10)
	Register(L, r, func(i int) (rgba [4]uint8) {
		if i == 1 {
			return [4]uint8{255, 128, 64, 255}
		}
		return [4]uint8{0, 0, 0, 255}
	}, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil)

	// Call rf.clear_i(1) - uses palette index
	err := L.DoString(`rf.clear_i(1)`)
	if err != nil {
		t.Fatalf("rf.clear_i should not error: %v", err)
	}

	// Verify pixels were cleared with palette color
	pix := r.Pixels()
	if pix[0] != 255 || pix[1] != 128 || pix[2] != 64 {
		t.Fatalf("clear_i should set pixels with palette color")
	}
}

func TestClearIFunction(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(10, 10)
	called := false
	colorByIndex := func(i int) (rgba [4]uint8) {
		called = true
		if i == 1 {
			return [4]uint8{255, 255, 255, 255}
		}
		return [4]uint8{0, 0, 0, 255}
	}
	Register(L, r, colorByIndex, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil)

	// Call rf.clear_i(1)
	err := L.DoString(`rf.clear_i(1)`)
	if err != nil {
		t.Fatalf("rf.clear_i should not error: %v", err)
	}

	if !called {
		t.Fatalf("colorByIndex should be called")
	}

	// Verify pixels were cleared with white
	pix := r.Pixels()
	if pix[0] != 255 || pix[1] != 255 || pix[2] != 255 {
		t.Fatalf("clear_i should use palette color")
	}
}

func TestInputFunctions(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(10, 10)
	Register(L, r, func(i int) (rgba [4]uint8) { return [4]uint8{0, 0, 0, 255} }, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil)

	// Test btn (should return boolean)
	err := L.DoString(`local result = rf.btn(0); if type(result) ~= "boolean" then error("btn should return boolean") end`)
	if err != nil {
		t.Fatalf("rf.btn should work: %v", err)
	}

	// Test btnp (should return boolean)
	err = L.DoString(`local result = rf.btnp(1); if type(result) ~= "boolean" then error("btnp should return boolean") end`)
	if err != nil {
		t.Fatalf("rf.btnp should work: %v", err)
	}
}

func TestPrintFunctions(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(100, 100)
	Register(L, r, func(i int) (rgba [4]uint8) { return [4]uint8{0, 0, 0, 255} }, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil)

	// Test print with cursor/color state
	err := L.DoString(`rf.cursor(10, 20); rf.color(5); rf.print("HELLO")`)
	if err != nil {
		t.Fatalf("rf.print with cursor/color should not error: %v", err)
	}

	// Test print_xy
	err = L.DoString(`rf.print_xy(10, 20, "TEST", 1)`)
	if err != nil {
		t.Fatalf("rf.print_xy should not error: %v", err)
	}

	// Test print with explicit coordinates
	err = L.DoString(`rf.print("WORLD", 30, 40, 2)`)
	if err != nil {
		t.Fatalf("rf.print with explicit coords should not error: %v", err)
	}
}

func TestLuaInvalidParameters(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(100, 100)
	Register(L, r, func(i int) (rgba [4]uint8) { return [4]uint8{0, 0, 0, 255} }, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil)

	// Test clear_i with wrong number of args
	err := L.DoString(`rf.clear_i()`) // missing args
	if err == nil {
		t.Fatalf("rf.clear_i with insufficient args should error")
	}

	err = L.DoString(`rf.clear_i("invalid")`) // wrong type
	if err == nil {
		t.Fatalf("rf.clear_i with wrong type should error")
	}

	// Test cursor/color functions
	err = L.DoString(`rf.cursor(10, 20)`)
	if err != nil {
		t.Fatalf("rf.cursor should not error: %v", err)
	}

	err = L.DoString(`rf.color(5)`)
	if err != nil {
		t.Fatalf("rf.color should not error: %v", err)
	}

	// Test btn with wrong type
	err = L.DoString(`rf.btn("invalid")`)
	if err == nil {
		t.Fatalf("rf.btn with wrong type should error")
	}

	// Test invalid button index (negative, out of range)
	err = L.DoString(`rf.btn(-1)`) // should work but return false
	if err != nil {
		t.Fatalf("rf.btn(-1) should not error, just return false")
	}

	err = L.DoString(`rf.btn(999)`) // should work but return false
	if err != nil {
		t.Fatalf("rf.btn(999) should not error, just return false")
	}

	// Test clear_i with out of range index
	err = L.DoString(`rf.clear_i(999)`) // out of range
	// Note: Lua CheckInt allows any int, colorByIndex handles it
	if err != nil {
		t.Logf("rf.clear_i with out of range index failed: %v (may be expected)", err)
	}
}

func TestLuaEdgeCases(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(100, 100)
	Register(L, r, func(i int) (rgba [4]uint8) {
		if i < 0 || i >= 50 {
			return [4]uint8{0, 0, 0, 255}
		}
		return [4]uint8{255, 255, 255, 255}
	}, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil)

	// Test with nil/empty values
	err := L.DoString(`rf.print_xy(0, 0, "", 0)`) // empty string
	if err != nil {
		t.Fatalf("rf.print_xy with empty string should work")
	}

	// Test with very large coordinates
	err = L.DoString(`rf.print_xy(999999, 999999, "TEST", 0)`)
	if err != nil {
		t.Fatalf("rf.print_xy with large coords should work (will clip)")
	}

	// Test with negative coordinates
	err = L.DoString(`rf.print_xy(-100, -100, "TEST", 0)`)
	if err != nil {
		t.Fatalf("rf.print_xy with negative coords should work (will clip)")
	}

	// Test clear_i with invalid index
	err = L.DoString(`rf.clear_i(-1)`)
	if err != nil {
		t.Fatalf("rf.clear_i(-1) should work (colorByIndex handles it)")
	}

	err = L.DoString(`rf.clear_i(999)`)
	if err != nil {
		t.Fatalf("rf.clear_i(999) should work (colorByIndex handles it)")
	}

	// Test drawing primitives with edge cases
	err = L.DoString(`rf.pset(-100, -100, 0)`)
	if err != nil {
		t.Fatalf("rf.pset with negative coords should work")
	}

	err = L.DoString(`rf.line(0, 0, 1000, 1000, 0)`)
	if err != nil {
		t.Fatalf("rf.line with large coords should work")
	}

	err = L.DoString(`rf.rect(-10, -10, 1000, 1000, 0)`)
	if err != nil {
		t.Fatalf("rf.rect with out of bounds should work")
	}

	err = L.DoString(`rf.circ(50, 50, -10, 0)`) // negative radius
	if err != nil {
		t.Fatalf("rf.circ with negative radius should work")
	}

	err = L.DoString(`rf.circ(50, 50, 999999, 0)`) // huge radius
	if err != nil {
		t.Fatalf("rf.circ with huge radius should work")
	}
}

func TestLuaInvalidFunctions(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(100, 100)
	Register(L, r, func(i int) (rgba [4]uint8) { return [4]uint8{0, 0, 0, 255} }, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil)

	// Test calling non-existent function
	err := L.DoString(`rf.nonexistent()`)
	if err == nil {
		t.Fatalf("calling non-existent function should error")
	}

	// Test calling with nil (should error in Lua)
	err = L.DoString(`local nilval = nil; rf.clear(nilval, 128, 64)`)
	if err == nil {
		t.Fatalf("calling with nil should error")
	}
}

func TestLuaMusicEdgeCases(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(100, 100)
	Register(L, r, func(i int) (rgba [4]uint8) { return [4]uint8{0, 0, 0, 255} }, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil)

	// Test music with empty table
	err := L.DoString(`rf.music({}, 120, 0.3)`)
	if err != nil {
		t.Fatalf("rf.music with empty table should work")
	}

	// Test music with invalid table (non-string values)
	err = L.DoString(`rf.music({1, 2, 3}, 120, 0.3)`)
	if err != nil {
		t.Fatalf("rf.music with non-string values should work (will skip)")
	}

	// Test music with nil table (should error)
	err = L.DoString(`rf.music(nil, 120, 0.3)`)
	if err == nil {
		t.Fatalf("rf.music with nil should error")
	}

	// Test music with invalid BPM
	err = L.DoString(`rf.music({"4C1"}, -10, 0.3)`) // negative BPM
	if err != nil {
		t.Fatalf("rf.music with negative BPM should work (audio handles it)")
	}

	err = L.DoString(`rf.music({"4C1"}, 0, 0.3)`) // zero BPM
	if err != nil {
		t.Fatalf("rf.music with zero BPM should work (audio uses default)")
	}

	// Test music with invalid gain
	err = L.DoString(`rf.music({"4C1"}, 120, -1.0)`) // negative gain
	if err != nil {
		t.Fatalf("rf.music with negative gain should work")
	}

	err = L.DoString(`rf.music({"4C1"}, 120, 999.0)`) // very large gain
	if err != nil {
		t.Fatalf("rf.music with large gain should work")
	}
}

func TestLuaAudioEdgeCases(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(100, 100)
	Register(L, r, func(i int) (rgba [4]uint8) { return [4]uint8{0, 0, 0, 255} }, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil)

	// Test sfx with invalid name
	err := L.DoString(`rf.sfx("nonexistent")`)
	if err != nil {
		t.Fatalf("rf.sfx with invalid name should work (does nothing)")
	}

	// Test sfx with empty string
	err = L.DoString(`rf.sfx("")`)
	if err != nil {
		t.Fatalf("rf.sfx with empty string should work")
	}

	// Test tone with invalid parameters
	err = L.DoString(`rf.tone(-100, 0.1, 0.3)`) // negative frequency
	if err != nil {
		t.Fatalf("rf.tone with negative freq should work")
	}

	err = L.DoString(`rf.tone(0, 0.1, 0.3)`) // zero frequency
	if err != nil {
		t.Fatalf("rf.tone with zero freq should work")
	}

	err = L.DoString(`rf.tone(440, -0.1, 0.3)`) // negative duration
	if err != nil {
		t.Fatalf("rf.tone with negative duration should work")
	}

	err = L.DoString(`rf.tone(440, 999999, 0.3)`) // very long duration
	if err != nil {
		t.Fatalf("rf.tone with long duration should work")
	}

	// Test noise with invalid parameters
	err = L.DoString(`rf.noise(-0.1, 0.3)`)
	if err != nil {
		t.Fatalf("rf.noise with negative duration should work")
	}

	err = L.DoString(`rf.noise(0, 0.3)`)
	if err != nil {
		t.Fatalf("rf.noise with zero duration should work")
	}
}

func TestLuaNilColorByIndex(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(100, 100)

	// Test with nil colorByIndex (should panic, but we test it doesn't crash setup)
	// Actually, Register requires non-nil, so we test with a function that returns zeros
	colorByIndex := func(i int) (rgba [4]uint8) {
		return [4]uint8{0, 0, 0, 0} // zero alpha
	}

	Register(L, r, colorByIndex, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil)

	// Should work even with zero alpha
	err := L.DoString(`rf.clear_i(0)`)
	if err != nil {
		t.Fatalf("rf.clear_i should work with any colorByIndex function")
	}
}

func TestLuaNilSetPalette(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(100, 100)

	// Test with nil setPalette
	Register(L, r, func(i int) (rgba [4]uint8) { return [4]uint8{0, 0, 0, 255} }, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil)

	// palette_set should work even if setPalette is nil
	err := L.DoString(`rf.palette_set("default")`)
	if err != nil {
		t.Fatalf("rf.palette_set should work even with nil callback")
	}
}

func TestRegisterWithDev(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	colorByIndex := func(i int) (rgba [4]uint8) { return [4]uint8{0, 0, 0, 255} }
	setPalette := func(name string) {}

	// Create a mock dev mode handler
	mockDevMode := &mockDevModeHandler{enabled: true}

	RegisterWithDev(L, r, colorByIndex, setPalette, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil, mockDevMode)

	// Verify rf table exists
	rf := L.GetGlobal("rf")
	if rf == lua.LNil {
		t.Fatal("rf table should exist")
	}

	// Verify debug functions are available (when dev mode is enabled)
	printh := L.GetField(rf.(*lua.LTable), "printh")
	if printh == lua.LNil {
		t.Error("rf.printh should exist when dev mode is enabled")
	}

	stat := L.GetField(rf.(*lua.LTable), "stat")
	if stat == lua.LNil {
		t.Error("rf.stat should exist when dev mode is enabled")
	}

	timeFunc := L.GetField(rf.(*lua.LTable), "time")
	if timeFunc == lua.LNil {
		t.Error("rf.time should exist when dev mode is enabled")
	}

	// Test calling debug functions
	err := L.DoString(`rf.printh("test message")`)
	if err != nil {
		t.Errorf("rf.printh failed: %v", err)
	}
}

// mockDevModeHandler implements DevModeHandler for testing
type mockDevModeHandler struct {
	enabled bool
	logs    []string
	stats   interface{}
}

func (m *mockDevModeHandler) IsEnabled() bool {
	return m.enabled
}

func (m *mockDevModeHandler) AddDebugLog(msg string) {
	m.logs = append(m.logs, msg)
}

func (m *mockDevModeHandler) GetStats() interface{} {
	if m.stats == nil {
		return struct {
			FPS         float64
			FrameCount  int64
			LuaMemory   int64
			LoadTime    interface{}
			LastReload  interface{}
			ReloadCount int
		}{
			FPS:        60.0,
			FrameCount: 1000,
			LuaMemory:  1024,
		}
	}
	return m.stats
}
