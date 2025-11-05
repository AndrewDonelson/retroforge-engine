-- Play State Module for Isometric Tilemap Demo

local charX = 240
local charY = 135
local direction = "down"
local isWalking = false
local animSpeed = 1.0

-- Get selected level from level_select_state
local function getSelectedLevel()
  -- Access the selected level from level_select_state
  -- For now, we'll use a global or default to "isometric"
  if _G.selected_level_index then
    return _G.selected_level_index
  end
  return 2  -- Default to isometric (level 2)
end

local function getCurrentMapName()
  local level_idx = getSelectedLevel()
  if level_idx == 1 then
    return "test_normal"  -- Normal tilemap
  else
    return "test"  -- Isometric tilemap
  end
end

function _INIT()
  -- Module initialization
end

function _ENTER()
  -- Reset state when entering
  charX = 240
  charY = 135
  direction = "down"
  isWalking = false
  animSpeed = 1.0
  
  -- Default to isometric if no level selected
  if not _G.selected_level_index then
    _G.selected_level_index = 2  -- Default to isometric
  end
  
  -- Initialize animation when entering state
  pcall(function()
    rf.playAnimation("character_walk", "walk_cycle")
    rf.setAnimationSpeed("character_walk", animSpeed)
  end)
end

function _HANDLE_INPUT()
  -- Button 0: SELECT - Return to menu
  if rf.btnp(0) then
    rf.sfx("stopall")
    if game then
      game.changeState("menu")
    end
    return
  end
  
  -- Button 1: START - Return to menu
  if rf.btnp(1) then
    rf.sfx("stopall")
    if game then
      game.changeState("menu")
    end
    return
  end
end

function _UPDATE(dt)
  -- Handle movement
  if rf.btn(2) then
    direction = "up"
    isWalking = true
    charY = math.max(8, charY - 60 * dt)
  elseif rf.btn(3) then
    direction = "down"
    isWalking = true
    charY = math.min(262, charY + 60 * dt)
  elseif rf.btn(4) then
    direction = "left"
    isWalking = true
    charX = math.max(8, charX - 60 * dt)
  elseif rf.btn(5) then
    direction = "right"
    isWalking = true
    charX = math.min(472, charX + 60 * dt)
  else
    isWalking = false
  end
end

function _DRAW()
  rf.clear_i(0)
  
  -- Draw tilemap (normal or isometric based on selection)
  local mapName = getCurrentMapName()
  
  -- Calculate center offset for tilemap
  -- Screen is 480x270, tilemap is 10x10 tiles
  -- Normal tiles: 16x16 pixels, Isometric tiles: 16x8 pixels (after conversion)
  local level_idx = getSelectedLevel()
  local offsetX, offsetY
  
  if level_idx == 1 then
    -- Normal tilemap: 10 tiles × 16 pixels = 160 pixels
    local mapWidth = 10 * 16
    local mapHeight = 10 * 16
    offsetX = (480 - mapWidth) / 2
    offsetY = (270 - mapHeight) / 2
  else
    -- Isometric tilemap: center the diamond shape
    -- For 10x10 isometric map with 16x8 tiles (after conversion):
    -- Formula: screenX = offsetX + (mapX-mapY)*(tile.Width/2)
    --          screenY = offsetY + (mapX+mapY)*(tile.Height/2)
    -- 
    -- Center of diamond in map coordinates is at (mapX=4.5, mapY=4.5)
    -- This maps to: screenX = offsetX + 0, screenY = offsetY + 36
    -- To center at screen center (240, 135):
    offsetX = 240  -- Center X of diamond aligns with screen center X
    offsetY = 135 - 36  -- Center Y of diamond (at Y=36 relative to offset) aligns with screen center Y
  end
  
  rf.drawTilemap(mapName, offsetX, offsetY)
  
  -- Draw character based on state
  if isWalking then
    -- Draw walking animation sprite
    pcall(function() rf.spr("character_walk", charX, charY, "walk_cycle") end)
  else
    -- Draw idle frame sprite based on direction
    local frameName = "idle_" .. direction
    pcall(function() rf.spr("character_idle", charX, charY, frameName) end)
  end
end

function _EXIT()
  pcall(function() rf.stopAnimation("character_walk") end)
end

function _DONE()
end

