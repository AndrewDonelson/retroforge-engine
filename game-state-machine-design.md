# Game State Machine Design

## Overview

A flexible state machine for a 2D game engine written in Go with Lua bindings. The state machine manages game flow through different states, supports both simple state transitions and state stacking for overlays, includes two built-in states for engine branding and credits, and provides an optional shared context for data passing between states.

## Core Concepts

### State Lifecycle

Every state follows a complete lifecycle:

1. **Initialize** - Called once when state is first created (load assets, allocate resources)
2. **Enter** - Called every time the state becomes active (reset variables, start music)
3. **HandleInput** - Process user input while active
4. **Update** - Update game logic each frame while active
5. **Draw** - Render the state
6. **Exit** - Called when leaving the state (pause music, but keep in memory)
7. **Shutdown** - Called once when state is destroyed (free assets, final cleanup)

### State Management Modes

**ChangeState** - Complete state transition
- Old state exits and is removed from the stack
- New state becomes active
- Use for: menu → playing, playing → game over

**Push/Pop** - State stacking for overlays
- Push: Add new state on top, previous state pauses but stays in memory
- Pop: Remove current state, return to previous state
- Use for: playing → pause menu, playing → inventory

### Shared Context

An optional key-value store that allows states to share data without being tightly coupled:
- **Optional** - Only use when you need inter-state communication
- **Persistent** - Data survives state transitions
- **Flexible** - Store any type of data (strings, numbers, tables, booleans)
- **Use cases** - Pass level selection, share player stats, store global settings

## Architecture

### Go Side - Core Components

```go
type State interface {
    Initialize(sm *StateMachine) error
    Enter(sm *StateMachine)
    HandleInput(sm *StateMachine)
    Update(dt float64)
    Draw()
    Exit(sm *StateMachine)
    Shutdown()
}

type StateMachine struct {
    stateStack       []*State           // Stack of active states
    stateRegistry    map[string]*State  // All registered states
    initialized      map[string]bool    // Track which states are initialized
    context          map[string]interface{} // Shared context for data passing
    isDebug          bool
    shouldExit       bool
    
    // Built-in states
    engineSplash     *EngineSplashState
    credits          *CreditsState
}
```

### State Stack Visualization

```
Single State (ChangeState):
┌─────────────┐
│ PlayingState│ ← Active
└─────────────┘

Stacked States (PushState):
┌─────────────┐
│ PausedState │ ← Active (top of stack)
├─────────────┤
│ PlayingState│ ← Paused (still in memory)
└─────────────┘

Shared Context (accessible to all states):
┌─────────────────────────────┐
│ context: {                  │
│   "player_health": 100,     │
│   "current_level": 3,       │
│   "difficulty": "hard"      │
│ }                           │
└─────────────────────────────┘
```

## Built-in States

### 1. Engine Splash State

**Purpose:** Display engine branding before any developer content

**Behavior:**
- First state shown in release builds
- Skipped entirely in debug builds
- Auto-transitions after 2-3 seconds
- Displays engine logo/name

**Lifecycle:**
```
Release Build: __engine_splash → developer_splash → main_menu
Debug Build:   (skipped) → main_menu
```

### 2. Credits State

**Purpose:** Display credits before game exit

**Features:**
- Built-in engine credits (engine name, version, developer)
- User-added credits (game developers, artists, sound, etc.)
- Auto-scrolling display
- Exit on any input (keyboard, mouse, controller)

**API for Adding Credits:**
```lua
game.addCredit(category, name, role)
```

**Categories:** "Developer", "Artist", "Sound", "Music", "Special Thanks", etc.

## State Machine Methods

### Core Methods

```go
// Initialization
func NewStateMachine(isDebugBuild bool) *StateMachine

// State Registration
func (sm *StateMachine) RegisterState(name string, luaCallbacks LuaCallbacks)
func (sm *StateMachine) UnregisterState(name string)

// Simple Transitions
func (sm *StateMachine) ChangeState(name string)

// State Stacking
func (sm *StateMachine) PushState(name string)
func (sm *StateMachine) PopState()
func (sm *StateMachine) PopAllStates()

// Main Loop (called by engine)
func (sm *StateMachine) HandleInput()
func (sm *StateMachine) Update(dt float64)
func (sm *StateMachine) Draw()

// Shared Context
func (sm *StateMachine) SetContext(key string, value interface{})
func (sm *StateMachine) GetContext(key string) (interface{}, bool)
func (sm *StateMachine) HasContext(key string) bool
func (sm *StateMachine) ClearContext(key string)
func (sm *StateMachine) ClearAllContext()

// Credits API
func (sm *StateMachine) AddCreditEntry(category, name, role string)

// Control
func (sm *StateMachine) RequestExit()
func (sm *StateMachine) ShouldExit() bool
```

