-- Pole Position
-- A classic arcade racing game with procedural tracks

-- Color constants (using Super Mario 50 palette: 0-49)
COLOR_BLACK = 0
COLOR_WHITE = 1
COLOR_DARKEST_GRAY = 24  -- Darkest gray for wheels
COLOR_DARK_GRAY = 25     -- Gray for dimmed text
COLOR_GRAY = 25          -- Gray for roads/backgrounds
COLOR_MEDIUM_GRAY = 25   -- Gray for unselected menu items
COLOR_LIGHT_GRAY = 32    -- Light gray for subtitles
COLOR_RED = 4            -- Red for cars, highlights
COLOR_GREEN = 3          -- Green for success/position
COLOR_BLUE = 48          -- Blue for menu highlights
COLOR_YELLOW = 6         -- Yellow for checkpoints, speed
COLOR_CYAN = 31          -- Light Cyan / Aqua Blue for highlights

-- Shared game state
selected_track = 1
lap_time = 0.0
best_lap_time = 999.9
current_position = 0
total_cars = 8
race_started = false
race_finished = false

-- Track data (10 procedural tracks)
tracks = {}
for i = 1, 10 do
  tracks[i] = {
    name = "Track " .. i,
    difficulty = i,
    seed = i * 1000 + 12345
  }
end

-- Track names
tracks[1].name = "Desert Highway"
tracks[2].name = "Mountain Pass"
tracks[3].name = "Coastal Road"
tracks[4].name = "Forest Trail"
tracks[5].name = "City Streets"
tracks[6].name = "Alpine Route"
tracks[7].name = "Canyon Run"
tracks[8].name = "Countryside"
tracks[9].name = "Twilight Drive"
tracks[10].name = "Championship"

-- Track generation function (procedural) - Top-down view
function generate_track(track_index)
  local track = tracks[track_index]
  local segments = {}
  local seed = track.seed
  local difficulty = track.difficulty
  
  -- Simple RNG based on seed
  local function rnd()
    seed = (1103515245 * seed + 12345) % 2147483648
    return seed
  end
  
  local function frand(a, b)
    return a + (rnd() % 10000) / 10000 * (b - a)
  end
  
  -- Generate track segments (for top-down, we just need lane count info)
  -- Not used in top-down view, but kept for compatibility
  track.segments = {}
  track.num_segments = 100
  return track
end

-- Initialize all tracks
for i = 1, 10 do
  generate_track(i)
end

function _INIT()
  rf.palette_set("Super Mario 50")
  
  -- Import state modules
  rf.import("splash_state.lua")
  rf.import("menu_state.lua")
  rf.import("track_select_state.lua")
  rf.import("play_state.lua")
  rf.import("pause_state.lua")
  rf.import("results_state.lua")
end

