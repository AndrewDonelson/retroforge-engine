-- Tron Light Cycles game

local state = "menu" -- menu | playing | gameover | victory
local menu_idx = 1
local level = 1
local best_level = 1
local score = 0

-- Grid settings
local GRID_WIDTH = 480
local GRID_HEIGHT = 270
local CELL_SIZE = 3 -- pixels per grid cell
local GRID_W = math.floor(GRID_WIDTH / CELL_SIZE) -- 160 cells
local GRID_H = math.floor(GRID_HEIGHT / CELL_SIZE) -- 90 cells

-- Directions: 0=up, 1=right, 2=down, 3=left
local DIR_UP = 0
local DIR_RIGHT = 1
local DIR_DOWN = 2
local DIR_LEFT = 3

-- Color indices (using Super Mario 50 palette: 0-49)
local COLOR_BLACK = 0
local COLOR_WHITE = 1
local COLOR_CYAN = 31 -- Light Cyan / Aqua Blue
local COLOR_GRAY = 25 -- Gray for dimmed text
local COLOR_RED = 4 -- Red for game over
local COLOR_BLUE = 48 -- Blue for menu highlights

-- Light cycle color definitions: {head, first_trail, rest_trail}
-- Pattern: Each color has 3 shades - shadow (darkest), base (middle), highlight (brightest)
-- Head uses base (middle), first trail uses highlight (brightest), rest uses shadow (darkest)
-- 
-- Player: Blue shades
--   47 = shadow (darker), 48 = base (middle), 49 = highlight (brighter)
local PLAYER_COLORS = {head = 48, trail1 = 49, trail2 = 47}

-- Enemy colors (3 shades each: shadow, base, highlight)
-- Enemy 1: Orange shades - 2=shadow, 3=base, 4=highlight
-- Enemy 2: Yellow-Green shades - 8=shadow, 9=base, 10=highlight  
-- Enemy 3: Dark Blue shades - 14=shadow, 15=base, 16=highlight
-- Enemy 4-6: Will cycle through these 3 color sets
local ENEMY_COLOR_SETS = {
  {head = 3, trail1 = 4, trail2 = 2},    -- Enemy 1: Orange (base=3, highlight=4, shadow=2)
  {head = 9, trail1 = 10, trail2 = 8}, -- Enemy 2: Yellow-Green (base=9, highlight=10, shadow=8)
  {head = 15, trail1 = 16, trail2 = 14}, -- Enemy 3: Dark Blue (base=15, highlight=16, shadow=14)
}

-- Game grid (true = occupied, false = empty)
local grid = {}
for y=0,GRID_H-1 do
  grid[y] = {}
  for x=0,GRID_W-1 do
    grid[y][x] = false
  end
end

-- Light cycle structure
local function create_cycle(x, y, dir, colors)
  return {
    x = x,
    y = y,
    dir = dir,
    colors = colors, -- {head, trail1, trail2}
    trail = {}, -- Array of {x, y} positions
    alive = true
  }
end

-- Player and enemies
local player = nil
local enemies = {}
local num_enemies = 0

-- Game parameters (scale with level)
local base_speed = 5.0 -- moves per second
local speed = base_speed
local base_trail_length = 20
local trail_length = base_trail_length
local move_timer = 0.0

-- Level seed for deterministic placement
local level_seed = 0

-- Random number generator using level seed
local function srnd()
  level_seed = (1103515245*level_seed + 12345) % 2147483648
  return level_seed
end

local function frand(a, b)
  return a + (srnd() % 10000) / 10000 * (b - a)
end

local function randi(a, b)
  return math.floor(frand(a, b + 1))
end

-- Get number of enemies for a level
local function get_enemy_count(lvl)
  return math.min(6, math.floor((lvl - 1) / 5) + 1)
end

