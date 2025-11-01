# RetroForge Module-Based State System

**Version:** 1.0  
**Date:** October 31, 2025  
**Status:** Design Complete - Ready for Implementation

## Overview

The Module-Based State System extends RetroForge's existing GameStateMachine with a convention-based approach for defining game states. Instead of manually creating state tables and registering them, developers can create separate `.lua` files that are automatically loaded and registered as states.

### Motivation

**Traditional State Registration (Current System):**
```lua
-- main.lua
local MenuState = {
  selectedOption = 1,
  
  initialize = function(sm)
    -- Setup
  end,
  
  enter = function(sm)
    selectedOption = 1
  end,
  
  handleInput = function(sm)
    -- Input handling
  end,
  
  update = function(dt)
    -- Update logic
  end,
  
  draw = function()
    -- Rendering
  end,
  
  exit = function(sm)
    -- Cleanup
  end
}

game.registerState("menu", MenuState)
```

**Module-Based System (New):**
```lua
-- main.lua
rf.import("menu_state.lua")  -- Auto-registers as "menu" state
game.changeState("menu")

-- menu_state.lua (separate file)
local selectedOption = 1

function _INIT()
  -- Setup
end

function _ENTER()
  selectedOption = 1
end

function _HANDLE_INPUT()
  -- Input handling
end

function _UPDATE(dt)
  -- Update logic
end

function _DRAW()
  -- Rendering
end

function _EXIT()
  -- Cleanup
end

function _DONE()
  -- Final teardown
end
```

### Benefits

1. **Separation of Concerns**: Each state lives in its own file
2. **Convention Over Configuration**: Standardized function names
3. **Automatic Registration**: No manual `game.registerState()` calls
4. **Module Persistence**: Module-level variables persist across enter/exit cycles
5. **Shared Context**: All states access a common context object
6. **Cleaner main.lua**: Entry point stays minimal and focused

---

## Technical Specification

### File Naming Convention

State modules use the following naming patterns:

| Filename | Registered State Name |
|----------|----------------------|
| `menu_state.lua` | `"menu"` |
| `playing_state.lua` | `"playing"` |
| `game_over_state.lua` | `"game_over"` |
| `pause.lua` | `"pause"` |
| `shop.lua` | `"shop"` |

**Pattern:** `{name}_state.lua` or `{name}.lua` → state name is `{name}`

### Required Functions

Every state module **must** implement these five functions:

```lua
function _INIT()
  -- Called once when state is first created
  -- Use for: Loading assets, one-time setup, resource allocation
end

function _UPDATE(dt)
  -- Called every frame while state is active
  -- Parameter: dt (delta time in seconds)
  -- Use for: Game logic, animations, timers
end

function _DRAW()
  -- Called every frame while state is active
  -- Use for: Rendering graphics, UI, effects
end

function _HANDLE_INPUT()
  -- Called every frame while state is active (before update)
  -- Use for: Processing player input, button presses
end

function _DONE()
  -- Called once when state is destroyed
  -- Use for: Final cleanup, saving data, freeing resources
end
```

### Optional Functions

State modules **may** implement these optional functions:

```lua
function _ENTER()
  -- Called every time state becomes active
  -- Use for: Resetting per-session state, starting music, showing UI
end

function _EXIT()
  -- Called every time state becomes inactive
  -- Use for: Pausing, saving temporary data, stopping sounds
end
```

### Function Lifecycle

```
State Created
    │
    └──> _INIT() ────────────────┐ (Once)
                                  │
    ┌─────────────────────────────┘
    │
    └──> _ENTER() ────────────────┐ (Every activation)
                                  │
    ┌─────────────────────────────┘
    │
    ├──> _HANDLE_INPUT() ─────────┤
    ├──> _UPDATE(dt) ─────────────┤ (Every frame while active)
    └──> _DRAW() ─────────────────┤
                                  │
    ┌─────────────────────────────┘
    │
    └──> _EXIT() ─────────────────┐ (Every deactivation)
                                  │
    ┌─────────────────────────────┘
    │
    └──> _DONE() ─────────────────┐ (Once on destroy)
                                  
State Destroyed
```

