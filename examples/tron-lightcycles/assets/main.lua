-- Tron Light Cycles game
-- Using GameStateMachine with module-based states

-- Shared game state (accessible by all state modules)
level = 1
best_level = 1
score = 0
countdown = 3.0
move_timer = 0.0

-- Player and enemies (shared)
player = nil
enemies = {}
num_enemies = 0

-- Grid settings (shared constants)
GRID_WIDTH = 480
GRID_HEIGHT = 270
CELL_SIZE = 3 -- pixels per grid cell
GRID_W = math.floor(GRID_WIDTH / CELL_SIZE) -- 160 cells
GRID_H = math.floor(GRID_HEIGHT / CELL_SIZE) -- 90 cells

-- Directions: 0=up, 1=right, 2=down, 3=left
DIR_UP = 0
DIR_RIGHT = 1
DIR_DOWN = 2
DIR_LEFT = 3

-- Color indices (using Super Mario 50 palette: 0-49)
COLOR_BLACK = 0
COLOR_WHITE = 1
COLOR_CYAN = 31 -- Light Cyan / Aqua Blue
COLOR_GRAY = 25 -- Gray for dimmed text
COLOR_RED = 4 -- Red for game over
COLOR_BLUE = 48 -- Blue for menu highlights

-- Light cycle trail color definitions for Super Mario 50 palette
-- Yellow (player): 8=highlight, 9=base, 10=shadow
-- Red (enemies): 2=highlight, 3=base, 4=shadow
PLAYER_TRAIL_COLORS = {highlight = 8, base = 9, shadow = 10}
ENEMY_TRAIL_COLORS = {highlight = 2, base = 3, shadow = 4}

-- Game grid (true = occupied, false = empty) - shared
grid = {}
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

-- (Player and enemies already declared in shared section above)

-- Game parameters (scale with level) - shared
base_speed = 5.0 -- moves per second
speed = base_speed
base_trail_length = 20
trail_length = base_trail_length

-- Level seed for deterministic placement - shared
level_seed = 0

-- Random number generator using level seed (shared functions)
function srnd()
  level_seed = (1103515245*level_seed + 12345) % 2147483648
  return level_seed
end

function frand(a, b)
  return a + (srnd() % 10000) / 10000 * (b - a)
end

function randi(a, b)
  return math.floor(frand(a, b + 1))
end

-- Get number of enemies for a level (shared function)
function get_enemy_count(lvl)
  return math.min(6, math.floor((lvl - 1) / 5) + 1)
end

-- Initialize level (shared function)
function init_level(lvl)
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
  player = create_cycle(player_x, player_y, DIR_UP, {}) -- Colors no longer used, kept for compatibility
  player.trail = {}
  grid[player_y][player_x] = true
  
  -- Place enemies
  enemies = {}
  for i=1,num_enemies do
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
          local enemy = create_cycle(ex, ey, edir, {}) -- Colors no longer used, kept for compatibility
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

-- Convert grid position to screen coordinates (shared function)
function grid_to_screen(gx, gy)
  return gx * CELL_SIZE + math.floor(CELL_SIZE / 2), 
         gy * CELL_SIZE + math.floor(CELL_SIZE / 2)
end

-- Convert screen coordinates to grid position (shared function)
function screen_to_grid(sx, sy)
  return math.floor(sx / CELL_SIZE), math.floor(sy / CELL_SIZE)
end

-- Check if position is valid and not occupied (shared function)
function can_move(gx, gy)
  if gx < 0 or gx >= GRID_W or gy < 0 or gy >= GRID_H then
    return false
  end
  return not grid[gy][gx]
end

-- Move a light cycle (shared function)
function move_cycle(cycle)
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

-- AI for enemy cycles (smarter: prefer straight, only turn when necessary) - shared function
function update_enemy_ai(enemy)
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

function _INIT()
  rf.palette_set("Super Mario 50")
  
  -- Import state modules
  rf.import("splash_state.lua")
  rf.import("menu_state.lua")
  rf.import("play_state.lua")
  rf.import("pause_state.lua")
  rf.import("gameover_state.lua")
  
  -- Initialize game variables
  level = 1
  score = 0
  best_level = 1
end

-- (Update functions are now in state modules)

-- Draw a 4-pixel-wide trail segment behind a cycle position (shared function)
-- dir: direction cycle is moving (0=up, 1=right, 2=down, 3=left)
-- gx, gy: grid position
-- colors: {highlight, base, shadow} color indices
function draw_trail_segment(gx, gy, dir, colors)
  local sx, sy = grid_to_screen(gx, gy)
  
  -- Trail is 4 pixels wide, perpendicular to movement direction
  -- Pattern: shadow (outer), base, base (center), highlight (outer)
  if dir == DIR_UP or dir == DIR_DOWN then
    -- Moving vertically, trail is horizontal (left-right)
    -- Draw 4 pixels centered: shadow, base, base, highlight
    rf.pset(sx - 1, sy, colors.shadow)   -- Left edge
    rf.pset(sx, sy, colors.base)          -- Left center
    rf.pset(sx + 1, sy, colors.base)      -- Right center
    rf.pset(sx + 2, sy, colors.highlight) -- Right edge
  else
    -- Moving horizontally, trail is vertical (up-down)
    -- Draw 4 pixels centered: shadow, base, base, highlight
    rf.pset(sx, sy - 1, colors.shadow)    -- Top edge
    rf.pset(sx, sy, colors.base)          -- Top center
    rf.pset(sx, sy + 1, colors.base)      -- Bottom center
    rf.pset(sx, sy + 2, colors.highlight) -- Bottom edge
  end
