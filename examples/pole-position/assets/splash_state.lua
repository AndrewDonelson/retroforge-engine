-- Splash State - Animated title screen with scrolling road

local splash_time = 0
local splash_duration = 4.0
local music_started = false
local road_offset = 0
local car_x = 0

function _INIT()
end

function _ENTER()
  splash_time = 0
  music_started = false
  road_offset = 0
  car_x = -100
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
  
  -- Any button (0-10) skips splash
  -- Button 0: SELECT
  if rf.btnp(0) then
    rf.sfx("stopall") -- Stop all audio before transition
    if game then
      game.changeState("menu")
    end
    return
  end
  
  -- Button 1: START
  if rf.btnp(1) then
    rf.sfx("stopall") -- Stop all audio before transition
    if game then
      game.changeState("menu")
    end
    return
  end
  
  -- Button 2: UP
  if rf.btnp(2) then
    rf.sfx("stopall") -- Stop all audio before transition
    if game then
      game.changeState("menu")
    end
    return
  end
  
  -- Button 3: DOWN
  if rf.btnp(3) then
    rf.sfx("stopall") -- Stop all audio before transition
    if game then
      game.changeState("menu")
    end
    return
  end
  
  -- Button 4: LEFT
  if rf.btnp(4) then
    rf.sfx("stopall") -- Stop all audio before transition
    if game then
      game.changeState("menu")
    end
    return
  end
  
  -- Button 5: RIGHT
  if rf.btnp(5) then
    rf.sfx("stopall") -- Stop all audio before transition
    if game then
      game.changeState("menu")
    end
    return
  end
  
  -- Button 6: A
  if rf.btnp(6) then
    rf.sfx("stopall") -- Stop all audio before transition
    if game then
      game.changeState("menu")
    end
    return
  end
  
  -- Button 7: B
  if rf.btnp(7) then
    rf.sfx("stopall") -- Stop all audio before transition
    if game then
      game.changeState("menu")
    end
    return
  end
  
  -- Button 8: X
  if rf.btnp(8) then
    rf.sfx("stopall") -- Stop all audio before transition
    if game then
      game.changeState("menu")
    end
    return
  end
  
  -- Button 9: Y
  if rf.btnp(9) then
    rf.sfx("stopall") -- Stop all audio before transition
    if game then
      game.changeState("menu")
    end
    return
  end
  
  -- Button 10: TURBO
  if rf.btnp(10) then
    rf.sfx("stopall") -- Stop all audio before transition
    if game then
      game.changeState("menu")
    end
    return
  end
end

function _UPDATE(dt)
  splash_time = splash_time + dt
  road_offset = road_offset + dt * 100
  if road_offset > 480 then
    road_offset = road_offset - 480
  end
  
  -- Animate car moving across screen
  car_x = car_x + dt * 80
  if car_x > 580 then
    car_x = -100
  end
  
  -- Start music on first frame
  if not music_started then
    rf.music("menu_music")
    music_started = true
  end
  
  -- Auto-advance after duration
  if splash_time >= splash_duration then
    rf.sfx("stopall")
    if game then
      game.changeState("menu")
    end
  end
end

function _DRAW()
  -- Clear screen
  rf.clear_i(COLOR_BLACK)
  
  -- Draw scrolling road pattern (perspective lines)
  for i = 0, 20 do
    local y = (i * 25 + road_offset) % 540 - 20
    if y >= 0 and y < 270 then
      local width = 200 + (y / 270) * 280  -- Wider at bottom
      local x1 = 240 - width / 2
      local x2 = 240 + width / 2
      rf.line(x1, y, x2, y, COLOR_GRAY)
    end
  end
  
  -- Draw road center lines
  for i = 0, 10 do
    local y = (i * 50 + road_offset) % 540 - 20
    if y >= 0 and y < 270 then
      rf.line(240 - 10, y, 240 - 10, y + 20, COLOR_WHITE)
      rf.line(240 + 10, y, 240 + 10, y + 20, COLOR_WHITE)
    end
  end
  
  -- Draw title "POLE POSITION"
  local title_y = 80
  local title_scale = 1.0 + math.sin(splash_time * 2) * 0.1
  local title_size = math.floor(3 * title_scale)
  
  -- Title (using print_anchored for centering)
  rf.print_anchored("POLE POSITION", "topcenter", COLOR_WHITE)
  
  -- Subtitle (centered)
  local subtitle_text = "RACE TO THE FINISH"
  local subtitle_x = 240 - string.len(subtitle_text) * 3
  rf.print_xy(subtitle_x, title_y + 50, subtitle_text, COLOR_LIGHT_GRAY)
  
  -- Draw animated car
  if car_x > -50 and car_x < 530 then
    draw_car_simple(car_x, 180, COLOR_RED, 1.5)
  end
  
  -- Press any button text
  if math.floor(splash_time * 2) % 2 == 0 then
    local msg = "Press any key to continue"
    local msg_x = 240 - string.len(msg) * 3
    rf.print_xy(msg_x, 250, msg, COLOR_CYAN)
  end
end

function _EXIT()
  rf.sfx("stopall")
end

function _DONE()
  -- Shutdown (unused)
end

-- Simple car drawing function for splash
function draw_car_simple(x, y, color, scale)
  local w = 16 * scale
  local h = 8 * scale
  
  -- Car body
  rf.rectfill(x - w/2, y - h/2, x + w/2, y + h/2, color)
  -- Windshield
  rf.rectfill(x - w/4, y - h/3, x + w/4, y - h/6, COLOR_CYAN)
  -- Wheels
  rf.circfill(x - w/3, y + h/2, 3 * scale, COLOR_BLACK)
  rf.circfill(x + w/3, y + h/2, 3 * scale, COLOR_BLACK)
  rf.circfill(x - w/3, y - h/2, 3 * scale, COLOR_BLACK)
  rf.circfill(x + w/3, y - h/2, 3 * scale, COLOR_BLACK)
end

