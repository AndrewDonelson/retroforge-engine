-- Menu State Module for Animated Character Demo

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
end

function _HANDLE_INPUT()
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
  -- Initialize animation on first update
  if not isPaused and animSpeed == 1.0 then
    pcall(function()
      rf.playAnimation("character_walk", "walk_cycle")
      rf.setAnimationSpeed("character_walk", animSpeed)
    end)
    animSpeed = 1.1
  end
  
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
  rf.cls(14)
  rf.print("Animated Character Demo", 10, 10, 0)
  rf.print("Arrow Keys: Move", 10, 20, 0)
  rf.print("A Button: Pause/Resume", 10, 30, 0)
  rf.print("B Button: Stop Animation", 10, 40, 0)
  rf.print("X Button: Speed Up | Y Button: Slow Down", 10, 50, 0)
  rf.print(string.format("Speed: %.1fx", animSpeed), 10, 60, 0)
  if isPaused then
    rf.print("PAUSED", 10, 70, 2)
  end
  
  if isWalking then
    pcall(function() rf.spr("character_walk", charX, charY, "walk_cycle") end)
  else
    local frameName = "idle_" .. direction
    pcall(function() rf.spr("character_idle", charX, charY, frameName) end)
  end
  
  rf.rect(10, 220, 310, 50, 1)
  rf.rectb(10, 220, 310, 50, 0)
  rf.print("Frames Sprite: 4 directional states", 15, 225, 0)
  rf.print("Animation Sprite: 4-frame walk cycle", 15, 235, 0)
  rf.print("Press keys to control animation!", 15, 245, 0)
end

function _EXIT()
  pcall(function() rf.stopAnimation("character_walk") end)
end

function _DONE()
end
