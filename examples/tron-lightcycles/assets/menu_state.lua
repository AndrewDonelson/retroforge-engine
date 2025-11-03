-- Menu State Module for Tron Light Cycles

local menu_idx = 1
local menu_time = 0.0
local menu_music_started = false

function _INIT()
  -- Module initialization
end

function _ENTER()
  -- Reset menu when entering
  menu_idx = 1
  menu_time = 0.0
  menu_music_started = false
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
  if rf.btnp(0) then
    -- SELECT menu option
    if menu_idx == 1 then
      -- Start game - reset score and level
      rf.sfx("select")
      rf.sfx("stopall") -- Stop all audio including music
      level = 1
      score = 0
      init_level(level)
      countdown = 3.0
      game.changeState("play")
    elseif menu_idx == 2 then
      -- Exit
      rf.sfx("select")
      game.exit()
    end
  end
  
  -- Button 1: START
  if rf.btnp(1) then
    -- START bypasses menu and goes straight to PLAY
    rf.sfx("select")
    rf.sfx("stopall") -- Stop all audio including music
    level = 1
    score = 0
    init_level(level)
    countdown = 3.0
    game.changeState("play")
    return
  end
  
  -- Button 2: UP
  if rf.btnp(2) then
    menu_idx = math.max(1, menu_idx - 1)
    rf.sfx("select")
  end
  
  -- Button 3: DOWN
  if rf.btnp(3) then
    menu_idx = math.min(2, menu_idx + 1)
    rf.sfx("select")
  end
  
  -- Button 4: LEFT
  -- Not used in menu
  
  -- Button 5: RIGHT
  -- Not used in menu
  
  -- Button 6: A
  if rf.btnp(6) then
    -- A button also selects menu option (same as SELECT)
    if menu_idx == 1 then
      -- Start game - reset score and level
      rf.sfx("select")
      rf.sfx("stopall") -- Stop all audio including music
      level = 1
      score = 0
      init_level(level)
      countdown = 3.0
      game.changeState("play")
    elseif menu_idx == 2 then
      -- Exit
      rf.sfx("select")
      game.exit()
    end
  end
  
  -- Button 7: B
  -- Not used in menu
  
  -- Button 8: X
  -- Not used in menu
  
  -- Button 9: Y
  -- Not used in menu
  
  -- Button 10: TURBO
  -- Not used in menu
end

function _UPDATE(dt)
  menu_time = menu_time + dt
  
  -- Start menu music on first menu update
  if not menu_music_started then
    rf.music("menu_music")
    menu_music_started = true
  end
end

function _DRAW()
  -- Clear screen
  rf.clear_i(COLOR_BLACK)
  
  -- Animated background grid effect
  local grid_spacing = 8
  for y=0,GRID_HEIGHT-1,grid_spacing do
    for x=0,GRID_WIDTH-1,grid_spacing do
      local phase = (x + y + menu_time * 20) % 100
      if phase < 50 then
        local sx = x
        local sy = y
        rf.pset(sx, sy, COLOR_CYAN)
      end
    end
  end
  
  -- Title at top center
  rf.print_anchored("TRON", "topcenter", COLOR_BLUE)
  local cycles_y = 70
  local cycles_x = 240 - string.len("LIGHT CYCLES")*3
  rf.print_xy(cycles_x, cycles_y, "LIGHT CYCLES", COLOR_BLUE)
  
  -- Menu items (use blue for selected, gray for dimmed)
  local c1 = (menu_idx == 1) and COLOR_BLUE or COLOR_GRAY
  local c2 = (menu_idx == 2) and COLOR_BLUE or COLOR_GRAY
  
  local play_x = 240 - string.len("PLAY")*3
  local quit_x = 240 - string.len("QUIT")*3
  rf.print_xy(play_x, 110, "PLAY", c1)
  rf.print_xy(quit_x, 126, "QUIT", c2)
  
  -- Instructions
  local turn_x = 240 - string.len("Arrow keys: Turn")*3
  local select_x = 240 - string.len("SELECT/Enter: Select")*3
  local start_x = 240 - string.len("START: Play")*3
  rf.print_xy(turn_x, 160, "Arrow keys: Turn", COLOR_GRAY)
  rf.print_xy(select_x, 176, "SELECT/Enter: Select", COLOR_GRAY)
  rf.print_xy(start_x, 192, "START: Play", COLOR_GRAY)
  
  -- Best score
  local best_text = "Best Level: " .. tostring(best_level)
  local best_x = 240 - string.len(best_text)*3
  rf.print_xy(best_x, 210, best_text, COLOR_GRAY)
  
  -- Decorative light cycle trails
  local trail_time = menu_time * 0.5
  for i=1,3 do
    local trail_x = 120 + i * 80 + math.sin(trail_time + i) * 40
    local trail_y = 230 + math.cos(trail_time * 0.7 + i) * 20
    -- Use same colors as in game: base for head, highlight for first trail, shadow for rest
    local head_colors = {48, 3, 15} -- Player base, Enemy 1 base, Enemy 3 base
    local trail1_colors = {49, 4, 16} -- Player highlight, Enemy 1 highlight, Enemy 3 highlight
    local trail2_colors = {47, 2, 14} -- Player shadow, Enemy 1 shadow, Enemy 3 shadow
    rf.circfill(trail_x, trail_y, 3, head_colors[i]) -- Head (base)
    rf.circfill(trail_x - 5, trail_y - 2, 2, trail1_colors[i]) -- First trail segment (highlight)
    rf.circfill(trail_x - 10, trail_y - 4, 2, trail2_colors[i]) -- Rest of trail (shadow)
  end
end

function _EXIT()
  -- Stop music when leaving menu (if transitioning to play state)
  rf.sfx("stopall")
end

function _DONE()
  -- Shutdown (unused for this state)
end
