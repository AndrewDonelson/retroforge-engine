# RetroForge - Final Specification (v2.1)

**Date:** October 29, 2025  
**Status:** FINAL - Ready for Implementation  
**Last Updated:** Resolution and Platform Changes

---

## 🎯 Final Specifications

### Display
- **Resolution:** 480×270 (16:9 landscape) or 270×480 (9:16 portrait)
- **Aspect Ratio:** 16:9 (modern device optimized)
- **Scaling:** Integer multiples for pixel-perfect display
  - 1x = 480×270 (native)
  - 2x = 960×540 (qHD)
  - 3x = 1440×810
  - 4x = 1920×1080 (Full HD - perfect!)
- **Colors:** 50-color palettes (Black + White + 16 colors × 3 shades)

### Platform Support
✅ **Windows** - Single .exe executable  
✅ **macOS** - .app bundle  
✅ **Linux** - Single binary  
✅ **Android** - .apk via gomobile  
✅ **Web** - WASM in browser  
❌ **iOS** - Not supported (cost and restrictions)

### Core Features
- **Code:** 16,384 Lua tokens
- **Sprites:** 256 slots (4×4, 8×8, 16×16, 32×32 pixels)
- **Map:** 8 layers with parallax scrolling
- **Audio:** 8 channels, chip-tune synthesis (5 waveforms)
- **Physics:** Box2D rigid body simulation
- **Node System:** Godot-style scene graph (15+ node types)
- **API:** Dual (high-level nodes + low-level direct)

---

## 📝 Recent Changes (v2.1)

### Resolution Update
**Before:** 320×240 (4:3)  
**After:** 480×270 (16:9)

**Rationale:**
1. **Modern devices are 16:9** - Phones, tablets, monitors, TVs
2. **Perfect HD scaling** - 4x scaling = 1920×1080 exactly
3. **More screen space** - 50% more pixels (129,600 vs 76,800)
4. **Better for UI** - Wider aspect for HUD elements
5. **Landscape/Portrait modes** - 480×270 or 270×480

### Platform Addition
**Added:** Android support via gomobile  
**Excluded:** iOS (by design)

**Android Benefits:**
- Free development (no expensive account)
- Open ecosystem (side-loading allowed)
- Large install base (billions of devices)
- Good Go support (gomobile works well)

**iOS Exclusion Reasons:**
- $99/year developer account required
- Strict App Store policies
- Complex code signing
- Limited testing without paid account
- Not worth the cost/hassle for free project

---

## 🏗️ Architecture Summary

### Engine (Go 1.23)
```
Core Systems:
├── Scene Graph (Godot-style nodes)
├── Physics Engine (Box2D-go)
├── Audio System (3-tier: SoundManager, AudioPlayer, MusicPlayer)
├── Graphics (SDL2/OpenGL/OpenGL ES)
├── Input (Keyboard, Gamepad, Touch)
└── Lua VM (gopher-lua)

Build Targets:
├── Desktop (Go std) → .exe, .app, binary
├── Android (gomobile) → .apk
└── Web (TinyGo) → .wasm
```

### Web App (Next.js 16)
```
Features:
├── Code Editor (react-ace)
├── Sprite Editor (Canvas-based)
├── Map Editor (8 layers)
├── Audio Editor (Piano roll)
├── Cart Manager (Export .rfs/.rfe)
└── Live Preview (WASM runtime)
```

---

## 🎮 Node System

15+ built-in node types:

**Core Nodes:**
- Node, Node2D, PhysicsBody2D

**Visual Nodes:**
- Sprite, AnimatedSprite, TileMap, ParticleEmitter

**Physics Nodes:**
- StaticBody, RigidBody, KinematicBody, CollisionShape

**Utility Nodes:**
- Camera, AudioPlayer, MusicPlayer, Timer

**UI Nodes:**
- CanvasLayer

---

## 🎵 Audio System

**Three-Tier Architecture:**

1. **SoundManager** (Automatic)
   - Exists in every game
   - Global volume controls
   - Audio ducking
   - Channel allocation

2. **AudioPlayer** (Node)
   - Sound effects
   - Multiple simultaneous
   - Spatial audio

3. **MusicPlayer** (Node)
   - Music tracks
   - Single track at a time
   - Automatic cross-fading

**Synthesis:**
- Waveforms: Square, Triangle, Sawtooth, Sine, Noise
- Envelope: ADSR
- Effects: Vibrato, Arpeggio, Pitch slide

---

## 📦 Distribution

