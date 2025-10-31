-- Moon Lander demo (controllable)

local ground_y = 235
local pad_x0, pad_x1 = 200, 280
local pad_y = ground_y-2

local ship = {
  x=240, y=60, vx=0, vy=0,
  angle=0, -- radians
  fuel=150,
  size=6,
  alive=true, landed=false
}

-- Level-scaled parameters (initialized by set_level)
local G = 1.0          -- gravity px/s^2
local THRUST = 8.0     -- thrust accel px/s^2
local ROT = 2.0        -- rad/s

local countdown = 5.0
local level = 1
local score = 0
local best_level = 1
local best_score = 0
local state = "menu" -- menu | playing
local menu_idx = 1

-- 50 levels, procedurally generated for now
local levels = {}
local level_seed = 0
local time_s = 0
local stars = {}
local heightmap = {}
local TWO_PI = 6.283185307179586
for i=1,50 do
  local t = {}
  local rough = (i-1)/49 -- 0..1
  t.padw = 80 - math.floor(rough*60) -- 80..20
  t.G = 0.25 + 0.75*rough
  t.THRUST = 3.0 + 13.0*rough
  t.ROT = 2.0
  t.land_speed = 12 + math.floor(rough*16)
  t.seed = i*413
  levels[i] = t
end

local function terrain_base_y(x)
  local L = levels[level]
  local rough = (level-1)/49
  local base = math.sin((x+L.seed)*0.06)*(8+rough*18)
  local micro = math.sin((x+L.seed*3)*0.23)*3
  return ground_y - 10 - math.floor(base + micro)
end

local function terrain_y(x)
  if x >= pad_x0 and x <= pad_x1 then return pad_y end
  return heightmap[x] or terrain_base_y(x)
end

local function set_level(idx)
  level = math.max(1, math.min(50, idx))
  local L = levels[level]
  G = L.G; THRUST = L.THRUST; ROT = L.ROT; land_speed = L.land_speed
  level_seed = L.seed
  local function srnd() level_seed = (1103515245*level_seed + 12345) % 2147483648; return level_seed end
  local function frand(a,b) return a + (srnd() % 10000)/10000*(b-a) end
  local padw = L.padw
  local bestx, bestSlope = 240, 1e9
  for _=1,80 do
    local halfw = math.floor(padw/2)
    local cx = math.floor(frand(20+halfw, 460-halfw))
    local x0 = cx - halfw; local x1 = cx + halfw
    local y0 = terrain_base_y(x0); local y1 = terrain_base_y(x1)
    local slope = math.abs(y1 - y0)
    if slope < bestSlope then bestSlope = slope; bestx = cx end
    if slope <= 2 then bestx = cx; break end
  end
  local bestHalf = math.floor(padw/2)
  pad_x0 = bestx - bestHalf; pad_x1 = bestx + bestHalf
  local mid = math.floor((pad_x0+pad_x1)/2)
  pad_y = math.min(terrain_base_y(pad_x0), terrain_base_y(pad_x1), terrain_base_y(mid))
  for x=0,479 do heightmap[x] = terrain_base_y(x) end
  ship.x, ship.y, ship.vx, ship.vy, ship.angle = 240, 60, 0, 0, 0
  ship.fuel, ship.alive, ship.landed = 150, true, false
  countdown = 5.0
end

local function clamp(v, a, b) if v<a then return a elseif v>b then return b else return v end end

function _INIT()
  rf.palette_set("default")
  set_level(1)
  state = "menu"
end

local prevThrust = false
local function normalize_angle()
  if ship.angle > math.pi then ship.angle = ship.angle - TWO_PI end
  if ship.angle < -math.pi then ship.angle = ship.angle + TWO_PI end
end

local start_melody = {"4C1","4E1","4G2","R1","4E1","4C2"}
local function update_menu(dt)
  if rf.btnp(2) then menu_idx = math.max(1, menu_idx-1); rf.sfx("move") end
  if rf.btnp(3) then menu_idx = math.min(2, menu_idx+1); rf.sfx("move") end
  if rf.btnp(4) or rf.btnp(5) then
    if menu_idx == 1 then
      rf.sfx("select"); score, level = 0, 1; set_level(level); state = "playing"; rf.music(start_melody, 120, 0.28)
    else
      rf.sfx("select"); rf.quit()
    end
  end
end

local crash_phase = nil
local crash_timer = 0
local taps_tokens = {"3A2","3F2","3D2","3F3"}
local taps_started = false
local taps_total_dur = 0

local function start_crash_sequence()
  crash_phase = "crash"; crash_timer = 0; taps_started = false; taps_total_dur = 0
end

local function update_crash(dt)
  crash_timer = crash_timer + dt
  if crash_phase == "crash" then
    if crash_timer >= 0.25 then crash_phase = "taps"; crash_timer = 0 end
  elseif crash_phase == "taps" then
    if not taps_started then
      rf.music(taps_tokens, 72, 0.24)
      local sum = 0; for i=1,#taps_tokens do local tk=taps_tokens[i]; local last=string.sub(tk,-1); sum = sum + (tonumber(last) or 1) end
      taps_total_dur = (60/72)*sum + 0.3; taps_started = true
    end
    if crash_timer >= taps_total_dur or crash_timer >= 5.0 then state = "menu"; menu_idx = 1; crash_phase = nil end
  else
    if crash_timer >= 5.0 then state = "menu"; menu_idx = 1 end
  end
end

local function update_countdown(dt)
  local prevInt = math.ceil(countdown); countdown = math.max(0, countdown - dt); local nowInt = math.ceil(countdown)
  if nowInt < prevInt then rf.tone(600,0.06,0.3) end
  ship.y = ship.y + math.sin(time_s*2)*0.2
