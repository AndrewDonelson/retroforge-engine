# RetroForge - Project Kickoff & Implementation Plan

**Date:** October 29, 2025  
**Status:** Ready to Begin Implementation  
**Development Mode:** Claude builds with user guidance and review

---

## ✅ Completed: Design Phase

### What We've Accomplished
- ✅ Full technical specification
- ✅ Project name finalized: **RetroForge**
- ✅ File extensions defined: `.rfs` and `.rfe`
- ✅ Technology stack chosen (Go 1.23, Next.js 16, Box2D-go)
- ✅ Node system architecture designed (Godot-style)
- ✅ Physics integration planned (Box2D)
- ✅ Audio system architected (AudioPlayer, MusicPlayer, SoundManager)
- ✅ Complete API specifications documented
- ✅ Color palette system designed
- ✅ Map system detailed
- ✅ Cart format finalized
- ✅ Branding guidelines created

---

## 👨‍💻 Development Approach

**Claude will implement the engine** with user providing:
- Design feedback and direction
- Feature prioritization
- Testing and validation
- Course corrections as needed

**Realistic Timeline:**
- **Phase 1:** Engine Core (2-3 weeks)
- **Phase 2:** Node System & Physics (2-3 weeks)
- **Phase 3:** Web Application (3-4 weeks)
- **Phase 4:** Polish & Features (2-3 weeks)
- **Phase 5:** Android & Cross-Platform (1-2 weeks)
- **Total:** ~3 months of development

**Platform Support:**
- ✅ Desktop: Windows, macOS, Linux
- ✅ Mobile: Android
- ✅ Web: Browser (WASM)
- ❌ iOS: Not supported (cost and restrictions)

Development happens between sessions, with regular check-ins for review and guidance.

---

## 🎯 Immediate Next Steps (This Week)

### 1. Domain & Social Media Setup
**Priority:** HIGH  
**Time:** 2-4 hours

**Tasks:**
- [ ] Check and purchase domains
  - Primary: `retroforge.dev`
  - Secondary: `retroforge.com`
  - Docs: `docs.retroforge.dev`
- [ ] Create GitHub organization: `github.com/retroforge`
- [ ] Reserve social media handles
  - Twitter/X: @retroforge or @retroforge_dev
  - Discord: Create server
  - Reddit: r/retroforge (if available)

### 2. Repository Setup
**Priority:** HIGH  
**Time:** 2-3 hours

**Create repositories:**
```bash
retroforge-engine/          # Go engine
├── README.md
├── go.mod
├── LICENSE (MIT)
└── .github/
    └── workflows/

retroforge-web/             # Next.js web app
├── README.md
├── package.json
├── LICENSE (MIT)
└── .github/
    └── workflows/

retroforge-docs/            # Documentation
├── README.md
└── docs/

retroforge-examples/        # Example carts
└── carts/
    └── hello-world.rfs
```

**Repository setup checklist:**
- [ ] Create GitHub organization
- [ ] Create repositories with MIT license
- [ ] Add comprehensive README files
- [ ] Set up GitHub Actions for CI/CD
- [ ] Configure branch protection rules
- [ ] Add issue templates
- [ ] Create contributing guidelines

### 3. Define Color Palettes
**Priority:** MEDIUM  
**Time:** 4-6 hours

We have the palette structure (Black + White + 16 colors × 3 shades), now we need actual hex values.

**Required palettes (16+):**
1. Default (balanced, general purpose)
2. Sunset (complementary)
3. Forest (complementary)
4. Ocean (analogous)
5. Autumn (analogous)
6. Primary (triadic)
7. Vibrant (triadic)
8. Spring (seasonal)
9. Summer (seasonal)
10. Fall (seasonal)
11. Winter (seasonal)
12. Energetic (mood)
13. Calm (mood)
14. Mysterious (mood)
15. Cheerful (mood)
16. Neon (specialty)
17. Gameboy (specialty)
18. Grayscale (specialty)

**Action:** Create `palettes.json` with all hex values

---

## 📋 Phase 1: Engine Foundation (Weeks 1-4)

### Week 1: Project Setup & Base Node System

**Repository Structure:**
```
retroforge-engine/
├── cmd/
│   └── retroforge/
│       └── main.go
├── internal/
│   ├── node/          # Node system
│   ├── scene/         # Scene graph
│   ├── physics/       # Box2D integration
│   ├── graphics/      # Rendering
│   ├── audio/         # Audio system
│   ├── input/         # Input handling
│   ├── lua/           # Lua bindings
│   └── cart/          # Cart loading
├── pkg/
│   └── math/          # Vector math
├── go.mod
└── README.md
```

