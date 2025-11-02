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
  -- Navigate menu
  if rf.btnp(2) then  -- Up
    selected_option = selected_option - 1
    if selected_option < 1 then
      selected_option = #menu_options
    end
    rf.sfx("move")
  elseif rf.btnp(3) then  -- Down
    selected_option = selected_option + 1
    if selected_option > #menu_options then
      selected_option = 1
    end
    rf.sfx("move")
  end
  
  -- Select option
  if rf.btnp(4) or rf.btnp(5) then  -- Z or X or Enter
    if selected_option == 1 then
      -- Resume - pop pause state to return to play
      rf.sfx("select")
      game.popState()
    elseif selected_option == 2 then
      -- Quit to menu - pop all states and change to menu
      rf.sfx("select")
      game.popAllStates()
      game.changeState("menu")
    end
  end
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
  local resume_color = (selected_option == 1) and COLOR_WHITE or COLOR_DARK_GRAY
  local quit_color = (selected_option == 2) and COLOR_WHITE or COLOR_DARK_GRAY
  
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
  -- Cleanup when leaving pause
end

function _DONE()
  -- Shutdown cleanup
end

