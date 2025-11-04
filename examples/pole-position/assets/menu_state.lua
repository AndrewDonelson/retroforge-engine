-- Menu State - Main menu with navigation

local menu_selection = 1
local menu_options = {
  "Start Race",
  "Track Select",
  "Exit"
}

function _INIT()
end

function _ENTER()
  menu_selection = 1
end

function _HANDLE_INPUT()
  -- Button 0: SELECT
  if rf.btnp(0) or rf.btnp(6) then  -- SELECT or A
    rf.sfx("select")
    if menu_selection == 1 then
      -- Start Race - go to track select
      rf.sfx("stopall")
      game.changeState("track_select")
    elseif menu_selection == 2 then
      -- Track Select
      rf.sfx("stopall")
      game.changeState("track_select")
    elseif menu_selection == 3 then
      -- Exit
      rf.sfx("stopall")
      game.exit()
    end
    return
  end
  
  -- Button 1: START
  if rf.btnp(1) then
    rf.sfx("select")
    rf.sfx("stopall")
    -- Start immediately with first track
    selected_track = 1
    game.changeState("play")
    return
  end
  
  -- Button 2: UP
  if rf.btnp(2) then
    rf.sfx("select")
    menu_selection = menu_selection - 1
    if menu_selection < 1 then
      menu_selection = #menu_options
    end
    return
  end
  
  -- Button 3: DOWN
  if rf.btnp(3) then
    rf.sfx("select")
    menu_selection = menu_selection + 1
    if menu_selection > #menu_options then
      menu_selection = 1
    end
    return
  end
end

function _UPDATE(dt)
  -- Music should already be playing from splash
end

function _DRAW()
  rf.clear_i(COLOR_BLACK)
  
  -- Title (centered, with margin)
  rf.print_anchored("POLE POSITION", "topcenter", COLOR_WHITE)
  
  -- Menu options (centered)
  local start_y = 120
  for i = 1, #menu_options do
    local y = start_y + (i - 1) * 45  -- Increased spacing
    local color = COLOR_MEDIUM_GRAY
    local option_text = menu_options[i]
    local text_x = 240 - string.len(option_text) * 3  -- Center text
    if i == menu_selection then
      color = COLOR_CYAN
      rf.print_xy(text_x - 20, y, ">", color)
    end
    rf.print_xy(text_x, y, option_text, color)
  end
  
  -- Instructions (centered, bottom with margin)
  local inst1 = "UP/DOWN: Navigate"
  local inst1_x = 240 - string.len(inst1) * 3
  rf.print_xy(inst1_x, 250, inst1, COLOR_GRAY)
  local inst2 = "SELECT/A: Choose  START: Quick Start"
  local inst2_x = 240 - string.len(inst2) * 3
  rf.print_xy(inst2_x, 235, inst2, COLOR_GRAY)
end

function _EXIT()
  rf.sfx("stopall")
end

function _DONE()
  -- Shutdown (unused)
end

