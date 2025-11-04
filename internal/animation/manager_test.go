package animation

import (
	"testing"

	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
)

func TestNewAnimationState(t *testing.T) {
	state := NewAnimationState("test_sprite")
	if state == nil {
		t.Fatal("NewAnimationState should not return nil")
	}
	if state.SpriteName != "test_sprite" {
		t.Errorf("Expected sprite name 'test_sprite', got %s", state.SpriteName)
	}
	if state.Playing {
		t.Error("New animation state should not be playing")
	}
	if state.Speed != 1.0 {
		t.Errorf("Expected default speed 1.0, got %f", state.Speed)
	}
	if state.Direction != 1 {
		t.Errorf("Expected default direction 1, got %d", state.Direction)
	}
}

func TestAnimationState_Reset(t *testing.T) {
	state := NewAnimationState("test")
	state.Playing = true
	state.Paused = true
	state.AnimationName = "walk"
	state.CurrentFrame = 5
	state.ElapsedTime = 100
	
	state.Reset()
	
	if state.Playing {
		t.Error("Reset should clear playing flag")
	}
	if state.Paused {
		t.Error("Reset should clear paused flag")
	}
	if state.AnimationName != "" {
		t.Error("Reset should clear animation name")
	}
	if state.CurrentFrame != 0 {
		t.Error("Reset should reset frame to 0")
	}
	if state.ElapsedTime != 0 {
		t.Error("Reset should reset elapsed time")
	}
}

func TestPlayAnimation(t *testing.T) {
	spriteData := &cartio.SpriteData{
		Type: cartio.SpriteTypeAnimation,
		Frames: []cartio.SpriteFrame{
			{Name: "frame1", Pixels: [][]int{{1}}, Duration: 100},
			{Name: "frame2", Pixels: [][]int{{2}}, Duration: 100},
		},
		Animations: []cartio.AnimationSequence{
			{
				Name:      "walk",
				FrameRefs: []string{"frame1", "frame2"},
				Speed:     1.0,
				Loop:      true,
				LoopType:  "forward",
			},
		},
	}
	
	state := NewAnimationState("test")
	err := PlayAnimation(state, spriteData, "walk")
	if err != nil {
		t.Errorf("PlayAnimation should succeed: %v", err)
	}
	if !state.Playing {
		t.Error("Animation should be playing after PlayAnimation")
	}
	if state.AnimationName != "walk" {
		t.Errorf("Expected animation name 'walk', got %s", state.AnimationName)
	}
	if state.CurrentFrame != 0 {
		t.Errorf("Expected starting frame 0, got %d", state.CurrentFrame)
	}
	if state.Speed != 1.0 {
		t.Errorf("Expected speed 1.0, got %f", state.Speed)
	}
	
	// Test with reverse animation
	spriteData2 := &cartio.SpriteData{
		Type: cartio.SpriteTypeAnimation,
		Frames: []cartio.SpriteFrame{
			{Name: "frame1", Pixels: [][]int{{1}}, Duration: 100},
			{Name: "frame2", Pixels: [][]int{{2}}, Duration: 100},
		},
		Animations: []cartio.AnimationSequence{
			{
				Name:      "reverse",
				FrameRefs: []string{"frame1", "frame2"},
				LoopType:  "reverse",
			},
		},
	}
	state2 := NewAnimationState("test")
	err = PlayAnimation(state2, spriteData2, "reverse")
	if err != nil {
		t.Errorf("PlayAnimation should succeed: %v", err)
	}
	if state2.Direction != -1 {
		t.Errorf("Expected direction -1 for reverse, got %d", state2.Direction)
	}
	if state2.CurrentFrame != 1 {
		t.Errorf("Expected starting at last frame (1) for reverse, got %d", state2.CurrentFrame)
	}
}

func TestPlayAnimation_Errors(t *testing.T) {
	state := NewAnimationState("test")
	
	// Test nil sprite data
	err := PlayAnimation(state, nil, "walk")
	if err == nil {
		t.Error("PlayAnimation should error with nil sprite data")
	}
	
	// Test non-animation sprite
	staticSprite := &cartio.SpriteData{Type: cartio.SpriteTypeStatic}
	err = PlayAnimation(state, staticSprite, "walk")
	if err == nil {
		t.Error("PlayAnimation should error with non-animation sprite")
	}
	
	// Test missing animation
	animSprite := &cartio.SpriteData{
		Type:       cartio.SpriteTypeAnimation,
		Animations: []cartio.AnimationSequence{},
	}
	err = PlayAnimation(state, animSprite, "nonexistent")
	if err == nil {
		t.Error("PlayAnimation should error with nonexistent animation")
	}
}

