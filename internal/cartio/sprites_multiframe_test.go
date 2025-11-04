package cartio

import (
	"strings"
	"testing"
)

func TestValidateFrameName(t *testing.T) {
	tests := []struct {
		name    string
		frameName string
		wantErr bool
	}{
		{"valid simple", "idle", false},
		{"valid with underscore", "walk_left", false},
		{"valid with hyphen", "walk-left", false},
		{"valid with numbers", "frame1", false},
		{"valid mixed", "frame_1-left", false},
		{"valid starts with underscore", "_private", false},
		{"valid max length", strings.Repeat("a", 64), false},
		
		{"invalid empty", "", true},
		{"invalid too long", string(make([]byte, 65)), true},
		{"invalid starts with number", "1frame", true},
		{"invalid starts with hyphen", "-frame", true},
		{"invalid contains space", "frame name", true},
		{"invalid contains special char", "frame@name", true},
		{"invalid contains dot", "frame.name", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFrameName(tt.frameName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFrameName(%q) error = %v, wantErr %v", tt.frameName, err, tt.wantErr)
			}
		})
	}
}

func TestValidateSpriteData_Static(t *testing.T) {
	// Valid static sprite
	sprite := &SpriteData{
		Width:  16,
		Height: 16,
		Type:   SpriteTypeStatic,
		Pixels: make([][]int, 16),
		IsUI:   false,
	}
	for i := range sprite.Pixels {
		sprite.Pixels[i] = make([]int, 16)
	}
	
	if err := ValidateSpriteData(sprite, "test"); err != nil {
		t.Errorf("valid static sprite should pass validation: %v", err)
	}
	
	// Invalid: missing pixels
	sprite2 := &SpriteData{
		Width:  16,
		Height: 16,
		Type:   SpriteTypeStatic,
		Pixels: nil,
		IsUI:   false,
	}
	if err := ValidateSpriteData(sprite2, "test"); err == nil {
		t.Error("static sprite without pixels should fail validation")
	}
	
	// Invalid: wrong pixel dimensions
	sprite3 := &SpriteData{
		Width:  16,
		Height: 16,
		Type:   SpriteTypeStatic,
		Pixels: make([][]int, 15),
		IsUI:   false,
	}
	if err := ValidateSpriteData(sprite3, "test"); err == nil {
		t.Error("static sprite with wrong pixel height should fail validation")
	}
}

func TestValidateSpriteData_Frames(t *testing.T) {
	// Valid frames sprite
	sprite := &SpriteData{
		Width:  16,
		Height: 16,
		Type:   SpriteTypeFrames,
		Frames: []SpriteFrame{
			{
				Name: "idle",
				Pixels: make([][]int, 16),
			},
			{
				Name: "walk",
				Pixels: make([][]int, 16),
			},
		},
		IsUI: false,
	}
	for i := range sprite.Frames {
		for j := range sprite.Frames[i].Pixels {
			sprite.Frames[i].Pixels[j] = make([]int, 16)
		}
	}
	
	if err := ValidateSpriteData(sprite, "test"); err != nil {
		t.Errorf("valid frames sprite should pass validation: %v", err)
	}
	
	// Invalid: no frames
	sprite2 := &SpriteData{
		Width:  16,
		Height: 16,
		Type:   SpriteTypeFrames,
		Frames: []SpriteFrame{},
		IsUI:   false,
	}
	if err := ValidateSpriteData(sprite2, "test"); err == nil {
		t.Error("frames sprite without frames should fail validation")
	}
	
	// Invalid: duplicate frame names
	sprite3 := &SpriteData{
		Width:  16,
		Height: 16,
		Type:   SpriteTypeFrames,
		Frames: []SpriteFrame{
			{
				Name: "idle",
				Pixels: make([][]int, 16),
			},
			{
				Name: "idle", // Duplicate
				Pixels: make([][]int, 16),
			},
		},
		IsUI: false,
	}
	for i := range sprite3.Frames {
		for j := range sprite3.Frames[i].Pixels {
			sprite3.Frames[i].Pixels[j] = make([]int, 16)
		}
	}
	if err := ValidateSpriteData(sprite3, "test"); err == nil {
		t.Error("frames sprite with duplicate names should fail validation")
	}
	
	// Invalid: invalid frame name
	sprite4 := &SpriteData{
		Width:  16,
		Height: 16,
		Type:   SpriteTypeFrames,
		Frames: []SpriteFrame{
			{
				Name: "1invalid", // Starts with number
				Pixels: make([][]int, 16),
			},
		},
		IsUI: false,
	}
	for i := range sprite4.Frames {
		for j := range sprite4.Frames[i].Pixels {
			sprite4.Frames[i].Pixels[j] = make([]int, 16)
		}
	}
	if err := ValidateSpriteData(sprite4, "test"); err == nil {
		t.Error("frames sprite with invalid frame name should fail validation")
	}
}

