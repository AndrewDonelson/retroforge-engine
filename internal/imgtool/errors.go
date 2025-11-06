package imgtool

import "fmt"

// Error codes
const (
	ErrCodeInvalidImage      = 1000
	ErrCodeInvalidPalette    = 1001
	ErrCodeInvalidDimensions = 1002
	ErrCodeProcessingFailed  = 1003
	ErrCodeInvalidOptions    = 1004
)

// ImgToolError is a custom error type for imgtool operations
type ImgToolError struct {
	Code    int
	Message string
	Cause   error
}

func (e *ImgToolError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func (e *ImgToolError) Unwrap() error {
	return e.Cause
}

// Predefined errors
var (
	ErrInvalidPNG = &ImgToolError{
		Code:    ErrCodeInvalidImage,
		Message: "invalid PNG file",
	}

	ErrInvalidPaletteSize = &ImgToolError{
		Code:    ErrCodeInvalidPalette,
		Message: "game palette must have exactly 48 colors (built-in colors 0-15 are separate)",
	}

	ErrInvalidDimensions = &ImgToolError{
		Code:    ErrCodeInvalidDimensions,
		Message: "invalid image dimensions",
	}

	ErrColorQuantization = &ImgToolError{
		Code:    ErrCodeProcessingFailed,
		Message: "color quantization failed",
	}

	ErrScalingFailed = &ImgToolError{
		Code:    ErrCodeProcessingFailed,
		Message: "image scaling failed",
	}

	ErrInvalidHexFormat = &ImgToolError{
		Code:    ErrCodeInvalidOptions,
		Message: "invalid hex color format (expected #RRGGBB)",
	}

	ErrInvalidPixelData = &ImgToolError{
		Code:    ErrCodeInvalidDimensions,
		Message: "pixel data dimensions do not match sprite dimensions",
	}
)

// NewInvalidImageError creates a new invalid image error
func NewInvalidImageError(cause error) error {
	return &ImgToolError{
		Code:    ErrCodeInvalidImage,
		Message: "invalid image",
		Cause:   cause,
	}
}

// NewProcessingError creates a new processing error
func NewProcessingError(msg string, cause error) error {
	return &ImgToolError{
		Code:    ErrCodeProcessingFailed,
		Message: msg,
		Cause:   cause,
	}
}

