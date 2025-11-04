package luabind

import (
	"testing"

	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
	"github.com/AndrewDonelson/retroforge-engine/internal/rendersoft"
	lua "github.com/yuin/gopher-lua"
)

func TestNewSpriteFrames(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	colorByIndex := func(i int) (rgba [4]uint8) { return [4]uint8{255, 255, 255, 255} }
	setPalette := func(name string) {}
	spritesMap := make(cartio.SpriteMap)

	Register(L, r, colorByIndex, setPalette, make(cartio.SFXMap), make(cartio.MusicMap), spritesMap, nil, nil)

	// Test creating a new frames sprite
	err := L.DoString(`
		local sprite = rf.newSpriteFrames("frames_sprite", 16, 16)
		if sprite == nil then
			error("sprite should not be nil")
		end
		if sprite.width ~= 16 then
			error("width should be 16")
		end
		if sprite.height ~= 16 then
			error("height should be 16")
		end
		if sprite.type ~= "frames" then
			error("type should be 'frames'")
		end
		if sprite.isUI ~= true then
			error("isUI should default to true")
		end
	`)

	if err != nil {
		t.Fatalf("Lua error: %v", err)
	}

	// Verify sprite exists in map
	sprite, ok := spritesMap["frames_sprite"]
	if !ok {
		t.Fatal("Sprite should exist in map")
	}
	if sprite.Type != cartio.SpriteTypeFrames {
		t.Errorf("Expected type frames, got %s", sprite.Type)
	}
	if sprite.Width != 16 || sprite.Height != 16 {
		t.Errorf("Expected 16x16, got %dx%d", sprite.Width, sprite.Height)
	}
	if len(sprite.Frames) != 0 {
		t.Errorf("Expected 0 frames initially, got %d", len(sprite.Frames))
	}
}

func TestNewSpriteFrames_SizeValidation(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	colorByIndex := func(i int) (rgba [4]uint8) { return [4]uint8{255, 255, 255, 255} }
	setPalette := func(name string) {}
	spritesMap := make(cartio.SpriteMap)

	Register(L, r, colorByIndex, setPalette, make(cartio.SFXMap), make(cartio.MusicMap), spritesMap, nil, nil)

	// Test invalid size (too small)
	err := L.DoString(`rf.newSpriteFrames("invalid", 1, 1)`)
	if err == nil {
		t.Error("Expected error for size 1x1")
	}

	// Test valid gameplay sprite size
	err = L.DoString(`rf.newSpriteFrames("gameplay", 16, 16, false)`)
	if err != nil {
		t.Errorf("Should allow 16x16 gameplay sprite: %v", err)
	}

	// Test invalid gameplay sprite size (too large)
	err = L.DoString(`rf.newSpriteFrames("invalid2", 64, 64, false)`)
	if err == nil {
		t.Error("Expected error for 64x64 gameplay sprite")
	}

	// Test valid UI sprite size
	err = L.DoString(`rf.newSpriteFrames("ui", 64, 64, true)`)
	if err != nil {
		t.Errorf("Should allow 64x64 UI sprite: %v", err)
	}

	// Test invalid UI sprite size (not divisible by 2)
	err = L.DoString(`rf.newSpriteFrames("invalid3", 3, 3, true)`)
	if err == nil {
		t.Error("Expected error for 3x3 UI sprite (not divisible by 2)")
	}
}

func TestAddSpriteFrame(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	colorByIndex := func(i int) (rgba [4]uint8) { return [4]uint8{255, 255, 255, 255} }
	setPalette := func(name string) {}
	spritesMap := make(cartio.SpriteMap)

	Register(L, r, colorByIndex, setPalette, make(cartio.SFXMap), make(cartio.MusicMap), spritesMap, nil, nil)

	// Create frames sprite and add a frame
	err := L.DoString(`
		rf.newSpriteFrames("test", 4, 4)
		
		-- Create frame pixels table
		local pixels = {}
		for y = 1, 4 do
			pixels[y] = {}
			for x = 1, 4 do
				pixels[y][x] = (y - 1) * 4 + (x - 1)
			end
		end
		
		rf.addSpriteFrame("test", "frame1", pixels)
	`)

	if err != nil {
		t.Fatalf("Lua error: %v", err)
	}

	// Verify frame was added
	sprite, ok := spritesMap["test"]
	if !ok {
		t.Fatal("Sprite should exist")
	}
	if len(sprite.Frames) != 1 {
		t.Errorf("Expected 1 frame, got %d", len(sprite.Frames))
	}
	if sprite.Frames[0].Name != "frame1" {
		t.Errorf("Expected frame name 'frame1', got '%s'", sprite.Frames[0].Name)
	}
	if len(sprite.Frames[0].Pixels) != 4 {
		t.Errorf("Expected 4 rows, got %d", len(sprite.Frames[0].Pixels))
	}

	// Test duplicate frame name
	err = L.DoString(`
		local pixels = {}
		for y = 1, 4 do
			pixels[y] = {}
			for x = 1, 4 do
				pixels[y][x] = 0
			end
		end
		rf.addSpriteFrame("test", "frame1", pixels)
	`)
	if err == nil {
		t.Error("Expected error for duplicate frame name")
	}

	// Test invalid frame name
	err = L.DoString(`
		local pixels = {}
		for y = 1, 4 do
			pixels[y] = {}
			for x = 1, 4 do
				pixels[y][x] = 0
			end
		end
		rf.addSpriteFrame("test", "1invalid", pixels)
	`)
	if err == nil {
		t.Error("Expected error for invalid frame name starting with number")
	}
}

