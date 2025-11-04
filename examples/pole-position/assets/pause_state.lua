-- Pause State - Pause overlay

local pause_selection = 1
local pause_options = {
  "Resume",
  "Quit to Menu"
}

function _INIT()
end

function _ENTER()
  pause_selection = 1
end

function _HANDLE_INPUT()
  -- Button 0: SELECT or A (6)
  if rf.btnp(0) or rf.btnp(6) then
    rf.sfx("select")
    if pause_selection == 1 then
      -- Resume
      game.popState()
    elseif pause_selection == 2 then
      -- Quit to Menu
      rf.sfx("stopall")
      game.changeState("menu")
    end
    return
  end
  
  -- Button 1: START
  if rf.btnp(1) then
    rf.sfx("select")
    game.popState()  -- Resume
    return
  end
  
  -- Button 2: UP
  if rf.btnp(2) then
    rf.sfx("select")
    pause_selection = pause_selection - 1
    if pause_selection < 1 then
      pause_selection = #pause_options
    end
    return
  end
  
  -- Button 3: DOWN
  if rf.btnp(3) then
    rf.sfx("select")
    pause_selection = pause_selection + 1
    if pause_selection > #pause_options then
      pause_selection = 1
    end
    return
  end
end

function _UPDATE(dt)
end

function _DRAW()
  -- Draw semi-transparent overlay (simulated with dark fill)
  rf.rectfill(0, 0, 480, 270, COLOR_BLACK)
  
  -- Pause title (centered)
  rf.print_anchored("PAUSED", "middlecenter", COLOR_WHITE)
  
  -- Menu options (centered)
  local start_y = 140
  for i = 1, #pause_options do
    local y = start_y + (i - 1) * 45  -- Increased spacing
    local color = COLOR_MEDIUM_GRAY
    local option_text = pause_options[i]
    local text_x = 240 - string.len(option_text) * 3  -- Center text
    if i == pause_selection then
      color = COLOR_CYAN
      rf.print_xy(text_x - 20, y, ">", color)
    end
    rf.print_xy(text_x, y, option_text, color)
  end
  
  -- Instructions (centered, bottom with margin)
  rf.print_anchored("UP/DOWN: Navigate  SELECT/A: Choose", "bottomcenter", COLOR_GRAY)
end

function _EXIT()
end

function _DONE()
  -- Shutdown (unused)
end

