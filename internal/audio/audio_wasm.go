//go:build js && wasm

package audio

import "syscall/js"

func Init() error                                         { return nil }
func Close()                                              {}
func PlayTone(freq float64, durSec float64, gain float64) { PlaySine(freq, durSec, gain) }
func PlaySine(freq float64, durSec float64, gain float64) {
	if fn := js.Global().Get("rf_audio_playSine"); fn.Truthy() {
		fn.Invoke(freq, durSec, gain)
	}
}
func PlayNoise(durSec float64, gain float64) {
	if fn := js.Global().Get("rf_audio_playNoise"); fn.Truthy() {
		fn.Invoke(durSec, gain)
	}
}
func PlaySFX(name string, args ...string) {}
func StopSFX(name string)                 {}
func StopAll() {
	if fn := js.Global().Get("rf_audio_stopAll"); fn.Truthy() {
		fn.Invoke()
	}
}
func PlayMusic(tokens []string, bpm int, gain float64) {
	PlayNotes(tokens, float64(bpm), gain)
}
func PlayNotes(tokens []string, bpm float64, gain float64) {
	if fn := js.Global().Get("rf_audio_playNotes"); fn.Truthy() {
		arr := js.Global().Get("Array").New(len(tokens))
		for i, tok := range tokens {
			arr.SetIndex(i, tok)
		}
		fn.Invoke(arr, bpm, gain)
	}
}
func Thrust(on bool) {
	if fn := js.Global().Get("rf_audio_thrust"); fn.Truthy() {
		fn.Invoke(on)
	} else if !on {
		StopAll()
	}
}