**Example Flow:**
1. `menu_state.lua` imported → `_INIT()` called
2. `game.changeState("menu")` → `_ENTER()` called
3. Game loop runs → `_HANDLE_INPUT()`, `_UPDATE(dt)`, `_DRAW()` called each frame
4. `game.changeState("playing")` → `_EXIT()` called on menu, `_ENTER()` called on playing
5. State unregistered → `_DONE()` called

---

## API Reference

### Core Functions

#### `rf.import(filename)`

Loads a Lua module file and automatically registers it as a game state.

**Parameters:**
- `filename` (string): Path to the state module file (relative to cart root)

**Returns:**
- (string): The registered state name

**Behavior:**
1. Loads the file in an isolated environment
2. Validates that all required functions are defined
3. Extracts state name from filename
4. Wraps module functions in state machine callbacks
5. Calls `game.registerState()` automatically

**Example:**
```lua
-- main.lua
function _init()
  rf.import("menu_state.lua")      -- Registers "menu"
  rf.import("playing_state.lua")   -- Registers "playing"
  rf.import("game_over.lua")       -- Registers "game_over"
  
  game.changeState("menu")
end
```

**Error Cases:**
- File not found → Runtime error with filename
- Missing required function → Runtime error listing missing function
- Syntax error in module → Runtime error with line number

### Context Object

#### `context`

A global table accessible to all state modules for sharing data.

**Usage:**
```lua
-- main.lua
context = {
  player = {x = 100, y = 100, lives = 3},
  score = 0,
  level = 1,
  settings = {volume = 0.8, difficulty = "normal"}
}

-- Any state module can access:
-- playing_state.lua
function _UPDATE(dt)
  context.player.x = context.player.x + 2
  context.score = context.score + 1
end

-- game_over_state.lua  
function _DRAW()
  rf.print("Final Score: " .. context.score, 100, 100, 7)
end
```

**Best Practices:**
- Define structure in `main.lua` before importing states
- Use nested tables for organization (player, settings, game data)
- Don't store per-state temporary data here (use module-level variables)

### State Machine Access

All state modules have access to the GameStateMachine API:

```lua
-- Inside any state module
function _HANDLE_INPUT()
  if rf.btnp(4) then
    game.changeState("playing")     -- Switch to new state
    game.pushState("pause")         -- Overlay state
    game.popState()                 -- Return to previous
    game.getState()                 -- Get current state name
    game.getStackDepth()            -- Get stack size
  end
end
```

See main API_REFERENCE.md for complete GameStateMachine documentation.

---

## Module Environment

### Scope and Persistence

State modules execute in an isolated environment with special properties:

```lua
-- menu_state.lua

-- MODULE SCOPE (persists across enter/exit)
local selected_option = 1
local menu_items = {"PLAY", "OPTIONS", "QUIT"}
local animation_timer = 0

function _ENTER()
  -- Reset per-session state
  selected_option = 1
  animation_timer = 0
  
  -- But menu_items persists (not reset)
end

function _UPDATE(dt)
  -- Module-level variables persist
  animation_timer = animation_timer + dt
end
```

**What persists:**
- Variables declared at module level (outside functions)
- Function definitions
- Tables and nested data structures

**What resets:**
- Variables assigned in `_ENTER()`
- `context` object (shared across all states)

### Variable Resolution Order

When accessing a variable, the engine checks in this order:

1. **Module scope** - Variables defined in the state file
2. **Context** - `context.player`, `context.score`, etc.
3. **Globals** - `rf.*`, `game.*`, built-in Lua functions

```lua
-- menu_state.lua
local selected = 1  -- Module scope (highest priority)

function _UPDATE(dt)
  selected = selected + 1        -- Uses module scope
  context.score = context.score + 1  -- Uses context
  rf.print("Hello", 0, 0, 7)     -- Uses global rf
end
```

