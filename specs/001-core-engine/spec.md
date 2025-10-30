# RetroForge Engine - Core Engine Specification

**Spec ID:** 001-core-engine  
**Version:** 1.0  
**Date:** October 30, 2025  
**Status:** Ready for Implementation

---

## 🎯 Overview

This specification defines the core RetroForge Engine implementation, including the Go runtime, Lua integration, node system, physics engine, audio system, and graphics rendering.

## 📋 User Stories

### As a Game Developer
- I want to create retro games using Lua scripting so that I can focus on game logic
- I want to use a node system so that I can organize my game objects hierarchically
- I want physics simulation so that I can create realistic interactions
- I want audio capabilities so that I can add sound effects and music
- I want cross-platform deployment so that my games work everywhere

### As a Retro Game Enthusiast
- I want to play RetroForge games on any platform so that I can enjoy them anywhere
- I want fast loading times so that I can start playing immediately
- I want smooth 60 FPS performance so that games feel responsive
- I want authentic retro aesthetics so that games feel like classic consoles

## 🏗️ Technical Requirements

### Core Engine
- **Language**: Go 1.23+
- **Runtime**: gopher-lua (Lua 5.1)
- **Graphics**: SDL2 + OpenGL
- **Physics**: Box2D-go
- **Audio**: Custom synthesis engine
- **Platforms**: Windows, macOS, Linux, Android, Web (WASM)

### Performance Targets
- **Frame Rate**: 60 FPS stable
- **Frame Budget**: <16.67ms per frame
- **Memory Usage**: <64MB WASM heap
- **Load Time**: <2 seconds for typical cart
- **Audio Latency**: <20ms

### API Requirements
- **Dual API**: High-level nodes + low-level direct
- **Lua Integration**: Full API exposed to Lua
- **Documentation**: Comprehensive with examples
- **Error Handling**: Clear error messages
- **Debugging**: Built-in debugging tools

## 🎮 Core Features

### 1. Node System
```go
type Node interface {
    AddChild(child Node)
    RemoveChild(child Node)
    GetParent() Node
    GetChildren() []Node
    Update(deltaTime float64)
    Draw(renderer *Renderer)
}
```

**Node Types**:
- `Node` - Base node class
- `Node2D` - 2D spatial node
- `Sprite` - Sprite rendering
- `AnimatedSprite` - Animated sprites
- `TileMap` - Tilemap rendering
- `Camera` - Viewport control
- `PhysicsBody2D` - Physics integration
- `AudioPlayer` - Sound effects
- `MusicPlayer` - Music tracks

### 2. Physics Engine
```go
type PhysicsWorld struct {
    world *box2d.B2World
    bodies map[int]*PhysicsBody2D
}

type PhysicsBody2D struct {
    body *box2d.B2Body
    bodyType BodyType
    collisionLayer int
    collisionMask int
}
```

**Body Types**:
- `StaticBody` - Immovable objects (walls, platforms)
- `RigidBody` - Physics-simulated objects
- `KinematicBody` - Player-controlled objects

### 3. Audio System
```go
type AudioSystem struct {
    soundManager *SoundManager
    audioPlayers []*AudioPlayer
    musicPlayer *MusicPlayer
}

type SoundManager struct {
    masterVolume float64
    sfxVolume float64
    musicVolume float64
    duckingEnabled bool
}
```

**Audio Features**:
- 8 simultaneous audio channels
- Chip-tune synthesis (5 waveforms)
- ADSR envelope controls
- Audio ducking
- Spatial audio support

### 4. Graphics System
```go
type Renderer struct {
    window *sdl.Window
    context sdl.GLContext
    shaderProgram uint32
    projectionMatrix mat4
}

type Sprite struct {
    texture uint32
    width int
    height int
    uvCoords [4]vec2
}
```

**Graphics Features**:
- 480×270 resolution (16:9)
- 50-color palette system
- Sprite rendering (4×4 to 32×32)
- Tilemap rendering (8 layers)
- Parallax scrolling
- Camera system

### 5. Lua Integration
```go
type LuaVM struct {
    L *lua.LState
    api *LuaAPI
}

type LuaAPI struct {
    graphics *GraphicsAPI
    input *InputAPI
    audio *AudioAPI
    physics *PhysicsAPI
    math *MathAPI
}
```

**Lua API**:
- Complete RetroForge API exposed
- Node system integration
- Event handling
- Memory management
- Error reporting

## 🔧 Implementation Details

### Project Structure
```
retroforge-engine/
├── cmd/
│   └── retroforge/
│       └── main.go
├── internal/
│   ├── engine/
│   │   ├── engine.go
│   │   └── game_loop.go
│   ├── graphics/
│   │   ├── renderer.go
│   │   ├── sprite.go
│   │   └── tilemap.go
│   ├── audio/
│   │   ├── audio_system.go
│   │   ├── sound_manager.go
│   │   └── synthesis.go
│   ├── physics/
│   │   ├── physics_world.go
│   │   └── physics_body.go
│   ├── nodes/
│   │   ├── node.go
│   │   ├── node2d.go
│   │   └── sprite.go
│   └── lua/
│       ├── lua_vm.go
│       ├── api.go
│       └── bindings.go
├── pkg/
│   ├── cart/
│   │   ├── loader.go
│   │   └── format.go
│   └── platform/
│       ├── sdl2.go
│       └── wasm.go
└── examples/
    ├── hello_world/
    └── platformer/
```

