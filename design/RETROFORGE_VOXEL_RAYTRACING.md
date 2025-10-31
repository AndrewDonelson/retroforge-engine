# RetroForge - Voxel & Raytracing Design Document

**Version:** 1.0  
**Date:** October 30, 2025  
**Status:** Phase 2 Feature Specification  
**Integration:** Post-Core Engine (Weeks 8-11)

---

## 🎯 Overview

This document specifies the integration of **voxel-based rendering** and **2.5D raytracing** capabilities into RetroForge, enabling developers to create retro-style 3D games reminiscent of Doom, Duke Nukem, and early 3D games. This feature maintains RetroForge's creative constraints while adding significant new creative possibilities.

### Core Philosophy
- **Simulated 3D**: No true 3D engine - everything is 2.5D raytracing
- **Retro Aesthetic**: Low resolution, pixelated, classic 3D game feel
- **Creative Constraints**: Voxel limitations encourage creativity
- **Performance First**: Must maintain 60 FPS on target platforms

---

## 🎮 Feature Goals

### Primary Objectives
1. **Voxel World Creation**: Build and edit 3D environments using voxels
2. **2.5D Raytracing**: Render 3D-like scenes using classic raytracing techniques
3. **Retro 3D Aesthetic**: Achieve Doom/Duke Nukem visual style
4. **Node Integration**: Seamlessly integrate with existing node system
5. **Performance**: Maintain 60 FPS within memory constraints

### Target Use Cases
- **Retro FPS Games**: Doom-style shooters with voxel environments
- **Voxel Builders**: Minecraft-style creative games
- **Puzzle Games**: 3D spatial puzzles with retro graphics
- **Adventure Games**: 3D exploration with 2D sprites
- **Educational Games**: 3D concepts with simple controls

---

## 🏗️ Technical Architecture

### System Overview

```
┌─────────────────────────────────────────────────┐
│              Voxel & Raytracing System          │
│                                                  │
│  ┌──────────────┐  ┌─────────────────────────┐  │
│  │ Voxel World  │  │   Raytracing Engine     │  │
│  │   Manager    │  │                         │  │
│  │              │  │  ┌─────────────────────┐ │  │
│  │ - Voxel Data │  │  │   Ray Caster        │ │  │
│  │ - Materials  │  │  │   - DDA Algorithm   │ │  │
│  │ - Lighting   │  │  │   - Hit Detection   │ │  │
│  │ - Animations │  │  │   - Texture Mapping │ │  │
│  └──────────────┘  │  └─────────────────────┘ │  │
│                    │                         │  │
│  ┌──────────────┐  │  ┌─────────────────────┐ │  │
│  │   Materials  │  │  │   Lighting System   │ │  │
│  │   System     │  │  │                     │ │  │
│  │              │  │  │  - Directional      │ │  │
│  │ - Textures   │  │  │  - Point Lights     │ │  │
│  │ - Properties │  │  │  - Ambient          │ │  │
│  │ - Shaders    │  │  │  - Shadows          │ │  │
│  └──────────────┘  │  └─────────────────────┘ │  │
└─────────────────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────┐
│            RetroForge Node System               │
│                                                  │
│  ┌──────────────┐  ┌─────────────────────────┐  │
│  │ VoxelWorld   │  │   RaytraceCamera        │  │
│  │   (Node)     │  │      (Node)             │  │
│  │              │  │                         │  │
│  │ - Voxel Data │  │  - Position/Rotation    │  │
│  │ - Materials  │  │  - FOV/Resolution      │  │
│  │ - Lighting   │  │  - Rendering Pipeline   │  │
│  └──────────────┘  │  └─────────────────────┘ │  │
│                    │                         │  │
│  ┌──────────────┐  │  ┌─────────────────────┐ │  │
│  │ VoxelMesh    │  │  │   VoxelLight        │ │  │
│  │   (Node)     │  │  │      (Node)         │ │  │
│  │              │  │  │                     │ │  │
│  │ - Static Geo │  │  │  - Light Properties │ │  │
│  │ - Animations │  │  │  - Color/Intensity  │ │  │
│  └──────────────┘  │  └─────────────────────┘ │  │
└─────────────────────────────────────────────────┘
```

### Core Components

#### 1. Voxel World Manager
**Purpose**: Manages voxel data, materials, and world state

**Key Features**:
- 3D voxel grid storage (64×32×64 default)
- Material system (8-16 materials max)
- Lighting calculations
- Voxel animations
- Memory-efficient data structures

#### 2. Raytracing Engine
**Purpose**: Renders 3D-like scenes using 2.5D raytracing