### State Machine Reference

Every state module automatically receives:

```lua
sm        -- State machine instance (available in callbacks)
context   -- Alias to sm.context (convenient access)
```

**Usage:**
```lua
function _INIT()
  -- sm is available
  rf.printh("State machine has " .. game.getStackDepth() .. " states")
end

function _ENTER()
  -- context is available
  context.music_volume = 0.8
end
```

---

## Implementation Guide

### Basic State Module Template

```lua
-- state_template.lua

-- Module-level state (persists)
local local_var = 0

function _INIT()
  -- One-time setup
  -- Load assets, allocate resources
end

function _ENTER()
  -- Called every activation
  -- Reset temporary state
end

function _HANDLE_INPUT()
  -- Process input every frame
  -- Check button presses
end

function _UPDATE(dt)
  -- Update game logic every frame
  -- dt = delta time in seconds
end

function _DRAW()
  -- Render graphics every frame
  -- Draw sprites, UI, effects
end

function _EXIT()
  -- Called every deactivation
  -- Pause, save temp data
end

function _DONE()
  -- One-time teardown
  -- Save data, free resources
end
```

### main.lua Structure

```lua
-- main.lua - Entry point

-- Global context (shared across all states)
context = {
  player = {},
  game_data = {},
  settings = {}
}

function _init()
  -- Load saved data
  load_save_data()
  
  -- Import all state modules
  rf.import("splash_state.lua")
  rf.import("menu_state.lua")
  rf.import("playing_state.lua")
  rf.import("pause_state.lua")
  rf.import("game_over_state.lua")
  rf.import("settings_state.lua")
  
  -- Start at splash screen
  game.changeState("splash")
end

function load_save_data()
  context.settings.volume = rf.peek(0) / 255
  context.high_score = rf.peek4(4)
end

-- Note: main.lua doesn't need _update, _draw, etc.
-- The state machine handles everything through active states
```

---

## Usage Examples

### Example 1: Simple Menu State

```lua
-- menu_state.lua

-- Persistent module state
local menu_items = {"START GAME", "OPTIONS", "CREDITS", "QUIT"}
local selected = 1
local blink_timer = 0

function _INIT()
  rf.printh("Menu initialized")
end

function _ENTER()
  -- Reset selection when returning to menu
  selected = 1
  blink_timer = 0
  rf.music("menu_theme")
end

function _HANDLE_INPUT()
  -- Navigate menu
  if rf.btnp(2) then  -- Up
    selected = selected - 1
    if selected < 1 then selected = #menu_items end
    rf.sfx("cursor")
  elseif rf.btnp(3) then  -- Down
    selected = selected + 1
    if selected > #menu_items then selected = 1 end
    rf.sfx("cursor")
  elseif rf.btnp(4) then  -- A button
    rf.sfx("select")
    
    if selected == 1 then
      game.changeState("playing")
    elseif selected == 2 then
      game.pushState("options")
    elseif selected == 3 then
      game.pushState("credits")
    elseif selected == 4 then
      rf.quit()
    end
  end
end

function _UPDATE(dt)
  blink_timer = blink_timer + dt
end

function _DRAW()
  rf.clear_i(0)
  
  -- Title
  rf.print_anchored("RETRO QUEST", "topcenter", 11)
  
  -- Menu items
  for i, item in ipairs(menu_items) do
    local y = 100 + (i - 1) * 25
    local color = 7
    local prefix = "  "
    
    if i == selected then
      -- Blinking cursor
      if blink_timer % 1.0 < 0.5 then
        prefix = "> "
      else
        prefix = "  "
      end
      color = 11
    end
    
    rf.print_xy(180, y, prefix .. item, color)
  end
  
  -- Footer
  rf.print_xy(10, 260, "v1.0", 5)
end

function _EXIT()
  rf.music("stopall")
end

function _DONE()
  rf.printh("Menu destroyed")
end
```

