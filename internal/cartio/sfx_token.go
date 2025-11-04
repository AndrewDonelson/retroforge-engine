package cartio

import (
	"fmt"
	"regexp"
	"strconv"
)

// Token format: [SQTWNL](\d{3})D(\d+)G(\d) or R(\d+)
// Examples:
//   S044D04G3 = Sine, 4400Hz (044*100), 0.25s (4/16), gain 0.375 (3/8)
//   Q220D02G2 = Square, 22000Hz, 0.125s, gain 0.25
//   R02 = Rest for 2/16 = 0.125s

var (
	tokenPattern = regexp.MustCompile(`^([SQTWNL])(\d{3})D(\d+)G(\d)$`)
	restPattern  = regexp.MustCompile(`^R(\d+)$`)
	stopallPattern = regexp.MustCompile(`^STOPALL$`)
)

// ParseToken parses a token string and returns waveform parameters
// Returns: type, frequency (Hz), duration (seconds), gain (0.0-1.0), error
func ParseToken(token string) (waveType string, freq float64, duration float64, gain float64, err error) {
	// Check for stopall token
	if stopallPattern.MatchString(token) {
		return "stopall", 0, 0, 0, nil
	}
	
	// Check for rest token
	if restPattern.MatchString(token) {
		matches := restPattern.FindStringSubmatch(token)
		if len(matches) != 2 {
			return "", 0, 0, 0, fmt.Errorf("invalid rest token: %s", token)
		}
		durSteps, _ := strconv.Atoi(matches[1])
		return "rest", 0, float64(durSteps) / 16.0, 0, nil
	}

	// Parse waveform token
	matches := tokenPattern.FindStringSubmatch(token)
	if len(matches) != 5 {
		return "", 0, 0, 0, fmt.Errorf("invalid token format: %s", token)
	}

	typeCode := matches[1]
	freqStr := matches[2]
	durStr := matches[3]
	gainStr := matches[4]

	// Parse frequency (0-127, multiply by 100 to get Hz)
	freqVal, err := strconv.Atoi(freqStr)
	if err != nil {
		return "", 0, 0, 0, fmt.Errorf("invalid frequency in token: %s", token)
	}
	freq = float64(freqVal) * 100.0

	// Parse duration (1-16, divide by 16 to get seconds)
	durVal, err := strconv.Atoi(durStr)
	if err != nil {
		return "", 0, 0, 0, fmt.Errorf("invalid duration in token: %s", token)
	}
	duration = float64(durVal) / 16.0

	// Parse gain (1-8, divide by 8 to get 0.0-1.0)
	gainVal, err := strconv.Atoi(gainStr)
	if err != nil {
		return "", 0, 0, 0, fmt.Errorf("invalid gain in token: %s", token)
	}
	gain = float64(gainVal) / 8.0

	// Map type code to waveform type
	switch typeCode {
	case "S":
		waveType = "sine"
	case "Q":
		waveType = "sine" // Square - map to sine for now (engine can add square support later)
	case "T":
		waveType = "sine" // Triangle - map to sine for now
	case "W":
		waveType = "sine" // Sawtooth - map to sine for now
	case "N":
		waveType = "noise"
	case "L":
		waveType = "thrust" // Loop = thrust (looped sound)
	default:
		return "", 0, 0, 0, fmt.Errorf("unknown waveform type code: %s", typeCode)
	}

	return waveType, freq, duration, gain, nil
}

