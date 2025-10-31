package cartio

// SFXDefinition represents a single sound effect
type SFXDefinition struct {
	Type     string  `json:"type"`     // "sine", "noise", "thrust", "stopall"
	Freq     float64 `json:"freq"`     // Frequency (Hz) for sine/thrust, 0 for noise
	Duration float64 `json:"duration"` // Duration in seconds (0 for looped sounds like thrust)
	Gain     float64 `json:"gain"`     // Gain/volume (0.0 to 1.0, typically 0.2-0.4)
}

// SFXMap maps sound effect names to their definitions
type SFXMap map[string]SFXDefinition