**Core Tasks:**
- [ ] Initialize Go module
- [ ] Set up SDL2 window and OpenGL context
- [ ] Create basic rendering loop (60 FPS)
- [ ] Implement base Node class
- [ ] Implement Node2D class
- [ ] Create Scene root node
- [ ] Basic scene graph traversal (update/draw)

**Deliverable:** Window opens, scene graph works, can add/remove nodes

### Week 2: Physics Integration (Box2D)

**Box2D Setup:**
- [ ] Add `github.com/ByteArena/box2d` dependency
- [ ] Create PhysicsWorld wrapper
- [ ] Implement PhysicsBody2D base class
- [ ] Implement StaticBody, RigidBody, KinematicBody
- [ ] Implement CollisionShape (rect, circle, capsule, polygon)
- [ ] Sync Box2D transforms to node transforms
- [ ] Collision callbacks (on_body_entered, on_body_exited)
- [ ] Physics debug rendering

**Deliverable:** Physics bodies fall with gravity, collide correctly

### Week 3: Lua Integration & Node API

**gopher-lua Setup:**
- [ ] Add gopher-lua dependency
- [ ] Create Lua VM wrapper
- [ ] Implement lifecycle callbacks (_init, _ready, _update, _draw)
- [ ] Bind Node classes to Lua
- [ ] Bind PhysicsBody classes to Lua
- [ ] Implement vec2 type and operations
- [ ] Implement Scene global table
- [ ] Implement Input global table

**API Bindings:**
```lua
-- Should work by end of week 3
player = RigidBody.new({
  position = vec2(240, 135),
  width = 16,
  height = 16
})
Scene:add_child(player)
```

**Deliverable:** Can create and control nodes from Lua

### Week 4: Graphics & First Complete Example

**Graphics System:**
- [ ] Color palette loading
- [ ] Sprite loading from PNG
- [ ] Implement Sprite node
- [ ] Sprite rendering with transforms
- [ ] Camera node implementation
- [ ] Camera following target
- [ ] Layer/z-index sorting
- [ ] Direct API graphics functions (cls, pset, spr, rect, circ, etc.)

**Example Cart: Physics Demo**
```lua
function _init()
  -- Create ground
  ground = StaticBody.new({
    position = vec2(160, 220),
    width = 320,
    height = 20
  })
  Scene:add_child(ground)
  
  -- Create falling boxes
  for i = 1, 10 do
    local box = RigidBody.new({
      position = vec2(100 + i * 10, 50),
      width = 16,
      height = 16
    })
    
    local sprite = Sprite.new({sprite_index = 1})
    box:add_child(sprite)
    
    Scene:add_child(box)
  end
  
  -- Create player
  player = RigidBody.new({
    position = vec2(160, 100),
    width = 16,
    height = 16,
    fixed_rotation = true
  })
  
  local sprite = Sprite.new({sprite_index = 2})
  player:add_child(sprite)
  Scene:add_child(player)
  
  -- Camera follows player
  camera = Camera.new()
  camera:follow(player)
  Scene:add_child(camera)
end

function player:_update()
  if Input.is_action_pressed("left") then
    self:apply_force(vec2(-500, 0))
  end
  if Input.is_action_pressed("right") then
    self:apply_force(vec2(500, 0))
  end
  if Input.is_action_just_pressed("jump") then
    self:apply_impulse(vec2(0, -300))
  end
end
```

**Deliverable:** Complete physics demo cart runs with node system

---

## 📋 Phase 2: Audio & Advanced Nodes (Weeks 5-7)

### Week 5: Audio System Foundation

**SoundManager Implementation:**
- [ ] Create SoundManager singleton
- [ ] Audio channel management (8 channels)
- [ ] Volume controls (master, sfx, music)
- [ ] Audio ducking system
- [ ] Simple audio synthesis (square, triangle, sine, sawtooth, noise)
- [ ] ADSR envelope implementation
- [ ] Sound effect playback

**Deliverable:** Can play simple chip-tune sounds

### Week 6: AudioPlayer & MusicPlayer Nodes

**AudioPlayer Node:**
- [ ] Multiple simultaneous SFX playback
- [ ] Per-instance volume/pitch control
- [ ] Loop support
- [ ] Bus routing

**MusicPlayer Node:**
- [ ] Single music track playback
- [ ] Cross-fade implementation
- [ ] Pattern sequencer
- [ ] Tempo control
- [ ] Seamless looping

**Deliverable:** Full audio system works with SFX and music

### Week 7: Advanced Nodes

