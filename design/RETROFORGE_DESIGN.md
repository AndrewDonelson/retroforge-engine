# RetroForge Fantasy Console - Design Document

**Version:** 0.3  
**Date:** October 29, 2025  
**Status:** Design Phase - Ready for Implementation

## Project Overview

**RetroForge** is a modern fantasy console inspired by PICO-8, designed with a clean separation between the runtime engine and development tools. The name captures both the retro aesthetic and the creative act of forging games. The project aims to provide creative constraints that encourage finished projects while offering more capability than 8-bit consoles.

### Core Technology Stack

- **Engine:** Go 1.23, compiled to WASM for web and native executables
- **Frontend:** Next.js 16, TypeScript, TailwindCSS
- **Runtime:** gopher-lua (Lua 5.1)
- **File Formats:** 
  - `.rfs` (RetroForge Source) - Editable source carts
  - `.rfe` (RetroForge Executable) - Distributed executable carts

### Philosophy

- **Open Source Model**: All carts are distributed as readable Lua code + assets
- **Learn by Example**: Embrace community sharing and remixing
- **Creative Constraints**: Limitations spark creativity and help developers finish projects
- **Modern Architecture**: Clean separation between engine (Go) and tooling (Web)

---

## Technical Specifications

### Display

- **Resolution:** 480×270 (landscape) or 270×480 (portrait)
- **Aspect Ratio:** 16:9 (modern device friendly)
- **Scaling:** Integer scaling for pixel-perfect display on modern screens
  - 2x = 960×540 (qHD)
  - 3x = 1440×810
  - 4x = 1920×1080 (Full HD)
- **Color Depth:** 256 colors (predefined palettes + custom option)
- **Color Palette Structure:**
  - Black (1 color)
  - White (1 color)
  - 16 base colors × 3 shades each (base, highlight, shadow) = 48 colors
  - Total per palette: 50 colors
  - Multiple predefined palettes based on color theory
  - Users can create custom palettes
- **Refresh Rate:** 60 FPS

### Code

- **Language:** Lua 5.1
- **Token Limit:** 16,384 tokens
- **Runtime:** gopher-lua (pure Go implementation)
- **Distribution:** Plain text Lua files

### Sprites

- **Total Slots:** 256 sprite slots
- **Supported Sizes:** 4×4, 8×8, 16×16, 32×32 pixels (1:1 aspect ratio)
- **Slot Allocation:**
  - 4×4 sprite = 1 slot
  - 8×8 sprite = 1 slot
  - 16×16 sprite = 4 slots (2×2 grid)
  - 32×32 sprite = 16 slots (4×4 grid)
- **Format:** PNG sprite sheet with JSON metadata

### Map System

- **Tile Size:** 16×16 pixels (fixed)
- **Map Dimensions:** Configurable per cart (default 128×128 tiles)
- **Layers:** 8 layers (0-7)
  - **Layer 0:** UI layer (no parallax, always on top)
  - **Layers 1-7:** Game world layers (parallax support, 1=foreground → 7=background)
- **Parallax:** Each layer can have different scroll speeds
- **Tile Flags:** 8 boolean flags per tile for collision/properties
  - Flag 0: Solid (collision)
  - Flag 1-7: User-definable (hazard, water, ice, etc.)
- **Data Storage:** Array of tile indices + flag data

### Audio

- **Channels:** 8 simultaneous audio channels
- **Sound Effects:** 128 definable patterns
- **Music:** Pattern-based sequencer
- **Synthesis:** Pure synthesis (no samples)
- **Waveforms:** Square, Triangle, Sawtooth, Sine, Noise
- **Volume Envelope:** ADSR (Attack, Decay, Sustain, Release)
- **Effects:**
  - Vibrato (pitch wobble)
  - Arpeggio (rapid note cycling for chord emulation)
  - Pitch slide (portamento)
- **Philosophy:** Authentic chip-tune aesthetic with small memory footprint

### Audio System Architecture

**Three-tier audio system:**

1. **AudioPlayer** (Sound Effects Node)
   - Plays individual sound effects
   - Multiple instances can play simultaneously
   - Positioned in scene graph
   - Volume, pitch, loop control per instance

2. **MusicPlayer** (Music Node)
   - Plays music tracks
   - Only one music track at a time
   - Automatic cross-fading between tracks
   - Tempo, volume control
   - Seamless looping

3. **SoundManager** (Automatic Singleton)
   - Exists in every game automatically
   - No code required to create
   - Manages all AudioPlayer and MusicPlayer instances
   - Global volume controls (master, sfx, music)
   - Audio ducking (lower music when SFX plays)
   - Handles channel allocation
   - Accessible via `Sound` global

### Physics System

- **Engine:** Box2D-go (port of Box2D)
- **Features:**
  - Rigid body dynamics
  - Collision detection (AABB, circles, polygons)
  - Joints and constraints
  - Raycasting
  - Continuous collision detection
  - Sleeping/active body optimization
- **Body Types:**
  - Static (walls, platforms - no physics)
  - Dynamic (affected by forces and gravity)
  - Kinematic (player-controlled, not affected by forces)
- **Collision Layers:** 8 collision layers with mask filtering
- **Performance:** Optimized for 60 FPS with reasonable object counts

### Node System (GameObject Architecture)

RetroForge uses a **scene graph / node system** similar to Godot, providing:
- Object hierarchy (parent/child relationships)
- Automatic update and draw order
- Component-based design
- Built-in common game objects
- Extensible via Lua

**Base Node Types:**
```
Node (abstract base)
├── Node2D (spatial nodes with transform)
│   ├── Sprite
│   ├── AnimatedSprite
│   ├── TileMap
│   ├── Camera
│   ├── ParallaxBackground
│   └── CanvasLayer
├── PhysicsBody2D (Box2D integration)
│   ├── StaticBody
│   ├── RigidBody
│   └── KinematicBody
├── CollisionShape
├── AudioPlayer (sound effects)
├── MusicPlayer (music tracks)
├── Timer
└── ParticleEmitter
```