### Example 2: Playing State with Game Logic

```lua
-- playing_state.lua

-- Persistent game state
local enemies = {}
local bullets = {}
local spawn_timer = 0
local wave = 1

function _INIT()
  rf.printh("Playing state initialized")
end

function _ENTER()
  -- Reset game for new play session
  context.player = {
    x = 240,
    y = 200,
    lives = 3,
    invincible = 0
  }
  context.score = 0
  
  enemies = {}
  bullets = {}
  spawn_timer = 0
  wave = 1
  
  rf.music("gameplay")
end

function _HANDLE_INPUT()
  local p = context.player
  
  -- Movement
  if rf.btn(0) then p.x = p.x - 3 end  -- Left
  if rf.btn(1) then p.x = p.x + 3 end  -- Right
  if rf.btn(2) then p.y = p.y - 3 end  -- Up
  if rf.btn(3) then p.y = p.y + 3 end  -- Down
  
  -- Clamp to screen
  p.x = rf.mid(p.x, 8, 472)
  p.y = rf.mid(p.y, 8, 262)
  
  -- Shoot
  if rf.btnp(4) then
    table.insert(bullets, {x = p.x, y = p.y - 10, vy = -5})
    rf.sfx("shoot")
  end
  
  -- Pause
  if rf.btnp(5) then
    game.pushState("pause")
  end
end

function _UPDATE(dt)
  local p = context.player
  
  -- Update invincibility
  if p.invincible > 0 then
    p.invincible = p.invincible - dt
  end
  
  -- Update bullets
  for i = #bullets, 1, -1 do
    local b = bullets[i]
    b.y = b.y + b.vy
    
    -- Remove off-screen bullets
    if b.y < 0 then
      table.remove(bullets, i)
    end
  end
  
  -- Update enemies
  for i = #enemies, 1, -1 do
    local e = enemies[i]
    e.y = e.y + e.vy
    
    -- Remove off-screen enemies
    if e.y > 270 then
      table.remove(enemies, i)
    end
    
    -- Check collision with bullets
    for j = #bullets, 1, -1 do
      local b = bullets[j]
      if check_collision(e, b, 12) then
        table.remove(enemies, i)
        table.remove(bullets, j)
        context.score = context.score + 10
        rf.sfx("explosion")
        break
      end
    end
    
    -- Check collision with player
    if p.invincible <= 0 and check_collision(e, p, 12) then
      p.lives = p.lives - 1
      p.invincible = 2.0
      table.remove(enemies, i)
      rf.sfx("hurt")
      
      if p.lives <= 0 then
        game.changeState("game_over")
      end
    end
  end
  
  -- Spawn enemies
  spawn_timer = spawn_timer + dt
  local spawn_rate = math.max(0.5, 2.0 - wave * 0.1)
  
  if spawn_timer > spawn_rate then
    spawn_timer = 0
    table.insert(enemies, {
      x = rf.rnd(460) + 10,
      y = -10,
      vy = 1 + wave * 0.2
    })
  end
  
  -- Wave progression
  if context.score > wave * 100 then
    wave = wave + 1
    rf.sfx("wave_complete")
  end
end

function _DRAW()
  rf.clear_i(1)
  
  local p = context.player
  
  -- Draw player (flashing when invincible)
  if p.invincible <= 0 or (p.invincible * 10) % 2 < 1 then
    rf.spr("player", p.x - 8, p.y - 8)
  end
  
  -- Draw bullets
  for _, b in ipairs(bullets) do
    rf.rectfill(b.x - 2, b.y - 4, b.x + 2, b.y + 4, 11)
  end
  
  -- Draw enemies
  for _, e in ipairs(enemies) do
    rf.spr("enemy", e.x - 8, e.y - 8)
  end
  
  -- HUD
  rf.print_xy(10, 10, "SCORE: " .. context.score, 7)
  rf.print_xy(10, 20, "LIVES: " .. p.lives, 7)
  rf.print_xy(10, 30, "WAVE: " .. wave, 7)
end

function _EXIT()
  -- Pause music (for pause screen overlay)
end

function _DONE()
  -- Save high score
  if context.score > context.high_score then
    context.high_score = context.score
    rf.poke4(4, context.high_score)
  end
end

-- Helper function (module scope)
function check_collision(a, b, radius)
  local dx = a.x - b.x
  local dy = a.y - b.y
  return (dx * dx + dy * dy) < (radius * radius)
end
```