### File Formats
- **`.rfs`** - RetroForge Source (development)
- **`.rfe`** - RetroForge Executable (distribution)

Both are ZIP archives; .rfe is compressed and signed.

### Build Matrix

| Platform | Command | Output | Size Est. |
|----------|---------|---------|-----------|
| Windows | `go build` | retroforge.exe | ~15MB |
| macOS | `go build` | RetroForge.app | ~15MB |
| Linux | `go build` | retroforge | ~15MB |
| Android | `gomobile bind` | retroforge.apk | ~20MB |
| Web | `tinygo build` | retroforge.wasm | ~2MB |

---

## 🗓️ Development Timeline

**Total Time:** ~14 weeks (3.5 months)

**Phase 1:** Engine Core (Weeks 1-4)
- Node system
- Physics (Box2D)
- Basic rendering
- Input handling

**Phase 2:** Audio & Nodes (Weeks 5-7)
- 3-tier audio system
- Advanced nodes
- Particle effects

**Phase 3:** Web App (Weeks 8-10)
- Next.js setup
- Editors (code, sprite, map, audio)
- WASM integration

**Phase 4:** Polish (Weeks 11-12)
- Example carts
- Documentation
- Performance optimization

**Phase 5:** Android (Weeks 13-14)
- gomobile setup
- Touch input
- APK builds
- Testing

---

## 📐 Technical Specifications

### Constraints
- **Code:** 16,384 tokens (2x PICO-8)
- **Sprites:** 256 slots
- **Map:** 128×128 tiles (configurable)
- **Memory:** 64MB WASM heap
- **Frame Budget:** 16.67ms (60 FPS)

### Performance Targets
- **Frame Rate:** Stable 60 FPS
- **Physics:** 100+ active bodies
- **Nodes:** 500+ in scene
- **Audio:** 8 simultaneous channels
- **Load Time:** <2 seconds

---

## 🎨 Example Code

### High-Level (Node API)
```lua
function _init()
  -- Create physics player
  player = RigidBody.new({
    position = vec2(240, 135),
    width = 16,
    height = 16
  })
  
  -- Add sprite
  local sprite = Sprite.new({sprite_index = 1})
  player:add_child(sprite)
  
  Scene:add_child(player)
  
  -- Start music
  Sound:play_music(1, 2.0)  -- Track 1, 2s crossfade
end

function player:_update()
  if Input.is_action_just_pressed("jump") then
    self:apply_impulse(vec2(0, -300))
    Sound:play_sfx(5)  -- Jump sound
  end
end
```

### Low-Level (Direct API)
```lua
x, y = 240, 135
vy = 0

function _update()
  -- Manual physics
  vy = vy + 0.5
  y = y + vy
  
  if btn(4) and y > 220 then
    vy = -10
    sfx(5)
  end
end

function _draw()
  cls(0)
  spr(1, x, y)
end
```

---

## 📚 Documentation Files

All documents updated to v2.1:

1. **RETROFORGE_DESIGN.md** (48KB)
   - Complete technical specification
   - Full API reference
   - All systems documented

2. **RETROFORGE_NODE_ARCHITECTURE.md** (15KB)
   - Node system deep dive
   - Physics integration details
   - Example implementations

3. **RETROFORGE_KICKOFF.md** (17KB)
   - 14-week implementation plan
   - Phase 5 added (Android)
   - Task checklists

4. **RETROFORGE_BRANDING.md** (9KB)
   - Brand identity
   - Marketing strategy
   - Community plan

5. **RETROFORGE_QUICK_REF.md** (4KB)
   - One-page overview
   - Quick specs reference

6. **RETROFORGE_INDEX.md** (5KB)
   - Master navigation
   - Document guide

7. **RETROFORGE_REVISION_SUMMARY.md** (12KB)
   - All changes explained
   - Before/after comparisons

8. **RETROFORGE_FINAL_SPEC.md** (This file)
   - Final specifications
   - Recent changes summary

---

## ✅ Design Checklist

**Core Design:**
- [x] Project name (RetroForge)
- [x] File extensions (.rfs, .rfe)
- [x] Display resolution (480×270, 16:9)
- [x] Platform support (Win, Mac, Linux, Android, Web)
- [x] Node system (Godot-style)
- [x] Physics engine (Box2D)
- [x] Audio system (3-tier)
- [x] API design (dual level)
- [x] Cart format (ZIP-based)

