# RetroForge Sound Tool Design Document

## TODO: Sound Tool (`soundtool`)

This document outlines the design for a new command-line tool and library for analyzing, validating, and processing RetroForge audio assets (SFX and Music).

---

## Executive Summary

### Core Capabilities
- **Duration Calculation**: Automatically calculate play length from music tokens and BPM
- **Validation**: Validate SFX and Music JSON files against RetroForge specifications
- **Loop Detection**: Identify and suggest loop points for music tracks
- **Format Conversion**: Convert between different audio formats if needed
- **Analysis**: Provide detailed statistics on audio assets

### Access Methods
- **CLI Tool** (`cmd/soundtool`): Command-line interface for batch processing
- **Go Package** (`internal/soundtool`): Reusable library for integration
- **Future**: Lua API bindings (optional, for runtime audio analysis)

---

## Architecture Overview

### Directory Structure
```
retroforge-engine/
├── cmd/
│   └── soundtool/
│       ├── main.go
│       └── commands/
│           ├── analyze.go      # Analyze audio files
│           ├── validate.go      # Validate JSON files
│           ├── duration.go      # Calculate durations
│           └── loop.go          # Find loop points
├── internal/
│   └── soundtool/
│       ├── soundtool.go         # Public API
│       ├── models.go            # Data structures
│       ├── calculator.go        # Duration calculation
│       ├── validator.go         # Validation logic
│       ├── loopfinder.go        # Loop detection
│       └── errors.go            # Error types
```

### Component Diagram
```
┌─────────────────┐
│   CLI (cobra)   │
└────────┬────────┘
         │
┌────────▼────────┐
│  soundtool API  │
└────────┬────────┘
         │
    ┌────┴────┬──────────┬──────────┐
    │         │          │          │
┌───▼───┐ ┌──▼───┐ ┌────▼────┐ ┌──▼────┐
│Calc.  │ │Valid.│ │Loop    │ │Analyze│
└───────┘ └──────┘ └─────────┘ └───────┘
```

### Interface Design Philosophy
- **Stateless**: Functions are pure where possible
- **Composable**: Small, focused functions that can be combined
- **Error-First**: Clear error messages with context

---

## Core Package API

### Public Interface

```go
package soundtool

// CalculateDuration calculates the total duration of a music track
// tokens: Array of music tokens (e.g., ["5E2", "5G#2", "R1"])
// bpm: Beats per minute
// Returns: Total beats and duration in seconds
func CalculateDuration(tokens []string, bpm float64) (totalBeats int, durationSeconds float64)

// ValidateMusic validates a music definition
func ValidateMusic(music *MusicDefinition) error

// ValidateSFX validates a sound effect definition
func ValidateSFX(sfx *SFXDefinition) error

// FindLoopPoint suggests optimal loop points in a music track
func FindLoopPoint(tokens []string, bpm float64) (startToken, endToken int, confidence float64)

// AnalyzeTrack provides detailed statistics about a music track
func AnalyzeTrack(music *MusicDefinition) *TrackAnalysis
```

### Options Structures

```go
type DurationOptions struct {
    BPM         float64 // Beats per minute (default: 120)
    Precision   int     // Decimal places for duration (default: 2)
}

type LoopOptions struct {
    MinLoopLength int     // Minimum loop length in beats (default: 4)
    MaxSearch     int     // Maximum tokens to search (default: 64)
    Tolerance     float64 // Musical similarity tolerance (default: 0.1)
}

type ValidationOptions struct {
    StrictMode bool   // Enable strict validation (default: false)
    Warnings   bool   // Show warnings for non-fatal issues (default: true)
}
```

---

## Data Models

### Core Go Structures

```go
// MusicDefinition matches cartio.MusicDefinition
type MusicDefinition struct {
    Tokens []string `json:"tokens"` // Array of note tokens
    BPM    float64  `json:"bpm"`    // Beats per minute
    Gain   float64  `json:"gain"`   // Volume/gain (0.0-1.0)
    Loop   bool     `json:"loop"`   // NEW: Whether track should loop (optional)
    Length float64  `json:"length"` // NEW: Cached duration in seconds (optional)
}

// SFXDefinition matches cartio.SFXDefinition
type SFXDefinition struct {
    Type     string  `json:"type"`     // "sine", "noise", "thrust", "stopall"
    Freq     float64 `json:"freq"`    // Frequency (Hz)
    Duration float64 `json:"duration"` // Duration in seconds
    Gain     float64 `json:"gain"`     // Volume/gain (0.0-1.0)
}

// TrackAnalysis provides detailed statistics
type TrackAnalysis struct {
    TotalBeats      int     `json:"totalBeats"`
    DurationSeconds float64 `json:"durationSeconds"`
    DurationMinutes float64 `json:"durationMinutes"`
    TokenCount      int     `json:"tokenCount"`
    NoteCount       int     `json:"noteCount"`
    RestCount       int     `json:"restCount"`
    AverageNoteLength float64 `json:"averageNoteLength"`
    LongestNote      string   `json:"longestNote"`
    ShortestNote     string   `json:"shortestNote"`
    HasLoopPoint     bool     `json:"hasLoopPoint"`
    SuggestedLoopStart int    `json:"suggestedLoopStart,omitempty"`
    SuggestedLoopEnd   int    `json:"suggestedLoopEnd,omitempty"`
}
```