func TestValidateSpriteData_Animation(t *testing.T) {
	// Valid animation sprite
	sprite := &SpriteData{
		Width:  16,
		Height: 16,
		Type:   SpriteTypeAnimation,
		Frames: []SpriteFrame{
			{
				Name: "frame1",
				Pixels: make([][]int, 16),
			},
			{
				Name: "frame2",
				Pixels: make([][]int, 16),
			},
		},
		Animations: []AnimationSequence{
			{
				Name:      "walk",
				FrameRefs: []string{"frame1", "frame2"},
				Speed:     1.0,
				Loop:      true,
				LoopType:  "forward",
			},
		},
		IsUI: false,
	}
	for i := range sprite.Frames {
		for j := range sprite.Frames[i].Pixels {
			sprite.Frames[i].Pixels[j] = make([]int, 16)
		}
	}
	
	if err := ValidateSpriteData(sprite, "test"); err != nil {
		t.Errorf("valid animation sprite should pass validation: %v", err)
	}
	
	// Invalid: no animations
	sprite2 := &SpriteData{
		Width:  16,
		Height: 16,
		Type:   SpriteTypeAnimation,
		Frames: []SpriteFrame{
			{Name: "frame1", Pixels: make([][]int, 16)},
		},
		Animations: []AnimationSequence{},
		IsUI: false,
	}
	for i := range sprite2.Frames {
		for j := range sprite2.Frames[i].Pixels {
			sprite2.Frames[i].Pixels[j] = make([]int, 16)
		}
	}
	if err := ValidateSpriteData(sprite2, "test"); err == nil {
		t.Error("animation sprite without animations should fail validation")
	}
	
	// Invalid: invalid loop type
	sprite3 := &SpriteData{
		Width:  16,
		Height: 16,
		Type:   SpriteTypeAnimation,
		Frames: []SpriteFrame{
			{Name: "frame1", Pixels: make([][]int, 16)},
		},
		Animations: []AnimationSequence{
			{
				Name:      "walk",
				FrameRefs: []string{"frame1"},
				LoopType:  "invalid",
			},
		},
		IsUI: false,
	}
	for i := range sprite3.Frames {
		for j := range sprite3.Frames[i].Pixels {
			sprite3.Frames[i].Pixels[j] = make([]int, 16)
		}
	}
	if err := ValidateSpriteData(sprite3, "test"); err == nil {
		t.Error("animation sprite with invalid loop type should fail validation")
	}
	
	// Invalid: references unknown frame
	sprite4 := &SpriteData{
		Width:  16,
		Height: 16,
		Type:   SpriteTypeAnimation,
		Frames: []SpriteFrame{
			{Name: "frame1", Pixels: make([][]int, 16)},
		},
		Animations: []AnimationSequence{
			{
				Name:      "walk",
				FrameRefs: []string{"frame1", "unknown"},
			},
		},
		IsUI: false,
	}
	for i := range sprite4.Frames {
		for j := range sprite4.Frames[i].Pixels {
			sprite4.Frames[i].Pixels[j] = make([]int, 16)
		}
	}
	if err := ValidateSpriteData(sprite4, "test"); err == nil {
		t.Error("animation sprite with unknown frame reference should fail validation")
	}
}

