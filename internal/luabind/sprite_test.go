package luabind

import (
	"fmt"
	"testing"

	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
	"github.com/AndrewDonelson/retroforge-engine/internal/rendersoft"
	lua "github.com/yuin/gopher-lua"
)

func TestNewSprite(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	colorByIndex := func(i int) (rgba [4]uint8) { return [4]uint8{255, 255, 255, 255} }
	setPalette := func(name string) {}
	spritesMap := make(cartio.SpriteMap)

	Register(L, r, colorByIndex, setPalette, make(cartio.SFXMap), make(cartio.MusicMap), spritesMap, nil, nil)

	// Test creating a new sprite
	err := L.DoString(`
		local sprite = rf.newSprite("test_sprite", 16, 16)
		if sprite == nil then
			error("sprite should not be nil")
		end
		if sprite.width ~= 16 then
			error("width should be 16")
		end
		if sprite.height ~= 16 then
			error("height should be 16")
		end
		if sprite.isUI ~= true then
			error("isUI should default to true")
		end
		if sprite.lifetime ~= 0 then
			error("lifetime should default to 0")
		end
		if sprite.maxSpawn ~= 0 then
			error("maxSpawn should default to 0")
		end
	`)

	if err != nil {
		t.Fatalf("Lua error: %v", err)
	}

	// Verify sprite exists in map
	sprite, ok := spritesMap["test_sprite"]
	if !ok {
		t.Fatal("sprite should be in map")
	}
	if sprite.Width != 16 || sprite.Height != 16 {
		t.Error("sprite dimensions incorrect")
	}
	if !sprite.IsUI {
		t.Error("sprite IsUI should be true")
	}
}

func TestSpritePset(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	colorByIndex := func(i int) (rgba [4]uint8) { return [4]uint8{255, 255, 255, 255} }
	setPalette := func(name string) {}
	spritesMap := make(cartio.SpriteMap)

	Register(L, r, colorByIndex, setPalette, make(cartio.SFXMap), make(cartio.MusicMap), spritesMap, nil, nil)

	err := L.DoString(`
		rf.newSprite("test", 8, 8)
		rf.sprite_pset("test", 4, 4, 7)
	`)

	if err != nil {
		t.Fatalf("Lua error: %v", err)
	}

	sprite := spritesMap["test"]
	if sprite.Pixels[4][4] != 7 {
		t.Errorf("pixel should be 7, got %d", sprite.Pixels[4][4])
	}
}

func TestSpriteLine(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	colorByIndex := func(i int) (rgba [4]uint8) { return [4]uint8{255, 255, 255, 255} }
	setPalette := func(name string) {}
	spritesMap := make(cartio.SpriteMap)

	Register(L, r, colorByIndex, setPalette, make(cartio.SFXMap), make(cartio.MusicMap), spritesMap, nil, nil)

	err := L.DoString(`
		rf.newSprite("test", 16, 16)
		rf.sprite_line("test", 0, 0, 15, 15, 5)
	`)

	if err != nil {
		t.Fatalf("Lua error: %v", err)
	}

	sprite := spritesMap["test"]
	// Check diagonal line pixels
	if sprite.Pixels[0][0] != 5 {
		t.Error("line start pixel not set")
	}
	if sprite.Pixels[15][15] != 5 {
		t.Error("line end pixel not set")
	}
}

func TestSpriteRect(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	colorByIndex := func(i int) (rgba [4]uint8) { return [4]uint8{255, 255, 255, 255} }
	setPalette := func(name string) {}
	spritesMap := make(cartio.SpriteMap)

	Register(L, r, colorByIndex, setPalette, make(cartio.SFXMap), make(cartio.MusicMap), spritesMap, nil, nil)

	err := L.DoString(`
		rf.newSprite("test", 16, 16)
		rf.sprite_rect("test", 2, 2, 13, 13, 3)
	`)

	if err != nil {
		t.Fatalf("Lua error: %v", err)
	}

	sprite := spritesMap["test"]
	// Check corners
	if sprite.Pixels[2][2] != 3 {
		t.Error("top-left corner not set")
	}
	if sprite.Pixels[2][13] != 3 {
		t.Error("top-right corner not set")
	}
	if sprite.Pixels[13][2] != 3 {
		t.Error("bottom-left corner not set")
	}
	if sprite.Pixels[13][13] != 3 {
		t.Error("bottom-right corner not set")
	}
}

