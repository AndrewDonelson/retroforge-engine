-- Galaxy Simulation
-- A physics-based spiral galaxy simulation:
-- - Module-based state system (rf.import)
-- - Game state machine (game.changeState, built-in splash/credits)
-- - Automatic sprite pooling (star sprite with maxSpawn = 2048)
-- - Central gravity physics (point gravity simulation)
-- - 2048 stars orbiting to form spiral galaxy structure

-- Import state modules
-- Credits are set up in menu_state.lua _INIT() to ensure game object is available
local menu_state = rf.import("menu_state.lua")
local play_state = rf.import("play_state.lua")

-- Set initial state context (optional) - wrapped in check for safety
  if game then
    game.setContext("demo_name", "Galaxy Simulation")
    game.setContext("features", {
      "Module-based states",
      "Automatic sprite pooling",
      "Central gravity physics",
      "Spiral galaxy formation"
    })
  end

-- The game will:
-- 1. Show engine splash (auto, handled by engine)
-- 2. Transition to menu (handled by engine after splash)
-- 3. Menu -> Play (when Start selected)
-- 4. Play -> Credits (when timer expires)
-- 5. Credits -> Exit (when any key pressed)

