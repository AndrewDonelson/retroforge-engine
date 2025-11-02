-- Play State Module for Moon Lander

-- State-local variables (persist across enter/exit)
local time_s = 0
local prevThrust = false

function _INIT()
  -- Module initialization
end

function _ENTER()
  -- Reset state when entering
  time_s = 0
  prevThrust = false
  -- Ship and level are set via set_level() when game starts
end

function _HANDLE_INPUT()
  -- Handle pause (ESC/button 3/Down) - push pause state overlay
  if rf.btnp(3) then
    rf.sfx("move")
    game.pushState("pause")
    return
  end
  
  -- Normal game input (when ship alive)
  if not ship.alive then return end
  
  -- Ship controls are handled in _UPDATE
end

local function normalize_angle()
  if ship.angle > math.pi then ship.angle = ship.angle - TWO_PI end
  if ship.angle < -math.pi then ship.angle = ship.angle + TWO_PI end
end

function _UPDATE(dt)
  time_s = time_s + dt
  
  -- Handle crash sequence
  if not ship.alive then
    update_crash(dt)
    return
  end
  
  -- Note: Pause is handled by pause_state overlay via PushState/PopState
  -- When pause state is on top of stack, play state's Update won't be called
  
  -- Handle countdown
  if countdown > 0 then
    update_countdown(dt)
    return
  end
  
  -- Update gameplay
  update_play(dt)
end

function update_countdown(dt)
  local prevInt = math.ceil(countdown)
  countdown = math.max(0, countdown - dt)
  local nowInt = math.ceil(countdown)
  if nowInt < prevInt then rf.tone(600, 0.06, 0.3) end
  ship.y = ship.y + math.sin(time_s*2)*0.2
end

function update_crash(dt)
  crash_timer = crash_timer + dt
  if crash_phase == "crash" then
    if crash_timer >= 0.25 then crash_phase = "taps"; crash_timer = 0 end
  elseif crash_phase == "taps" then
    if not taps_started then
      rf.music("taps")
      taps_total_dur = 4.0
      taps_started = true
    end
    if crash_timer >= taps_total_dur or crash_timer >= 5.0 then
      game.changeState("menu")
      crash_phase = nil
    end
  else
    if crash_timer >= 5.0 then
      game.changeState("menu")
    end
  end
end

function update_play(dt)
  -- Ship controls
  if rf.btn(0) then ship.angle = ship.angle + ROT*dt end
  if rf.btn(1) then ship.angle = ship.angle - ROT*dt end
  normalize_angle()
  
  local thrust = rf.btn(2) and ship.fuel > 0
  if thrust then
    local ax = math.sin(ship.angle) * THRUST
    local ay = -math.cos(ship.angle) * THRUST
    ship.vx = ship.vx + ax*dt
    ship.vy = ship.vy + ay*dt
    ship.fuel = clamp(ship.fuel - 25*dt, 0, 999)
  end
  
  if thrust ~= prevThrust then
    rf.sfx("thrust", thrust and "on" or "off")
    prevThrust = thrust
  end
  ship.thrusting = thrust
  
  -- Gravity
  ship.vy = ship.vy + G*dt
  if ship.vy > 300 then ship.vy = 300 end
  
  -- Friction
  ship.vx = ship.vx * (1 - 0.5*dt)
  
  -- Update position
  ship.x = ship.x + ship.vx
  ship.y = ship.y + ship.vy
  
  -- Screen boundaries
  if ship.x < 0 then ship.x = 0; ship.vx = 0
  elseif ship.x > 479 then ship.x = 479; ship.vx = 0
  end
  
  -- Ground collision
  local ground_here = terrain_y(math.floor(ship.x))
  if ship.y >= (ground_here - ship.size) then
    ship.y = ground_here - ship.size
    local speed = math.sqrt(ship.vx*ship.vx + ship.vy*ship.vy)
    local vy_abs = math.abs(ship.vy)
    local angle_ok = math.abs(ship.angle) < 0.2
    local on_pad = ship.x >= pad_x0 and ship.x <= pad_x1
    
    if vy_abs < (land_speed or 18) * 0.6 and speed < (land_speed or 18) and angle_ok and on_pad then
      -- Successful landing
      ship.landed = true
      score = score + math.floor(100 + (ship.fuel*2) + math.max(0, (land_speed - speed)*5))
      countdown = 3.0
      best_level = math.max(best_level, level)
      best_score = math.max(best_score, score)
      rf.sfx("land")
      set_level(level + 1)
    else
      -- Crash
      rf.sfx("thrust", "off")
      rf.sfx("stopall")
      ship.alive = false
      crash_timer = 0
      start_crash_sequence()
      rf.sfx("crash")
    end
    ship.vx, ship.vy = 0, 0
  end