func TestAddSpriteAnimation(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	colorByIndex := func(i int) (rgba [4]uint8) { return [4]uint8{255, 255, 255, 255} }
	setPalette := func(name string) {}
	spritesMap := make(cartio.SpriteMap)

	Register(L, r, colorByIndex, setPalette, make(cartio.SFXMap), make(cartio.MusicMap), spritesMap, nil, nil)

	// Create animation sprite (type="animation")
	err := L.DoString(`
		-- First create as frames sprite, then we'll add animation
		-- Actually, we need to manually set type to animation
		-- For now, let's test with a sprite that has type animation
	`)
	
	// Create sprite with frames first
	err = L.DoString(`
		rf.newSpriteFrames("anim_sprite", 4, 4)
		
		-- Add frames
		local pixels1 = {}
		for y = 1, 4 do
			pixels1[y] = {}
			for x = 1, 4 do
				pixels1[y][x] = 1
			end
		end
		rf.addSpriteFrame("anim_sprite", "frame1", pixels1)
		
		local pixels2 = {}
		for y = 1, 4 do
			pixels2[y] = {}
			for x = 1, 4 do
				pixels2[y][x] = 2
			end
		end
		rf.addSpriteFrame("anim_sprite", "frame2", pixels2)
	`)
	
	if err != nil {
		t.Fatalf("Failed to create frames: %v", err)
	}
	
	// Manually set sprite type to animation for testing
	sprite := spritesMap["anim_sprite"]
	sprite.Type = cartio.SpriteTypeAnimation
	spritesMap["anim_sprite"] = sprite

	// Now add animation
	err = L.DoString(`
		local frameRefs = {"frame1", "frame2"}
		rf.addSpriteAnimation("anim_sprite", "walk", frameRefs, 1.0, true, "forward")
	`)

	if err != nil {
		t.Fatalf("Lua error: %v", err)
	}

	// Verify animation was added
	sprite = spritesMap["anim_sprite"]
	if len(sprite.Animations) != 1 {
		t.Errorf("Expected 1 animation, got %d", len(sprite.Animations))
	}
	if sprite.Animations[0].Name != "walk" {
		t.Errorf("Expected animation name 'walk', got '%s'", sprite.Animations[0].Name)
	}
	if len(sprite.Animations[0].FrameRefs) != 2 {
		t.Errorf("Expected 2 frame refs, got %d", len(sprite.Animations[0].FrameRefs))
	}

	// Test duplicate animation name
	err = L.DoString(`
		local frameRefs = {"frame1"}
		rf.addSpriteAnimation("anim_sprite", "walk", frameRefs)
	`)
	if err == nil {
		t.Error("Expected error for duplicate animation name")
	}

	// Test invalid frame reference
	err = L.DoString(`
		local frameRefs = {"nonexistent"}
		rf.addSpriteAnimation("anim_sprite", "invalid", frameRefs)
	`)
	if err == nil {
		t.Error("Expected error for invalid frame reference")
	}

	// Test invalid loop type
	err = L.DoString(`
		local frameRefs = {"frame1"}
		rf.addSpriteAnimation("anim_sprite", "invalid_loop", frameRefs, 1.0, true, "invalid")
	`)
	if err == nil {
		t.Error("Expected error for invalid loop type")
	}
}

func TestSprWithFrames(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	colorByIndex := func(i int) (rgba [4]uint8) { return [4]uint8{255, 255, 255, 255} }
	setPalette := func(name string) {}
	spritesMap := make(cartio.SpriteMap)

	Register(L, r, colorByIndex, setPalette, make(cartio.SFXMap), make(cartio.MusicMap), spritesMap, nil, nil)

	// Create frames sprite and add frames
	err := L.DoString(`
		rf.newSpriteFrames("player", 4, 4)
		
		local pixels1 = {}
		for y = 1, 4 do
			pixels1[y] = {}
			for x = 1, 4 do
				pixels1[y][x] = 1
			end
		end
		rf.addSpriteFrame("player", "left", pixels1)
		
		local pixels2 = {}
		for y = 1, 4 do
			pixels2[y] = {}
			for x = 1, 4 do
				pixels2[y][x] = 2
			end
		end
		rf.addSpriteFrame("player", "right", pixels2)
	`)

	if err != nil {
		t.Fatalf("Failed to create frames sprite: %v", err)
	}

	// Test drawing with frame name
	err = L.DoString(`rf.spr("player", 10, 10, "left")`)
	if err != nil {
		t.Errorf("Failed to draw frame: %v", err)
	}

	// Test drawing without frame name (should error)
	err = L.DoString(`rf.spr("player", 10, 10)`)
	if err == nil {
		t.Error("Expected error when drawing frames sprite without frame name")
	}

	// Test invalid frame name
	err = L.DoString(`rf.spr("player", 10, 10, "nonexistent")`)
	if err == nil {
		t.Error("Expected error for nonexistent frame")
	}
}

