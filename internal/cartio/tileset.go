package cartio

import (
	"encoding/json"
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
}

// TilesetMap maps tile names to their data
type TilesetMap map[string]TileData

// TilesetData represents a complete tileset with metadata
type TilesetData struct {
	IsISO    bool       `json:"isISO"`    // If true, tileset is isometric (renders using isometric transformation)
	Seed     int        `json:"seed,omitempty"` // Seed for deterministic tile variation (rotation/flipping) - only for normal tiles
	Tiles    TilesetMap `json:"tiles"`     // Map of tile names to tile data (or root-level tiles if isISO not present)
}

// tilesetJSON is used for JSON unmarshaling (handles both old format and new format)
type tilesetJSON struct {
	IsISO bool       `json:"isISO"`
	Seed  int        `json:"seed,omitempty"`
	Tiles TilesetMap `json:"tiles"`
	// Support old format where tiles are at root level
	TilesetMap
}

// UnmarshalJSON handles both old format (tiles at root) and new format (tiles in "tiles" field)
func (ts *TilesetData) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Check if "isISO" field exists
	if isISOData, hasIsISO := raw["isISO"]; hasIsISO {
		var isISO bool
		if err := json.Unmarshal(isISOData, &isISO); err != nil {
			return err
		}
		ts.IsISO = isISO

		// Check if "seed" field exists
		if seedData, hasSeed := raw["seed"]; hasSeed {
			if err := json.Unmarshal(seedData, &ts.Seed); err != nil {
				return err
			}
		}

		// Check if "tiles" field exists
		if tilesData, hasTiles := raw["tiles"]; hasTiles {
			if err := json.Unmarshal(tilesData, &ts.Tiles); err != nil {
				return err
			}
		} else {
			// New format with isISO but tiles at root (shouldn't happen, but handle it)
			ts.Tiles = make(TilesetMap)
			for key, value := range raw {
				if key != "isISO" {
					var tile TileData
					if err := json.Unmarshal(value, &tile); err == nil {
						ts.Tiles[key] = tile
					}
				}
			}
		}
	} else {
		// Old format: tiles at root level, no isISO flag
		ts.IsISO = false
		// Check for seed in old format too
		if seedData, hasSeed := raw["seed"]; hasSeed {
			if err := json.Unmarshal(seedData, &ts.Seed); err != nil {
				// Seed parsing error is non-fatal
			}
		}
		if err := json.Unmarshal(data, &ts.Tiles); err != nil {
			return err
		}
	}

	return nil
}

// MarshalJSON writes the tileset in the new format
func (ts TilesetData) MarshalJSON() ([]byte, error) {
	result := map[string]interface{}{
		"isISO": ts.IsISO,
		"tiles": ts.Tiles,
	}
	if ts.Seed != 0 {
		result["seed"] = ts.Seed
	}
	return json.Marshal(result)
}

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
// isIsometric: if true, allows larger tiles (up to 128x128) for isometric 2.5D tiles
func ValidateTileData(tile *TileData, tileName string, isIsometric bool) error {
	// Validate dimensions
	// Isometric tiles need larger dimensions (64x48, 128x96, etc.) for 2.5D effect
	if isIsometric {
		// Isometric tiles: allow up to 128x128 (for 2.5D tiles with side faces)
		if tile.Width < 2 || tile.Height < 2 {
			return fmt.Errorf("tile '%s': dimensions must be at least 2x2, got %dx%d", tileName, tile.Width, tile.Height)
		}
		if tile.Width > 128 || tile.Height > 128 {
			return fmt.Errorf("tile '%s': isometric tile dimensions cannot exceed 128x128, got %dx%d", tileName, tile.Width, tile.Height)
		}
	} else {
		// Normal tiles: use standard sprite size validation (2x2 to 32x32)
		if err := ValidateSpriteSize(tile.Width, tile.Height, false); err != nil {
			return fmt.Errorf("tile '%s': %w", tileName, err)
		}
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

