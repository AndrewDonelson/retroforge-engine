-- game.lua
-- Simple multiplayer platformer example
-- Based on RetroForge.V2.md design document

-- Game state tables
players = {}
platforms = {}
score = {}
state = "menu"  -- "menu" or "playing"
menu_selected = 1
prev_was_grounded = {}  -- Track previous grounded state for landing sound
prev_vx = {}  -- Track previous velocity for movement sound
winner = nil  -- Winner player ID
game_over = false
spectator_target = nil  -- Which player ID is being spectated (nil = first alive)
prev_tab_pressed = false  -- Track Tab key state
exit_button_hover = false  -- Track if mouse/pointer is over EXIT button

function _INIT()
  rf.palette_set("SNES 50")
  state = "menu"
  
  -- Create platforms (host only) - SNES colors: vibrant reds, oranges, yellows, greens, blues
  -- Platforms are thinner now (h=8 instead of 20)
  -- SNES palette structure: 0=black, 1=white, then 16 hues × 3 shades (highlight/base/shadow)
  if rf.is_multiplayer() and rf.is_host() then
    platforms = {
      {x=0, y=262, w=480, h=8, color=38},      -- Ground platform
      {x=80, y=220, w=100, h=8, color=23},    -- Green platform
      {x=300, y=220, w=100, h=8, color=35},   -- Blue platform
      {x=180, y=180, w=120, h=8, color=20},   -- Yellow/orange platform
      {x=50, y=140, w=80, h=8, color=17},     -- Purple platform
      {x=350, y=140, w=80, h=8, color=11},   -- Red platform
      {x=150, y=100, w=100, h=8, color=23},   -- Green platform
      {x=330, y=100, w=100, h=8, color=35},   -- Blue platform
      {x=100, y=60, w=80, h=8, color=20},    -- Yellow platform
      {x=300, y=60, w=80, h=8, color=17},     -- Purple platform
      {x=200, y=30, w=80, h=8, color=2}      -- Red platform (near top)
    }
    rf.network_sync(platforms, "slow")  -- Platforms never move
  elseif not rf.is_multiplayer() then
    -- Solo mode: also create platforms
    platforms = {
      {x=0, y=262, w=480, h=8, color=38},      -- Ground platform
      {x=80, y=220, w=100, h=8, color=23},     -- Green platform
      {x=300, y=220, w=100, h=8, color=35},    -- Blue platform
      {x=180, y=180, w=120, h=8, color=20},    -- Yellow/orange platform
      {x=50, y=140, w=80, h=8, color=17},      -- Purple platform
      {x=350, y=140, w=80, h=8, color=11},     -- Red platform
      {x=150, y=100, w=100, h=8, color=23},    -- Green platform
      {x=330, y=100, w=100, h=8, color=35},    -- Blue platform
      {x=100, y=60, w=80, h=8, color=20},      -- Yellow platform
      {x=300, y=60, w=80, h=8, color=17},      -- Purple platform
      {x=200, y=30, w=80, h=8, color=2}       -- Red platform (near top)
    }
  end
  
  -- Reset game state
  winner = nil
  game_over = false
  spectator_target = nil
  prev_tab_pressed = false
  exit_button_hover = false
  
  -- Create player for each connected player (or just one for solo)
  -- SNES player colors: red, blue, green, yellow, purple, orange
  local player_colors = {2, 14, 23, 20, 17, 11}  -- SNES vibrant colors per player
  local player_count = rf.is_multiplayer() and rf.player_count() or 1
  for i = 1, player_count do
    players[i] = {
      x = 50 + (i-1) * 100,
      y = 254,  -- Start on bottom platform (262 - 8)
      vx = 0,
      vy = 0,
      sprite = i % 4,  -- Cycle through sprite indices
      health = 100,
      on_ground = true,  -- Start on ground
      alive = true,  -- Track if player is alive
      color = player_colors[((i-1) % #player_colors) + 1]  -- Different vibrant colors per player
    }
    score[i] = 0
    prev_was_grounded[i] = true
    prev_vx[i] = 0
  end
  
  -- Register for automatic synchronization (multiplayer only)
  if rf.is_multiplayer() then
    rf.network_sync(players, "fast")     -- Player positions update smoothly (includes alive state)
    rf.network_sync(score, "slow")       -- Score updates don't need high frequency
  end
end

function update_menu(dt)
  if rf.btnp(2) then 
    menu_selected = math.max(1, menu_selected - 1)
    rf.sfx("move")
  end
  if rf.btnp(3) then 
    menu_selected = math.min(2, menu_selected + 1)
    rf.sfx("move")
  end
  if rf.btnp(4) or rf.btnp(5) then  -- Z or X
    rf.sfx("select")
    if menu_selected == 1 then
      state = "playing"
      rf.music("bgm")
    else
      rf.quit()
    end
  end
end

function get_alive_players()
  -- Get list of alive player IDs
  local alive = {}
  for id, p in pairs(players) do
    if p.alive then
      table.insert(alive, id)
    end
  end
  table.sort(alive)  -- Sort by ID
  return alive
end

function cycle_spectator_target()
  -- Cycle to next alive player
  local alive = get_alive_players()
  if #alive == 0 then
    spectator_target = nil
    return
  end
  
  if spectator_target == nil then
    spectator_target = alive[1]
    return
  end
  
  -- Find current target index
  local current_idx = 1
  for i, id in ipairs(alive) do
    if id == spectator_target then
      current_idx = i
      break
    end
  end
  
  -- Move to next
  current_idx = current_idx + 1
  if current_idx > #alive then
    current_idx = 1
  end
  spectator_target = alive[current_idx]
end

function update_spectator_mode()
  -- Check if local player is dead
  -- In solo mode, my_player_id might not work, so check player 1
  local my_id = 1
  if rf.is_multiplayer() then
    my_id = rf.my_player_id()
  end
  local my_player = players[my_id]
  
  if not my_player or my_player.alive then
    prev_tab_pressed = false  -- Reset tab state when not spectating
    return false  -- Not in spectator mode
  end
  
  -- Dead player: handle spectator controls
  -- Initialize spectator target if needed
  if spectator_target == nil then
    local alive = get_alive_players()
    if #alive > 0 then
      spectator_target = alive[1]
    else
      -- No alive players - can't spectate
      return false
    end
  end
  
  -- Cycle spectator target: Hold Up (2) + Press X (5) = Tab replacement
  local tab_pressed = rf.btn(2) and rf.btnp(5)
  
  if tab_pressed and not prev_tab_pressed then
    cycle_spectator_target()
    rf.sfx("move")
  end
  prev_tab_pressed = tab_pressed
  
  -- EXIT button: Hold Down (3) + Press X (5)
  if rf.btn(3) and rf.btnp(5) then
    rf.sfx("select")
    state = "menu"
    game_over = false
    winner = nil
    spectator_target = nil
    prev_tab_pressed = false
    rf.music("stopall")
    return true
  end
  
  return false
end

function _UPDATE(dt)
  if state == "menu" then
    update_menu(dt)
    -- Note: ESC quitting is handled at SDL level, but players should use QUIT menu option
    return
  end
  
  -- Check spectator mode for dead players
  if update_spectator_mode() then
    return  -- Exited to menu
  end
  
  -- Allow restart when game over (press Z/X)
  if game_over then
    if rf.btnp(4) or rf.btnp(5) then
      -- Restart game
      state = "menu"
      game_over = false
      winner = nil
      spectator_target = nil
      rf.music("stopall")
    end
    return  -- Stop gameplay when game over
  end
  
  -- Check for win condition (player reached top)
  for id, p in pairs(players) do
    if p.alive and p.y <= 10 then  -- Reached top of screen
      winner = id
      game_over = true
      score[id] = score[id] + 100  -- Bonus for winning
      rf.music("stopall")
      break
    end
  end
  
  if game_over then
    return  -- Stop gameplay when game over
  end
  
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
  local id = 1
  
  -- If dead, enter spectator mode (though in solo there's no one to spectate)
  if not p.alive then
    -- In solo mode, if dead, can restart
    if rf.btnp(4) or rf.btnp(5) then
      state = "menu"
      game_over = false
      winner = nil
      rf.music("stopall")
    end
    return
  end
  
  -- Store previous state for sound detection
  local was_grounded = p.on_ground
  local old_vx = p.vx
  local old_y = p.y
  
  -- Apply inputs
  if rf.btn(1) then p.vx = 3 end     -- right
  if rf.btn(0) then p.vx = -3 end    -- left
  if rf.btn(4) and p.on_ground then  -- jump
    p.vy = -5  -- Reduced jump strength by 50% (was -10)
    p.on_ground = false
    rf.sfx("jump")
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
      -- Play landing sound if wasn't grounded before
      if not was_grounded then
        rf.sfx("land")
      end
    end
  end
  
  -- Wrap screen
  if p.x < 0 then p.x = 480 end
  if p.x > 480 then p.x = 0 end
  
  -- Check if player fell off bottom
  if p.y > 270 then
    rf.sfx("fall")
    p.y = 100
    p.x = 50
    p.on_ground = false
    score[1] = score[1] - 10  -- Penalty
  end
  
  -- Check win condition (reached top)
  if p.y <= 10 then
    winner = id
    game_over = true
    score[id] = score[id] + 100
    rf.music("stopall")
  end
end

function check_head_collision(jumper, landee)
  -- Check if jumper landed on top of landee (head jump = elimination)
  -- Jumper must be falling down and landee must be below
  if jumper.vy <= 0 then return false end  -- Jumper not falling
  if not landee.alive then return false end  -- Landee must be alive
  if not jumper.alive then return false end  -- Jumper must be alive
  
  -- Check if jumper is directly above landee (landing on head)
  local jumper_bottom = jumper.y + 16
  local jumper_left = jumper.x
  local jumper_right = jumper.x + 16
  local landee_top = landee.y
  local landee_left = landee.x
  local landee_right = landee.x + 16
  
  -- Jumper's bottom must be near landee's top, and horizontally overlapping
  -- More lenient collision - any horizontal overlap counts
  if jumper_bottom >= landee_top - 1 and jumper_bottom <= landee_top + 10 and
     jumper_left < landee_right and jumper_right > landee_left then
    return true
  end
  
  return false
end

function update_host()
  -- Update each player based on their inputs
  for id = 1, rf.player_count() do
    local p = players[id]
    
    if not p.alive then 
      -- Dead players fall off screen
      p.y = p.y + 2
      if p.y > 270 then
        p.y = 270  -- Keep them off screen
      end
    else
      -- Store previous state for sound detection
      local was_grounded = p.on_ground
      local old_vx = p.vx
      local old_y = p.y
      
      -- Apply inputs (5/sec, but smooth with interpolation)
      if rf.btn(id, 0) then p.vx = -3 end    -- left
      if rf.btn(id, 1) then p.vx = 3 end     -- right
      if rf.btn(id, 4) and p.on_ground then  -- jump
        p.vy = -5  -- Reduced jump strength by 50% (was -10)
        p.on_ground = false
        rf.sfx("jump")  -- Host plays sounds for all players
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
          -- Play landing sound if wasn't grounded before
          if not was_grounded then
            rf.sfx("land")
          end
        end
      end
      
      -- Check head collisions (jumping on other players' heads)
      -- Only check if this player is falling (jumping down on someone)
      if p.vy > 0 then  -- Only when falling
        for other_id, other_p in pairs(players) do
          if other_id ~= id and check_head_collision(p, other_p) then
            -- Jumper lands on head - landee dies
            other_p.alive = false
            rf.sfx("fall")  -- Death sound
            score[other_id] = score[other_id] - 50  -- Death penalty
            -- Give jumper a bounce
            p.vy = -6
            p.on_ground = false
            score[id] = score[id] + 10  -- Bonus for eliminating opponent
            break
          end
        end
      end
      
      -- Wrap screen
      if p.x < 0 then p.x = 480 end
      if p.x > 480 then p.x = 0 end
      
      -- Check if player fell off bottom
      if p.y > 270 then
        rf.sfx("fall")
        p.y = 100
        p.x = 50 + (id-1) * 100
        p.on_ground = false
        score[id] = score[id] - 10  -- Penalty
      end
      
      -- Check win condition (reached top)
      if p.y <= 10 then
        winner = id
        game_over = true
        score[id] = score[id] + 100
        rf.music("stopall")
        break
      end
    end
  end
  
  -- Engine automatically syncs players and score tables!
end

function draw_menu()
  -- SNES-style menu with vibrant colors
  rf.print_anchored("PLATFORMER", "topcenter", 20)  -- Bright yellow
  
  local c1 = (menu_selected == 1) and 20 or 1  -- Yellow when selected, white otherwise
  local c2 = (menu_selected == 2) and 20 or 1
  
  local play_x = 240 - string.len("PLAY") * 3
  local quit_x = 240 - string.len("QUIT") * 3
  rf.print_xy(play_x, 100, "PLAY", c1)
  rf.print_xy(quit_x, 116, "QUIT", c2)
  
  -- Selection arrow
  if menu_selected == 1 then
    rf.print_xy(play_x - 15, 100, ">", 20)
  else
    rf.print_xy(quit_x - 15, 116, ">", 20)
  end
  
  local controls_text = "Left/Right: Move   Z: Jump"
  local controls_x = 240 - string.len(controls_text) * 3
  rf.print_xy(controls_x, 140, controls_text, 14)  -- Light blue
  
  local menu_text = "Up/Down: Select   Z/X: Confirm"
  local menu_x = 240 - string.len(menu_text) * 3
  rf.print_xy(menu_x, 156, menu_text, 14)
  
  if rf.is_multiplayer() then
    local mp_text = "MULTIPLAYER MODE"
    local mp_x = 240 - string.len(mp_text) * 3
    rf.print_xy(mp_x, 176, mp_text, 17)  -- Purple
  end
end

function draw_background()
  -- SNES-style gradient sky - use blues from SNES palette
  -- SNES has nice blues around indices 35-37 (blue hues)
  rf.rectfill(0, 0, 479, 50, 37)   -- Dark blue (SNES shadow blue)
  rf.rectfill(0, 50, 479, 100, 36) -- Medium blue (SNES base blue)
  rf.rectfill(0, 100, 479, 135, 35) -- Light blue (SNES highlight blue)
  
  -- Draw SNES-style clouds (more defined)
  for i = 1, 4 do
    local cloud_x = (i * 120) % 480
    local cloud_y = 25 + (i * 15) % 35
    -- Main cloud body (white/light)
    rf.circfill(cloud_x, cloud_y, 10, 1)
    rf.circfill(cloud_x + 6, cloud_y, 8, 1)
    rf.circfill(cloud_x - 6, cloud_y, 8, 1)
    rf.circfill(cloud_x, cloud_y - 3, 7, 1)
    -- Cloud shadow (darker)
    rf.circfill(cloud_x + 2, cloud_y + 2, 6, 37)
  end
end

function draw_player(p, id)
  -- SNES-style player character - more defined and colorful
  local x, y = math.floor(p.x), math.floor(p.y)
  local color = p.color or 2
  
  -- Body base (larger, more rounded)
  rf.circfill(x + 8, y + 14, 7, color)
  rf.circfill(x + 8, y + 9, 6, color)
  rf.circfill(x + 8, y + 6, 5, color)
  
  -- Shadow/shading for depth (darker shade)
  rf.circfill(x + 6, y + 14, 5, color + 2)  -- Left shadow
  
  -- Eyes (white first, then black pupils)
  rf.circfill(x + 7, y + 6, 2, 1)  -- Left eye white
  rf.circfill(x + 9, y + 6, 2, 1)  -- Right eye white
  
  local eye_offset = 0
  if p.vx > 0.5 then eye_offset = 1  -- Looking right
  elseif p.vx < -0.5 then eye_offset = -1  -- Looking left
  end
  rf.circfill(x + 7 + eye_offset, y + 6, 1, 0)  -- Left pupil
  rf.circfill(x + 9 + eye_offset, y + 6, 1, 0)  -- Right pupil
  
  -- Highlight on top (bright highlight color)
  rf.circfill(x + 7, y + 4, 2, color - 3)
  
  -- Feet when on ground
  if p.on_ground then
    rf.rectfill(x + 4, y + 16, x + 6, y + 17, color + 2)  -- Left foot
    rf.rectfill(x + 10, y + 16, x + 12, y + 17, color + 2)  -- Right foot
  end
end

function draw_platform(plat)
  local color = plat.color or 38
  local highlight = math.max(2, color - 3)  -- Lighter shade (clamp to valid range)
  local shadow = math.min(49, color + 2)     -- Darker shade (clamp to valid range)
  
  -- Main platform body (thinner now - h=8)
  rf.rectfill(plat.x, plat.y, plat.x + plat.w, plat.y + plat.h, color)
  
  -- Top highlight stripe (bright edge)
  rf.rectfill(plat.x, plat.y, plat.x + plat.w, plat.y + 1, highlight)
  
  -- Left/right edge highlights (thin)
  rf.rectfill(plat.x, plat.y, plat.x + 1, plat.y + plat.h, highlight)
  rf.rectfill(plat.x + plat.w - 1, plat.y, plat.x + plat.w, plat.y + plat.h, highlight)
  
  -- Bottom shadow for depth
  rf.rectfill(plat.x + 1, plat.y + plat.h, plat.x + plat.w + 1, plat.y + plat.h + 1, shadow)
end

function _DRAW()
  -- Both host and non-host render the same way
  rf.clear_i(0)  -- Clear to black
  
  if state == "menu" then
    draw_menu()
    return
  end
  
  -- Check if we're in spectator mode
  local my_id = 1
  if rf.is_multiplayer() then
    my_id = rf.my_player_id()
  end
  local my_player = players[my_id]
  local is_spectating = my_player and not my_player.alive
  
  -- Note: Camera would be implemented here, but for now we draw relative to screen
  -- Camera follows spectated player or own player
  
  -- Draw background
  draw_background()
  
  -- Draw platforms with SNES-style visuals (thinner now)
  for _, plat in ipairs(platforms) do
    draw_platform(plat)
  end
  
  -- Draw all players with SNES-style graphics
  for id, p in pairs(players) do
    if p.alive then
      draw_player(p, id)
      
      -- Highlight spectated player
      if is_spectating and id == spectator_target then
        -- Draw arrow above spectated player
        local arrow_x = p.x + 8
        local arrow_y = p.y - 12
        rf.print("^", arrow_x - 3, arrow_y, 20)  -- Yellow arrow
        rf.line(arrow_x - 3, arrow_y + 5, arrow_x, arrow_y + 9, 20)
        rf.line(arrow_x + 3, arrow_y + 5, arrow_x, arrow_y + 9, 20)
      end
    else
      -- Draw dead player as faded/ghost
      local x, y = math.floor(p.x), math.floor(p.y)
      rf.circfill(x + 8, y + 14, 7, 37)  -- Faded blue
      rf.circfill(x + 8, y + 9, 6, 37)
      rf.circfill(x + 8, y + 6, 5, 37)
      rf.print("X", p.x + 4, p.y - 8, 2)  -- X mark
    end
    
    -- Show player ID above sprite with outline effect
    local my_display_id = 1
    if rf.is_multiplayer() then
      my_display_id = rf.my_player_id()
    end
    local text_color = (id == my_display_id) and 20 or 1
    if not p.alive then text_color = 2 end  -- Red for dead
    if is_spectating and id == spectator_target then text_color = 20 end  -- Yellow for spectated
    local outline_color = 0
    -- Outline text for visibility
    for ox = -1, 1 do
      for oy = -1, 1 do
        if ox ~= 0 or oy ~= 0 then
          rf.print(tostring(id), p.x + 4 + ox, p.y - 8 + oy, outline_color)
        end
      end
    end
    rf.print(tostring(id), p.x + 4, p.y - 8, text_color)
  end
  
  -- Draw spectator UI if dead
  if is_spectating then
    -- Spectator message at top
    local spec_text = "SPECTATING"
    if spectator_target then
      spec_text = "SPECTATING P" .. tostring(spectator_target)
    end
    rf.print_xy(240 - string.len(spec_text) * 3, 20, spec_text, 14)  -- Blue
    rf.print_xy(240 - string.len("UP+X: Cycle") * 3, 35, "UP+X: Cycle", 2)  -- Red
    
    -- EXIT button at bottom center
    local exit_text = "EXIT"
    local exit_x = 240 - string.len(exit_text) * 3
    local exit_y = 250
    local exit_w = string.len(exit_text) * 6 + 8
    local exit_h = 12
    
    -- Check if hovering/pressing exit button
    local exit_active = rf.btn(3) and rf.btnp(5)  -- Down + X pressed
    
    -- Button background (highlighted if active)
    local bg_color = exit_active and 35 or 37  -- Light blue if active, dark blue otherwise
    rf.rectfill(exit_x - 4, exit_y - 2, exit_x + exit_w - 4, exit_y + exit_h - 2, bg_color)
    rf.rect(exit_x - 4, exit_y - 2, exit_x + exit_w - 4, exit_y + exit_h - 2, 20)  -- Yellow border
    
    -- Button text
    local exit_color = exit_active and 20 or 1  -- Yellow if active, white otherwise
    rf.print_xy(exit_x, exit_y + 2, exit_text, exit_color)
    
    -- Instructions below button
    rf.print_xy(exit_x - 25, exit_y + 14, "DOWN+X: Exit", 14)
  end
  
  -- Draw goal line at top (win condition)
  rf.rectfill(0, 0, 480, 3, 20)  -- Yellow goal line
  
  -- Draw UI panel with SNES-style border
  rf.rectfill(0, 0, 125, 65, 37)  -- Dark blue background
  rf.rect(0, 0, 125, 65, 35)      -- Light blue border
  rf.rectfill(2, 2, 123, 63, 0)  -- Black inner
  
  -- Draw scores with better styling
  for id, s in pairs(score) do
    local score_y = 8 + (id-1) * 12
    local player_alive = players[id] and players[id].alive
    local label_color = player_alive and 20 or 2  -- Yellow if alive, red if dead
    rf.print("P" .. tostring(id), 8, score_y, label_color)  -- Colored label
    rf.print(": " .. tostring(s), 20, score_y, 1)   -- White score
  end
  
  -- Show game over / winner message
  if game_over and winner then
    local win_text = "PLAYER " .. tostring(winner) .. " WINS!"
    local win_x = 240 - string.len(win_text) * 3
    rf.print_xy(win_x, 120, win_text, 20)  -- Yellow winner text
    rf.print_xy(win_x - 50, 140, "PRESS Z/X TO RESTART", 14)  -- Blue instruction
  end
  
  -- Show connection info at bottom with SNES styling
  local info_y = 248
  if rf.is_multiplayer() then
    rf.print("Players: " .. tostring(rf.player_count()), 10, info_y, 14)  -- Light blue
    if rf.is_host() then
      rf.print("HOST", 10, info_y + 10, 20)  -- Yellow
    else
      local display_id = rf.my_player_id()
      rf.print("PLAYER " .. tostring(display_id), 10, info_y + 10, 17)  -- Purple
    end
  else
    rf.print("SOLO MODE", 10, info_y, 23)  -- Green
    if is_spectating then
      rf.print("DEAD - Use QUIT to exit", 10, info_y + 10, 2)  -- Red
    end
  end
  
  -- Show goal instruction
  if not game_over then
    rf.print("REACH TOP!", 360, 5, 20)  -- Yellow instruction at top right
    rf.print("JUMP ON HEADS", 340, 15, 2)  -- Red warning about head jumps
  end
end

function check_collision(player, platform)
  -- Simple AABB collision (adjusted for thinner platforms)
  return player.x < platform.x + platform.w and
         player.x + 16 > platform.x and
         player.y < platform.y + platform.h and
         player.y + 16 > platform.y and
         player.vy >= 0  -- Only collide when falling
end