### Dependencies
```go
// go.mod
module github.com/retroforge/retroforge-engine

go 1.23

require (
    github.com/yuin/gopher-lua v0.0.0-20220504180219-658193537a64
    github.com/ByteArena/box2d v1.0.2
    github.com/veandco/go-sdl2 v0.4.25
    github.com/go-gl/gl v0.0.0-20211210172815-726fda9656d6
    github.com/go-gl/glfw/v3.3/glfw v0.0.0-20221017161538-93cebf72946b
)
```

### Build Targets
```bash
# Desktop builds
go build -o retroforge ./cmd/retroforge

# WASM build
GOOS=js GOARCH=wasm go build -o retroforge.wasm ./cmd/retroforge

# Android build (using gomobile)
gomobile bind -target=android ./pkg/mobile
```

## 🧪 Testing Strategy

### Unit Tests
- **Node system** - Hierarchy and lifecycle
- **Physics** - Collision detection and response
- **Audio** - Synthesis and playback
- **Graphics** - Rendering and transformations
- **Lua API** - All exposed functions

### Integration Tests
- **Cart loading** - Various cart formats
- **Cross-platform** - All target platforms
- **Performance** - Frame rate and memory usage
- **Audio** - Latency and quality
- **Physics** - Complex scenarios

### Example Carts
- **Hello World** - Basic functionality
- **Platformer** - Physics and input
- **Audio Demo** - Sound and music
- **Sprite Animation** - Graphics features
- **Voxel Demo** - 3D raytracing (Phase 2)

## 📚 Documentation Requirements

### API Documentation
- **Complete reference** for all functions
- **Code examples** for common patterns
- **Performance notes** for optimization
- **Platform differences** where applicable
- **Migration guides** for updates

### User Guides
- **Getting started** tutorial
- **Node system** guide
- **Physics** tutorial
- **Audio** composition guide
- **Performance** optimization tips

### Developer Documentation
- **Architecture** overview
- **Contributing** guidelines
- **Build** instructions
- **Testing** procedures
- **Release** process

## 🚀 Success Criteria

### Technical Goals
- ✅ **60 FPS** on all target platforms
- ✅ **<64MB** memory usage for typical carts
- ✅ **<2 seconds** cart loading time
- ✅ **Zero crashes** in stable releases
- ✅ **100% API** test coverage

### User Experience Goals
- ✅ **<5 minutes** to first running cart
- ✅ **Clear documentation** for all features
- ✅ **Helpful error messages** for common issues
- ✅ **Smooth performance** across all platforms
- ✅ **Intuitive APIs** for common tasks

### Quality Goals
- ✅ **Clean code** following Go best practices
- ✅ **Comprehensive testing** at all levels
- ✅ **Performance optimization** for all features
- ✅ **Cross-platform** compatibility
- ✅ **Maintainable** architecture

---

## 📋 Implementation Checklist

### Phase 1: Core Foundation (Week 1-2)
- [ ] Go project setup with dependencies
- [ ] Basic SDL2 window and rendering loop
- [ ] Lua VM integration with gopher-lua
- [ ] Basic node system implementation
- [ ] Simple sprite rendering

### Phase 2: Node System (Week 3-4)
- [ ] Complete node hierarchy implementation
- [ ] Node2D with transform system
- [ ] Sprite and AnimatedSprite nodes
- [ ] Camera system
- [ ] Scene graph management

### Phase 3: Physics Integration (Week 5-6)
- [ ] Box2D integration
- [ ] PhysicsBody2D node types
- [ ] Collision detection and response
- [ ] Physics world management
- [ ] Joint and constraint support

### Phase 4: Audio System (Week 7-8)
- [ ] Audio synthesis engine
- [ ] SoundManager implementation
- [ ] AudioPlayer and MusicPlayer nodes
- [ ] Waveform generation (5 types)
- [ ] ADSR envelope controls

### Phase 5: Graphics System (Week 9-10)
- [ ] Complete sprite rendering
- [ ] Tilemap system with 8 layers
- [ ] Parallax scrolling
- [ ] Color palette system
- [ ] Camera and viewport controls

### Phase 6: Lua API (Week 11-12)
- [ ] Complete Lua API bindings
- [ ] Node system Lua integration
- [ ] Event handling system
- [ ] Memory management
- [ ] Error reporting and debugging

### Phase 7: Testing & Optimization (Week 13-14)
- [ ] Comprehensive test suite
- [ ] Performance optimization
- [ ] Cross-platform testing
- [ ] Documentation completion
- [ ] Example carts creation

---

**This specification provides the complete technical foundation for implementing the RetroForge Engine core functionality. All implementation should follow this specification and the project constitution.**

---

*"Forge Your Retro Dreams" - RetroForge Engine Core* 🔨✨
