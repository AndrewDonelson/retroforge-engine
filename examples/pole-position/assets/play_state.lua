-- Play State - Top-down racing gameplay

-- Game state
local player_car = {
  x = 0.0,  -- Horizontal position (screen coordinates)
  lane = 2,  -- Current lane (1-4, with 2-3 being center lanes)
  speed = 0.0,  -- Vertical speed (forward)
  max_speed = 6.0,
  acceleration = 0.3,
  deceleration = 0.2,
  spin_timer = 0.0,  -- Spin-out timer (0 = normal, >0 = spinning)
  air_timer = 0.0,  -- Air timer (0 = on ground, >0 = in air)
  air_angle = 0.0,  -- Rotation angle when in air
  crashed = false
}

local track_scroll_x = 0.0  -- Track horizontal scroll position
local track = nil
local traffic = {}
local obstacles = {}  -- Oil slicks and potholes
local countdown = 3.0
local engine_sound_playing = false
local lap_start_time = 0.0
local lap_number = 1
local total_laps = 3

-- Lane positions (centered on screen)
local LANE_WIDTH = 80
local LANE_1_X = 120  -- Left lane
local LANE_2_X = 200  -- Left-center
local LANE_3_X = 280  -- Right-center
local LANE_4_X = 360  -- Right lane

function _INIT()
end

function _ENTER()
  -- Initialize player
  player_car.x = LANE_2_X  -- Start in left-center lane
  player_car.lane = 2
  player_car.speed = 0.0
  player_car.spin_timer = 0.0
  player_car.air_timer = 0.0
  player_car.air_angle = 0.0
  player_car.crashed = false
  
  -- Ensure track is selected (default to first track if not set)
  if not selected_track or selected_track < 1 or selected_track > #tracks then
    selected_track = 1
  end
  
  -- Load selected track
  track = tracks[selected_track]
  if not track.segments then
    generate_track(selected_track)
    track = tracks[selected_track]
  end
  
  -- Reset track scroll
  track_scroll_x = 0.0
  
  -- Initialize traffic (spawn ahead of player)
  traffic = {}
  for i = 1, 8 do
    local car = {
      x = 0,  -- Will be set based on lane
      y = 100 + i * 80,  -- Ahead of player
      lane = math.random(1, 4),  -- Random lane
      speed = 2.0 + math.random() * 2.0,  -- Random speed
      sprite_name = "f1_car_" .. ((i % 3) + 2)  -- Use f1_car_2, f1_car_3, or f1_car_4
    }
    car.x = get_lane_x(car.lane)
    traffic[i] = car
  end
  
  -- Initialize obstacles (oil slicks and potholes)
  obstacles = {}
  for i = 1, 15 do
    local obstacle = {
      x = get_lane_x(math.random(1, 4)),
      y = 150 + i * 120,  -- Ahead of player
      type = math.random() < 0.5 and "oil" or "pothole"
    }
    obstacles[i] = obstacle
  end
  
  -- Reset race state
  countdown = 3.0
  race_started = false
  race_finished = false
  lap_time = 0.0
  lap_start_time = 0.0
  lap_number = 1
  current_position = 1
  
  -- Stop all audio
  rf.sfx("stopall")
  engine_sound_playing = false
end

function get_lane_x(lane)
  if lane == 1 then return LANE_1_X
  elseif lane == 2 then return LANE_2_X
  elseif lane == 3 then return LANE_3_X
  elseif lane == 4 then return LANE_4_X
  else return 240  -- Center
  end
end

function _HANDLE_INPUT()
  -- Button 0: SELECT - Pause or Restart if crashed
  if rf.btnp(0) then
    if player_car.crashed then
      -- Restart race
      rf.sfx("select")
      _ENTER()  -- Re-initialize
      return
    elseif race_started and not race_finished then
      rf.sfx("select")
      game.pushState("pause")
      return
    end
  end
  
  -- Button 1: START - Not used in play
  
  -- Button 2: UP - Accelerate
  if race_started and not race_finished and not player_car.crashed and player_car.spin_timer <= 0 and player_car.air_timer <= 0 then
    if rf.btn(2) then
      player_car.speed = math.min(player_car.max_speed, 
                                  player_car.speed + player_car.acceleration)
    end
  end
  
  -- Button 3: DOWN - Brake
  if race_started and not race_finished and not player_car.crashed and player_car.spin_timer <= 0 and player_car.air_timer <= 0 then
    if rf.btn(3) then
      player_car.speed = math.max(0.0, 
                                  player_car.speed - player_car.deceleration * 2)
    end
  end
  
  -- Button 4: LEFT - Change lane left
  if race_started and not race_finished and not player_car.crashed and player_car.spin_timer <= 0 and player_car.air_timer <= 0 then
    if rf.btnp(4) then
      player_car.lane = math.max(1, player_car.lane - 1)
      player_car.x = get_lane_x(player_car.lane)
    end
  end
  
  -- Button 5: RIGHT - Change lane right
  if race_started and not race_finished and not player_car.crashed and player_car.spin_timer <= 0 and player_car.air_timer <= 0 then
    if rf.btnp(5) then
      player_car.lane = math.min(4, player_car.lane + 1)
      player_car.x = get_lane_x(player_car.lane)
    end
  end