**Scene Structure:**
- Each cart has a root Scene node
- Nodes can be nested (parent/child hierarchy)
- Nodes automatically update and draw based on tree order
- Z-index controls draw order within layers

### Persistent Storage

- **Save Data:** 1024 bytes per cart
- **Format:** Key-value storage (64 slots × 16 bytes each)

---

## Architecture

### System Overview

```
┌─────────────────────────────────────────────────┐
│              Web Application                     │
│  (Development Environment)                       │
│                                                  │
│  ┌──────────┐ ┌──────────┐ ┌────────────────┐  │
│  │   Code   │ │  Sprite  │ │  Map Editor    │  │
│  │  Editor  │ │  Editor  │ │                │  │
│  └──────────┘ └──────────┘ └────────────────┘  │
│                                                  │
│  ┌──────────┐ ┌──────────┐ ┌────────────────┐  │
│  │  Sound   │ │  Music   │ │ Cart Manager   │  │
│  │  Editor  │ │  Tracker │ │                │  │
│  └──────────┘ └──────────┘ └────────────────┘  │
│                                                  │
│  ┌─────────────────────────────────────────┐   │
│  │         Live Preview / Testing          │   │
│  └─────────────────────────────────────────┘   │
└─────────────────────────────────────────────────┘
                       │
                       │ WASM Bridge
                       ▼
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
│  ┌──────────┐  ┌──────────┐  ┌──────────┐     │
│  │  Input   │  │  Sound   │  │ Collision│     │
│  │  System  │  │ Manager  │  │ Detection│     │
│  └──────────┘  └──────────┘  └──────────┘     │
│                                                  │
│  ┌─────────────────────────────────────────┐   │
│  │     Platform Layer (SDL2/OpenGL)        │   │
│  └─────────────────────────────────────────┘   │
└─────────────────────────────────────────────────┘
```

### Components

#### Go Runtime Engine

**Responsibilities:**
- Load and execute cart files
- Provide Lua API bindings
- Handle rendering (2D graphics primitives)
- Process audio synthesis and playback
- Manage input from keyboard/gamepad
- Persistent storage management

**Target Platforms:**
- Desktop (Windows, macOS, Linux)
- Web (WASM - future)
- Mobile (iOS, Android - future)

#### Web Application

**Responsibilities:**
- Code editing with syntax highlighting
- Visual sprite editor (multi-size support)
- Tilemap/map editor
- Sound effect designer
- Music pattern sequencer
- Cart export and management
- Live preview/testing interface

**Technology Stack:**
- Frontend Framework: React (with Vite)
- Code Editor: react-ace (Ace Editor for React)
- Canvas API for sprite/map editors
- Web Audio API for audio preview
- WASM for engine runtime in browser

---

## Cart Format

### File Structure

Carts use two file formats:
- **`.rfs`** (RetroForge Source) - Editable source format for development
- **`.rfe`** (RetroForge Executable) - Compressed, signed format for distribution

Both formats are **ZIP archives** with the same internal structure. The `.rfe` format includes compression and digital signatures.

```
game.rfe (compressed archive)
├── manifest.json       # Metadata, signature, and file manifest
├── main.lua           # Entry point (required)
├── code/              # Additional Lua modules (optional)
│   └── helpers.lua
├── levels/            # Map data (one JSON per level)
│   ├── level_1.json
│   ├── level_2.json
│   └── level_3.json
├── sprites.png        # Sprite sheet image
├── sprites.json       # Sprite definitions and metadata
└── audio.json         # Sound effects and music patterns
```

### File Format Details

**`.rfs` (RetroForge Source)**
- Uncompressed ZIP archive
- Used during development
- Easy to inspect and debug
- Can be version controlled (unzip to directory)
- Editable in the web app

**`.rfe` (RetroForge Executable)**
- Compressed ZIP archive (deflate)
- SHA-256 hash for integrity
- Optional RSA signature for authenticity
- Used for distribution
- Smaller file size
- Prevents casual tampering while remaining inspectable

### Compression
- **Format:** ZIP compression (standard deflate)
- **Benefits:** 
  - Smaller file sizes for distribution
  - Industry-standard format
  - Can be inspected with standard tools (for learning/debugging)
  - Fast decompression

### Digital Signature
- **Purpose:** Verify cart authenticity and integrity
- **Implementation:** 
  - SHA-256 hash of cart contents
  - Optional: RSA signature from cart author
  - Prevents tampering while allowing inspection of code
- **manifest.json includes:**
  ```json
  {
    "signature": {
      "hash": "sha256_hash_of_contents",
      "author_signature": "optional_rsa_signature",
      "timestamp": "2025-10-29T12:00:00Z"
    }
  }
  ```

### Metadata (manifest.json)

```json
{
  "version": "1.0",
  "title": "Game Title",
  "author": "Author Name",
  "description": "Game description",
  "created": "2025-10-29",
  "modified": "2025-10-29",
  "tags": ["action", "arcade"],
  "orientation": "landscape",
  "palette": "sunset",
  "signature": {
    "hash": "sha256_hash_of_contents",
    "author_signature": "optional_rsa_signature",
    "timestamp": "2025-10-29T12:00:00Z"
  }
}
```

### Map Data Format

Each level is stored as a **separate JSON file** in the `levels/` directory. This approach provides:
- ✅ Easy level management (add/remove levels without affecting others)
- ✅ Smaller file size per level (faster loading/editing)
- ✅ Clear organization
- ✅ Parallel loading capability

