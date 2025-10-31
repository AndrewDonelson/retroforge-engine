package luabind

import (
	"testing"

	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
	"github.com/AndrewDonelson/retroforge-engine/internal/rendersoft"
	lua "github.com/yuin/gopher-lua"
)

func TestBitwiseShl(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	Register(L, r, func(i int) (rgba [4]uint8) { return [4]uint8{0, 0, 0, 255} }, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil)

	// Test basic shift left
	err := L.DoString(`local result = rf.shl(1, 3); if result ~= 8 then error("Expected 8, got " .. result) end`)
	if err != nil {
		t.Fatalf("shl(1, 3) should equal 8: %v", err)
	}

	// Test shift by 0
	err = L.DoString(`local result = rf.shl(5, 0); if result ~= 5 then error("Expected 5, got " .. result) end`)
	if err != nil {
		t.Fatalf("shl(5, 0) should equal 5: %v", err)
	}

	// Test shift by 1
	err = L.DoString(`local result = rf.shl(10, 1); if result ~= 20 then error("Expected 20, got " .. result) end`)
	if err != nil {
		t.Fatalf("shl(10, 1) should equal 20: %v", err)
	}

	// Test negative shift (should shift right)
	err = L.DoString(`local result = rf.shl(8, -2); if result ~= 2 then error("Expected 2, got " .. result) end`)
	if err != nil {
		t.Fatalf("shl(8, -2) should equal 2: %v", err)
	}

	// Test large shift (should clamp to 0)
	err = L.DoString(`local result = rf.shl(1, 100); if result ~= 0 then error("Expected 0 for large shift, got " .. result) end`)
	if err != nil {
		t.Fatalf("shl(1, 100) should equal 0: %v", err)
	}

	// Test edge case: shift by 63 (may overflow signed int, that's expected)
	err = L.DoString(`local result = rf.shl(1, 63); if result == 0 then error("Shift by 63 should not be 0") end`)
	if err != nil {
		t.Fatalf("shl(1, 63) should work: %v", err)
	}
}

func TestBitwiseShr(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	Register(L, r, func(i int) (rgba [4]uint8) { return [4]uint8{0, 0, 0, 255} }, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil)

	// Test basic shift right
	err := L.DoString(`local result = rf.shr(8, 3); if result ~= 1 then error("Expected 1, got " .. result) end`)
	if err != nil {
		t.Fatalf("shr(8, 3) should equal 1: %v", err)
	}

	// Test shift by 0
	err = L.DoString(`local result = rf.shr(5, 0); if result ~= 5 then error("Expected 5, got " .. result) end`)
	if err != nil {
		t.Fatalf("shr(5, 0) should equal 5: %v", err)
	}

	// Test shift by 1
	err = L.DoString(`local result = rf.shr(20, 1); if result ~= 10 then error("Expected 10, got " .. result) end`)
	if err != nil {
		t.Fatalf("shr(20, 1) should equal 10: %v", err)
	}

	// Test negative shift (should shift left)
	err = L.DoString(`local result = rf.shr(2, -2); if result ~= 8 then error("Expected 8, got " .. result) end`)
	if err != nil {
		t.Fatalf("shr(2, -2) should equal 8: %v", err)
	}

	// Test large shift (should be 0 for positive, -1 for negative)
	err = L.DoString(`local result = rf.shr(100, 100); if result ~= 0 then error("Expected 0 for large shift, got " .. result) end`)
	if err != nil {
		t.Fatalf("shr(100, 100) should equal 0: %v", err)
	}

	// Test sign extension with negative number
	err = L.DoString(`local result = rf.shr(-1, 1); if result ~= -1 then error("Expected -1 for shr(-1, 1), got " .. result) end`)
	if err != nil {
		t.Fatalf("shr(-1, 1) should equal -1 (sign extended): %v", err)
	}
}

