-- Pause State Module for Moon Lander
-- Overlay state that pauses the game when pushed

local selected_option = 1
local menu_options = {
  "Resume",
  "Quit"
}

function _INIT()
  -- Module initialization
end

function _ENTER()
  -- Reset menu when entering
  selected_option = 1
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
    -- SELECT confirms menu selection
    if selected_option == 1 then
      -- Resume - pop pause state to return to play
      rf.sfx("select")
      game.popState()
    elseif selected_option == 2 then
      -- Quit to menu - pop all states and change to menu
      rf.sfx("select")
      rf.sfx("stopall") -- Stop all audio
      game.popAllStates()
      game.changeState("menu")
    end
  end
  
  -- Button 1: START
  if rf.btnp(1) then
    -- START resumes immediately (bypasses menu)
    rf.sfx("select")
    game.popState()
  end
  
  -- Button 2: UP
  if rf.btnp(2) then
    selected_option = selected_option - 1
    if selected_option < 1 then
      selected_option = #menu_options
    end
    rf.sfx("select")
  end
  
  -- Button 3: DOWN
  if rf.btnp(3) then
    selected_option = selected_option + 1
    if selected_option > #menu_options then
      selected_option = 1
    end
    rf.sfx("select")
  end
  
  -- Button 4: LEFT
  -- Not used in pause menu
  
  -- Button 5: RIGHT
  -- Not used in pause menu
  
  -- Button 6: A
  if rf.btnp(6) then
    -- A button also confirms menu selection (same as SELECT)
    if selected_option == 1 then
      -- Resume - pop pause state to return to play
      rf.sfx("select")
      game.popState()
    elseif selected_option == 2 then
      -- Quit to menu - pop all states and change to menu
      rf.sfx("select")
      rf.sfx("stopall") -- Stop all audio
      game.popAllStates()
      game.changeState("menu")
    end
  end
  
  -- Button 7: B
  -- Not used in pause menu
  
  -- Button 8: X
  -- Not used in pause menu
  
  -- Button 9: Y
  -- Not used in pause menu
  
  -- Button 10: TURBO
  -- Not used in pause menu
end

function _UPDATE(dt)
  -- Pause state doesn't need update logic
end

function _DRAW()
  -- Draw semi-transparent overlay
  rf.rectfill(0, 0, 479, 269, COLOR_BLACK)
  
  -- Draw bordered box for pause menu
  local box_x0, box_y0, box_x1, box_y1 = 140, 80, 340, 190
  rf.rectfill(box_x0, box_y0, box_x1, box_y1, COLOR_GRAY)
  rf.rect(box_x0-2, box_y0-2, box_x1+2, box_y1+2, COLOR_WHITE)
  
  -- Pause title
  rf.print_anchored("PAUSED", "topcenter", COLOR_LIGHT_GRAY)
  
  -- Menu options
  -- Use medium gray (30) for unselected, white for selected
  local COLOR_MEDIUM_GRAY = 30
  local resume_color = (selected_option == 1) and COLOR_WHITE or COLOR_MEDIUM_GRAY
  local quit_color = (selected_option == 2) and COLOR_WHITE or COLOR_MEDIUM_GRAY
  
  local resume_x = 240 - (string.len("RESUME") * 3)
  local quit_x = 240 - (string.len("QUIT") * 3)
  
  -- Selection indicator
  if selected_option == 1 then
    rf.print_xy(resume_x - 18, 125, ">", COLOR_WHITE)
  else
    rf.print_xy(quit_x - 18, 141, ">", COLOR_WHITE)
  end
  
  rf.print_xy(resume_x, 125, "RESUME", resume_color)
  rf.print_xy(quit_x, 141, "QUIT", quit_color)
  
  -- Instructions
  rf.print_xy(240 - string.len("Up/Down: Navigate | Enter/O/X: Confirm")*3, 210, "Up/Down: Navigate | Enter/O/X: Confirm", COLOR_GRAY)
end

function _EXIT()
  -- Stop all audio when exiting pause state (safety cleanup)
  rf.sfx("stopall")
end

function _DONE()
  -- Shutdown cleanup
end