### Example 3: Pause State (Overlay)

```lua
-- pause_state.lua

local options = {"RESUME", "RESTART", "QUIT TO MENU"}
local selected = 1

function _INIT()
end

function _ENTER()
  selected = 1
  rf.music("stopall")  -- Pause game music
end

function _HANDLE_INPUT()
  if rf.btnp(2) then  -- Up
    selected = selected - 1
    if selected < 1 then selected = #options end
    rf.sfx("cursor")
  elseif rf.btnp(3) then  -- Down
    selected = selected + 1
    if selected > #options then selected = 1 end
    rf.sfx("cursor")
  elseif rf.btnp(4) then  -- A button
    rf.sfx("select")
    
    if selected == 1 then
      game.popState()  -- Resume game
    elseif selected == 2 then
      game.popState()  -- Pop pause
      game.changeState("playing")  -- Restart
    elseif selected == 3 then
      game.changeState("menu")  -- Back to menu
    end
  elseif rf.btnp(5) then  -- Pause button
    game.popState()  -- Resume
  end
end

function _UPDATE(dt)
  -- Nothing to update
end

function _DRAW()
  -- Draw previous state (dimmed game)
  game.drawPreviousState()
  
  -- Semi-transparent overlay
  rf.rectfill(0, 0, 480, 270, 0)  -- Assuming dark color
  
  -- Pause menu
  rf.print_anchored("PAUSED", "topcenter", 7)
  
  for i, option in ipairs(options) do
    local y = 120 + (i - 1) * 25
    local prefix = (i == selected) and "> " or "  "
    local color = (i == selected) and 11 or 7
    rf.print_xy(180, y, prefix .. option, color)
  end
  
  -- Instructions
  rf.print_anchored("Press START to resume", "bottomcenter", 5)
end

function _EXIT()
  -- Resume game music
  rf.music("gameplay")
end

function _DONE()
end
```

### Example 4: Game Over State

```lua
-- game_over_state.lua

local wait_timer = 0
local new_high_score = false

function _INIT()
end

function _ENTER()
  wait_timer = 0
  
  -- Check for new high score
  if context.score > context.high_score then
    context.high_score = context.score
    new_high_score = true
    rf.poke4(4, context.high_score)
    rf.sfx("high_score")
  else
    new_high_score = false
  end
  
  rf.music("game_over")
end

function _HANDLE_INPUT()
  -- Wait 2 seconds before allowing input
  if wait_timer > 2.0 then
    if rf.btnp(4) then  -- A button
      game.changeState("menu")
    end
  end
end

function _UPDATE(dt)
  wait_timer = wait_timer + dt
end

function _DRAW()
  rf.clear_i(0)
  
  -- Game Over text
  rf.print_anchored("GAME OVER", "topcenter", 8)
  
  -- Scores
  rf.print_anchored("YOUR SCORE", "middlecenter", 7)
  rf.print_anchored(tostring(context.score), "middlecenter", 11)
  rf.cursor(240, 160)
  
  if new_high_score then
    rf.print_anchored("NEW HIGH SCORE!", "middlecenter", 11)
    rf.cursor(240, 180)
  else
    rf.print_anchored("HIGH SCORE: " .. context.high_score, "middlecenter", 5)
    rf.cursor(240, 180)
  end
  
  -- Continue prompt (after delay)
  if wait_timer > 2.0 then
    if wait_timer % 1.0 < 0.5 then
      rf.print_anchored("PRESS A TO CONTINUE", "bottomcenter", 7)
    end
  end
end

function _EXIT()
  rf.music("stopall")
end

function _DONE()
end
```