**Implement:**
- [ ] AnimatedSprite (sprite animations)
- [ ] TileMap (map rendering with layers)
- [ ] ParticleEmitter (particle effects)
- [ ] Timer (delayed actions)
- [ ] CollisionShape refinements
- [ ] Node groups
- [ ] Raycasting API

**Deliverable:** All core nodes implemented and tested

---

## 📋 Phase 3: Web Application (Weeks 8-10)

### Week 5: Next.js Setup

**Project Initialization**
```bash
npx create-next-app@latest retroforge-web \
  --typescript \
  --tailwind \
  --app \
  --src-dir
```

**Initial Structure**
```
retroforge-web/
├── src/
│   ├── app/
│   │   ├── layout.tsx
│   │   ├── page.tsx (landing)
│   │   ├── create/
│   │   │   └── page.tsx (editor)
│   │   └── play/
│   │       └── [id]/
│   │           └── page.tsx
│   ├── components/
│   │   ├── Editor/
│   │   │   ├── CodeEditor.tsx
│   │   │   ├── SpriteEditor.tsx
│   │   │   └── MapEditor.tsx
│   │   └── UI/
│   ├── lib/
│   │   ├── cart.ts
│   │   └── wasm.ts
│   └── types/
│       └── cart.ts
└── public/
    └── engine.wasm
```

**Tasks:**
- [ ] Initialize Next.js project
- [ ] Set up TailwindCSS
- [ ] Install dependencies (react-ace, etc.)
- [ ] Create basic layout
- [ ] Set up routing

### Week 6: Code Editor

**react-ace Integration**
- [ ] Install react-ace and ace-builds
- [ ] Create CodeEditor component
- [ ] Configure Lua syntax highlighting
- [ ] Add RetroForge API autocomplete
- [ ] Implement save/load functionality
- [ ] Add line numbers, bracket matching
- [ ] Theme support (light/dark)

**Features:**
- [ ] Real-time syntax validation
- [ ] API documentation tooltips
- [ ] Keyboard shortcuts
- [ ] Multiple file tabs (for code/ directory)

### Week 7: Sprite Editor

**Canvas-based Editor**
- [ ] Drawing tools (pencil, fill, line, rect, circle)
- [ ] Color palette selector
- [ ] Sprite grid view
- [ ] Individual sprite editor
- [ ] Import/export PNG
- [ ] Undo/redo system

**Features:**
- [ ] Multi-size sprite support (4×4, 8×8, 16×16, 32×32)
- [ ] Copy/paste sprites
- [ ] Flip horizontal/vertical
- [ ] Preview at actual size

### Week 8: WASM Integration

**Compile Engine to WASM**
```bash
# Using TinyGo
tinygo build -o public/engine.wasm -target wasm ./engine
```

**JavaScript Bridge**
- [ ] Load WASM module
- [ ] Create JavaScript ↔ WASM interface
- [ ] Pass cart data to WASM
- [ ] Handle input from browser
- [ ] Render to HTML Canvas
- [ ] Audio output via Web Audio API

**Deliverable:** Cart runs in browser

---

## 📋 Phase 4: Complete Feature Set & Polish (Weeks 11-12)

### Week 9: Map Editor
- [ ] 8-layer map editor
- [ ] Tile placement
- [ ] Parallax configuration
- [ ] Tile flag editor
- [ ] Multiple level support

### Week 10: Audio Editor
- [ ] Piano roll interface
- [ ] Waveform selection
- [ ] ADSR envelope editor
- [ ] Pattern sequencer
- [ ] Real-time preview

### Week 11: Polish & Testing
- [ ] Complete Lua API implementation
- [ ] Add all missing graphics functions
- [ ] Map system with parallax
- [ ] Persistent storage (dget/dset)
- [ ] Performance optimization

### Week 12: Example Carts & Documentation
- [ ] Create 5 example carts
- [ ] Write API documentation
- [ ] Create tutorial series
- [ ] User guide
- [ ] Video walkthroughs

---

## 🛠️ Development Environment Setup

### Required Tools

**Go Development:**
```bash
# Install Go 1.23
# macOS
brew install go@1.23

# Linux
wget https://go.dev/dl/go1.23.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.linux-amd64.tar.gz

# Install TinyGo (for WASM)
brew install tinygo  # macOS
# or follow: https://tinygo.org/getting-started/install/
```

**SDL2:**
```bash
# macOS
brew install sdl2 sdl2_image sdl2_ttf

# Linux (Ubuntu/Debian)
sudo apt-get install libsdl2-dev libsdl2-image-dev libsdl2-ttf-dev

# Windows
# Download from: https://www.libsdl.org/download-2.0.php
```

