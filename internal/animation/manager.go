package animation

import (
	"fmt"

	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
)

// AnimationState tracks the runtime animation state for a sprite instance
type AnimationState struct {
	SpriteName    string                   // Name of the sprite
	AnimationName string                   // Name of the active animation (empty if none)
	Sequence      *cartio.AnimationSequence // Reference to animation sequence
	CurrentFrame  int                      // Current index in frameRefs array
	ElapsedTime   int64                    // Milliseconds elapsed in current frame
	Playing       bool                     // Whether animation is currently playing
	Paused        bool                     // Whether animation is paused
	Speed         float64                  // Current speed multiplier (from animation or runtime override)
	Direction     int                      // Direction for pingpong: 1 = forward, -1 = reverse
}

// NewAnimationState creates a new animation state
func NewAnimationState(spriteName string) *AnimationState {
	return &AnimationState{
		SpriteName:    spriteName,
		AnimationName: "",
		Sequence:      nil,
		CurrentFrame:  0,
		ElapsedTime:   0,
		Playing:       false,
		Paused:        false,
		Speed:         1.0,
		Direction:     1, // Start forward
	}
}

// Reset resets the animation state to initial values
func (s *AnimationState) Reset() {
	s.AnimationName = ""
	s.Sequence = nil
	s.CurrentFrame = 0
	s.ElapsedTime = 0
	s.Playing = false
	s.Paused = false
	s.Speed = 1.0
	s.Direction = 1
}

// UpdateAnimationState updates a single animation state based on elapsed time
// Returns the current frame name and whether the frame changed
func UpdateAnimationState(state *AnimationState, spriteData *cartio.SpriteData, deltaTimeMs int64) (string, bool) {
	if state == nil || !state.Playing || state.Paused || state.Sequence == nil {
		// Not playing or paused - return current frame name if available
		if state != nil && state.Sequence != nil && len(state.Sequence.FrameRefs) > 0 && state.CurrentFrame < len(state.Sequence.FrameRefs) {
			return state.Sequence.FrameRefs[state.CurrentFrame], false
		}
		return "", false
	}

	// Check for empty frame refs
	if len(state.Sequence.FrameRefs) == 0 {
		return "", false
	}

	frameChanged := false
	oldFrame := state.CurrentFrame

	// Get current frame reference
	if state.CurrentFrame >= len(state.Sequence.FrameRefs) {
		// Invalid frame index - reset to 0
		state.CurrentFrame = 0
		state.ElapsedTime = 0
		return state.Sequence.FrameRefs[0], true
	}

	currentFrameName := state.Sequence.FrameRefs[state.CurrentFrame]

	// Find frame duration from sprite frames
	var frameDuration int64 = 100 // Default 100ms if not found
	if spriteData != nil {
		for _, frame := range spriteData.Frames {
			if frame.Name == currentFrameName {
				frameDuration = int64(frame.Duration)
				break
			}
		}
	}

	// Apply speed multiplier
	adjustedDuration := int64(float64(frameDuration) / state.Speed)
	if adjustedDuration < 1 {
		adjustedDuration = 1 // Minimum 1ms
	}

	// Update elapsed time
	state.ElapsedTime += deltaTimeMs

	// Check if we should advance to next frame
	if state.ElapsedTime >= adjustedDuration {
		// Advance frame
		state.ElapsedTime = 0 // Reset elapsed time for new frame

		// Handle loop types
		if state.Sequence.LoopType == "reverse" {
			// Reverse: play backwards
			state.Direction = -1
			state.CurrentFrame--
			if state.CurrentFrame < 0 {
				if state.Sequence.Loop {
					state.CurrentFrame = len(state.Sequence.FrameRefs) - 1
				} else {
					// Stop at first frame
					state.CurrentFrame = 0
					state.Playing = false
				}
			}
		} else if state.Sequence.LoopType == "pingpong" {
			// Pingpong: play forward, then reverse, then forward again
			state.CurrentFrame += state.Direction
			if state.CurrentFrame >= len(state.Sequence.FrameRefs) {
				// Reached end - reverse direction
				state.Direction = -1
				state.CurrentFrame = len(state.Sequence.FrameRefs) - 2 // Go back one
				if state.CurrentFrame < 0 {
					state.CurrentFrame = 0
				}
				if !state.Sequence.Loop {
					state.Playing = false
				}
			} else if state.CurrentFrame < 0 {
				// Reached beginning - forward direction
				state.Direction = 1
				state.CurrentFrame = 1
				if state.CurrentFrame >= len(state.Sequence.FrameRefs) {
					state.CurrentFrame = len(state.Sequence.FrameRefs) - 1
				}
				if !state.Sequence.Loop {
					state.Playing = false
				}
			}
		} else {
			// Forward (default)
			state.Direction = 1
			state.CurrentFrame++
			if state.CurrentFrame >= len(state.Sequence.FrameRefs) {
				if state.Sequence.Loop {
					state.CurrentFrame = 0 // Loop back to start
				} else {
					// Stop at last frame
					state.CurrentFrame = len(state.Sequence.FrameRefs) - 1
					state.Playing = false
				}
			}
		}

		frameChanged = true
	}

	// Return current frame name
	if state.CurrentFrame < len(state.Sequence.FrameRefs) {
		newFrameName := state.Sequence.FrameRefs[state.CurrentFrame]
		return newFrameName, frameChanged || (oldFrame != state.CurrentFrame)
	}

	return currentFrameName, frameChanged
}

