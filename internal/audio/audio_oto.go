//go:build !js && !sdl

package audio

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/ebitengine/oto/v3"
)

var (
	otoCtx      *oto.Context
	otoPlayer   *oto.Player
	audioMu     sync.Mutex
	initialized bool
	audioChan   chan []byte
)

// audioSource implements io.Reader for oto
type audioSource struct {
	data chan []byte
	buf  []byte
	pos  int
}

func (s *audioSource) Read(p []byte) (int, error) {
	if len(s.buf) == 0 {
		select {
		case s.buf = <-s.data:
			s.pos = 0
		default:
			// No data available, return silence
			for i := range p {
				p[i] = 0
			}
			return len(p), nil
		}
	}

	n := copy(p, s.buf[s.pos:])
	s.pos += n
	if s.pos >= len(s.buf) {
		s.buf = nil
		s.pos = 0
	}

	// Fill remainder with silence if needed
	if n < len(p) {
		for i := n; i < len(p); i++ {
			p[i] = 0
		}
	}
	return len(p), nil
}

func Init() error {
	audioMu.Lock()
	defer audioMu.Unlock()
	if initialized {
		return nil
	}

	// Initialize oto context
	op := &oto.NewContextOptions{
		SampleRate:   44100,
		ChannelCount: 1,
		Format:       oto.FormatSignedInt16LE,
	}

	var err error
	otoCtx, ready, err := oto.NewContext(op)
	if err != nil {
		return err
	}
	<-ready

	// Create audio source with buffered channel
	audioChan = make(chan []byte, 10) // Buffer 10 chunks
	source := &audioSource{data: audioChan}

	// Create player from source
	otoPlayer = otoCtx.NewPlayer(source)
	if otoPlayer == nil {
		return fmt.Errorf("failed to create oto player")
	}

	// Start playing
	otoPlayer.Play()

	go mixerLoop()
	initialized = true
	return nil
}

func mixerLoop() {
	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()
	bufSamples := 44100 / 50 // ~20ms

	for range ticker.C {
		if otoPlayer == nil {
			continue
		}

		mu.Lock()
		f32 := make([]float32, bufSamples)
		dt := 1.0 / 44100.0
		for _, v := range voices {
			for i := 0; i < bufSamples; i++ {
				var s float64
				if v.kind == "sine" || v.kind == "loop" {
					s = math.Sin(2*math.Pi*v.phase) * v.gain
					v.phase += v.freq * dt
					if v.phase > 1 {
						v.phase -= 1
					}
				} else {
					s = (randFloat()*2 - 1) * v.gain
				}
				f32[i] += float32(s)
				if v.tleft > 0 {
					v.tleft -= dt
				}
			}
		}
		// Cull finished (keep looped voices)
		n := voices[:0]
		for _, v := range voices {
			if v.tleft <= 0 && v.kind != "loop" {
				continue
			}
			n = append(n, v)
		}
		voices = n
		mu.Unlock()

		// Convert to int16 PCM and send to audio channel
		samples := float32ToS16Bytes(f32)
		if len(samples) > 0 && audioChan != nil {
			select {
			case audioChan <- samples:
			default:
				// Channel full, skip this chunk
			}
		}
	}
}

