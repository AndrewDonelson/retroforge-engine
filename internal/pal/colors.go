package pal

import "image/color"

// Default50: index 0=black, 1=white, rest simple gradient placeholders.
var Default50 = func() []color.RGBA {
    out := make([]color.RGBA, 50)
    out[0] = color.RGBA{0,0,0,255}
    out[1] = color.RGBA{255,255,255,255}
    // fill remaining with a repeating pattern
    idx := 2
    for i := 0; i < 16 && idx < 50; i++ {
        base := uint8((i*15)%256)
        for s := 0; s < 3 && idx < 50; s++ {
            out[idx] = color.RGBA{base, uint8((int(base)+40)%256), uint8((int(base)+80)%256), 255}
            idx++
        }
    }
    for idx < 50 {
        v := uint8((idx*5)&255)
        out[idx] = color.RGBA{v,v,v,255}
        idx++
    }
    return out
}()

type Manager struct {
    current []color.RGBA
}

func NewManager() *Manager { return &Manager{current: append([]color.RGBA{}, Default50...)} }

func (m *Manager) Color(i int) color.RGBA {
    if i < 0 || i >= len(m.current) { return m.current[0] }
    return m.current[i]
}

func (m *Manager) Set(name string) {
    // TODO: support multiple named palettes; for now only default
    _ = name
    m.current = append([]color.RGBA{}, Default50...)
}