// GetCurrentFrameName returns the name of the current frame in the animation
func GetCurrentFrameName(state *AnimationState) string {
	if state == nil || state.Sequence == nil || len(state.Sequence.FrameRefs) == 0 {
		return ""
	}
	if state.CurrentFrame >= 0 && state.CurrentFrame < len(state.Sequence.FrameRefs) {
		return state.Sequence.FrameRefs[state.CurrentFrame]
	}
	return ""
}

// PlayAnimation starts playing an animation for a sprite instance
func PlayAnimation(state *AnimationState, spriteData *cartio.SpriteData, animationName string) error {
	if spriteData == nil {
		return fmt.Errorf("sprite data is nil")
	}
	if spriteData.Type != cartio.SpriteTypeAnimation {
		return fmt.Errorf("sprite type '%s' does not support animations", spriteData.Type)
	}

	// Find animation sequence
	anim, err := spriteData.GetAnimation(animationName)
	if err != nil {
		return fmt.Errorf("animation '%s' not found: %w", animationName, err)
	}

	// Set up animation state
	state.AnimationName = animationName
	state.Sequence = anim
	state.CurrentFrame = 0
	state.ElapsedTime = 0
	state.Playing = true
	state.Paused = false
	state.Speed = anim.Speed
	if state.Speed <= 0 {
		state.Speed = 1.0
	}

	// Set direction based on loop type
	if anim.LoopType == "reverse" {
		state.Direction = -1
		state.CurrentFrame = len(anim.FrameRefs) - 1 // Start at end for reverse
	} else {
		state.Direction = 1 // Forward or pingpong start forward
	}

	return nil
}

// PauseAnimation pauses the current animation
func PauseAnimation(state *AnimationState) {
	if state != nil {
		state.Paused = true
	}
}

// ResumeAnimation resumes a paused animation
func ResumeAnimation(state *AnimationState) {
	if state != nil {
		state.Paused = false
	}
}

// StopAnimation stops the current animation
func StopAnimation(state *AnimationState) {
	if state != nil {
		state.Playing = false
		state.Paused = false
		state.CurrentFrame = 0
		state.ElapsedTime = 0
		// Keep sequence and animation name for potential resume
	}
}

// SetAnimationSpeed sets the speed multiplier for an animation
func SetAnimationSpeed(state *AnimationState, speed float64) {
	if state != nil {
		if speed > 0 {
			state.Speed = speed
		} else {
			state.Speed = 1.0 // Default to 1.0 if invalid
		}
	}
}

// SetAnimationFrame sets the current frame index in the animation
func SetAnimationFrame(state *AnimationState, frameIndex int) error {
	if state == nil {
		return fmt.Errorf("animation state is nil")
	}
	if state.Sequence == nil {
		return fmt.Errorf("no animation sequence active")
	}
	if frameIndex < 0 || frameIndex >= len(state.Sequence.FrameRefs) {
		return fmt.Errorf("frame index %d out of range [0, %d)", frameIndex, len(state.Sequence.FrameRefs))
	}

	state.CurrentFrame = frameIndex
	state.ElapsedTime = 0 // Reset elapsed time when manually setting frame
	return nil
}

// GetAnimationFrame returns the current frame index
func GetAnimationFrame(state *AnimationState) int {
	if state == nil || state.Sequence == nil {
		return -1
	}
	return state.CurrentFrame
}