**Level File Structure (level_1.json):**
```json
{
  "name": "Level 1 - Forest",
  "width": 128,
  "height": 128,
  "tile_size": 16,
  "layers": [
    {
      "id": 0,
      "name": "UI",
      "parallax": 0,
      "visible": true,
      "data": [0, 0, 0, 1, 2, ...]  // Flat array, row-major order
    },
    {
      "id": 1,
      "name": "Foreground",
      "parallax": 1.0,
      "visible": true,
      "data": [...]
    },
    {
      "id": 2,
      "name": "Midground",
      "parallax": 0.7,
      "visible": true,
      "data": [...]
    },
    {
      "id": 7,
      "name": "Background",
      "parallax": 0.3,
      "visible": true,
      "data": [...]
    }
  ],
  "tile_flags": {
    "0": 0,    // Empty tile, no flags
    "1": 1,    // Solid ground (flag 0 set)
    "2": 3,    // Solid + hazard (flags 0 and 1 set)
    "3": 0     // Decoration, no collision
  },
  "spawn_points": [
    {"x": 64, "y": 64, "type": "player"},
    {"x": 100, "y": 50, "type": "enemy"}
  ]
}
```

**Data Storage:**
- Tiles stored as flat array (row-major order)
- Index calculation: `index = y * width + x`
- Empty tiles represented as `0`
- Sparse areas compress well with ZIP

### Color Palette Format

Palettes define the 50 colors available to the cart:

```json
{
  "name": "default",
  "colors": [
    "#000000",  // Black
    "#FFFFFF",  // White
    // Red family
    "#8B0000",  // Shadow
    "#FF0000",  // Base
    "#FF6B6B",  // Highlight
    // Blue family
    "#00008B",  // Shadow
    "#0000FF",  // Base
    "#6B6BFF",  // Highlight
    // ... 14 more color families (16 total)
  ]
}
```

### Predefined Palettes

The system includes multiple color theory-based palettes for users to choose from:

#### **Complementary Palettes** (Opposing colors on color wheel)
- **sunset**: Warm oranges/reds vs cool blues/purples
- **forest**: Greens vs magentas/reds

#### **Analogous Palettes** (Adjacent colors on color wheel)
- **ocean**: Blues, teals, and blue-greens
- **autumn**: Reds, oranges, and yellows

#### **Triadic Palettes** (Three evenly-spaced colors)
- **primary**: Red, yellow, blue based
- **vibrant**: High saturation triadic scheme

#### **Seasonal Palettes**
- **spring**: Pastels, light greens, pinks, soft yellows
- **summer**: Bright, saturated warm colors
- **fall**: Earth tones, oranges, browns, deep reds
- **winter**: Cool blues, whites, icy colors

#### **Mood-Based Palettes**
- **energetic**: High saturation, warm colors, strong contrasts
- **calm**: Low saturation, cool colors, gentle gradients
- **mysterious**: Dark purples, deep blues, blacks with accent colors
- **cheerful**: Bright yellows, oranges, light colors

#### **Specialty Palettes**
- **neon**: Bright cyberpunk colors with deep blacks
- **gameboy**: Green monochrome (4 shades × many hues)
- **grayscale**: Full grayscale range (black, white, 48 grays)
- **custom**: User-defined palette

Each palette maintains the structure: **Black + White + 16 colors × 3 shades = 50 colors**

---

## Lua API

RetroForge provides two API styles:

1. **Node API** (High-level, object-oriented) - Recommended for most games
2. **Direct API** (Low-level, PICO-8-style) - For simple carts or direct control

Both can be mixed in the same cart.

---

## Node System API

### Core Concepts

Every cart has a **Scene** that contains a tree of **Nodes**. Nodes automatically update and draw in tree order.

```lua
-- Basic cart structure with nodes
function _init()
  -- Create player
  player = RigidBody.new({
    position = vec2(240, 135),
    width = 16,
    height = 16
  })
  
  -- Add sprite to player
  local sprite = Sprite.new({sprite_index = 1})
  player:add_child(sprite)
  
  -- Add to scene
  Scene:add_child(player)
end

-- Nodes update automatically!
-- Override to add custom behavior
function player:_update()
  if Input.is_action_pressed("left") then
    self:apply_force(vec2(-100, 0))
  end
end
```

### Node (Base Class)

**Properties:**
```lua
node.name              -- string: node name
node.parent            -- Node: parent node (or nil)
node.children          -- array: child nodes
node.active            -- bool: is node active?
```

**Methods:**
```lua
node:add_child(child_node)           -- Add child node
node:remove_child(child_node)        -- Remove child node
node:get_child(index)                -- Get child by index
node:get_node(path)                  -- Get node by path ("../Player")
node:queue_free()                    -- Remove from tree next frame
node:is_in_tree()                    -- Check if added to scene
```

**Virtual Methods (Override in Lua):**
```lua
function node:_ready()               -- Called when added to tree
  -- Initialize node
end

function node:_update()              -- Called every frame
  -- Update logic
end

function node:_draw()                -- Called every frame (after update)
  -- Custom drawing
end

function node:_on_removed()          -- Called when removed from tree
  -- Cleanup
end
```

### Node2D (Spatial Node)

**Inherits:** Node  
**Use for:** Any node that needs position/rotation/scale

**Properties:**
```lua
node.position          -- vec2: position in pixels
node.rotation          -- number: rotation in radians
node.scale             -- vec2: scale (1.0 = normal)
node.z_index           -- number: draw order (higher = front)
node.global_position   -- vec2: world position (read-only)
```

**Methods:**
```lua
node:translate(offset)               -- Move by offset (vec2)
node:rotate(angle)                   -- Rotate by angle (radians)
node:look_at(target_pos)             -- Point toward position
node:get_global_transform()          -- Get world transform matrix
```

### Sprite

**Inherits:** Node2D  
**Use for:** Drawing sprites

