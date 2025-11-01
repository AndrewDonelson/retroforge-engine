-- Play State Module
-- Physics demo with bouncing balls and sprite pooling

-- Module-level variables (persist across enter/exit)
local balls = {}
local ball_count = 0
local balls_max = 1024
local lifetime_min = 30
local lifetime_max = 120
local spawn_timer = 0
local spawn_delay = 0.1  -- Spawn a ball every 0.2 seconds
local play_timer = 0
-- No time limit - user controls when to exit
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

-- Physics bodies
local bodies = {}
local world_bounds = {
  left = 10,
  right = 470,
  top = 10,
  bottom = 260
}

function _INIT()
  -- Create world boundaries using physics API with high restitution for bouncing
  -- Box2D SetAsBox uses half-width/half-height, so positions need to account for that
  -- Left wall: position at x = world_bounds.left (10), extends 5 pixels each side
  local left_x = world_bounds.left
  local left_body = rf.physics_create_body("static", left_x, 135)
  rf.physics_body_add_box(left_body, 10, 280, 0, 1.0, 0.1)
  
  -- Right wall: position at x = world_bounds.right (470), extends 5 pixels each side
  local right_x = world_bounds.right
  local right_body = rf.physics_create_body("static", right_x, 135)
  rf.physics_body_add_box(right_body, 10, 280, 0, 1.0, 0.1)
  
  -- Top wall: position at y = world_bounds.top (10), extends 5 pixels up/down
  local top_y = world_bounds.top
  local top_body = rf.physics_create_body("static", 240, top_y)
  rf.physics_body_add_box(top_body, 480, 10, 0, 1.0, 0.1)
  
  -- Bottom wall: position at y = world_bounds.bottom (260), extends 5 pixels up/down
  local bottom_y = world_bounds.bottom
  local bottom_body = rf.physics_create_body("static", 240, bottom_y)
  rf.physics_body_add_box(bottom_body, 480, 10, 0, 1.0, 0.1)
  
  table.insert(bodies, left_body)
  table.insert(bodies, right_body)
  table.insert(bodies, top_body)
  table.insert(bodies, bottom_body)
end

function _ENTER()
  -- Reset when entering play state
  balls = {}
  ball_count = 0
  spawn_timer = 0
  play_timer = 0
  stats.max_objects = 0
  stats.total_spawned = 0
  fps_counter = 0
  fps_accumulator = 0
  
  -- Don't spawn initial ball here - spawn it in first UPDATE
  -- This ensures physics step runs after ball is created
end

function _HANDLE_INPUT()
  -- Return to menu: Press Z to go back to menu
  if rf.btnp(4) then  -- Z button
    game.changeState("menu")  -- Go back to menu, NOT exit
    return
  end
  
  -- Spawn ball on arrow keys from that direction (buttons: 0=Left, 1=Right, 2=Up, 3=Down)
  -- Only spawn if under limit (active ball count)
  if ball_count < balls_max then
    if rf.btnp(2) then  -- Up arrow
      spawnBallFromDirection("up")
    elseif rf.btnp(3) then  -- Down arrow
      spawnBallFromDirection("down")
    elseif rf.btnp(0) then  -- Left arrow
      spawnBallFromDirection("left")
    elseif rf.btnp(1) then  -- Right arrow
      spawnBallFromDirection("right")
    elseif rf.btnp(5) then  -- X button spawns from center
      spawnBall()
    end
  end
end

function spawnBall()
  -- Check maxSpawn limit (balls_max) - don't spawn if already at limit
  if ball_count >= balls_max then
    return  -- Don't spawn, wait for balls to expire
  end
  
  -- Spawn a new ball from center
  local x = 240  -- Center X
  local y = 135  -- Center Y
  
  -- Create physics body
  local body = rf.physics_create_body("dynamic", x, y)
  
  -- Create circle fixture with high restitution (bounciness) and low friction
  rf.physics_body_add_circle(body, 3, 1.0, 0.9, 0.1)
  
  -- Apply random initial velocity
  local vx = (math.random() - 0.5) * 300  -- -150 to 150
  local vy = (math.random() - 0.5) * 300  -- -150 to 150
  rf.physics_body_set_velocity(body, vx, vy)
  
  -- Random lifetime between 15 and 30 seconds
  local lifetime = math.random() * lifetime_min + lifetime_max
  
  -- Store ball data
  local color = colors[color_index]
  color_index = color_index + 1
  if color_index > #colors then
    color_index = 1
  end
  
  table.insert(balls, {
    body = body,
    color = color,
    lifetime = lifetime,
    age = 0
  })
  
  ball_count = ball_count + 1
  stats.total_spawned = stats.total_spawned + 1
  stats.max_objects = math.max(stats.max_objects, ball_count)
