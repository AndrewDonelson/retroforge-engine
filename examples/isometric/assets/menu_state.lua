-- Menu State Module for Isometric Tilemap Demo

local selected_option = 1
local menu_options = {
  "Start Game",
  "Level Select",
  "Exit"
}

function _INIT()
  -- Module initialization (runs once)
  game.addCredit("Game", "Isometric Tilemap Demo", "Isometric Tilesets & Tilemaps")
  game.addCredit("Engine Features", "Module System", "rf.import()")
  game.addCredit("Engine Features", "State Machine", "game.* API")
  game.addCredit("Engine Features", "Tilesets", "Static, Frames, Animation")
  game.addCredit("Engine Features", "Tilemaps", "Grid-based maps")
  game.addCredit("Engine Features", "Isometric Rendering", "2.5D tiles")
end

function _ENTER()
  selected_option = 1
end

function _HANDLE_INPUT()
  -- Button 0: SELECT
  if rf.btnp(0) then
    if selected_option == 1 then
      rf.sfx("stopall")
      game.changeState("play")
    elseif selected_option == 2 then
      rf.sfx("stopall")
      game.changeState("level_select")
    elseif selected_option == 3 then
      rf.sfx("select")
      rf.sfx("stopall")
      game.exit()
    end
  end

  -- Button 1: START
  if rf.btnp(1) then
    rf.sfx("stopall")
    game.changeState("play")
    return
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

  -- Button 6: A
  if rf.btnp(6) then
    if selected_option == 1 then
      rf.sfx("stopall")
      game.changeState("play")
    elseif selected_option == 2 then
      rf.sfx("stopall")
      game.changeState("level_select")
    elseif selected_option == 3 then
      rf.sfx("select")
      rf.sfx("stopall")
      game.exit()
    end
  end
end

function _UPDATE(dt)
end

function _DRAW()
  rf.clear_i(0)

  -- Draw title at topcenter
  rf.print_anchored("ISOMETRIC TILEMAP DEMO", "topcenter", 15)

  -- Draw subtitle
  local subtitle_y = 20
  local subtitle_w = #"Isometric Tilesets & Tilemaps" * 6
  rf.print_xy((480 - subtitle_w) / 2, subtitle_y, "Isometric Tilesets & Tilemaps", 7)

  -- Draw menu options
  local menu_base_y = 120
  for i, option in ipairs(menu_options) do
    local y = menu_base_y + (i - 1) * 30
    local color = 7

    if i == selected_option then
      color = 15
      local selected_text = "> " .. option .. " <"
      local text_w = #selected_text * 6
      rf.print_xy((480 - text_w) / 2, y, selected_text, color)
    else
      local text_w = #option * 6
      rf.print_xy((480 - text_w) / 2, y, option, color)
    end
  end

  -- Draw instructions
  rf.print_anchored("Arrow Keys: Navigate", "bottomcenter", 6)
  local inst2_y = 250
  local inst2_w = #"Z/X: Select" * 6
  rf.print_xy((480 - inst2_w) / 2, inst2_y, "Z/X: Select", 6)
end

function _EXIT()
  rf.sfx("stopall")
end

function _DONE()
end

