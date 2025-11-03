-- Game Over State Module for Tron Light Cycles

local gameover_timer = 0.0

function _INIT()
  -- Module initialization
end

function _ENTER()
  -- Reset timer when entering
  gameover_timer = 0.0
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
  
  -- Wait 3 seconds before accepting input
  if gameover_timer >= 3.0 then
    -- Button 0: SELECT
    if rf.btnp(0) then
      rf.sfx("select")
      rf.sfx("stopall") -- Stop taps music
      game.changeState("menu")
      return
    end
    
    -- Button 1: START
    if rf.btnp(1) then
      rf.sfx("select")
      rf.sfx("stopall") -- Stop taps music
      game.changeState("menu")
      return
    end
    
    -- Button 2: UP
    -- Not used in gameover state
    
    -- Button 3: DOWN
    -- Not used in gameover state
    
    -- Button 4: LEFT
    -- Not used in gameover state
    
    -- Button 5: RIGHT
    -- Not used in gameover state
    
    -- Button 6: A
    if rf.btnp(6) then
      rf.sfx("select")
      rf.sfx("stopall") -- Stop taps music
      game.changeState("menu")
      return
    end
    
    -- Button 7: B
    -- Not used in gameover state
    
    -- Button 8: X
    -- Not used in gameover state
    
    -- Button 9: Y
    -- Not used in gameover state
    
    -- Button 10: TURBO
    -- Not used in gameover state
  end
end

function _UPDATE(dt)
  gameover_timer = gameover_timer + dt
end

function _DRAW()
  -- Clear screen
  rf.clear_i(COLOR_BLACK)
  
  -- Redraw game state
  for y=0,GRID_H-1,5 do
    for x=0,GRID_W-1,5 do
      local sx, sy = grid_to_screen(x, y)
      rf.pset(sx, sy, COLOR_WHITE)
    end
  end
  
  -- Draw cycles (frozen state)
  for i=1,#enemies do
    draw_cycle(enemies[i], false)
  end
  
  if player then
    draw_cycle(player, true)
  end
  
  -- Draw HUD
  draw_hud()
  
  -- Draw game over message
  local gameover_x = 240 - string.len("GAME OVER")*3
  local level_text = "Level: " .. tostring(level)
  local level_x = 240 - string.len(level_text)*3
  local score_text = "Score: " .. tostring(score)
  local score_x = 240 - string.len(score_text)*3
  rf.print_xy(gameover_x, 110, "GAME OVER", COLOR_RED)
  rf.print_xy(level_x, 130, level_text, COLOR_WHITE)
  rf.print_xy(score_x, 150, score_text, COLOR_WHITE)
  
  if gameover_timer >= 3.0 then
    local continue_text = "Press SELECT/START/A to continue"
    local continue_x = 240 - string.len(continue_text)*3
    rf.print_xy(continue_x, 170, continue_text, COLOR_GRAY)
  end
end

function _EXIT()
  -- Stop all audio when leaving gameover (taps music should already be stopped, but ensure cleanup)
  rf.sfx("stopall")
end

function _DONE()
  -- Shutdown cleanup
end
