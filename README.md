# RetroForge Engine

**RetroForge Fantasy Console - Go Engine**

A modern fantasy console engine built in Go, designed for creating retro-style games with modern development tools. RetroForge provides a complete runtime for 2D games with built-in multiplayer support, physics, audio, and cross-platform deployment.

## 🎯 Project Overview

RetroForge Engine is the core runtime engine that powers the RetroForge fantasy console. It provides:

- **Cross-platform support**: Windows, macOS, Linux, Android, Web (WASM)
- **Lua scripting**: gopher-lua integration for game logic
- **Node system**: Godot-style scene graph architecture
- **Physics engine**: Box2D integration for realistic physics
- **Audio system**: 3-tier audio with chip-tune synthesis
- **Graphics**: 2D rendering with 480×270 resolution, 50-color palette
- **Multiplayer**: WebRTC-based networking for up to 6 players
- **Voxel raytracing**: 2.5D raytracing for retro 3D games (Phase 2)

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────┐
│            Go Runtime Engine                     │
│                                                  │
│  ┌─────────────────┐  ┌────────────────────┐   │
│  │  gopher-lua VM  │  │   Cart Loader      │   │
│  └─────────────────┘  └────────────────────┘   │
│                                                  │
│  ┌──────────────────────────────────────────┐  │
│  │         Scene Graph / Node System        │  │
│  │  - Node hierarchy                        │  │
│  │  - Automatic update/draw                 │  │
│  │  - Transform propagation                 │  │
│  └──────────────────────────────────────────┘  │
│                                                  │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐     │
│  │ Graphics │  │  Audio   │  │ Physics  │     │
│  │  System  │  │  System  │  │ (Box2D)  │     │
│  └──────────┘  └──────────┘  └──────────┘     │
│                                                  │
│  ┌──────────┐  ┌──────────────────────────┐    │
│  │ Network │  │     Platform Layer        │    │
│  │ (WebRTC)│  │    (SDL2/OpenGL/WASM)      │    │
│  └──────────┘  └──────────────────────────┘    │
└─────────────────────────────────────────────────┘
```

## 🚀 Quick Start

### Prerequisites

- Go 1.23+
- SDL2 development libraries (for desktop builds)
- CGO enabled

### Installation

```bash
git clone https://github.com/retroforge/retroforge-engine.git
cd retroforge-engine
go mod download
go build -o retroforge ./cmd/retroforge
```

### Running a Cart

**Desktop:**
```bash
./retroforge path/to/game.rf
```

**Development Mode (Hot Reload):**
```bash
make run-dev FOLDER=examples/multiplayer-platformer
```

**WASM Build (for Web):**
```bash
make wasm
```

## 📦 Cart Format

RetroForge carts are ZIP archives (`.rf` files) containing:

- `manifest.json` - Cart metadata (title, author, description, multiplayer settings)
- `assets/main.lua` - Main game script (Lua code)
- `assets/sfx.json` - Sound effects definitions
- `assets/music.json` - Music track definitions
- `assets/sprites.json` - Sprite definitions
- Additional assets as needed

## 🎮 Example Games

The engine includes several example games:

- **Hello World** - Minimal example with centered text
- **Moon Lander** - Lunar landing game with physics, HUD, and music
- **Tron Light Cycles** - Classic Tron-style game with increasing difficulty
- **Multiplayer Platformer** - Demo showcasing multiplayer features with up to 6 players

Run examples with:
```bash
make run-dev FOLDER=examples/moon-lander
```

## 🔌 Core Features

### Graphics
- **Resolution**: 480×270 (higher than PICO-8's 128×128)
- **Palette**: 50 colors (compared to PICO-8's 16)
- **Drawing primitives**: Lines, rectangles, circles, ellipses, triangles, polygons
- **Sprites**: JSON-based sprite system with named sprites
- **Tilemap**: 256×256 tilemap support
- **Camera**: Viewport/camera system
- **Clipping**: Rectangular clipping regions

### Audio
- **8 audio channels** (compared to PICO-8's 4)
- **5 waveforms**: Sine, square, triangle, sawtooth, noise
- **JSON-based definitions**: Easy to create and modify sounds
- **Music support**: Pattern-based music system

### Physics
- **Box2D integration**: Full rigid body physics
- **Body types**: Static, dynamic, kinematic
- **Shapes**: Boxes, circles, and more
- **Forces and impulses**: Realistic physics interactions

### Input
- **8 standard buttons**: Left, Right, Up, Down, Z, X, plus extended buttons
- **Edge detection**: `btnp()` for just-pressed detection
- **Multiplayer input**: Host can check other players' inputs

### Memory & Persistence
- **2MB runtime memory**: Access via `poke`/`peek` functions
- **64KB cart storage**: Persistent storage (2x PICO-8's 32KB)
- **Cart persistence**: Save/load game state with `cstore()`/`reload()`

### State Machine (NEW!)
- **Flexible state management**: Register states with lifecycle callbacks
- **State stacking**: Push/pop states for overlays (pause menus, inventory)
- **Shared context**: Pass data between states without tight coupling
- **Built-in states**: Engine splash (debug-skippable) and credits screens
- **Complete lifecycle**: Initialize, Enter, HandleInput, Update, Draw, Exit, Shutdown

### Multiplayer (NEW!)
- **Up to 6 players** via WebRTC networking
- **Automatic synchronization**: Register tables for sync with 3-tier system
- **Host authority**: One player controls game logic, prevents conflicts
- **Star topology**: All players connect directly to host
- **Sync tiers**:
  - **Fast**: 30-60 updates/second (player positions)
  - **Moderate**: 15 updates/second (powerups, items)
  - **Slow**: 5 updates/second (scores, UI)
- **Lua API**: Simple API for checking multiplayer state and managing sync

### Development Tools
- **Hot reload**: Auto-reload when files change (development mode)
- **Debug tools**: `printh()`, `stat()`, `time()` functions
- **File watching**: Automatic reload of assets and manifest

## 📚 Documentation

- **[API Reference](API_REFERENCE.md)** - Complete API documentation
- **[PICO-8 Comparison](PICO8_COMPARISON.md)** - Feature-by-feature comparison
- **[Multiplayer Design](RetroForge.V2.md)** - Complete multiplayer architecture
- **[Project Constitution](memory/constitution.md)** - Development principles and standards

## 🛠️ Development

This project uses spec-driven development with the following structure:

```
├── specs/                    # Feature specifications
├── memory/                   # Project constitution
├── scripts/                  # Development scripts
├── templates/                # Code templates
├── examples/                 # Example games
├── cmd/                      # Command-line tools
│   ├── retroforge/          # Main engine binary
│   ├── wasm/                # WASM build target
│   └── cartbundle/          # Cart bundling tool
└── internal/                 # Internal packages
```

### Available Commands

- `make build` - Build desktop binary
- `make wasm` - Build WASM for web
- `make run-dev FOLDER=path` - Run with hot reload
- `make test` - Run tests
- `make clean` - Clean build artifacts

## 🎮 Example Cart

```lua
function _INIT()
  rf.palette_set("SNES 50")
  player = {
    x = 240,
    y = 135,
    vx = 0,
    vy = 0
  }
