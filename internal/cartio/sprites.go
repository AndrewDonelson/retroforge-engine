package cartio

import (
	"fmt"
	"strings"
)

// SpriteType indicates the type of sprite data
type SpriteType string

const (
	SpriteTypeStatic    SpriteType = "static"    // Single frame sprite
	SpriteTypeFrames    SpriteType = "frames"    // Multiple named frames (states)
	SpriteTypeAnimation SpriteType = "animation" // Animated sequence
)

// MountPoint represents a point within a sprite where projectiles/thrusters originate
type MountPoint struct {
	X    int    `json:"x"`              // X coordinate within sprite bounds
	Y    int    `json:"y"`              // Y coordinate within sprite bounds
	Name string `json:"name,omitempty"` // Optional name for accessing by name in Lua
}

// SpriteFrame represents a single frame within a multi-frame sprite
type SpriteFrame struct {
	Name     string  `json:"name"`              // Frame name (e.g., "idle", "left", "right")
	Pixels   [][]int `json:"pixels"`            // Frame pixel data
	Duration int     `json:"duration"`          // Frame duration in milliseconds (for animations)
}

// AnimationSequence represents an animation sequence
type AnimationSequence struct {
	Name      string   `json:"name"`              // Animation name (e.g., "walk", "shoot")
	FrameRefs []string `json:"frameRefs"`        // Array of frame names in sequence
	Speed     float64  `json:"speed"`             // Animation speed multiplier (default 1.0)
	Loop      bool     `json:"loop"`              // Whether to loop animation
	LoopType  string   `json:"loopType"`          // "forward", "reverse", or "pingpong" (default: "forward")
}

// SpriteData represents a single sprite (can be static, multi-frame, or animated)
type SpriteData struct {
	Width        int                `json:"width"`        // Sprite width (for static/frame reference)
	Height       int                `json:"height"`       // Sprite height (for static/frame reference)
	Pixels       [][]int            `json:"pixels,omitempty"` // Static sprite pixels (if type is static)
	Type         SpriteType         `json:"type"`         // "static", "frames", or "animation"
	Frames       []SpriteFrame      `json:"frames,omitempty"` // Named frames (if type is frames or animation)
	Animations   []AnimationSequence `json:"animations,omitempty"` // Animation sequences (if type is animation)
	UseCollision bool               `json:"useCollision"` // Enable collision detection with other sprites
	MountPoints  []MountPoint        `json:"mountPoints"`  // Array of mount points (e.g., for bullets, thrusters)
	IsUI         bool               `json:"isUI"`         // If true, sprite is UI element and not affected by physics
	Lifetime     int                `json:"lifetime"`     // Lifetime in milliseconds (0 = no lifetime limit)
	MaxSpawn     int                `json:"maxSpawn"`     // Maximum instances that can be spawned simultaneously (0 = no limit)
}

// SpriteMap maps sprite names to their data
type SpriteMap map[string]SpriteData

// ValidateSpriteSize validates sprite dimensions based on new size rules:
// - Minimum size: 2x2 for all sprites
// - Gameplay sprites (isUI=false): 2x2 to 32x32 (any size)
// - UI sprites (isUI=true): 2x2 to 256x256 (both dimensions must be divisible by 2)
func ValidateSpriteSize(width, height int, isUI bool) error {
	// Minimum size: 2x2
	if width < 2 || height < 2 {
		return fmt.Errorf("sprite dimensions must be at least 2x2, got %dx%d", width, height)
	}

	if isUI {
		// UI sprites: 2-256, both dimensions divisible by 2
		if width > 256 || height > 256 {
			return fmt.Errorf("UI sprite dimensions cannot exceed 256, got %dx%d", width, height)
		}
		if width%2 != 0 || height%2 != 0 {
			return fmt.Errorf("UI sprite dimensions must be divisible by 2, got %dx%d", width, height)
		}
	} else {
		// Gameplay sprites: 2-32
		if width > 32 || height > 32 {
			return fmt.Errorf("gameplay sprite dimensions cannot exceed 32x32, got %dx%d", width, height)
		}
	}
	return nil
}