-- Initialize level
local function init_level(lvl)
  level = lvl
  level_seed = lvl * 7919 -- prime multiplier for variation
  
  -- Calculate difficulty scaling
  local difficulty = (lvl - 1) / 50 -- 0 to ~1 over 50 levels
  speed = base_speed * (1.0 + difficulty * 2.0) -- 1x to 3x speed
  trail_length = base_trail_length + math.floor(difficulty * 40) -- 20 to 60 trail length
  
  -- Clear grid and trails from all cycles
  for y=0,GRID_H-1 do
    for x=0,GRID_W-1 do
      grid[y][x] = false
    end
  end
  
  -- Clear all trail arrays
  if player then player.trail = {} end
  for i=1,#enemies do
    if enemies[i] then enemies[i].trail = {} end
  end
  
  -- Calculate number of enemies
  num_enemies = get_enemy_count(level)
  
  -- Place player at random bottom position, facing up
  local player_x = randi(5, GRID_W - 6)
  local player_y = GRID_H - 3
  player = create_cycle(player_x, player_y, DIR_UP, PLAYER_COLORS)
  player.trail = {}
  grid[player_y][player_x] = true
  
  -- Place enemies
  enemies = {}
  for i=1,num_enemies do
    local enemy_colors = ENEMY_COLOR_SETS[((i-1) % #ENEMY_COLOR_SETS) + 1]
    local placed = false
    local attempts = 0
    while not placed and attempts < 100 do
      local ex = randi(5, GRID_W - 6)
      local ey = randi(3, 15) -- Top area
      local edir = DIR_DOWN -- Start moving down
      
      -- Check if position is free
      if not grid[ey][ex] then
        local too_close = false
        -- Don't place too close to player
        local dx = math.abs(ex - player_x)
        local dy = math.abs(ey - player_y)
        if dx < 20 and dy < 40 then
          too_close = true
        end
        -- Don't place too close to other enemies
        for j=1,#enemies do
          local e = enemies[j]
          local dx2 = math.abs(ex - e.x)
          local dy2 = math.abs(ey - e.y)
          if dx2 < 15 and dy2 < 30 then
            too_close = true
            break
          end
        end
        
        if not too_close then
          local enemy = create_cycle(ex, ey, edir, enemy_colors)
          enemy.trail = {}
          grid[ey][ex] = true
          table.insert(enemies, enemy)
          placed = true
        end
      end
      attempts = attempts + 1
    end
  end
  
  move_timer = 0.0
end

-- Convert grid position to screen coordinates
local function grid_to_screen(gx, gy)
  return gx * CELL_SIZE + math.floor(CELL_SIZE / 2), 
         gy * CELL_SIZE + math.floor(CELL_SIZE / 2)
end

-- Convert screen coordinates to grid position
local function screen_to_grid(sx, sy)
  return math.floor(sx / CELL_SIZE), math.floor(sy / CELL_SIZE)
end

-- Check if position is valid and not occupied
local function can_move(gx, gy)
  if gx < 0 or gx >= GRID_W or gy < 0 or gy >= GRID_H then
    return false
  end
  return not grid[gy][gx]
end

-- Move a light cycle
local function move_cycle(cycle)
  if not cycle.alive then return false end
  
  local new_x, new_y = cycle.x, cycle.y
  
  -- Calculate new position based on direction
  if cycle.dir == DIR_UP then
    new_y = cycle.y - 1
  elseif cycle.dir == DIR_RIGHT then
    new_x = cycle.x + 1
  elseif cycle.dir == DIR_DOWN then
    new_y = cycle.y + 1
  elseif cycle.dir == DIR_LEFT then
    new_x = cycle.x - 1
  end
  
  -- Check collision with walls or any trail (including own trail)
  if not can_move(new_x, new_y) then
    cycle.alive = false
    -- Clear current position from grid since cycle is dead
    grid[cycle.y][cycle.x] = false
    return false
  end
  
  -- Old position becomes part of the trail - it stays marked in grid (already true)
  table.insert(cycle.trail, {x = cycle.x, y = cycle.y})
  
  -- Limit trail length (remove oldest trail segments from array and grid)
  while #cycle.trail > trail_length do
    local old = table.remove(cycle.trail, 1)
    -- Clear from grid (safe - if another cycle needed it, collision would have happened)
    grid[old.y][old.x] = false
  end
  
  -- Update position
  cycle.x, cycle.y = new_x, new_y
  
  -- Safety check: new position should be empty (already checked by can_move)
  if grid[cycle.y][cycle.x] then
    cycle.alive = false
    grid[cycle.y][cycle.x] = false
    return false
  end
  
  -- Mark new position as occupied (cycle head)
  grid[cycle.y][cycle.x] = true
  
  return true
end

-- AI for enemy cycles (smarter: prefer straight, only turn when necessary)
local function update_enemy_ai(enemy)
  if not enemy.alive then return end
  
  -- First, check if current direction is still safe
  local can_continue = false
  local test_x, test_y = enemy.x, enemy.y
  if enemy.dir == DIR_UP then test_y = enemy.y - 1
  elseif enemy.dir == DIR_RIGHT then test_x = enemy.x + 1
  elseif enemy.dir == DIR_DOWN then test_y = enemy.y + 1
  elseif enemy.dir == DIR_LEFT then test_x = enemy.x - 1
  end
  
  if can_move(test_x, test_y) then
    can_continue = true
  end
  
  -- If we can continue straight, only turn occasionally (10% chance)
  if can_continue and frand(0, 1) > 0.1 then
    return -- Keep going straight
  end
  
  -- Otherwise, need to turn - find best options
  local dx = player.x - enemy.x
  local dy = player.y - enemy.y
  local options = {}
  
  -- Check each direction
  for dir=0,3 do
    if dir == (enemy.dir + 2) % 4 then
      -- Can't reverse
    else
      test_x, test_y = enemy.x, enemy.y
      if dir == DIR_UP then test_y = enemy.y - 1
      elseif dir == DIR_RIGHT then test_x = enemy.x + 1
      elseif dir == DIR_DOWN then test_y = enemy.y + 1
      elseif dir == DIR_LEFT then test_x = enemy.x - 1
      end
      
      if can_move(test_x, test_y) then
        local dist = math.sqrt((test_x - player.x)^2 + (test_y - player.y)^2)
        -- Prefer current direction slightly to reduce constant turning
        local bonus = (dir == enemy.dir) and -5 or 0
        table.insert(options, {dir = dir, dist = dist, score = dist + bonus})
      end
    end
  end
  
  if #options > 0 then
    -- Sort by score (prefer closer to player, with slight bonus for continuing straight)
    table.sort(options, function(a, b)
      local bias_a = frand(0, 30) -- Less randomness for smarter AI
      local bias_b = frand(0, 30)
      return (a.score + bias_a) < (b.score + bias_b)
    end)
    
    -- Take the best option (or sometimes second best for slight unpredictability)
    local choice_idx = 1
    if frand(0, 1) < 0.15 and #options > 1 then -- 15% chance to take second best
      choice_idx = 2
    end
    
    enemy.dir = options[choice_idx].dir
  end
end

local countdown = 3.0
local gameover_timer = 0.0

function _INIT()
  rf.palette_set("default")
  state = "menu"
  level = 1
  score = 0
  menu_time = 0.0
  
  -- Play futuristic menu music
  rf.music("menu_music")
end

local function update_menu(dt)
  if rf.btnp(2) then -- Up
    menu_idx = math.max(1, menu_idx - 1)
    rf.sfx("move")
  end
  if rf.btnp(3) then -- Down
    menu_idx = math.min(2, menu_idx + 1)
    rf.sfx("move")
  end
  if rf.btnp(4) or rf.btnp(5) then -- Select
    if menu_idx == 1 then
      rf.sfx("select")
      level = 1
      score = 0
      init_level(level)
      state = "playing"
      countdown = 3.0
      rf.music("start_melody")
    else
      rf.sfx("select")
      rf.quit()
    end
  end
end

local function update_countdown(dt)
  local prev_int = math.ceil(countdown)
  countdown = math.max(0, countdown - dt)
  local now_int = math.ceil(countdown)
  if now_int < prev_int and now_int > 0 then
    rf.tone(600, 0.1, 0.3)
  end
end

local function update_gameplay(dt)
  -- Handle player input
  if player.alive then
    if rf.btnp(0) then -- Left
      if player.dir ~= DIR_RIGHT then
        player.dir = DIR_LEFT
        rf.sfx("move")
      end
    elseif rf.btnp(1) then -- Right
      if player.dir ~= DIR_LEFT then
        player.dir = DIR_RIGHT
        rf.sfx("move")
      end
    elseif rf.btnp(2) then -- Up
      if player.dir ~= DIR_DOWN then
        player.dir = DIR_UP
        rf.sfx("move")
      end
    elseif rf.btnp(3) then -- Down
      if player.dir ~= DIR_UP then
        player.dir = DIR_DOWN
        rf.sfx("move")
      end
    end
  end
  
  -- Update move timer
  move_timer = move_timer + dt * speed
  
  if move_timer >= 1.0 then
    move_timer = move_timer - 1.0
    
    -- Update AI for enemies
    for i=1,#enemies do
      if enemies[i].alive then
        update_enemy_ai(enemies[i])
      end
    end
    
    -- Move player
    if player.alive then
      if not move_cycle(player) then
        -- Player crashed
        rf.sfx("crash")
        state = "gameover"
        gameover_timer = 0.0
        best_level = math.max(best_level, level)
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
    
    if player.alive and all_enemies_dead then
      rf.sfx("land") -- Victory sound
      score = score + level * 100
      level = level + 1
      init_level(level)
      countdown = 2.0
    end
  end
end

local function update_gameover(dt)
  gameover_timer = gameover_timer + dt
  if gameover_timer >= 3.0 then
    if rf.btnp(4) or rf.btnp(5) then
      rf.sfx("select")
      state = "menu"
      menu_idx = 1
    end
  end
end

function _UPDATE(dt)
  if state == "menu" then
    menu_time = menu_time + dt
    update_menu(dt)
  elseif state == "playing" then
    if countdown > 0 then
      update_countdown(dt)
    else
      update_gameplay(dt)
    end
  elseif state == "gameover" then
    update_gameover(dt)
  end
end

local function draw_cycle(cycle)
  if not cycle.alive then return end
  
  -- Draw trail with gradient
  -- First trail segment uses trail1 color, rest use trail2 color
  for i=1,#cycle.trail do
    local pos = cycle.trail[i]
    local sx, sy = grid_to_screen(pos.x, pos.y)
    local trail_color = cycle.colors.trail2 -- Default: rest of trail
    if i == #cycle.trail then
      -- First trail segment (closest to head) uses trail1
      trail_color = cycle.colors.trail1
    end
    rf.circfill(sx, sy, 2, trail_color)
  end
  
  -- Draw cycle head
  local sx, sy = grid_to_screen(cycle.x, cycle.y)
  rf.circfill(sx, sy, 3, cycle.colors.head)
end

local menu_time = 0.0

local function draw_menu()
  -- Animated background grid effect
  local grid_spacing = 8
  for y=0,GRID_HEIGHT-1,grid_spacing do
    for x=0,GRID_WIDTH-1,grid_spacing do
      local phase = (x + y + menu_time * 20) % 100
      if phase < 50 then
        local sx = x
        local sy = y
        local alpha = math.floor(phase / 50 * 60 + 20)
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
  local select_x = 240 - string.len("O/X/Enter: Select")*3
  rf.print_xy(turn_x, 160, "Arrow keys: Turn", COLOR_GRAY)
  rf.print_xy(select_x, 176, "O/X/Enter: Select", COLOR_GRAY)
  
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

local function draw_hud()
  rf.print_xy(2, 2, "LEVEL: " .. tostring(level), COLOR_WHITE)
  rf.print_xy(2, 12, "SCORE: " .. tostring(score), COLOR_WHITE)
  rf.print_xy(2, 22, "ENEMIES: " .. tostring(num_enemies), COLOR_WHITE)
end

local function draw_playing()
  -- Draw grid background (optional subtle grid)
  for y=0,GRID_H-1,5 do
    for x=0,GRID_W-1,5 do
      local sx, sy = grid_to_screen(x, y)
      rf.pset(sx, sy, COLOR_WHITE)
    end
  end
  
  -- Draw enemy cycles
  for i=1,#enemies do
    draw_cycle(enemies[i])
  end
  
  -- Draw player cycle
  if player then
    draw_cycle(player)
  end
  
  -- Draw HUD
  draw_hud()
  
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

local function draw_gameover()
  -- Redraw game state
  for y=0,GRID_H-1,5 do
    for x=0,GRID_W-1,5 do
      local sx, sy = grid_to_screen(x, y)
      rf.pset(sx, sy, COLOR_WHITE)
    end
  end
  
  for i=1,#enemies do
    draw_cycle(enemies[i])
  end
  
  if player then
    draw_cycle(player)
  end
  
  draw_hud()
  
  local gameover_x = 240 - string.len("GAME OVER")*3
  local level_text = "Level: " .. tostring(level)
  local level_x = 240 - string.len(level_text)*3
  local score_text = "Score: " .. tostring(score)
  local score_x = 240 - string.len(score_text)*3
  rf.print_xy(gameover_x, 110, "GAME OVER", COLOR_RED)
  rf.print_xy(level_x, 130, level_text, COLOR_WHITE)
  rf.print_xy(score_x, 150, score_text, COLOR_WHITE)
  if gameover_timer >= 3.0 then
    local continue_text = "Press O/X/Enter to continue"
    local continue_x = 240 - string.len(continue_text)*3
    rf.print_xy(continue_x, 170, continue_text, COLOR_GRAY)
  end
end

function _DRAW()
  rf.clear_i(COLOR_BLACK)
  
  if state == "menu" then
    draw_menu()
  elseif state == "playing" then
    draw_playing()
  elseif state == "gameover" then
    draw_gameover()
  end
end