### Validation Methods

```go
func (m *MusicDefinition) Validate() error
func (s *SFXDefinition) Validate() error

// Helper methods
func (m *MusicDefinition) CalculateDuration() (int, float64)
func (m *MusicDefinition) UpdateLength() // Calculates and sets Length field
```

---

## Duration Calculation

### Algorithm

The duration calculation follows this formula:
```
Duration (seconds) = (Total Beats / BPM) × 60
```

**Token Parsing:**
- Notes: `"5E2"` → octave 5, note E, duration 2 beats
- Rests: `"R1"`, `"R2"`, etc. → rest for N beats
- Default: If no duration specified, assume 1 beat

**Example:**
```json
{
  "tokens": ["5G1", "5C2", "5E2", "R1"],
  "bpm": 100
}
```

Calculation:
- `5G1` = 1 beat
- `5C2` = 2 beats
- `5E2` = 2 beats
- `R1` = 1 beat
- Total: 6 beats
- Duration: (6 / 100) × 60 = 3.6 seconds

### Implementation

```go
func CalculateDuration(tokens []string, bpm float64) (int, float64) {
    totalBeats := 0
    
    for _, token := range tokens {
        beats := parseTokenDuration(token)
        totalBeats += beats
    }
    
    durationSeconds := (float64(totalBeats) / bpm) * 60.0
    return totalBeats, durationSeconds
}

func parseTokenDuration(token string) int {
    // Handle rests
    if strings.HasPrefix(token, "R") {
        beats, _ := strconv.Atoi(token[1:])
        if beats == 0 {
            return 1 // Default rest
        }
        return beats
    }
    
    // Find last digit (duration)
    re := regexp.MustCompile(`(\d+)$`)
    match := re.FindString(token)
    if match != "" {
        beats, _ := strconv.Atoi(match)
        return beats
    }
    
    return 1 // Default to 1 beat
}
```

---

## Loop Property

### Design

The `loop` property in `MusicDefinition` indicates whether a track should automatically loop when played. This is a boolean flag:

```json
{
  "menu_music": {
    "tokens": ["5E2", "5G#2", "5B2", "..."],
    "bpm": 120,
    "gain": 0.25,
    "loop": true  // Track will loop automatically
  }
}
```

### Behavior
- `loop: true` - Music automatically restarts when finished
- `loop: false` or omitted - Music plays once and stops
- Default: `false` (backward compatible)

### Loop Point Detection (Future)

```go
// FindLoopPoint analyzes a track for natural loop points
// Looks for:
// 1. Repeated patterns (e.g., same sequence appears twice)
// 2. Musical resolution points (ends on tonic, major intervals)
// 3. Phrase boundaries (rests followed by similar patterns)
func FindLoopPoint(tokens []string, bpm float64) (start, end int, confidence float64)
```

---

## Length Property

### Design

The `length` property is a cached duration value in seconds. This allows:
- Quick lookup without recalculation
- JSON files to document track length
- Tools to detect when length doesn't match calculated value

```json
{
  "menu_music": {
    "tokens": ["5E2", "5G#2", "5B2", "..."],
    "bpm": 120,
    "gain": 0.25,
    "length": 63.0  // Cached duration in seconds
  }
}
```

### Auto-Calculation

Tools can automatically calculate and update the `length` field:
```go
func (m *MusicDefinition) UpdateLength() {
    _, duration := CalculateDuration(m.Tokens, m.BPM)
    m.Length = duration
}
```

### Validation

If `length` is present, validate it matches calculated duration:
```go
func (m *MusicDefinition) ValidateLength() error {
    _, calculated := CalculateDuration(m.Tokens, m.BPM)
    if m.Length > 0 && math.Abs(m.Length - calculated) > 0.01 {
        return fmt.Errorf("length mismatch: declared %.2f, calculated %.2f", 
            m.Length, calculated)
    }
    return nil
}
```

---

## CLI Tool

### Structure

```
soundtool
├── analyze <music.json>     # Analyze all tracks
├── duration <track>         # Calculate duration for a track
├── validate <file.json>     # Validate SFX/Music JSON
├── loop <track>             # Find loop points
└── update-lengths <file>    # Auto-calculate and update length fields
```

### Command Examples

```bash
# Calculate duration for a specific track
soundtool duration music.json menu_music

# Analyze all tracks in a music file
soundtool analyze music.json

# Validate music file
soundtool validate music.json

# Auto-calculate and update length fields
soundtool update-lengths music.json

# Find loop points for a track
soundtool loop music.json menu_music
```

### Usage

