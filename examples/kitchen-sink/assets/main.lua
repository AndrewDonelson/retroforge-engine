-- Kitchen Sink Demo
-- Comprehensive showcase of all RetroForge Engine features:
-- - Module-based state system (rf.import)
-- - Game state machine (game.changeState, built-in splash/credits)
-- - Automatic sprite pooling (ball sprite with maxSpawn > 10)
-- - Physics engine (Box2D integration)
-- - Stats display (FPS, memory, object counts)

-- Import state modules
-- Credits are set up in menu_state.lua _INIT() to ensure game object is available
local menu_state = rf.import("menu_state.lua")
local play_state = rf.import("play_state.lua")

-- Set initial state context (optional) - wrapped in check for safety
if game then
  game.setContext("demo_name", "Kitchen Sink")
  game.setContext("features", {
    "Module-based states",
    "Automatic sprite pooling",
    "Physics simulation",
    "Stats display"
  })
end

-- The game will:
-- 1. Show engine splash (auto, handled by engine)
-- 2. Transition to menu (handled by engine after splash)
-- 3. Menu -> Play (when Start selected)
-- 4. Play -> Credits (when timer expires)
-- 5. Credits -> Exit (when any key pressed)