---

## Best Practices

### Module Organization

**DO:**
```lua
-- playing_state.lua

-- Group related module state at the top
local player_state = {x = 0, y = 0, vx = 0, vy = 0}
local enemies = {}
local powerups = {}

-- Helper functions after state
function spawn_enemy()
  table.insert(enemies, {x = rf.rnd(480), y = 0})
end

function check_powerup_collision()
  -- Logic here
end

-- Lifecycle functions at the bottom
function _INIT()
  -- ...
end

function _UPDATE(dt)
  -- ...
end
```

**DON'T:**
```lua
-- Mixing lifecycle and helpers randomly
function _UPDATE(dt)
  -- ...
end

function spawn_enemy()
  -- ...
end

function _INIT()
  -- ...
end

function _DRAW()
  -- ...
end
```

### State Transitions

**DO - Clean transitions:**
```lua
function _HANDLE_INPUT()
  if rf.btnp(4) then
    -- Save any necessary state to context
    context.previous_state = "menu"
    
    -- Transition
    game.changeState("playing")
  end
end
```

**DON'T - Transition during critical operations:**
```lua
function _UPDATE(dt)
  for i = #enemies, 1, -1 do
    -- BAD: State transition during enemy loop
    if enemies[i].health <= 0 then
      game.changeState("game_over")  -- DON'T DO THIS
      table.remove(enemies, i)
    end
  end
end

-- BETTER: Set a flag, transition in _HANDLE_INPUT
local game_over_flag = false

function _UPDATE(dt)
  for i = #enemies, 1, -1 do
    if enemies[i].health <= 0 then
      game_over_flag = true
    end
  end
end

function _HANDLE_INPUT()
  if game_over_flag then
    game.changeState("game_over")
  end
end
```

### Context Usage

**DO - Organize context in main.lua:**
```lua
-- main.lua
context = {
  player = {x = 0, y = 0, lives = 3, score = 0},
  game = {level = 1, difficulty = "normal"},
  settings = {volume = 0.8, particles = true},
  persistent = {high_score = 0, unlocks = {}}
}
```

**DON'T - Create context structure in states:**
```lua
-- playing_state.lua
function _INIT()
  -- DON'T: Context should be defined in main.lua
  context.player = {}
  context.enemies = {}
end
```

### Module State Persistence

**DO - Understand what persists:**
```lua
-- Module-level state persists
local total_playthroughs = 0
local best_time = 999999

function _ENTER()
  -- This increments across multiple play sessions
  total_playthroughs = total_playthroughs + 1
  
  -- Per-session state
  local session_time = 0  -- This does NOT persist
end
```

**DON'T - Assume _ENTER resets everything:**
```lua
local score = 0

function _ENTER()
  -- This does NOT reset score to 0 automatically
  -- You must explicitly reset it
  score = 0
end
```

### Error Handling

**DO - Validate critical data:**
```lua
function _INIT()
  -- Validate required assets exist
  if not rf.sprite("player") then
    rf.printh("ERROR: Missing player sprite!")
  end
end

function _UPDATE(dt)
  -- Guard against nil values
  if context.player and context.player.lives > 0 then
    -- Update logic
  end
end
```

**DON'T - Assume everything is valid:**
```lua
function _UPDATE(dt)
  -- BAD: Will crash if context.player is nil
  context.player.x = context.player.x + 1
end
```

---

## Advanced Patterns

### State Communication

States can pass data through context:

```lua
-- menu_state.lua
function _HANDLE_INPUT()
  if rf.btnp(4) then
    context.selected_difficulty = selected
    game.changeState("playing")
  end
end

-- playing_state.lua
function _ENTER()
  local difficulty = context.selected_difficulty or 1
  enemy_speed = difficulty * 2
end
```

### Substate Pattern

Use module-level functions for complex state logic:

