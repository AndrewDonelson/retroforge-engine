-- Animated Character Demo
-- Demonstrates multi-frame sprites and animations
-- Features:
-- - Frames sprite with directional states (up, down, left, right)
-- - Animation sprite with walk cycle
-- - Animation controls (play, pause, speed adjustment)
-- - Module-based state system (rf.import)

-- Import state modules
-- Credits are set up in menu_state.lua _INIT() to ensure game object is available
local menu_state = rf.import("menu_state.lua")
local play_state = rf.import("play_state.lua")

-- Set initial state context (optional) - wrapped in check for safety
if game then
  game.setContext("demo_name", "Animated Character Demo")
  game.setContext("features", {
    "Multi-frame sprites",
    "Animation sprites",
    "Animation controls",
    "State machine"
  })
end

-- The game will:
-- 1. Show engine splash (auto, handled by engine)
-- 2. Transition to menu (handled by engine after splash)
-- 3. Menu -> Play (when Start selected)
-- 4. Play -> Credits (when timer expires)
-- 5. Credits -> Exit (when any key pressed)