end

function spawnBallFromDirection(direction)
  -- Check maxSpawn limit (balls_max)
  if ball_count >= balls_max then
    return
  end
  
  -- Screen center
  local center_x = 240
  local center_y = 135
  
  -- Spawn position based on direction
  local x, y
  local vx, vy
  
  if direction == "up" then
    -- Spawn from top, move toward center (down)
    x = center_x + (math.random() - 0.5) * 100  -- Random X near center
    y = world_bounds.top + 5  -- Top edge
    -- Velocity toward center (down and toward center)
    vx = (center_x - x) * 0.3 + (math.random() - 0.5) * 50  -- Toward center X with randomness
    vy = math.random() * 150 + 100  -- Down toward center (100-250)
  elseif direction == "down" then
    -- Spawn from bottom, move toward center (up)
    x = center_x + (math.random() - 0.5) * 100
    y = world_bounds.bottom - 5  -- Bottom edge
    vx = (center_x - x) * 0.3 + (math.random() - 0.5) * 50
    vy = -(math.random() * 150 + 100)  -- Up toward center
  elseif direction == "left" then
    -- Spawn from left, move toward center (right)
    x = world_bounds.left + 5  -- Left edge
    y = center_y + (math.random() - 0.5) * 100
    vx = math.random() * 150 + 100  -- Right toward center
    vy = (center_y - y) * 0.3 + (math.random() - 0.5) * 50  -- Toward center Y
  elseif direction == "right" then
    -- Spawn from right, move toward center (left)
    x = world_bounds.right - 5  -- Right edge
    y = center_y + (math.random() - 0.5) * 100
    vx = -(math.random() * 150 + 100)  -- Left toward center
    vy = (center_y - y) * 0.3 + (math.random() - 0.5) * 50
  else
    -- Invalid direction, use center spawn
    return spawnBall()
  end
  
  -- Create physics body
  local body = rf.physics_create_body("dynamic", x, y)
  
  -- Create circle fixture
  rf.physics_body_add_circle(body, 3, 1.0, 0.9, 0.1)
  
  -- Apply velocity toward center
  rf.physics_body_set_velocity(body, vx, vy)
  
  -- Random lifetime
  local lifetime = math.random() * lifetime_min + lifetime_max
  
  -- Store ball data
  local color = colors[color_index]
  color_index = color_index + 1
  if color_index > #colors then
    color_index = 1
  end
  
  table.insert(balls, {
    body = body,
    color = color,
    lifetime = lifetime,
    age = 0
  })
  
  ball_count = ball_count + 1
  stats.total_spawned = stats.total_spawned + 1
  stats.max_objects = math.max(stats.max_objects, ball_count)
end