// ValidateFrameName validates a frame name according to naming rules
func ValidateFrameName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("frame name cannot be empty")
	}
	if len(name) > 64 {
		return fmt.Errorf("frame name cannot exceed 64 characters, got %d", len(name))
	}
	
	// Must start with letter or underscore
	first := name[0]
	if !((first >= 'a' && first <= 'z') || (first >= 'A' && first <= 'Z') || first == '_') {
		return fmt.Errorf("frame name must start with a letter or underscore, got '%c'", first)
	}
	
	// Only alphanumeric, underscore, and hyphen allowed
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || 
			(char >= '0' && char <= '9') || char == '_' || char == '-') {
			return fmt.Errorf("frame name contains invalid character '%c' (only alphanumeric, underscore, hyphen allowed)", char)
		}
	}
	
	return nil
}

// ValidateSpriteData validates a complete sprite data structure
func ValidateSpriteData(sprite *SpriteData, spriteName string) error {
	// Validate dimensions
	if err := ValidateSpriteSize(sprite.Width, sprite.Height, sprite.IsUI); err != nil {
		return fmt.Errorf("sprite '%s': %w", spriteName, err)
	}
	
	// Default type to "static" if not specified (backward compatibility)
	if sprite.Type == "" {
		sprite.Type = SpriteTypeStatic
	}
	
	// Validate type
	if sprite.Type != SpriteTypeStatic && sprite.Type != SpriteTypeFrames && sprite.Type != SpriteTypeAnimation {
		return fmt.Errorf("sprite '%s': invalid type '%s', must be 'static', 'frames', or 'animation'", spriteName, sprite.Type)
	}
	
	// Validate based on type
	switch sprite.Type {
	case SpriteTypeStatic:
		// Static sprites must have pixels
		if sprite.Pixels == nil || len(sprite.Pixels) == 0 {
			return fmt.Errorf("sprite '%s': static sprite must have pixels", spriteName)
		}
		// Validate pixel dimensions
		if len(sprite.Pixels) != sprite.Height {
			return fmt.Errorf("sprite '%s': pixel array height (%d) does not match sprite height (%d)", spriteName, len(sprite.Pixels), sprite.Height)
		}
		for i, row := range sprite.Pixels {
			if len(row) != sprite.Width {
				return fmt.Errorf("sprite '%s': row %d width (%d) does not match sprite width (%d)", spriteName, i, len(row), sprite.Width)
			}
		}
		
	case SpriteTypeFrames:
		// Frames sprites must have frames array
		if len(sprite.Frames) == 0 {
			return fmt.Errorf("sprite '%s': frames sprite must have at least one frame", spriteName)
		}
		
		// Validate each frame
		frameNames := make(map[string]bool)
		for i, frame := range sprite.Frames {
			// Validate frame name
			if err := ValidateFrameName(frame.Name); err != nil {
				return fmt.Errorf("sprite '%s' frame %d: %w", spriteName, i, err)
			}
			
			// Check for duplicate frame names
			if frameNames[frame.Name] {
				return fmt.Errorf("sprite '%s': duplicate frame name '%s'", spriteName, frame.Name)
			}
			frameNames[frame.Name] = true
			
			// Validate frame pixels
			if frame.Pixels == nil || len(frame.Pixels) == 0 {
				return fmt.Errorf("sprite '%s' frame '%s': must have pixels", spriteName, frame.Name)
			}
			if len(frame.Pixels) != sprite.Height {
				return fmt.Errorf("sprite '%s' frame '%s': pixel array height (%d) does not match sprite height (%d)", spriteName, frame.Name, len(frame.Pixels), sprite.Height)
			}
			for j, row := range frame.Pixels {
				if len(row) != sprite.Width {
					return fmt.Errorf("sprite '%s' frame '%s' row %d: width (%d) does not match sprite width (%d)", spriteName, frame.Name, j, len(row), sprite.Width)
				}
			}
		}
		
	case SpriteTypeAnimation:
		// Animation sprites must have frames and animations
		if len(sprite.Frames) == 0 {
			return fmt.Errorf("sprite '%s': animation sprite must have at least one frame", spriteName)
		}
		if len(sprite.Animations) == 0 {
			return fmt.Errorf("sprite '%s': animation sprite must have at least one animation", spriteName)
		}
		
		// Validate frames (same as frames type)
		frameNames := make(map[string]bool)
		for i, frame := range sprite.Frames {
			if err := ValidateFrameName(frame.Name); err != nil {
				return fmt.Errorf("sprite '%s' frame %d: %w", spriteName, i, err)
			}
			if frameNames[frame.Name] {
				return fmt.Errorf("sprite '%s': duplicate frame name '%s'", spriteName, frame.Name)
			}
			frameNames[frame.Name] = true
			
			if frame.Pixels == nil || len(frame.Pixels) == 0 {
				return fmt.Errorf("sprite '%s' frame '%s': must have pixels", spriteName, frame.Name)
			}
			if len(frame.Pixels) != sprite.Height {
				return fmt.Errorf("sprite '%s' frame '%s': pixel array height (%d) does not match sprite height (%d)", spriteName, frame.Name, len(frame.Pixels), sprite.Height)
			}
			for j, row := range frame.Pixels {
				if len(row) != sprite.Width {
					return fmt.Errorf("sprite '%s' frame '%s' row %d: width (%d) does not match sprite width (%d)", spriteName, frame.Name, j, len(row), sprite.Width)
				}
			}
		}
		
		// Validate animations
		animationNames := make(map[string]bool)
		for i, anim := range sprite.Animations {
			// Validate animation name
			if err := ValidateFrameName(anim.Name); err != nil {
				return fmt.Errorf("sprite '%s' animation %d: %w", spriteName, i, err)
			}
			
			if animationNames[anim.Name] {
				return fmt.Errorf("sprite '%s': duplicate animation name '%s'", spriteName, anim.Name)
			}
			animationNames[anim.Name] = true
			
			// Validate frame references
			if len(anim.FrameRefs) == 0 {
				return fmt.Errorf("sprite '%s' animation '%s': must have at least one frame reference", spriteName, anim.Name)
			}
			
			for _, frameRef := range anim.FrameRefs {
				if !frameNames[frameRef] {
					return fmt.Errorf("sprite '%s' animation '%s': references unknown frame '%s'", spriteName, anim.Name, frameRef)
				}
			}
			
			// Validate loop type
			if anim.LoopType != "" && anim.LoopType != "forward" && anim.LoopType != "reverse" && anim.LoopType != "pingpong" {
				return fmt.Errorf("sprite '%s' animation '%s': invalid loopType '%s', must be 'forward', 'reverse', or 'pingpong'", spriteName, anim.Name, anim.LoopType)
			}
			
			// Default loop type to "forward" if not specified
			if anim.LoopType == "" {
				anim.LoopType = "forward"
			}
			
			// Validate speed (must be positive)
			if anim.Speed <= 0 {
				anim.Speed = 1.0 // Default to 1.0
			}
		}
	}
	
	return nil
}

