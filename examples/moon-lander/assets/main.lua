-- Moon Lander
-- A lunar landing game using GameStateMachine

-- Color constants (shared across all states)
COLOR_BLACK = 0          -- Background
COLOR_WHITE = 1          -- Ship, title, menu selected, messages
COLOR_GRAY = 8           -- Menu unselected, instructions, terrain, HUD ALT text
COLOR_LIGHT_GRAY = 32    -- Thrust flame, HUD VY text
COLOR_DARK_GRAY = 5      -- HUD FUEL text
COLOR_LANDINGPAD = 26    -- Landing pad
COLOR_STARS = 45         -- Stars

-- Shared game state (accessible by all state modules)
ground_y = 235
pad_x0, pad_x1 = 200, 280 
pad_y = ground_y - 2

ship = {
  x = 240, y = 60, vx = 0, vy = 0,
  angle = 0, -- radians
  fuel = 100,
  size = 6,
  alive = true,
  landed = false
}

-- Level-scaled parameters (initialized by set_level)
G = 1.0          -- gravity px/s^2
THRUST = 4.0     -- thrust accel px/s^2
ROT = 2.0        -- rad/s
land_speed = 18

countdown = 5.0
level = 1
score = 0
best_level = 1
best_score = 0

-- 50 levels, procedurally generated
levels = {}
level_seed = 0
stars = {}
heightmap = {}
TWO_PI = 6.283185307179586

for i = 1, 50 do
  local t = {}
  local rough = (i-1) / 49 -- 0..1
  t.padw = 80 - math.floor(rough * 60) -- 80..20
  t.G = 0.25 + 0.75 * rough
  t.THRUST = 2.0 + 6.0 * rough
  t.ROT = 2.0
  t.land_speed = 12 + math.floor(rough * 16)
  t.seed = i * 413
  levels[i] = t
end

-- Shared helper functions (accessible by all state modules)
function terrain_base_y(x)
  local L = levels[level]
  local rough = (level-1) / 49
  local base = math.sin((x + L.seed) * 0.06) * (8 + rough * 18)
  local micro = math.sin((x + L.seed * 3) * 0.23) * 3
  return ground_y - 10 - math.floor(base + micro)
end

function terrain_y(x)
  if x >= pad_x0 and x <= pad_x1 then return pad_y end
  return heightmap[x] or terrain_base_y(x)
end

function set_level(idx)
  level = math.max(1, math.min(50, idx))
  local L = levels[level]
  G = L.G
  THRUST = L.THRUST
  ROT = L.ROT
  land_speed = L.land_speed
  level_seed = L.seed
  
  local function srnd()
    level_seed = (1103515245 * level_seed + 12345) % 2147483648
    return level_seed
  end
  local function frand(a, b)
    return a + (srnd() % 10000) / 10000 * (b - a)
  end
  
  local padw = L.padw
  local bestx, bestSlope = 240, 1e9
  for _ = 1, 80 do
    local halfw = math.floor(padw / 2)
    local cx = math.floor(frand(20 + halfw, 460 - halfw))
    local x0 = cx - halfw
    local x1 = cx + halfw
    local y0 = terrain_base_y(x0)
    local y1 = terrain_base_y(x1)
    local slope = math.abs(y1 - y0)
    if slope < bestSlope then
      bestSlope = slope
      bestx = cx
    end
    if slope <= 2 then
      bestx = cx
      break
    end
  end
  
  local bestHalf = math.floor(padw / 2)
  pad_x0 = bestx - bestHalf
  pad_x1 = bestx + bestHalf
  local mid = math.floor((pad_x0 + pad_x1) / 2)
  pad_y = math.min(terrain_base_y(pad_x0), terrain_base_y(pad_x1), terrain_base_y(mid))
  
  for x = 0, 479 do
    heightmap[x] = terrain_base_y(x)
  end
  
  ship.x, ship.y, ship.vx, ship.vy, ship.angle = 240, 60, 0, 0, 0
  ship.fuel, ship.alive, ship.landed = 150, true, false
  countdown = 5.0
  
  -- Generate stars (random count 40-80 based on level seed)
  stars = {}
  local starSeed = L.seed * 17
  local function starRnd()
    starSeed = (1103515245 * starSeed + 12345) % 2147483648
    return starSeed
  end
  local function starFrand(a, b)
    return a + (starRnd() % 10000) / 10000 * (b - a)
  end
  local numStars = math.floor(starFrand(40, 81)) -- 40 to 80 inclusive
  for i = 1, numStars do
    stars[i] = {
      x = math.floor(starFrand(0, 480)),
      y = math.floor(starFrand(10, 200)),
      phase = starFrand(0, TWO_PI),
      speed = starFrand(0.5, 2.0)
    }
  end
end

function clamp(v, a, b)
  if v < a then return a elseif v > b then return b else return v end
end

-- Crash sequence variables (shared)
crash_phase = nil
crash_timer = 0
taps_started = false
taps_total_dur = 0

function start_crash_sequence()
  crash_phase = "crash"
  crash_timer = 0
  taps_started = false
  taps_total_dur = 0
end

function _INIT()
  -- Palette is set from manifest.json (Grayscale 50), don't override it
  
  -- Import state modules
  rf.import("splash_state.lua")
  rf.import("menu_state.lua")
  rf.import("play_state.lua")
  rf.import("pause_state.lua")
  
  -- Initialize first level
  set_level(1)
end