### Implementation Details

#### ChangeState
```go
func (sm *StateMachine) ChangeState(name string) {
    // Exit and remove all states from stack
    for len(sm.stateStack) > 0 {
        top := sm.stateStack[len(sm.stateStack)-1]
        top.Exit(sm)
        sm.stateStack = sm.stateStack[:len(sm.stateStack)-1]
    }
    
    // Get new state
    newState := sm.stateRegistry[name]
    
    // Initialize if first time
    if !sm.initialized[name] {
        newState.Initialize(sm)
        sm.initialized[name] = true
    }
    
    // Add to stack and enter
    sm.stateStack = append(sm.stateStack, newState)
    newState.Enter(sm)
}
```

#### PushState
```go
func (sm *StateMachine) PushState(name string) {
    // Exit current top state (but keep on stack)
    if len(sm.stateStack) > 0 {
        top := sm.stateStack[len(sm.stateStack)-1]
        top.Exit(sm)
    }
    
    // Get new state
    newState := sm.stateRegistry[name]
    
    // Initialize if first time
    if !sm.initialized[name] {
        newState.Initialize(sm)
        sm.initialized[name] = true
    }
    
    // Add to stack and enter
    sm.stateStack = append(sm.stateStack, newState)
    newState.Enter(sm)
}
```

#### PopState
```go
func (sm *StateMachine) PopState() {
    if len(sm.stateStack) == 0 {
        return
    }
    
    // Exit and remove top state
    top := sm.stateStack[len(sm.stateStack)-1]
    top.Exit(sm)
    sm.stateStack = sm.stateStack[:len(sm.stateStack)-1]
    
    // Re-enter previous state
    if len(sm.stateStack) > 0 {
        previous := sm.stateStack[len(sm.stateStack)-1]
        previous.Enter(sm)
    }
}
```

#### Context Management
```go
func (sm *StateMachine) SetContext(key string, value interface{}) {
    sm.context[key] = value
}

func (sm *StateMachine) GetContext(key string) (interface{}, bool) {
    value, exists := sm.context[key]
    return value, exists
}

func (sm *StateMachine) HasContext(key string) bool {
    _, exists := sm.context[key]
    return exists
}

func (sm *StateMachine) ClearContext(key string) {
    delete(sm.context, key)
}

func (sm *StateMachine) ClearAllContext() {
    sm.context = make(map[string]interface{})
}
```

## Lua Integration

### State Module Format

```lua
local MyState = {
    -- Called once on first creation
    initialize = function(sm)
        -- Load assets
        -- Allocate resources
        -- One-time setup
    end,
    
    -- Called every time state becomes active
    enter = function(sm)
        -- Reset variables
        -- Start music
        -- Setup for this run
        
        -- Optional: Read from shared context
        local data = game.getContext("some_key")
    end,
    
    -- Called each frame for input
    handleInput = function(sm)
        -- Process keyboard/mouse/controller
    end,
    
    -- Called each frame for logic
    update = function(dt)
        -- Update game logic
        -- dt is delta time in seconds
    end,
    
    -- Called each frame for rendering
    draw = function()
        -- Render the state
    end,
    
    -- Called when leaving state
    exit = function(sm)
        -- Pause music
        -- Save temporary data
        -- State may be resumed later
        
        -- Optional: Write to shared context
        game.setContext("some_key", someValue)
    end,
    
    -- Called once on destruction
    shutdown = function()
        -- Free assets
        -- Final cleanup
        -- State will not be used again
    end
}

-- Register the state
game.registerState("my_state", MyState)
```

### Lua API

```lua
-- State Management
game.registerState(name, stateModule)
game.unregisterState(name)
game.changeState(name)
game.pushState(name)
game.popState()
game.popAllStates()

-- Shared Context
game.setContext(key, value)           -- Set a value in context
local value = game.getContext(key)    -- Get a value (returns nil if not found)
local exists = game.hasContext(key)   -- Check if key exists
game.clearContext(key)                -- Remove a specific key
game.clearAllContext()                -- Clear all context data

-- Credits
game.addCredit(category, name, role)

-- Control
game.exit() -- Transition to credits then exit

-- Utility (optional feature for drawing previous state)
game.drawPreviousState() -- Draw state underneath in stack
```

