#!/bin/bash

# RetroForge Engine - Setup Plan Script
# This script initializes the development plan for the RetroForge Engine

set -e

echo "🔨 Setting up RetroForge Engine development plan..."

# Check if we're in the right directory
if [ ! -f "memory/constitution.md" ]; then
    echo "❌ Error: Please run this script from the retroforge-engine root directory"
    exit 1
fi

# Create initial plan if it doesn't exist
if [ ! -f "specs/001-core-engine/plan.md" ]; then
    echo "📋 Creating initial development plan..."
    cat > "specs/001-core-engine/plan.md" << 'EOF'
# RetroForge Engine - Development Plan

**Plan ID:** 001-core-engine  
**Version:** 1.0  
**Date:** October 30, 2025  
**Status:** Ready for Implementation

---

## 🎯 Implementation Overview

This plan outlines the implementation of the core RetroForge Engine, including the Go runtime, Lua integration, node system, physics engine, audio system, and graphics rendering.

## 📋 Implementation Phases

### Phase 1: Foundation (Weeks 1-2)
**Goal**: Establish core Go project structure and basic functionality

**Tasks**:
- [ ] Initialize Go module with dependencies
- [ ] Set up SDL2 window and rendering loop
- [ ] Integrate gopher-lua for scripting
- [ ] Implement basic node system
- [ ] Create simple sprite rendering

**Deliverables**:
- Working Go application with SDL2
- Basic Lua script execution
- Simple node hierarchy
- Sprite rendering to screen

### Phase 2: Node System (Weeks 3-4)
**Goal**: Complete node system implementation

**Tasks**:
- [ ] Implement complete node hierarchy
- [ ] Add Node2D with transform system
- [ ] Create Sprite and AnimatedSprite nodes
- [ ] Implement Camera system
- [ ] Add scene graph management

**Deliverables**:
- Complete node system
- Transform propagation
- Scene graph management
- Camera controls

### Phase 3: Physics Integration (Weeks 5-6)
**Goal**: Integrate Box2D physics engine

**Tasks**:
- [ ] Integrate Box2D-go library
- [ ] Implement PhysicsBody2D node types
- [ ] Add collision detection and response
- [ ] Create physics world management
- [ ] Add joint and constraint support

**Deliverables**:
- Physics simulation working
- Collision detection
- Physics body types
- Joint system

### Phase 4: Audio System (Weeks 7-8)
**Goal**: Implement audio synthesis and playback

**Tasks**:
- [ ] Create audio synthesis engine
- [ ] Implement SoundManager
- [ ] Add AudioPlayer and MusicPlayer nodes
- [ ] Create waveform generation (5 types)
- [ ] Add ADSR envelope controls

**Deliverables**:
- Audio synthesis working
- Sound effects playback
- Music playback
- Audio ducking

### Phase 5: Graphics System (Weeks 9-10)
**Goal**: Complete graphics rendering system

**Tasks**:
- [ ] Implement complete sprite rendering
- [ ] Add tilemap system with 8 layers
- [ ] Create parallax scrolling
- [ ] Implement color palette system
- [ ] Add camera and viewport controls

**Deliverables**:
- Complete graphics system
- Tilemap rendering
- Parallax scrolling
- Color palette support

### Phase 6: Lua API (Weeks 11-12)
**Goal**: Complete Lua API integration

**Tasks**:
- [ ] Implement complete Lua API bindings
- [ ] Add node system Lua integration
- [ ] Create event handling system
- [ ] Add memory management
- [ ] Implement error reporting and debugging

**Deliverables**:
- Complete Lua API
- Event system
- Memory management
- Debugging tools

### Phase 7: Testing & Optimization (Weeks 13-14)
**Goal**: Comprehensive testing and performance optimization

**Tasks**:
- [ ] Create comprehensive test suite
- [ ] Optimize performance
- [ ] Test cross-platform compatibility
- [ ] Complete documentation
- [ ] Create example carts

**Deliverables**:
- Test suite
- Performance benchmarks
- Cross-platform builds
- Documentation
- Example games