func TestSpriteRectfill(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	colorByIndex := func(i int) (rgba [4]uint8) { return [4]uint8{255, 255, 255, 255} }
	setPalette := func(name string) {}
	spritesMap := make(cartio.SpriteMap)

	Register(L, r, colorByIndex, setPalette, make(cartio.SFXMap), make(cartio.MusicMap), spritesMap, nil, nil)

	err := L.DoString(`
		rf.newSprite("test", 16, 16)
		rf.sprite_rectfill("test", 4, 4, 11, 11, 8)
	`)

	if err != nil {
		t.Fatalf("Lua error: %v", err)
	}

	sprite := spritesMap["test"]
	// Check filled area
	for y := 4; y <= 11; y++ {
		for x := 4; x <= 11; x++ {
			if sprite.Pixels[y][x] != 8 {
				t.Errorf("pixel at (%d, %d) should be 8, got %d", x, y, sprite.Pixels[y][x])
				return
			}
		}
	}
}

func TestSpriteCirc(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	colorByIndex := func(i int) (rgba [4]uint8) { return [4]uint8{255, 255, 255, 255} }
	setPalette := func(name string) {}
	spritesMap := make(cartio.SpriteMap)

	Register(L, r, colorByIndex, setPalette, make(cartio.SFXMap), make(cartio.MusicMap), spritesMap, nil, nil)

	err := L.DoString(`
		rf.newSprite("test", 16, 16)
		rf.sprite_circ("test", 8, 8, 6, 9)
	`)

	if err != nil {
		t.Fatalf("Lua error: %v", err)
	}

	sprite := spritesMap["test"]
	// Check circle outline points (at least some should be set)
	found := false
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			if sprite.Pixels[y][x] == 9 {
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	if !found {
		t.Error("circle outline should have some pixels set")
	}
}

func TestSpriteCircfill(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	colorByIndex := func(i int) (rgba [4]uint8) { return [4]uint8{255, 255, 255, 255} }
	setPalette := func(name string) {}
	spritesMap := make(cartio.SpriteMap)

	Register(L, r, colorByIndex, setPalette, make(cartio.SFXMap), make(cartio.MusicMap), spritesMap, nil, nil)

	err := L.DoString(`
		rf.newSprite("test", 16, 16)
		rf.sprite_circfill("test", 8, 8, 5, 10)
	`)

	if err != nil {
		t.Fatalf("Lua error: %v", err)
	}

	sprite := spritesMap["test"]
	// Center should be filled
	if sprite.Pixels[8][8] != 10 {
		t.Error("circle center should be filled")
	}
}

func TestSetSpriteProperty(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	colorByIndex := func(i int) (rgba [4]uint8) { return [4]uint8{255, 255, 255, 255} }
	setPalette := func(name string) {}
	spritesMap := make(cartio.SpriteMap)

	Register(L, r, colorByIndex, setPalette, make(cartio.SFXMap), make(cartio.MusicMap), spritesMap, nil, nil)

	err := L.DoString(`
		rf.newSprite("test", 8, 8)
		rf.setSpriteProperty("test", "isUI", false)
		rf.setSpriteProperty("test", "useCollision", true)
		rf.setSpriteProperty("test", "lifetime", 5000)
		rf.setSpriteProperty("test", "maxSpawn", 20)
	`)

	if err != nil {
		t.Fatalf("Lua error: %v", err)
	}

	sprite := spritesMap["test"]
	if sprite.IsUI != false {
		t.Error("IsUI should be false")
	}
	if sprite.UseCollision != true {
		t.Error("UseCollision should be true")
	}
	if sprite.Lifetime != 5000 {
		t.Error("Lifetime should be 5000")
	}
	if sprite.MaxSpawn != 20 {
		t.Error("MaxSpawn should be 20")
	}
}

func TestSpriteProperties(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	colorByIndex := func(i int) (rgba [4]uint8) { return [4]uint8{255, 255, 255, 255} }
	setPalette := func(name string) {}
	spritesMap := make(cartio.SpriteMap)

	Register(L, r, colorByIndex, setPalette, make(cartio.SFXMap), make(cartio.MusicMap), spritesMap, nil, nil)

	err := L.DoString(`
		rf.newSprite("test", 8, 8)
		local sprite = rf.sprite("test")
		if sprite.isUI ~= true then
			error("isUI should be true")
		end
		if sprite.lifetime ~= 0 then
			error("lifetime should be 0")
		end
		if sprite.maxSpawn ~= 0 then
			error("maxSpawn should be 0")
		end
	`)

	if err != nil {
		t.Fatalf("Lua error: %v", err)
	}
}

func TestSpriteErrors(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	colorByIndex := func(i int) (rgba [4]uint8) { return [4]uint8{255, 255, 255, 255} }
	setPalette := func(name string) {}
	spritesMap := make(cartio.SpriteMap)

	Register(L, r, colorByIndex, setPalette, make(cartio.SFXMap), make(cartio.MusicMap), spritesMap, nil, nil)

	// Test invalid dimensions - too small (< 2x2)
	err := L.DoString(`rf.newSprite("test", 1, 10)`)
	if err == nil {
		t.Error("should error on width < 2")
	}

	err = L.DoString(`rf.newSprite("test", 10, 1)`)
	if err == nil {
		t.Error("should error on height < 2")
	}

	err = L.DoString(`rf.newSprite("test", 1, 1)`)
	if err == nil {
		t.Error("should error on 1x1 sprite")
	}

	err = L.DoString(`rf.newSprite("test", 0, 10)`)
	if err == nil {
		t.Error("should error on width = 0")
	}

	err = L.DoString(`rf.newSprite("test", 10, 0)`)
	if err == nil {
		t.Error("should error on height = 0")
	}

	err = L.DoString(`rf.newSprite("test", -1, 10)`)
	if err == nil {
		t.Error("should error on negative width")
	}

	err = L.DoString(`rf.newSprite("test", 10, -1)`)
	if err == nil {
		t.Error("should error on negative height")
	}

	// Test gameplay sprite too large (> 32x32)
	err = L.DoString(`rf.newSprite("test", 33, 32, false)`)
	if err == nil {
		t.Error("should error on gameplay sprite width > 32")
	}

	err = L.DoString(`rf.newSprite("test", 32, 33, false)`)
	if err == nil {
		t.Error("should error on gameplay sprite height > 32")
	}

	err = L.DoString(`rf.newSprite("test", 33, 33, false)`)
	if err == nil {
		t.Error("should error on gameplay sprite 33x33")
	}

	// Test UI sprite too large (> 256)
	err = L.DoString(`rf.newSprite("test", 257, 256, true)`)
	if err == nil {
		t.Error("should error on UI sprite width > 256")
	}

	err = L.DoString(`rf.newSprite("test", 256, 257, true)`)
	if err == nil {
		t.Error("should error on UI sprite height > 256")
	}

	// Test UI sprite odd dimensions (not divisible by 2)
	err = L.DoString(`rf.newSprite("test", 3, 4, true)`)
	if err == nil {
		t.Error("should error on UI sprite with odd width")
	}

	err = L.DoString(`rf.newSprite("test", 4, 3, true)`)
	if err == nil {
		t.Error("should error on UI sprite with odd height")
	}

	err = L.DoString(`rf.newSprite("test", 3, 3, true)`)
	if err == nil {
		t.Error("should error on UI sprite with odd dimensions")
	}

	// Test sprite not found
	err = L.DoString(`rf.sprite_pset("nonexistent", 0, 0, 1)`)
	if err == nil {
		t.Error("should error on nonexistent sprite")
	}
}

func TestSpriteSizeValidation(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	colorByIndex := func(i int) (rgba [4]uint8) { return [4]uint8{255, 255, 255, 255} }
	setPalette := func(name string) {}
	spritesMap := make(cartio.SpriteMap)

	Register(L, r, colorByIndex, setPalette, make(cartio.SFXMap), make(cartio.MusicMap), spritesMap, nil, nil)

	// Test valid gameplay sprites (2x2 to 32x32)
	testCases := []struct {
		name     string
		width    int
		height   int
		isUI     bool
		shouldErr bool
	}{
		{"gameplay_2x2", 2, 2, false, false},
		{"gameplay_8x8", 8, 8, false, false},
		{"gameplay_16x16", 16, 16, false, false},
		{"gameplay_32x32", 32, 32, false, false},
		{"gameplay_3x5", 3, 5, false, false}, // Any size 2-32
		{"gameplay_7x11", 7, 11, false, false},
		{"gameplay_15x23", 15, 23, false, false},
		{"gameplay_31x2", 31, 2, false, false},
		{"gameplay_2x31", 2, 31, false, false},
		{"ui_2x2", 2, 2, true, false},
		{"ui_4x4", 4, 4, true, false},
		{"ui_8x8", 8, 8, true, false},
		{"ui_16x16", 16, 16, true, false},
		{"ui_32x32", 32, 32, true, false},
		{"ui_64x64", 64, 64, true, false},
		{"ui_128x128", 128, 128, true, false},
		{"ui_256x256", 256, 256, true, false},
		{"ui_2x256", 2, 256, true, false},
		{"ui_256x2", 256, 2, true, false},
		{"ui_4x128", 4, 128, true, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			spritesMap := make(cartio.SpriteMap)
			Register(L, r, colorByIndex, setPalette, make(cartio.SFXMap), make(cartio.MusicMap), spritesMap, nil, nil)
			
			var luaCode string
			if tc.isUI {
				luaCode = fmt.Sprintf(`rf.newSprite("%s", %d, %d, true)`, tc.name, tc.width, tc.height)
			} else {
				luaCode = fmt.Sprintf(`rf.newSprite("%s", %d, %d, false)`, tc.name, tc.width, tc.height)
			}
			
			err := L.DoString(luaCode)
			if tc.shouldErr {
				if err == nil {
					t.Errorf("expected error for %s (%dx%d, isUI=%v)", tc.name, tc.width, tc.height, tc.isUI)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %s (%dx%d, isUI=%v): %v", tc.name, tc.width, tc.height, tc.isUI, err)
				}
				// Verify sprite was created
				if _, ok := spritesMap[tc.name]; !ok {
					t.Errorf("sprite %s should exist in map", tc.name)
				}
			}
		})
	}
}

