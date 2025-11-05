package cartio

import (
	"fmt"
	"strings"
)

// TileData represents a single tile (similar to SpriteData but without isUI/lifetime)
// Can be static, frames, or animation type
type TileData struct {
	Width      int                `json:"width"`      // Tile width
	Height     int                `json:"height"`     // Tile height
	Pixels     [][]int            `json:"pixels,omitempty"` // Static tile pixels (if type is static)
	Type       SpriteType         `json:"type"`       // "static", "frames", or "animation"
	Frames     []SpriteFrame      `json:"frames,omitempty"` // Named frames (if type is frames or animation)
	Animations []AnimationSequence `json:"animations,omitempty"` // Animation sequences (if type is animation)
	UseCollision bool             `json:"useCollision"` // Enable collision detection
	MountPoints  []MountPoint     `json:"mountPoints"`  // Array of mount points
	IsISO        bool             `json:"isISO"`         // If true, convert to isometric when loaded
}

// TilesetMap maps tile names to their data
type TilesetMap map[string]TileData

// NormalizeTileData ensures tile data is in a consistent state
func NormalizeTileData(tile *TileData) {
	// Default type to "static" if not specified
	if tile.Type == "" {
		tile.Type = SpriteTypeStatic
	}

	// Normalize animations (same as sprites)
	for i := range tile.Animations {
		if tile.Animations[i].LoopType == "" {
			tile.Animations[i].LoopType = "forward"
		}
		if tile.Animations[i].Speed <= 0 {
			tile.Animations[i].Speed = 1.0
		}
		// Trim whitespace from frame refs
		for j := range tile.Animations[i].FrameRefs {
			tile.Animations[i].FrameRefs[j] = strings.TrimSpace(tile.Animations[i].FrameRefs[j])
		}
	}
}

