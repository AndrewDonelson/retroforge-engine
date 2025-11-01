-- Galaxy Simulation State
-- Creates a spiral galaxy using central gravity physics

-- Module-level variables
local stars = {}
local star_count = 0
local stars_max = 2048
local spawn_timer = 0
local spawn_delay = 0.025  -- Spawn a star every 0.05 seconds (fast)
local play_timer = 0
local colors = {11, 12, 13, 14, 15, 3, 4, 5, 6, 7}  -- Various colors
local color_index = 1
local stats = {
  fps = 0,
  memory = 0,
  objects = 0,
  max_objects = 0,
  total_spawned = 0
}
local fps_counter = 0
local fps_accumulator = 0

-- Galaxy center (screen center)
local center_x = 240
local center_y = 135

-- Gravity parameters
local gravity_strength = 800.0  -- Strength of central gravity (balanced for stable orbits)
local min_orbital_radius = 30   -- Minimum distance from center to spawn
local max_orbital_radius = 200  -- Maximum initial distance from center
local central_mass = 3000.0     -- Virtual central mass for orbital calculations

-- Screen bounds for spawning
local screen_width = 480
local screen_height = 270

function _INIT()
  -- No boundaries - open space for galaxy formation
end

function _ENTER()
  -- Reset when entering play state
  stars = {}
  star_count = 0
  spawn_timer = 0
  play_timer = 0
  stats.max_objects = 0
  stats.total_spawned = 0
  fps_counter = 0
  fps_accumulator = 0
  color_index = 1
end

function _HANDLE_INPUT()
  -- Return to menu: Press Z to go back to menu
  if rf.btnp(4) then  -- Z button
    game.changeState("menu")
    return
  end
end