**Properties:**
```lua
sprite.sprite_index    -- number: sprite index to draw
sprite.flip_h          -- bool: flip horizontally
sprite.flip_v          -- bool: flip vertically
sprite.modulate        -- color: tint color (default: white)
sprite.centered        -- bool: draw centered on position
sprite.offset          -- vec2: drawing offset
```

**Methods:**
```lua
sprite:set_sprite(index)             -- Set sprite index
```

### AnimatedSprite

**Inherits:** Sprite  
**Use for:** Sprite animations

**Properties:**
```lua
anim.animation         -- string: current animation name
anim.frame             -- number: current frame
anim.playing           -- bool: is playing?
anim.speed_scale       -- number: playback speed multiplier
```

**Methods:**
```lua
anim:play(animation_name, from_start)     -- Play animation
anim:stop()                                -- Stop animation
anim:set_animation(name)                   -- Set current animation
anim:add_animation(name, frames, fps)     -- Define animation
-- frames: array of sprite indices
-- fps: frames per second
```

**Signals:**
```lua
function anim:on_animation_finished()     -- Override
  -- Called when animation completes
end
```

### PhysicsBody2D (Base for physics nodes)

**Inherits:** Node2D  
**Use for:** Base class for physics objects

**Properties:**
```lua
body.collision_layer   -- number: which layer (0-7)
body.collision_mask    -- number: which layers to collide with
```

**Methods:**
```lua
body:add_collision_shape(shape)          -- Add collision shape
body:remove_collision_shape(shape)       -- Remove collision shape
body:get_colliding_bodies()              -- Array of bodies touching
```

### StaticBody

