package cartio

import "fmt"

// SFXMap maps sound effect names to token arrays
// Format: {"soundname": ["S044D04G3"]} or {"soundname": ["S044D04G3", "Q220D02G2"]}
// Legacy format is still supported for backward compatibility (see GetSFXTokens)
type SFXMap map[string]interface{}

// GetSFXTokens extracts token array from SFXMap for a given sound name
// Handles both new format (array of tokens) and legacy format (object with type/freq/duration/gain)
func GetSFXTokens(sfxMap SFXMap, name string) []string {
	value, ok := sfxMap[name]
	if !ok {
		return nil
	}
	
	switch v := value.(type) {
	case []interface{}:
		// New format: array of tokens
		tokens := make([]string, 0, len(v))
		for _, tok := range v {
			if str, ok := tok.(string); ok {
				tokens = append(tokens, str)
			}
		}
		return tokens
	case []string:
		// Direct string array
		return v
	case map[string]interface{}:
		// Legacy format: convert to tokens
		return legacyToTokens(v)
	}
	
	return nil
}

// legacyToTokens converts legacy SFX format to tokens
func legacyToTokens(legacy map[string]interface{}) []string {
	sfxType, _ := legacy["type"].(string)
	
	if sfxType == "stopall" {
		return []string{"STOPALL"}
	}
	
	freq, _ := legacy["freq"].(float64)
	if freq < 0 {
		freq = 0
	}
	duration, _ := legacy["duration"].(float64)
	if duration < 0 {
		duration = 0
	}
	gain, _ := legacy["gain"].(float64)
	if gain < 0 {
		gain = 0
	}
	if gain > 1 {
		gain = 1
	}
	
	// Convert to token format
	freqVal := int(freq / 100)
	if freqVal > 127 {
		freqVal = 127
	}
	durVal := int(duration * 16)
	if durVal < 1 {
		durVal = 1
	}
	if durVal > 16 {
		durVal = 16
	}
	gainVal := int(gain * 8)
	if gainVal < 1 {
		gainVal = 1
	}
	if gainVal > 8 {
		gainVal = 8
	}
	
	var code string
	switch sfxType {
	case "sine":
		code = "S"
	case "noise":
		code = "N"
	case "thrust":
		code = "L"
	default:
		code = "S"
	}
	
	return []string{fmt.Sprintf("%s%03dD%02dG%d", code, freqVal, durVal, gainVal)}
}
