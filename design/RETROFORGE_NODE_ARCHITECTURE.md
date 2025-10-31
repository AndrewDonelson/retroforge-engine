# RetroForge - Node System & Physics Architecture

**Version:** 1.0  
**Date:** October 29, 2025  
**Status:** Specification Complete

---

## Overview

RetroForge uses a **Godot-inspired node system** combined with **Box2D physics** to provide a powerful yet approachable game development environment. This document details the architecture and implementation strategy.

---

## Philosophy

### Why a Node System?

1. **Familiar to modern developers** - Unity and Godot users feel at home
2. **Reduces boilerplate** - Less manual management of game objects
3. **Automatic organization** - Scene graph handles update/draw order
4. **Component composition** - Build complex behaviors from simple nodes
5. **Still accessible** - Low-level API available for direct control

### Dual API Approach

RetroForge provides **two complementary APIs**:

**Node API** (High-level, recommended)
```lua
player = RigidBody.new({position = vec2(240, 135)})
player:apply_impulse(vec2(0, -300))  -- Jump
Scene:add_child(player)
```

**Direct API** (Low-level, PICO-8-style)
```lua
x, y = 240, 135
vy = -5  -- Jump
spr(1, x, y)
```

Both can be mixed in the same cart!

---

## Scene Graph Architecture

### Hierarchy

```
Scene (root)
├── Camera
├── Player (RigidBody)
│   ├── Sprite
│   ├── CollisionShape
│   └── AudioPlayer
├── TileMap (layer 1)
├── TileMap (layer 2)
├── Enemy1 (RigidBody)
│   ├── AnimatedSprite
│   └── CollisionShape
├── Enemy2 (RigidBody)
│   └── ...
└── UI (CanvasLayer)
    ├── ScoreLabel
    └── HealthBar
```

### Node Lifecycle

```
1. Node created (new)
2. Node added to tree (add_child)
3. _ready() called
4. Every frame:
   a. _update() called (depth-first)
   b. _draw() called (depth-first with z-index)
5. Node removed (queue_free)
6. _on_removed() called
```

### Transform Propagation

Nodes inherit transforms from parents:
- **Local transform:** Position/rotation/scale relative to parent
- **Global transform:** Actual world position (calculated)

```lua
-- Parent at (100, 50)
parent.position = vec2(100, 50)

-- Child at local (20, 10)
child.position = vec2(20, 10)

-- Child's global position is (120, 60)
print(child.global_position)  -- vec2(120, 60)
```

---

## Node Type Hierarchy

### Abstract Base Classes

**Node** (abstract)
- Base for all nodes
- Name, parent, children
- Active state
- Virtual methods: _ready(), _update(), _draw()

**Node2D** (abstract, extends Node)
- Adds spatial properties
- Position, rotation, scale
- Z-index for draw order
- Transform calculations

**PhysicsBody2D** (abstract, extends Node2D)
- Base for physics nodes
- Collision layers/masks
- Shape management
- Collision callbacks

### Concrete Node Classes

```
Renderable Nodes:
├── Sprite              - Single sprite
├── AnimatedSprite      - Animated sprite sequences
├── TileMap             - Map layer rendering
├── ParticleEmitter     - Particle effects
└── CanvasLayer         - UI layer

Physics Nodes:
├── StaticBody          - Static geometry
├── RigidBody           - Dynamic physics
└── KinematicBody       - Player-controlled

Utility Nodes:
├── Camera              - Viewport control
├── CollisionShape      - Collision geometry
├── AudioPlayer         - Sound effects
├── MusicPlayer         - Music tracks
└── Timer               - Delayed actions
```

---

## Physics Integration (Box2D)

### Box2D Overview

Box2D is a proven 2D physics engine used in:
- Angry Birds
- Limbo
- Crayon Physics Deluxe
- Hundreds of other games

**Features we use:**
- Rigid body dynamics
- Collision detection
- Joints and constraints
- Raycasting
- Contact filtering
- Continuous collision detection

### Physics World

