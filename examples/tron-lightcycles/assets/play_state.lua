-- Play State Module for Tron Light Cycles

-- State-local variables
local move_sound_playing = false
local boost_sound_playing = false
local boost_amount = 100.0  -- 0-100%
local is_boosting = false
local player_move_timer = 0.0  -- Separate timer for player movement when boosting

function _INIT()
  -- Module initialization
end

function _ENTER()
  -- Reset state when entering
  move_sound_playing = false
  boost_sound_playing = false
  boost_amount = 100.0  -- Start with full boost
  is_boosting = false
  countdown = 3.0
  move_timer = 0.0
  player_move_timer = 0.0
end

function _HANDLE_INPUT()
  -- ============================================================================
  -- STANDARD INPUT HANDLER TEMPLATE
  -- ============================================================================
  -- Universal 11-Button Input System
  -- Lua receives ONLY button indices (0-10), not keys.
  -- Key mapping happens at engine level:
  --   - WASM: Controller sends button indices directly
  --   - Desktop: Engine maps keyboard to button indices
  --
  -- Button Index Reference:
  --   0 = SELECT
  --   1 = START
  --   2 = UP
  --   3 = DOWN
  --   4 = LEFT
  --   5 = RIGHT
  --   6 = A
  --   7 = B
  --   8 = X
  --   9 = Y
  --  10 = TURBO
  -- ============================================================================
  
  -- Button 0: SELECT
  if countdown <= 0 and rf.btnp(0) then
    -- SELECT pauses the game
    rf.sfx("select")
    game.pushState("pause")
    return
  end
  
  -- Button 1: START
  -- Not used in play state
  
  -- Button 2: UP
  -- Movement: Ignore backwards movement (prevent turning 180 degrees)
  if player and player.alive and rf.btnp(2) then
    if player.dir ~= DIR_DOWN then
      player.dir = DIR_UP
      rf.sfx("turn")
    end
  end
  
  -- Button 3: DOWN
  -- Movement: Ignore backwards movement (prevent turning 180 degrees)
  if player and player.alive and rf.btnp(3) then
    if player.dir ~= DIR_UP then
      player.dir = DIR_DOWN
      rf.sfx("turn")
    end
  end
  
  -- Button 4: LEFT
  -- Movement: Ignore backwards movement (prevent turning 180 degrees)
  if player and player.alive and rf.btnp(4) then
    if player.dir ~= DIR_RIGHT then
      player.dir = DIR_LEFT
      rf.sfx("turn")
    end
  end
  
  -- Button 5: RIGHT
  -- Movement: Ignore backwards movement (prevent turning 180 degrees)
  if player and player.alive and rf.btnp(5) then
    if player.dir ~= DIR_LEFT then
      player.dir = DIR_RIGHT
      rf.sfx("turn")
    end
  end
  
  -- Button 6: A
  -- Not used in play state
  
  -- Button 7: B
  -- Not used in play state
  
  -- Button 8: X
  -- Not used in play state
  
  -- Button 9: Y
  -- Not used in play state
  
  -- Button 10: TURBO
  -- Handled in _UPDATE() for boost system
end

local function update_countdown(dt)
  local prev_int = math.ceil(countdown)
  countdown = math.max(0, countdown - dt)
  local now_int = math.ceil(countdown)
  if now_int < prev_int and now_int > 0 then
    rf.tone(600, 0.1, 0.3)
  end
end