```lua
-- playing_state.lua
local current_phase = "intro"

function _UPDATE(dt)
  if current_phase == "intro" then
    update_intro(dt)
  elseif current_phase == "gameplay" then
    update_gameplay(dt)
  elseif current_phase == "boss" then
    update_boss(dt)
  end
end

function update_intro(dt)
  -- Intro logic
  if intro_complete then
    current_phase = "gameplay"
  end
end

function update_gameplay(dt)
  -- Normal gameplay
  if boss_triggered then
    current_phase = "boss"
  end
end

function update_boss(dt)
  -- Boss fight logic
end
```

### Shared Utilities

Create a utilities module (not a state):

```lua
-- utils.lua (not imported as state)
function distance(x1, y1, x2, y2)
  local dx = x2 - x1
  local dy = y2 - y1
  return math.sqrt(dx * dx + dy * dy)
end

function circle_collision(a, b, radius)
  return distance(a.x, a.y, b.x, b.y) < radius
end

-- Load in main.lua
dofile("utils.lua")

-- Use in any state
function _UPDATE(dt)
  if circle_collision(player, enemy, 16) then
    -- Handle collision
  end
end
```

---

## Debugging and Development

### Debug Output

```lua
function _INIT()
  rf.printh("State initialized: " .. debug.getinfo(1).source)
end

function _UPDATE(dt)
  if rf.btn(8) then  -- Debug button
    rf.printh("Player: " .. context.player.x .. ", " .. context.player.y)
    rf.printh("Score: " .. context.score)
  end
end
```

### State Tracing

```lua
-- Add to each state for debugging
local state_name = "menu"  -- Set this in each module

function _ENTER()
  rf.printh("ENTER: " .. state_name)
end

function _EXIT()
  rf.printh("EXIT: " .. state_name)
end
```

### Runtime Statistics

```lua
function _DRAW()
  if rf.stat then  -- Development mode only
    rf.print_xy(10, 250, "FPS: " .. rf.stat(0), 7)
    rf.print_xy(10, 260, "State: " .. game.getState(), 7)
  end
end
```

---

## Error Reference

### Common Errors

**1. Missing Required Function**
```
ERROR: menu_state.lua must define function _UPDATE()
```
**Solution:** Ensure all five required functions are defined, even if empty.

**2. File Not Found**
```
ERROR: Failed to load menu_state.lua: file not found
```
**Solution:** Check filename spelling and that file is in cart root directory.

**3. Syntax Error**
```
ERROR: Failed to execute playing_state.lua: unexpected symbol near 'end'
```
**Solution:** Fix Lua syntax error in the state file.

**4. Context Not Defined**
```
ERROR: attempt to index global 'context' (a nil value)
```
**Solution:** Define `context = {}` in main.lua before importing states.

**5. Invalid State Name**
```
ERROR: State 'manu' not found (did you mean 'menu'?)
```
**Solution:** Check spelling in `game.changeState()` call.

---

## Performance Considerations

### Module Loading

- Modules are loaded once when `rf.import()` is called
- Module code is compiled and cached
- Minimal performance impact after initial load

### Memory Usage

- Module-level variables persist for the entire cart lifetime
- Use `_DONE()` to clean up large tables and release memory
- Consider clearing unused data in `_EXIT()` for memory-intensive states

### Frame Budget

Each state should aim for:
- `_HANDLE_INPUT()` < 1ms
- `_UPDATE()` < 8ms
- `_DRAW()` < 6ms

Use `rf.stat(0)` to monitor frame rate in development mode.

---

## Migration Guide

### Converting Existing State Tables

**Before:**
```lua
local MenuState = {
  initialize = function(sm)
    -- Setup
  end,
  
  update = function(dt)
    -- Logic
  end,
  
  draw = function()
    -- Render
  end
}

game.registerState("menu", MenuState)
```