function spawnStar()
  -- Check maxSpawn limit
  if star_count >= stars_max then
    return
  end
  
  -- Spawn in spiral pattern to create galaxy structure
  -- Mix of random edge spawns and spiral arm spawns
  local spawn_x, spawn_y
  local use_spiral = math.random() < 0.7  -- 70% spawn in spiral arms, 30% random edge
  
  if use_spiral then
    -- Spawn along spiral arm (logarithmic spiral)
    local spiral_base_angle = (stats.total_spawned * 0.15) % (2 * math.pi)  -- Rotating spiral
    local r = min_orbital_radius + math.random() * (max_orbital_radius - min_orbital_radius)
    
    -- Logarithmic spiral: θ = b * log(r/a) + θ0
    local a = min_orbital_radius
    local b = 0.2  -- Spiral tightness
    local theta = spiral_base_angle + b * math.log(r / a + 1)
    
    spawn_x = center_x + r * math.cos(theta)
    spawn_y = center_y + r * math.sin(theta)
    
    -- Ensure within screen bounds
    spawn_x = math.max(10, math.min(screen_width - 10, spawn_x))
    spawn_y = math.max(10, math.min(screen_height - 10, spawn_y))
  else
    -- Random edge spawn
    local edge = math.random(4)
    if edge == 1 then  -- Top
      spawn_x = math.random(0, screen_width)
      spawn_y = 5
    elseif edge == 2 then  -- Right
      spawn_x = screen_width - 5
      spawn_y = math.random(0, screen_height)
    elseif edge == 3 then  -- Bottom
      spawn_x = math.random(0, screen_width)
      spawn_y = screen_height - 5
    else  -- Left
      spawn_x = 5
      spawn_y = math.random(0, screen_height)
    end
  end
  
  -- Create physics body
  local body = rf.physics_create_body("dynamic", spawn_x, spawn_y)
  
  -- Disable Box2D gravity (we'll use custom central gravity)
  rf.physics_body_set_gravity_scale(body, 0.0)
  
  -- Create circle fixture - very small, low density
  rf.physics_body_add_circle(body, 1.5, 0.05, 0.0, 0.0)  -- Very light, no bounce, no friction
  
  -- Calculate distance to center
  local dx = center_x - spawn_x
  local dy = center_y - spawn_y
  local dist_to_center = math.sqrt(dx * dx + dy * dy)
  
  -- Normalize radial direction
  local dir_x = dx / dist_to_center
  local dir_y = dy / dist_to_center
  
  -- Calculate tangential direction (perpendicular, for orbital motion)
  local tan_x = -dir_y
  local tan_y = dir_x
  
  -- Calculate orbital velocity for stable circular orbit
  -- v = sqrt(GM/r) for circular orbit
  -- For stable orbit: centripetal force = gravitational force
  -- mv²/r = GMm/r² => v = sqrt(GM/r)
  -- Stars closer to center (smaller r) have HIGHER orbital velocity (faster)
  local GM = gravity_strength * central_mass  -- Gravitational parameter
  local orbital_velocity = math.sqrt(GM / dist_to_center)
  
  -- Stars closer to center orbit faster (differential rotation - Milky Way characteristic)
  -- Keep orbital velocity close to ideal for stable circular orbits
  -- Add slight randomness to create variety in orbits (slightly elliptical)
  local speed_variation = 0.9 + math.random() * 0.2  -- 90-110% of orbital speed
  local tangential_speed = orbital_velocity * speed_variation
  
  -- NO radial component - pure tangential (orbital) motion
  -- This creates stable circular/elliptical orbits that maintain distance from center
  -- Stars will orbit rather than spiral inward
  local radial_speed = 0
  
  -- Combine velocities: pure tangential (orbital motion perpendicular to radius)
  local vx = tan_x * tangential_speed
  local vy = tan_y * tangential_speed
  
  rf.physics_body_set_velocity(body, vx, vy)
  
  -- Random size factor
  local size_factor = 0.6 + math.random() * 0.4  -- 0.6 to 1.0
  
  -- Store star data
  local color = colors[color_index]
  color_index = color_index + 1
  if color_index > #colors then
    color_index = 1
  end
  
  table.insert(stars, {
    body = body,
    color = color,
    size = size_factor
  })
  
  star_count = star_count + 1
  stats.total_spawned = stats.total_spawned + 1
  stats.max_objects = math.max(stats.max_objects, star_count)
end

function _UPDATE(dt)
  -- Spawn stars automatically to build the galaxy
  spawn_timer = spawn_timer + dt
  if spawn_timer >= spawn_delay then
    spawn_timer = spawn_timer - spawn_delay
    if star_count < stars_max then
      spawnStar()
    end
  end
  
  -- Apply central gravity to all stars (inverse square law)
  -- This creates the centripetal force that keeps stars in orbit
  for _, star in ipairs(stars) do
    local x, y = rf.physics_body_get_position(star.body)
    
    -- Calculate direction to center
    local dx = center_x - x
    local dy = center_y - y
    local dist_sq = dx * dx + dy * dy
    local dist = math.sqrt(dist_sq)
    
    -- Avoid division by zero and extreme forces at center
    if dist > 2 then
      -- Normalize direction
      local dir_x = dx / dist
      local dir_y = dy / dist
      
      -- Apply gravitational acceleration using inverse square law
      -- F = GM/r^2, so acceleration = GM/r^2
      -- This provides the centripetal force needed for circular/elliptical orbits
      local GM = gravity_strength * central_mass
      local force_scale = GM / (dist * dist)
      
      -- Get current velocity
      local vx, vy = rf.physics_body_get_velocity(star.body)
      
      -- Apply gravity acceleration (centripetal force toward center)
      -- This pulls stars toward center, while their tangential velocity keeps them orbiting
      vx = vx + dir_x * force_scale * dt
      vy = vy + dir_y * force_scale * dt
      
      -- Update velocity
      rf.physics_body_set_velocity(star.body, vx, vy)
    end
  end
  
  -- Update stats
  stats.objects = star_count
  play_timer = play_timer + dt
  
  -- Calculate FPS
  fps_accumulator = fps_accumulator + dt
  fps_counter = fps_counter + 1
  if fps_accumulator >= 0.5 then
    stats.fps = fps_counter / fps_accumulator
    fps_accumulator = 0
    fps_counter = 0
  end
  
  -- Try to get memory from dev mode if available
  if rf.stat then
    stats.memory = rf.stat(2) or 0
  end
  
  -- Remove stars that are way out of bounds (emergency cleanup)
  for i = #stars, 1, -1 do
    local star = stars[i]
    local x, y = rf.physics_body_get_position(star.body)
    
    -- Check if way out of bounds
    if x < -100 or x > screen_width + 100 or y < -100 or y > screen_height + 100 then
      rf.physics_body_destroy(star.body)
      table.remove(stars, i)
      star_count = star_count - 1
    end
  end
end

function _DRAW()
  -- Clear screen (black for space)
  rf.clear_i(0)
  
  -- Draw galaxy center (bright core)
  rf.circfill(center_x, center_y, 3, 15)
  rf.circfill(center_x, center_y, 2, 14)
  rf.circfill(center_x, center_y, 1, 11)
  
  -- Draw stars (sort by distance to draw closer ones on top)
  local star_list = {}
  for _, star in ipairs(stars) do
    local x, y = rf.physics_body_get_position(star.body)
    local dx = center_x - x
    local dy = center_y - y
    local dist_sq = dx * dx + dy * dy
    table.insert(star_list, {star = star, dist_sq = dist_sq, x = x, y = y})
  end
  
  -- Sort by distance (farther first, so closer stars draw on top)
  table.sort(star_list, function(a, b) return a.dist_sq > b.dist_sq end)
  
  -- Draw stars
  for _, item in ipairs(star_list) do
    local star = item.star
    local x, y = item.x, item.y
    
    -- Draw star with size variation
    local draw_radius = math.max(1, math.floor(1.5 * star.size + 0.5))
    rf.circfill(math.floor(x), math.floor(y), draw_radius, star.color)
    
    -- Optional: draw small trail for fast-moving stars
    local vx, vy = rf.physics_body_get_velocity(star.body)
    local speed = math.sqrt(vx * vx + vy * vy)
    if speed > 40 and item.dist_sq > 100 then  -- Only for fast stars away from center
      local trail_length = math.min(2, speed / 40)
      local trail_x = x - (vx / speed) * trail_length
      local trail_y = y - (vy / speed) * trail_length
      rf.circfill(math.floor(trail_x), math.floor(trail_y), 1, star.color)
    end
  end
  
  -- Draw title
  rf.print_xy(10, 5, "GALAXY SIMULATION", 15)
  
  -- Draw stats on the right side
  local stats_x = 300
  local stats_y = 15
  rf.print_xy(stats_x, stats_y, string.format("FPS: %.1f", stats.fps), 7)
  stats_y = stats_y + 15
  rf.print_xy(stats_x, stats_y, string.format("Time: %.1f", play_timer), 7)
  stats_y = stats_y + 15
  rf.print_xy(stats_x, stats_y, "Z: Menu", 6)
  
  -- Draw info at bottom
  rf.print_anchored("Spiral Galaxy Formation", "bottomcenter", 11)
  rf.print_anchored("(Stars: " .. star_count .. "/" .. stars_max .. ")", "bottomright", 11)
end

function _EXIT()
  -- Cleanup when leaving
  for _, star in ipairs(stars) do
    rf.physics_body_destroy(star.body)
  end
  stars = {}
  star_count = 0
end

function _DONE()
  -- Shutdown cleanup
end