func TestPauseResumeAnimation(t *testing.T) {
	state := NewAnimationState("test")
	state.Playing = true
	
	PauseAnimation(state)
	if !state.Paused {
		t.Error("PauseAnimation should set paused flag")
	}
	if !state.Playing {
		t.Error("PauseAnimation should not clear playing flag")
	}
	
	ResumeAnimation(state)
	if state.Paused {
		t.Error("ResumeAnimation should clear paused flag")
	}
}

func TestStopAnimation(t *testing.T) {
	state := NewAnimationState("test")
	state.Playing = true
	state.Paused = true
	state.CurrentFrame = 5
	state.ElapsedTime = 100
	state.AnimationName = "walk"
	
	StopAnimation(state)
	
	if state.Playing {
		t.Error("StopAnimation should clear playing flag")
	}
	if state.Paused {
		t.Error("StopAnimation should clear paused flag")
	}
	if state.CurrentFrame != 0 {
		t.Error("StopAnimation should reset frame to 0")
	}
	if state.ElapsedTime != 0 {
		t.Error("StopAnimation should reset elapsed time")
	}
}

func TestSetAnimationSpeed(t *testing.T) {
	state := NewAnimationState("test")
	
	SetAnimationSpeed(state, 2.0)
	if state.Speed != 2.0 {
		t.Errorf("Expected speed 2.0, got %f", state.Speed)
	}
	
	SetAnimationSpeed(state, 0.5)
	if state.Speed != 0.5 {
		t.Errorf("Expected speed 0.5, got %f", state.Speed)
	}
	
	// Test invalid speed (should default to 1.0)
	SetAnimationSpeed(state, -1.0)
	if state.Speed != 1.0 {
		t.Errorf("Expected speed 1.0 for invalid speed, got %f", state.Speed)
	}
	
	SetAnimationSpeed(state, 0.0)
	if state.Speed != 1.0 {
		t.Errorf("Expected speed 1.0 for zero speed, got %f", state.Speed)
	}
}

func TestSetAnimationFrame(t *testing.T) {
	state := NewAnimationState("test")
	state.Sequence = &cartio.AnimationSequence{
		FrameRefs: []string{"frame1", "frame2", "frame3"},
	}
	
	err := SetAnimationFrame(state, 1)
	if err != nil {
		t.Errorf("SetAnimationFrame should succeed: %v", err)
	}
	if state.CurrentFrame != 1 {
		t.Errorf("Expected frame 1, got %d", state.CurrentFrame)
	}
	if state.ElapsedTime != 0 {
		t.Error("SetAnimationFrame should reset elapsed time")
	}
	
	// Test invalid frame index
	err = SetAnimationFrame(state, 10)
	if err == nil {
		t.Error("SetAnimationFrame should error with out-of-range index")
	}
	
	err = SetAnimationFrame(state, -1)
	if err == nil {
		t.Error("SetAnimationFrame should error with negative index")
	}
	
	// Test nil state
	err = SetAnimationFrame(nil, 0)
	if err == nil {
		t.Error("SetAnimationFrame should error with nil state")
	}
	
	// Test nil sequence
	state2 := NewAnimationState("test")
	err = SetAnimationFrame(state2, 0)
	if err == nil {
		t.Error("SetAnimationFrame should error with nil sequence")
	}
}

func TestGetAnimationFrame(t *testing.T) {
	state := NewAnimationState("test")
	state.Sequence = &cartio.AnimationSequence{
		FrameRefs: []string{"frame1", "frame2"},
	}
	state.CurrentFrame = 1
	
	frame := GetAnimationFrame(state)
	if frame != 1 {
		t.Errorf("Expected frame 1, got %d", frame)
	}
	
	// Test nil state
	frame = GetAnimationFrame(nil)
	if frame != -1 {
		t.Errorf("Expected -1 for nil state, got %d", frame)
	}
	
	// Test nil sequence
	state2 := NewAnimationState("test")
	frame = GetAnimationFrame(state2)
	if frame != -1 {
		t.Errorf("Expected -1 for nil sequence, got %d", frame)
	}
}