end

function draw_stars()
  for i = 1, #stars do
    local s = stars[i]
    local brightness = math.floor(128 + 127 * math.sin(time_s * s.speed + s.phase))
    brightness = math.max(128, math.min(255, brightness))
    rf.pset(s.x, s.y, COLOR_STARS)
  end
end

function draw_ship()
  local s = ship.size
  local sin, cos = math.sin(ship.angle), math.cos(ship.angle)
  local p0x, p0y = 0, -s
  local p1x, p1y = -s*0.7, s
  local p2x, p2y = s*0.7, s
  local function tx(x, y)
    return math.floor(ship.x + x*cos + y*sin + 0.5), math.floor(ship.y - x*sin + y*cos + 0.5)
  end
  local a1x, a1y = tx(p0x, p0y)
  local a2x, a2y = tx(p1x, p1y)
  local a3x, a3y = tx(p2x, p2y)
  rf.line(a1x, a1y, a2x, a2y, COLOR_WHITE)
  rf.line(a2x, a2y, a3x, a3y, COLOR_WHITE)
  rf.line(a3x, a3y, a1x, a1y, COLOR_WHITE)
  if ship.thrusting then
    local fx, fy = tx(0, s+2)
    rf.circfill(fx, fy, 2, COLOR_LIGHT_GRAY)
  end
end

function draw_hud()
  rf.print_anchored("MOON LANDER", "topcenter", COLOR_LIGHT_GRAY)
  rf.print_xy(2, 2, "ALT:"..string.format("%3.0f", ground_y-ship.y), COLOR_LIGHT_GRAY)
  rf.print_xy(2, 10, "VY :"..string.format("%4.1f", ship.vy), COLOR_LIGHT_GRAY)
  rf.print_xy(2, 18, "FUEL:"..string.format("%3.0f", ship.fuel), COLOR_LIGHT_GRAY)
  rf.print_xy(380, 2, "LVL:"..tostring(level), COLOR_LIGHT_GRAY)
  rf.print_xy(380, 10, "SCORE:"..tostring(score), COLOR_LIGHT_GRAY)
end

function draw_level()
  local prevx, prevy = 0, terrain_y(0)
  for x = 0, 479, 2 do
    local ty = terrain_y(x)
    rf.rectfill(x, ty, x, 269, COLOR_GRAY)
    rf.line(prevx, prevy, x, ty, COLOR_GRAY)
    prevx, prevy = x, ty
  end
  rf.rectfill(pad_x0, pad_y, pad_x1, pad_y+2, COLOR_LANDINGPAD)
end

function draw_messages()
  if countdown > 0 and ship.alive and not ship.landed then
    rf.print_anchored("GET READY:"..tostring(math.ceil(countdown)), "middlecenter", COLOR_WHITE)
    return
  end
  if ship.landed then
    rf.print_anchored("LANDED!", "middlecenter", COLOR_WHITE)
    return
  end
  if not ship.alive and not ship.landed then
    rf.print_xy(240 - string.len("CRASHED")*3, 120, "CRASHED", COLOR_WHITE)
    rf.print_xy(240 - string.len("RETURNING TO MENU...")*3, 140, "RETURNING TO MENU...", COLOR_GRAY)
  end
end

function _DRAW()
  -- Clear screen
  rf.clear_i(COLOR_BLACK)
  
  -- Draw game scene
  -- Note: When pause state is active, it will be drawn on top via state stack
  draw_stars()
  draw_hud()
  draw_level()
  draw_ship()
  draw_messages()
end

function _EXIT()
  -- Cleanup when leaving
end

function _DONE()
  -- Shutdown cleanup
end

