-- game.lua
-- Simple multiplayer platformer example
-- Based on RetroForge.V2.md design document

-- Game state tables
players = {}
platforms = {}
score = {}

function _init()
  -- Create platforms (host only)
  if rf.is_multiplayer() and rf.is_host() then
    platforms = {
      {x=0, y=250, w=480, h=20},
      {x=100, y=200, w=100, h=20},
      {x=300, y=150, w=100, h=20},
      {x=200, y=100, w=100, h=20}
    }
    rf.network_sync(platforms, "slow")  -- Platforms never move
  elseif not rf.is_multiplayer() then
    -- Solo mode: also create platforms
    platforms = {
      {x=0, y=250, w=480, h=20},
      {x=100, y=200, w=100, h=20},
      {x=300, y=150, w=100, h=20},
      {x=200, y=100, w=100, h=20}
    }
  end
  
  -- Create player for each connected player (or just one for solo)
  local player_count = rf.is_multiplayer() and rf.player_count() or 1
  for i = 1, player_count do
    players[i] = {
      x = 50 + (i-1) * 100,
      y = 100,
      vx = 0,
      vy = 0,
      sprite = i % 4,  -- Cycle through sprite indices
      health = 100,
      on_ground = false
    }
    score[i] = 0
  end
  
  -- Register for automatic synchronization (multiplayer only)
  if rf.is_multiplayer() then
    rf.network_sync(players, "fast")     -- Player positions update smoothly
    rf.network_sync(score, "slow")       -- Score updates don't need high frequency
  end
end

function _update()
  if rf.is_multiplayer() and rf.is_host() then
    -- HOST: Run game logic for all players
    update_host()
  elseif rf.is_multiplayer() then
    -- NON-HOST: Just send inputs (engine handles automatically)
    -- Normal btn() calls work for local player
    -- Engine will send input to host automatically
  else
    -- SOLO MODE: Run game logic normally
    update_solo()
  end
end

function update_solo()
  -- Single player mode - same logic as host but for one player
  local p = players[1]
  
  -- Apply inputs
  if rf.btn(1) then p.vx = 3 end     -- right
  if rf.btn(0) then p.vx = -3 end    -- left
  if rf.btn(4) and p.on_ground then  -- jump
    p.vy = -10
    p.on_ground = false
  end
  
  -- Apply physics
  p.vy = p.vy + 0.5  -- gravity
  p.x = p.x + p.vx
  p.y = p.y + p.vy
  
  -- Apply friction
  p.vx = p.vx * 0.9
  
  -- Check platform collisions
  p.on_ground = false
  for _, plat in ipairs(platforms) do
    if check_collision(p, plat) then
      p.y = plat.y - 16
      p.vy = 0
      p.on_ground = true
    end
  end
  
  -- Wrap screen
  if p.x < 0 then p.x = 480 end
  if p.x > 480 then p.x = 0 end
  
  -- Check if player fell off bottom
  if p.y > 270 then
    p.y = 100
    p.x = 50
    score[1] = score[1] - 10  -- Penalty
  end
end

function update_host()
  -- Update each player based on their inputs
  for id = 1, rf.player_count() do
    local p = players[id]
    
    -- Apply inputs (5/sec, but smooth with interpolation)
    if rf.btn(id, 0) then p.vx = -3 end    -- left
    if rf.btn(id, 1) then p.vx = 3 end     -- right
    if rf.btn(id, 4) and p.on_ground then  -- jump
      p.vy = -10
      p.on_ground = false
    end
    
    -- Apply physics
    p.vy = p.vy + 0.5  -- gravity
    p.x = p.x + p.vx
    p.y = p.y + p.vy
    
    -- Apply friction
    p.vx = p.vx * 0.9
    
    -- Check platform collisions
    p.on_ground = false
    for _, plat in ipairs(platforms) do
      if check_collision(p, plat) then
        p.y = plat.y - 16
        p.vy = 0
        p.on_ground = true
      end
    end
    
    -- Wrap screen
    if p.x < 0 then p.x = 480 end
    if p.x > 480 then p.x = 0 end
    
    -- Check if player fell off bottom
    if p.y > 270 then
      p.y = 100
      p.x = 50 + (id-1) * 100
      score[id] = score[id] - 10  -- Penalty
    end
  end
  
  -- Engine automatically syncs players and score tables!
end

function _draw()
  -- Both host and non-host render the same way
  rf.clear_i(0)  -- Clear to black
  
  -- Draw platforms
  for _, plat in ipairs(platforms) do
    rf.rectfill(plat.x, plat.y, plat.x + plat.w, plat.y + plat.h, 7)
  end
  
  -- Draw all players
  for id, p in pairs(players) do
    -- Simple sprite representation (colored box)
    local color_idx = 2 + (id % 6)  -- Different colors per player
    rf.rectfill(p.x, p.y, p.x + 16, p.y + 16, color_idx)
    
    -- Show player ID above sprite
    local color = (rf.is_multiplayer() and id == rf.my_player_id()) and 11 or 6
    rf.print(tostring(id), p.x + 4, p.y - 8, color)
  end
  
  -- Draw scores
  for id, s in pairs(score) do
    rf.print("P" .. tostring(id) .. ": " .. tostring(s), 10, 10 + (id-1) * 10, 7)
  end
  
  -- Show connection info
  if rf.is_multiplayer() then
    rf.print("Players: " .. tostring(rf.player_count()), 10, 250, 7)
    if rf.is_host() then
      rf.print("HOST", 10, 260, 10)
    else
      rf.print("PLAYER " .. tostring(rf.my_player_id()), 10, 260, 9)
    end
  end
end

function check_collision(player, platform)
  -- Simple AABB collision
  return player.x < platform.x + platform.w and
         player.x + 16 > platform.x and
         player.y < platform.y + platform.h and
         player.y + 16 > platform.y and
         player.vy >= 0  -- Only collide when falling
end