end

local function update_play(dt)
  if rf.btn(0) then ship.angle = ship.angle + ROT*dt end
  if rf.btn(1) then ship.angle = ship.angle - ROT*dt end
  normalize_angle()
  local thrust = rf.btn(2) and ship.fuel>0
  if thrust then
    local ax = math.sin(ship.angle) * THRUST; local ay = -math.cos(ship.angle) * THRUST
    ship.vx = ship.vx + ax*dt; ship.vy = ship.vy + ay*dt
    ship.fuel = clamp(ship.fuel - 25*dt, 0, 999)
  end
  if thrust ~= prevThrust then rf.sfx("thrust", thrust and "on" or "off"); prevThrust = thrust end
  ship.thrusting = thrust
  ship.vy = ship.vy + G*dt; if ship.vy > 300 then ship.vy = 300 end
  ship.vx = ship.vx * (1 - 0.5*dt)
  ship.x = ship.x + ship.vx; ship.y = ship.y + ship.vy
  if ship.x < 0 then ship.x=0; ship.vx=0 elseif ship.x > 479 then ship.x=479; ship.vx=0 end
  local ground_here = terrain_y(math.floor(ship.x))
  if ship.y >= (ground_here - ship.size) then
    ship.y = ground_here - ship.size
    local speed = math.sqrt(ship.vx*ship.vx + ship.vy*ship.vy)
    local vy_abs = math.abs(ship.vy)
    local angle_ok = math.abs(ship.angle) < 0.2
    local on_pad = ship.x >= pad_x0 and ship.x <= pad_x1
    -- Check vertical velocity separately and use a stricter threshold
    if vy_abs < (land_speed or 18) * 0.6 and speed < (land_speed or 18) and angle_ok and on_pad then
      ship.landed = true; score = score + math.floor(100 + (ship.fuel*2) + math.max(0, (land_speed - speed)*5))
      countdown = 3.0; best_level = math.max(best_level, level); best_score = math.max(best_score, score)
      rf.sfx("land"); set_level(level+1)
    else
      rf.sfx("thrust","off"); rf.sfx("stopall"); ship.alive = false; crash_timer = 0; start_crash_sequence(); rf.sfx("crash")
    end
    ship.vx, ship.vy = 0,0
  end
end

function _UPDATE(dt)
  time_s = time_s + dt
  if state == "menu" then update_menu(dt); return end
  if not ship.alive then update_crash(dt); return end
  if countdown > 0 then update_countdown(dt); return end
  update_play(dt)
end

local function draw_ship()
  local s = ship.size
  local sin, cos = math.sin(ship.angle), math.cos(ship.angle)
  local p0x, p0y = 0, -s
  local p1x, p1y = -s*0.7, s
  local p2x, p2y = s*0.7, s
  local function tx(x,y) return math.floor(ship.x + x*cos + y*sin + 0.5), math.floor(ship.y - x*sin + y*cos + 0.5) end
  local a1x,a1y = tx(p0x,p0y); local a2x,a2y = tx(p1x,p1y); local a3x,a3y = tx(p2x,p2y)
  rf.line(a1x,a1y,a2x,a2y,1); rf.line(a2x,a2y,a3x,a3y,1); rf.line(a3x,a3y,a1x,a1y,1)
  if ship.thrusting then local fx,fy = tx(0,s+2); rf.circfill_rgb(fx,fy,2,255,180,0) end
end

function _DRAW()
  rf.clear_i(0)
  if state == "menu" then draw_menu(); return end
  draw_hud(); draw_level(); draw_ship(); draw_messages()
end

function draw_menu()
  rf.print_center("MOON LANDER", 70, 255,255,255)
  local sel = {r=255,g=255,b=255}; local dim = {r=160,g=160,b=160}
  local c1 = (menu_idx==1) and sel or dim; local c2 = (menu_idx==2) and sel or dim
  rf.print_center("PLAY", 100, c1.r,c1.g,c1.b)
  rf.print_center("QUIT", 116, c2.r,c2.g,c2.b)
  rf.print_center("Up/Down to select, O/X/Enter to confirm", 140, 200,200,200)
  rf.print_center("Controls: Left/Right Rotate, Up Thrust", 156, 200,200,200)
  rf.print_center("Best Level: "..tostring(best_level).."   Best Score: "..tostring(best_score), 176, 200,200,200)
end

function draw_hud()
  rf.print_center("MOON LANDER", 10, 255,255,255)
  rf.print_xy(2, 2,   "ALT:"..string.format("%3.0f", ground_y-ship.y), 2)
  rf.print_xy(2, 10,  "VY :"..string.format("%4.1f", ship.vy), 3)
  rf.print_xy(2, 18,  "FUEL:"..string.format("%3.0f", ship.fuel), 4)
  rf.print_xy(380, 2,  "LVL:"..tostring(level), 1)
  rf.print_xy(380, 10, "SCORE:"..tostring(score), 1)
end

function draw_level()
  local prevx, prevy = 0, terrain_y(0)
  for x=0,479,2 do local ty = terrain_y(x); rf.rectfill(x, ty, x, 269, 2); rf.line(prevx, prevy, x, ty, 2); prevx, prevy = x, ty end
  rf.rectfill(pad_x0, pad_y, pad_x1, pad_y+2, 1)
end

function draw_messages()
  if countdown > 0 and ship.alive and not ship.landed then rf.print_center("GET READY:"..tostring(math.ceil(countdown)), 130, 255,255,255); return end
  if ship.landed then rf.print_center("LANDED!", 130, 255,255,255); return end
  if not ship.alive and not ship.landed then rf.print_center("CRASHED", 120, 255,255,255); rf.print_center("RETURNING TO MENU...", 140, 200,200,200) end
end