end

function _UPDATE(dt)
  -- Input
  if rf.btn(0) then player.vx = -3 end  -- Left
  if rf.btn(1) then player.vx = 3 end   -- Right
  if rf.btn(4) and player.on_ground then  -- Jump
    player.vy = -10
    player.on_ground = false
  end
  
  -- Physics
  player.vy = player.vy + 0.5  -- Gravity
  player.x = player.x + player.vx
  player.y = player.y + player.vy
  player.vx = player.vx * 0.9   -- Friction
  
  -- Ground collision
  if player.y > 250 then
    player.y = 250
    player.vy = 0
    player.on_ground = true
  end
end

function _DRAW()
  rf.clear_i(0)  -- Clear to black
  rf.rectfill(0, 250, 480, 270, 38)  -- Ground
  rf.circfill(player.x, player.y, 8, 2)  -- Player
end
```

## 🌐 Multiplayer Example

```lua
function _INIT()
  players = {}
  score = {}
  
  -- Check if multiplayer
  local player_count = rf.is_multiplayer() and rf.player_count() or 1
  
  for i = 1, player_count do
    players[i] = { x = 50 + (i-1) * 100, y = 135, alive = true }
    score[i] = 0
  end
  
  -- Register for sync (multiplayer only)
  if rf.is_multiplayer() then
    rf.network_sync(players, "fast")   -- Smooth movement
    rf.network_sync(score, "slow")     -- Less frequent updates
  end
end

function _UPDATE(dt)
  if rf.is_multiplayer() and rf.is_host() then
    -- Host runs game logic for all players
    for id = 1, rf.player_count() do
      local p = players[id]
      if rf.btn(id, 0) then p.x = p.x - 3 end
      if rf.btn(id, 1) then p.x = p.x + 3 end
      -- Engine automatically syncs players table!
    end
  elseif rf.is_multiplayer() then
    -- Non-host: just send inputs (handled automatically)
    -- Players table is automatically updated by engine
  else
    -- Solo mode: normal game logic
  end
end

function _DRAW()
  rf.clear_i(0)
  for id, p in pairs(players) do
    if p.alive then
      rf.circfill(p.x, p.y, 8, id + 1)
    end
  end
end
```

## 📊 Feature Comparison with PICO-8

| Feature | PICO-8 | RetroForge |
|---------|--------|------------|
| **Resolution** | 128×128 | 480×270 ✅ |
| **Palette** | 16 colors | 50 colors ✅ |
| **Cart Size** | 32 KB | 64 KB ✅ |
| **Audio Channels** | 4 | 8 ✅ |
| **Physics** | Manual | Box2D ✅ |
| **Multiplayer** | Built-in | WebRTC (up to 6) ✅ |
| **Platforms** | Desktop + Web | Desktop + Web + Android ✅ |

See [PICO8_COMPARISON.md](PICO8_COMPARISON.md) for detailed comparison.

## 🤝 Contributing

See [memory/constitution.md](memory/constitution.md) for development principles and standards.

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

---

*"Forge Your Retro Dreams" - RetroForge Engine* 🔨✨
