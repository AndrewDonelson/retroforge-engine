package tile2iso

import (
	"errors"
	"testing"
)

func TestErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"ErrNilTexture", ErrNilTexture, "texture cannot be nil"},
		{"ErrInvalidDimensions", ErrInvalidDimensions, "invalid tile dimensions"},
		{"ErrInvalidLightingMode", ErrInvalidLightingMode, "unknown lighting mode"},
		{"ErrInvalidSpriteType", ErrInvalidSpriteType, "sprite type not supported"},
		{"ErrSpriteNotFound", ErrSpriteNotFound, "sprite not found"},
		{"ErrInvalidFrame", ErrInvalidFrame, "invalid frame"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.want {
				t.Errorf("error message = %q, want %q", tt.err.Error(), tt.want)
			}

			// Check that errors are properly comparable
			if !errors.Is(tt.err, tt.err) {
				t.Error("error should be equal to itself")
			}
		})
	}
}