func TestGetCurrentFrameName(t *testing.T) {
	state := NewAnimationState("test")
	state.Sequence = &cartio.AnimationSequence{
		FrameRefs: []string{"frame1", "frame2"},
	}
	state.CurrentFrame = 0
	
	name := GetCurrentFrameName(state)
	if name != "frame1" {
		t.Errorf("Expected 'frame1', got %s", name)
	}
	
	state.CurrentFrame = 1
	name = GetCurrentFrameName(state)
	if name != "frame2" {
		t.Errorf("Expected 'frame2', got %s", name)
	}
	
	// Test nil state
	name = GetCurrentFrameName(nil)
	if name != "" {
		t.Errorf("Expected empty string for nil state, got %s", name)
	}
	
	// Test nil sequence
	state2 := NewAnimationState("test")
	name = GetCurrentFrameName(state2)
	if name != "" {
		t.Errorf("Expected empty string for nil sequence, got %s", name)
	}
}

func TestUpdateAnimationState_NotPlaying(t *testing.T) {
	state := NewAnimationState("test")
	state.Playing = false
	
	spriteData := &cartio.SpriteData{
		Type: cartio.SpriteTypeAnimation,
	}
	
	frameName, changed := UpdateAnimationState(state, spriteData, 100)
	if frameName != "" {
		t.Errorf("Expected empty frame name when not playing, got %s", frameName)
	}
	if changed {
		t.Error("Frame should not change when not playing")
	}
	
	// Test paused
	state.Playing = true
	state.Paused = true
	state.Sequence = &cartio.AnimationSequence{
		FrameRefs: []string{"frame1"},
	}
	frameName, changed = UpdateAnimationState(state, spriteData, 100)
	if changed {
		t.Error("Frame should not change when paused")
	}
}

func TestUpdateAnimationState_ForwardLoop(t *testing.T) {
	spriteData := &cartio.SpriteData{
		Type: cartio.SpriteTypeAnimation,
		Frames: []cartio.SpriteFrame{
			{Name: "frame1", Pixels: [][]int{{1}}, Duration: 100},
			{Name: "frame2", Pixels: [][]int{{2}}, Duration: 100},
		},
	}
	
	state := NewAnimationState("test")
	state.Playing = true
	state.Sequence = &cartio.AnimationSequence{
		FrameRefs: []string{"frame1", "frame2"},
		Loop:      true,
		LoopType:  "forward",
	}
	state.CurrentFrame = 0
	state.ElapsedTime = 0
	
	// Update with 50ms - should not advance
	frameName, changed := UpdateAnimationState(state, spriteData, 50)
	if changed {
		t.Error("Frame should not change after 50ms (duration is 100ms)")
	}
	if frameName != "frame1" {
		t.Errorf("Expected 'frame1', got %s", frameName)
	}
	
	// Update with another 60ms (total 110ms) - should advance
	frameName, changed = UpdateAnimationState(state, spriteData, 60)
	if !changed {
		t.Error("Frame should change after exceeding duration")
	}
	if frameName != "frame2" {
		t.Errorf("Expected 'frame2', got %s", frameName)
	}
	if state.CurrentFrame != 1 {
		t.Errorf("Expected frame index 1, got %d", state.CurrentFrame)
	}
	
	// Advance past end - should loop
	state.ElapsedTime = 100
	frameName, changed = UpdateAnimationState(state, spriteData, 100)
	if !changed {
		t.Error("Frame should change when looping")
	}
	if frameName != "frame1" {
		t.Errorf("Expected 'frame1' after loop, got %s", frameName)
	}
	if state.CurrentFrame != 0 {
		t.Errorf("Expected frame index 0 after loop, got %d", state.CurrentFrame)
	}
}

func TestUpdateAnimationState_ForwardNoLoop(t *testing.T) {
	spriteData := &cartio.SpriteData{
		Type: cartio.SpriteTypeAnimation,
		Frames: []cartio.SpriteFrame{
			{Name: "frame1", Pixels: [][]int{{1}}, Duration: 100},
			{Name: "frame2", Pixels: [][]int{{2}}, Duration: 100},
		},
	}
	
	state := NewAnimationState("test")
	state.Playing = true
	state.Sequence = &cartio.AnimationSequence{
		FrameRefs: []string{"frame1", "frame2"},
		Loop:      false,
		LoopType:  "forward",
	}
	state.CurrentFrame = 1
	state.ElapsedTime = 50
	
	// Advance past last frame - should stop immediately
	frameName, changed := UpdateAnimationState(state, spriteData, 60)
	if !changed {
		t.Error("Frame should change to stop")
	}
	if state.Playing {
		t.Error("Animation should stop immediately when reaching end without loop")
	}
	if frameName != "frame2" {
		t.Errorf("Expected 'frame2' (last frame), got %s", frameName)
	}
	if state.CurrentFrame != 1 {
		t.Errorf("Expected frame index 1 (last frame), got %d", state.CurrentFrame)
	}
	
	// Next update should not advance (already stopped)
	frameName, changed = UpdateAnimationState(state, spriteData, 100)
	if changed {
		t.Error("Frame should not change when animation is stopped")
	}
	if state.Playing {
		t.Error("Animation should remain stopped")
	}
}