function _UPDATE(dt)
  -- Spawn initial ball on first frame if none exist
  -- Physics.Step() runs BEFORE Update(), so ball created here will move next frame
  -- That's fine - it will start moving immediately on second frame
  if ball_count == 0 then
    spawnBall()
  end
  
  -- Auto-spawn balls periodically (only if under limit)
  spawn_timer = spawn_timer + dt
  if spawn_timer >= spawn_delay then
    spawn_timer = spawn_timer - spawn_delay
    -- Only check active ball_count, total_spawned is just a stat counter (no limit)
    if ball_count < balls_max then
      spawnBall()
    end
  end
  
  -- Update physics (engine handles this, but we can update stats)
  stats.objects = ball_count
  
  -- Update play timer (just for display, no limit)
  play_timer = play_timer + dt
  
  -- Calculate FPS manually (more reliable than dev mode stat)
  fps_accumulator = fps_accumulator + dt
  fps_counter = fps_counter + 1
  if fps_accumulator >= 0.5 then  -- Update FPS every 0.5 seconds
    stats.fps = fps_counter / fps_accumulator
    fps_accumulator = 0
    fps_counter = 0
  end
  
  -- Try to get memory from dev mode if available
  if rf.stat then
    stats.memory = rf.stat(2) or 0  -- Memory is stat(2)
  end
  
  -- Update ball lifetimes and remove expired balls
  for i = #balls, 1, -1 do
    local ball = balls[i]
    
    -- Update age
    ball.age = ball.age + dt
    
    -- Check if lifetime expired
    if ball.age >= ball.lifetime then
      -- Destroy physics body and remove ball
      rf.physics_body_destroy(ball.body)
      table.remove(balls, i)
      ball_count = ball_count - 1
    else
      -- Check if way out of bounds (emergency cleanup)
      local x, y = rf.physics_body_get_position(ball.body)
      if y > 300 or x < -50 or x > 530 then
        rf.physics_body_destroy(ball.body)
        table.remove(balls, i)
        ball_count = ball_count - 1
      end
    end
  end
  
  -- No time limit - game runs until user exits
end

function _DRAW()
  -- Clear screen
  rf.clear_i(0)
  
  -- Draw boundaries
  rf.rect(world_bounds.left, world_bounds.top, 
          world_bounds.right, world_bounds.bottom, 6)
  
  -- Draw balls
  for _, ball in ipairs(balls) do
    local x, y = rf.physics_body_get_position(ball.body)
    
    -- Draw using sprite (pooled sprite) - ball sprite is still 12x12 but we draw smaller
    -- Radius is now 3, so offset by 6 (sprite center) minus 3 (radius) = 3
    rf.spr("ball", math.floor(x - 6), math.floor(y - 6))
    
    -- Draw colored circle overlay with new radius (3 instead of 6)
    rf.circfill(math.floor(x), math.floor(y), 3, ball.color)
  end
  
  -- Draw title at topleft (with padding)
  rf.print_xy(10, 5, "KITCHEN SINK DEMO - PLAYING", 15)
  
  -- Draw stats on the right side to avoid overlapping with physics box
  local stats_x = 300  -- Right side of screen
  local stats_y = 15
  rf.print_xy(stats_x, stats_y, string.format("FPS: %.1f", stats.fps), 7)
  stats_y = stats_y + 15
  rf.print_xy(stats_x, stats_y, string.format("Memory: %d KB", math.floor(stats.memory / 1024)), 7)
  stats_y = stats_y + 15
  rf.print_xy(stats_x, stats_y, string.format("Objects: %d", stats.objects), 7)
  stats_y = stats_y + 15
  rf.print_xy(stats_x, stats_y, string.format("Max Objects: %d", stats.max_objects), 7)
  stats_y = stats_y + 15
  rf.print_xy(stats_x, stats_y, string.format("Total Spawned: %d", stats.total_spawned), 7)
  stats_y = stats_y + 15
  rf.print_xy(stats_x, stats_y, string.format("Time: %.1f", play_timer), 7)
  stats_y = stats_y + 15
  rf.print_xy(stats_x, stats_y, "Arrows: Spawn | Z: Menu | X: Center", 6)
  
  -- Draw status message at topright
  rf.print_anchored("POOLED SPRITES", "topright", 11)
  
  -- Draw info at bottomleft
  rf.print_anchored("Balls use pooled sprites!", "bottomleft", 11)
  
  -- Draw info at bottomright (show active limit, not total spawn limit)
  rf.print_anchored("(Active: " .. ball_count .. "/" ..balls_max..")", "bottomright", 11)
  
  -- Draw center indicator (optional)
  rf.print_anchored("PHYSICS DEMO", "middlecenter", 3)
end

function _EXIT()
  -- Cleanup physics bodies when leaving
  for _, ball in ipairs(balls) do
    -- Bodies will be cleaned up by engine
  end
  balls = {}
  ball_count = 0
end

function _DONE()
  -- Shutdown cleanup
  bodies = {}
end