**Node.js & NPM:**
```bash
# Install Node.js 20+ LTS
# macOS
brew install node

# Or use nvm
nvm install 20
nvm use 20
```

**Recommended IDE:**
- **VS Code** with extensions:
  - Go (official)
  - TailwindCSS IntelliSense
  - Prettier
  - ESLint
  - TypeScript

### Clone and Setup

```bash
# Clone repositories
git clone https://github.com/retroforge/engine.git retroforge-engine
git clone https://github.com/retroforge/web.git retroforge-web

# Engine setup
cd retroforge-engine
go mod download
go build

# Web setup
cd ../retroforge-web
npm install
npm run dev
```

---

## 📊 Success Metrics

### MVP Success Criteria
- [ ] Engine runs at stable 60 FPS
- [ ] Can create and run a simple cart
- [ ] Web editor loads and works
- [ ] WASM version runs in browser
- [ ] All core graphics API implemented
- [ ] Input handling works
- [ ] At least 1 complete example cart

### Beta Launch Criteria
- [ ] All 8 API categories complete
- [ ] Sprite editor fully functional
- [ ] Map editor with 8 layers
- [ ] Audio editor with piano roll
- [ ] 5+ example carts
- [ ] Documentation complete
- [ ] 10 alpha testers give positive feedback

### Public Launch Criteria
- [ ] All features complete and polished
- [ ] Performance optimized
- [ ] Comprehensive documentation
- [ ] Video tutorials
- [ ] Landing page live
- [ ] Cart repository functional
- [ ] Community Discord active

---

## 📋 Phase 5: Android & Cross-Platform (Weeks 13-14)

### Week 13: Android Build Setup

**gomobile Setup:**
- [ ] Install gomobile (`go install golang.org/x/mobile/cmd/gomobile@latest`)
- [ ] Initialize gomobile (`gomobile init`)
- [ ] Create Android-specific rendering layer
- [ ] Touch input handling for Android
- [ ] OpenGL ES context setup
- [ ] OpenSL ES audio backend

**Android Build Configuration:**
```bash
# Build APK
gomobile bind -target=android \
  -o retroforge.aar \
  github.com/retroforge/engine
```

**Deliverable:** Engine builds as Android APK

### Week 14: Android Testing & Optimization

**Testing:**
- [ ] Test on various Android devices
- [ ] Performance optimization for mobile
- [ ] Battery usage optimization
- [ ] Touch input refinement
- [ ] Orientation handling (landscape/portrait)
- [ ] Android permissions setup

**Distribution:**
- [ ] Generate signed APK
- [ ] Test sideloading
- [ ] Prepare for Play Store (optional)
- [ ] Create installation instructions

**Screen Scaling for Mobile:**
- Base resolution: 480×270 (landscape) or 270×480 (portrait)
- Auto-scale to device resolution
- Maintain aspect ratio
- Touch input mapping

**Deliverable:** RetroForge runs on Android devices

---

## 📋 Platform Build Matrix

| Platform | Build Tool | Output | Distribution |
|----------|-----------|---------|--------------|
| Windows | Go 1.23 | retroforge.exe | Direct download |
| macOS | Go 1.23 | RetroForge.app | Direct download / DMG |
| Linux | Go 1.23 | retroforge | Direct download / AppImage |
| Android | gomobile | retroforge.apk | Sideload / Play Store |
| Web | TinyGo | retroforge.wasm | Hosted on site |
| iOS | ❌ | - | Not supported |

---

## 🎮 First Example Cart Ideas

1. **Hello World** - Moving circle
2. **Pong** - Classic paddle game
3. **Space Shooter** - Simple shmup
4. **Platformer** - Basic jumping game
5. **Breakout** - Brick breaking
6. **Snake** - Classic snake game
7. **Maze** - Procedural maze walker
8. **Music Demo** - Show off audio system

---

## 🤝 Getting Help

### Documentation
- Go: https://go.dev/doc/
- gopher-lua: https://github.com/yuin/gopher-lua
- SDL2: https://wiki.libsdl.org/
- Next.js: https://nextjs.org/docs
- TinyGo: https://tinygo.org/docs/

### Community
- Create issues on GitHub for bugs
- Use Discussions for questions
- Join Discord for real-time help
- Share progress on Twitter/X

---

## 🚀 Let's Build!

RetroForge is ready to move from design to development. With clear specifications, a solid tech stack, and a phased approach, we're set up for success.

**First commit:** Initialize repositories and begin Week 1 tasks!

**Next document needed:** Technical implementation guide for the Go engine.

---

*Ready to forge some retro magic? Let's do this! 🔨✨*