## Complete Typical Game State Flow

```
START
  │
  ├─ Release Build
  │    │
  │    1. __engine_splash (built-in, auto)
  │    ↓ (auto-transition after 2-3s)
  │
  ├─ Debug Build
  │    (skips splash)
  │
  2. developer_splash (user-created, optional)
  ↓ (auto-transition or any key)
  
  3. main_menu (user-created)
  ├─→ 4. settings (user-created)
  │   └→ back to main_menu
  │
  ├─→ 5. level_select (user-created)
  │   │   Context: Set current_level, difficulty
  │   └→ playing
  │
  ├─→ 6. credits (user triggers)
  │   └→ EXIT
  │
  └─→ 7. playing (user-created)
      │   Context: Read current_level, difficulty
      ├─→ 8. paused (PUSH - overlay)
      │   ├→ resume (POP back to playing)
      │   ├→ settings (PUSH on top)
      │   └→ quit (CHANGE to main_menu)
      │
      ├─→ 9. inventory (PUSH - overlay)
      │   │   Context: Share player_inventory, player_gold
      │   └→ close (POP back to playing)
      │
      ├─→ 10. game_over (CHANGE)
      │   │   Context: Set final_score, enemies_killed, time_played
      │   ├→ retry (CHANGE to playing)
      │   └→ menu (CHANGE to main_menu)
      │
      └─→ 11. level_complete (CHANGE)
          │   Context: Set level_stars, completion_time
          ├→ next level (CHANGE to playing)
          └→ menu (CHANGE to main_menu)

  12. __credits (built-in)
  ↓ (any input)
  
  EXIT
```

## Usage Examples

### Example 1: Main Menu State

```lua
local MainMenuState = {
    selectedOption = 1,
    options = {"Play", "Settings", "Credits", "Exit"},
    
    initialize = function(sm)
        -- Load menu background
        -- Load menu music
    end,
    
    enter = function(sm)
        selectedOption = 1
        audio.playMusic("menu_music")
    end,
    
    handleInput = function(sm)
        if input.isPressed("up") then
            selectedOption = math.max(1, selectedOption - 1)
        elseif input.isPressed("down") then
            selectedOption = math.min(#options, selectedOption + 1)
        elseif input.isPressed("enter") then
            if selectedOption == 1 then
                game.changeState("playing")
            elseif selectedOption == 2 then
                game.pushState("settings") -- Overlay settings
            elseif selectedOption == 3 then
                game.changeState("__credits")
            elseif selectedOption == 4 then
                game.exit()
            end
        end
    end,
    
    update = function(dt)
        -- Animate menu items
    end,
    
    draw = function()
        -- Draw background
        -- Draw menu options
        -- Highlight selected
    end,
    
    exit = function(sm)
        audio.stopMusic()
    end,
    
    shutdown = function()
        -- Free menu assets
    end
}

game.registerState("main_menu", MainMenuState)
```

### Example 2: Level Selection with Context

```lua
local LevelSelectState = {
    selectedLevel = 1,
    selectedDifficulty = "normal",
    
    initialize = function(sm)
        -- Load level previews
    end,
    
    enter = function(sm)
        -- Check if player has unlocked levels
        local unlockedLevels = game.getContext("unlocked_levels") or 1
        selectedLevel = math.min(selectedLevel, unlockedLevels)
    end,
    
    handleInput = function(sm)
        if input.isPressed("left") then
            selectedLevel = math.max(1, selectedLevel - 1)
        elseif input.isPressed("right") then
            selectedLevel = math.min(maxLevels, selectedLevel + 1)
        elseif input.isPressed("space") then
            -- Cycle difficulty
            if selectedDifficulty == "easy" then
                selectedDifficulty = "normal"
            elseif selectedDifficulty == "normal" then
                selectedDifficulty = "hard"
            else
                selectedDifficulty = "easy"
            end
        elseif input.isPressed("enter") then
            -- Pass data to playing state via context
            game.setContext("current_level", selectedLevel)
            game.setContext("level_difficulty", selectedDifficulty)
            game.setContext("level_start_time", os.time())
            
            game.changeState("playing")
        elseif input.isPressed("escape") then
            game.changeState("main_menu")
        end
    end,
    
    draw = function()
        -- Draw level previews
        -- Draw difficulty selector
    end,
    
    exit = function(sm)
        -- Cleanup
    end,
    
    shutdown = function()
        -- Free assets
    end
}

game.registerState("level_select", LevelSelectState)
```

