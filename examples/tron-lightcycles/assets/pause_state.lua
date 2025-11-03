-- Pause State Module for Tron Light Cycles
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
  -- Pause menu doesn't need update logic
end

function _DRAW()
  -- Semi-transparent overlay (simulate with gray box)
  -- Draw a gray box over the play state
  for y=80,190 do
    for x=160,320 do
      rf.pset(x, y, 25) -- COLOR_GRAY
    end
  end
  
  -- Draw menu title
  local pause_x = 240 - string.len("PAUSED")*3
  rf.print_xy(pause_x, 100, "PAUSED", COLOR_WHITE)
  
  -- Draw menu options
  local c1 = (selected_option == 1) and COLOR_WHITE or COLOR_GRAY
  local c2 = (selected_option == 2) and COLOR_WHITE or COLOR_GRAY
  
  local resume_x = 240 - string.len("Resume")*3
  local quit_x = 240 - string.len("Quit")*3
  rf.print_xy(resume_x, 140, "Resume", c1)
  rf.print_xy(quit_x, 160, "Quit", c2)
  
  -- Instructions
  rf.print_xy(240 - string.len("Up/Down: Navigate | SELECT/A: Confirm | START: Resume")*3, 210, "Up/Down: Navigate | SELECT/A: Confirm | START: Resume", COLOR_GRAY)
end

function _EXIT()
  -- Pause menu exit (unused)
end

function _DONE()
  -- Shutdown (unused)
end