// GetFramePixels returns the pixel data for a specific frame name
func (s *SpriteData) GetFramePixels(frameName string) ([][]int, error) {
	switch s.Type {
	case SpriteTypeStatic:
		return s.Pixels, nil
	case SpriteTypeFrames, SpriteTypeAnimation:
		for _, frame := range s.Frames {
			if frame.Name == frameName {
				return frame.Pixels, nil
			}
		}
		return nil, fmt.Errorf("frame '%s' not found in sprite", frameName)
	default:
		return nil, fmt.Errorf("unknown sprite type: %s", s.Type)
	}
}

// HasFrame checks if a frame with the given name exists
func (s *SpriteData) HasFrame(frameName string) bool {
	if s.Type == SpriteTypeStatic {
		return frameName == "" || frameName == "default"
	}
	for _, frame := range s.Frames {
		if frame.Name == frameName {
			return true
		}
	}
	return false
}

// GetAnimation returns an animation sequence by name
func (s *SpriteData) GetAnimation(animationName string) (*AnimationSequence, error) {
	if s.Type != SpriteTypeAnimation {
		return nil, fmt.Errorf("sprite type '%s' does not support animations", s.Type)
	}
	for i := range s.Animations {
		if s.Animations[i].Name == animationName {
			return &s.Animations[i], nil
		}
	}
	return nil, fmt.Errorf("animation '%s' not found", animationName)
}

// NormalizeSpriteData ensures sprite data is in a consistent state
func NormalizeSpriteData(sprite *SpriteData) {
	// Default type to "static" if not specified
	if sprite.Type == "" {
		sprite.Type = SpriteTypeStatic
	}
	
	// Normalize animations
	for i := range sprite.Animations {
		if sprite.Animations[i].LoopType == "" {
			sprite.Animations[i].LoopType = "forward"
		}
		if sprite.Animations[i].Speed <= 0 {
			sprite.Animations[i].Speed = 1.0
		}
		// Trim whitespace from frame refs
		for j := range sprite.Animations[i].FrameRefs {
			sprite.Animations[i].FrameRefs[j] = strings.TrimSpace(sprite.Animations[i].FrameRefs[j])
		}
	}
}