func TestGetFramePixels(t *testing.T) {
	// Static sprite
	staticSprite := &SpriteData{
		Type:   SpriteTypeStatic,
		Pixels: [][]int{{1, 2}, {3, 4}},
	}
	pixels, err := staticSprite.GetFramePixels("")
	if err != nil {
		t.Errorf("GetFramePixels should work for static sprite: %v", err)
	}
	if len(pixels) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(pixels))
	}
	
	// Frames sprite
	framesSprite := &SpriteData{
		Type: SpriteTypeFrames,
		Frames: []SpriteFrame{
			{Name: "idle", Pixels: [][]int{{5, 6}}},
			{Name: "walk", Pixels: [][]int{{7, 8}}},
		},
	}
	pixels, err = framesSprite.GetFramePixels("idle")
	if err != nil {
		t.Errorf("GetFramePixels should find frame: %v", err)
	}
	if pixels[0][0] != 5 {
		t.Errorf("Expected pixel 5, got %d", pixels[0][0])
	}
	
	// Animation sprite
	animSprite := &SpriteData{
		Type: SpriteTypeAnimation,
		Frames: []SpriteFrame{
			{Name: "frame1", Pixels: [][]int{{9, 10}}},
		},
	}
	pixels, err = animSprite.GetFramePixels("frame1")
	if err != nil {
		t.Errorf("GetFramePixels should find frame in animation sprite: %v", err)
	}
	if pixels[0][0] != 9 {
		t.Errorf("Expected pixel 9, got %d", pixels[0][0])
	}
	
	// Frame not found
	_, err = framesSprite.GetFramePixels("nonexistent")
	if err == nil {
		t.Error("GetFramePixels should error for nonexistent frame")
	}
	
	// Unknown sprite type
	unknownSprite := &SpriteData{
		Type: SpriteType("unknown"),
	}
	_, err = unknownSprite.GetFramePixels("frame1")
	if err == nil {
		t.Error("GetFramePixels should error for unknown sprite type")
	}
}

func TestHasFrame(t *testing.T) {
	// Static sprite
	staticSprite := &SpriteData{Type: SpriteTypeStatic}
	if !staticSprite.HasFrame("") {
		t.Error("Static sprite should have default frame")
	}
	if !staticSprite.HasFrame("default") {
		t.Error("Static sprite should have default frame")
	}
	
	// Frames sprite
	framesSprite := &SpriteData{
		Type: SpriteTypeFrames,
		Frames: []SpriteFrame{
			{Name: "idle"},
			{Name: "walk"},
		},
	}
	if !framesSprite.HasFrame("idle") {
		t.Error("Should have idle frame")
	}
	if framesSprite.HasFrame("nonexistent") {
		t.Error("Should not have nonexistent frame")
	}
}

func TestGetAnimation(t *testing.T) {
	sprite := &SpriteData{
		Type: SpriteTypeAnimation,
		Animations: []AnimationSequence{
			{Name: "walk", FrameRefs: []string{"frame1"}},
			{Name: "run", FrameRefs: []string{"frame2"}},
		},
	}
	
	anim, err := sprite.GetAnimation("walk")
	if err != nil {
		t.Errorf("GetAnimation should find animation: %v", err)
	}
	if anim.Name != "walk" {
		t.Errorf("Expected walk animation, got %s", anim.Name)
	}
	
	_, err = sprite.GetAnimation("nonexistent")
	if err == nil {
		t.Error("GetAnimation should error for nonexistent animation")
	}
	
	// Non-animation sprite
	staticSprite := &SpriteData{Type: SpriteTypeStatic}
	_, err = staticSprite.GetAnimation("walk")
	if err == nil {
		t.Error("GetAnimation should error for non-animation sprite")
	}
}

func TestNormalizeSpriteData(t *testing.T) {
	// Test default type
	sprite := &SpriteData{Width: 16, Height: 16}
	NormalizeSpriteData(sprite)
	if sprite.Type != SpriteTypeStatic {
		t.Errorf("Expected default type static, got %s", sprite.Type)
	}
	
	// Test animation defaults
	sprite2 := &SpriteData{
		Type: SpriteTypeAnimation,
		Animations: []AnimationSequence{
			{
				Name:      "walk",
				FrameRefs: []string{" frame1 ", " frame2 "}, // With whitespace
				LoopType:  "", // Empty
				Speed:     0,  // Invalid
			},
		},
	}
	NormalizeSpriteData(sprite2)
	if sprite2.Animations[0].LoopType != "forward" {
		t.Errorf("Expected default loopType forward, got %s", sprite2.Animations[0].LoopType)
	}
	if sprite2.Animations[0].Speed != 1.0 {
		t.Errorf("Expected default speed 1.0, got %f", sprite2.Animations[0].Speed)
	}
	// Check whitespace trimmed
	if sprite2.Animations[0].FrameRefs[0] != "frame1" {
		t.Errorf("Expected trimmed frameRef, got %q", sprite2.Animations[0].FrameRefs[0])
	}
	
	// Test pingpong loop type
	sprite3 := &SpriteData{
		Type: SpriteTypeAnimation,
		Animations: []AnimationSequence{
			{
				Name:      "walk",
				FrameRefs: []string{"frame1"},
				LoopType:  "pingpong",
			},
		},
	}
	NormalizeSpriteData(sprite3)
	if sprite3.Animations[0].LoopType != "pingpong" {
		t.Errorf("Expected pingpong loopType, got %s", sprite3.Animations[0].LoopType)
	}
	
	// Test reverse loop type
	sprite4 := &SpriteData{
		Type: SpriteTypeAnimation,
		Animations: []AnimationSequence{
			{
				Name:      "walk",
				FrameRefs: []string{"frame1"},
				LoopType:  "reverse",
			},
		},
	}
	NormalizeSpriteData(sprite4)
	if sprite4.Animations[0].LoopType != "reverse" {
		t.Errorf("Expected reverse loopType, got %s", sprite4.Animations[0].LoopType)
	}
}

