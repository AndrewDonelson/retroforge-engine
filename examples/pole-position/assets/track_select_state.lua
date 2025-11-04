-- Track Select State - Choose from 10 tracks

local scroll_offset = 0
local max_scroll = 0

function _INIT()
end

function _ENTER()
  selected_track = 1
  scroll_offset = 0
  -- Calculate max scroll (show 5 tracks at a time)
  max_scroll = math.max(0, (#tracks - 5) * 50)
end

function _HANDLE_INPUT()
  -- Button 0: SELECT or A (6)
  if rf.btnp(0) or rf.btnp(6) then
    rf.sfx("select")
    rf.sfx("stopall")
    game.changeState("play")
    return
  end
  
  -- Button 1: START
  if rf.btnp(1) then
    rf.sfx("select")
    rf.sfx("stopall")
    game.changeState("play")
    return
  end
  
  -- Button 2: UP
  if rf.btnp(2) then
    rf.sfx("select")
    selected_track = selected_track - 1
    if selected_track < 1 then
      selected_track = #tracks
    end
    -- Update scroll
    if selected_track * 50 - scroll_offset < 50 then
      scroll_offset = math.max(0, selected_track * 50 - 50)
    end
    return
  end
  
  -- Button 3: DOWN
  if rf.btnp(3) then
    rf.sfx("select")
    selected_track = selected_track + 1
    if selected_track > #tracks then
      selected_track = 1
    end
    -- Update scroll
    if selected_track * 50 - scroll_offset > 200 then
      scroll_offset = math.min(max_scroll, selected_track * 50 - 200)
    end
    return
  end
end

function _UPDATE(dt)
end

function _DRAW()
  rf.clear_i(COLOR_BLACK)
  
  -- Title (centered, with margin)
  rf.print_anchored("SELECT TRACK", "topcenter", COLOR_WHITE)
  
  -- Track list (centered)
  local start_y = 70  -- Increased margin from top
  local visible_start = math.floor(scroll_offset / 50) + 1
  local visible_end = math.min(#tracks, visible_start + 4)
  
  for i = visible_start, visible_end do
    local y = start_y + (i - visible_start) * 45 - (scroll_offset % 50)  -- Increased spacing
    if y >= 60 and y < 240 then
      local track = tracks[i]
      local color = COLOR_MEDIUM_GRAY
      local prefix = "  "
      
      if i == selected_track then
        color = COLOR_CYAN
        prefix = "> "
        -- Highlight background (centered)
        local highlight_w = 360
        local highlight_x = 240 - highlight_w / 2
        rf.rectfill(highlight_x, y - 5, highlight_x + highlight_w, y + 35, COLOR_DARK_GRAY)
      end
      
      -- Track number and name (centered)
      local track_text = prefix .. i .. ". " .. track.name
      local track_x = 240 - string.len(track_text) * 3
      rf.print_xy(track_x, y, track_text, color)
      
      -- Difficulty indicator (right aligned)
      local diff_text = ""
      for j = 1, track.difficulty do
        diff_text = diff_text .. "*"
      end
      local diff_x = 400 - string.len(diff_text) * 3
      rf.print_xy(diff_x, y, diff_text, COLOR_YELLOW)
    end
  end
  
  -- Instructions (centered, bottom with margin)
  local inst1 = "UP/DOWN: Select Track"
  local inst1_x = 240 - string.len(inst1) * 3
  rf.print_xy(inst1_x, 250, inst1, COLOR_GRAY)
  local inst2 = "SELECT/A: Start Race  START: Quick Start"
  local inst2_x = 240 - string.len(inst2) * 3
  rf.print_xy(inst2_x, 235, inst2, COLOR_GRAY)
end

function _EXIT()
end

function _DONE()
  -- Shutdown (unused)
end

