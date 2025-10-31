# PICO-8 vs RetroForge Engine Comparison

A comprehensive feature-by-feature comparison between PICO-8 (the most popular fantasy console) and RetroForge Engine.

---

## 📊 Feature Comparison Overview

| Feature Category | PICO-8 | RetroForge | Notes |
|-----------------|--------|------------|-------|
| **Display** | 128×128 | 480×270 | RetroForge has higher resolution |
| **Palette** | 16 colors | 50 colors | RetroForge has larger palette |
| **Cart Size Limit** | 32 KB | 64 KB | RetroForge 2x PICO-8 capacity |
| **Built-in Editor** | ✅ Yes (all-in-one) | ❌ No (separate webapp) | PICO-8 has integrated IDE |
| **Sprite System** | ✅ Sprite sheet (8×8) | ✅ JSON-based sprites | Different approaches |
| **Tilemap** | ✅ Yes | ✅ Yes | Both implemented |
| **Camera/Viewport** | ✅ Yes | ✅ Yes | Both implemented |
| **Physics** | ❌ Manual | ✅ Box2D | RetroForge has physics |
| **Audio Channels** | 4 channels | 8 channels | RetroForge more channels |
| **Memory** | 2 MB limit | No strict limit* | RetroForge more flexible |
| **Multiplayer/Networking** | ✅ Built-in (up to 4 players) | ✅ WebRTC (up to 6 players) | RetroForge uses modern WebRTC |
| **State Machine** | ❌ No | ✅ Yes | RetroForge has flexible state management |
| **Export Targets** | Desktop + HTML5 | Desktop + WASM + Android | RetroForge has mobile |

\* *Limits may exist but are not strictly enforced*

---

## 🎮 Graphics Functions

### PICO-8 Graphics API

```lua
-- Screen
cls([color])                    -- Clear screen
camera([x, y])                  -- Set camera offset
pal([c0, c1], [p])              -- Color swap/remap

-- Drawing Primitives
pset(x, y, [color])             -- Set pixel
pget(x, y)                      -- Get pixel
line(x0, y0, x1, y1, [color])  -- Line
rect(x0, y0, x1, y1, [color])   -- Rectangle outline
rectfill(x0, y0, x1, y1, [color]) -- Filled rectangle
circ(x, y, r, [color])          -- Circle outline
circfill(x, y, r, [color])     -- Filled circle
elli(x, y, rx, ry, [color])     -- Ellipse outline
ellifill(x, y, rx, ry, [color]) -- Filled ellipse

-- Sprites
spr(n, x, y, [w, h, flip_x, flip_y]) -- Draw sprite
sspr(sx, sy, sw, sh, dx, dy, [dw, dh, flip_x, flip_y]) -- Sprite region

-- Text
print(str, [x, y, [color]])     -- Print text (supports \n, \#)
cursor(x, y)                     -- Set text cursor
color([color])                  -- Set text color

-- Clipping & Masking
clip([x, y, w, h])              -- Set clipping rectangle
```

### RetroForge Graphics API

```lua
-- Screen
rf.clear_i(index)               -- Clear screen
rf.camera([x, y])               -- Set camera offset (no args = reset)
rf.clip([x, y, w, h])          -- Set clipping rectangle (no args = disable)

-- Drawing Primitives
rf.pset(x, y, index)            -- Set pixel
rf.pget(x, y)                   -- Get pixel (returns {r, g, b, a} table)
rf.line(x0, y0, x1, y1, index)  -- Line
rf.rect(x0, y0, x1, y1, index)  -- Rectangle outline
rf.rectfill(x0, y0, x1, y1, index) -- Filled rectangle
rf.circ(x, y, radius, index)    -- Circle outline
rf.circfill(x, y, radius, index) -- Filled circle
rf.elli(x, y, rx, ry, index)    -- Ellipse outline
rf.ellifill(x, y, rx, ry, index) -- Filled ellipse

-- Additional Shapes (RetroForge only)
rf.triangle(x, y, radius, [filled], index)
rf.diamond(x, y, radius, [filled], index)
rf.square(x, y, radius, [filled], index)
rf.pentagon(x, y, radius, [filled], index)
rf.hexagon(x, y, radius, [filled], index)
rf.star(x, y, radius, [filled], index)

-- Text
rf.print(text, [x, y, index]) -- Print text (PICO-8-like, uses cursor/color state if args omitted)
rf.print_xy(x, y, text, [index]) -- Print at position (uses color state if index omitted)
rf.print_anchored(text, anchor, index) -- Print with anchor (9 positions)
rf.cursor([x, y]) -- Set text cursor (no args = reset)
rf.color([index]) -- Set text color (no args = reset)

-- Sprites
rf.sprite(name)                 -- Get sprite data (returns table)
rf.spr(name, x, y, [flip_x, flip_y]) -- Draw sprite by name
rf.sspr(name, sx, sy, sw, sh, dx, dy, [dw, dh, flip_x, flip_y]) -- Sprite region

-- Color Remapping
rf.pal([c0, c1, p])             -- Remap color index (no args = reset)

-- Palette
rf.palette_set(name)            -- Set active palette
```