One Box2D world per cart:
- Gravity: Configurable (default: vec2(0, 980) - 980 pixels/s²)
- Timestep: Fixed at 1/60 second
- Velocity iterations: 8
- Position iterations: 3

### Body Types

**StaticBody**
- Doesn't move
- Zero mass
- Used for: Walls, platforms, obstacles
- Box2D type: b2_staticBody

**RigidBody**
- Full physics simulation
- Affected by forces and gravity
- Used for: Crates, balls, ragdolls
- Box2D type: b2_dynamicBody

**KinematicBody**
- Moved by code, not physics
- Infinite mass
- Used for: Player characters, moving platforms
- Box2D type: b2_kinematicBody

### Collision Shapes

Supported shapes:
1. **Rectangle** - Most common, fast
2. **Circle** - Very fast, good for balls
3. **Capsule** - Good for characters
4. **Polygon** - Up to 8 vertices, convex only

```lua
-- Create a physics body with shape
player = RigidBody.new({
  position = vec2(240, 135),
  width = 16,
  height = 16
})

-- Add circular collision
local shape = CollisionShape.circle(8)
player:add_collision_shape(shape)
```

### Collision Layers

8 collision layers (0-7) with mask filtering:

```lua
-- Player on layer 0, collides with layers 0, 1, 2
player.collision_layer = 0  -- Which layer I'm on
player.collision_mask = 7   -- Binary: 0b00000111 (layers 0,1,2)

-- Enemy on layer 1, collides with layers 0, 1
enemy.collision_layer = 1
enemy.collision_mask = 3    -- Binary: 0b00000011 (layers 0,1)
```

**Common layer setup:**
- Layer 0: Player
- Layer 1: Enemies
- Layer 2: World (walls, platforms)
- Layer 3: Collectibles
- Layer 4: Triggers
- Layer 5-7: Custom

### Physics Callbacks

**Collision Detection:**
```lua
function player:on_body_entered(other)
  if other:is_in_group("enemy") then
    self:take_damage()
  end
end

function player:on_body_exited(other)
  -- Collision ended
end
```

**Sensor Bodies:**
```lua
-- Create trigger zone (no collision response)
trigger = StaticBody.new({
  position = vec2(200, 100),
  is_sensor = true  -- Detects but doesn't block
})

function trigger:on_body_entered(other)
  print("Player entered trigger!")
end
```

---

## Audio System Architecture

### Three-Tier Design

```
┌─────────────────────────────────────┐
│       SoundManager (Singleton)      │
│  - Global volume controls           │
│  - Channel allocation               │
│  - Audio ducking                    │
│  - Quick play functions             │
└─────────────────────────────────────┘
           │                │
           ▼                ▼
    ┌────────────┐   ┌──────────────┐
    │ AudioPlayer│   │ MusicPlayer  │
    │  (Node)    │   │   (Node)     │
    │            │   │              │
    │ - SFX      │   │ - Music      │
    │ - Multiple │   │ - Single     │
    │ - Spatial  │   │ - Crossfade  │
    └────────────┘   └──────────────┘
```

### SoundManager (Automatic)

Created automatically in every cart. No initialization needed.

**Responsibilities:**
1. Master volume control
2. SFX/Music bus management
3. Channel allocation (8 channels)
4. Audio ducking
5. Quick play shortcuts

**Global access:**
```lua
Sound.master_volume = 0.8
Sound.sfx_volume = 1.0
Sound.music_volume = 0.7

Sound:play_sfx(5)  -- Quick SFX playback
Sound:play_music(1, 1.0)  -- Play track 1 with 1s crossfade
```

### AudioPlayer (Node)

For sound effects. Multiple can play simultaneously.

**Features:**
- Position-based volume (closer = louder)
- Per-instance volume/pitch
- Looping support
- Bus routing (sfx/master)

**Usage:**
```lua
-- Create and play
sfx_player = AudioPlayer.new({sound_index = 5, volume = 0.8})
player:add_child(sfx_player)  -- Child of player for spatial audio
sfx_player:play()

-- One-shot quick play
Sound:play_sfx(5)
```

### MusicPlayer (Node)

For music tracks. Only one track at a time.

