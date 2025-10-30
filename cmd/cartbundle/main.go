package main

import (
    "bytes"
    _ "embed"
    "flag"
    "github.com/AndrewDonelson/retroforge-engine/internal/engine"
    "github.com/AndrewDonelson/retroforge-engine/internal/sdlrun"
)

//go:embed cart.rf
var cartBytes []byte

func main() {
    scale := flag.Int("scale", 3, "window scale")
    flag.Parse()
    e := engine.New(60)
    defer e.Close()
    if err := e.LoadCartFromReader(bytes.NewReader(cartBytes), int64(len(cartBytes))); err != nil { panic(err) }
    if err := sdlrun.RunWindow(e, *scale); err != nil { panic(err) }
}