### Example 3: Playing State with Context

```lua
local PlayingState = {
    player = nil,
    enemies = {},
    score = 0,
    level = 1,
    difficulty = "normal",
    startTime = 0,
    
    initialize = function(sm)
        -- Load game assets
    end,
    
    enter = function(sm)
        -- Get level info from context
        level = game.getContext("current_level") or 1
        difficulty = game.getContext("level_difficulty") or "normal"
        startTime = game.getContext("level_start_time") or os.time()
        
        -- Get persistent player data from context
        local playerHealth = game.getContext("player_health") or 100
        local playerGold = game.getContext("player_gold") or 0
        local playerInventory = game.getContext("player_inventory") or {}
        
        -- Initialize game
        player = createPlayer(playerHealth, playerGold, playerInventory)
        enemies = spawnEnemies(level, difficulty)
        score = 0
        
        audio.playMusic("game_music")
    end,
    
    handleInput = function(sm)
        if input.isPressed("escape") then
            -- Save current state to context before pausing
            game.setContext("pause_timestamp", os.time())
            game.pushState("paused")
        elseif input.isPressed("i") then
            -- Share inventory data with overlay
            game.setContext("player_inventory", player.inventory)
            game.setContext("player_gold", player.gold)
            game.pushState("inventory")
        end
        
        player:handleInput()
    end,
    
    update = function(dt)
        player:update(dt)
        
        for _, enemy in ipairs(enemies) do
            enemy:update(dt)
        end
        
        -- Check win condition
        if #enemies == 0 then
            -- Store level completion data
            local timeTaken = os.time() - startTime
            game.setContext("level_complete_time", timeTaken)
            game.setContext("level_complete_score", score)
            game.setContext("level_complete_stars", calculateStars(score, timeTaken))
            
            -- Unlock next level
            local unlockedLevels = game.getContext("unlocked_levels") or 1
            game.setContext("unlocked_levels", math.max(unlockedLevels, level + 1))
            
            game.changeState("level_complete")
        end
        
        -- Check game over
        if player.health <= 0 then
            -- Store final stats
            game.setContext("final_score", score)
            game.setContext("enemies_killed", totalKills)
            game.setContext("time_played", os.time() - startTime)
            game.setContext("player_level", player.level)
            
            game.changeState("game_over")
        end
    end,
    
    draw = function()
        -- Draw game world
        player:draw()
        for _, enemy in ipairs(enemies) do
            enemy:draw()
        end
        -- Draw HUD
    end,
    
    exit = function(sm)
        -- Save persistent player data back to context
        game.setContext("player_health", player.health)
        game.setContext("player_gold", player.gold)
        game.setContext("player_inventory", player.inventory)
    end,
    
    shutdown = function()
        -- Free game assets
    end
}

game.registerState("playing", PlayingState)
```

### Example 4: Pause Menu State (Overlay)

```lua
local PausedState = {
    selectedOption = 1,
    options = {"Resume", "Settings", "Quit to Menu"},
    pauseDuration = 0,
    
    initialize = function(sm)
        -- Load pause UI assets
    end,
    
    enter = function(sm)
        selectedOption = 1
        
        -- Check pause reason from context
        local pauseTime = game.getContext("pause_timestamp") or os.time()
        pauseDuration = 0
        
        audio.pauseMusic()
    end,
    
    handleInput = function(sm)
        if input.isPressed("escape") then
            game.popState() -- Return to playing
        elseif input.isPressed("up") then
            selectedOption = math.max(1, selectedOption - 1)
        elseif input.isPressed("down") then
            selectedOption = math.min(#options, selectedOption + 1)
        elseif input.isPressed("enter") then
            if selectedOption == 1 then
                -- Resume
                game.clearContext("pause_timestamp")
                game.popState()
            elseif selectedOption == 2 then
                -- Stack settings on top of pause
                game.pushState("settings")
            elseif selectedOption == 3 then
                -- Quit to menu
                game.clearContext("pause_timestamp")
                
                -- Ask user to confirm (optional)
                game.setContext("confirm_action", "quit_to_menu")
                game.pushState("confirm_dialog")
            end
        end
    end,
    
    update = function(dt)
        pauseDuration = pauseDuration + dt
    end,
    
    draw = function()
        -- Draw previous state (dimmed playing state)
        game.drawPreviousState()
        
        -- Draw semi-transparent overlay
        graphics.setColor(0, 0, 0, 0.5)
        graphics.rectangle("fill", 0, 0, screen.width, screen.height)
        
        -- Draw pause menu
        graphics.setColor(1, 1, 1, 1)
        graphics.print("PAUSED", 100, 50)
        
        -- Draw menu options
        for i, option in ipairs(options) do
            if i == selectedOption then
                graphics.setColor(1, 1, 0, 1) -- Yellow for selected
            else
                graphics.setColor(1, 1, 1, 1) -- White for others
            end
            graphics.print(option, 100, 100 + i * 30)
        end
    end,
    
    exit = function(sm)
        -- If resuming, unpause music
        if selectedOption == 1 then
            audio.resumeMusic()
        end
    end,
    
    shutdown = function()
        -- Cleanup
    end
}

game.registerState("paused", PausedState)
```