end

function _UPDATE(dt)
  -- Countdown
  if countdown > 0 then
    local prev_count = math.ceil(countdown)
    countdown = countdown - dt
    local now_count = math.ceil(countdown)
    
    -- Play countdown sound for each number
    if now_count < prev_count and now_count > 0 then
      rf.sfx("select")
    end
    
    if countdown <= 0 then
      race_started = true
      lap_start_time = 0.0
      rf.sfx("countdown")
      -- Start race music
      local music_track = "race_music_1"
      if selected_track % 3 == 2 then
        music_track = "race_music_2"
      elseif selected_track % 3 == 0 then
        music_track = "race_music_3"
      end
      rf.music(music_track)
    end
    return
  end
  
  if not race_started or race_finished or player_car.crashed then
    return
  end
  
  -- Update lap time
  lap_time = lap_time + dt
  lap_start_time = lap_start_time + dt
  
  -- Natural deceleration
  if player_car.speed > 0 then
    player_car.speed = math.max(0.0, player_car.speed - 0.05)
  end
  
  -- Update spin-out timer
  if player_car.spin_timer > 0 then
    player_car.spin_timer = player_car.spin_timer - dt
    if player_car.spin_timer <= 0 then
      player_car.spin_timer = 0.0
      -- Resume control
    end
  end
  
  -- Update air timer
  if player_car.air_timer > 0 then
    player_car.air_timer = player_car.air_timer - dt
    player_car.air_angle = player_car.air_angle + dt * 360 * 2  -- Rotate 2 full rotations per second
    
    if player_car.air_timer <= 0 then
      player_car.air_timer = 0.0
      -- Land facing random direction
      player_car.air_angle = math.random() * 360
      -- Resume control
    end
  end
  
  -- Move track horizontally (slowly) - left/right movement
  track_scroll_x = track_scroll_x + dt * 30  -- Slow horizontal scroll
  if track_scroll_x > 1000 then
    track_scroll_x = track_scroll_x - 1000  -- Wrap to prevent overflow
  end
  
  -- Move traffic cars forward
  for i = 1, #traffic do
    local car = traffic[i]
    car.y = car.y - car.speed * 10 * dt
    
    -- Respawn if behind player
    if car.y < -50 then
      car.y = 400 + math.random() * 200
      car.lane = math.random(1, 4)
      car.x = get_lane_x(car.lane)
    end
  end
  
  -- Move obstacles forward
  for i = 1, #obstacles do
    local obstacle = obstacles[i]
    obstacle.y = obstacle.y - player_car.speed * 10 * dt
    
    -- Respawn if behind player
    if obstacle.y < -50 then
      obstacle.y = 400 + math.random() * 200
      obstacle.x = get_lane_x(math.random(1, 4))
      obstacle.type = math.random() < 0.5 and "oil" or "pothole"
    end
  end
  
  -- Check collision with traffic
  for i = 1, #traffic do
    local car = traffic[i]
    local dist_y = math.abs(car.y - 220)  -- Player is at y=220
    local dist_x = math.abs(car.x - player_car.x)
    
    if dist_y < 30 and dist_x < 30 then
      -- CRASH!
      rf.sfx("crash")
      player_car.crashed = true
      player_car.speed = 0.0
      break
    end
  end
  
  -- Check collision with obstacles
  if player_car.spin_timer <= 0 and player_car.air_timer <= 0 then
    for i = 1, #obstacles do
      local obstacle = obstacles[i]
      local dist_y = math.abs(obstacle.y - 220)
      local dist_x = math.abs(obstacle.x - player_car.x)
      
      if dist_y < 20 and dist_x < 40 then
        if obstacle.type == "oil" then
          -- Oil slick - spin out
          rf.sfx("crash")
          player_car.spin_timer = 1.5
          player_car.speed = math.max(0.0, player_car.speed - 2.0)
        elseif obstacle.type == "pothole" then
          -- Pothole - thrown in air
          rf.sfx("crash")
          player_car.air_timer = 1.0
          player_car.air_angle = math.random() * 360
          player_car.speed = math.max(0.0, player_car.speed - 1.0)
        end
        -- Remove obstacle after hitting
        obstacle.y = -100
        break
      end
    end
  end
  
  -- Engine sound
  if player_car.speed > 0.1 then
    if not engine_sound_playing then
      rf.sfx("engine", "on")
      engine_sound_playing = true
    end
  else
    if engine_sound_playing then
      rf.sfx("engine", "off")
      engine_sound_playing = false
    end
  end
  
  -- Check lap completion
  if lap_start_time >= 30.0 then  -- 30 seconds per lap
    lap_number = lap_number + 1
    lap_start_time = 0.0
    
    if lap_number > total_laps then
      race_finished = true
      rf.sfx("stopall")
      rf.sfx("checkpoint")
      game.changeState("results")
    else
      rf.sfx("checkpoint")
    end
  end
  
  -- Calculate position (simplified)
  current_position = 1
  for i = 1, #traffic do
    if traffic[i].y > 220 then
      current_position = current_position + 1
    end
  end