end

-- Draw trail line connecting two grid positions (shared function)
-- This fills gaps between segments
function draw_trail_line(gx1, gy1, gx2, gy2, dir, colors)
  local sx1, sy1 = grid_to_screen(gx1, gy1)
  local sx2, sy2 = grid_to_screen(gx2, gy2)
  
  -- Draw line between points with 4-pixel width
  local dx = sx2 - sx1
  local dy = sy2 - sy1
  local steps = math.max(math.abs(dx), math.abs(dy))
  
  if dir == DIR_UP or dir == DIR_DOWN then
    -- Vertical movement, horizontal trail
    for step = 0, steps do
      local t = steps > 0 and (step / steps) or 0
      local x = math.floor(sx1 + dx * t)
      local y = math.floor(sy1 + dy * t)
      rf.pset(x - 1, y, colors.shadow)
      rf.pset(x, y, colors.base)
      rf.pset(x + 1, y, colors.base)
      rf.pset(x + 2, y, colors.highlight)
    end
  else
    -- Horizontal movement, vertical trail
    for step = 0, steps do
      local t = steps > 0 and (step / steps) or 0
      local x = math.floor(sx1 + dx * t)
      local y = math.floor(sy1 + dy * t)
      rf.pset(x, y - 1, colors.shadow)
      rf.pset(x, y, colors.base)
      rf.pset(x, y + 1, colors.base)
      rf.pset(x, y + 2, colors.highlight)
    end
  end
end

-- Draw rotated sprite manually (since rf.spr doesn't support rotation) - shared function
function draw_sprite_rotated(name, x, y, angle)
  local sprite = rf.sprite(name)
  if not sprite then return end
  
  local w, h = sprite.width, sprite.height
  local cx, cy = x + w / 2, y + h / 2 -- Center point
  
  -- Convert angle (0=up, 90=right, 180=down, 270=left) to radians
  local rad = (angle * math.pi) / 180
  local cos_a = math.cos(rad)
  local sin_a = math.sin(rad)
  
  -- Draw sprite pixels with rotation
  for sy = 0, h - 1 do
    for sx = 0, w - 1 do
      local px = sx - w / 2 + 0.5
      local py = sy - h / 2 + 0.5
      
      -- Rotate around center
      local rx = px * cos_a - py * sin_a
      local ry = px * sin_a + py * cos_a
      
      local draw_x = math.floor(cx + rx)
      local draw_y = math.floor(cy + ry)
      
      local color_idx = sprite.pixels[sy + 1][sx + 1] -- Lua is 1-indexed
      if color_idx >= 0 then -- -1 is transparent
        rf.pset(draw_x, draw_y, color_idx)
      end
    end
  end
end

function draw_cycle(cycle, is_player)
  if not cycle.alive then return end
  
  -- Determine trail colors
  local trail_colors = is_player and PLAYER_TRAIL_COLORS or ENEMY_TRAIL_COLORS
  
  -- Draw continuous trail (no gaps)
  if #cycle.trail > 0 then
    -- Draw segments and lines between them to fill gaps
    for i=1,#cycle.trail do
      local pos = cycle.trail[i]
      -- Determine direction based on movement from this segment to next (or current for last)
      local trail_dir = cycle.dir
      if i < #cycle.trail then
        -- Determine direction from this segment to next
        local next_pos = cycle.trail[i+1]
        local dx = next_pos.x - pos.x
        local dy = next_pos.y - pos.y
        if dx > 0 then trail_dir = DIR_RIGHT
        elseif dx < 0 then trail_dir = DIR_LEFT
        elseif dy > 0 then trail_dir = DIR_DOWN
        else trail_dir = DIR_UP
        end
        
        -- Draw line between this segment and next to fill gap
        draw_trail_line(pos.x, pos.y, next_pos.x, next_pos.y, trail_dir, trail_colors)
      else
        -- Last segment - draw line to current cycle position
        draw_trail_line(pos.x, pos.y, cycle.x, cycle.y, cycle.dir, trail_colors)
      end
    end
  end
  
  -- Draw cycle sprite rotated to face direction
  local sx, sy = grid_to_screen(cycle.x, cycle.y)
  local sprite_name = is_player and "player" or "enemy"
  
  -- Convert direction to angle: 0=up, 90=right, 180=down, 270=left
  local angle = cycle.dir * 90
  
  -- Center 16x16 sprite on grid position (sprite origin is top-left, so offset by -8)
  draw_sprite_rotated(sprite_name, sx - 8, sy - 8, angle)
end

-- Draw HUD (shared function)
function draw_hud()
  -- Score and level at top left
  local score_text = "SCORE: " .. tostring(score)
  local level_text = "LEVEL: " .. tostring(level)
  rf.print_xy(2, 2, score_text, COLOR_WHITE)
  rf.print_xy(2, 10, level_text, COLOR_WHITE)
end

-- (Drawing functions are now in state modules)

