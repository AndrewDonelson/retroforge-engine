package tile2iso

import "errors"

var (
	// ErrNilTexture indicates a nil texture was provided
	ErrNilTexture = errors.New("texture cannot be nil")
	// ErrInvalidDimensions indicates invalid tile dimensions
	ErrInvalidDimensions = errors.New("invalid tile dimensions")
	// ErrInvalidLightingMode indicates an unknown lighting mode
	ErrInvalidLightingMode = errors.New("unknown lighting mode")
	// ErrInvalidSpriteType indicates sprite type is not supported for this operation
	ErrInvalidSpriteType = errors.New("sprite type not supported")
	// ErrSpriteNotFound indicates a required sprite was not found
	ErrSpriteNotFound = errors.New("sprite not found")
	// ErrInvalidFrame indicates a frame was not found or invalid
	ErrInvalidFrame = errors.New("invalid frame")
)

