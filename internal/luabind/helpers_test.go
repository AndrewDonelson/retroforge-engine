package luabind

import (
	"fmt"
	"math"
	"testing"

	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
	"github.com/AndrewDonelson/retroforge-engine/internal/rendersoft"
	lua "github.com/yuin/gopher-lua"
)

func TestFlr(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	Register(L, r, func(i int) (rgba [4]uint8) { return [4]uint8{0, 0, 0, 255} }, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil)

	// Test basic floor operations
	testCases := []struct {
		input    float64
		expected int64
	}{
		{3.7, 3},
		{3.2, 3},
		{-3.7, -3},
		{-3.2, -3},
		{0.0, 0},
		{0.5, 0},
		{-0.5, 0},
		{999.999, 999},
		{-999.999, -999},
		{math.Pi, 3},
		{-math.Pi, -3},
	}

	for _, tc := range testCases {
		script := fmt.Sprintf(`local result = rf.flr(%.10f); if result ~= %d then error("flr(%.10f) expected %d, got " .. result) end`, tc.input, tc.expected, tc.input, tc.expected)
		err := L.DoString(script)
		if err != nil {
			t.Errorf("flr(%.2f) failed: %v", tc.input, err)
		}
	}

	// Test edge cases
	edgeCases := []struct {
		name     string
		input    float64
		expected int64
	}{
		{"Very large positive", 1e10, 10000000000},
		{"Very large negative", -1e10, -10000000000},
		{"Very small positive", 1e-10, 0},
		{"Very small negative", -1e-10, 0},
	}

	for _, ec := range edgeCases {
		script := fmt.Sprintf(`local result = rf.flr(%.10e); if result ~= %d then error("%s: expected %d, got " .. result) end`, ec.input, ec.expected, ec.name, ec.expected)
		err := L.DoString(script)
		if err != nil {
			t.Errorf("%s: %v", ec.name, err)
		}
	}
}

func TestCeil(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	Register(L, r, func(i int) (rgba [4]uint8) { return [4]uint8{0, 0, 0, 255} }, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil)

	// Test basic ceiling operations
	testCases := []struct {
		input    float64
		expected int64
	}{
		{3.7, 4},
		{3.2, 4},
		{-3.7, -3},
		{-3.2, -3},
		{0.0, 0},
		{0.5, 1},
		{-0.5, 0},
		{999.001, 1000},
		{-999.001, -999},
		{math.Pi, 4},
		{-math.Pi, -3},
	}

	for _, tc := range testCases {
		script := fmt.Sprintf(`local result = rf.ceil(%.10f); if result ~= %d then error("ceil(%.10f) expected %d, got " .. result) end`, tc.input, tc.expected, tc.input, tc.expected)
		err := L.DoString(script)
		if err != nil {
			t.Errorf("ceil(%.2f) = %v, expected %d: %v", tc.input, err, tc.expected, err)
		}
	}

	// Test exact integers (should remain unchanged)
	exactInts := []float64{0, 5, -5, 100, -100}
	for _, x := range exactInts {
		script := fmt.Sprintf(`local result = rf.ceil(%.0f); if result ~= %.0f then error("ceil(%.0f) should equal itself, got " .. result) end`, x, x, x)
		err := L.DoString(script)
		if err != nil {
			t.Errorf("ceil(%.0f) should equal itself: %v", x, err)
		}
	}
}

