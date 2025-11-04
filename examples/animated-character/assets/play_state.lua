-- Play State Module for Animated Character Demo

local charX = 240
local charY = 135
local direction = "down"
local isWalking = false
local animSpeed = 1.0
local isPaused = false

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
  isPaused = false
  
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
  
  -- Button 6: A - Toggle pause/resume
  if rf.btnp(6) then
    if isPaused then
      pcall(function() rf.resumeAnimation("character_walk") end)
      isPaused = false
    else
      pcall(function() rf.pauseAnimation("character_walk") end)
      isPaused = true
    end
  end
  
  -- Button 7: B - Stop animation
  if rf.btnp(7) then
    pcall(function() rf.stopAnimation("character_walk") end)
    isPaused = false
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
  
  -- Animation speed controls
  if rf.btn(8) then
    animSpeed = math.min(3.0, animSpeed + 0.5 * dt)
    pcall(function() rf.setAnimationSpeed("character_walk", animSpeed) end)
  elseif rf.btn(9) then
    animSpeed = math.max(0.1, animSpeed - 0.5 * dt)
    pcall(function() rf.setAnimationSpeed("character_walk", animSpeed) end)
  end
end

function _DRAW()
  rf.clear_i(0)
  rf.print("Animated Character Demo", 10, 10, 15)
  rf.print("Arrow Keys: Move", 10, 20, 7)
  rf.print("A Button: Pause/Resume", 10, 30, 7)
  rf.print("B Button: Stop Animation", 10, 40, 7)
  rf.print("X Button: Speed Up | Y Button: Slow Down", 10, 50, 7)
  rf.print(string.format("Speed: %.1fx", animSpeed), 10, 60, 7)
  if isPaused then
    rf.print("PAUSED", 10, 70, 2)
  end
  
  -- Draw character based on state
  if isWalking then
    -- Draw walking animation sprite
    pcall(function() rf.spr("character_walk", charX, charY, "walk_cycle") end)
  else
    -- Draw idle frame sprite based on direction
    local frameName = "idle_" .. direction
    pcall(function() rf.spr("character_idle", charX, charY, frameName) end)
  end
  
  -- Info box at bottom (centered)
  local box_x = 10
  local box_y = 200
  local box_w = 460
  local box_h = 60
  rf.rect(box_x, box_y, box_x + box_w, box_y + box_h, 1)
  rf.rectb(box_x, box_y, box_x + box_w, box_y + box_h, 0)
  rf.print("Frames Sprite: 4 directional states", box_x + 5, box_y + 5, 7)
  rf.print("Animation Sprite: 4-frame walk cycle", box_x + 5, box_y + 20, 7)
  rf.print("Press keys to control animation!", box_x + 5, box_y + 35, 7)
end

function _EXIT()
  pcall(function() rf.stopAnimation("character_walk") end)
end

function _DONE()
end