### Missing in RetroForge

⚠️ **Tilemap system** - Implemented with `map()` / `mget()` / `mset()`, but different API than PICO-8

---

## 🎹 Audio Functions

### PICO-8 Audio API

```lua
-- Sound Effects
sfx(n, [channel], [offset])     -- Play SFX by index
sfx(n, -1)                      -- Stop all SFX

-- Music
music([n], [fade_len], [channel_mask]) -- Play music track
stop([channel_mask])            -- Stop music

-- Audio Data Access
poke(0x5f50, v)                 -- Set SFX volume
poke(0x5f51, v)                 -- Set music volume
```

### RetroForge Audio API

```lua
-- Sound Effects
rf.sfx(name, [action])          -- Play SFX by name (from sfx.json)
rf.sfx("sine", freq, duration, gain) -- Inline SFX
rf.tone(frequency, duration, [gain])  -- Play tone
rf.noise(duration, [gain])      -- Play noise

-- Music
rf.music(name)                  -- Play track by name (from music.json)
rf.music(tokens, bpm, gain)     -- Inline music definition
```

### Differences

| Feature | PICO-8 | RetroForge |
|---------|--------|------------|
| **SFX Format** | Index-based (0-63) | Name-based (JSON) |
| **Music Format** | Index-based (0-63) | Name-based (JSON) |
| **Channels** | 4 channels | 8 channels |
| **SFX Editor** | Built-in tracker | JSON-based |
| **Music Editor** | Built-in tracker | JSON-based |
| **Volume Control** | `poke()` memory | Not exposed yet |
| **Channel Selection** | Manual channel pick | Automatic assignment |

---

## 🎮 Input Functions

### PICO-8 Input API

```lua
btn([i], [p])                   -- Button held (player 0-7)
btnp([i], [p])                  -- Button just pressed
```

**Button Indexes:**
- 0 = ❌ (O/X)
- 1 = ❌ (O/X) (alternative)
- 2 = ⬆️
- 3 = ⬇️
- 4 = ⬅️
- 5 = ➡️
- 6 = ❌ (O/X) (alternative)
- 7 = ❌ (O/X) (alternative)

### RetroForge Input API

```lua
rf.btn(button)                  -- Button held
rf.btnp(button)                 -- Button just pressed
rf.btnr(button)                 -- Button just released (RetroForge only)
```

**Button Indexes:**
- 0-15: Standard buttons
- 16+: Extended buttons

### Differences

✅ **RetroForge has `btnr()`** - Button release detection (not in PICO-8)  
✅ **RetroForge supports more buttons** - 0-15+ vs PICO-8's 0-7  
✅ **Multi-player Input** - WebRTC-based networking with up to 6 players  

---

## 🗺️ Map/Tilemap System

### PICO-8 Map System

```lua
-- Tilemap Functions
mget(x, y)                      -- Get tile at map coordinate
mset(x, y, v)                   -- Set tile at map coordinate
map(cel_x, cel_y, sx, sy, cel_w, cel_h, [layer]) -- Draw map region
fget(index, [flag])             -- Get sprite flag
fset(index, flag, [val])        -- Set sprite flag
```