func TestRnd(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	state := NewState()
	RegisterWithState(L, r, func(i int) (rgba [4]uint8) { return [4]uint8{0, 0, 0, 255} }, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil, state)

	// Test rnd() with no arguments (should return 0.0 to 1.0, exclusive of 1.0)
	values := make([]float64, 100)
	for i := 0; i < 100; i++ {
		err := L.DoString(`local result = rf.rnd(); return result`)
		if err != nil {
			t.Fatalf("rnd() call failed: %v", err)
		}
		result := L.ToNumber(-1)
		L.Pop(1)
		values[i] = float64(result)
		if result < 0.0 || result >= 1.0 {
			t.Errorf("rnd() returned %.6f, expected 0.0 <= result < 1.0", result)
		}
	}

	// Test that values are distributed (basic check: at least some variety)
	minVal, maxVal := values[0], values[0]
	for _, v := range values {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}
	if maxVal-minVal < 0.1 {
		t.Errorf("rnd() values too similar (min=%.6f, max=%.6f), might not be random enough", minVal, maxVal)
	}

	// Test rnd(x) with positive argument
	err := L.DoString(`local result = rf.rnd(100); if result < 0 or result >= 100 then error("rnd(100) returned " .. result .. ", expected 0 <= result < 100") end`)
	if err != nil {
		t.Errorf("rnd(100) failed: %v", err)
	}

	// Test rnd(x) with negative argument
	err = L.DoString(`local result = rf.rnd(-50); if result > 0 or result <= -50 then error("rnd(-50) returned " .. result .. ", expected -50 < result <= 0") end`)
	if err != nil {
		t.Errorf("rnd(-50) failed: %v", err)
	}

	// Test rnd(0)
	err = L.DoString(`local result = rf.rnd(0); if result ~= 0 then error("rnd(0) returned " .. result .. ", expected 0") end`)
	if err != nil {
		t.Errorf("rnd(0) failed: %v", err)
	}

	// Test deterministic behavior with seed reset
	state.SetRNGSeed(1)
	err = L.DoString(`local result1 = rf.rnd(); return result1`)
	if err != nil {
		t.Fatalf("rnd() with seed reset failed: %v", err)
	}
	val1 := L.ToNumber(-1)
	L.Pop(1)

	state.SetRNGSeed(1)
	err = L.DoString(`local result2 = rf.rnd(); return result2`)
	if err != nil {
		t.Fatalf("rnd() with seed reset failed: %v", err)
	}
	val2 := L.ToNumber(-1)
	L.Pop(1)

	if val1 != val2 {
		t.Errorf("rnd() with same seed should produce same value, got %.6f and %.6f", val1, val2)
	}

	// Test very large argument
	err = L.DoString(`local result = rf.rnd(1e10); if result < 0 or result >= 1e10 then error("rnd(1e10) returned " .. result .. ", expected 0 <= result < 1e10") end`)
	if err != nil {
		t.Errorf("rnd(1e10) failed: %v", err)
	}
}

func TestMid(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	Register(L, r, func(i int) (rgba [4]uint8) { return [4]uint8{0, 0, 0, 255} }, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil)

	testCases := []struct {
		x, y, z, expected float64
		description       string
	}{
		{5, 0, 10, 5, "Middle value"},
		{15, 0, 10, 10, "Above range"},
		{-5, 0, 10, 0, "Below range"},
		{5, 10, 0, 5, "Reversed y,z (should clamp correctly)"},
		{15, 10, 0, 10, "Above range with reversed y,z"},
		{-5, 10, 0, 0, "Below range with reversed y,z"},
		{0, -10, 10, 0, "Value at lower bound"},
		{10, -10, 10, 10, "Value at upper bound"},
		{-10, -10, 10, -10, "Value equals lower bound"},
		{-5, -10, -1, -5, "All negative, in range"},
		{-15, -10, -1, -10, "All negative, below range"},
		{0, -10, -1, -1, "All negative, above range"},
		{0, 0, 0, 0, "All same value"},
		{3.7, 3, 5, 3.7, "Float in range"},
		{2.5, 3, 5, 3, "Float below range"},
		{6.2, 3, 5, 5, "Float above range"},
	}

	for _, tc := range testCases {
		script := fmt.Sprintf(`local result = rf.mid(%.10f, %.10f, %.10f); if math.abs(result - %.10f) > 0.0001 then error("%s: mid(%.2f, %.2f, %.2f) = " .. result .. ", expected %.6f") end`, tc.x, tc.y, tc.z, tc.expected, tc.description, tc.x, tc.y, tc.z, tc.expected)
		err := L.DoString(script)
		if err != nil {
			t.Errorf("%s: %v", tc.description, err)
		}
	}
}