## 🎯 Success Criteria

### Technical Goals
- ✅ 60 FPS on all target platforms
- ✅ <64MB memory usage for typical carts
- ✅ <2 seconds cart loading time
- ✅ Zero crashes in stable releases
- ✅ 100% API test coverage

### User Experience Goals
- ✅ <5 minutes to first running cart
- ✅ Clear documentation for all features
- ✅ Helpful error messages for common issues
- ✅ Smooth performance across all platforms
- ✅ Intuitive APIs for common tasks

## 🚀 Next Steps

1. **Review this plan** with the development team
2. **Set up development environment** with all dependencies
3. **Begin Phase 1** implementation
4. **Regular checkpoints** to ensure progress
5. **Continuous testing** throughout development

---

**This plan provides the roadmap for implementing the RetroForge Engine core functionality. All development should follow this plan and the project constitution.**

---

*"Forge Your Retro Dreams" - RetroForge Engine Plan* 🔨✨
EOF
    echo "✅ Initial development plan created"
else
    echo "ℹ️  Development plan already exists"
fi

# Create tasks breakdown if it doesn't exist
if [ ! -f "specs/001-core-engine/tasks.md" ]; then
    echo "📋 Creating task breakdown..."
    cat > "specs/001-core-engine/tasks.md" << 'EOF'
# RetroForge Engine - Task Breakdown

**Task ID:** 001-core-engine  
**Version:** 1.0  
**Date:** October 30, 2025  
**Status:** Ready for Implementation

---

## 📋 Task Overview

This document breaks down the core engine implementation into specific, actionable tasks that can be executed in the correct order.

## 🎯 User Story 1: Basic Engine Foundation

### Task 1.1: Project Setup
**Priority**: High  
**Estimated Time**: 4 hours  
**Dependencies**: None

**Description**: Initialize Go project with all required dependencies

**Tasks**:
- [ ] Create go.mod file with RetroForge module
- [ ] Add all required dependencies (gopher-lua, box2d-go, go-sdl2, etc.)
- [ ] Set up basic project structure
- [ ] Create main.go entry point
- [ ] Add build scripts for all platforms

**Files to Create**:
- `go.mod`
- `cmd/retroforge/main.go`
- `internal/engine/engine.go`
- `scripts/build.sh`

**Acceptance Criteria**:
- [ ] Project builds successfully
- [ ] All dependencies resolve
- [ ] Basic project structure in place

### Task 1.2: SDL2 Window Setup
**Priority**: High  
**Estimated Time**: 6 hours  
**Dependencies**: Task 1.1

**Description**: Create basic SDL2 window and rendering loop

**Tasks**:
- [ ] Initialize SDL2
- [ ] Create window with 480×270 resolution
- [ ] Set up OpenGL context
- [ ] Create basic rendering loop
- [ ] Handle window events

**Files to Create**:
- `internal/graphics/renderer.go`
- `internal/platform/sdl2.go`
- `internal/engine/game_loop.go`

**Acceptance Criteria**:
- [ ] Window opens and displays
- [ ] Rendering loop runs at 60 FPS
- [ ] Window events handled properly

### Task 1.3: Lua Integration
**Priority**: High  
**Estimated Time**: 8 hours  
**Dependencies**: Task 1.2

**Description**: Integrate gopher-lua for script execution

**Tasks**:
- [ ] Initialize Lua VM
- [ ] Create basic Lua API bindings
- [ ] Implement script loading
- [ ] Add error handling
- [ ] Create test Lua script

**Files to Create**:
- `internal/lua/lua_vm.go`
- `internal/lua/api.go`
- `internal/lua/bindings.go`
- `examples/hello_world.lua`

**Acceptance Criteria**:
- [ ] Lua scripts can be loaded and executed
- [ ] Basic API functions work
- [ ] Error handling works properly

## 🎯 User Story 2: Node System Implementation

### Task 2.1: Base Node Class
**Priority**: High  
**Estimated Time**: 6 hours  
**Dependencies**: Task 1.3

**Description**: Implement base Node class and hierarchy

