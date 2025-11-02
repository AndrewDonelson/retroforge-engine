-- Game splash screen - plays for 5 seconds before menu

-- Color constants (must match main.lua)
local COLOR_BLACK = 0
local COLOR_WHITE = 1
local COLOR_GRAY = 8
local COLOR_LIGHT_GRAY = 32
local COLOR_DARK_GRAY = 5
local COLOR_LANDINGPAD = 26
local COLOR_STARS = 45

local splash_time = 0
local splash_duration = 5.0
local stars_splash = {}
local TWO_PI = 6.283185307179586

-- Generate stars for splash screen
local function generate_splash_stars()
  stars_splash = {}
  local star_seed = 12345
  local function star_rnd() star_seed = (1103515245*star_seed + 12345) % 2147483648; return star_seed end
  local function star_frand(a,b) return a + (star_rnd() % 10000)/10000*(b-a) end
  for i=1,80 do
    stars_splash[i] = {
      x = star_frand(0, 480),
      y = star_frand(10, 200),
      phase = star_frand(0, TWO_PI),
      speed = star_frand(0.3, 1.5),
      size = math.floor(star_frand(1, 3))
    }
  end
end

-- Draw detailed lander (more detailed than in-game)
local function draw_detailed_lander(x, y, angle)
  local cos_a = math.cos(angle)
  local sin_a = math.sin(angle)
  
  -- Main body (triangle)
  local body_size = 10
  local p0x, p0y = 0, -body_size
  local p1x, p1y = -body_size*0.7, body_size
  local p2x, p2y = body_size*0.7, body_size
  
  local function rotate(px, py)
    return math.floor(px*cos_a - py*sin_a + 0.5), math.floor(px*sin_a + py*cos_a + 0.5)
  end
  
  local a0x, a0y = rotate(p0x, p0y)
  local a1x, a1y = rotate(p1x, p1y)
  local a2x, a2y = rotate(p2x, p2y)
  
  -- Draw body outline (filled)
  rf.tri(x + a0x, y + a0y, x + a1x, y + a1y, x + a2x, y + a2y, COLOR_GRAY)
  rf.line(x + a0x, y + a0y, x + a1x, y + a1y, COLOR_WHITE)
  rf.line(x + a1x, y + a1y, x + a2x, y + a2y, COLOR_WHITE)
  rf.line(x + a2x, y + a2y, x + a0x, y + a0y, COLOR_WHITE)
  
  -- Draw cockpit window
  local win_x, win_y = rotate(0, -body_size*0.5)
  rf.circfill(x + win_x, y + win_y, 4, COLOR_LIGHT_GRAY)
  rf.circ(x + win_x, y + win_y, 4, COLOR_WHITE)
  -- Window highlight
  rf.circfill(x + win_x - 1, y + win_y - 1, 2, COLOR_WHITE)
  
  -- Draw landing legs (more detailed)
  local leg1x, leg1y = rotate(-body_size*0.9, body_size*0.9)
  local leg2x, leg2y = rotate(body_size*0.9, body_size*0.9)
  rf.line(x + a1x, y + a1y, x + leg1x, y + leg1y, COLOR_WHITE)
  rf.line(x + a2x, y + a2y, x + leg2x, y + leg2y, COLOR_WHITE)
  -- Foot pads
  rf.circfill(x + leg1x, y + leg1y, 2, COLOR_GRAY)
  rf.circfill(x + leg2x, y + leg2y, 2, COLOR_GRAY)
  
  -- Draw engine bell (detailed)
  local eng_x, eng_y = rotate(0, body_size*1.2)
  rf.circfill(x + eng_x, y + eng_y, 3, COLOR_DARK_GRAY)
  rf.circ(x + eng_x, y + eng_y, 3, COLOR_GRAY)
  -- Engine nozzle
  rf.circfill(x + eng_x, y + eng_y, 2, COLOR_LIGHT_GRAY)
  
  -- Draw side panels/details
  local panel1x, panel1y = rotate(-body_size*0.3, body_size*0.3)
  local panel2x, panel2y = rotate(body_size*0.3, body_size*0.3)
  rf.rectfill(x + panel1x - 1, y + panel1y - 1, x + panel1x + 1, y + panel1y + 1, COLOR_GRAY)
  rf.rectfill(x + panel2x - 1, y + panel2y - 1, x + panel2x + 1, y + panel2y + 1, COLOR_GRAY)
end