**After:**
```lua
-- menu_state.lua

function _INIT()
  -- Setup
end

function _UPDATE(dt)
  -- Logic
end

function _DRAW()
  -- Render
end

function _HANDLE_INPUT()
  -- Add input handling
end

function _DONE()
  -- Add cleanup
end

-- In main.lua
rf.import("menu_state.lua")
```

### Handling Module State

**Before:**
```lua
local MenuState = {
  selectedOption = 1,  -- Table-level state
  
  enter = function(sm)
    MenuState.selectedOption = 1
  end
}
```

**After:**
```lua
-- menu_state.lua
local selected = 1  -- Module-level state (cleaner)

function _ENTER()
  selected = 1
end
```

---

## Implementation Notes

### Engine Requirements

The engine must implement `rf.import()` with:

1. **Module Environment Creation**
   - Isolated Lua environment per module
   - Inherit `rf.*` and `game.*` APIs
   - Access to `context` table

2. **Function Validation**
   - Verify all required functions exist
   - Provide helpful error messages

3. **State Name Extraction**
   - Strip `_state.lua` suffix
   - Strip `.lua` extension
   - Use remaining string as state name

4. **Automatic Registration**
   - Wrap module functions in state machine callbacks
   - Call `game.registerState()` internally

### Lua Implementation Pseudocode

```lua
function rf.import(filename)
  -- Create module environment
  local module_env = {}
  setmetatable(module_env, {
    __index = function(t, key)
      return context[key] or _G[key]
    end
  })
  
  -- Load file
  local chunk = loadfile(filename, "t", module_env)
  pcall(chunk)
  
  -- Validate
  for _, func in ipairs({"_INIT", "_UPDATE", "_DRAW", "_HANDLE_INPUT", "_DONE"}) do
    assert(type(module_env[func]) == "function", 
           filename .. " must define " .. func .. "()")
  end
  
  -- Extract state name
  local state_name = filename:gsub("_state%.lua$", ""):gsub("%.lua$", "")
  
  -- Create and register state
  game.registerState(state_name, {
    initialize = module_env._INIT,
    enter = module_env._ENTER,
    handleInput = function(sm)
      module_env.sm = sm
      module_env.context = sm.context
      return module_env._HANDLE_INPUT()
    end,
    update = module_env._UPDATE,
    draw = module_env._DRAW,
    exit = module_env._EXIT,
    shutdown = module_env._DONE
  })
  
  return state_name
end
```

---

## Future Enhancements

### Potential Features

1. **Hot Reloading**
   - Reload state modules without restarting cart
   - Preserve module state across reloads
   - Development mode only

2. **State Parameters**
   ```lua
   game.changeState("level", {level_id = 3, difficulty = "hard"})
   
   -- level_state.lua
   function _ENTER()
     local params = sm.params
     load_level(params.level_id)
   end
   ```

3. **State History**
   ```lua
   -- Access previous state name
   function _ENTER()
     local previous = sm.previous_state
     if previous == "pause" then
       -- Resumed from pause
     end
   end
   ```

4. **Async State Loading**
   ```lua
   rf.import_async("large_level.lua", function()
     game.changeState("large_level")
   end)
   ```

5. **State Groups**
   ```lua
   rf.import_folder("states/")  -- Import all .lua files in folder
   ```

---

## Conclusion

The Module-Based State System provides a clean, intuitive way to organize game states in RetroForge. By following convention over configuration, developers can focus on game logic rather than boilerplate code.

### Key Takeaways

- ✅ One file per state for clean organization
- ✅ Convention-based function names (_INIT, _UPDATE, etc.)
- ✅ Automatic registration with `rf.import()`
- ✅ Shared context for cross-state communication
- ✅ Module-level persistence for state data
- ✅ Full access to GameStateMachine API

### Getting Started

1. Define `context` in main.lua
2. Create state files with required functions
3. Import states with `rf.import()`
4. Start game with `game.changeState()`

Happy game development! 🎮

---

**Document Version:** 1.0  
**Last Updated:** October 31, 2025  
**Status:** Ready for Implementation