**Features:**
- Automatic cross-fading
- Seamless looping
- Tempo control
- Position seeking
- Fade in/out

**Usage:**
```lua
-- Create music player (usually once per cart)
music = MusicPlayer.new()
Scene:add_child(music)

-- Play track with 2-second crossfade
music:play(1, 2.0)

-- Change music dynamically
music:play(2, 1.5)  -- Fades out track 1, fades in track 2
```

### Audio Ducking

Automatically lowers music when sound effects play:

```lua
-- Enable ducking: music drops to 30% when SFX plays
Sound:enable_ducking(true, 0.3)

-- Now when any SFX plays:
-- 1. Music volume drops to 30%
-- 2. SFX plays
-- 3. Music returns to 100% after SFX ends
```

---

## Implementation Plan

### Phase 1: Core Node System

**Week 1: Base Node Classes**
```go
// Go implementation
package node

type Node struct {
    name     string
    parent   *Node
    children []*Node
    active   bool
    luaRef   int  // Reference to Lua object
}

type Node2D struct {
    Node
    position Vec2
    rotation float64
    scale    Vec2
    zIndex   int
}
```

**Lua bindings:**
```lua
-- Expose to Lua via gopher-lua
L.SetGlobal("Node", L.NewFunction(nodeConstructor))
L.SetGlobal("Scene", sceneTable)
```

### Phase 2: Physics Integration

**Week 2: Box2D Wrapper**
```go
package physics

import "github.com/ByteArena/box2d"

type PhysicsWorld struct {
    world *box2d.B2World
    bodies map[int]*box2d.B2Body
}

type PhysicsBody2D struct {
    Node2D
    body *box2d.B2Body
    bodyType BodyType
}
```

**Synchronization:**
```go
func (pb *PhysicsBody2D) Update() {
    // Sync Box2D position to node position
    pos := pb.body.GetPosition()
    pb.position.X = pos.X
    pb.position.Y = pos.Y
    pb.rotation = pb.body.GetAngle()
}
```

### Phase 3: Audio System

**Week 3: Audio Architecture**
```go
package audio

type SoundManager struct {
    masterVolume float64
    sfxVolume    float64
    musicVolume  float64
    players      []*AudioPlayer
    musicPlayer  *MusicPlayer
    channels     [8]Channel
}

type AudioPlayer struct {
    Node
    soundIndex int
    volume     float64
    channel    *Channel
}

type MusicPlayer struct {
    Node
    trackIndex  int
    fadeTime    float64
    currentTrack *Track
    nextTrack   *Track
}
```

### Phase 4: Scene Graph Rendering

**Week 4: Render Pipeline**
```go
func (s *Scene) Render() {
    // 1. Sort nodes by z-index
    // 2. Apply camera transform
    // 3. Traverse tree depth-first
    // 4. Call node._draw() for each
    
    for _, layer := range s.getSortedLayers() {
        for _, node := range layer.nodes {
            if node.active {
                node.Render()
            }
        }
    }
}
```

---

## Example Game

### Platformer with Physics

```lua
function _init()
  -- Create player
  player = RigidBody.new({
    position = vec2(240, 135),
    width = 16,
    height = 16,
    collision_layer = 0,
    collision_mask = 6,  -- Collide with world (bit 1) and enemies (bit 2)
    fixed_rotation = true
  })
  
  -- Add sprite to player
  local sprite = Sprite.new({sprite_index = 1})
  player:add_child(sprite)
  
  -- Add collision shape
  local shape = CollisionShape.rect(16, 16)
  player:add_collision_shape(shape)
  
  -- Add to scene
  Scene:add_child(player)
  
  -- Create camera that follows player
  camera = Camera.new()
  camera:follow(player, 0.1)  -- Smooth following
  Scene:add_child(camera)
  
  -- Load level
  tilemap = TileMap.new({
    level_name = "level_1",
    layer_id = 1
  })
  Scene:add_child(tilemap)
  
  -- Background layer with parallax
  bg = TileMap.new({
    level_name = "level_1",
    layer_id = 7,
    parallax = 0.5  -- Moves at half speed
  })
  Scene:add_child(bg)
  
  -- Start music
  music = MusicPlayer.new()
  Scene:add_child(music)
  music:play(1)
  
  -- Enable audio ducking
  Sound:enable_ducking(true, 0.3)
end

-- Override player's update function
function player:_update()
  local move_force = 500
  
  if Input.is_action_pressed("left") then
    self:apply_force(vec2(-move_force, 0))
  end
  
  if Input.is_action_pressed("right") then
    self:apply_force(vec2(move_force, 0))
  end
  
  if Input.is_action_just_pressed("jump") and self:is_on_floor() then
    self:apply_impulse(vec2(0, -300))
    Sound:play_sfx(5)  -- Jump sound
  end
end

-- Handle collisions
function player:on_body_entered(other)
  if other:is_in_group("enemy") then
    self:take_damage()
    Sound:play_sfx(10)  -- Hit sound
  elseif other:is_in_group("coin") then
    other:queue_free()
    score = score + 10
    Sound:play_sfx(3)  -- Coin sound
  end
end
```