```bash
# Build
make soundtool-release

# Analyze music file
./soundtool analyze examples/tron-lightcycles/assets/music.json

# Output:
# menu_music:
#   Total beats: 126
#   Duration: 1m 3.00s (63.00 seconds)
#   Token count: 66
#   Suggested loop: tokens 0-65 (confidence: 0.85)
#
# taps:
#   Total beats: 57
#   Duration: 0m 34.20s (34.20 seconds)
#   Token count: 37
#   No loop point detected
```

---

## Validation Strategy

### Music Validation
- Tokens must be valid format (note or rest)
- BPM must be positive (> 0)
- Gain must be 0.0-1.0
- If `length` is present, must match calculated duration (within tolerance)
- Warn if track is too short (< 1 second) or too long (> 10 minutes)

### SFX Validation
- Type must be one of: "sine", "noise", "thrust", "stopall"
- Frequency must be >= 0 (0 allowed for noise/stopall)
- Duration must be positive (> 0)
- Gain must be 0.0-1.0

---

## Testing Strategy

### Unit Tests
- Duration calculation with various token formats
- Validation of valid and invalid inputs
- Loop detection algorithms
- Edge cases (empty tracks, single notes, very long tracks)

### Integration Tests
- Full analysis pipeline
- JSON file parsing and validation
- CLI command execution

### Coverage Requirements
- Minimum 80% code coverage
- 100% coverage for critical path (duration calculation, validation)

---

## Implementation Phases

### Phase 1: Core Library (Week 1)
- [ ] Data models (`models.go`)
- [ ] Duration calculation (`calculator.go`)
- [ ] Basic validation (`validator.go`)
- [ ] Error types (`errors.go`)

### Phase 2: CLI Tool (Week 1-2)
- [ ] Cobra command structure
- [ ] `analyze` command
- [ ] `duration` command
- [ ] `validate` command

### Phase 3: Advanced Features (Week 2)
- [ ] `length` property support
- [ ] `loop` property support
- [ ] `update-lengths` command
- [ ] Loop point detection (`loopfinder.go`)

### Phase 4: Polish (Week 2-3)
- [ ] Comprehensive tests
- [ ] Documentation
- [ ] Integration with build system
- [ ] Example usage in README

---

## Success Criteria

### Functional
- ✅ Accurately calculates duration from tokens and BPM
- ✅ Validates music and SFX JSON files
- ✅ Can update `length` fields automatically
- ✅ Supports `loop` property

### Quality
- ✅ Clear error messages
- ✅ Comprehensive test coverage
- ✅ Well-documented API

### Performance
- ✅ Analyzes 100+ tracks in < 1 second
- ✅ Memory efficient (no unnecessary allocations)

### User Experience
- ✅ Simple CLI interface
- ✅ Helpful output format
- ✅ Integration with existing workflow

---

## Dependencies

### Required Go Packages
- `github.com/spf13/cobra` - CLI framework (same as imgtool)
- Standard library: `regexp`, `strconv`, `strings`, `math`

### RetroForge Engine Packages
- `github.com/AndrewDonelson/retroforge-engine/internal/cartio` - For data structure compatibility

---

## API Documentation Template

### GoDoc Example
```go
// CalculateDuration calculates the total playback duration of a music track.
// 
// Parameters:
//   - tokens: Array of music tokens (e.g., ["5E2", "5G#2", "R1"])
//   - bpm: Beats per minute (must be > 0)
//
// Returns:
//   - totalBeats: Total number of beats in the track
//   - durationSeconds: Total duration in seconds
//
// Example:
//   beats, secs := CalculateDuration([]string{"5G1", "5C2", "R1"}, 100)
//   // beats = 4, secs = 2.4
func CalculateDuration(tokens []string, bpm float64) (totalBeats int, durationSeconds float64)
```

---

## Error Handling

### Error Types

```go
type SoundToolError struct {
    Code    int
    Message string
    Cause   error
}

var (
    ErrInvalidToken    = &SoundToolError{Code: 1000, Message: "invalid token format"}
    ErrInvalidBPM      = &SoundToolError{Code: 1001, Message: "BPM must be > 0"}
    ErrInvalidGain     = &SoundToolError{Code: 1002, Message: "gain must be 0.0-1.0"}
    ErrDurationMismatch = &SoundToolError{Code: 1003, Message: "declared length doesn't match calculated"}
)
```

---

## Acceptance Checklist

- [ ] Duration calculation is accurate for all test cases
- [ ] Validation catches all invalid inputs
- [ ] CLI tool is easy to use
- [ ] `length` property is automatically calculated
- [ ] `loop` property is supported
- [ ] Documentation is complete
- [ ] Tests pass with 80%+ coverage
- [ ] Integration with Makefile
- [ ] Example usage in README

---

## Document Control

**Status**: TODO / Design Phase  
**Version**: 0.1.0  
**Last Updated**: 2024-11-02  
**Author**: RetroForge Team  
**Related**: See `imgtool` design document for similar structure

---

## Notes

- This tool is inspired by the `imgtool` design and follows similar patterns
- Duration calculation formula: `(Total Beats / BPM) × 60 = seconds`
- `length` and `loop` properties are optional for backward compatibility
- Future enhancements may include audio visualization and MIDI export