### Example 5: Inventory State (Overlay with Shared Context)

```lua
local InventoryState = {
    selectedSlot = 1,
    inventory = {},
    gold = 0,
    
    initialize = function(sm)
        -- Load inventory UI
    end,
    
    enter = function(sm)
        -- Get shared inventory data from context
        inventory = game.getContext("player_inventory") or {}
        gold = game.getContext("player_gold") or 0
        selectedSlot = 1
    end,
    
    handleInput = function(sm)
        if input.isPressed("escape") or input.isPressed("i") then
            -- Close inventory
            game.popState()
        elseif input.isPressed("left") then
            selectedSlot = math.max(1, selectedSlot - 1)
        elseif input.isPressed("right") then
            selectedSlot = math.min(#inventory, selectedSlot + 1)
        elseif input.isPressed("enter") then
            -- Use/equip item
            if selectedSlot <= #inventory then
                local item = inventory[selectedSlot]
                useItem(item)
                
                -- Update context with modified inventory
                game.setContext("player_inventory", inventory)
            end
        elseif input.isPressed("delete") then
            -- Drop item
            if selectedSlot <= #inventory then
                table.remove(inventory, selectedSlot)
                selectedSlot = math.min(selectedSlot, #inventory)
                
                -- Update context
                game.setContext("player_inventory", inventory)
            end
        end
    end,
    
    update = function(dt)
        -- Animate inventory UI
    end,
    
    draw = function()
        -- Draw previous state (dimmed)
        game.drawPreviousState()
        
        -- Draw inventory overlay
        graphics.setColor(0, 0, 0, 0.7)
        graphics.rectangle("fill", 50, 50, 700, 500)
        
        -- Draw items
        graphics.setColor(1, 1, 1, 1)
        graphics.print("Gold: " .. gold, 70, 70)
        
        for i, item in ipairs(inventory) do
            local x = 70 + ((i - 1) % 5) * 100
            local y = 100 + math.floor((i - 1) / 5) * 100
            
            if i == selectedSlot then
                graphics.setColor(1, 1, 0, 1)
            else
                graphics.setColor(1, 1, 1, 1)
            end
            
            drawItem(item, x, y)
        end
    end,
    
    exit = function(sm)
        -- Save changes back to context
        game.setContext("player_inventory", inventory)
        game.setContext("player_gold", gold)
    end,
    
    shutdown = function()
        -- Cleanup
    end
}

game.registerState("inventory", InventoryState)
```

### Example 6: Game Over State with Stats

