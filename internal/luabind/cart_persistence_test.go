package luabind

import (
	"testing"

	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
	"github.com/AndrewDonelson/retroforge-engine/internal/rendersoft"
	lua "github.com/yuin/gopher-lua"
)

func TestCstore(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	state := NewState()

	RegisterWithState(L, r, func(i int) (rgba [4]uint8) {
		return [4]uint8{0, 0, 0, 255}
	}, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil, state)

	// Set up test data in runtime memory
	runtimeMem := state.GetMemory()
	for i := 0; i < 100; i++ {
		runtimeMem[i] = byte(i)
	}

	// Test cstore: copy from runtime memory to cart storage
	err := L.DoString(`rf.cstore(0, 0, 100)`)
	if err != nil {
		t.Fatalf("Failed to cstore: %v", err)
	}

	// Verify data was copied to cart storage
	cartStore := state.GetCartStore()
	for i := 0; i < 100; i++ {
		if cartStore[i] != byte(i) {
			t.Errorf("Expected cartStore[%d] = %d, got %d", i, i, cartStore[i])
		}
	}

	// Test partial copy
	runtimeMem[200] = 0xFF
	runtimeMem[201] = 0xEE
	runtimeMem[202] = 0xDD

	err = L.DoString(`rf.cstore(10, 200, 3)`)
	if err != nil {
		t.Fatalf("Failed to cstore partial: %v", err)
	}

	if cartStore[10] != 0xFF || cartStore[11] != 0xEE || cartStore[12] != 0xDD {
		t.Errorf("Partial copy failed: got [%d, %d, %d]", cartStore[10], cartStore[11], cartStore[12])
	}
}

func TestReload(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	state := NewState()

	RegisterWithState(L, r, func(i int) (rgba [4]uint8) {
		return [4]uint8{0, 0, 0, 255}
	}, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil, state)

	// Set up test data in cart storage
	cartStore := state.GetCartStore()
	for i := 0; i < 100; i++ {
		cartStore[i] = byte(i + 100)
	}

	// Clear runtime memory first
	runtimeMem := state.GetMemory()
	for i := 0; i < 100; i++ {
		runtimeMem[i] = 0
	}

	// Test reload: copy from cart storage to runtime memory
	err := L.DoString(`rf.reload(0, 0, 100)`)
	if err != nil {
		t.Fatalf("Failed to reload: %v", err)
	}

	// Verify data was copied to runtime memory
	for i := 0; i < 100; i++ {
		if runtimeMem[i] != byte(i+100) {
			t.Errorf("Expected runtimeMem[%d] = %d, got %d", i, i+100, runtimeMem[i])
		}
	}

	// Test partial reload
	cartStore[50] = 0xAA
	cartStore[51] = 0xBB
	cartStore[52] = 0xCC

	err = L.DoString(`rf.reload(500, 50, 3)`)
	if err != nil {
		t.Fatalf("Failed to reload partial: %v", err)
	}

	if runtimeMem[500] != 0xAA || runtimeMem[501] != 0xBB || runtimeMem[502] != 0xCC {
		t.Errorf("Partial reload failed: got [%d, %d, %d]", runtimeMem[500], runtimeMem[501], runtimeMem[502])
	}
}

func TestCstoreReloadRoundTrip(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	state := NewState()

	RegisterWithState(L, r, func(i int) (rgba [4]uint8) {
		return [4]uint8{0, 0, 0, 255}
	}, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil, state)

	// Set up test data
	runtimeMem := state.GetMemory()
	testData := []byte{0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE, 0xF0}
	copy(runtimeMem[1000:1008], testData)

	// Store to cart
	err := L.DoString(`rf.cstore(2000, 1000, 8)`)
	if err != nil {
		t.Fatalf("Failed to cstore: %v", err)
	}

	// Clear runtime memory
	for i := 1000; i < 1008; i++ {
		runtimeMem[i] = 0
	}

	// Reload from cart
	err = L.DoString(`rf.reload(1000, 2000, 8)`)
	if err != nil {
		t.Fatalf("Failed to reload: %v", err)
	}

	// Verify round trip
	for i := 0; i < 8; i++ {
		if runtimeMem[1000+i] != testData[i] {
			t.Errorf("Round trip failed at index %d: expected %02X, got %02X", i, testData[i], runtimeMem[1000+i])
		}
	}
}