func TestValidateSpriteData_BackwardCompatibility(t *testing.T) {
	// Sprite without type should default to static
	sprite := &SpriteData{
		Width:  16,
		Height: 16,
		Pixels: make([][]int, 16),
		IsUI:   false,
	}
	for i := range sprite.Pixels {
		sprite.Pixels[i] = make([]int, 16)
	}
	
	if err := ValidateSpriteData(sprite, "test"); err != nil {
		t.Errorf("Old sprite format should be valid (defaults to static): %v", err)
	}
	if sprite.Type != SpriteTypeStatic {
		t.Errorf("Expected default type static, got %s", sprite.Type)
	}
}

func TestValidateSpriteData_InvalidType(t *testing.T) {
	sprite := &SpriteData{
		Width:  16,
		Height: 16,
		Type:   SpriteType("invalid"),
		Pixels: make([][]int, 16),
		IsUI:   false,
	}
	for i := range sprite.Pixels {
		sprite.Pixels[i] = make([]int, 16)
	}
	
	if err := ValidateSpriteData(sprite, "test"); err == nil {
		t.Error("Sprite with invalid type should fail validation")
	}
}

func TestValidateSpriteData_AnimationEdgeCases(t *testing.T) {
	// Test empty frame refs
	sprite := &SpriteData{
		Width:  16,
		Height: 16,
		Type:   SpriteTypeAnimation,
		Frames: []SpriteFrame{
			{Name: "frame1", Pixels: make([][]int, 16)},
		},
		Animations: []AnimationSequence{
			{
				Name:      "walk",
				FrameRefs: []string{}, // Empty
			},
		},
		IsUI: false,
	}
	for i := range sprite.Frames {
		for j := range sprite.Frames[i].Pixels {
			sprite.Frames[i].Pixels[j] = make([]int, 16)
		}
	}
	if err := ValidateSpriteData(sprite, "test"); err == nil {
		t.Error("Animation with empty frame refs should fail validation")
	}
	
	// Test duplicate animation names
	sprite2 := &SpriteData{
		Width:  16,
		Height: 16,
		Type:   SpriteTypeAnimation,
		Frames: []SpriteFrame{
			{Name: "frame1", Pixels: make([][]int, 16)},
		},
		Animations: []AnimationSequence{
			{
				Name:      "walk",
				FrameRefs: []string{"frame1"},
			},
			{
				Name:      "walk", // Duplicate
				FrameRefs: []string{"frame1"},
			},
		},
		IsUI: false,
	}
	for i := range sprite2.Frames {
		for j := range sprite2.Frames[i].Pixels {
			sprite2.Frames[i].Pixels[j] = make([]int, 16)
		}
	}
	if err := ValidateSpriteData(sprite2, "test"); err == nil {
		t.Error("Animation with duplicate names should fail validation")
	}
	
	// Test invalid animation name
	sprite3 := &SpriteData{
		Width:  16,
		Height: 16,
		Type:   SpriteTypeAnimation,
		Frames: []SpriteFrame{
			{Name: "frame1", Pixels: make([][]int, 16)},
		},
		Animations: []AnimationSequence{
			{
				Name:      "1invalid", // Invalid name
				FrameRefs: []string{"frame1"},
			},
		},
		IsUI: false,
	}
	for i := range sprite3.Frames {
		for j := range sprite3.Frames[i].Pixels {
			sprite3.Frames[i].Pixels[j] = make([]int, 16)
		}
	}
	if err := ValidateSpriteData(sprite3, "test"); err == nil {
		t.Error("Animation with invalid name should fail validation")
	}
}

