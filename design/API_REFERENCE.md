# RetroForge Engine API Reference

**All colors must use palette indices (0-49, -1 for transparent). RGB functions are not available.**

## Graphics

### Screen Clearing
- `rf.clear_i(index)` - Clear screen using palette color index (0-49, -1 for transparent)

### Text
- `rf.print(text, [x, y, index])` - Print text. If x/y/index omitted, uses cursor/color state. Returns cursor after text.
- `rf.print_anchored(text, anchor, index)` - Print text using anchor position and palette color index
  - Anchors: `"topleft"`, `"topcenter"`, `"topright"`, `"middleleft"`, `"middlecenter"`, `"middleright"`, `"bottomleft"`, `"bottomcenter"`, `"bottomright"`
- `rf.print_xy(x, y, text, [index])` - Print text at exact (x, y) position. If index omitted, uses color state.
- `rf.cursor([x, y])` - Set text cursor position. Call with no arguments to reset cursor.
- `rf.color([index])` - Set text color index. Call with no arguments to reset color.

### Primitives
- `rf.pset(x, y, index)` - Set pixel using palette color index
- `rf.line(x0, y0, x1, y1, index)` - Draw line using palette color index
- `rf.rect(x0, y0, x1, y1, index)` - Draw rectangle outline using palette color index
- `rf.rectfill(x0, y0, x1, y1, index)` - Draw filled rectangle using palette color index
- `rf.circ(x, y, radius, index)` - Draw circle outline using palette color index
- `rf.circfill(x, y, radius, index)` - Draw filled circle using palette color index

### Shape Primitives
- `rf.triangle(x, y, radius, [filled], index)` - Draw triangle (pointing up). `filled` is optional (default: false)
- `rf.diamond(x, y, radius, [filled], index)` - Draw diamond shape. `filled` is optional (default: false)
- `rf.square(x, y, radius, [filled], index)` - Draw square centered at (x, y). `filled` is optional (default: false)
- `rf.pentagon(x, y, radius, [filled], index)` - Draw pentagon. `filled` is optional (default: false)
- `rf.hexagon(x, y, radius, [filled], index)` - Draw hexagon. `filled` is optional (default: false)
- `rf.star(x, y, radius, [filled], index)` - Draw 10-point star. `filled` is optional (default: false)

### Ellipse Drawing
- `rf.elli(x, y, rx, ry, index)` - Draw ellipse outline with radii rx, ry
- `rf.ellifill(x, y, rx, ry, index)` - Draw filled ellipse with radii rx, ry

### Pixel Reading
- `rf.pget(x, y)` - Get pixel color at (x, y). Returns table with `{r, g, b, a}` values (0-255)

### Clipping
- `rf.clip([x, y, w, h])` - Set clipping rectangle. Call with no arguments to disable clipping
- Clipping restricts all drawing operations to the specified rectangle

### Camera
- `rf.camera([x, y])` - Set camera offset. All drawing operations are offset by (x, y). Call with no arguments to reset (0, 0)

### Sprite Drawing
- `rf.spr(name, x, y, [flip_x, flip_y])` - Draw sprite by name at position (x, y). Optional horizontal/vertical flipping
- `rf.sspr(name, sx, sy, sw, sh, dx, dy, [dw, dh, flip_x, flip_y])` - Draw sprite region. Scales from source (sx, sy, sw, sh) to destination (dx, dy, dw, dh)

### Tilemap
- `rf.mget(x, y)` - Get tile value at map coordinate (x, y). Returns tile index (0 = empty)
- `rf.mset(x, y, v)` - Set tile at map coordinate (x, y) to value v
- `rf.map(cel_x, cel_y, sx, sy, cel_w, cel_h)` - Draw tilemap region. cel_x/cel_y = tile coordinates, sx/sy = screen position, cel_w/cel_h = tiles to draw

### Color Remapping
- `rf.pal([c0, c1, p])` - Remap color index. `pal(c0, c1)` maps color c0 to c1. `pal()` with no args resets all remapping
- `p` parameter (optional, default true) enables/disables the remap
- Color remapping affects all drawing operations that use palette indices

### Memory API
- `rf.poke(addr, val)` - Write byte value to memory address
- `rf.peek(addr)` - Read byte value from memory address. Returns 0 if out of bounds
- `rf.poke2(addr, val)` - Write 16-bit value (little-endian)
- `rf.peek2(addr)` - Read 16-bit value (little-endian)
- `rf.poke4(addr, val)` - Write 32-bit value (little-endian)
- `rf.peek4(addr)` - Read 32-bit value (little-endian)
- Memory size: 2MB (like PICO-8)