-- Draw Earth on horizon
local function draw_earth(progress)
  local earth_y = 200 - math.floor(progress * 20) -- Earth rises as splash progresses
  local earth_size = 60
  local earth_x = 240
  
  -- Earth disk with gradient effect (using gray scale palette)
  -- Outer edge (darker)
  rf.circ(earth_x, earth_y, earth_size, COLOR_GRAY)
  rf.circfill(earth_x, earth_y, earth_size, COLOR_GRAY)
  
  -- Inner highlights (lighter gray)
  rf.circfill(earth_x, earth_y, earth_size - 10, COLOR_LIGHT_GRAY)
  rf.circfill(earth_x, earth_y, earth_size - 20, COLOR_WHITE)
  
  -- Earth continents (darker patches)
  rf.circfill(earth_x - 10, earth_y - 5, 8, COLOR_DARK_GRAY)
  rf.circfill(earth_x + 12, earth_y + 8, 6, COLOR_DARK_GRAY)
  rf.circfill(earth_x - 8, earth_y + 12, 5, COLOR_DARK_GRAY)
  rf.circfill(earth_x + 5, earth_y - 10, 7, COLOR_DARK_GRAY)
  
  -- Earth atmosphere glow (subtle)
  rf.circ(earth_x, earth_y, earth_size + 2, COLOR_LIGHT_GRAY)
end

-- Draw terrain with landing pad
local function draw_splash_terrain(progress)
  local horizon_y = 220 + math.floor(progress * 15) -- Terrain rises as splash progresses
  
  -- Draw terrain silhouette
  for x = 0, 479, 2 do
    local height_variation = math.sin(x * 0.02) * 8 + math.sin(x * 0.05) * 4
    local ty = horizon_y + math.floor(height_variation)
    if ty < 269 then
      rf.rectfill(x, ty, x+1, 269, COLOR_GRAY)
    end
  end
  
  -- Draw landing pad (centered)
  local pad_w = 100
  local pad_x0 = 240 - pad_w/2
  local pad_x1 = 240 + pad_w/2
  local pad_y = horizon_y - 5
  rf.rectfill(pad_x0, pad_y, pad_x1, pad_y+3, COLOR_LANDINGPAD)
  rf.rect(pad_x0-1, pad_y-1, pad_x1+1, pad_y+4, COLOR_WHITE)
end

-- Draw animated stars
local function draw_splash_stars(time)
  for i = 1, #stars_splash do
    local s = stars_splash[i]
    local brightness = math.floor(128 + 127 * math.sin(time * s.speed + s.phase))
    brightness = math.max(128, math.min(255, brightness))
    
    -- Draw star based on size
    if s.size >= 2 then
      rf.pset(s.x, s.y, COLOR_STARS)
      rf.pset(s.x+1, s.y, COLOR_STARS)
      rf.pset(s.x, s.y+1, COLOR_STARS)
    else
      rf.pset(s.x, s.y, COLOR_STARS)
    end
  end
end

-- State module functions (required by module system)

function _INIT()
  generate_splash_stars()
end

function _ENTER()
  splash_time = 0
  -- Start music (4 seconds)
  rf.music("space_odyssey")
end

function _HANDLE_INPUT()
  -- Any input skips splash
  if rf.btnp(0) or rf.btnp(1) or rf.btnp(2) or rf.btnp(3) or rf.btnp(4) or rf.btnp(5) then
    game.changeState("menu")
  end
end

function _UPDATE(dt)
  splash_time = splash_time + dt
  
  -- Auto-transition after duration
  if splash_time >= splash_duration then
    game.changeState("menu")
  end
end

function _DRAW()
  -- Clear screen
  rf.clear_i(COLOR_BLACK)
  
  local progress = math.min(1.0, splash_time / splash_duration)
  
  -- Draw stars
  draw_splash_stars(splash_time)
  
  -- Draw Earth on horizon
  draw_earth(progress)
  
  -- Draw terrain
  draw_splash_terrain(progress)
  
  -- Draw title at top (always visible)
  rf.print_anchored("MOON LANDER", "topcenter", COLOR_WHITE)
  
  -- Draw detailed lander (animated: floating and slightly rotating)
  local lander_x = 240
  local lander_y = 110 + math.sin(splash_time * 2) * 5 -- Floating animation
  local lander_angle = math.sin(splash_time * 0.5) * 0.1 -- Subtle rotation
  draw_detailed_lander(lander_x, lander_y, lander_angle)
  
  -- Show subtitle text (appears after 1 second, positioned lower)
  local subtitle_alpha = math.min(1.0, (splash_time - 1.0) / 0.5)
  if subtitle_alpha > 0 then
    rf.print_xy(240 - string.len("Land safely on the moon")*3, 35, "Land safely on the moon", COLOR_LIGHT_GRAY)
  end
end

function _EXIT()
  -- Stop music
  rf.music_stop()
end

function _DONE()
  -- Shutdown (unused)
end

