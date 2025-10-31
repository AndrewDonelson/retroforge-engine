//go:build !js && !wasm

package main

import (
	"bytes"
	_ "embed"
	"flag"

	"github.com/AndrewDonelson/retroforge-engine/internal/engine"
	"github.com/AndrewDonelson/retroforge-engine/internal/sdlrun"
)

// This file is provided as a placeholder and gets overwritten during bundle builds.
// The linter may complain but the file exists at compile time.
//
//go:embed cart.rf
var cartBytes []byte

func main() {
	scale := flag.Int("scale", 3, "window scale")
	flag.Parse()

	// Safety check: ensure embedded cart data exists
	if len(cartBytes) == 0 {
		panic("embedded cart.rf is empty - this should not happen. Ensure the file exists during build.")
	}

	e := engine.New(60)
	defer e.Close()
	if err := e.LoadCartFromReader(bytes.NewReader(cartBytes), int64(len(cartBytes))); err != nil {
		panic(err)
	}
	if err := sdlrun.RunWindow(e, *scale); err != nil {
		panic(err)
	}
}