end

function _DRAW()
  rf.clear_i(COLOR_BLACK)
  
  -- Draw top-down track (road with lanes)
  draw_topdown_track()
  
  -- Draw obstacles
  draw_obstacles()
  
  -- Draw traffic
  draw_traffic_topdown()
  
  -- Draw player car
  draw_player_car_topdown()
  
  -- Draw HUD
  draw_hud()
  
  -- Countdown overlay
  if countdown > 0 then
    local cd_text = math.ceil(countdown)
    if cd_text == 0 then
      cd_text = "GO!"
    end
    local text = tostring(cd_text)
    local text_x = 240 - string.len(text) * 9
    rf.print_xy(text_x, 100, text, COLOR_YELLOW)
  end
  
  -- Crash message
  if player_car.crashed then
    local crash_text = "CRASH!"
    local crash_x = 240 - string.len(crash_text) * 9
    rf.print_xy(crash_x, 120, crash_text, COLOR_RED)
    local restart_text = "Press SELECT to restart"
    local restart_x = 240 - string.len(restart_text) * 3
    rf.print_xy(restart_x, 135, restart_text, COLOR_WHITE)
  end
end

function draw_topdown_track()
  -- Draw road surface (gray)
  rf.rectfill(0, 0, 480, 270, COLOR_GRAY)
  
  -- Draw lane dividers (white lines) - scroll horizontally
  local base_lane_xs = {LANE_1_X, LANE_2_X, LANE_3_X, LANE_4_X}
  for i = 1, #base_lane_xs do
    local base_x = base_lane_xs[i]
    -- Slight horizontal movement based on scroll
    local x = base_x + math.sin(track_scroll_x * 0.01) * 10  -- Slight sway
    -- Dashed line effect (vertical lines)
    for y = 0, 270, 20 do
      local line_y = (y + math.floor(track_scroll_x * 0.5)) % 270
      rf.rectfill(x - 1, line_y, x + 1, line_y + 10, COLOR_WHITE)
    end
  end
  
  -- Draw road edges (with slight movement)
  local edge_sway = math.sin(track_scroll_x * 0.01) * 5
  rf.rectfill(0, 0, 10 + edge_sway, 270, COLOR_WHITE)  -- Left edge
  rf.rectfill(470 - edge_sway, 0, 480, 270, COLOR_WHITE)  -- Right edge
end

function draw_obstacles()
  for i = 1, #obstacles do
    local obstacle = obstacles[i]
    if obstacle.y >= 0 and obstacle.y < 270 then
      if obstacle.type == "oil" then
        -- Draw oil slick (dark circle)
        rf.circfill(obstacle.x, obstacle.y, 15, COLOR_DARKEST_GRAY)
        rf.circ(obstacle.x, obstacle.y, 15, COLOR_BLACK)
      elseif obstacle.type == "pothole" then
        -- Draw pothole (dark rectangle)
        rf.rectfill(obstacle.x - 12, obstacle.y - 8, obstacle.x + 12, obstacle.y + 8, COLOR_BLACK)
        rf.rect(obstacle.x - 12, obstacle.y - 8, obstacle.x + 12, obstacle.y + 8, COLOR_DARKEST_GRAY)
      end
    end
  end
end

function draw_traffic_topdown()
  for i = 1, #traffic do
    local car = traffic[i]
    if car.y >= 0 and car.y < 270 then
      -- Draw car sprite (top-down, scaled smaller)
      draw_sprite_scaled(car.sprite_name, car.x, car.y, 0.3)
    end
  end
end

function draw_player_car_topdown()
  local x = player_car.x
  local y = 220  -- Fixed at bottom of screen
  
  -- Draw spinning effect if in spin-out
  if player_car.spin_timer > 0 then
    -- Draw multiple car images for spin effect
    local spin_angle = (1.5 - player_car.spin_timer) * 360 * 2
    for i = 0, 3 do
      local angle = spin_angle + i * 90
      local rad = (angle * math.pi) / 180
      local offset_x = math.cos(rad) * 5
      local offset_y = math.sin(rad) * 5
      draw_sprite_scaled("f1_car_1", x + offset_x, y + offset_y, 0.35)
    end
  -- Draw in air effect if thrown up
  elseif player_car.air_timer > 0 then
    local air_y = y - (player_car.air_timer * 30)  -- Move up
    -- Draw rotated car
    draw_sprite_scaled_rotated("f1_car_1", x, air_y, player_car.air_angle, 0.35)
  else
    -- Normal drawing
    draw_sprite_scaled("f1_car_1", x, y, 0.35)
  end