**Key Features**:
- DDA (Digital Differential Analyzer) algorithm
- Texture mapping
- Lighting calculations
- Shadow casting
- Performance optimizations

#### 3. Node System Integration
**Purpose**: Integrate with existing RetroForge node hierarchy

**New Node Types**:
- `VoxelWorld` - 3D world container
- `RaytraceCamera` - 3D camera with raytracing
- `VoxelMesh` - Static voxel geometry
- `VoxelLight` - Light sources
- `VoxelEntity` - Animated voxel objects

---

## 📊 Technical Specifications

### Performance Constraints

| Metric | Target | Rationale |
|--------|--------|-----------|
| **Resolution** | 480×270 | RetroForge standard |
| **Frame Rate** | 60 FPS | Smooth gameplay |
| **Frame Budget** | 16.67ms | 60 FPS requirement |
| **Memory Usage** | <8MB | Within 64MB WASM limit |
| **Voxel World Size** | 64×32×64 | Balance detail vs performance |
| **Voxel Size** | 8×8×8 pixels | Retro aesthetic |
| **Max Materials** | 16 | Memory efficiency |
| **Max Lights** | 8 | Performance limit |

### Memory Layout

```
Voxel World (64×32×64):
├── Voxel Data: 131,072 bytes (1 byte per voxel)
├── Material Data: 16×256 bytes = 4,096 bytes
├── Lighting Data: 131,072 bytes (1 byte per voxel)
├── Texture Atlas: 1,048,576 bytes (512×512×4)
├── Raytracing Buffers: 2,073,600 bytes (480×270×4×4)
└── Total: ~3.4MB (well within 64MB limit)
```

---

## 🎨 Visual Design

### Retro 3D Aesthetic

**Resolution**: 480×270 (16:9) - perfect for retro 3D
- **Doom**: 320×200 (4:3)
- **Duke Nukem 3D**: 320×200 (4:3)
- **RetroForge**: 480×270 (16:9) - 50% more pixels!

**Voxel Size**: 8×8×8 pixels
- Creates chunky, blocky aesthetic
- Maintains pixel-perfect rendering
- Fits retro game feel

**Color Palette**: Use existing 50-color system
- 16 base materials × 3 shades each
- Consistent with RetroForge branding
- Easy to create cohesive art

### Lighting Model

**Simple but Effective**:
- **Ambient**: Base lighting level
- **Directional**: Sun/moon light
- **Point Lights**: Torches, lamps, etc.
- **Shadows**: Simple shadow casting

**No Complex Shading**:
- Flat shading per voxel face
- No normal mapping
- No specular highlights
- Maintains retro aesthetic

---

## 🛠️ Implementation Plan

### Phase 1: Core Voxel System (Week 8)

#### Week 8.1: Voxel Data Structures
```go
// Voxel data structure
type Voxel struct {
    Material uint8  // 0-15 (16 materials max)
    Light    uint8  // 0-15 (lighting level)
    Flags    uint8  // Animation, special properties
}

// Voxel world
type VoxelWorld struct {
    Voxels   [64][32][64]Voxel
    Materials [16]Material
    Lights    []VoxelLight
    Dirty     bool  // Needs re-render
}
```

#### Week 8.2: Basic Voxel Operations
```go
// Core voxel functions
func (vw *VoxelWorld) SetVoxel(x, y, z int, material uint8)
func (vw *VoxelWorld) GetVoxel(x, y, z int) Voxel
func (vw *VoxelWorld) IsSolid(x, y, z int) bool
func (vw *VoxelWorld) Clear()
```

### Phase 2: Raytracing Engine (Week 9)

#### Week 9.1: DDA Algorithm Implementation
```go
// Raytracing core
type Ray struct {
    Origin    Vec3
    Direction Vec3
    MaxDist   float64
}

type HitInfo struct {
    Hit       bool
    Distance  float64
    Position  Vec3
    Normal    Vec3
    Material  uint8
    UV        Vec2
}

func (vw *VoxelWorld) CastRay(ray Ray) HitInfo
```

#### Week 9.2: Camera and Rendering
```go
// Raytracing camera
type RaytraceCamera struct {
    Position  Vec3
    Rotation  Vec3
    FOV       float64
    Resolution Vec2
    RenderBuffer []Color
}

func (cam *RaytraceCamera) Render(world *VoxelWorld)
```

### Phase 3: Node Integration (Week 10)

#### Week 10.1: Node System Extension
```lua
-- New node types
VoxelWorld.new({
  size = vec3(64, 32, 64),
  voxel_size = 8
})

RaytraceCamera.new({
  position = vec3(32, 16, 0),
  direction = vec3(0, 0, 1),
  fov = 90
})

VoxelMesh.new({
  voxel_data = {...},
  material = 1
})

VoxelLight.new({
  position = vec3(10, 5, 10),
  color = color(255, 255, 0),
  intensity = 1.0
})
```