func TestPlayAnimation(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	colorByIndex := func(i int) (rgba [4]uint8) { return [4]uint8{255, 255, 255, 255} }
	setPalette := func(name string) {}
	spritesMap := make(cartio.SpriteMap)

	Register(L, r, colorByIndex, setPalette, make(cartio.SFXMap), make(cartio.MusicMap), spritesMap, nil, nil)

	// Create animation sprite
	err := L.DoString(`
		rf.newSpriteFrames("anim_sprite", 4, 4)
		
		local pixels1 = {}
		for y = 1, 4 do
			pixels1[y] = {}
			for x = 1, 4 do
				pixels1[y][x] = 1
			end
		end
		rf.addSpriteFrame("anim_sprite", "frame1", pixels1)
		
		local pixels2 = {}
		for y = 1, 4 do
			pixels2[y] = {}
			for x = 1, 4 do
				pixels2[y][x] = 2
			end
		end
		rf.addSpriteFrame("anim_sprite", "frame2", pixels2)
	`)
	
	// Set type to animation
	sprite := spritesMap["anim_sprite"]
	sprite.Type = cartio.SpriteTypeAnimation
	spritesMap["anim_sprite"] = sprite

	// Add animation
	err = L.DoString(`
		local frameRefs = {"frame1", "frame2"}
		rf.addSpriteAnimation("anim_sprite", "walk", frameRefs, 1.0, true, "forward")
	`)

	// Test playing animation
	err = L.DoString(`rf.playAnimation("anim_sprite", "walk")`)
	if err != nil {
		t.Errorf("Failed to play animation: %v", err)
	}

	// Test playing nonexistent animation
	err = L.DoString(`rf.playAnimation("anim_sprite", "nonexistent")`)
	if err == nil {
		t.Error("Expected error for nonexistent animation")
	}

	// Test playing on non-animation sprite
	err = L.DoString(`rf.playAnimation("nonexistent", "walk")`)
	if err == nil {
		t.Error("Expected error for nonexistent sprite")
	}
}

func TestAnimationControl(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	r := rendersoft.New(480, 270)
	colorByIndex := func(i int) (rgba [4]uint8) { return [4]uint8{255, 255, 255, 255} }
	setPalette := func(name string) {}
	spritesMap := make(cartio.SpriteMap)

	Register(L, r, colorByIndex, setPalette, make(cartio.SFXMap), make(cartio.MusicMap), spritesMap, nil, nil)

	// Create animation sprite
	err := L.DoString(`
		rf.newSpriteFrames("test", 4, 4)
		local pixels = {}
		for y = 1, 4 do
			pixels[y] = {}
			for x = 1, 4 do
				pixels[y][x] = 1
			end
		end
		rf.addSpriteFrame("test", "frame1", pixels)
	`)
	
	sprite := spritesMap["test"]
	sprite.Type = cartio.SpriteTypeAnimation
	spritesMap["test"] = sprite

	err = L.DoString(`
		local frameRefs = {"frame1"}
		rf.addSpriteAnimation("test", "anim", frameRefs)
		rf.playAnimation("test", "anim")
	`)

	// Test pause
	err = L.DoString(`rf.pauseAnimation("test")`)
	if err != nil {
		t.Errorf("Failed to pause: %v", err)
	}

	// Test resume
	err = L.DoString(`rf.resumeAnimation("test")`)
	if err != nil {
		t.Errorf("Failed to resume: %v", err)
	}

	// Test stop
	err = L.DoString(`rf.stopAnimation("test")`)
	if err != nil {
		t.Errorf("Failed to stop: %v", err)
	}

	// Test set speed
	err = L.DoString(`rf.setAnimationSpeed("test", 2.0)`)
	if err != nil {
		t.Errorf("Failed to set speed: %v", err)
	}

	// Test set frame
	err = L.DoString(`rf.setAnimationFrame("test", 0)`)
	if err != nil {
		t.Errorf("Failed to set frame: %v", err)
	}

	// Test get frame
	err = L.DoString(`
		local frame = rf.getAnimationFrame("test")
		if frame ~= 0 then
			error("Expected frame 0, got " .. frame)
		end
	`)
	if err != nil {
		t.Errorf("Failed to get frame: %v", err)
	}

	// Test get frame for sprite with no animation
	err = L.DoString(`
		local frame = rf.getAnimationFrame("nonexistent")
		if frame ~= -1 then
			error("Expected -1 for nonexistent sprite, got " .. frame)
		end
	`)
	if err != nil {
		t.Errorf("Failed to handle nonexistent sprite: %v", err)
	}
}