```lua
local GameOverState = {
    score = 0,
    kills = 0,
    time = 0,
    playerLevel = 1,
    isNewHighScore = false,
    
    initialize = function(sm)
        -- Load game over UI
    end,
    
    enter = function(sm)
        -- Retrieve stats from context
        score = game.getContext("final_score") or 0
        kills = game.getContext("enemies_killed") or 0
        time = game.getContext("time_played") or 0
        playerLevel = game.getContext("player_level") or 1
        
        -- Check for high score
        local highScore = game.getContext("high_score") or 0
        if score > highScore then
            isNewHighScore = true
            game.setContext("high_score", score)
            saveHighScore(score)
        end
        
        audio.playSound("game_over")
    end,
    
    handleInput = function(sm)
        if input.isPressed("enter") then
            -- Retry - start new game with same level
            game.clearContext("final_score")
            game.clearContext("enemies_killed")
            game.clearContext("time_played")
            
            game.changeState("playing")
        elseif input.isPressed("escape") then
            -- Return to menu - clear game session data
            game.clearContext("final_score")
            game.clearContext("enemies_killed")
            game.clearContext("time_played")
            game.clearContext("player_health")
            game.clearContext("player_gold")
            game.clearContext("player_inventory")
            
            game.changeState("main_menu")
        end
    end,
    
    update = function(dt)
        -- Animate game over screen
    end,
    
    draw = function()
        graphics.setColor(0.2, 0, 0, 1)
        graphics.rectangle("fill", 0, 0, screen.width, screen.height)
        
        graphics.setColor(1, 1, 1, 1)
        graphics.print("GAME OVER", 200, 100, 0, 2, 2)
        
        if isNewHighScore then
            graphics.setColor(1, 1, 0, 1)
            graphics.print("NEW HIGH SCORE!", 200, 150)
        end
        
        graphics.setColor(1, 1, 1, 1)
        graphics.print("Score: " .. score, 200, 200)
        graphics.print("Enemies Killed: " .. kills, 200, 230)
        graphics.print("Time: " .. formatTime(time), 200, 260)
        graphics.print("Level: " .. playerLevel, 200, 290)
        
        graphics.print("Press ENTER to retry", 200, 350)
        graphics.print("Press ESC for menu", 200, 380)
    end,
    
    exit = function(sm)
        -- Cleanup
    end,
    
    shutdown = function()
        -- Free assets
    end
}

game.registerState("game_over", GameOverState)
```

### Example 7: Settings State with Global Settings

```lua
local SettingsState = {
    selectedOption = 1,
    options = {"Volume", "Difficulty", "Fullscreen", "Back"},
    volume = 0.8,
    difficulty = "normal",
    fullscreen = false,
    
    initialize = function(sm)
        -- Load settings UI
    end,
    
    enter = function(sm)
        -- Load settings from context (or defaults)
        volume = game.getContext("setting_volume") or 0.8
        difficulty = game.getContext("setting_difficulty") or "normal"
        fullscreen = game.getContext("setting_fullscreen") or false
        
        selectedOption = 1
    end,
    
    handleInput = function(sm)
        if input.isPressed("up") then
            selectedOption = math.max(1, selectedOption - 1)
        elseif input.isPressed("down") then
            selectedOption = math.min(#options, selectedOption + 1)
        elseif input.isPressed("left") or input.isPressed("right") then
            local direction = input.isPressed("left") and -1 or 1
            
            if selectedOption == 1 then
                -- Adjust volume
                volume = math.max(0, math.min(1, volume + direction * 0.1))
                game.setContext("setting_volume", volume)
                audio.setMasterVolume(volume)
            elseif selectedOption == 2 then
                -- Cycle difficulty
                local difficulties = {"easy", "normal", "hard"}
                local currentIndex = 1
                for i, diff in ipairs(difficulties) do
                    if diff == difficulty then
                        currentIndex = i
                        break
                    end
                end
                
                currentIndex = currentIndex + direction
                if currentIndex < 1 then currentIndex = #difficulties end
                if currentIndex > #difficulties then currentIndex = 1 end
                
                difficulty = difficulties[currentIndex]
                game.setContext("setting_difficulty", difficulty)
            elseif selectedOption == 3 then
                -- Toggle fullscreen
                fullscreen = not fullscreen
                game.setContext("setting_fullscreen", fullscreen)
                window.setFullscreen(fullscreen)
            end
        elseif input.isPressed("enter") or input.isPressed("escape") then
            if selectedOption == 4 or input.isPressed("escape") then
                -- Back - pop state
                game.popState()
            end
        end
    end,
    
    update = function(dt)
        -- Animate settings UI
    end,
    
    draw = function()
        -- Draw previous state if pushed on top
        if game.getStackDepth() > 1 then
            game.drawPreviousState()
            graphics.setColor(0, 0, 0, 0.5)
            graphics.rectangle("fill", 0, 0, screen.width, screen.height)
        end
        
        -- Draw settings menu
        graphics.setColor(1, 1, 1, 1)
        graphics.print("SETTINGS", 200, 50, 0, 2, 2)
        
        local y = 150
        for i, option in ipairs(options) do
            if i == selectedOption then
                graphics.setColor(1, 1, 0, 1)
            else
                graphics.setColor(1, 1, 1, 1)
            end
            
            local displayText = option
            if i == 1 then
                displayText = option .. ": " .. math.floor(volume * 100) .. "%"
            elseif i == 2 then
                displayText = option .. ": " .. difficulty
            elseif i == 3 then
                displayText = option .. ": " .. (fullscreen and "ON" or "OFF")
            end
            
            graphics.print(displayText, 200, y)
            y = y + 40
        end
    end,
    
    exit = function(sm)
        -- Settings are already saved to context
    end,
    
    shutdown = function()
        -- Cleanup
    end
}

game.registerState("settings", SettingsState)
```