func TestUpdateAnimationState_Reverse(t *testing.T) {
	spriteData := &cartio.SpriteData{
		Type: cartio.SpriteTypeAnimation,
		Frames: []cartio.SpriteFrame{
			{Name: "frame1", Pixels: [][]int{{1}}, Duration: 100},
			{Name: "frame2", Pixels: [][]int{{2}}, Duration: 100},
		},
	}
	
	state := NewAnimationState("test")
	state.Playing = true
	state.Sequence = &cartio.AnimationSequence{
		FrameRefs: []string{"frame1", "frame2"},
		Loop:      true,
		LoopType:  "reverse",
	}
	state.CurrentFrame = 1
	state.Direction = -1
	state.ElapsedTime = 50
	
	// Advance backwards
	frameName, changed := UpdateAnimationState(state, spriteData, 60)
	if !changed {
		t.Error("Frame should change when advancing backwards")
	}
	if frameName != "frame1" {
		t.Errorf("Expected 'frame1', got %s", frameName)
	}
	if state.CurrentFrame != 0 {
		t.Errorf("Expected frame index 0, got %d", state.CurrentFrame)
	}
	
	// Advance past beginning with loop
	state.ElapsedTime = 100
	frameName, changed = UpdateAnimationState(state, spriteData, 100)
	if !changed {
		t.Error("Frame should change when looping backwards")
	}
	if state.CurrentFrame != 1 {
		t.Errorf("Expected frame index 1 after reverse loop, got %d", state.CurrentFrame)
	}
}

func TestUpdateAnimationState_PingPong(t *testing.T) {
	spriteData := &cartio.SpriteData{
		Type: cartio.SpriteTypeAnimation,
		Frames: []cartio.SpriteFrame{
			{Name: "frame1", Pixels: [][]int{{1}}, Duration: 100},
			{Name: "frame2", Pixels: [][]int{{2}}, Duration: 100},
			{Name: "frame3", Pixels: [][]int{{3}}, Duration: 100},
		},
	}
	
	state := NewAnimationState("test")
	state.Playing = true
	state.Sequence = &cartio.AnimationSequence{
		FrameRefs: []string{"frame1", "frame2", "frame3"},
		Loop:      true,
		LoopType:  "pingpong",
	}
	state.CurrentFrame = 1
	state.Direction = 1
	state.ElapsedTime = 50
	
	// Advance forward to end
	state.CurrentFrame = 2
	state.ElapsedTime = 50
	_, changed := UpdateAnimationState(state, spriteData, 60)
	if !changed {
		t.Error("Frame should change when hitting end in pingpong")
	}
	if state.Direction != -1 {
		t.Error("Direction should reverse to -1 at end")
	}
	if state.CurrentFrame != 1 {
		t.Errorf("Expected frame index 1 after pingpong reverse, got %d", state.CurrentFrame)
	}
	
	// Advance backwards to beginning
	state.CurrentFrame = 0
	state.Direction = -1
	state.ElapsedTime = 50
	_, changed = UpdateAnimationState(state, spriteData, 60)
	if !changed {
		t.Error("Frame should change when hitting beginning in pingpong")
	}
	if state.Direction != 1 {
		t.Error("Direction should reverse to 1 at beginning")
	}
	if state.CurrentFrame != 1 {
		t.Errorf("Expected frame index 1 after pingpong forward, got %d", state.CurrentFrame)
	}
}

