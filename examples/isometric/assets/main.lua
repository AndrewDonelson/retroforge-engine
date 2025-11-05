-- Isometric Tilemap Demo
-- Demonstrates isometric tilemaps with animated character walking

-- Import state modules
local menu_state = rf.import("menu_state.lua")
local level_select_state = rf.import("level_select_state.lua")
local play_state = rf.import("play_state.lua")

-- Set initial state context (optional) - wrapped in check for safety
if game then
  game.setContext("demo_name", "Isometric Tilemap Demo")
  game.setContext("features", {
    "Isometric tilemaps",
    "Tilesets",
    "Animated character",
    "State machine"
  })
end

-- The game will:
-- 1. Show engine splash (auto, handled by engine)
-- 2. Transition to menu (handled by engine after splash)
-- 3. Menu -> Play (when Start selected)
-- 4. Play shows isometric tilemap with walking character