### Example 8: Adding Credits

```lua
-- In your game's initialization (main.lua)

-- Add your game's credits
game.addCredit("Developer", "Jane Doe", "Lead Developer")
game.addCredit("Developer", "John Smith", "Gameplay Programmer")
game.addCredit("Artist", "Alice Johnson", "Character Artist")
game.addCredit("Artist", "Bob Brown", "Environment Artist")
game.addCredit("Sound", "SoundWorks Inc.", "Music Composition")
game.addCredit("Sound", "AudioTeam", "Sound Effects")
game.addCredit("Special Thanks", "PlayTesters United", "QA Testing")
game.addCredit("Special Thanks", "Mom & Dad", "Support")

-- Engine credits are automatically included:
-- - Engine name and version
-- - Engine developer
-- - Third-party libraries used
```

## Shared Context Best Practices

### When to Use Context

**✅ Good Use Cases:**
- Passing data between states (level selection → playing)
- Storing game session data (score, stats, progress)
- Global settings (volume, difficulty, controls)
- Player persistent data (health, inventory, gold)
- Temporary state communication (pause reasons, last menu option)
- Unlockable content tracking (levels, achievements)

**❌ Avoid Using For:**
- Large objects that should be managed by states themselves
- Data that only one state needs
- Transient frame-by-frame data
- As a replacement for proper state design
- Complex objects that need serialization

### Naming Conventions

Use descriptive prefixes for clarity:

```lua
-- Player data
game.setContext("player_health", 100)
game.setContext("player_gold", 500)
game.setContext("player_inventory", items)
game.setContext("player_level", 5)

-- Level data
game.setContext("current_level", 3)
game.setContext("level_difficulty", "hard")
game.setContext("unlocked_levels", 10)

-- Session data
game.setContext("session_score", 1000)
game.setContext("session_kills", 25)
game.setContext("session_time", 300)

-- Settings
game.setContext("setting_volume", 0.8)
game.setContext("setting_fullscreen", true)
game.setContext("setting_difficulty", "normal")

-- UI state
game.setContext("ui_last_menu", "options")
game.setContext("ui_confirm_action", "quit")

-- Temporary flags
game.setContext("pause_timestamp", os.time())
game.setContext("level_start_time", os.time())
```

### Cleanup Strategies

**Clear specific data when no longer needed:**
```lua
local GameOverState = {
    exit = function(sm)
        -- Clean up game session data
        game.clearContext("session_score")
        game.clearContext("session_time")
        game.clearContext("session_kills")
    end
}
```

**Clear all temporary data on new game:**
```lua
local MainMenuState = {
    handleInput = function(sm)
        if selectedOption == "New Game" then
            -- Clear all previous game data
            game.clearAllContext()
            
            -- Initialize new game defaults
            game.setContext("player_health", 100)
            game.setContext("player_gold", 0)
            game.setContext("player_inventory", {})
            game.setContext("current_level", 1)
            
            game.changeState("playing")
        end
    end
}
```

**Preserve settings across sessions:**
```lua
-- At game start, load saved settings
local savedVolume = loadFromFile("volume") or 0.8
game.setContext("setting_volume", savedVolume)

-- At game exit, save settings
local CreditsState = {
    exit = function(sm)
        -- Save settings before exiting
        local volume = game.getContext("setting_volume")
        saveToFile("volume", volume)
    end
}
```

### Context Usage Patterns

**Pattern 1: Data Producer → Data Consumer**
```lua
-- LevelSelectState produces data
game.setContext("current_level", selectedLevel)
game.changeState("playing")

-- PlayingState consumes data
local level = game.getContext("current_level") or 1
```