func TestValidateSpriteData_FramesEdgeCases(t *testing.T) {
	// Test frame with wrong pixel dimensions
	sprite := &SpriteData{
		Width:  16,
		Height: 16,
		Type:   SpriteTypeFrames,
		Frames: []SpriteFrame{
			{
				Name: "idle",
				Pixels: make([][]int, 15), // Wrong height
			},
		},
		IsUI: false,
	}
	for i := range sprite.Frames {
		for j := range sprite.Frames[i].Pixels {
			sprite.Frames[i].Pixels[j] = make([]int, 16)
		}
	}
	if err := ValidateSpriteData(sprite, "test"); err == nil {
		t.Error("Frame with wrong pixel height should fail validation")
	}
	
	// Test frame with wrong row width
	sprite2 := &SpriteData{
		Width:  16,
		Height: 16,
		Type:   SpriteTypeFrames,
		Frames: []SpriteFrame{
			{
				Name: "idle",
				Pixels: make([][]int, 16),
			},
		},
		IsUI: false,
	}
	for i := range sprite2.Frames {
		sprite2.Frames[i].Pixels[i] = make([]int, 15) // Wrong width
	}
	if err := ValidateSpriteData(sprite2, "test"); err == nil {
		t.Error("Frame with wrong row width should fail validation")
	}
	
	// Test frame with empty pixels
	sprite3 := &SpriteData{
		Width:  16,
		Height: 16,
		Type:   SpriteTypeFrames,
		Frames: []SpriteFrame{
			{
				Name: "idle",
				Pixels: nil, // Empty
			},
		},
		IsUI: false,
	}
	if err := ValidateSpriteData(sprite3, "test"); err == nil {
		t.Error("Frame with empty pixels should fail validation")
	}
	
	// Test animation sprite with frame that has wrong pixel dimensions
	sprite4 := &SpriteData{
		Width:  16,
		Height: 16,
		Type:   SpriteTypeAnimation,
		Frames: []SpriteFrame{
			{
				Name: "frame1",
				Pixels: make([][]int, 15), // Wrong height
			},
		},
		Animations: []AnimationSequence{
			{
				Name:      "walk",
				FrameRefs: []string{"frame1"},
			},
		},
		IsUI: false,
	}
	for i := range sprite4.Frames {
		for j := range sprite4.Frames[i].Pixels {
			if j < len(sprite4.Frames[i].Pixels) {
				sprite4.Frames[i].Pixels[j] = make([]int, 16)
			}
		}
	}
	if err := ValidateSpriteData(sprite4, "test"); err == nil {
		t.Error("Animation sprite with frame wrong pixel height should fail validation")
	}
	
	// Test animation sprite with frame that has wrong row width
	sprite5 := &SpriteData{
		Width:  16,
		Height: 16,
		Type:   SpriteTypeAnimation,
		Frames: []SpriteFrame{
			{
				Name: "frame1",
				Pixels: make([][]int, 16),
			},
		},
		Animations: []AnimationSequence{
			{
				Name:      "walk",
				FrameRefs: []string{"frame1"},
			},
		},
		IsUI: false,
	}
	for i := range sprite5.Frames {
		sprite5.Frames[i].Pixels[0] = make([]int, 15) // Wrong width on first row
		for j := 1; j < len(sprite5.Frames[i].Pixels); j++ {
			sprite5.Frames[i].Pixels[j] = make([]int, 16)
		}
	}
	if err := ValidateSpriteData(sprite5, "test"); err == nil {
		t.Error("Animation sprite with frame wrong row width should fail validation")
	}
	
	// Test animation with all valid loop types
	sprite6 := &SpriteData{
		Width:  16,
		Height: 16,
		Type:   SpriteTypeAnimation,
		Frames: []SpriteFrame{
			{Name: "frame1", Pixels: make([][]int, 16)},
			{Name: "frame2", Pixels: make([][]int, 16)},
		},
		Animations: []AnimationSequence{
			{Name: "forward", FrameRefs: []string{"frame1", "frame2"}, LoopType: "forward"},
			{Name: "reverse", FrameRefs: []string{"frame1", "frame2"}, LoopType: "reverse"},
			{Name: "pingpong", FrameRefs: []string{"frame1", "frame2"}, LoopType: "pingpong"},
		},
		IsUI: false,
	}
	for i := range sprite6.Frames {
		for j := range sprite6.Frames[i].Pixels {
			sprite6.Frames[i].Pixels[j] = make([]int, 16)
		}
	}
	if err := ValidateSpriteData(sprite6, "test"); err != nil {
		t.Errorf("Valid animation sprite with all loop types should pass: %v", err)
	}
}

