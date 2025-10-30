# RetroForge Engine

**RetroForge Fantasy Console - Go Engine**

A modern fantasy console engine built in Go, designed for creating retro-style games with modern development tools.

## 🎯 Project Overview

RetroForge Engine is the core runtime engine that powers the RetroForge fantasy console. It provides:

- **Cross-platform support**: Windows, macOS, Linux, Android, Web (WASM)
- **Lua scripting**: gopher-lua integration for game logic
- **Node system**: Godot-style scene graph architecture
- **Physics engine**: Box2D integration for realistic physics
- **Audio system**: 3-tier audio with chip-tune synthesis
- **Graphics**: 2D rendering with 480×270 resolution
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
│  ┌─────────────────────────────────────────┐   │
│  │     Platform Layer (SDL2/OpenGL)        │   │
│  └─────────────────────────────────────────┘   │
└─────────────────────────────────────────────────┘
```

## 🚀 Quick Start

### Prerequisites

- Go 1.23+
- SDL2 development libraries
- CGO enabled

### Installation

```bash
git clone https://github.com/retroforge/retroforge-engine.git
cd retroforge-engine
go mod download
go build -o retroforge .
```

### Running a Cart

```bash
./retroforge path/to/game.rfe
```

## 📚 Documentation

- [Design Document](../design/RETROFORGE_DESIGN.md) - Complete technical specification
- [Node Architecture](../design/RETROFORGE_NODE_ARCHITECTURE.md) - Scene graph details
- [Voxel Raytracing](../design/RETROFORGE_VOXEL_RAYTRACING.md) - 3D features (Phase 2)

## 🛠️ Development

This project uses [Spec-Driven Development](https://github.com/github/spec-kit) with the following structure:

```
├── specs/                    # Feature specifications
├── memory/                   # Project constitution
├── scripts/                  # Development scripts
├── templates/                # Code templates
└── src/                      # Source code
```

### Available Commands

- `./scripts/setup-plan.sh` - Initialize development plan
- `./scripts/create-new-feature.sh` - Create new feature spec
- `./scripts/check-prerequisites.sh` - Verify development environment

## 🎮 Example Cart

```lua
function _init()
  player = RigidBody.new({
    position = vec2(240, 135),
    width = 16,
    height = 16
  })
  
  local sprite = Sprite.new({sprite_index = 1})
  player:add_child(sprite)
  Scene:add_child(player)
end

function player:_update()
  if Input.is_action_just_pressed("jump") then
    self:apply_impulse(vec2(0, -300))
  end
end
```

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🤝 Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines.

---

*"Forge Your Retro Dreams" - RetroForge Engine* 🔨✨
