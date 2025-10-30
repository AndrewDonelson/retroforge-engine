package config

import (
    "os"
    "strconv"
)

// Config holds engine-wide defaults.
type Config struct {
    ScreenWidth  int
    ScreenHeight int
    TargetFPS    int
    PaletteName  string
}

func Defaults() Config {
    return Config{
        ScreenWidth:  480,
        ScreenHeight: 270,
        TargetFPS:    60,
        PaletteName:  "RetroForge 50",
    }
}

// Load merges environment variables onto Defaults.
// Supported envs: RETROFORGE_WIDTH, RETROFORGE_HEIGHT, RETROFORGE_FPS, RETROFORGE_PALETTE
func Load() Config {
    c := Defaults()
    if v := getenvInt("RETROFORGE_WIDTH"); v > 0 { c.ScreenWidth = v }
    if v := getenvInt("RETROFORGE_HEIGHT"); v > 0 { c.ScreenHeight = v }
    if v := getenvInt("RETROFORGE_FPS"); v > 0 { c.TargetFPS = v }
    if v := os.Getenv("RETROFORGE_PALETTE"); v != "" { c.PaletteName = v }
    return c
}

func getenvInt(key string) int {
    s := os.Getenv(key)
    if s == "" { return 0 }
    n, _ := strconv.Atoi(s)
    return n
}


