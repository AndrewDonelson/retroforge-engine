# Animated Character Demo

A cute demonstration of RetroForge's multi-frame sprite and animation system.

## Features Demonstrated

- **Frames Sprite** (`character_idle`): Contains 4 directional states (up, down, left, right)
  - Each frame can be selected individually using `rf.spr(name, x, y, frameName)`
  
- **Animation Sprite** (`character_walk`): Contains a 4-frame walk cycle animation
  - Uses `rf.playAnimation()` to start the animation
  - Animation loops automatically with configurable speed
  - Demonstrates pause, resume, stop, and speed controls

## Controls

- **Arrow Keys (Buttons 2-5)**: Move the character (triggers walk animation)
- **A Button (Button 6)**: Pause/Resume animation
- **B Button (Button 7)**: Stop animation
- **X Button (Button 8)**: Speed up animation
- **Y Button (Button 9)**: Slow down animation

## Sprite Creation

All sprites were created using primitive shapes (circles and ellipses) to form a cute blob character. The character uses the Pastel 50 palette for a soft, engaging appearance.

## Code Highlights

```lua
-- main.lua: Import state module
local menu_state = rf.import("menu_state.lua")

-- menu_state.lua: State module structure
function _INIT() end
function _ENTER()
    -- Start animation when entering state
    rf.playAnimation("character_walk", "walk_cycle")
end
function _HANDLE_INPUT()
    -- Handle button presses (instant actions)
    -- Button 6 = A, Button 7 = B
    if rf.btnp(6) then  -- A button
        rf.pauseAnimation("character_walk")
    end
end
function _UPDATE(dt)
    -- Handle movement (continuous actions using dt)
    -- Button 2 = UP, Button 3 = DOWN, Button 4 = LEFT, Button 5 = RIGHT
    if rf.btn(2) then  -- UP button
        charY = charY - 60 * dt
    end
end
function _DRAW()
    -- Switch between frames sprite and animation sprite
    if isWalking then
        rf.spr("character_walk", x, y, "walk_cycle")
    else
        rf.spr("character_idle", x, y, "idle_" .. direction)
    end
end
function _EXIT() end
function _DONE() end

-- Animation control API
rf.playAnimation("character_walk", "walk_cycle")
rf.pauseAnimation("character_walk")
rf.resumeAnimation("character_walk")
rf.setAnimationSpeed("character_walk", speed)
```

This example showcases:
- **Module-based state system** (`rf.import()`)
- **Multi-frame sprites** (directional states)
- **Animation sprites** (walk cycle)
- **Animation controls** (play, pause, speed)

