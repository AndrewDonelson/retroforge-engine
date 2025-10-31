package cartio

// MusicDefinition represents a music track
type MusicDefinition struct {
	Tokens []string `json:"tokens"` // Array of note tokens like ["4C1","4E1","R1"]
	BPM    float64  `json:"bpm"`    // Default BPM (optional, can be overridden)
	Gain   float64  `json:"gain"`   // Default gain/volume (optional, can be overridden)
}

// MusicMap maps music track names to their definitions
type MusicMap map[string]MusicDefinition
