-- Level Select State - Choose between Normal and Isometric tilemap

local selected_level = 1
local levels = {
  {
    name = "Normal",
    description = "Orthogonal tilemap",
    tilesetName = "terrain_normal",
    mapName = "test_normal"
  },
  {
    name = "Isometric",
    description = "Isometric tilemap (2.5D)",
    tilesetName = "terrain",
    mapName = "test"
  }
}

function _INIT()
end

function _ENTER()
  selected_level = 1
end

function _HANDLE_INPUT()
  -- Button 0: SELECT or A (6)
  if rf.btnp(0) or rf.btnp(6) then
    rf.sfx("select")
    rf.sfx("stopall")
    -- Store selected level index globally for play_state to access
    _G.selected_level_index = selected_level
    game.changeState("play")
    return
  end
  
  -- Button 1: START
  if rf.btnp(1) then
    rf.sfx("select")
    rf.sfx("stopall")
    -- Store selected level index globally for play_state to access
    _G.selected_level_index = selected_level
    game.changeState("play")
    return
  end
  
  -- Button 2: UP
  if rf.btnp(2) then
    rf.sfx("select")
    selected_level = selected_level - 1
    if selected_level < 1 then
      selected_level = #levels
    end
    return
  end
  
  -- Button 3: DOWN
  if rf.btnp(3) then
    rf.sfx("select")
    selected_level = selected_level + 1
    if selected_level > #levels then
      selected_level = 1
    end
    return
  end
end

function _UPDATE(dt)
end

function _DRAW()
  rf.clear_i(0)
  
  -- Draw title
  rf.print_anchored("SELECT TILEMAP VIEW", "topcenter", 15)
  
  -- Draw level options
  local start_y = 80
  for i = 1, #levels do
    local y = start_y + (i - 1) * 60
    local level = levels[i]
    local color = 7
    local desc_color = 6
    
    if i == selected_level then
      color = 15
      desc_color = 14
      local selected_text = "> " .. level.name .. " <"
      local text_w = #selected_text * 6
      rf.print_xy((480 - text_w) / 2, y, selected_text, color)
    else
      local text_w = #level.name * 6
      rf.print_xy((480 - text_w) / 2, y, level.name, color)
    end
    
    -- Draw description
    local desc_y = y + 15
    local desc_w = #level.description * 6
    rf.print_xy((480 - desc_w) / 2, desc_y, level.description, desc_color)
  end
  
  -- Instructions
  rf.print_anchored("Arrow Keys: Navigate", "bottomcenter", 6)
  local inst2_y = 250
  local inst2_w = #"SELECT/A: Choose" * 6
  rf.print_xy((480 - inst2_w) / 2, inst2_y, "SELECT/A: Choose", 6)
end

function _EXIT()
end

function _DONE()
end