**Tasks**:
- [ ] Create Node interface
- [ ] Implement base Node struct
- [ ] Add parent/child relationships
- [ ] Implement add/remove child methods
- [ ] Add node lifecycle methods

**Files to Create**:
- `internal/nodes/node.go`
- `internal/nodes/node_interface.go`

**Acceptance Criteria**:
- [ ] Node hierarchy works
- [ ] Parent/child relationships work
- [ ] Lifecycle methods called correctly

### Task 2.2: Node2D Implementation
**Priority**: High  
**Estimated Time**: 8 hours  
**Dependencies**: Task 2.1

**Description**: Implement Node2D with transform system

**Tasks**:
- [ ] Create Node2D struct
- [ ] Implement position, rotation, scale
- [ ] Add transform calculations
- [ ] Implement global/local transforms
- [ ] Add transform propagation

**Files to Create**:
- `internal/nodes/node2d.go`
- `internal/math/transform.go`

**Acceptance Criteria**:
- [ ] Transform calculations work
- [ ] Global transforms calculated correctly
- [ ] Transform propagation works

### Task 2.3: Sprite Node
**Priority**: Medium  
**Estimated Time**: 6 hours  
**Dependencies**: Task 2.2

**Description**: Implement Sprite node for rendering

**Tasks**:
- [ ] Create Sprite struct
- [ ] Implement texture loading
- [ ] Add sprite rendering
- [ ] Implement sprite properties
- [ ] Add sprite management

**Files to Create**:
- `internal/nodes/sprite.go`
- `internal/graphics/sprite.go`

**Acceptance Criteria**:
- [ ] Sprites render correctly
- [ ] Sprite properties work
- [ ] Texture loading works

## 🎯 User Story 3: Physics Integration

### Task 3.1: Box2D Integration
**Priority**: High  
**Estimated Time**: 8 hours  
**Dependencies**: Task 2.3

**Description**: Integrate Box2D physics engine

**Tasks**:
- [ ] Initialize Box2D world
- [ ] Create physics body types
- [ ] Implement collision detection
- [ ] Add physics body management
- [ ] Create physics world updates

**Files to Create**:
- `internal/physics/physics_world.go`
- `internal/physics/physics_body.go`
- `internal/physics/collision.go`

**Acceptance Criteria**:
- [ ] Physics world initializes
- [ ] Bodies can be created
- [ ] Collision detection works
- [ ] Physics simulation runs

### Task 3.2: Physics Nodes
**Priority**: High  
**Estimated Time**: 10 hours  
**Dependencies**: Task 3.1

**Description**: Create physics node types

**Tasks**:
- [ ] Implement PhysicsBody2D base class
- [ ] Create StaticBody node
- [ ] Create RigidBody node
- [ ] Create KinematicBody node
- [ ] Add collision shape support

**Files to Create**:
- `internal/nodes/physics_body2d.go`
- `internal/nodes/static_body.go`
- `internal/nodes/rigid_body.go`
- `internal/nodes/kinematic_body.go`

**Acceptance Criteria**:
- [ ] All physics body types work
- [ ] Collision shapes work
- [ ] Physics simulation works

## 🎯 User Story 4: Audio System

### Task 4.1: Audio Synthesis
**Priority**: Medium  
**Estimated Time**: 12 hours  
**Dependencies**: Task 3.2

**Description**: Implement audio synthesis engine

**Tasks**:
- [ ] Create waveform generators
- [ ] Implement ADSR envelope
- [ ] Add audio effects
- [ ] Create audio buffer management
- [ ] Add audio playback

**Files to Create**:
- `internal/audio/synthesis.go`
- `internal/audio/waveforms.go`
- `internal/audio/envelope.go`
- `internal/audio/audio_buffer.go`

**Acceptance Criteria**:
- [ ] Waveforms generate correctly
- [ ] ADSR envelope works
- [ ] Audio playback works
- [ ] Effects work

### Task 4.2: Audio Nodes
**Priority**: Medium  
**Estimated Time**: 8 hours  
**Dependencies**: Task 4.1

**Description**: Create audio node types