**Features:**
- 128×64 tile map
- 8×8 tiles
- Sprite flags for collision/custom data
- Layered rendering support

### RetroForge Map System

```lua
-- Tilemap Functions
rf.mget(x, y)                      -- Get tile at map coordinate
rf.mset(x, y, v)                   -- Set tile at map coordinate
rf.map(cel_x, cel_y, sx, sy, cel_w, cel_h) -- Draw map region
```

**Features:**
- 256×256 tile map (default, larger than PICO-8)
- Flexible tile rendering
- Camera system integration ✅

---

## 🖼️ Sprite System

### PICO-8 Sprite System

```lua
-- Sprite Sheet
-- 128×128 pixel sheet, 8×8 sprites = 16×16 grid
spr(n, x, y, [w, h, flip_x, flip_y]) -- Draw sprite
sspr(sx, sy, sw, sh, dx, dy, [dw, dh, flip_x, flip_y]) -- Sprite region

-- Sprite Flags
fget(index, [flag])             -- Get flag (collision, etc.)
fset(index, flag, [val])       -- Set flag
```

**Features:**
- Fixed 8×8 pixel sprites
- Built-in flip support
- Sprite flags for game logic
- Integrated sprite editor

### RetroForge Sprite System

```lua
-- Built-in Sprite Drawing
rf.spr("player", x, y)                    -- Draw sprite by name
rf.spr("player", x, y, true, false)       -- With horizontal flip
rf.sspr("player", sx, sy, sw, sh, dx, dy) -- Draw sprite region

-- Sprite Data Access (optional)
local spr = rf.sprite("player")
-- Returns: {width, height, pixels, useCollision, mountPoints}
-- Manual drawing still possible via spr.pixels array if needed
```

**Features:**
- Flexible sizes (8×8 to 32×32)
- JSON-based storage
- Collision metadata (`useCollision`)
- Mount points for projectiles/thrusters
- Named mount points (access by index or name)

**Differences:**
- PICO-8: Fixed 8×8, built-in rendering
- RetroForge: Flexible sizes (8×8 to 32×32), built-in rendering via `rf.spr()` / `rf.sspr()`
- RetroForge: More metadata (collision, mount points)
- PICO-8: Sprite flags for custom data

---

## 💾 Memory & Storage

### PICO-8 Memory

```lua
-- Cart Data
-- Fixed 32 KB cart size limit
-- Sprites: 8 KB (128×128, 1 byte per pixel)
-- Map: 8 KB (128×64, 1 byte per tile) - Note: RetroForge uses 256×256
-- SFX: 2 KB (64 sounds)
-- Music: 1 KB (64 songs)
-- Code: ~15 KB remaining

-- Runtime Memory
-- 2 MB RAM limit
poke(addr, val)                 -- Write byte to memory
peek(addr)                      -- Read byte from memory
poke2(addr, val)                -- Write 16-bit value
peek2(addr)                     -- Read 16-bit value
poke4(addr, val)                -- Write 32-bit value
peek4(addr)                     -- Read 32-bit value

-- Cart Data Access
cstore(dest_addr, src_addr, len) -- Copy from runtime memory to cart storage
reload(dest_addr, src_addr, len) -- Copy from cart storage to runtime memory
```

### RetroForge Memory

```lua
-- Cart Data
-- 64 KB cart size limit (2x PICO-8's 32 KB)
-- JSON-based storage (sprites.json, sfx.json, music.json)
-- Lua code in main.lua

-- Runtime Memory
rf.poke(addr, val)                 -- Write byte to memory
rf.peek(addr)                      -- Read byte from memory
rf.poke2(addr, val)                -- Write 16-bit value
rf.peek2(addr)                     -- Read 16-bit value
rf.poke4(addr, val)                -- Write 32-bit value
rf.peek4(addr)                     -- Read 32-bit value

-- Cart Data Access
rf.cstore(dest_addr, src_addr, len) -- Copy from runtime memory to cart storage
rf.reload(dest_addr, src_addr, len) -- Copy from cart storage to runtime memory
```

