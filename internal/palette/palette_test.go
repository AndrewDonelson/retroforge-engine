package palette

import "testing"

func TestNewAndValidate(t *testing.T) {
    // Build 50-color palette: black, white, then 16 hues × 3 shades (fake values, valid hex)
    colors := make([]string, 0, 50)
    colors = append(colors, "#000000", "#ffffff")
    hues := []string{"#ff4d4d", "#ff914d", "#ffd84d", "#b6ff4d", "#4dd487", "#36d8c7", "#4dd5ff", "#66bfff", "#6f88ff", "#8a75ff", "#b478ff", "#ff6fb1", "#ff7fa0", "#a8795a", "#a0b15a", "#38bdf8"}
    for _, base := range hues {
        // we don't compute shades here; just repeat base to satisfy size and hex validation
        colors = append(colors, base, base, base)
    }

    if len(colors) != 50 {
        t.Fatalf("expected 50, got %d", len(colors))
    }

    p, err := New("test", colors)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if err := p.IsValid(); err != nil {
        t.Fatalf("palette invalid: %v", err)
    }
}

func TestInvalidHex(t *testing.T) {
    colors := make([]string, 50)
    for i := range colors {
        colors[i] = "#000000"
    }
    colors[10] = "nothex"
    if _, err := New("bad", colors); err == nil {
        t.Fatalf("expected error for invalid hex")
    }
}