### Cart Persistence
- `rf.cstore(dest_addr, src_addr, len)` - Copy `len` bytes from runtime memory (src_addr) to cart storage (dest_addr)
- `rf.reload(dest_addr, src_addr, len)` - Copy `len` bytes from cart storage (src_addr) to runtime memory (dest_addr)
- Cart storage size: 64KB (2x PICO-8's 32KB)
- Used for saving/loading game state that persists across cart reloads
- Addresses are validated and clamped to available bounds

### Physics (Box2D Integration)
- `rf.physics_create_body(type, x, y)` - Create physics body. `type` = `"static"`, `"dynamic"`, or `"kinematic"`. Returns body ID
- `rf.physics_body_add_box(body_id, width, height, [density])` - Add box fixture to body
- `rf.physics_body_add_circle(body_id, radius, [density])` - Add circle fixture to body
- `rf.physics_body_set_position(body_id, x, y)` - Set body position
- `rf.physics_body_get_position(body_id)` - Get body position. Returns x, y
- `rf.physics_body_set_velocity(body_id, vx, vy)` - Set body linear velocity
- `rf.physics_body_get_velocity(body_id)` - Get body velocity. Returns vx, vy
- `rf.physics_body_apply_force(body_id, fx, fy, px, py)` - Apply force at point (px, py)
- `rf.physics_body_destroy(body_id)` - Destroy body and remove from world
- Physics world has default gravity (0, 9.8) - downward like real physics

## Input

### `rf.btn(button)`
Check if button is pressed. Returns boolean.
- Button 0-15: Standard buttons
- Button 16+: Extended buttons

### `rf.btnp(button)`
Check if button was just pressed this frame (edge-triggered). Returns boolean.

### `rf.btnr(button)`
Check if button was just released this frame. Returns boolean.

## Audio

### `rf.sfx(name, ...)`
Play sound effect by name (from `sfx.json`). Fallback to inline parameters if not found.
- `rf.sfx("jump")` - Play named SFX
- `rf.sfx("sine", freq, duration, gain)` - Inline SFX definition
- Types: `"sine"`, `"noise"`, `"thrust"`, `"stopall"`

### `rf.music(name, ...)`
Play music track by name (from `music.json`). Fallback to inline parameters if not found.
- `rf.music("menu")` - Play named track
- `rf.music(tokens, bpm, gain)` - Inline music definition
- `tokens`: Array of note strings (e.g., `{"4C1", "4E1", "R1"}`)
- `bpm`: Beats per minute (optional)
- `gain`: Volume 0.0-1.0 (optional)

## Sprites

### `rf.sprite(name)`
Get sprite data by name (from `sprites.json`). Returns table:
```lua
{
  width = 16,
  height = 16,
  pixels = {{row1}, {row2}, ...},  -- 2D array of color indices
  useCollision = true,
  mountPoints = {
    [1] = {x = 8, y = 8, name = "thrust"},  -- Access by index
    ["thrust"] = {x = 8, y = 8, name = "thrust"}  -- Access by name
  }
}
```

## Palette

### `rf.palette_set(name)`
Set active palette by name (e.g., `"RetroForge 50"`, `"Grayscale 50"`).

## System

### `rf.quit()`
Request application quit.

## State Machine API

The state machine provides a flexible system for managing game flow through different states (menus, playing, pause screens, etc.). Use `game.*` functions to register and manage states.

### State Registration

#### `game.registerState(name, stateTable)`
Registers a new state with the state machine.

**Parameters:**
- `name` (string): Unique name for the state
- `stateTable` (table): Table containing state lifecycle callbacks:
  - `initialize(sm)` (optional): Called once when state is first created
  - `enter(sm)` (optional): Called every time state becomes active
  - `handleInput(sm)` (optional): Called each frame for input processing
  - `update(dt)` (optional): Called each frame for game logic (dt is delta time in seconds)
  - `draw()` (optional): Called each frame for rendering
  - `exit(sm)` (optional): Called when leaving the state
  - `shutdown()` (optional): Called once when state is destroyed

**Example:**
```lua
local MenuState = {
  selectedOption = 1,
  
  initialize = function(sm)
    -- Load assets, allocate resources (called once)
  end,
  
  enter = function(sm)
    selectedOption = 1
    -- Reset variables, start music (called every time state becomes active)
  end,
  
  handleInput = function(sm)
    if rf.btnp(2) then -- Up
      selectedOption = math.max(1, selectedOption - 1)
    elseif rf.btnp(3) then -- Down
      selectedOption = math.min(3, selectedOption + 1)
    elseif rf.btnp(4) then -- A button
      if selectedOption == 1 then
        game.changeState("playing")
      end
    end
  end,
  
  update = function(dt)
    -- Animate menu items
  end,
  
  draw = function()
    rf.clear_i(0)
    rf.print("MENU", 100, 50, 7)
    -- Draw menu options
  end,
  
  exit = function(sm)
    -- Pause music, save temporary data
  end,
  
  shutdown = function()
    -- Free assets (called once when state is destroyed)
  end
}

game.registerState("menu", MenuState)
```

#### `game.unregisterState(name)`
Removes a state from the registry. Cannot unregister states that are currently in the stack.

### State Transitions

#### `game.changeState(name)`
Replaces all states in the stack with a new state. Use for complete transitions (e.g., menu → playing).

#### `game.pushState(name)`
Adds a new state on top of the current stack. Previous state pauses but stays in memory. Use for overlays (e.g., playing → pause menu).

#### `game.popState()`
Removes the top state from the stack and returns to the previous state. Use to close overlays (e.g., close pause menu → resume playing).

#### `game.popAllStates()`
Removes all states from the stack.

### Shared Context

States can share data through a persistent context that survives state transitions.

#### `game.setContext(key, value)`
Stores a value in the shared context. Value can be string, number, boolean, or table.

#### `game.getContext(key)`
Retrieves a value from the shared context. Returns `nil` if key doesn't exist.

#### `game.hasContext(key)`
Returns `true` if the key exists in context, `false` otherwise.

#### `game.clearContext(key)`
Removes a specific key from the context.

#### `game.clearAllContext()`
Clears all context data.

**Example:**
```lua
-- In level select state
game.setContext("current_level", 3)
game.setContext("difficulty", "hard")
game.changeState("playing")

-- In playing state
local level = game.getContext("current_level") or 1
local difficulty = game.getContext("difficulty") or "normal"
```

### Credits API

#### `game.addCredit(category, name, role)`
Adds a credit entry to be displayed in the credits state. Categories: "Developer", "Artist", "Sound", "Music", "Special Thanks", etc.

**Example:**
```lua
game.addCredit("Developer", "Jane Doe", "Lead Developer")
game.addCredit("Artist", "John Smith", "Character Artist")
```

### Control

#### `game.exit()`
Transitions to credits state, then exits the game. Credits state will display all added credits before exit.

### Utility

#### `game.drawPreviousState()`
Draws the state underneath the current state in the stack. Useful for overlays (e.g., pause menu showing dimmed game behind it).

#### `game.getStackDepth()`
Returns the number of states currently in the stack.

**Example - Pause Menu Overlay:**
```lua
local PauseState = {
  draw = function()
    -- Draw the previous state (dimmed game)
    game.drawPreviousState()
    
    -- Draw semi-transparent overlay
    rf.rectfill(0, 0, 480, 270, 0) -- Assuming color 0 with alpha
    
    -- Draw pause menu
    rf.print("PAUSED", 200, 100, 7)
  end
}
```

## Helper Functions

PICO-8-style utility functions for common operations.

### `rf.flr(x)`
Floor function: Truncates `x` towards zero (same as `math.floor` in Lua).

### `rf.ceil(x)`
Ceiling function: Rounds `x` up to the nearest integer.

### `rf.rnd([x])`
Random number generator (PICO-8-compatible LCG).
- `rf.rnd()` - Returns a random float between 0.0 and 1.0 (exclusive of 1.0)
- `rf.rnd(x)` - Returns a random float between 0.0 and `x` (exclusive of `x`)
- `rf.rnd(-x)` - Returns a random float between `-x` and 0.0 (exclusive of `-x`)
- Uses deterministic seeded random number generator (same seed = same sequence)

### `rf.mid(x, y, z)`
Clamp function: Returns `x` clamped between `y` and `z`.
- If `y > z`, they are automatically swapped
- Returns `y` if `x < y`, `z` if `x > z`, otherwise `x`

### `rf.sgn(x)`
Sign function: Returns the sign of `x`.
- Returns `1` if `x > 0`
- Returns `0` if `x == 0`
- Returns `-1` if `x < 0`

### `rf.chr(n)`
Convert number to character: Returns a single-character string with byte value `n`.
- `n` is clamped to 0-255 range
- Returns the byte as a single-character string

### `rf.ord(c)`
Convert character to number: Returns the byte value of the first character in string `c`.
- Returns `0` for empty strings
- Returns the first byte value (0-255) of the string

## Bitwise Operations

PICO-8-compatible bitwise functions for low-level bit manipulation.

### `rf.shl(x, y)`
Bitwise shift left: `x << y`. Returns `x` shifted left by `y` bits.
- Negative `y` performs shift right
- Shifts beyond 63 bits result in 0

### `rf.shr(x, y)`
Bitwise shift right (arithmetic): `x >> y`. Returns `x` shifted right by `y` bits with sign extension.
- Negative `y` performs shift left
- Shifts beyond 63 bits: 0 for positive numbers, -1 for negative (sign-extended)

### `rf.band(x, y)`
Bitwise AND: `x & y`. Returns the bitwise AND of `x` and `y`.

### `rf.bor(x, y)`
Bitwise OR: `x | y`. Returns the bitwise OR of `x` and `y`.

### `rf.bxor(x, y)`
Bitwise XOR: `x ^ y`. Returns the bitwise XOR (exclusive OR) of `x` and `y`.

### `rf.bnot(x)`
Bitwise NOT: `~x`. Returns the bitwise complement (one's complement) of `x`.

**Note:** All bitwise operations work with 64-bit signed integers. Floating point numbers are converted to integers.

## Debug Functions (Development Mode Only)

These functions are **only available when running from a folder** (not in packed .rf files).

### `rf.printh(str)`
Print debug message to console. Only works in development mode.

### `rf.stat(n)`
Get system statistics. Only works in development mode.
- `rf.stat(0)` - Current FPS
- `rf.stat(1)` - Frame count
- `rf.stat(2)` - Lua memory usage (bytes)
- `rf.stat(3)` - Last load time (milliseconds)
- `rf.stat(4)` - Last reload time (Unix timestamp)
- `rf.stat(5)` - Total reload count

### `rf.time()`
Get current time in seconds (Unix timestamp). Only works in development mode.

**Usage:** Run the engine with `--folder <path>` flag to enable development mode:
```bash
retroforge --folder examples/helloworld --window
```
When running from a folder, files are automatically reloaded when modified.

## Multiplayer API

RetroForge supports multiplayer games with up to 6 players via WebRTC networking. The engine handles all networking automatically - you just need to register tables for synchronization and check if you're the host.

### Connection Information

- `rf.is_multiplayer()` - Returns `true` if game is in multiplayer mode, `false` for solo play
- `rf.player_count()` - Returns total number of connected players (1-6). Always >= 1 (includes local player)
- `rf.my_player_id()` - Returns the local player's ID (1-6). Host is always player_id = 1
- `rf.is_host()` - Returns `true` if local player is the host. Only host runs game logic in `_update()`

### State Synchronization

- `rf.network_sync(table, tier)` - Register a Lua table for automatic synchronization
  - `table`: Lua table to synchronize (e.g., `players`, `score`)
  - `tier`: Synchronization frequency - `"fast"` (30-60/sec), `"moderate"` (15/sec), or `"slow"` (5/sec)
  - Example: `rf.network_sync(players, "fast")` - Sync player positions smoothly
  - Example: `rf.network_sync(score, "slow")` - Sync score updates less frequently
  - Only host can modify tables that aren't keyed by player_id
  - Tables keyed by player_id: only that player can modify their entry

- `rf.network_unsync(table)` - Unregister a table from synchronization
  - Stops syncing the table
  - Useful for cleanup or dynamic objects

### Input Handling (Host Only)

- `rf.btn(player_id, button)` - Check another player's button state (host only)
  - `player_id`: 1-6
  - `button`: 0-15 (button index)
  - Returns `true` if that player is pressing the button
  - Only works in host's code
  - Non-hosts use normal `rf.btn(button)` for local player

**Example:**
```lua
function _update()
  if rf.is_host() then
    -- Host can check all players' inputs
    for id = 1, rf.player_count() do
      if rf.btn(id, 0) then  -- Check if Player id is pressing left
        players[id].vx = -2
      end
    end
  else
    -- Non-host: just use normal input
    if rf.btn(0) then
      -- Local input handling (if needed)
    end
  end
end
```

### Multiplayer Example

```lua
-- Simple multiplayer platformer
players = {}
score = {}

function _init()
  local player_count = rf.is_multiplayer() and rf.player_count() or 1
  
  -- Create player for each connected player
  for i = 1, player_count do
    players[i] = {x = 100*i, y = 100, vx = 0, vy = 0}
    score[i] = 0
  end
  
  -- Register for sync (multiplayer only)
  if rf.is_multiplayer() then
    rf.network_sync(players, "fast")   -- Smooth movement
    rf.network_sync(score, "slow")     -- Less frequent updates
  end
end

function _update()
  if rf.is_multiplayer() and rf.is_host() then
    -- Host runs game logic
    for id = 1, rf.player_count() do
      if rf.btn(id, 0) then players[id].vx = -3 end
      if rf.btn(id, 1) then players[id].vx = 3 end
      players[id].x = players[id].x + players[id].vx
    end
    -- Engine automatically syncs players table!
  else
    -- Solo mode or non-host (engine handles input sync automatically)
  end
end

function _draw()
  rf.clear_i(0)
  for id, p in pairs(players) do
    rf.rectfill(p.x, p.y, p.x+16, p.y+16, 7)
    rf.print("P" .. id, p.x, p.y-10, 7)
  end
end
```

**Note:** Multiplayer functionality requires WebRTC connections which are established through the RetroForge webapp platform. See `RetroForge.V2.md` for complete multiplayer architecture details.