**Tasks**:
- [ ] Implement SoundManager
- [ ] Create AudioPlayer node
- [ ] Create MusicPlayer node
- [ ] Add audio ducking
- [ ] Implement spatial audio

**Files to Create**:
- `internal/audio/sound_manager.go`
- `internal/nodes/audio_player.go`
- `internal/nodes/music_player.go`

**Acceptance Criteria**:
- [ ] Sound effects play
- [ ] Music plays
- [ ] Audio ducking works
- [ ] Spatial audio works

## 🎯 User Story 5: Graphics System

### Task 5.1: Complete Sprite Rendering
**Priority**: High  
**Estimated Time**: 10 hours  
**Dependencies**: Task 4.2

**Description**: Complete sprite rendering system

**Tasks**:
- [ ] Implement sprite batching
- [ ] Add sprite animations
- [ ] Create sprite atlas support
- [ ] Add sprite effects
- [ ] Optimize sprite rendering

**Files to Create**:
- `internal/graphics/sprite_batch.go`
- `internal/graphics/sprite_atlas.go`
- `internal/nodes/animated_sprite.go`

**Acceptance Criteria**:
- [ ] Sprite batching works
- [ ] Animations work
- [ ] Atlas support works
- [ ] Performance is good

### Task 5.2: Tilemap System
**Priority**: Medium  
**Estimated Time**: 12 hours  
**Dependencies**: Task 5.1

**Description**: Implement tilemap rendering

**Tasks**:
- [ ] Create tilemap data structures
- [ ] Implement tilemap rendering
- [ ] Add 8-layer support
- [ ] Create parallax scrolling
- [ ] Add tile flags system

**Files to Create**:
- `internal/graphics/tilemap.go`
- `internal/nodes/tilemap.go`
- `internal/graphics/parallax.go`

**Acceptance Criteria**:
- [ ] Tilemaps render correctly
- [ ] 8 layers work
- [ ] Parallax scrolling works
- [ ] Tile flags work

## 🎯 User Story 6: Lua API Integration

### Task 6.1: Complete Lua API
**Priority**: High  
**Estimated Time**: 16 hours  
**Dependencies**: Task 5.2

**Description**: Complete Lua API bindings

**Tasks**:
- [ ] Bind all node types to Lua
- [ ] Add input handling
- [ ] Implement math functions
- [ ] Add utility functions
- [ ] Create event system

**Files to Create**:
- `internal/lua/node_bindings.go`
- `internal/lua/input_bindings.go`
- `internal/lua/math_bindings.go`
- `internal/lua/event_system.go`

**Acceptance Criteria**:
- [ ] All APIs work in Lua
- [ ] Input handling works
- [ ] Math functions work
- [ ] Events work

### Task 6.2: Cart Loading
**Priority**: High  
**Estimated Time**: 8 hours  
**Dependencies**: Task 6.1

**Description**: Implement cart loading system

**Tasks**:
- [ ] Create cart format parser
- [ ] Implement asset loading
- [ ] Add cart validation
- [ ] Create cart manager
- [ ] Add error handling

**Files to Create**:
- `pkg/cart/loader.go`
- `pkg/cart/format.go`
- `pkg/cart/validator.go`
- `internal/engine/cart_manager.go`

**Acceptance Criteria**:
- [ ] Carts load correctly
- [ ] Assets load correctly
- [ ] Validation works
- [ ] Error handling works

## 🎯 User Story 7: Testing and Optimization

### Task 7.1: Unit Tests
**Priority**: High  
**Estimated Time**: 20 hours  
**Dependencies**: Task 6.2

**Description**: Create comprehensive unit tests

**Tasks**:
- [ ] Test all node types
- [ ] Test physics system
- [ ] Test audio system
- [ ] Test graphics system
- [ ] Test Lua API

**Files to Create**:
- `internal/nodes/node_test.go`
- `internal/physics/physics_test.go`
- `internal/audio/audio_test.go`
- `internal/graphics/graphics_test.go`
- `internal/lua/lua_test.go`