func TestCstoreEdgeCases(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	state := NewState()

	RegisterWithState(L, r, func(i int) (rgba [4]uint8) {
		return [4]uint8{0, 0, 0, 255}
	}, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil, state)

	runtimeMem := state.GetMemory()
	runtimeMem[0] = 0x42

	// Test negative addresses (should do nothing)
	err := L.DoString(`rf.cstore(-1, 0, 10)`)
	if err != nil {
		t.Fatalf("cstore with negative dest should not error: %v", err)
	}

	err = L.DoString(`rf.cstore(0, -1, 10)`)
	if err != nil {
		t.Fatalf("cstore with negative src should not error: %v", err)
	}

	// Test negative length (should do nothing)
	err = L.DoString(`rf.cstore(0, 0, -1)`)
	if err != nil {
		t.Fatalf("cstore with negative length should not error: %v", err)
	}

	// Test out of bounds source address
	err = L.DoString(`rf.cstore(0, 3000000, 10)`)
	if err != nil {
		t.Fatalf("cstore with out of bounds src should not error: %v", err)
	}

	// Test out of bounds destination address
	err = L.DoString(`rf.cstore(100000, 0, 10)`)
	if err != nil {
		t.Fatalf("cstore with out of bounds dest should not error: %v", err)
	}

	// Test zero length
	err = L.DoString(`rf.cstore(0, 0, 0)`)
	if err != nil {
		t.Fatalf("cstore with zero length should not error: %v", err)
	}

	// Test length that exceeds source
	runtimeMem[100] = 0xAA
	runtimeMem[101] = 0xBB
	err = L.DoString(`rf.cstore(0, 100, 1000000)`)
	if err != nil {
		t.Fatalf("cstore with excessive length should clamp: %v", err)
	}

	cartStore := state.GetCartStore()
	// Should have copied available data
	if cartStore[0] != 0xAA || cartStore[1] != 0xBB {
		t.Errorf("Expected clamped copy, got [%d, %d]", cartStore[0], cartStore[1])
	}
}

func TestReloadEdgeCases(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	state := NewState()

	RegisterWithState(L, r, func(i int) (rgba [4]uint8) {
		return [4]uint8{0, 0, 0, 255}
	}, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil, state)

	cartStore := state.GetCartStore()
	cartStore[0] = 0x42

	// Test negative addresses (should do nothing)
	err := L.DoString(`rf.reload(-1, 0, 10)`)
	if err != nil {
		t.Fatalf("reload with negative dest should not error: %v", err)
	}

	err = L.DoString(`rf.reload(0, -1, 10)`)
	if err != nil {
		t.Fatalf("reload with negative src should not error: %v", err)
	}

	// Test negative length (should do nothing)
	err = L.DoString(`rf.reload(0, 0, -1)`)
	if err != nil {
		t.Fatalf("reload with negative length should not error: %v", err)
	}

	// Test out of bounds source address
	err = L.DoString(`rf.reload(0, 100000, 10)`)
	if err != nil {
		t.Fatalf("reload with out of bounds src should not error: %v", err)
	}

	// Test out of bounds destination address
	err = L.DoString(`rf.reload(3000000, 0, 10)`)
	if err != nil {
		t.Fatalf("reload with out of bounds dest should not error: %v", err)
	}

	// Test zero length
	err = L.DoString(`rf.reload(0, 0, 0)`)
	if err != nil {
		t.Fatalf("reload with zero length should not error: %v", err)
	}

	// Test length that exceeds source
	cartStore[100] = 0xAA
	cartStore[101] = 0xBB
	err = L.DoString(`rf.reload(0, 100, 1000000)`)
	if err != nil {
		t.Fatalf("reload with excessive length should clamp: %v", err)
	}

	runtimeMem := state.GetMemory()
	// Should have copied available data
	if runtimeMem[0] != 0xAA || runtimeMem[1] != 0xBB {
		t.Errorf("Expected clamped copy, got [%d, %d]", runtimeMem[0], runtimeMem[1])
	}
}

func TestCstoreReloadOverlap(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	state := NewState()

	RegisterWithState(L, r, func(i int) (rgba [4]uint8) {
		return [4]uint8{0, 0, 0, 255}
	}, nil, make(cartio.SFXMap), make(cartio.MusicMap), make(cartio.SpriteMap), nil, state)

	// Test overlapping regions (copy should handle correctly)
	runtimeMem := state.GetMemory()
	for i := 0; i < 20; i++ {
		runtimeMem[i] = byte(i)
	}

	// Store with overlap
	err := L.DoString(`rf.cstore(5, 0, 15)`)
	if err != nil {
		t.Fatalf("Failed to cstore with overlap: %v", err)
	}

	cartStore := state.GetCartStore()
	// Verify overlap was handled correctly
	for i := 0; i < 15; i++ {
		expected := byte(i)
		if cartStore[5+i] != expected {
			t.Errorf("Overlap copy failed at cart[%d]: expected %d, got %d", 5+i, expected, cartStore[5+i])
		}
	}
}