// ValidateTileData validates a complete tile data structure
func ValidateTileData(tile *TileData, tileName string) error {
	// Validate dimensions (same as sprites, but tiles are always gameplay sprites)
	if err := ValidateSpriteSize(tile.Width, tile.Height, false); err != nil {
		return fmt.Errorf("tile '%s': %w", tileName, err)
	}

	// Default type to "static" if not specified
	if tile.Type == "" {
		tile.Type = SpriteTypeStatic
	}

	// Validate type
	if tile.Type != SpriteTypeStatic && tile.Type != SpriteTypeFrames && tile.Type != SpriteTypeAnimation {
		return fmt.Errorf("tile '%s': invalid type '%s', must be 'static', 'frames', or 'animation'", tileName, tile.Type)
	}

	// Validate based on type (same validation logic as sprites)
	switch tile.Type {
	case SpriteTypeStatic:
		if tile.Pixels == nil || len(tile.Pixels) == 0 {
			return fmt.Errorf("tile '%s': static tile must have pixels", tileName)
		}
		if len(tile.Pixels) != tile.Height {
			return fmt.Errorf("tile '%s': pixel array height (%d) does not match tile height (%d)", tileName, len(tile.Pixels), tile.Height)
		}
		for i, row := range tile.Pixels {
			if len(row) != tile.Width {
				return fmt.Errorf("tile '%s': row %d width (%d) does not match tile width (%d)", tileName, i, len(row), tile.Width)
			}
		}

	case SpriteTypeFrames:
		if len(tile.Frames) == 0 {
			return fmt.Errorf("tile '%s': frames tile must have at least one frame", tileName)
		}
		frameNames := make(map[string]bool)
		for i, frame := range tile.Frames {
			if err := ValidateFrameName(frame.Name); err != nil {
				return fmt.Errorf("tile '%s' frame %d: %w", tileName, i, err)
			}
			if frameNames[frame.Name] {
				return fmt.Errorf("tile '%s': duplicate frame name '%s'", tileName, frame.Name)
			}
			frameNames[frame.Name] = true
			if frame.Pixels == nil || len(frame.Pixels) == 0 {
				return fmt.Errorf("tile '%s' frame '%s': must have pixels", tileName, frame.Name)
			}
			if len(frame.Pixels) != tile.Height {
				return fmt.Errorf("tile '%s' frame '%s': pixel array height (%d) does not match tile height (%d)", tileName, frame.Name, len(frame.Pixels), tile.Height)
			}
			for j, row := range frame.Pixels {
				if len(row) != tile.Width {
					return fmt.Errorf("tile '%s' frame '%s' row %d: width (%d) does not match tile width (%d)", tileName, frame.Name, j, len(row), tile.Width)
				}
			}
		}

	case SpriteTypeAnimation:
		if len(tile.Frames) == 0 {
			return fmt.Errorf("tile '%s': animation tile must have at least one frame", tileName)
		}
		if len(tile.Animations) == 0 {
			return fmt.Errorf("tile '%s': animation tile must have at least one animation", tileName)
		}
		frameNames := make(map[string]bool)
		for i, frame := range tile.Frames {
			if err := ValidateFrameName(frame.Name); err != nil {
				return fmt.Errorf("tile '%s' frame %d: %w", tileName, i, err)
			}
			if frameNames[frame.Name] {
				return fmt.Errorf("tile '%s': duplicate frame name '%s'", tileName, frame.Name)
			}
			frameNames[frame.Name] = true
			if frame.Pixels == nil || len(frame.Pixels) == 0 {
				return fmt.Errorf("tile '%s' frame '%s': must have pixels", tileName, frame.Name)
			}
			if len(frame.Pixels) != tile.Height {
				return fmt.Errorf("tile '%s' frame '%s': pixel array height (%d) does not match tile height (%d)", tileName, frame.Name, len(frame.Pixels), tile.Height)
			}
			for j, row := range frame.Pixels {
				if len(row) != tile.Width {
					return fmt.Errorf("tile '%s' frame '%s' row %d: width (%d) does not match tile width (%d)", tileName, frame.Name, j, len(row), tile.Width)
				}
			}
		}
		animationNames := make(map[string]bool)
		for i, anim := range tile.Animations {
			if err := ValidateFrameName(anim.Name); err != nil {
				return fmt.Errorf("tile '%s' animation %d: %w", tileName, i, err)
			}
			if animationNames[anim.Name] {
				return fmt.Errorf("tile '%s': duplicate animation name '%s'", tileName, anim.Name)
			}
			animationNames[anim.Name] = true
			if len(anim.FrameRefs) == 0 {
				return fmt.Errorf("tile '%s' animation '%s': must have at least one frame reference", tileName, anim.Name)
			}
			for _, frameRef := range anim.FrameRefs {
				if !frameNames[frameRef] {
					return fmt.Errorf("tile '%s' animation '%s': references unknown frame '%s'", tileName, anim.Name, frameRef)
				}
			}
			if anim.LoopType != "" && anim.LoopType != "forward" && anim.LoopType != "reverse" && anim.LoopType != "pingpong" {
				return fmt.Errorf("tile '%s' animation '%s': invalid loopType '%s', must be 'forward', 'reverse', or 'pingpong'", tileName, anim.Name, anim.LoopType)
			}
			if anim.LoopType == "" {
				anim.LoopType = "forward"
			}
			if anim.Speed <= 0 {
				anim.Speed = 1.0
			}
		}
	}

	return nil
}

// GetFramePixels returns the pixel data for a specific frame name (same as SpriteData)
func (t *TileData) GetFramePixels(frameName string) ([][]int, error) {
	switch t.Type {
	case SpriteTypeStatic:
		return t.Pixels, nil
	case SpriteTypeFrames, SpriteTypeAnimation:
		for _, frame := range t.Frames {
			if frame.Name == frameName {
				return frame.Pixels, nil
			}
		}
		return nil, fmt.Errorf("frame '%s' not found in tile", frameName)
	default:
		return nil, fmt.Errorf("unknown tile type: %s", t.Type)
	}
}