func TestUpdateAnimationState_SpeedMultiplier(t *testing.T) {
	spriteData := &cartio.SpriteData{
		Type: cartio.SpriteTypeAnimation,
		Frames: []cartio.SpriteFrame{
			{Name: "frame1", Pixels: [][]int{{1}}, Duration: 100},
		},
	}
	
	state := NewAnimationState("test")
	state.Playing = true
	state.Sequence = &cartio.AnimationSequence{
		FrameRefs: []string{"frame1", "frame2"},
	}
	state.Speed = 2.0 // Double speed
	state.CurrentFrame = 0
	state.ElapsedTime = 0
	
	// With 2x speed, 50ms should be enough (100ms / 2 = 50ms)
	_, changed := UpdateAnimationState(state, spriteData, 50)
	if !changed {
		t.Error("Frame should advance faster with speed multiplier")
	}
	
	// Test slow speed (0.5x)
	state.Speed = 0.5
	state.CurrentFrame = 0
	state.ElapsedTime = 0
	_, changed = UpdateAnimationState(state, spriteData, 100)
	if changed {
		t.Error("Frame should not advance with 0.5x speed at 100ms (needs 200ms)")
	}
	
	_, changed = UpdateAnimationState(state, spriteData, 110)
	if !changed {
		t.Error("Frame should advance after 210ms total with 0.5x speed")
	}
}

func TestUpdateAnimationState_InvalidFrameIndex(t *testing.T) {
	spriteData := &cartio.SpriteData{
		Type: cartio.SpriteTypeAnimation,
		Frames: []cartio.SpriteFrame{
			{Name: "frame1", Pixels: [][]int{{1}}, Duration: 100},
		},
	}
	
	state := NewAnimationState("test")
	state.Playing = true
	state.Sequence = &cartio.AnimationSequence{
		FrameRefs: []string{"frame1"},
	}
	state.CurrentFrame = 10 // Invalid index
	
	frameName, changed := UpdateAnimationState(state, spriteData, 100)
	if !changed {
		t.Error("Frame should change when resetting invalid index")
	}
	if state.CurrentFrame != 0 {
		t.Errorf("Expected frame reset to 0, got %d", state.CurrentFrame)
	}
	if frameName != "frame1" {
		t.Errorf("Expected 'frame1', got %s", frameName)
	}
}

func TestUpdateAnimationState_DefaultFrameDuration(t *testing.T) {
	// Test with sprite data that doesn't have the frame
	spriteData := &cartio.SpriteData{
		Type: cartio.SpriteTypeAnimation,
		Frames: []cartio.SpriteFrame{
			{Name: "other", Pixels: [][]int{{1}}, Duration: 50},
		},
	}
	
	state := NewAnimationState("test")
	state.Playing = true
	state.Sequence = &cartio.AnimationSequence{
		FrameRefs: []string{"frame1"},
	}
	state.CurrentFrame = 0
	state.ElapsedTime = 0
	
	// Should use default 100ms duration
	_, changed := UpdateAnimationState(state, spriteData, 50)
	if changed {
		t.Error("Frame should not change before default duration")
	}
	
	_, changed = UpdateAnimationState(state, spriteData, 60)
	if !changed {
		t.Error("Frame should change after default duration")
	}
}

func TestGetCurrentFrameName_EdgeCases(t *testing.T) {
	// Test with out of bounds frame index
	state := NewAnimationState("test")
	state.Sequence = &cartio.AnimationSequence{
		FrameRefs: []string{"frame1", "frame2"},
	}
	state.CurrentFrame = 5 // Out of bounds
	
	name := GetCurrentFrameName(state)
	if name != "" {
		t.Errorf("Expected empty string for out-of-bounds frame, got %s", name)
	}
	
	// Test with negative frame index
	state.CurrentFrame = -1
	name = GetCurrentFrameName(state)
	if name != "" {
		t.Errorf("Expected empty string for negative frame, got %s", name)
	}
}

func TestUpdateAnimationState_EdgeCases(t *testing.T) {
	// Test with nil state
	frameName, changed := UpdateAnimationState(nil, nil, 100)
	if frameName != "" {
		t.Error("Expected empty frame name for nil state")
	}
	if changed {
		t.Error("Expected no change for nil state")
	}
	
	// Test with playing but nil sequence
	state := NewAnimationState("test")
	state.Playing = true
	state.Sequence = nil
	frameName, changed = UpdateAnimationState(state, nil, 100)
	if frameName != "" {
		t.Error("Expected empty frame name for nil sequence")
	}
	if changed {
		t.Error("Expected no change for nil sequence")
	}
	
	// Test with sequence but empty frameRefs
	state2 := NewAnimationState("test")
	state2.Playing = true
	state2.Sequence = &cartio.AnimationSequence{
		FrameRefs: []string{},
	}
	frameName, changed = UpdateAnimationState(state2, nil, 100)
	if frameName != "" {
		t.Error("Expected empty frame name for empty frameRefs")
	}
	
	// Test with sequence but frameRefs index out of bounds
	state3 := NewAnimationState("test")
	state3.Playing = true
	state3.Sequence = &cartio.AnimationSequence{
		FrameRefs: []string{"frame1"},
	}
	state3.CurrentFrame = 5 // Out of bounds
	frameName, changed = UpdateAnimationState(state3, nil, 100)
	if !changed {
		t.Error("Should reset invalid frame index")
	}
	if state3.CurrentFrame != 0 {
		t.Errorf("Expected frame reset to 0, got %d", state3.CurrentFrame)
	}
}