**Pattern 2: Bidirectional State Communication**
```lua
-- PlayingState shares inventory
game.setContext("player_inventory", inventory)
game.pushState("shop")

-- ShopState modifies inventory
local inventory = game.getContext("player_inventory")
table.insert(inventory, purchasedItem)
game.setContext("player_inventory", inventory)
game.popState()

-- PlayingState reads updated inventory
local inventory = game.getContext("player_inventory")
```

**Pattern 3: Persistent Data Across State Changes**
```lua
-- Save before changing state
PlayingState.exit = function(sm)
    game.setContext("player_health", player.health)
    game.setContext("player_gold", player.gold)
end

-- Restore when entering state
PlayingState.enter = function(sm)
    player.health = game.getContext("player_health") or 100
    player.gold = game.getContext("player_gold") or 0
end
```

**Pattern 4: Conditional Behavior Based on Context**
```lua
local MainMenuState = {
    enter = function(sm)
        -- Show "Continue" only if game in progress
        if game.hasContext("player_health") then
            showContinueButton = true
        else
            showContinueButton = false
        end
    end
}
```

### Context Gotchas and Solutions

**Problem: Stale data**
```lua
-- BAD: Reading context only once
local PlayingState = {
    gold = game.getContext("player_gold") or 0  -- This is nil at load time!
}

-- GOOD: Read in enter()
local PlayingState = {
    gold = 0,
    
    enter = function(sm)
        gold = game.getContext("player_gold") or 0
    end
}
```

**Problem: Forgetting to clear temporary data**
```lua
-- BAD: Context pollution
game.setContext("pause_timestamp", os.time())
-- ... never cleared

-- GOOD: Clear when done
PausedState.exit = function(sm)
    game.clearContext("pause_timestamp")
end
```

**Problem: Complex object synchronization**
```lua
-- BAD: Modifying context directly
local inventory = game.getContext("player_inventory")
table.insert(inventory, item)  -- Context not updated!

-- GOOD: Set context after modification
local inventory = game.getContext("player_inventory")
table.insert(inventory, item)
game.setContext("player_inventory", inventory)
```

## Design Decisions Summary

### Initialize vs Enter
- **Initialize**: One-time setup, heavy operations (load assets)
- **Enter**: Reset state for this session (reset score, position)

### Exit vs Shutdown
- **Exit**: Temporary pause, may return later (when pushed)
- **Shutdown**: Permanent cleanup, won't return (when changed or unregistered)

### ChangeState vs Push/Pop
- **ChangeState**: Complete replacement, old state destroyed
- **Push/Pop**: Stack states, previous state frozen in memory

### Context vs State Variables
- **Context**: Share data between states, persistent across transitions
- **State Variables**: Internal to state, destroyed with state

### Built-in States
- **Engine Splash**: Branding, skipped in debug
- **Credits**: Forced final state before exit, customizable

### Debug vs Release
- **Debug**: Skip engine splash, start directly with user states
- **Release**: Show engine splash first

## State Machine Guarantees

1. **Initialization happens once** - Before first Enter()
2. **Shutdown happens once** - After final Exit() when unregistered or changed
3. **Enter always follows Exit** - States are properly cleaned up
4. **Only top state is active** - HandleInput/Update/Draw called on stack top only
5. **Previous states preserved when pushed** - Can resume exactly where left off
6. **Credits always before exit** - Ensures proper attribution
7. **Context persists across transitions** - Data available until explicitly cleared

## Future Considerations

### Optional Features to Add Later
- State transitions with fade effects
- State event system (broadcast events to all states)
- State history tracking (navigate back/forward)
- State groups (disable/enable multiple states)
- Conditional state transitions with guards
- Async state loading with loading screens
- State serialization (save/load game state)
- Context change listeners (notify when context values change)

### Performance Notes
- Limit stack depth (recommend max 3-4 states)
- Monitor memory with many states registered
- Consider state pooling for frequently created/destroyed states
- Profile Draw() when rendering multiple stacked states
- Context is lightweight (hash map), but avoid storing large objects
- Clear context regularly to prevent memory leaks

### Architecture Benefits

1. **Decoupling** - States don't need direct references to each other
2. **Flexibility** - Easy to add/remove states without breaking others
3. **Testability** - States can be tested in isolation
4. **Maintainability** - Clear separation of concerns
5. **Reusability** - States can be reused across projects
6. **Optional Complexity** - Use advanced features (push/pop, context) only when needed