#### Week 10.2: Lua API Integration
```lua
-- High-level voxel API
function _init()
  -- Create voxel world
  world = VoxelWorld.new({
    size = vec3(64, 32, 64),
    voxel_size = 8
  })
  
  -- Add some voxels
  world:set_voxel(10, 5, 10, 1)  -- stone
  world:set_voxel(11, 5, 10, 2)  -- grass
  
  -- Create raytracing camera
  camera = RaytraceCamera.new({
    position = vec3(32, 16, 0),
    direction = vec3(0, 0, 1)
  })
  
  Scene:add_child(world)
  Scene:add_child(camera)
end

function _draw()
  -- Raytrace the world
  camera:render(world)
end
```

### Phase 4: Advanced Features (Week 11)

#### Week 11.1: Lighting System
```go
// Lighting calculations
type VoxelLight struct {
    Position  Vec3
    Color     Color
    Intensity float64
    Type      LightType  // Directional, Point, Ambient
}

func (vw *VoxelWorld) CalculateLighting()
func (vw *VoxelWorld) CastShadow(light VoxelLight, pos Vec3) bool
```

#### Week 11.2: Performance Optimization
- **Frustum Culling**: Only render visible voxels
- **Level of Detail**: Reduce detail for distant voxels
- **Dirty Flagging**: Only re-render changed areas
- **Memory Pooling**: Reuse ray objects

---

## 📚 Open Source Libraries

### Primary Libraries

#### 1. **Go-Raytracing Libraries**
- **github.com/kingsleyliao/ray-tracer**: Go ray tracer with realistic lighting
- **github.com/markphelps/raytracer**: Educational Go ray tracer implementation
- **Custom DDA Implementation**: Optimized for voxel raytracing

#### 2. **Voxel Processing Libraries**
- **github.com/klauspost/compress**: Efficient data compression for voxel data
- **github.com/ByteArena/box2d**: Physics integration for voxel objects
- **Custom Voxel Engine**: Built specifically for RetroForge constraints

#### 3. **Graphics Libraries**
- **github.com/veandco/go-sdl2**: SDL2 bindings (already in use)
- **github.com/go-gl/gl**: OpenGL bindings for advanced rendering
- **Custom Software Renderer**: Optimized for 2.5D raytracing

### Integration Strategy

**Phase 1**: Use existing Go libraries as reference
**Phase 2**: Implement custom optimized versions
**Phase 3**: Integrate with RetroForge's existing systems

---

## 🎮 Example Games

### 1. Retro FPS (Doom-style)
```lua
function _init()
  -- Create maze world
  world = VoxelWorld.new({size = vec3(64, 32, 64)})
  
  -- Build walls
  for x = 0, 63 do
    for z = 0, 63 do
      if (x + z) % 8 == 0 then
        world:set_voxel(x, 0, z, 1)  -- floor
        world:set_voxel(x, 8, z, 1)  -- ceiling
      end
    end
  end
  
  -- Create camera
  camera = RaytraceCamera.new({
    position = vec3(32, 4, 0),
    direction = vec3(0, 0, 1)
  })
  
  Scene:add_child(world)
  Scene:add_child(camera)
end

function _update()
  -- Move camera
  if Input.is_action_pressed("left") then
    camera:rotate_y(-0.05)
  end
  if Input.is_action_pressed("right") then
    camera:rotate_y(0.05)
  end
  if Input.is_action_pressed("forward") then
    camera:move_forward(0.5)
  end
end

function _draw()
  camera:render(world)
end
```

### 2. Voxel Builder (Minecraft-style)
```lua
function _init()
  world = VoxelWorld.new({size = vec3(32, 16, 32)})
  camera = RaytraceCamera.new()
  
  -- Build a simple house
  for x = 10, 20 do
    for z = 10, 20 do
      world:set_voxel(x, 0, z, 1)  -- floor
      world:set_voxel(x, 8, z, 1)  -- ceiling
    end
  end
  
  -- Add walls
  for y = 0, 8 do
    for x = 10, 20 do
      world:set_voxel(x, y, 10, 2)  -- wall
      world:set_voxel(x, y, 20, 2)  -- wall
    end
    for z = 10, 20 do
      world:set_voxel(10, y, z, 2)  -- wall
      world:set_voxel(20, y, z, 2)  -- wall
    end
  end
end
```