func TestUpdateAnimationState_ReverseNoLoop(t *testing.T) {
	spriteData := &cartio.SpriteData{
		Type: cartio.SpriteTypeAnimation,
		Frames: []cartio.SpriteFrame{
			{Name: "frame1", Pixels: [][]int{{1}}, Duration: 100},
			{Name: "frame2", Pixels: [][]int{{2}}, Duration: 100},
		},
	}
	
	state := NewAnimationState("test")
	state.Playing = true
	state.Sequence = &cartio.AnimationSequence{
		FrameRefs: []string{"frame1", "frame2"},
		Loop:      false,
		LoopType:  "reverse",
	}
	state.CurrentFrame = 0
	state.Direction = -1
	state.ElapsedTime = 50
	
	// Advance past beginning - should stop
	frameName, changed := UpdateAnimationState(state, spriteData, 60)
	if !changed {
		t.Error("Frame should change to stop")
	}
	if state.Playing {
		t.Error("Animation should stop when reaching beginning without loop")
	}
	if frameName != "frame1" {
		t.Errorf("Expected 'frame1' (first frame), got %s", frameName)
	}
	if state.CurrentFrame != 0 {
		t.Errorf("Expected frame index 0 (first frame), got %d", state.CurrentFrame)
	}
}

func TestUpdateAnimationState_PingPongNoLoop(t *testing.T) {
	spriteData := &cartio.SpriteData{
		Type: cartio.SpriteTypeAnimation,
		Frames: []cartio.SpriteFrame{
			{Name: "frame1", Pixels: [][]int{{1}}, Duration: 100},
			{Name: "frame2", Pixels: [][]int{{2}}, Duration: 100},
		},
	}
	
	state := NewAnimationState("test")
	state.Playing = true
	state.Sequence = &cartio.AnimationSequence{
		FrameRefs: []string{"frame1", "frame2"},
		Loop:      false,
		LoopType:  "pingpong",
	}
	state.CurrentFrame = 1
	state.Direction = 1
	state.ElapsedTime = 50
	
	// Advance to end - should stop (no loop)
	_, changed := UpdateAnimationState(state, spriteData, 60)
	if !changed {
		t.Error("Frame should change when hitting end")
	}
	if state.Playing {
		t.Error("Animation should stop when reaching end without loop in pingpong")
	}
	
	// Test pingpong with single frame
	state2 := NewAnimationState("test")
	state2.Playing = true
	state2.Sequence = &cartio.AnimationSequence{
		FrameRefs: []string{"frame1"},
		Loop:      true,
		LoopType:  "pingpong",
	}
	state2.CurrentFrame = 0
	state2.Direction = 1
	state2.ElapsedTime = 50
	
	// With single frame, should wrap around
	_, changed = UpdateAnimationState(state2, spriteData, 60)
	if state2.CurrentFrame != 0 {
		t.Errorf("Expected frame 0 for single frame pingpong, got %d", state2.CurrentFrame)
	}
}

func TestUpdateAnimationState_MinimumDuration(t *testing.T) {
	spriteData := &cartio.SpriteData{
		Type: cartio.SpriteTypeAnimation,
		Frames: []cartio.SpriteFrame{
			{Name: "frame1", Pixels: [][]int{{1}}, Duration: 1}, // Very short
		},
	}
	
	state := NewAnimationState("test")
	state.Playing = true
	state.Sequence = &cartio.AnimationSequence{
		FrameRefs: []string{"frame1", "frame2"},
	}
	state.Speed = 100.0 // Very high speed (should result in < 1ms duration)
	state.CurrentFrame = 0
	state.ElapsedTime = 0
	
	// Should enforce minimum 1ms duration
	_, changed := UpdateAnimationState(state, spriteData, 1)
	if !changed {
		t.Error("Frame should advance even with very high speed (minimum 1ms)")
	}
}