**Inherits:** PhysicsBody2D  
**Use for:** Walls, platforms (don't move)

Static bodies don't move and aren't affected by forces. Use for level geometry.

```lua
wall = StaticBody.new({
  position = vec2(100, 200),
  width = 64,
  height = 16
})
```

### RigidBody

**Inherits:** PhysicsBody2D  
**Use for:** Physics-simulated objects (affected by gravity and forces)

**Properties:**
```lua
body.velocity          -- vec2: current velocity (read-only)
body.mass              -- number: mass in kg
body.gravity_scale     -- number: gravity multiplier (1.0 = normal)
body.linear_damping    -- number: velocity damping (0-1)
body.angular_damping   -- number: rotation damping (0-1)
body.friction          -- number: surface friction (0-1)
body.bounce            -- number: bounciness/restitution (0-1)
body.can_sleep         -- bool: allow physics sleeping
body.fixed_rotation    -- bool: prevent rotation
```

**Methods:**
```lua
body:apply_force(force_vec2)             -- Apply continuous force
body:apply_impulse(impulse_vec2)         -- Apply instant impulse
body:apply_torque(torque)                -- Apply rotational force
body:set_velocity(velocity_vec2)         -- Set velocity directly
body:is_sleeping()                       -- Is body at rest?
```

**Collision Signals:**
```lua
function body:on_body_entered(other_body)
  -- Called when collision starts
end

function body:on_body_exited(other_body)
  -- Called when collision ends
end
```

### KinematicBody

**Inherits:** PhysicsBody2D  
**Use for:** Player-controlled characters (not affected by physics forces)

**Methods:**
```lua
body:move_and_collide(velocity_vec2)     -- Move with collision detection
body:move_and_slide(velocity_vec2)       -- Move and slide along surfaces
body:is_on_floor()                       -- Touching floor?
body:is_on_ceiling()                     -- Touching ceiling?
body:is_on_wall()                        -- Touching wall?
body:get_floor_normal()                  -- Floor surface normal
body:get_slide_collision_count()         -- Number of slide collisions
```

### CollisionShape

**Inherits:** Node  
**Use for:** Defining collision boundaries (child of PhysicsBody2D)

**Shape Types:**
```lua
-- Rectangle
shape = CollisionShape.rect(width, height)

-- Circle
shape = CollisionShape.circle(radius)

-- Capsule
shape = CollisionShape.capsule(width, height)

-- Polygon (up to 8 points)
shape = CollisionShape.polygon({
  vec2(0, 0),
  vec2(16, 0),
  vec2(16, 16),
  vec2(0, 16)
})
```

### Camera

**Inherits:** Node2D  
**Use for:** Viewport control and camera movement

**Properties:**
```lua
camera.zoom            -- number: zoom level (1.0 = normal, 2.0 = 2x)
camera.follow_target   -- Node2D: node to follow (or nil)
camera.follow_smooth   -- number: smoothing factor (0-1, 0=instant)
camera.limits          -- table: {left, top, right, bottom}
camera.drag_h          -- number: horizontal drag margin
camera.drag_v          -- number: vertical drag margin
```

**Methods:**
```lua
camera:follow(node)                      -- Start following node
camera:unfollow()                        -- Stop following
camera:shake(intensity, duration)        -- Camera shake effect
camera:screen_to_world(screen_pos)       -- Convert screen to world coords
camera:world_to_screen(world_pos)        -- Convert world to screen coords
```

### TileMap

**Inherits:** Node2D  
**Use for:** Rendering map layers

**Properties:**
```lua
tilemap.level_name     -- string: which level to load
tilemap.layer_id       -- number: which layer (0-7)
tilemap.parallax       -- number: parallax factor (1.0 = normal)
tilemap.tile_size      -- number: tile size (default: 16)
```

**Methods:**
```lua
tilemap:get_tile(x, y)                   -- Get tile at grid position
tilemap:set_tile(x, y, tile_id)          -- Set tile at grid position
tilemap:world_to_map(world_pos)          -- Convert world to tile coords
tilemap:map_to_world(map_pos)            -- Convert tile to world coords
```

### AudioPlayer (Sound Effects)

**Inherits:** Node  
**Use for:** Playing sound effects (multiple can play simultaneously)

**Properties:**
```lua
audio.sound_index      -- number: which SFX to play
audio.volume           -- number: volume (0.0-1.0)
audio.pitch            -- number: pitch multiplier (1.0 = normal)
audio.loop             -- bool: loop the sound?
audio.bus              -- string: audio bus ("sfx" or "master")
```

**Methods:**
```lua
audio:play(sound_index)                  -- Play sound (optional index)
audio:stop()                             -- Stop playing
audio:pause()                            -- Pause
audio:resume()                           -- Resume from pause
audio:is_playing()                       -- Check if playing
```

### MusicPlayer (Music Tracks)

**Inherits:** Node  
**Use for:** Playing music (only one track at a time, auto cross-fade)

**Properties:**
```lua
music.track_index      -- number: which music track
music.volume           -- number: volume (0.0-1.0)
music.bpm              -- number: tempo (read-only)
music.loop             -- bool: loop the track?
music.bus              -- string: audio bus ("music" or "master")
```

**Methods:**
```lua
music:play(track_index, crossfade_time)  -- Play track with fade
music:stop(fade_time)                    -- Stop with fade out
music:pause()                            -- Pause
music:resume()                           -- Resume
music:is_playing()                       -- Check if playing
music:seek(position)                     -- Seek to position (seconds)
music:get_position()                     -- Get current position
```

**Note:** Only one MusicPlayer can play at a time. Starting a new track automatically fades out the current one.

### SoundManager (Automatic Global)

**Global singleton:** Accessible as `Sound`  
**Created automatically:** Exists without any code  
**Use for:** Global audio control

**Properties:**
```lua
Sound.master_volume    -- number: master volume (0.0-1.0)
Sound.sfx_volume       -- number: SFX volume (0.0-1.0)
Sound.music_volume     -- number: music volume (0.0-1.0)
Sound.muted            -- bool: is all audio muted?
```

**Methods:**
```lua
Sound:play_sfx(sound_index, volume, pitch)   -- Quick SFX playback
Sound:play_music(track_index, crossfade)     -- Quick music playback
Sound:stop_music(fade_time)                  -- Stop current music
Sound:stop_all_sfx()                         -- Stop all sound effects
Sound:set_bus_volume(bus_name, volume)       -- Set bus volume
Sound:get_bus_volume(bus_name)               -- Get bus volume
Sound:enable_ducking(enable, amount)         -- Auto-lower music when SFX plays
```

**Audio Ducking:**
When enabled, music automatically lowers volume when sound effects play, then returns to normal.

```lua
-- Enable ducking: music drops to 30% when SFX plays
Sound:enable_ducking(true, 0.3)
```

### Timer

**Inherits:** Node  
**Use for:** Delayed or repeated actions

**Properties:**
```lua
timer.wait_time        -- number: duration in seconds
timer.one_shot         -- bool: trigger once (true) or repeat (false)
timer.autostart        -- bool: start automatically
timer.paused           -- bool: is timer paused?
timer.time_left        -- number: remaining time (read-only)
```

**Methods:**
```lua
timer:start(time)                        -- Start timer (optional time)
timer:stop()                             -- Stop timer
timer:pause()                            -- Pause
timer:resume()                           -- Resume
timer:is_stopped()                       -- Check if stopped
```

**Signals:**
```lua
function timer:on_timeout()              -- Override
  -- Called when timer completes
end
```

### ParticleEmitter

**Inherits:** Node2D  
**Use for:** Particle effects (explosions, smoke, sparkles)

**Properties:**
```lua
emitter.emitting       -- bool: is emitting?
emitter.amount         -- number: particles per emission
emitter.lifetime       -- number: particle lifetime (seconds)
emitter.one_shot       -- bool: emit once or continuous
emitter.explosiveness  -- number: 0=steady, 1=all at once
emitter.speed          -- number: particle speed
emitter.spread         -- number: angle spread (degrees)
emitter.gravity        -- vec2: gravity applied to particles
emitter.damping        -- number: velocity damping (0-1)
emitter.sprite_index   -- number: sprite to use for particles
emitter.color_start    -- color: initial particle color
emitter.color_end      -- color: final particle color
emitter.scale_start    -- number: initial scale
emitter.scale_end      -- number: final scale
```

**Methods:**
```lua
emitter:emit(count)                      -- Emit particles
emitter:start()                          -- Start emitting
emitter:stop()                           -- Stop emitting
emitter:restart()                        -- Stop and start
```

---

## Direct API (Low-Level)

These functions work without the node system, for PICO-8-style direct control.

### Lifecycle Functions

```lua
function _init()
  -- Called once at startup
  -- Initialize game state
end

function _update()
  -- Called every frame (60 FPS)
  -- Update game logic
end

function _draw()
  -- Called every frame after _update()
  -- Render graphics
end
```

### Graphics API

```lua
-- Screen Management
cls(color)                              -- Clear screen with color
camera(x, y)                            -- Set camera offset
clip(x, y, w, h)                        -- Set clipping region
clip()                                  -- Reset clipping

-- Drawing Primitives
pset(x, y, color)                       -- Set pixel
pget(x, y)                              -- Get pixel color
line(x1, y1, x2, y2, color)            -- Draw line
rect(x, y, w, h, color)                -- Draw rectangle outline
rectfill(x, y, w, h, color)            -- Draw filled rectangle
circ(x, y, r, color)                   -- Draw circle outline
circfill(x, y, r, color)               -- Draw filled circle

-- Sprites
spr(n, x, y, w, h, flip_x, flip_y)     -- Draw sprite
  -- n: sprite index
  -- x, y: screen position
  -- w, h: width/height in sprites (default 1, 1)
  -- flip_x, flip_y: boolean flip flags (optional)

-- Map
map(cell_x, cell_y, screen_x, screen_y, cells_w, cells_h, layer)
  -- Draw portion of map to screen

-- Text
print(text, x, y, color)               -- Draw text
cursor(x, y)                           -- Set cursor position for print

-- Palette
pal(c0, c1)                            -- Swap color c0 with c1
pal()                                  -- Reset palette
palt(c, transparent)                   -- Set color transparency
palt()                                 -- Reset transparency
```

### Input API

```lua
-- Button States
btn(b)                                 -- Get button state (continuous)
btnp(b, hold, repeat)                  -- Get button pressed (single frame)
  -- b: button index
  --   0=up, 1=down, 2=left, 3=right
  --   4=A, 5=B, 6=X, 7=Y

-- Mouse (future)
mouse()                                -- Get mouse state {x, y, btn}
```

### Audio API (Direct)

**Note:** For most games, use AudioPlayer and MusicPlayer nodes instead. These direct functions are for quick prototyping or simple carts.

```lua
-- Quick Sound Effects
sfx(n, channel, offset, length)        -- Play sound effect
  -- n: sound index (-1 to stop)
  -- channel: audio channel (0-7, -1 for auto)
  -- offset: start position in pattern
  -- length: number of notes to play

-- Quick Music
music(pattern, fade_ms, channel_mask)  -- Play music pattern
  -- pattern: pattern index (-1 to stop)
  -- fade_ms: fade in/out time
  -- channel_mask: which channels to use

-- Global controls (shortcuts to SoundManager)
set_volume(volume)                     -- Set master volume
set_sfx_volume(volume)                 -- Set SFX volume
set_music_volume(volume)               -- Set music volume
```

**Recommended approach:**
```lua
-- Instead of direct API:
sfx(5, 0)

-- Use nodes:
local player = AudioPlayer.new({sound_index = 5})
Scene:add_child(player)
player:play()

-- Or use SoundManager:
Sound:play_sfx(5)
```

### Input API

**Button Constants:**
```lua
BTN_UP = 0
BTN_DOWN = 1
BTN_LEFT = 2
BTN_RIGHT = 3
BTN_A = 4        -- Action button
BTN_B = 5        -- Jump/Back button
BTN_X = 6        -- Secondary action
BTN_Y = 7        -- Secondary action
```

**Input Object (Global):**
```lua
-- Button states
Input.is_action_pressed("action_name")   -- Is button held down?
Input.is_action_just_pressed("action")   -- Was button just pressed?
Input.is_action_just_released("action")  -- Was button just released?

-- Direct button access
Input.get_button(button_id)              -- Get button state (bool)
Input.get_button_pressed(button_id)      -- Was pressed this frame?

-- Action mapping (define in cart)
Input.add_action("jump", BTN_A)          -- Map action to button
Input.add_action("shoot", BTN_X)

-- Mouse (future)
Input.get_mouse_position()               -- vec2: mouse position
Input.is_mouse_button_pressed(button)    -- Mouse button state
```

**Legacy functions (PICO-8 style):**
```lua
btn(b)                                 -- Button state (continuous)
btnp(b, hold, repeat)                  -- Button pressed (single frame)
  -- b: button index (0-7)
  -- hold: frames to hold before repeat (optional)
  -- repeat: repeat interval in frames (optional)
```

### Physics API

**Raycasting:**
```lua
-- Cast a ray and get first hit
hit = Physics.raycast(origin_vec2, direction_vec2, distance)
-- Returns: {body, point, normal, distance} or nil

-- Cast ray and get all hits
hits = Physics.raycast_all(origin, direction, distance)
-- Returns: array of hit results

-- Check if point is in body
body = Physics.point_query(point_vec2)
```

**Collision queries:**
```lua
-- Get all bodies in area
bodies = Physics.area_query(center, radius)

-- Get all bodies in rectangle
bodies = Physics.rect_query(x, y, width, height)
```

**World settings:**
```lua
Physics.set_gravity(x, y)              -- Set world gravity
Physics.get_gravity()                  -- Get world gravity (vec2)
Physics.set_timestep(step)             -- Physics timestep (default: 1/60)
```

### Vector Math (vec2)

**Constructor:**
```lua
v = vec2(x, y)                         -- Create vector
v = vec2.zero()                        -- (0, 0)
v = vec2.one()                         -- (1, 1)
v = vec2.up()                          -- (0, -1)
v = vec2.down()                        -- (0, 1)
v = vec2.left()                        -- (-1, 0)
v = vec2.right()                       -- (1, 0)
```

**Operations:**
```lua
v1 + v2                                -- Add vectors
v1 - v2                                -- Subtract
v * scalar                             -- Multiply by scalar
v / scalar                             -- Divide by scalar
-v                                     -- Negate

v:length()                             -- Get magnitude
v:normalized()                         -- Get unit vector
v:dot(other)                           -- Dot product
v:distance_to(other)                   -- Distance to other vector
v:angle()                              -- Angle in radians
v:rotated(angle)                       -- Rotate by angle
v:lerp(target, weight)                 -- Linear interpolation
```

**Properties:**
```lua
v.x                                    -- X component
v.y                                    -- Y component
```

### Scene API

**Global Scene object:**
```lua
Scene:add_child(node)                  -- Add node to root
Scene:remove_child(node)               -- Remove from root
Scene:get_node(path)                   -- Get node by path
Scene:find_nodes_by_name(name)         -- Find all nodes with name
Scene:clear()                          -- Remove all nodes
```

**Scene tree queries:**
```lua
Scene:get_nodes_in_group(group)        -- Get all nodes in group
Scene:call_group(group, method, ...)   -- Call method on group
```

**Node groups:**
```lua
node:add_to_group(group_name)          -- Add node to group
node:remove_from_group(group_name)     -- Remove from group
node:is_in_group(group_name)           -- Check if in group
```

### Audio Data Format

Sound effects and music are defined in JSON format with the following structure:

```json
{
  "sfx": [
    {
      "id": 0,
      "name": "jump",
      "waveform": "square",
      "speed": 16,
      "notes": [
        {"pitch": 48, "volume": 15, "effect": "none"},
        {"pitch": 52, "volume": 12, "effect": "none"},
        {"pitch": 55, "volume": 8, "effect": "none"}
      ],
      "envelope": {
        "attack": 0.01,
        "decay": 0.1,
        "sustain": 0.6,
        "release": 0.2
      },
      "vibrato": {
        "enabled": false,
        "speed": 4,
        "depth": 0.05
      }
    }
  ],
  "music": [
    {
      "id": 0,
      "name": "level_1",
      "bpm": 120,
      "patterns": [
        {
          "channel": 0,
          "sfx_sequence": [0, 1, 2, 0]
        }
      ]
    }
  ]
}
```

### Map API

```lua
-- Map Access
mget(x, y, layer)                      -- Get map tile at position
mset(x, y, tile, layer)                -- Set map tile at position

-- Tile Flags
fget(tile, flag)                       -- Get flag value for tile
fset(tile, flag, value)                -- Set flag value for tile
  -- flag: 0-7
  -- value: boolean (true/false)

-- Parallax
camera(x, y, layer)                    -- Set camera position for layer
  -- layer: 0-7 (omit for all layers)
```

### Memory/Storage API

```lua
-- Persistent Data
cartdata(id)                           -- Set cart data ID (call in _init)
dget(index)                            -- Get persistent value (0-63)
dset(index, value)                     -- Set persistent value
```

### Utility API

```lua
-- Math
abs(x)                                 -- Absolute value
min(x, y)                              -- Minimum
max(x, y)                              -- Maximum
mid(x, y, z)                           -- Middle value
flr(x)                                 -- Floor
ceil(x)                                -- Ceiling
sgn(x)                                 -- Sign (-1, 0, 1)
sqrt(x)                                -- Square root
sin(x)                                 -- Sine (0-1 input = full cycle)
cos(x)                                 -- Cosine (0-1 input = full cycle)
atan2(dx, dy)                          -- Arctangent

-- Random
rnd(max)                               -- Random number [0, max)
srand(seed)                            -- Set random seed

-- Time
time()                                 -- Seconds since cart start

-- Strings
sub(str, start, end)                   -- Substring
#str                                   -- String length
```

### Table API

```lua
-- Standard Lua table operations available
add(table, value)                      -- Add to end of array
del(table, value)                      -- Remove value from array
count(table)                           -- Count elements
all(table)                             -- Iterator for arrays
foreach(table, func)                   -- Call function for each element
```

---

## Development Roadmap

### Phase 1: Core Engine (MVP)
- [ ] Go project setup with gopher-lua
- [ ] Basic window and rendering loop (SDL2)
- [ ] Implement core graphics API (primitives, sprites)
- [ ] Basic input handling
- [ ] Cart loader (JSON format)
- [ ] Lifecycle functions (_init, _update, _draw)

### Phase 2: Complete Engine
- [ ] Audio system implementation
- [ ] Map system
- [ ] Persistent storage
- [ ] Complete Lua API
- [ ] Performance optimization

### Phase 3: Web Application
- [ ] Project setup and architecture
- [ ] Code editor integration
- [ ] Sprite editor
- [ ] Map editor
- [ ] Sound/music tools
- [ ] Cart export/import

### Phase 4: Polish & Distribution
- [ ] Example carts and documentation
- [ ] Tutorial system
- [ ] Cart sharing platform
- [ ] Cross-platform builds
- [ ] Performance profiling tools

---

## Design Decisions Made

### ✅ Core Architecture

**Node System:** Godot-style scene graph
- Hierarchical node tree with parent/child relationships
- Automatic update and draw ordering
- Transform propagation through tree
- Component-based via node composition
- Both high-level (nodes) and low-level (direct) APIs available

**Physics Engine:** Box2D-go
- Industry-standard physics (Box2D port)
- Rigid body dynamics with realistic simulation
- Collision detection and response
- Joints, constraints, and raycasting
- Three body types: Static, Dynamic, Kinematic
- 8 collision layers with mask filtering

### ✅ Audio Architecture

**Three-tier system:**

1. **AudioPlayer** (Node) - Sound effects
   - Multiple instances play simultaneously
   - Positioned in scene graph
   - Per-instance volume/pitch control
   
2. **MusicPlayer** (Node) - Music tracks
   - Only one track plays at a time
   - Automatic cross-fading between tracks
   - Tempo and loop control
   
3. **SoundManager** (Automatic Global) - Master control
   - Created automatically in every cart
   - No code needed to use
   - Global volume controls (master, sfx, music)
   - Audio ducking (lower music when SFX plays)
   - Channel allocation management
   - Accessible via `Sound` global

### ✅ Color Palette
- **System:** Predefined palettes based on color theory + custom palette option
- **Structure:** Black + White + 16 colors with 3 shades each (base, highlight, shadow)
- **Total colors per palette:** 50 colors
- **Predefined palettes:** 16+ palettes across complementary, analogous, triadic, seasonal, mood-based, and specialty categories
- **See:** "Predefined Palettes" section for complete list

### ✅ Map System
- **Tile Size:** 16×16 pixels (fixed)
- **Layers:** 8 layers with parallax support
  - Layer 0: UI layer (no parallax)
  - Layers 1-7: Game world (parallax enabled, 1=near, 7=far)
- **Collision:** Tile flags system (8 flags per tile)
  - Flag 0 reserved for collision detection
  - Flags 1-7 user-definable
- **Data Format:** Separate JSON file per level
  - Benefits: Easy level management, smaller files, faster loading
  - Location: `levels/` directory in cart

### ✅ Audio System
- **Pure Synthesis** (no samples)
- **Waveforms:** Square, Triangle, Sawtooth, Sine, Noise
- **Envelope:** ADSR (Attack, Decay, Sustain, Release)
- **Effects:** Vibrato, Arpeggio, Pitch slide
- **Editor UI:** Visual piano-roll interface
  - **Why:** Most intuitive for non-programmers
  - **Features:** Click to place notes, drag to adjust length, visual waveform preview
  - **Inspiration:** GarageBand, FL Studio piano roll (simplified)
  - **Alternative view:** Simple pattern grid for quick chip-tune creation

### ✅ Web App ↔ Engine Communication
- **WASM engine** compiled from Go using TinyGo
- Runs directly in browser with JavaScript interop
- No separate server needed for development
- Direct JavaScript ↔ WASM communication via Web APIs

### ✅ WASM Performance Targets
- **Frame Rate:** Stable 60 FPS (16.67ms per frame)
- **Frame Budget:**
  - Update logic: 8ms max
  - Rendering: 6ms max
  - Audio: 2ms max
- **Memory Limit:** 64MB WASM heap (generous for constrained games)
- **Audio Latency:** <20ms for responsive sound
- **Load Time:** <2 seconds for average cart (before decompression)
- **Optimization Strategy:**
  - Pre-allocate common objects
  - Batch rendering calls
  - Use TypedArrays for pixel/audio buffers
  - Minimize GC pressure

### ✅ Build & Distribution
- **Engine Builds:**
  - Windows: Single .exe executable
  - macOS: .app bundle
  - Linux: Single executable binary
  - Android: .apk package (via gomobile)
  - Web: WASM module for browsers
- **iOS:** Not supported (expensive developer account required, restrictive policies)
- **Cart Formats:** 
  - `.rfs` (RetroForge Source) - Editable development format
  - `.rfe` (RetroForge Executable) - Compressed distribution format
- **Compression:** ZIP format with deflate algorithm
- **Security:** SHA-256 hash + optional RSA signature for verification
- **Package Manager:** Built into web application
- **Repository:** Online cart repository integrated with web app and editors

### ✅ Technology Stack
- **Frontend:** Next.js 16, TypeScript, TailwindCSS
- **Backend/Engine:** Go 1.23, gopher-lua
- **Compilation:** TinyGo for WASM, standard Go for native builds
- **Code Editor:** react-ace with Lua syntax highlighting

---

## Open Questions

### 1. UI Component Library
- **Confirmed:** TailwindCSS for styling
- **Need to decide:** Component library
  - **shadcn/ui** (Radix + Tailwind) - Modern, accessible, customizable ✅ Recommended
  - **Headless UI** - Minimal, fully unstyled
  - **Custom components only** - Maximum control, more work

### 2. Sprite Editor Features
### 2. Sprite Editor Features
What functionality should the sprite editor include?
- Basic drawing tools (pencil, fill, line, rectangle, circle)
- Animation preview and timeline
- Onion skinning for animation
- Import from external tools (PNG import)
- Copy/paste between sprites
- Sprite flipping/rotating tools
- Color palette selector

### 3. Distribution Platform Details
- **Hosting:** Self-hosted vs cloud service (Vercel, Netlify)?
- **Authentication:** Email/password, OAuth (GitHub, Google)?
- **Social features:** Likes, comments, follows, user profiles?
- **Monetization:** Free only, or optional creator tips/donations?
- **Moderation:** Community reporting, automated filters?

### 4. Documentation Strategy
How should we teach users to use the console?
- Interactive tutorials built into web app?
- Separate documentation site (like docs.retroforge.dev)?
- Video tutorial series?
- Example cart library with commented code?
- API reference with live examples?
- Community wiki?

---

## Technology Choices

### Go Engine
- **Language:** Go 1.23
- **Lua VM:** gopher-lua (pure Go, Lua 5.1 implementation)
- **Physics:** box2d-go (port of Box2D physics engine)
- **Graphics:** SDL2 + OpenGL (desktop), WebGL (WASM), OpenGL ES (Android)
- **Audio:** oto or beep library (desktop), Web Audio API (WASM), OpenSL ES (Android)
- **Build Tools:**
  - Standard Go toolchain for native builds
  - TinyGo for WASM compilation
  - gomobile for Android APK builds
- **Key Dependencies:**
  - `github.com/yuin/gopher-lua` - Lua VM
  - `github.com/ByteArena/box2d` - Physics engine
  - `github.com/veandco/go-sdl2` - SDL2 bindings (desktop)
  - `github.com/hajimehoshi/oto` - Audio playback
  - `golang.org/x/mobile` - Android/mobile support
- **Target Platforms:**
  - Desktop: Windows, macOS, Linux (single executables)
  - Mobile: Android (APK)
  - Web: WASM module for browser-based runtime
  - iOS: Not supported (development restrictions and costs)

### Web Application
- **Framework:** Next.js 16 (React-based)
- **Language:** TypeScript
- **Styling:** TailwindCSS
- **Code Editor:** react-ace (Ace Editor for React)
  - Lua syntax highlighting
  - Autocomplete for RetroForge API
  - Custom themes
  - Line numbers, bracket matching
- **UI Components:** TBD (shadcn/ui recommended)
- **Canvas Editors:** HTML5 Canvas API for sprite/map editors
- **Audio:** Web Audio API for preview and WASM audio output
- **State Management:** React Context / Zustand / Jotai (TBD)
- **Package Manager:** Custom cart repository with versioning
- **Build Tool:** Next.js built-in (Turbopack)

---

## References

- PICO-8: https://www.lexaloffle.com/pico-8.php
- gopher-lua: https://github.com/yuin/gopher-lua
- Fantasy Console Wiki: https://github.com/paladin-t/fantasy
- TIC-80: https://tic80.com/

---

## Notes

This is a living document. As development progresses, sections will be updated with implementation details, decisions, and lessons learned.