**Documentation:**
- [x] Design document
- [x] Architecture guide
- [x] Implementation plan
- [x] Branding guide
- [x] Quick reference
- [x] Index/navigation
- [x] Revision history
- [x] Final spec (this)

**Ready for Development:**
- [x] Tech stack finalized
- [x] All dependencies identified
- [x] Build tools determined
- [x] Examples written
- [x] Timeline established

---

## 🚀 Next Steps

### Immediate (This Week)
1. Purchase domains
   - retroforge.dev
   - retroforge.com
2. Create GitHub organization
3. Set up repositories
4. Install development tools

### Week 1 Start
1. Go project initialization
2. SDL2 setup
3. Base Node class
4. Simple rendering test

### Milestone 1 (Week 4)
- Physics demo cart working
- Boxes fall, player jumps
- Camera follows player

---

## 🎯 Success Criteria

### Technical
- ✅ Runs at stable 60 FPS
- ✅ All platforms build successfully
- ✅ Node system works as designed
- ✅ Physics behaves realistically
- ✅ Audio plays without latency

### User Experience
- ✅ Easy to learn (clear docs)
- ✅ Quick to prototype (node system)
- ✅ Fun to use (immediate feedback)
- ✅ Sharable carts (repository)
- ✅ Cross-platform (play anywhere)

### Community
- ✅ Active Discord server
- ✅ Regular game jams
- ✅ Example cart library
- ✅ Tutorial series
- ✅ Featured creators

---

## 📊 Comparison with PICO-8

| Feature | PICO-8 | RetroForge | Difference |
|---------|---------|------------|------------|
| Resolution | 128×128 | 480×270 | 6x more pixels |
| Aspect | 1:1 | 16:9 | Modern displays |
| Colors | 16 | 50 | 3x more colors |
| Code | 8,192 tokens | 16,384 tokens | 2x more code |
| Sprites | 256 (8×8) | 256 (4-32px) | Variable sizes |
| Physics | None | Box2D | Full simulation |
| Node System | No | Yes | Modern architecture |
| Platforms | Win/Mac/Linux/Web | +Android | Mobile support |
| Price | $15 | Free | Open source |

**RetroForge advantages:**
- More capable (resolution, physics, nodes)
- Free and open source
- Modern architecture
- Mobile support (Android)

**PICO-8 advantages:**
- More established community
- Simpler (fewer features = easier to learn)
- Integrated IDE
- Large cart library

---

## 🎮 Target Audience

### Primary Users
- **Indie developers** - Want to finish games quickly
- **Game jam participants** - Need rapid prototyping
- **Retro enthusiasts** - Love pixel art and chip-tunes
- **Beginners** - Learning game development
- **Educators** - Teaching programming

### Use Cases
- Weekend game projects
- Game jam entries
- Prototyping ideas
- Learning game development
- Retro game nostalgia
- Mobile game development
- Portfolio projects

---

## 💡 Design Philosophy

**1. Creative Constraints**
Limitations spark creativity and help finish projects

**2. Modern + Retro**
Modern tools (nodes, physics) with retro aesthetic

**3. Accessible**
Both high-level (nodes) and low-level (direct) APIs

**4. Open**
Free, open source, learn from others' code

**5. Cross-Platform**
Build once, run everywhere (except iOS)

---

## 🔗 Key Links

**Development:**
- Go: https://go.dev
- Box2D: https://box2d.org
- gomobile: https://pkg.go.dev/golang.org/x/mobile
- gopher-lua: https://github.com/yuin/gopher-lua
- Next.js: https://nextjs.org

**Inspiration:**
- PICO-8: https://www.lexaloffle.com/pico-8.php
- Godot: https://godotengine.org
- TIC-80: https://tic80.com

---

## 📞 Contact & Community (Future)

**Domains:** retroforge.dev, retroforge.com  
**GitHub:** github.com/retroforge  
**Discord:** RetroForge Community  
**Twitter/X:** @retroforge  

---

## 🎉 Final Notes

**RetroForge v2.1 is COMPLETE and READY for implementation!**

All specifications finalized:
✅ Resolution optimized for modern devices  
✅ Platform support expanded (Android)  
✅ Node system designed  
✅ Physics integrated  
✅ Audio architected  
✅ Documentation comprehensive  
✅ Timeline realistic  

**Time to start building the future of retro game development!** 🔨✨

---

**Version History:**
- v1.0 - Initial design
- v2.0 - Node system + physics + audio
- v2.1 - Resolution update + Android support (FINAL)

**This is the definitive specification. Let's build it!**
