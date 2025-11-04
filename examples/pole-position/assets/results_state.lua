-- Results State - Race results and lap times

local result_selection = 1
local result_options = {
  "Race Again",
  "Next Track",
  "Main Menu"
}

function _INIT()
end

function _ENTER()
  result_selection = 1
  
  -- Update best lap time if this was better
  if lap_time < best_lap_time then
    best_lap_time = lap_time
  end
end

function _HANDLE_INPUT()
  -- Button 0: SELECT or A (6)
  if rf.btnp(0) or rf.btnp(6) then
    rf.sfx("select")
    if result_selection == 1 then
      -- Race Again (same track)
      rf.sfx("stopall")
      game.changeState("play")
    elseif result_selection == 2 then
      -- Next Track
      selected_track = selected_track + 1
      if selected_track > #tracks then
        selected_track = 1
      end
      rf.sfx("stopall")
      game.changeState("play")
    elseif result_selection == 3 then
      -- Main Menu
      rf.sfx("stopall")
      game.changeState("menu")
    end
    return
  end
  
  -- Button 1: START
  if rf.btnp(1) then
    rf.sfx("select")
    -- Quick restart
    rf.sfx("stopall")
    game.changeState("play")
    return
  end
  
  -- Button 2: UP
  if rf.btnp(2) then
    rf.sfx("select")
    result_selection = result_selection - 1
    if result_selection < 1 then
      result_selection = #result_options
    end
    return
  end
  
  -- Button 3: DOWN
  if rf.btnp(3) then
    rf.sfx("select")
    result_selection = result_selection + 1
    if result_selection > #result_options then
      result_selection = 1
    end
    return
  end
end

function _UPDATE(dt)
end

function _DRAW()
  rf.clear_i(COLOR_BLACK)
  
  -- Title (centered, with margin)
  rf.print_anchored("RACE RESULTS", "topcenter", COLOR_WHITE)
  
  -- Track name (centered)
  local track_name = tracks[selected_track].name
  local track_text = "Track: " .. track_name
  local track_x = 240 - string.len(track_text) * 3
  rf.print_xy(track_x, 60, track_text, COLOR_CYAN)
  
  -- Position (centered)
  local pos_text = "Position: " .. current_position .. " / " .. total_cars
  local pos_color = COLOR_YELLOW
  if current_position == 1 then
    pos_color = COLOR_GREEN
  end
  local pos_x = 240 - string.len(pos_text) * 3
  rf.print_xy(pos_x, 90, pos_text, pos_color)
  
  -- Lap time (centered)
  local time_text = string.format("Lap Time: %.2f", lap_time)
  local time_x = 240 - string.len(time_text) * 3
  rf.print_xy(time_x, 120, time_text, COLOR_WHITE)
  
  -- Best lap time (centered)
  local best_text = string.format("Best Time: %.2f", best_lap_time)
  local best_x = 240 - string.len(best_text) * 3
  rf.print_xy(best_x, 140, best_text, COLOR_LIGHT_GRAY)
  
  -- Menu options (centered)
  local start_y = 170
  for i = 1, #result_options do
    local y = start_y + (i - 1) * 40  -- Increased spacing
    local color = COLOR_MEDIUM_GRAY
    local option_text = result_options[i]
    local text_x = 240 - string.len(option_text) * 3  -- Center text
    if i == result_selection then
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