**Features:**
- 2MB memory (same as PICO-8)
- 64KB cart storage (2x PICO-8's 32KB)
- Full poke/peek API support ✅
- Cart persistence API ✅

**Differences:**
- PICO-8: Fixed memory layout, direct memory access
- RetroForge: Flexible storage, JSON-based, full memory access API ✅
- PICO-8: Cart persistence via `cstore()` / `reload()`
- RetroForge: Cart persistence via `rf.cstore()` / `rf.reload()` ✅

---

## 📝 Development Environment

### PICO-8

✅ **All-in-one IDE:**
- Built-in code editor (with syntax highlighting)
- Sprite editor (pixel art)
- Map editor (tile placement)
- SFX editor (tracker-style)
- Music editor (tracker-style)
- Cart browser
- Splore (game sharing platform)

✅ **Hot reload:** Edit and test instantly  
✅ **Built-in debugger:** Step through code  
✅ **Command mode:** `save`, `load`, `folder`, `export`, etc.  

### RetroForge

❌ **No built-in editor** (separate webapp):
- Code editor (via webapp `/editor/code`)
- Sprite editor (via webapp `/editor/sprite`)
- Sound editor (via webapp `/editor/sound`)
- Music editor (via webapp `/editor/music`)
- Properties editor (via webapp `/editor/properties`)

✅ **Hot reload** (development mode only - when running from folder)  
✅ **Built-in debugger** (development mode only - `rf.printh()`, `rf.stat()`, `rf.time()`)  
✅ **Command mode** (via webapp: save, load, export functionality)

**Webapp Features:**
- Multi-file editing (main.lua, manifest.json, etc.)
- Palette editor
- Sprite management (add/delete/rename/duplicate)
- SFX/Music JSON editors
- Cart sharing (via Convex backend)
- Save/load/export commands (available in webapp UI)

**Development Mode Features (local only):**
- Hot reload: Automatic reload when files change (folder-based development)
- Debug logging: `rf.printh()` for console output
- Statistics: `rf.stat()` for FPS, memory, load times
- Time functions: `rf.time()` for timestamps
- File watching: Monitors `assets/` directory and `manifest.json`

---

## 🌐 Platform & Export

### PICO-8 Export Targets

✅ **Desktop:**
- Windows (standalone .exe)
- macOS (standalone .app)
- Linux (standalone binary)

✅ **Web:**
- HTML5 export (runs in browser)
- Embeds in web pages

✅ **Mobile:**
- ❌ No native mobile support
- Only via web export

### RetroForge Export Targets

✅ **Desktop:**
- Windows (standalone .exe)
- macOS (standalone .app)
- Linux (standalone binary)

✅ **Web:**
- WASM export (runs in browser)
- Webapp integration

✅ **Mobile:**
- ✅ Android (native support)
- ❌ iOS (not yet)

---

## 📚 API Completeness

### PICO-8 Full API Reference

```lua
-- Math
abs(x), atan2(y, x), cos(x), sin(x), sqrt(x), sgn(x)
flr(x), ceil(x), rnd(x), mid(x, y, z)
shl(x, y), shr(x, y), band(x, y), bor(x, y), bxor(x, y), bnot(x)

-- String
sub(s, start, [end]), chr(c), ord(s)
split(s, [sep]), join(t, [sep])

-- Table
add(t, v), del(t, v), all(t), foreach(t, f), count(t)

-- Utility
stat(n)                         -- Get system stats
time()                          -- Get time in seconds
printh(str, [filename])         -- Print to console/file

-- Graphics (see above)
-- Audio (see above)
-- Input (see above)
-- Map (see above)
```

### RetroForge API

```lua
-- Math
-- Uses standard Lua math library (math.abs, math.cos, etc.)
rf.flr(x)                       -- Floor: truncate towards zero
rf.ceil(x)                      -- Ceiling: round up
rf.rnd([x])                     -- Random number (0-1 or 0-x)
rf.mid(x, y, z)                 -- Clamp x between y and z
rf.sgn(x)                       -- Sign: -1, 0, or 1

-- String Helpers
rf.chr(n)                       -- Convert number to character
rf.ord(c)                       -- Convert character to number (first byte)

-- Bitwise Operations
rf.shl(x, y)                    -- Shift left: x << y
rf.shr(x, y)                    -- Shift right (arithmetic): x >> y
rf.band(x, y)                   -- Bitwise AND: x & y
rf.bor(x, y)                    -- Bitwise OR: x | y
rf.bxor(x, y)                   -- Bitwise XOR: x ^ y
rf.bnot(x)                      -- Bitwise NOT: ~x

-- String
-- Uses standard Lua string library (string.sub, etc.)

-- Table
-- Uses standard Lua table library (table.insert, etc.)

-- Utility
rf.quit()                       -- Request quit

-- Graphics (see above)
-- Audio (see above)
-- Input (see above)
-- Map (see above) ✅
```

**Available Utilities:**
✅ `rf.stat(n)` - System statistics (development mode only)
✅ `rf.time()` - Precise time functions (development mode only)
✅ `rf.printh(str)` - Debug output (development mode only)
✅ `rf.flr(x)`, `rf.ceil(x)`, `rf.rnd([x])`, `rf.mid(x, y, z)`, `rf.sgn(x)` - Math helper functions
✅ `rf.chr(n)`, `rf.ord(c)` - String helper functions
✅ `rf.shl(x, y)`, `rf.shr(x, y)` - Bitwise shift operations
✅ `rf.band(x, y)`, `rf.bor(x, y)`, `rf.bxor(x, y)`, `rf.bnot(x)` - Bitwise logical operations
- **Note:** Debug functions only work when running from a folder (development mode), not in packed `.rf` files
- **Note:** Helper functions and bitwise operations are always available (not dev-mode only)  

---

## 🎯 Unique RetroForge Features

✅ **Additional Shapes:** Triangle, Diamond, Square, Pentagon, Hexagon, Star  
✅ **Sprite Mount Points:** Named mount points for projectiles/thrusters  
✅ **Sprite Collision Metadata:** `useCollision` flag  
✅ **50-Color Palette:** More colors than PICO-8's 16  
✅ **Higher Resolution:** 480×270 vs 128×128  
✅ **Button Release:** `btnr()` for edge-triggered release  
✅ **8 Audio Channels:** vs PICO-8's 4  
✅ **Android Support:** Native mobile support  
✅ **JSON-Based Assets:** Flexible sprite/SFX/music storage  
✅ **Physics Engine:** Box2D integration (fully implemented)  
✅ **Feature Parity:** All core PICO-8 graphics/tilemap/camera/memory APIs implemented  
✅ **State Machine:** Flexible state management system with lifecycle callbacks, state stacking, and shared context  
✅ **Node System:** Godot-style scene graph (planned)  
✅ **Multiplayer:** WebRTC-based multiplayer support (up to 6 players) with automatic state sync  

---

## 📋 Summary: PICO-8 Feature Status in RetroForge

### ✅ Implemented Features

**Core Graphics & Rendering:**
- **✅ Tilemap System** - `rf.map()`, `rf.mget()`, `rf.mset()` functions fully implemented
- **✅ Camera System** - `rf.camera()` function for viewport control
- **✅ Sprite Drawing** - `rf.spr()` / `rf.sspr()` functions for drawing sprites
- **✅ Pixel Reading** - `rf.pget()` to read screen pixels
- **✅ Color Remapping** - `rf.pal()` for runtime color swapping
- **✅ Clipping** - `rf.clip()` for drawing regions
- **✅ Ellipse Drawing** - `rf.elli()` / `rf.ellifill()` functions

**Memory & Storage:**
- **✅ Memory Access** - `rf.poke()`, `rf.peek()` functions (2MB memory)
- **✅ Cart Persistence** - `rf.cstore()` / `rf.reload()` for saving/loading cart data (64KB)

**Development Tools:**
- **✅ Hot Reload** - Live code editing (development mode only - when running from folder)
- **✅ Debugger** - Built-in debugging tools (development mode only - `rf.printh()`, `rf.stat()`, `rf.time()`)
- **✅ Command Mode** - Save, load, export commands available in webapp
- **✅ System Stats** - `rf.stat()` function for system statistics (development mode only)
- **✅ Text Cursor** - `rf.cursor()` / `rf.color()` state management

**Implemented:**
- **✅ Multi-player Support** - WebRTC-based networking with up to 6 players
  - `rf.is_multiplayer()`, `rf.player_count()`, `rf.is_host()`, `rf.my_player_id()`
  - `rf.network_sync(table, tier)` for automatic state synchronization
  - `rf.btn(player_id, button)` for host to check other players' inputs
  - 3-tier sync system (fast/moderate/slow)
  - Host authority model with star topology

### ❌ Missing / Different Approach

1. **❌ Built-in IDE** - PICO-8's all-in-one desktop editor vs RetroForge's separate webapp
   - RetroForge provides equivalent functionality through web-based editors
   - Different architectural approach (browser-based vs desktop application)

---

## 🎯 Priority Recommendations for RetroForge

### ✅ Completed High Priority Features

1. **✅ Tilemap System** - Implemented!
   - `rf.mget(x, y)`, `rf.mset(x, y, v)`, `rf.map(...)` functions

2. **✅ Camera System** - Implemented!
   - `rf.camera(x, y)` function
   - Integration with all drawing operations

3. **✅ Sprite Drawing Functions** - Implemented!
   - `rf.spr(name, x, y, [flip_x, flip_y])`
   - `rf.sspr(...)` for sprite regions

4. **✅ Pixel Reading** - Implemented!
   - `rf.pget(x, y)` function

5. **✅ Clipping** - Implemented!
   - `rf.clip(x, y, w, h)` function

6. **✅ Color Remapping** - Implemented!
   - `rf.pal(c0, c1, p)` function

7. **✅ Ellipse Drawing** - Implemented!
   - `rf.elli(x, y, rx, ry, index)` / `rf.ellifill(...)`

8. **✅ Memory API** - Implemented!
   - `rf.poke(addr, val)`, `rf.peek(addr)` functions
   - `rf.poke2/peek2`, `rf.poke4/peek4` for 16/32-bit values

9. **✅ Physics Engine** - Box2D integrated!
   - Full physics body, fixture, and force API

### ✅ Completed Medium Priority Features

10. **✅ Hot Reload** - Implemented!
    - Automatic file watching for `assets/` directory and `manifest.json`
    - Auto-reload when files change (development mode only)
    - 500ms cooldown to prevent rapid reloads
    - Available when running with `--folder` flag

11. **✅ Debug Tools** - Implemented!
    - `rf.printh(str)` for console/debug logging
    - `rf.stat(n)` for system statistics (FPS, memory, load times, reload info)
    - `rf.time()` for precise timestamps
    - Debug log history with timestamps
    - Available in development mode only

### ✅ Completed Medium Priority Features (continued)

12. **✅ Cart Persistence** - Implemented!
    - `rf.cstore(dest_addr, src_addr, len)` - Copy from runtime memory to cart storage
    - `rf.reload(dest_addr, src_addr, len)` - Copy from cart storage to runtime memory
    - 64KB cart storage (2x PICO-8's 32KB)
    - Full address validation and bounds clamping

13. **✅ Text Cursor/Color State** - Implemented!
    - `rf.cursor([x, y])` - Set/reset text cursor position
    - `rf.color([index])` - Set/reset text color
    - `rf.print()` - PICO-8-like print function using cursor/color state
    - Cursor automatically advances after printing (handles newlines)

### ✅ Completed Low Priority Features

13. **✅ Command Mode** - Implemented in webapp!
    - Save, load, export functionality available via webapp UI
    - Cart management and sharing via Convex backend

### ✅ Completed Multiplayer Features

14. **✅ Multi-player Support** - Implemented!
    - `rf.is_multiplayer()` - Check if in multiplayer mode
    - `rf.player_count()` - Get number of players (1-6)
    - `rf.my_player_id()` - Get local player ID (1-6)
    - `rf.is_host()` - Check if local player is host
    - `rf.btn(player_id, button)` - Host can check other players' inputs
    - `rf.network_sync(table, tier)` - Register tables for automatic sync
    - `rf.network_unsync(table)` - Unregister tables from sync
    - 3-tier sync system: "fast" (30-60/sec), "moderate" (15/sec), "slow" (5/sec)
    - Host authority model with star topology
    - WebRTC-based real-time multiplayer support

15. **✅ State Machine** - Implemented!
    - `game.registerState(name, stateTable)` - Register states with lifecycle callbacks
    - `game.changeState(name)` - Complete state transitions
    - `game.pushState(name)` / `game.popState()` - State stacking for overlays
    - `game.setContext(key, value)` / `game.getContext(key)` - Shared context for data passing
    - `game.addCredit(category, name, role)` - Add credits for credits screen
    - Built-in engine splash and credits states
    - Complete lifecycle: Initialize, Enter, HandleInput, Update, Draw, Exit, Shutdown

### Low Priority (Coming Soon)

16. **Tilemap Layers** - Multi-layer support
    - 8-layer tilemap system (currently single layer)

---

## 📊 Feature Parity Score

| Category | PICO-8 | RetroForge | Parity |
|----------|--------|------------|--------|
| **Graphics Primitives** | 9 | 15 | ✅ Better (more shapes + ellipses) |
| **Sprite System** | Built-in | Built-in | ✅ Similar (name-based vs index-based) |
| **Tilemap** | ✅ | ✅ | ✅ Implemented |
| **Camera** | ✅ | ✅ | ✅ Implemented |
| **Clipping** | ✅ | ✅ | ✅ Implemented |
| **Color Remapping** | ✅ | ✅ | ✅ Implemented |
| **Pixel Reading** | ✅ | ✅ | ✅ Implemented |
| **Memory Access** | ✅ | ✅ | ✅ Implemented |
| **Cart Persistence** | ✅ | ✅ | ✅ Implemented |
| **Text Cursor/Color** | ✅ | ✅ | ✅ Implemented |
| **Physics** | ❌ | ✅ Box2D | ✅ RetroForge advantage |
| **Audio** | Built-in | JSON-based | ⚠️ Different approach |
| **Input** | 7 | 8 | ✅ Similar |
| **Development Tools** | ✅ IDE | ✅ Webapp + Dev Mode | ✅ Command mode available in webapp |
| **Command Mode** | ✅ Built-in | ✅ Webapp UI | ✅ Save/load/export available |
| **Multi-player** | ✅ Built-in (up to 4 players) | ✅ WebRTC (up to 6 players) | ✅ Implemented |
| **State Machine** | ❌ No | ✅ Yes | ✅ RetroForge advantage |
| **Platforms** | Desktop + Web | Desktop + Web + Android | ✅ More platforms |
| **Resolution** | 128×128 | 480×270 | ✅ Higher resolution |
| **Palette** | 16 colors | 50 colors | ✅ More colors |

**Overall Assessment:** RetroForge now has feature parity with PICO-8's core graphics, tilemap, camera, memory APIs, and multiplayer support. Hot reload and debug tools are available in development mode (when running from folders). Command mode features (save, load, export) are available in the webapp UI. Multiplayer support is implemented via WebRTC with automatic state synchronization. RetroForge offers additional features like state machine management, physics, higher resolution, larger palette, Android support, and modern WebRTC-based multiplayer (up to 6 players). The main remaining advantage for PICO-8 is the all-in-one IDE experience, though RetroForge's webapp provides equivalent functionality in a browser-based interface.

---

*Last Updated: RetroForge Engine now includes tilemap (256×256), camera, clipping, color remapping, pixel reading, sprite drawing, ellipse drawing, memory API, Box2D physics integration, hot reload (dev mode), debug tools (dev mode), cart persistence (cstore/reload, 64KB), text cursor/color state, command mode features (save/load/export via webapp), WebRTC-based multiplayer support (up to 6 players) with automatic state synchronization, and a flexible state machine system with lifecycle callbacks, state stacking, and shared context. Full feature parity achieved with PICO-8's core APIs, including multiplayer, plus additional features like state management and physics.*