---

## Performance Considerations

### Node Count Limits

**Target performance:**
- 60 FPS stable
- Up to 500 nodes in scene
- Up to 100 active physics bodies
- 8 simultaneous audio channels

### Optimization Strategies

**1. Object Pooling**
```lua
-- Pool for bullets
bullet_pool = {}

function get_bullet()
  if #bullet_pool > 0 then
    return table.remove(bullet_pool)
  end
  return RigidBody.new({...})
end

function return_bullet(bullet)
  bullet.active = false
  table.insert(bullet_pool, bullet)
end
```

**2. Spatial Partitioning**
- Quadtree for broad-phase collision
- Only update visible nodes
- Cull offscreen sprites

**3. Physics Optimization**
- Use StaticBody for non-moving objects
- Enable sleeping for inactive bodies
- Simplify collision shapes
- Limit contact points

**4. Audio Optimization**
- Pool AudioPlayer nodes
- Stop distant sounds
- Limit simultaneous SFX
- Stream music, don't load all

---

## Testing Strategy

### Unit Tests (Go)
```go
func TestNodeHierarchy(t *testing.T) {
    parent := NewNode()
    child := NewNode()
    parent.AddChild(child)
    assert.Equal(t, parent, child.Parent())
}

func TestPhysicsBodySync(t *testing.T) {
    body := NewRigidBody()
    body.ApplyImpulse(Vec2{0, -100})
    body.Update()
    assert.True(t, body.Position.Y < 0)
}
```

### Integration Tests (Lua)
```lua
-- Test cart
function test_player_jump()
  player = RigidBody.new({position = vec2(0, 0)})
  player:apply_impulse(vec2(0, -300))
  
  -- Simulate 1 frame
  player:_update()
  
  assert(player.velocity.y < 0, "Player should move up")
end
```

---

## Migration from PICO-8

For PICO-8 users, here's the mapping:

```lua
-- PICO-8 style (still works!)
function _update()
  if btn(2) then x = x - 2 end
  spr(1, x, y)
end

-- RetroForge node style (recommended)
function _init()
  player = Sprite.new({
    sprite_index = 1,
    position = vec2(x, y)
  })
  Scene:add_child(player)
end

function player:_update()
  if Input.is_action_pressed("left") then
    self.position.x = self.position.x - 2
  end
end
```

Both styles work! Use what feels comfortable.

---

## References

- **Box2D Manual:** https://box2d.org/documentation/
- **Godot Scene System:** https://docs.godotengine.org/en/stable/getting_started/step_by_step/nodes_and_scenes.html
- **Unity GameObjects:** https://docs.unity3d.com/Manual/class-GameObject.html
- **box2d-go:** https://github.com/ByteArena/box2d

---

## Next Steps

1. ✅ Architecture finalized
2. 🔲 Implement base Node class in Go
3. 🔲 Integrate Box2D-go
4. 🔲 Create Lua bindings
5. 🔲 Build example platformer
6. 🔲 Performance testing

*The node system brings modern game engine patterns to retro game development!* 🎮✨
