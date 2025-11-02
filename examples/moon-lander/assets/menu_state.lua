-- Menu State Module for Moon Lander

local selected_option = 1
local menu_options = {
  "Play",
  "Exit"
}

function _INIT()
  -- Module initialization (runs once)
  -- Set up credits here when game object is guaranteed to be available
  game.addCredit("Game", "Moon Lander", "Lunar landing game")
  game.addCredit("Gameplay", "Procedural Levels", "50 procedurally generated levels")
  game.addCredit("Gameplay", "Physics", "Gravity and thrust simulation")
  game.addCredit("Engine", "Module System", "rf.import()")
  game.addCredit("Engine", "State Machine", "game.* API")
end

function _ENTER()
  -- Reset menu when entering
  selected_option = 1
end

function _HANDLE_INPUT()
  -- Navigate menu (buttons: 0=Left, 1=Right, 2=Up, 3=Down, 4=Z, 5=X)
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
      -- Start game - reset score and level
      score = 0
      level = 1
      set_level(level)
      rf.sfx("select")
      rf.music("start_melody")
      game.changeState("play")
    elseif selected_option == 2 then
      -- Exit - transition to credits (which will then exit)
      rf.sfx("select")
      game.exit()
    end
  end
end

function _UPDATE(dt)
  -- Menu doesn't need update logic
end

function _DRAW()
  -- Clear screen
  rf.clear_i(COLOR_BLACK)
  
  -- Draw title at topcenter
  rf.print_anchored("MOON LANDER", "topcenter", COLOR_WHITE)
  
  -- Draw menu options (moved up to avoid overlap)
  local menu_base_y = 100
  for i, option in ipairs(menu_options) do
    local y = menu_base_y + (i - 1) * 30
    local color = (i == selected_option) and COLOR_WHITE or COLOR_GRAY
    
    if i == selected_option then
      local selected_text = "> " .. option .. " <"
      local text_w = #selected_text * 6
      rf.print_xy((480 - text_w) / 2, y, selected_text, color)
    else
      local text_w = #option * 6
      rf.print_xy((480 - text_w) / 2, y, option, color)
    end
  end
  
  -- Draw best scores (between menu and controls)
  local best_text = "Best Level: "..tostring(best_level).."   Best Score: "..tostring(best_score)
  rf.print_xy(240 - string.len(best_text)*3, 165, best_text, COLOR_GRAY)
  
  -- Draw instructions (moved lower to avoid overlap)
  rf.print_xy(240 - string.len("Up/Down to select, O/X/Enter to confirm")*3, 200, "Up/Down to select, O/X/Enter to confirm", COLOR_GRAY)
  rf.print_xy(240 - string.len("Controls: Left/Right Rotate, Up Thrust")*3, 216, "Controls: Left/Right Rotate, Up Thrust", COLOR_GRAY)
end

function _EXIT()
  -- Cleanup when leaving menu
end

function _DONE()
  -- Shutdown (unused for this state)
end

