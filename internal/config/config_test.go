package config

import "testing"

func TestDefaults(t *testing.T) {
    d := Defaults()
    if d.ScreenWidth != 480 || d.ScreenHeight != 270 || d.TargetFPS != 60 {
        t.Fatalf("unexpected defaults: %#v", d)
    }
}

func TestLoadFromEnv(t *testing.T) {
    t.Setenv("RETROFORGE_WIDTH", "800")
    t.Setenv("RETROFORGE_HEIGHT", "450")
    t.Setenv("RETROFORGE_FPS", "30")
    t.Setenv("RETROFORGE_PALETTE", "Neon 50")
    c := Load()
    if c.ScreenWidth != 800 || c.ScreenHeight != 450 || c.TargetFPS != 30 || c.PaletteName != "Neon 50" {
        t.Fatalf("env load failed: %#v", c)
    }
}


