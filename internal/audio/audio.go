//go:build !js
package audio

import (
    "math"
    "sync"
    "time"
    "github.com/veandco/go-sdl2/sdl"
    "strings"
)

type voice struct {
    kind string // "sine", "noise", or "loop"
    freq float64
    gain float64
    phase float64
    tleft float64 // seconds left; <=0 for loop (e.g., thrust)
}

var (
    mu sync.Mutex
    voices []*voice
    initialized bool
    dev sdl.AudioDeviceID
)

func Init() error {
    if initialized { return nil }
    if err := sdl.InitSubSystem(sdl.INIT_AUDIO); err != nil { return err }
    // Use signed 16-bit output for broader device support
    want := &sdl.AudioSpec{Freq: 44100, Format: sdl.AUDIO_S16LSB, Channels: 1, Samples: 1024}
    var have sdl.AudioSpec
    d, err := sdl.OpenAudioDevice("", false, want, &have, 0)
    if err != nil { return err }
    dev = d
    sdl.PauseAudioDevice(dev, false)
    go mixerLoop()
    initialized = true
    return nil
}

// simple RNG
var seed uint32 = 1
func randFloat() float64 { seed = 1664525*seed + 1013904223; return float64(seed&0xFFFF)/65535.0 }

func PlaySine(freq, dur, gain float64) { mu.Lock(); defer mu.Unlock(); voices = append(voices, &voice{kind:"sine", freq:freq, gain:gain, tleft:dur}) }
func PlayNoise(dur, gain float64)     { mu.Lock(); defer mu.Unlock(); voices = append(voices, &voice{kind:"noise", gain:gain, tleft:dur}) }

// Thrust on/off: looped low buzz
var thrustOn bool
func Thrust(on bool) {
    mu.Lock(); defer mu.Unlock()
    if on && !thrustOn {
        voices = append(voices, &voice{kind:"loop", freq:110, gain:0.2, tleft:-1})
        thrustOn = true
    } else if !on && thrustOn {
        // stop only looped voices; keep one-shots playing
        n := voices[:0]
        for _, v := range voices {
            if v.kind == "loop" { continue }
            n = append(n, v)
        }
        voices = n
        thrustOn = false
    }
}

// StopAll immediately stops all current sounds and clears queued audio.
func StopAll() {
    mu.Lock()
    voices = []*voice{}
    mu.Unlock()
    if dev != 0 {
        sdl.ClearQueuedAudio(dev)
    }
    thrustOn = false
}

func mixerLoop() {
    ticker := time.NewTicker(20 * time.Millisecond)
    defer ticker.Stop()
    bufSamples := 44100 / 50 // ~20ms
    for range ticker.C {
        if dev == 0 { continue }
        // prevent queue from growing without bound (> 300ms)
        if sdl.GetQueuedAudioSize(dev) > 44100*2*3/10 { // bytes
            sdl.ClearQueuedAudio(dev)
        }
        mu.Lock()
        f32 := make([]float32, bufSamples)
        dt := 1.0 / 44100.0
        for _, v := range voices {
            for i := 0; i < bufSamples; i++ {
                var s float64
                if v.kind == "sine" {
                    s = math.Sin(2*math.Pi*v.phase) * v.gain
                    v.phase += v.freq * dt
                    if v.phase > 1 { v.phase -= 1 }
                } else {
                    s = (randFloat()*2 - 1) * v.gain
                }
                f32[i] += float32(s)
                if v.tleft > 0 { v.tleft -= dt }
            }
        }
        // cull finished (keep looped voices)
        n := voices[:0]
        for _, v := range voices {
            if v.tleft <= 0 && v.kind != "loop" { continue }
            n = append(n, v)
        }
        voices = n
        mu.Unlock()
        // queue to device as int16 PCM
        sdl.QueueAudio(dev, float32ToS16Bytes(f32))
    }
}

func float32ToS16Bytes(f []float32) []byte {
    b := make([]byte, len(f)*2)
    for i, v := range f {
        // clamp and convert
        if v > 1 { v = 1 } else if v < -1 { v = -1 }
        s := int16(v * 32767)
        b[i*2+0] = byte(s)
        b[i*2+1] = byte(uint16(s) >> 8)
    }
    return b
}

// --- Simple note parser and sequence player ---

var noteOffsets = map[string]int{
    "C":0, "C#":1, "D":2, "D#":3, "E":4, "F":5, "F#":6,
    "G":7, "G#":8, "A":9, "A#":10, "B":11,
}

func noteToFreq(note string, defaultOctave int) (float64, bool) {
    // Formats: "1G#2" (octave-note-len) or "G#2" (note-len) or "R2" (rest)
    s := strings.ToUpper(strings.TrimSpace(note))
    if len(s) == 0 { return 0, false }
    if s[0] == 'R' { return 0, true }
    octave := defaultOctave
    pos := 0
    if s[0] >= '0' && s[0] <= '9' { octave = int(s[0]-'0'); pos++ }
    if pos >= len(s) { return 0, false }
    n := string(s[pos])
    pos++
    if pos < len(s) && s[pos] == '#' { n += "#"; pos++ }
    off, ok := noteOffsets[n]
    if !ok { return 0, false }
    // Compute frequency using A4 reference
    semitone := (octave-4)*12 + (off-9)
    f := 440.0 * math.Pow(2, float64(semitone)/12.0)
    return f, true
}

// PlayNotes plays a sequence of tokens at bpm with given gain.
// Tokens like: 1G#2, G2, R1 (rest), default octave 4 if not prefixed.
func PlayNotes(tokens []string, bpm float64, gain float64) {
    if bpm <= 0 { bpm = 120 }
    beat := 60.0 / bpm
    go func() {
        for _, t := range tokens {
            s := strings.ToUpper(strings.TrimSpace(t))
            if s == "" { continue }
            // length = last digit if present
            length := 1
            if last := s[len(s)-1]; last >= '0' && last <= '9' {
                length = int(last-'0')
                s = s[:len(s)-1]
            }
            dur := float64(length) * beat
            if s == "R" { time.Sleep(time.Duration(dur*1000)*time.Millisecond); continue }
            if f, ok := noteToFreq(s, 4); ok {
                PlaySine(f, dur*0.95, gain)
                time.Sleep(time.Duration(dur*1000)*time.Millisecond)
            } else {
                // unknown token, skip beat
                time.Sleep(time.Duration(beat*1000)*time.Millisecond)
            }
        }
    }()
}