func TestBitwiseBand(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	Register(L, r, func(i int) (rgba [4]uint8) { return [4]uint8{0, 0, 0, 255} }, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil)

	// Test basic AND
	err := L.DoString(`local result = rf.band(5, 3); if result ~= 1 then error("Expected 1, got " .. result) end`)
	if err != nil {
		t.Fatalf("band(5, 3) should equal 1: %v", err)
	}

	// Test AND with 0
	err = L.DoString(`local result = rf.band(5, 0); if result ~= 0 then error("Expected 0, got " .. result) end`)
	if err != nil {
		t.Fatalf("band(5, 0) should equal 0: %v", err)
	}

	// Test AND with same number
	err = L.DoString(`local result = rf.band(7, 7); if result ~= 7 then error("Expected 7, got " .. result) end`)
	if err != nil {
		t.Fatalf("band(7, 7) should equal 7: %v", err)
	}

	// Test AND with negative numbers
	err = L.DoString(`local result = rf.band(-1, 5); if result ~= 5 then error("Expected 5, got " .. result) end`)
	if err != nil {
		t.Fatalf("band(-1, 5) should equal 5: %v", err)
	}

	// Test edge case: large numbers
	err = L.DoString(`local result = rf.band(0xFFFFFFFF, 0xFFFF); if result ~= 0xFFFF then error("Expected 65535, got " .. result) end`)
	if err != nil {
		t.Fatalf("band with large numbers should work: %v", err)
	}
}

func TestBitwiseBor(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	Register(L, r, func(i int) (rgba [4]uint8) { return [4]uint8{0, 0, 0, 255} }, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil)

	// Test basic OR
	err := L.DoString(`local result = rf.bor(5, 3); if result ~= 7 then error("Expected 7, got " .. result) end`)
	if err != nil {
		t.Fatalf("bor(5, 3) should equal 7: %v", err)
	}

	// Test OR with 0
	err = L.DoString(`local result = rf.bor(5, 0); if result ~= 5 then error("Expected 5, got " .. result) end`)
	if err != nil {
		t.Fatalf("bor(5, 0) should equal 5: %v", err)
	}

	// Test OR with same number
	err = L.DoString(`local result = rf.bor(7, 7); if result ~= 7 then error("Expected 7, got " .. result) end`)
	if err != nil {
		t.Fatalf("bor(7, 7) should equal 7: %v", err)
	}

	// Test OR with negative numbers
	err = L.DoString(`local result = rf.bor(0, -1); if result ~= -1 then error("Expected -1, got " .. result) end`)
	if err != nil {
		t.Fatalf("bor(0, -1) should equal -1: %v", err)
	}
}

func TestBitwiseBxor(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	Register(L, r, func(i int) (rgba [4]uint8) { return [4]uint8{0, 0, 0, 255} }, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil)

	// Test basic XOR
	err := L.DoString(`local result = rf.bxor(5, 3); if result ~= 6 then error("Expected 6, got " .. result) end`)
	if err != nil {
		t.Fatalf("bxor(5, 3) should equal 6: %v", err)
	}

	// Test XOR with 0
	err = L.DoString(`local result = rf.bxor(5, 0); if result ~= 5 then error("Expected 5, got " .. result) end`)
	if err != nil {
		t.Fatalf("bxor(5, 0) should equal 5: %v", err)
	}

	// Test XOR with same number (should be 0)
	err = L.DoString(`local result = rf.bxor(7, 7); if result ~= 0 then error("Expected 0, got " .. result) end`)
	if err != nil {
		t.Fatalf("bxor(7, 7) should equal 0: %v", err)
	}

	// Test XOR with negative numbers
	err = L.DoString(`local result = rf.bxor(-1, 5); if result ~= -6 then error("Expected -6, got " .. result) end`)
	if err != nil {
		t.Fatalf("bxor(-1, 5) should equal -6: %v", err)
	}

	// Test XOR property: a XOR b XOR b = a
	err = L.DoString(`local a, b = 10, 7; local result = rf.bxor(rf.bxor(a, b), b); if result ~= a then error("XOR property failed") end`)
	if err != nil {
		t.Fatalf("XOR property should hold: %v", err)
	}
}