**Acceptance Criteria**:
- [ ] All tests pass
- [ ] Test coverage >90%
- [ ] Performance tests included

### Task 7.2: Integration Tests
**Priority**: High  
**Estimated Time**: 16 hours  
**Dependencies**: Task 7.1

**Description**: Create integration tests

**Tasks**:
- [ ] Test complete workflows
- [ ] Test cross-platform compatibility
- [ ] Test performance benchmarks
- [ ] Test memory usage
- [ ] Test error handling

**Files to Create**:
- `tests/integration/engine_test.go`
- `tests/integration/cart_test.go`
- `tests/performance/benchmarks.go`

**Acceptance Criteria**:
- [ ] All integration tests pass
- [ ] Performance targets met
- [ ] Memory usage within limits

### Task 7.3: Example Carts
**Priority**: Medium  
**Estimated Time**: 12 hours  
**Dependencies**: Task 7.2

**Description**: Create example carts

**Tasks**:
- [ ] Create Hello World cart
- [ ] Create Platformer demo
- [ ] Create Audio demo
- [ ] Create Sprite Animation demo
- [ ] Create Physics demo

**Files to Create**:
- `examples/hello_world/`
- `examples/platformer/`
- `examples/audio_demo/`
- `examples/sprite_animation/`
- `examples/physics_demo/`

**Acceptance Criteria**:
- [ ] All examples work
- [ ] Examples demonstrate features
- [ ] Examples are well-documented

## 🚀 Implementation Order

### Week 1-2: Foundation
1. Task 1.1: Project Setup
2. Task 1.2: SDL2 Window Setup
3. Task 1.3: Lua Integration

### Week 3-4: Node System
4. Task 2.1: Base Node Class
5. Task 2.2: Node2D Implementation
6. Task 2.3: Sprite Node

### Week 5-6: Physics
7. Task 3.1: Box2D Integration
8. Task 3.2: Physics Nodes

### Week 7-8: Audio
9. Task 4.1: Audio Synthesis
10. Task 4.2: Audio Nodes

### Week 9-10: Graphics
11. Task 5.1: Complete Sprite Rendering
12. Task 5.2: Tilemap System

### Week 11-12: Lua API
13. Task 6.1: Complete Lua API
14. Task 6.2: Cart Loading

### Week 13-14: Testing
15. Task 7.1: Unit Tests
16. Task 7.2: Integration Tests
17. Task 7.3: Example Carts

## 📊 Progress Tracking

### Checkpoints
- **Week 2**: Basic engine running
- **Week 4**: Node system working
- **Week 6**: Physics simulation working
- **Week 8**: Audio system working
- **Week 10**: Graphics system complete
- **Week 12**: Lua API complete
- **Week 14**: Testing complete

### Success Metrics
- [ ] All tasks completed on time
- [ ] All tests passing
- [ ] Performance targets met
- [ ] Documentation complete
- [ ] Examples working

---

**This task breakdown provides the detailed roadmap for implementing the RetroForge Engine. Each task should be completed in order, with regular checkpoints to ensure progress.**

---

*"Forge Your Retro Dreams" - RetroForge Engine Tasks* 🔨✨
EOF
    echo "✅ Task breakdown created"
else
    echo "ℹ️  Task breakdown already exists"
fi

echo ""
echo "🎉 RetroForge Engine development plan setup complete!"
echo ""
echo "📁 Files created:"
echo "  - specs/001-core-engine/plan.md"
echo "  - specs/001-core-engine/tasks.md"
echo ""
echo "🚀 Next steps:"
echo "  1. Review the development plan"
echo "  2. Set up your development environment"
echo "  3. Begin implementing Task 1.1: Project Setup"
echo "  4. Follow the task breakdown for implementation order"
echo ""
echo "📚 Documentation:"
echo "  - Constitution: memory/constitution.md"
echo "  - Specification: specs/001-core-engine/spec.md"
echo "  - Plan: specs/001-core-engine/plan.md"
echo "  - Tasks: specs/001-core-engine/tasks.md"
echo ""
echo "Happy coding! 🔨✨"