end

-- Draw scaled sprite (centered at x, y)
function draw_sprite_scaled(sprite_name, cx, cy, scale)
  local sprite = rf.sprite(sprite_name)
  if not sprite then return end
  
  local w = sprite.width
  local h = sprite.height
  local scaled_w = math.floor(w * scale)
  local scaled_h = math.floor(h * scale)
  
  -- Calculate top-left position (centered)
  local x = cx - scaled_w / 2
  local y = cy - scaled_h / 2
  
  -- Draw scaled pixels
  local x_scale = w / scaled_w
  local y_scale = h / scaled_h
  
  for dy = 0, scaled_h - 1 do
    for dx = 0, scaled_w - 1 do
      local src_x = math.floor(dx * x_scale)
      local src_y = math.floor(dy * y_scale)
      
      if src_x >= 0 and src_x < w and src_y >= 0 and src_y < h then
        local color_idx = sprite.pixels[src_y + 1][src_x + 1]  -- Lua is 1-indexed
        if color_idx >= 0 then  -- -1 is transparent
          rf.pset(math.floor(x + dx), math.floor(y + dy), color_idx)
        end
      end
    end
  end
end

-- Draw scaled and rotated sprite
function draw_sprite_scaled_rotated(sprite_name, cx, cy, angle, scale)
  local sprite = rf.sprite(sprite_name)
  if not sprite then return end
  
  local w = sprite.width
  local h = sprite.height
  local rad = (angle * math.pi) / 180
  local cos_a = math.cos(rad)
  local sin_a = math.sin(rad)
  
  local scaled_w = math.floor(w * scale)
  local scaled_h = math.floor(h * scale)
  
  -- Draw rotated pixels
  for sy = 0, scaled_h - 1 do
    for sx = 0, scaled_w - 1 do
      -- Rotate around center
      local px = sx - scaled_w / 2
      local py = sy - scaled_h / 2
      local rx = px * cos_a - py * sin_a
      local ry = px * sin_a + py * cos_a
      
      -- Map back to source
      local src_x = math.floor(rx / scale + w / 2)
      local src_y = math.floor(ry / scale + h / 2)
      
      if src_x >= 0 and src_x < w and src_y >= 0 and src_y < h then
        local color_idx = sprite.pixels[src_y + 1][src_x + 1]
        if color_idx >= 0 then
          rf.pset(math.floor(cx + px), math.floor(cy + py), color_idx)
        end
      end
    end
  end
end

function draw_hud()
  -- Top HUD
  -- Speed
  local speed_text = string.format("SPD: %d", math.floor(player_car.speed * 10))
  rf.print_xy(10, 10, speed_text, COLOR_WHITE)
  
  -- Lap time
  local time_text = string.format("TIME: %.1f", lap_time)
  rf.print_xy(10, 30, time_text, COLOR_WHITE)
  
  -- Lap counter
  local lap_text = string.format("LAP: %d/%d", lap_number, total_laps)
  local lap_x = 480 - string.len(lap_text) * 6 - 10
  rf.print_xy(lap_x, 10, lap_text, COLOR_WHITE)
  
  -- Position
  local pos_text = string.format("POS: %d/%d", current_position, total_cars)
  local pos_x = 480 - string.len(pos_text) * 6 - 10
  rf.print_xy(pos_x, 30, pos_text, COLOR_WHITE)
  
  -- Speed bar
  local bar_width = 100
  local bar_height = 10
  local bar_x = 10
  local bar_y = 50
  local speed_percent = player_car.speed / player_car.max_speed
  rf.rectfill(bar_x, bar_y, bar_x + bar_width, bar_y + bar_height, COLOR_DARK_GRAY)
  rf.rectfill(bar_x, bar_y, bar_x + bar_width * speed_percent, bar_y + bar_height, COLOR_GREEN)
  
  -- Spin-out indicator
  if player_car.spin_timer > 0 then
    rf.print_xy(10, 70, "SPIN OUT!", COLOR_RED)
  end
  
  -- Air indicator
  if player_car.air_timer > 0 then
    rf.print_xy(10, 70, "AIR!", COLOR_YELLOW)
  end
end

function _EXIT()
  rf.sfx("stopall")
  engine_sound_playing = false
end

function _DONE()
  -- Shutdown (unused)
end