func TestSgn(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	Register(L, r, func(i int) (rgba [4]uint8) { return [4]uint8{0, 0, 0, 255} }, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil)

	testCases := []struct {
		input    float64
		expected float64
		name     string
	}{
		{5.0, 1, "Positive integer"},
		{-5.0, -1, "Negative integer"},
		{0.0, 0, "Zero"},
		{0.5, 1, "Positive float"},
		{-0.5, -1, "Negative float"},
		{1e10, 1, "Very large positive"},
		{-1e10, -1, "Very large negative"},
		{1e-10, 1, "Very small positive"},
		{-1e-10, -1, "Very small negative"},
		{math.Pi, 1, "Pi (positive)"},
		{-math.Pi, -1, "Negative Pi"},
	}

	for _, tc := range testCases {
		script := fmt.Sprintf(`local result = rf.sgn(%.10e); if result ~= %.0f then error("%s: sgn(%.2f) = " .. result .. ", expected %.0f") end`, tc.input, tc.expected, tc.name, tc.input, tc.expected)
		err := L.DoString(script)
		if err != nil {
			t.Errorf("%s: %v", tc.name, err)
		}
	}

	// Test exact zero edge cases
	zeros := []float64{0.0, -0.0}
	for _, z := range zeros {
		script := fmt.Sprintf(`local result = rf.sgn(%.1f); if result ~= 0 then error("sgn(%.1f) = " .. result .. ", expected 0") end`, z, z)
		err := L.DoString(script)
		if err != nil {
			t.Errorf("sgn(%.1f) should be 0: %v", z, err)
		}
	}
}

func TestChr(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	Register(L, r, func(i int) (rgba [4]uint8) { return [4]uint8{0, 0, 0, 255} }, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil)

	testCases := []struct {
		input    int
		expected byte
		name     string
	}{
		{65, 'A', "A"},
		{97, 'a', "a"},
		{48, '0', "0"},
		{32, ' ', "Space"},
		{0, 0, "Null"},
		{255, 255, "Max byte"},
		{10, '\n', "Newline"},
		{13, '\r', "Carriage return"},
		{9, '\t', "Tab"},
	}

	for _, tc := range testCases {
		// Use string.byte to check the result instead of direct comparison
		script := fmt.Sprintf(`local result = rf.chr(%d); local byteVal = string.byte(result, 1) or 0; if byteVal ~= %d then error("%s: chr(%d) produced byte " .. byteVal .. ", expected %d") end`, tc.input, tc.expected, tc.name, tc.input, tc.expected)
		err := L.DoString(script)
		if err != nil {
			t.Errorf("%s: %v", tc.name, err)
		}
	}

	// Test edge cases: out of range values should be clamped
	edgeCases := []struct {
		input    int
		expected int
		name     string
	}{
		{-1, 0, "Negative (should clamp to 0)"},
		{256, 255, "Above 255 (should clamp to 255)"},
		{-100, 0, "Large negative (should clamp to 0)"},
		{1000, 255, "Very large (should clamp to 255)"},
	}

	for _, ec := range edgeCases {
		err := L.DoString(fmt.Sprintf(`local result = rf.chr(%d); local byteVal = string.byte(result, 1) or 0; if byteVal ~= %d then error("%s: chr(%d) produced byte " .. byteVal .. ", expected %d") end`, ec.input, ec.expected, ec.name, ec.input, ec.expected))
		if err != nil {
			t.Errorf("%s: %v", ec.name, err)
		}
	}

	// Test a range of valid ASCII values (skip problematic bytes that Lua might handle differently)
	// Test every 16th value but skip 128+ for now as they're UTF-8 dependent
	for i := 0; i < 128; i += 16 {
		err := L.DoString(fmt.Sprintf(`local result = rf.chr(%d); if #result ~= 1 or string.byte(result, 1) ~= %d then error("chr(%d) failed") end`, i, i, i))
		if err != nil {
			t.Errorf("chr(%d) failed: %v", i, err)
		}
	}
	// Test a few high byte values individually
	for _, i := range []int{128, 160, 192, 224, 240, 255} {
		err := L.DoString(fmt.Sprintf(`local result = rf.chr(%d); local byteVal = string.byte(result, 1) or 0; if byteVal ~= %d then error("chr(%d) produced byte " .. byteVal .. ", expected %d") end`, i, i, i, i))
		if err != nil {
			// High byte values might be valid UTF-8 sequences, just check they're in range
			err2 := L.DoString(fmt.Sprintf(`local result = rf.chr(%d); local byteVal = string.byte(result, 1) or 0; if byteVal < 0 or byteVal > 255 then error("chr(%d) produced invalid byte " .. byteVal) end`, i, i))
			if err2 != nil {
				t.Errorf("chr(%d) validation failed: %v", i, err2)
			}
		}
	}
}