function _UPDATE(dt)
  -- Handle countdown
  if countdown > 0 then
    update_countdown(dt)
    return
  end
  
  -- Update boost system
  local shift_pressed = false  -- Will be set via a new API or button check
  -- For now, we'll need to add SHIFT support - check if button 4 (Z) is held as shift
  -- Actually, let's use a dedicated check - we'll map SHIFT separately
  -- Since we don't have SHIFT mapped yet, let's use button 4 held as boost for now
  -- TODO: Add proper SHIFT key support
  
  -- Check for boost activation (SHIFT key held)
  -- We'll need to add rf.shift() function, but for now let's use button 4 as placeholder
  -- Actually, let me check if we can detect shift via modifier keys
  -- For now, I'll add a workaround: use button 4 + movement direction as boost
  
  -- Boost regeneration (when not boosting)
  if not is_boosting then
    -- Regenerate 10% per second (0-100% over 10 seconds)
    boost_amount = math.min(100.0, boost_amount + (10.0 * dt))
  end
  
  -- Boost: TURBO (10) increases speed by 50%
  local turbo_pressed = rf.btn(10)  -- TURBO (button 10)
  local wants_boost = turbo_pressed and player and player.alive and boost_amount > 0.0 and countdown <= 0
  
  -- Handle boost state
  if wants_boost and not is_boosting then
    -- Start boosting
    is_boosting = true
    -- Stop move sound, start boost sound
    if move_sound_playing then
      rf.sfx("move", "off")
      move_sound_playing = false
    end
    if not boost_sound_playing then
      rf.sfx("boost", "on")
      boost_sound_playing = true
    end
  elseif not wants_boost and is_boosting then
    -- Stop boosting
    is_boosting = false
    -- Stop boost sound, start move sound
    if boost_sound_playing then
      rf.sfx("boost", "off")
      boost_sound_playing = false
    end
    if player and player.alive and not move_sound_playing then
      rf.sfx("move", "on")
      move_sound_playing = true
    end
  end
  
  -- Consume boost while boosting (100% lasts 5 seconds)
  if is_boosting then
    local boost_depletion_rate = 20.0  -- 20% per second (100% / 5 seconds)
    boost_amount = math.max(0.0, boost_amount - (boost_depletion_rate * dt))
    if boost_amount <= 0.0 then
      -- Boost depleted, stop boosting
      is_boosting = false
      if boost_sound_playing then
        rf.sfx("boost", "off")
        boost_sound_playing = false
      end
      if player and player.alive and not move_sound_playing then
        rf.sfx("move", "on")
        move_sound_playing = true
      end
    end
  end
  
  -- Handle player alive state for sound
  if player and player.alive then
    if not is_boosting then
      -- Start continuous motorcycle sound if not already playing
      if not move_sound_playing then
        rf.sfx("move", "on")
        move_sound_playing = true
      end
    end
  else
    -- Stop all sounds if player is dead
    if move_sound_playing then
      rf.sfx("move", "off")
      move_sound_playing = false
    end
    if boost_sound_playing then
      rf.sfx("boost", "off")
      boost_sound_playing = false
    end
    is_boosting = false
  end
  
  -- Calculate effective speed (boost increases by 50%)
  -- Note: speed is in moves per second, so we need to scale the timer update accordingly
  local effective_speed = speed
  if is_boosting then
    effective_speed = speed * 1.5  -- 50% increase (e.g., 5.0 -> 7.5 moves/sec)
  end
  
  -- Update player move timer separately when boosting (allows player to move faster than enemies)
  if is_boosting and player and player.alive then
    player_move_timer = player_move_timer + dt * effective_speed
    
    -- Move player more frequently when boosting
    while player_move_timer >= 1.0 do
      player_move_timer = player_move_timer - 1.0
      if not move_cycle(player) then
        -- Player crashed - stop everything and transition to gameover
        rf.sfx("move", "off") -- Stop motorcycle sound
        if boost_sound_playing then
          rf.sfx("boost", "off")
        end
        move_sound_playing = false
        boost_sound_playing = false
        rf.sfx("crash")
        rf.sfx("stopall") -- Stop all sounds before transition
        best_level = math.max(best_level, level)
        rf.music("taps") -- Play taps for loss
        -- Change state immediately - this will exit play state and enter gameover
        game.changeState("gameover")
        return -- Exit update early since we're changing states
      end
    end
  end
  
  -- Update shared move timer (for normal movement and enemy movement)
  -- When boosting, player moves separately, so this only affects enemies
  move_timer = move_timer + dt * speed
  
  if move_timer >= 1.0 then
    move_timer = move_timer - 1.0
    
    -- Update AI for enemies
    for i=1,#enemies do
      if enemies[i].alive then
        update_enemy_ai(enemies[i])
      end
    end
    
    -- Move player (only if not boosting - when boosting, player moves via player_move_timer)
    if not is_boosting and player and player.alive then
      if not move_cycle(player) then
        -- Player crashed - stop everything and transition to gameover
        rf.sfx("move", "off") -- Stop motorcycle sound
        if boost_sound_playing then
          rf.sfx("boost", "off")
        end
        move_sound_playing = false
        boost_sound_playing = false
        rf.sfx("crash")
        rf.sfx("stopall") -- Stop all sounds before transition
        best_level = math.max(best_level, level)
        rf.music("taps") -- Play taps for loss
        -- Change state immediately - this will exit play state and enter gameover
        game.changeState("gameover")
        return -- Exit update early since we're changing states
      end
    end
    
    -- Move enemies
    for i=1,#enemies do
      if enemies[i].alive then
        if not move_cycle(enemies[i]) then
          -- Enemy crashed
          rf.sfx("crash")
        end
      end
    end
    
    -- Check if player won (all enemies dead)
    local all_enemies_dead = true
    for i=1,#enemies do
      if enemies[i].alive then
        all_enemies_dead = false
        break
      end
    end
    
    if player and player.alive and all_enemies_dead then
      rf.sfx("move", "off") -- Stop motorcycle sound
      if boost_sound_playing then
        rf.sfx("boost", "off")
        boost_sound_playing = false
      end
      move_sound_playing = false
      is_boosting = false
      boost_amount = 100.0  -- Reset boost for next level
      rf.sfx("won") -- Victory sound
      score = score + level * 100
      level = level + 1
      init_level(level)
      countdown = 2.0
    end
  end
