# RetroForge Engine Developer Guide

**RetroForge Version:** v1.0 Alpha  
**Guide Version:** 1.0  
**Last Updated:** October 31, 2025  
**Complete Guide for Building Retro-Style Games**

---

## Table of Contents

1. [What is RetroForge?](#what-is-retroforge)
2. [Core Concepts & Architecture](#core-concepts--architecture)
3. [Getting Started](#getting-started)
4. [Game Development Fundamentals](#game-development-fundamentals)
5. [API Reference](#api-reference)
6. [Advanced Features](#advanced-features)
7. [Best Practices & Patterns](#best-practices--patterns)
8. [Examples & Tutorials](#examples--tutorials)
9. [Troubleshooting](#troubleshooting)

---

## What is RetroForge?

### Overview

**RetroForge** is a modern fantasy console (inspired by PICO-8) that enables developers to create retro-style 2D games with modern development tools. It combines the charm and constraints of retro game development with the power of a full-featured engine.

### Key Characteristics

- **Fantasy Console**: Self-contained runtime with fixed capabilities (like PICO-8, TIC-80)
- **Modern Engine**: Built in Go with Box2D physics, WebRTC multiplayer, and cross-platform deployment
- **Lua Scripting**: Write games in Lua, a simple and powerful scripting language
- **Cart-Based**: Games are packaged as `.rf` cart files (ZIP archives)
- **Multiplayer-Ready**: Built-in WebRTC networking for up to 6 players

### Philosophy

RetroForge follows the **"constrained creativity"** philosophy of fantasy consoles:

1. **Fixed Constraints**: 480×270 resolution, 50-color palette, fixed memory limits
2. **Clear Capabilities**: Well-defined API with clear boundaries
3. **Creative Solutions**: Constraints inspire creative problem-solving
4. **Rapid Iteration**: Simple tools enable quick prototyping and iteration
5. **Community Focus**: Shareable cart files (.rf) work everywhere

### Why RetroForge?

**Compared to PICO-8:**
- ✅ Higher resolution (480×270 vs 128×128)
- ✅ More colors (50 vs 16)
- ✅ Built-in multiplayer (WebRTC vs manual networking)
- ✅ Modern development tools (hot reload, debugging)
- ✅ Cross-platform (Desktop, Web, Android)

**Compared to Full Game Engines:**
- ✅ Simpler API (fewer concepts to learn)
- ✅ Faster iteration (no complex build processes)
- ✅ Smaller file sizes (carts are typically 10-100KB)
- ✅ Built-in limitations (forces creativity)
- ✅ Retro aesthetic (pixel art, chip-tune audio)

---

## Core Concepts & Architecture

### How RetroForge Works

RetroForge operates on a **game loop** pattern:

```
┌─────────────────────────────────────┐
│  Engine Initialization               │
│  - Load cart (.rf file)              │
│  - Parse manifest.json               │
│  - Load assets (sprites, audio)      │
│  - Initialize Lua VM                 │
│  - Register state machine            │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  Game Loop (60 FPS)                 │
│  ┌──────────────────────────────┐  │
│  │ 1. Physics Step                │  │
│  │    - Physics.Step()           │  │
│  │    - Network frame update     │  │
│  └──────────────┬────────────────┘  │
│  ┌──────────────▼────────────────┐  │
│  │ 2. State Machine              │  │
│  │    - HandleInput()            │  │
│  │    - Update(dt)               │  │
│  │    - Draw()                    │  │
│  └──────────────┬────────────────┘  │
│  ┌──────────────▼────────────────┐  │
│  │ 3. Input State Step            │  │
│  │    - input.Step()             │  │
│  │    - Saves current frame for   │  │
│  │      next frame's btnp()       │  │
│  └──────────────────────────────┘  │
└─────────────────────────────────────┘
```

### Core Components

#### 1. **Lua Virtual Machine**
- Runs your game code
- Provides `rf.*` API for graphics, input, audio, etc.
- Provides `game.*` API for state management
- Isolated environment per cart (no global pollution)

#### 2. **State Machine**
- Manages game flow (menu → playing → game over)
- Supports state stacking (pause overlays)
- Built-in splash and credits screens
- Module-based system for organization

#### 3. **Renderer**
- 480×270 pixel buffer
- 50-color palette system
- Drawing primitives (lines, circles, rectangles)
- Sprite system with automatic pooling
- Camera and clipping support

#### 4. **Physics Engine (Box2D)**
- Rigid body physics
- Collision detection
- Forces and impulses
- Static, dynamic, and kinematic bodies

#### 5. **Audio System**
- 8 audio channels
- Chip-tune synthesis (5 waveforms)
- JSON-based sound effects and music
- Pattern-based music system

#### 6. **Multiplayer (WebRTC)**
- Up to 6 players
- Automatic state synchronization
- Host-authoritative architecture
- 3-tier sync system (fast/moderate/slow)

---

## Getting Started

### Prerequisites

- **Go 1.23+** (for engine development)
- **Text editor** (VS Code, Sublime, etc.)
- **Basic Lua knowledge** (helpful but not required)

### Installation

#### For Game Development (Web)

No installation needed! Just:
1. Visit the RetroForge webapp
2. Use the online editor
3. Export your cart (.rf file)

#### For Engine Development (Desktop)

```bash
git clone https://github.com/AndrewDonelson/retroforge-engine.git
cd retroforge-engine
go mod download
go build -o retroforge ./cmd/retroforge
```

### Your First Cart

Create a folder structure:

```
my-game/
├── manifest.json
└── assets/
    └── main.lua
```

**manifest.json:**
```json
{
  "title": "My First Game",
  "author": "Your Name",
  "description": "A simple game",
  "entry": "main.lua",
  "palette": "PICO-8"
}
```

**assets/main.lua:**
```lua
-- Simple hello world
function _INIT()
  -- Initialize (called once at start)
end

function _UPDATE(dt)
  -- Update game logic (dt = delta time in seconds)
end

function _DRAW()
  -- Draw everything
  rf.clear_i(0)  -- Clear to black (color 0)
  rf.print_anchored("HELLO WORLD!", "middlecenter", 15)  -- White text
end
```

### Running Your Cart

**Development Mode (Hot Reload):**
```bash
make run-dev FOLDER=my-game
```

**Pack as Cart:**
```bash
./retroforge -pack my-game
# Creates my-game.rf
```

**Run Packed Cart:**
```bash
./retroforge -cart my-game.rf -window
```

---

## Game Development Fundamentals

### The Game Loop

Every RetroForge game has three main functions:

#### `_INIT()`
Called **once** when the game starts. Use for:
- Initializing variables
- Loading resources
- Setting up game state

```lua
function _INIT()
  player = {x = 240, y = 135, vx = 0, vy = 0}
  score = 0
  level = 1
end
```

#### `_UPDATE(dt)`
Called **every frame** (60 times per second). Use for:
- Processing input
- Updating game logic
- Physics calculations
- AI behavior

```lua
function _UPDATE(dt)
  -- Move player
  if rf.btn(0) then  -- Left
    player.vx = -3
  elseif rf.btn(1) then  -- Right
    player.vx = 3
  else
    player.vx = 0
  end
  
  player.x = player.x + player.vx * dt * 60  -- Scale by dt for frame-independent movement
end
```

#### `_DRAW()`
Called **every frame** after `_UPDATE`. Use for:
- Clearing the screen
- Drawing sprites
- Drawing primitives
- Drawing UI

```lua
function _DRAW()
  rf.clear_i(0)  -- Always clear first!
  
  -- Draw player
  rf.circfill(player.x, player.y, 8, 7)  -- White circle
  
  -- Draw UI
  rf.print_xy(10, 10, "Score: " .. score, 15)
end
```

### Coordinate System

- **Origin**: Top-left corner is (0, 0)
- **X-axis**: Increases rightward (0 to 479)
- **Y-axis**: Increases downward (0 to 269)
- **Resolution**: 480×270 pixels

```
(0,0) ──────────────── (479,0)
  │                         │
  │                         │
  │         Screen          │
  │                         │
  │                         │
(0,269) ────────────── (479,269)
```

### Color System

RetroForge uses a **palette system** with 50 colors (indices 0-49):

- **Palette indices**: Use numbers (0-49) instead of RGB values
- **Transparent**: Use `-1` for transparent pixels
- **Palette selection**: Choose from predefined palettes (PICO-8, SNES, etc.)
- **Color remapping**: Use `rf.pal(c0, c1)` to remap colors

```lua
rf.palette_set("PICO-8")  -- Use PICO-8 palette
rf.rectfill(100, 100, 150, 150, 7)  -- Color 7 from palette
```

### Input System

8 standard buttons:

```
Buttons:
  0 = Left
  1 = Right
  2 = Up
  3 = Down
  4 = Z (or O)
  5 = X (or Cancel)
  6-15 = Extended buttons
```

**Functions:**
- `rf.btn(button)` - Check if button is held
- `rf.btnp(button)` - Check if button was just pressed (edge-triggered)
- `rf.btnr(button)` - Check if button was just released

```lua
function _UPDATE(dt)
  -- Movement (continuous)
  if rf.btn(0) then  -- Left held
    player.x = player.x - 3
  end
  
  -- Jump (one-time)
  if rf.btnp(4) then  -- Z just pressed
    player.vy = -10
    rf.sfx("jump")
  end
end
```

---

## API Reference

### Graphics API

#### Screen Operations

**`rf.clear_i(index)`**
Clear the entire screen to a palette color.

```lua
rf.clear_i(0)  -- Clear to black
```

#### Drawing Primitives

**`rf.pset(x, y, index)`** - Draw a single pixel
```lua
rf.pset(100, 100, 7)  -- White pixel
```

**`rf.line(x0, y0, x1, y1, index)`** - Draw a line
```lua
rf.line(0, 0, 479, 269, 7)  -- Diagonal line
```

**`rf.rect(x0, y0, x1, y1, index)`** - Draw rectangle outline
```lua
rf.rect(10, 10, 100, 50, 7)
```

**`rf.rectfill(x0, y0, x1, y1, index)`** - Draw filled rectangle
```lua
rf.rectfill(10, 10, 100, 50, 7)
```

**`rf.circ(x, y, radius, index)`** - Draw circle outline
```lua
rf.circ(240, 135, 50, 7)
```

**`rf.circfill(x, y, radius, index)`** - Draw filled circle
```lua
rf.circfill(240, 135, 50, 7)
```

**`rf.elli(x, y, rx, ry, index)`** - Draw ellipse outline
```lua
rf.elli(240, 135, 60, 30, 7)
```

**`rf.ellifill(x, y, rx, ry, index)`** - Draw filled ellipse
```lua
rf.ellifill(240, 135, 60, 30, 7)
```

#### Text Drawing

**`rf.print(text, [x, y, index])`** - Print text at position
```lua
rf.print("Hello", 10, 10, 7)  -- Positioned
rf.print("World", 7)  -- Uses cursor position
```

**`rf.print_xy(x, y, text, [index])`** - Print at exact position
```lua
rf.print_xy(100, 100, "Score: " .. score, 15)
```

**`rf.print_anchored(text, anchor, index)`** - Print using anchor
```lua
rf.print_anchored("GAME OVER", "middlecenter", 15)
-- Anchors: topleft, topcenter, topright, middleleft, middlecenter,
--          middleright, bottomleft, bottomcenter, bottomright
```

#### Sprites

**`rf.spr(name, x, y, [flip_x, flip_y])`** - Draw sprite
```lua
rf.spr("player", player.x, player.y)
rf.spr("enemy", enemy.x, enemy.y, true, false)  -- Flipped horizontally
```

**`rf.sspr(name, sx, sy, sw, sh, dx, dy, [dw, dh, flip_x, flip_y])`** - Draw sprite region
```lua
rf.sspr("tiles", 0, 0, 16, 16, 100, 100)  -- Draw 16×16 region at (100,100)
```

#### Camera & Clipping

**`rf.camera([x, y])`** - Set camera offset
```lua
rf.camera(player.x - 240, player.y - 135)  -- Follow player
rf.camera()  -- Reset to (0, 0)
```

**`rf.clip([x, y, w, h])`** - Set clipping rectangle
```lua
rf.clip(100, 100, 200, 100)  -- Restrict drawing to rectangle
rf.clip()  -- Disable clipping
```

### Sprite Creation & Editing

#### Creating Sprites

**`rf.newSprite(name, width, height)`** - Create empty sprite
```lua
local bullet = rf.newSprite("bullet", 8, 8)
```

#### Drawing to Sprites

**`rf.sprite_pset(sprite_name, x, y, index)`** - Set pixel in sprite
```lua
rf.sprite_pset("bullet", 4, 4, 11)  -- Yellow center
```

**`rf.sprite_circfill(sprite_name, x, y, radius, index)`** - Draw circle in sprite
```lua
rf.sprite_circfill("bullet", 4, 4, 3, 11)
```

**`rf.sprite_rectfill(sprite_name, x0, y0, x1, y1, index)`** - Draw rectangle in sprite
```lua
rf.sprite_rectfill("enemy", 0, 0, 15, 15, 8)
```

#### Sprite Properties

**`rf.setSpriteProperty(sprite_name, property, value)`** - Set sprite property
```lua
rf.setSpriteProperty("bullet", "isUI", false)  -- Physics-enabled
rf.setSpriteProperty("bullet", "useCollision", true)  -- Enable collision
rf.setSpriteProperty("bullet", "lifetime", 2000)  -- 2 second lifetime
rf.setSpriteProperty("bullet", "maxSpawn", 50)  -- Max 50 at once
```

**Properties:**
- `isUI` (boolean): If `true`, sprite is UI (no physics). Default: `true`
- `useCollision` (boolean): Enable physics collision. Default: `false`
- `lifetime` (number): Auto-destroy after milliseconds. Default: `0` (no limit)
- `maxSpawn` (number): Maximum simultaneous instances. Default: `0` (no limit)

**Automatic Sprite Pooling:**
- Sprites with `isUI=false` AND `maxSpawn > 10` are automatically pooled
- Improves performance by reusing instances
- Transparent to developer (no code changes needed)

### Physics API

#### Creating Bodies

**`rf.physics_create_body(type, x, y)`** - Create physics body
```lua
local body = rf.physics_create_body("dynamic", 240, 135)
-- Types: "static", "dynamic", "kinematic"
```

#### Adding Fixtures

**`rf.physics_body_add_box(body_id, width, height, [density, restitution, friction])`**
```lua
rf.physics_body_add_box(body, 16, 16, 1.0, 0.5, 0.3)
-- density: 1.0 = normal, restitution: 0.5 = bouncy, friction: 0.3 = medium
```

**`rf.physics_body_add_circle(body_id, radius, [density, restitution, friction])`**
```lua
rf.physics_body_add_circle(body, 8, 1.0, 0.8, 0.1)
```

#### Controlling Bodies

**`rf.physics_body_set_position(body_id, x, y)`** - Set position
```lua
rf.physics_body_set_position(player.body, 100, 100)
```

**`rf.physics_body_get_position(body_id)`** - Get position (returns x, y)
```lua
local x, y = rf.physics_body_get_position(player.body)
```

**`rf.physics_body_set_velocity(body_id, vx, vy)`** - Set velocity
```lua
rf.physics_body_set_velocity(player.body, 5, -10)
```

**`rf.physics_body_get_velocity(body_id)`** - Get velocity (returns vx, vy)
```lua
local vx, vy = rf.physics_body_get_velocity(player.body)
```

**`rf.physics_body_set_gravity_scale(body_id, scale)`** - Set gravity scale
```lua
rf.physics_body_set_gravity_scale(body, 0.0)  -- No gravity
```

**`rf.physics_body_apply_force(body_id, fx, fy, px, py)`** - Apply force
```lua
rf.physics_body_apply_force(body, 10, 0, 0, 0)  -- Push right
```

**`rf.physics_body_destroy(body_id)`** - Destroy body
```lua
rf.physics_body_destroy(enemy.body)
```

### Audio API

#### Sound Effects

**`rf.sfx(name, ...)`** - Play sound effect
```lua
rf.sfx("jump")  -- Named SFX from sfx.json
rf.sfx("sine", 440, 0.1, 0.5)  -- Inline: sine wave, 440Hz, 0.1s, 50% volume
```

**Waveforms:** `"sine"`, `"square"`, `"triangle"`, `"sawtooth"`, `"noise"`

#### Music

**`rf.music(name, ...)`** - Play music track
```lua
rf.music("theme")  -- Named track from music.json
rf.music({"4C1", "4E1", "4G1", "R1"}, 120, 0.7)  -- Inline pattern
```

### State Machine API

#### Module-Based States (Recommended)

**`rf.import(filename)`** - Import state module
```lua
local menu_state = rf.import("menu_state.lua")  -- Auto-registers as "menu"
local play_state = rf.import("play_state.lua")  -- Auto-registers as "play"
```

**State Module Functions:**
- `_INIT()` - Called once (required)
- `_UPDATE(dt)` - Called every frame (required)
- `_DRAW()` - Called every frame (required)
- `_HANDLE_INPUT()` - Called every frame (required)
- `_DONE()` - Called once on destroy (required)
- `_ENTER()` - Called when state becomes active (optional)
- `_EXIT()` - Called when state becomes inactive (optional)

#### State Transitions

**`game.changeState(name)`** - Replace all states with new state
```lua
game.changeState("playing")  -- Menu → Playing
```

**`game.pushState(name)`** - Add state on top (overlay)
```lua
game.pushState("pause")  -- Playing → Pause (playing stays in memory)
```

**`game.popState()`** - Remove top state
```lua
game.popState()  -- Close pause menu, return to playing
```

**`game.exit()`** - Exit game (shows credits, then quits)
```lua
game.exit()
```

#### Context API

**`game.setContext(key, value)`** - Store data in shared context
```lua
game.setContext("level", 5)
game.setContext("difficulty", "hard")
```

**`game.getContext(key)`** - Get data from context
```lua
local level = game.getContext("level") or 1
```

#### Credits API

**`game.addCredit(category, name, role)`** - Add credit entry
```lua
game.addCredit("Developer", "Jane Doe", "Lead Programmer")
game.addCredit("Artist", "John Smith", "Pixel Artist")
```

### Memory API

#### Runtime Memory (2MB)

**`rf.poke(addr, val)`** - Write byte
```lua
rf.poke(0x1000, 42)
```

**`rf.peek(addr)`** - Read byte
```lua
local val = rf.peek(0x1000)
```

**`rf.poke2(addr, val)`** / **`rf.peek2(addr)`** - 16-bit operations
**`rf.poke4(addr, val)`** / **`rf.peek4(addr)`** - 32-bit operations

#### Cart Storage (64KB)

**`rf.cstore(dest_addr, src_addr, len)`** - Save to cart storage
```lua
rf.cstore(0, 0x1000, 256)  -- Save 256 bytes from memory to cart
```

**`rf.reload(dest_addr, src_addr, len)`** - Load from cart storage
```lua
rf.reload(0x1000, 0, 256)  -- Load 256 bytes from cart to memory
```

### Multiplayer API

#### Connection Info

**`rf.is_multiplayer()`** - Check if in multiplayer mode
```lua
if rf.is_multiplayer() then
  -- Multiplayer logic
end
```

**`rf.player_count()`** - Get total player count (1-6)
```lua
local count = rf.player_count()
```

**`rf.my_player_id()`** - Get local player ID (1-6)
```lua
local my_id = rf.my_player_id()
```

**`rf.is_host()`** - Check if local player is host
```lua
if rf.is_host() then
  -- Host-only logic (game logic, collision, etc.)
end
```

#### Synchronization

**`rf.network_sync(table, tier)`** - Register table for sync
```lua
rf.network_sync(players, "fast")   -- Smooth movement (30-60/sec)
rf.network_sync(items, "moderate")  -- Moderate updates (15/sec)
rf.network_sync(score, "slow")     -- Infrequent updates (5/sec)
```

**`rf.network_unsync(table)`** - Unregister table
```lua
rf.network_unsync(old_table)
```

#### Input (Host Only)

**`rf.btn(player_id, button)`** - Check another player's input (host only)
```lua
if rf.is_host() then
  for id = 1, rf.player_count() do
    if rf.btn(id, 0) then  -- Player id pressing left
      players[id].x = players[id].x - 3
    end
  end
end
```

---

## Advanced Features

### State Machine Architecture

RetroForge uses a **stack-based state machine** that supports both simple transitions and overlays:

```
Menu State (ChangeState)
  ↓
Playing State (ChangeState)
  ↓
Pause State (PushState) ← Overlay
  ↓
Playing State (PopState) ← Return
  ↓
Game Over (ChangeState)
```

**Benefits:**
- Clean separation of concerns
- Easy to add pause menus, inventory screens, etc.
- Built-in splash and credits screens
- Module-based organization

### Module-Based State System

Instead of manually registering states, use `rf.import()` to automatically load and register states from separate files:

```lua
-- main.lua
rf.import("menu_state.lua")  -- Registers as "menu"
rf.import("play_state.lua")  -- Registers as "play"

-- menu_state.lua
local selected = 1

function _INIT()
  -- One-time setup
end

function _HANDLE_INPUT()
  if rf.btnp(2) then selected = selected - 1 end
  if rf.btnp(4) then game.changeState("play") end
end

function _UPDATE(dt)
  -- Update logic
end

function _DRAW()
  rf.clear_i(0)
  rf.print("MENU", 200, 100, 7)
end

function _DONE()
  -- Cleanup
end
```

### Automatic Sprite Pooling

When a sprite has `isUI=false` and `maxSpawn > 10`, RetroForge automatically creates a pool to reuse instances:

```lua
-- In sprites.json or via code
{
  "bullet": {
    "width": 4,
    "height": 4,
    "isUI": false,
    "maxSpawn": 100  -- Automatically pooled!
  }
}
```

**Benefits:**
- Reduces garbage collection
- Improves performance
- Transparent (no code changes needed)
- Automatic cleanup

### Multiplayer Architecture

RetroForge uses a **host-authoritative** model with **star topology**:

```
         Host (Player 1)
        /  |  |  |  \
       /   |  |  |   \
   P2   P3  P4  P5  P6
```

**How it works:**
1. Host runs game logic in `_UPDATE()`
2. Host checks all players' inputs via `rf.btn(player_id, button)`
3. Host updates game state (players table, etc.)
4. Engine automatically syncs registered tables to all clients
5. Clients receive updates and render locally

**Best Practices:**
- Only host modifies non-player data
- Use appropriate sync tiers (fast/moderate/slow)
- Keep sync tables simple (avoid nested complex structures)
- Test with network latency

---

## Best Practices & Patterns

### Code Organization

**1. Use Module-Based States**
```lua
-- ✅ Good: Separate files
rf.import("menu_state.lua")
rf.import("play_state.lua")

-- ❌ Avoid: Everything in main.lua
```

**2. Keep Functions Focused**
```lua
-- ✅ Good: Single responsibility
function updatePlayer()
  -- Player logic only
end

function updateEnemies()
  -- Enemy logic only
end

-- ❌ Avoid: Monolithic update function
```

**3. Use Descriptive Variable Names**
```lua
-- ✅ Good
local player_velocity_x = 5
local enemy_spawn_timer = 2.0

-- ❌ Avoid
local vx = 5
local t = 2.0
```

### Performance

**1. Minimize Draw Calls**
```lua
-- ✅ Good: Batch operations
function _DRAW()
  rf.clear_i(0)
  drawAllSprites()  -- Single function call
end

-- ❌ Avoid: Many individual calls in _DRAW
```

**2. Use Sprite Pooling**
```lua
-- ✅ Good: Auto-pooled sprites
rf.setSpriteProperty("bullet", "maxSpawn", 100)

-- ❌ Avoid: Creating many sprites manually
```

**3. Optimize Physics**
```lua
-- ✅ Good: Reuse bodies
local bullet_bodies = {}

-- ❌ Avoid: Creating/destroying bodies every frame
```

### State Management

**1. Use Context for Inter-State Data**
```lua
-- Menu state
game.setContext("selected_level", 3)
game.changeState("playing")

-- Playing state
local level = game.getContext("selected_level") or 1
```

**2. Clean Up on Exit**
```lua
function _EXIT()
  -- Stop music
  rf.music("stopall")
  
  -- Clear timers
  -- Reset state
end
```

**3. Use Push/Pop for Overlays**
```lua
-- Playing state
if rf.btnp(4) then
  game.pushState("pause")  -- Don't destroy playing state
end

-- Pause state
if rf.btnp(4) then
  game.popState()  -- Return to playing
end
```

### Multiplayer Patterns

**1. Host-Authoritative Logic**
```lua
function _UPDATE(dt)
  if rf.is_host() then
    -- Host runs ALL game logic
    updatePlayers()
    updateEnemies()
    checkCollisions()
    updateScore()
  end
  -- Non-hosts just receive updates
end
```

**2. Sync Tier Strategy**
```lua
-- Fast: Smooth movement
rf.network_sync(players, "fast")

-- Moderate: Items, powerups
rf.network_sync(items, "moderate")

-- Slow: UI, scores
rf.network_sync(ui_state, "slow")
```

**3. Player ID as Key**
```lua
-- ✅ Good: Each player owns their data
players[rf.my_player_id()] = {x = 100, y = 100}

-- ❌ Avoid: Modifying other players' data
players[other_id].x = 50  -- Only that player should modify
```

---

## Examples & Tutorials

### Example 1: Simple Platformer

```lua
-- main.lua
local player = {x = 100, y = 100, vx = 0, vy = 0, on_ground = false}

function _INIT()
  rf.palette_set("PICO-8")
end

function _UPDATE(dt)
  -- Input
  if rf.btn(0) then player.vx = -3
  elseif rf.btn(1) then player.vx = 3
  else player.vx = 0 end
  
  if rf.btnp(4) and player.on_ground then
    player.vy = -10
    player.on_ground = false
    rf.sfx("jump")
  end
  
  -- Physics
  player.vy = player.vy + 0.5  -- Gravity
  player.x = player.x + player.vx
  player.y = player.y + player.vy
  
  -- Ground collision
  if player.y > 250 then
    player.y = 250
    player.vy = 0
    player.on_ground = true
  end
end

function _DRAW()
  rf.clear_i(0)
  rf.rectfill(0, 250, 480, 270, 5)  -- Ground
  rf.circfill(player.x, player.y, 8, 7)  -- Player
end
```

### Example 2: State-Based Game

```lua
-- main.lua
rf.import("menu_state.lua")
rf.import("play_state.lua")

-- menu_state.lua
local selected = 1
local options = {"START", "QUIT"}

function _HANDLE_INPUT()
  if rf.btnp(2) then selected = 1 end
  if rf.btnp(3) then selected = 2 end
  if rf.btnp(4) then
    if selected == 1 then
      game.changeState("play")
    else
      game.exit()
    end
  end
end

function _DRAW()
  rf.clear_i(0)
  for i, opt in ipairs(options) do
    local color = (i == selected) and 15 or 7
    rf.print_xy(200, 100 + i*20, opt, color)
  end
end
```

### Example 3: Physics Game

```lua
-- main.lua
local ball = {}

function _INIT()
  local body = rf.physics_create_body("dynamic", 240, 100)
  rf.physics_body_add_circle(body, 10, 1.0, 0.9, 0.1)  -- Bouncy
  rf.physics_body_set_gravity_scale(body, 1.0)
  rf.physics_body_set_velocity(body, 5, 0)
  ball.body = body
end

function _UPDATE(dt)
  local x, y = rf.physics_body_get_position(ball.body)
  if y > 300 then
    rf.physics_body_destroy(ball.body)
    -- Reset ball
    local body = rf.physics_create_body("dynamic", 240, 100)
    rf.physics_body_add_circle(body, 10, 1.0, 0.9, 0.1)
    ball.body = body
  end
end

function _DRAW()
  rf.clear_i(0)
  local x, y = rf.physics_body_get_position(ball.body)
  rf.circfill(x, y, 10, 11)
end
```

---

## Troubleshooting

### Common Issues

**1. Game Not Drawing**
- Make sure `_DRAW()` clears the screen first: `rf.clear_i(0)`
- Check that drawing functions are called in `_DRAW()`, not `_UPDATE()`

**2. Input Not Working**
- Use `rf.btnp()` for one-time actions, `rf.btn()` for continuous
- Check button indices (0-5 for standard buttons)

**3. Physics Not Working**
- Make sure `useCollision=true` on sprites
- Check that bodies are created as `"dynamic"` type
- Verify gravity scale (0.0 disables gravity)

**4. Module Import Failing**
- Ensure file exists in `assets/` folder
- Check that all required functions are implemented
- Verify file is included in cart when packing

**5. Multiplayer Not Syncing**
- Only host should modify synced tables
- Use appropriate sync tiers
- Check that `rf.network_sync()` is called in `_INIT()`

### Debugging

**Development Mode:**
```bash
make run-dev FOLDER=my-game
```

**Debug Functions (Dev Mode Only):**
```lua
rf.printh("Debug message")  -- Print to console
local fps = rf.stat(0)  -- Get FPS
local memory = rf.stat(2)  -- Get memory usage
```

**Hot Reload:**
- Files are automatically reloaded when modified
- State machine and module states are preserved
- Assets (sprites, audio) are reloaded

---

## Additional Resources

- **API Reference**: Complete function reference
- **PICO-8 Comparison**: Feature-by-feature comparison
- **Multiplayer Design**: Complete multiplayer architecture
- **Examples**: Full working examples in `examples/` folder

---

**Happy Coding!** 🔨✨

*Forge Your Retro Dreams with RetroForge Engine*