### 3. 3D Puzzle Game
```lua
function _init()
  world = VoxelWorld.new({size = vec3(16, 8, 16)})
  camera = RaytraceCamera.new()
  
  -- Create puzzle blocks
  world:set_voxel(5, 1, 5, 3)   -- red block
  world:set_voxel(10, 1, 5, 4)  -- blue block
  world:set_voxel(5, 1, 10, 5)  -- green block
  
  -- Goal position
  goal_pos = vec3(10, 1, 10)
end

function _update()
  -- Move blocks with arrow keys
  if Input.is_action_just_pressed("left") then
    move_block(vec3(-1, 0, 0))
  end
end
```

---

## 🔧 Development Tools

### Voxel Editor Integration

**Web-based Voxel Editor**:
- 3D viewport with voxel placement
- Material selection palette
- Lighting preview
- Export to RetroForge format

**Features**:
- Click to place/remove voxels
- Material brush tool
- Copy/paste voxel regions
- Undo/redo support
- Real-time preview

### Asset Pipeline

**Voxel File Format**:
```json
{
  "version": "1.0",
  "size": [64, 32, 64],
  "voxel_size": 8,
  "materials": [
    {"id": 0, "name": "air", "color": "#000000"},
    {"id": 1, "name": "stone", "color": "#808080"},
    {"id": 2, "name": "grass", "color": "#00FF00"}
  ],
  "voxels": [0, 1, 2, ...]  // Flat array
}
```

**Import/Export**:
- MagicaVoxel .vox format support
- RetroForge .rfs/.rfe integration
- Texture atlas generation

---

## 📈 Performance Optimization

### Rendering Optimizations

**1. Frustum Culling**
```go
func (cam *RaytraceCamera) IsVoxelVisible(pos Vec3) bool {
    // Check if voxel is within camera frustum
    return cam.frustum.Contains(pos)
}
```

**2. Level of Detail**
```go
func (vw *VoxelWorld) GetLODLevel(distance float64) int {
    if distance < 32 { return 1 }  // Full detail
    if distance < 64 { return 2 }  // Half detail
    return 4  // Quarter detail
}
```

**3. Dirty Flagging**
```go
func (vw *VoxelWorld) SetVoxel(x, y, z int, material uint8) {
    vw.Voxels[x][y][z].Material = material
    vw.Dirty = true  // Mark for re-render
}
```

**4. Memory Pooling**
```go
var rayPool = sync.Pool{
    New: func() interface{} {
        return &Ray{}
    },
}

func GetRay() *Ray {
    return rayPool.Get().(*Ray)
}

func ReturnRay(r *Ray) {
    rayPool.Put(r)
}
```

### Memory Management

**Voxel Data Compression**:
- Run-length encoding for empty space
- Material indexing to reduce data size
- Lazy loading for large worlds

**Texture Optimization**:
- 512×512 texture atlas (1MB)
- Material-based texture coordinates
- Mipmapping for distant voxels

---

## 🧪 Testing Strategy

### Unit Tests
```go
func TestVoxelWorld_SetVoxel(t *testing.T) {
    world := NewVoxelWorld(64, 32, 64)
    world.SetVoxel(10, 5, 10, 1)
    
    voxel := world.GetVoxel(10, 5, 10)
    assert.Equal(t, uint8(1), voxel.Material)
}

func TestRaytracing_CastRay(t *testing.T) {
    world := NewVoxelWorld(64, 32, 64)
    world.SetVoxel(10, 5, 10, 1)
    
    ray := Ray{
        Origin: Vec3{0, 5, 0},
        Direction: Vec3{1, 0, 0},
        MaxDist: 20,
    }
    
    hit := world.CastRay(ray)
    assert.True(t, hit.Hit)
    assert.Equal(t, 10.0, hit.Distance)
}
```

### Performance Tests
```go
func BenchmarkVoxelRendering(b *testing.B) {
    world := NewVoxelWorld(64, 32, 64)
    camera := NewRaytraceCamera()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        camera.Render(world)
    }
}

func BenchmarkRayCasting(b *testing.B) {
    world := NewVoxelWorld(64, 32, 64)
    ray := Ray{Origin: Vec3{0, 0, 0}, Direction: Vec3{1, 0, 0}}
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        world.CastRay(ray)
    }
}
```

### Integration Tests
```lua
-- Test cart for voxel functionality
function test_voxel_placement()
  world = VoxelWorld.new({size = vec3(16, 8, 16)})
  world:set_voxel(5, 2, 5, 1)
  
  local voxel = world:get_voxel(5, 2, 5)
  assert(voxel.material == 1, "Voxel not placed correctly")
end

function test_raytracing()
  world = VoxelWorld.new({size = vec3(16, 8, 16)})
  world:set_voxel(10, 2, 5, 1)
  
  camera = RaytraceCamera.new({
    position = vec3(0, 2, 0),
    direction = vec3(1, 0, 0)
  })
  
  camera:render(world)
  -- Check that voxel is visible
end
```