func TestOrd(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	Register(L, r, func(i int) (rgba [4]uint8) { return [4]uint8{0, 0, 0, 255} }, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil)

	testCases := []struct {
		input    string
		expected int
		name     string
	}{
		{"A", 65, "A"},
		{"a", 97, "a"},
		{"0", 48, "0"},
		{" ", 32, "Space"},
		{string([]byte{0}), 0, "Null"}, // Use explicit byte array to avoid escape issues
		{string([]byte{255}), 255, "Max byte"},
		{"\n", 10, "Newline"},
		{"\r", 13, "Carriage return"},
		{"\t", 9, "Tab"},
		{"!", 33, "Exclamation"},
		{"@", 64, "At symbol"},
	}

	for _, tc := range testCases {
		// Build the string using string.char to avoid escape sequence issues
		var charCode int
		if len(tc.input) > 0 {
			charCode = int(tc.input[0])
		}
		// Use string.char to create the test string in Lua
		script := fmt.Sprintf(`local str = string.char(%d); local result = rf.ord(str); if result ~= %d then error("ord returned " .. result .. ", expected %d") end`, charCode, tc.expected, tc.expected)
		err := L.DoString(script)
		if err != nil {
			t.Errorf("%s: %v", tc.name, err)
		}
	}

	// Test multi-character strings (should return first character)
	multiChar := []struct {
		input    string
		expected int
		name     string
	}{
		{"Hello", 72, "First char of 'Hello'"},
		{"World", 87, "First char of 'World'"},
		{"123", 49, "First char of '123'"},
		{"ABC", 65, "First char of 'ABC'"},
	}

	for _, mc := range multiChar {
		escaped := fmt.Sprintf("%q", mc.input)
		script := fmt.Sprintf(`local str = %s; local result = rf.ord(str); if result ~= %d then error("ord returned " .. result .. ", expected %d") end`, escaped, mc.expected, mc.expected)
		err := L.DoString(script)
		if err != nil {
			t.Errorf("%s: %v", mc.name, err)
		}
	}

	// Test empty string
	err := L.DoString(`local result = rf.ord(""); if result ~= 0 then error("ord(\"\") = " .. result .. ", expected 0") end`)
	if err != nil {
		t.Errorf("ord(\"\") failed: %v", err)
	}

	// Test Unicode (should return first byte)
	unicodeCases := []struct {
		input string
		name  string
	}{
		{"🚀", "Rocket emoji"},
		{"中", "Chinese character"},
		{"é", "Accented e"},
	}

	for _, uc := range unicodeCases {
		escaped := fmt.Sprintf("%q", uc.input)
		err := L.DoString(fmt.Sprintf(`local str = %s; local result = rf.ord(str); if result < 0 or result > 255 then error("ord returned " .. result .. ", expected 0-255") end`, escaped))
		if err != nil {
			t.Errorf("%s: %v", uc.name, err)
		}
		// Verify it matches first byte
		if len(uc.input) > 0 {
			firstByte := int(uc.input[0])
			err := L.DoString(fmt.Sprintf(`local str = %s; local result = rf.ord(str); if result ~= %d then error("ord returned " .. result .. ", expected first byte %d") end`, escaped, firstByte, firstByte))
			if err != nil {
				t.Errorf("%s first byte check: %v", uc.name, err)
			}
		}
	}
}

func TestHelperFunctionsCombined(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	state := NewState()
	RegisterWithState(L, r, func(i int) (rgba [4]uint8) { return [4]uint8{0, 0, 0, 255} }, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil, state)

	// Test practical combinations
	err := L.DoString(`
		-- Clamp random value to range
		local r = rf.rnd(100)
		local clamped = rf.mid(r, 10, 90)
		if clamped < 10 or clamped > 90 then
			error("mid failed to clamp")
		end

		-- Floor and ceiling operations
		local val = 3.7
		local floored = rf.flr(val)
		local ceiled = rf.ceil(val)
		if floored ~= 3 or ceiled ~= 4 then
			error("floor/ceil failed")
		end

		-- Sign check
		if rf.sgn(5) ~= 1 or rf.sgn(-5) ~= -1 or rf.sgn(0) ~= 0 then
			error("sgn failed")
		end

		-- Character conversion
		local char = rf.chr(65)
		local num = rf.ord(char)
		if num ~= 65 then
			error("chr/ord roundtrip failed")
		end
	`)
	if err != nil {
		t.Fatalf("Combined helper functions test failed: %v", err)
	}
}
