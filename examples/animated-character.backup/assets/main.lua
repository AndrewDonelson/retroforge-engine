-- Kitchen Sink Demo
-- Comprehensive showcase of all RetroForge Engine features:
-- - Module-based state system (rf.import)
-- - Game state machine (game.changeState, built-in splash/credits)
-- - Automatic sprite pooling (ball sprite with maxSpawn > 10)
-- - Physics engine (Box2D integration)
-- - Stats display (FPS, memory, object counts)

-- Import state modules
-- Credits are set up in menu_state.lua _INIT() to ensure game object is available
local splash_state = rf.import("splash_state.lua")
local menu_state = rf.import("menu_state.lua")

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
-- 2. Transition to splash (handled by engine after splash)
-- 3. Splash -> Menu (after 3 seconds or any input)