end

function _DRAW()
  -- Clear screen
  rf.clear_i(COLOR_BLACK)
  
  -- Draw grid background (subtle darker dots, closer together)
  for y=0,GRID_H-1,3 do
    for x=0,GRID_W-1,3 do
      local sx, sy = grid_to_screen(x, y)
      rf.pset(sx, sy, 22) -- Dark gray (50% darker than white, index 22 in Super Mario 50 palette)
    end
  end
  
  -- Draw enemy cycles
  for i=1,#enemies do
    draw_cycle(enemies[i], false)
  end
  
  -- Draw player cycle
  if player then
    draw_cycle(player, true)
  end
  
  -- Draw HUD
  draw_hud()
  
  -- Draw boost label and bar at top center
  local boost_text = "BOOST"
  local boost_text_x = 240 - string.len(boost_text) * 3
  rf.print_xy(boost_text_x, 2, boost_text, COLOR_WHITE)
  
  -- Draw boost bar (smaller, centered under text)
  local bar_width = 120
  local bar_height = 4
  local bar_x = 240 - bar_width / 2
  local bar_y = 12
  
  -- Draw background bar (gray)
  rf.rectfill(bar_x, bar_y, bar_x + bar_width, bar_y + bar_height, COLOR_GRAY)
  
  -- Draw boost fill (cyan when boosting, blue when not)
  local fill_width = math.floor((boost_amount / 100.0) * bar_width)
  local boost_color = is_boosting and COLOR_CYAN or COLOR_BLUE
  if fill_width > 0 then
    rf.rectfill(bar_x, bar_y, bar_x + fill_width, bar_y + bar_height, boost_color)
  end
  
  -- Draw border
  rf.rect(bar_x - 1, bar_y - 1, bar_x + bar_width + 1, bar_y + bar_height + 1, COLOR_WHITE)
  
  -- Draw countdown
  if countdown > 0 then
    local level_text = "LEVEL " .. tostring(level)
    local level_x = 240 - string.len(level_text)*3
    rf.print_xy(level_x, 120, level_text, COLOR_WHITE)
    if math.ceil(countdown) > 0 then
      local ready_text = "GET READY: " .. tostring(math.ceil(countdown))
      local ready_x = 240 - string.len(ready_text)*3
      rf.print_xy(ready_x, 140, ready_text, COLOR_WHITE)
    else
      local go_x = 240 - string.len("GO!")*3
      rf.print_xy(go_x, 140, "GO!", COLOR_WHITE)
    end
  end
end

function _EXIT()
  -- Stop all audio when exiting play state
  rf.sfx("stopall")
  if move_sound_playing then
    move_sound_playing = false
  end
end

function _DONE()
  -- Shutdown cleanup
end