func TestBitwiseBnot(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	Register(L, r, func(i int) (rgba [4]uint8) { return [4]uint8{0, 0, 0, 255} }, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil)

	// Test NOT of 0 (should be -1 for signed integers)
	err := L.DoString(`local result = rf.bnot(0); if result ~= -1 then error("Expected -1, got " .. result) end`)
	if err != nil {
		t.Fatalf("bnot(0) should equal -1: %v", err)
	}

	// Test NOT of -1 (should be 0)
	err = L.DoString(`local result = rf.bnot(-1); if result ~= 0 then error("Expected 0, got " .. result) end`)
	if err != nil {
		t.Fatalf("bnot(-1) should equal 0: %v", err)
	}

	// Test NOT of 5
	err = L.DoString(`local result = rf.bnot(5); if result ~= -6 then error("Expected -6, got " .. result) end`)
	if err != nil {
		t.Fatalf("bnot(5) should equal -6: %v", err)
	}

	// Test NOT property: NOT (NOT x) = x
	err = L.DoString(`local x = 10; local result = rf.bnot(rf.bnot(x)); if result ~= x then error("Double NOT should equal original") end`)
	if err != nil {
		t.Fatalf("Double NOT property should hold: %v", err)
	}

	// Test NOT with negative number
	err = L.DoString(`local result = rf.bnot(-10); if result ~= 9 then error("Expected 9, got " .. result) end`)
	if err != nil {
		t.Fatalf("bnot(-10) should equal 9: %v", err)
	}
}

func TestBitwiseCombined(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	Register(L, r, func(i int) (rgba [4]uint8) { return [4]uint8{0, 0, 0, 255} }, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil)

	// Test combined operations
	err := L.DoString(`
		local a, b = 0x0F, 0xF0
		local and_result = rf.band(a, b)  -- Should be 0
		local or_result = rf.bor(a, b)     -- Should be 0xFF
		local xor_result = rf.bxor(a, b)   -- Should be 0xFF
		
		if and_result ~= 0 then error("AND failed") end
		if or_result ~= 255 then error("OR failed") end
		if xor_result ~= 255 then error("XOR failed") end
	`)
	if err != nil {
		t.Fatalf("Combined bitwise operations failed: %v", err)
	}

	// Test shift with bitwise operations
	err = L.DoString(`
		local x = 1
		local shifted = rf.shl(x, 8)  -- Should be 256
		local masked = rf.band(shifted, 0xFF)  -- Should be 0
		
		if shifted ~= 256 then error("Shift failed") end
		if masked ~= 0 then error("Mask failed") end
	`)
	if err != nil {
		t.Fatalf("Shift with bitwise operations failed: %v", err)
	}

	// Test extracting bits using shift and mask
	err = L.DoString(`
		local value = 0x1234
		local low_byte = rf.band(value, 0xFF)  -- Should be 0x34
		local high_byte = rf.shr(rf.band(value, 0xFF00), 8)  -- Should be 0x12
		
		if low_byte ~= 52 then error("Low byte extraction failed") end
		if high_byte ~= 18 then error("High byte extraction failed") end
	`)
	if err != nil {
		t.Fatalf("Bit extraction failed: %v", err)
	}
}

func TestBitwiseEdgeCases(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	Register(L, r, func(i int) (rgba [4]uint8) { return [4]uint8{0, 0, 0, 255} }, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil)

	// Test with floating point numbers (should convert to int)
	err := L.DoString(`local result = rf.band(5.7, 3.2); if result ~= 1 then error("Expected 1, got " .. result) end`)
	if err != nil {
		t.Fatalf("band with floats should work: %v", err)
	}

	// Test with large numbers (0xFFFFFFFF = 4294967295)
	err = L.DoString(`local result = rf.band(4294967295, 0xFFFFFFFF); if result ~= 4294967295 then error("Large number AND failed, got " .. result) end`)
	if err != nil {
		t.Fatalf("band with large numbers should work: %v", err)
	}

	// Test shift by negative large number
	err = L.DoString(`local result = rf.shl(1, -100); if result ~= 0 then error("Expected 0 for large negative shift") end`)
	if err != nil {
		t.Fatalf("shl with large negative shift should work: %v", err)
	}

	// Test zero in all operations
	err = L.DoString(`
		if rf.band(0, 5) ~= 0 then error("band(0, 5) should be 0") end
		if rf.bor(0, 5) ~= 5 then error("bor(0, 5) should be 5") end
		if rf.bxor(0, 5) ~= 5 then error("bxor(0, 5) should be 5") end
	`)
	if err != nil {
		t.Fatalf("Bitwise operations with zero should work: %v", err)
	}
}
