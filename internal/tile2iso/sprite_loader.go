package tile2iso

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
)

// LoadSpritesFromFile loads sprites from a sprites.json file
func LoadSpritesFromFile(path string) (cartio.SpriteMap, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read sprites file: %w", err)
	}

	var spriteMap cartio.SpriteMap
	if err := json.Unmarshal(data, &spriteMap); err != nil {
		return nil, fmt.Errorf("failed to parse sprites JSON: %w", err)
	}

	// Normalize and validate all sprites
	for spriteName, sprite := range spriteMap {
		cartio.NormalizeSpriteData(&sprite)
		if err := cartio.ValidateSpriteData(&sprite, spriteName); err != nil {
			return nil, fmt.Errorf("sprite '%s' validation error: %w", spriteName, err)
		}
		spriteMap[spriteName] = sprite
	}

	return spriteMap, nil
}

// GetSpritePixels extracts pixel data from a sprite, handling all sprite types
// For static sprites, frameName should be empty or "default"
// For frames/animation sprites, frameName should be the frame name
func GetSpritePixels(sprite *cartio.SpriteData, frameName string) ([][]int, error) {
	if sprite == nil {
		return nil, ErrNilTexture
	}

	switch sprite.Type {
	case cartio.SpriteTypeStatic:
		if frameName != "" && frameName != "default" {
			return nil, fmt.Errorf("static sprite does not support frame name '%s'", frameName)
		}
		return sprite.Pixels, nil

	case cartio.SpriteTypeFrames, cartio.SpriteTypeAnimation:
		if frameName == "" {
			// Default to first frame if no frame name provided
			if len(sprite.Frames) == 0 {
				return nil, ErrInvalidFrame
			}
			frameName = sprite.Frames[0].Name
		}

		pixels, err := sprite.GetFramePixels(frameName)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidFrame, err)
		}
		return pixels, nil

	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidSpriteType, sprite.Type)
	}
}

// GetSpriteFromMap retrieves a sprite from a sprite map by name
func GetSpriteFromMap(spriteMap cartio.SpriteMap, spriteName string) (*cartio.SpriteData, error) {
	sprite, exists := spriteMap[spriteName]
	if !exists {
		return nil, fmt.Errorf("%w: '%s'", ErrSpriteNotFound, spriteName)
	}
	return &sprite, nil
}