func TestSpriteIsUI(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	colorByIndex := func(i int) (rgba [4]uint8) { return [4]uint8{255, 255, 255, 255} }
	setPalette := func(name string) {}
	spritesMap := make(cartio.SpriteMap)

	Register(L, r, colorByIndex, setPalette, make(cartio.SFXMap), make(cartio.MusicMap), spritesMap, nil, nil)

	// Test default isUI=true
	err := L.DoString(`rf.newSprite("ui_default", 16, 16)`)
	if err != nil {
		t.Fatalf("Lua error: %v", err)
	}
	sprite, ok := spritesMap["ui_default"]
	if !ok {
		t.Fatal("sprite should exist")
	}
	if !sprite.IsUI {
		t.Error("default isUI should be true")
	}

	// Test explicit isUI=true
	err = L.DoString(`rf.newSprite("ui_explicit", 16, 16, true)`)
	if err != nil {
		t.Fatalf("Lua error: %v", err)
	}
	sprite, ok = spritesMap["ui_explicit"]
	if !ok {
		t.Fatal("sprite should exist")
	}
	if !sprite.IsUI {
		t.Error("explicit isUI=true should be true")
	}

	// Test explicit isUI=false
	err = L.DoString(`rf.newSprite("gameplay", 16, 16, false)`)
	if err != nil {
		t.Fatalf("Lua error: %v", err)
	}
	sprite, ok = spritesMap["gameplay"]
	if !ok {
		t.Fatal("sprite should exist")
	}
	if sprite.IsUI {
		t.Error("explicit isUI=false should be false")
	}
}