---

## 📚 Documentation Plan

### User Documentation

**1. Voxel Basics**
- What are voxels?
- How to create voxel worlds
- Material system explanation
- Lighting basics

**2. Raytracing Guide**
- Camera controls
- Movement and rotation
- Performance tips
- Common patterns

**3. Example Projects**
- Complete game tutorials
- Code walkthroughs
- Best practices
- Performance optimization

### API Reference

**VoxelWorld API**:
```lua
-- Creation
world = VoxelWorld.new({size = vec3(64, 32, 64)})

-- Voxel operations
world:set_voxel(x, y, z, material)
world:get_voxel(x, y, z)
world:is_solid(x, y, z)
world:clear()

-- Lighting
world:add_light(light)
world:remove_light(light)
world:calculate_lighting()
```

**RaytraceCamera API**:
```lua
-- Creation
camera = RaytraceCamera.new({
  position = vec3(0, 0, 0),
  direction = vec3(0, 0, 1),
  fov = 90
})

-- Movement
camera:move_forward(distance)
camera:move_right(distance)
camera:move_up(distance)
camera:rotate_y(angle)
camera:rotate_x(angle)

-- Rendering
camera:render(world)
```

---

## 🚀 Success Metrics

### Technical Goals
- ✅ 60 FPS stable on target platforms
- ✅ <8MB memory usage for voxel world
- ✅ <16.67ms frame budget maintained
- ✅ Seamless integration with node system

### User Experience Goals
- ✅ Easy to learn (clear documentation)
- ✅ Quick to prototype (simple API)
- ✅ Fun to use (immediate visual feedback)
- ✅ Creative constraints encourage finished projects

### Community Goals
- ✅ 10+ voxel-based example carts
- ✅ Active community discussions
- ✅ Tutorial series completion
- ✅ Performance optimization guides

---

## 🔄 Future Enhancements

### Phase 3 Possibilities
- **Voxel Animations**: Moving voxel objects
- **Advanced Lighting**: Multiple light types, shadows
- **Particle Effects**: 3D particle systems
- **Sound Integration**: 3D positional audio
- **Physics Integration**: Voxel physics with Box2D

### Advanced Features
- **Voxel Streaming**: Large worlds with loading
- **Multiplayer**: Shared voxel worlds
- **VR Support**: 3D headset integration
- **Advanced Shaders**: Custom material effects

---

## 📋 Implementation Checklist

### Week 8: Core Voxel System
- [ ] Voxel data structures
- [ ] Basic voxel operations
- [ ] Material system
- [ ] Memory management
- [ ] Unit tests

### Week 9: Raytracing Engine
- [ ] DDA algorithm implementation
- [ ] Ray-voxel intersection
- [ ] Camera system
- [ ] Basic rendering
- [ ] Performance optimization

### Week 10: Node Integration
- [ ] VoxelWorld node
- [ ] RaytraceCamera node
- [ ] VoxelMesh node
- [ ] VoxelLight node
- [ ] Lua API bindings

### Week 11: Advanced Features
- [ ] Lighting system
- [ ] Shadow casting
- [ ] Performance optimization
- [ ] Example games
- [ ] Documentation

---

## 🎉 Conclusion

The voxel and raytracing feature will significantly enhance RetroForge's capabilities while maintaining its retro aesthetic and creative constraints. By implementing 2.5D raytracing with voxels, RetroForge will offer a unique platform for creating retro-style 3D games that feel authentic to the classic era while leveraging modern development tools.

**Key Benefits**:
- 🎮 **Unique Positioning**: Only fantasy console with voxel raytracing
- 🎨 **Creative Constraints**: Voxel limitations encourage creativity
- ⚡ **Performance**: Optimized for 60 FPS on all platforms
- 🔧 **Easy to Use**: Simple API with powerful capabilities
- 📚 **Well Documented**: Comprehensive guides and examples

**This feature will make RetroForge the premier platform for retro 3D game development!** 🔨✨

---

**Document Version History**:
- v1.0 - Initial specification (October 30, 2025)

**Next Steps**:
1. Review and approve specification
2. Begin Week 8 implementation
3. Create example voxel games
4. Develop web-based voxel editor
5. Launch with community showcase

---

*"Forge Your Retro Dreams in 3D!" - RetroForge Voxel System* 🎮✨
