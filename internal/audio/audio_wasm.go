//go:build js && wasm

package audio

import "syscall/js"

func Init() error { return nil }
func Close()        {}
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
func PlaySFX(name string, args ...string)                 {}
func StopSFX(name string)                                 {}
func StopAll() {
    if fn := js.Global().Get("rf_audio_stopAll"); fn.Truthy() {
        fn.Invoke()
    }
}
func PlayMusic(tokens []string, bpm int, gain float64)    {}
func PlayNotes(tokens []string, bpm float64, gain float64) {}
func Thrust(on bool) {
    // Optional: low buzz can be emulated via JS if desired
    if !on { StopAll() }
}


