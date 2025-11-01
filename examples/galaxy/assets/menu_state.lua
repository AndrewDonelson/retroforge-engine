-- Menu State Module
-- Shows a menu with options

local selected_option = 1
local menu_options = {
  "Start Game",
  "Exit"
}

function _INIT()
  -- Module initialization (runs once)
  -- Set up credits here when game object is guaranteed to be available
        game.addCredit("Game", "Galaxy Simulation", "Physics Demo")
  game.addCredit("Physics", "Central Gravity", "Point gravity simulation")
  game.addCredit("Physics", "Orbital Mechanics", "Spiral galaxy formation")
  game.addCredit("Graphics", "Sprite Pooling", "2048 stars (automatic)")
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
  elseif rf.btnp(3) then  -- Down
    selected_option = selected_option + 1
    if selected_option > #menu_options then
      selected_option = 1
    end
  end
  
  -- Select option
  if rf.btnp(4) or rf.btnp(5) then  -- Z or X
    if selected_option == 1 then
      -- Start game
      game.changeState("play")
    elseif selected_option == 2 then
      -- Exit - transition to credits (which will then exit)
      game.exit()
    end
  end
end

function _UPDATE(dt)
  -- Menu doesn't need update logic
end

function _DRAW()
  -- Clear screen
  rf.clear_i(0)
  
  -- Draw title at topcenter
  rf.print_anchored("GALAXY SIMULATION", "topcenter", 15)
  
  -- Draw subtitle (centered manually with offset from top)
  local subtitle_y = 20
  local subtitle_text = "Spiral Galaxy Formation"
  local subtitle_w = #subtitle_text * 6
  rf.print_xy((480 - subtitle_w) / 2, subtitle_y, subtitle_text, 7)
  
  -- Draw menu options - use middlecenter for base, then offset manually
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
  
  -- Draw instructions at bottomcenter
  rf.print_anchored("Arrow Keys: Navigate", "bottomcenter", 6)
  
  -- Draw second instruction line (manually offset up from bottom)
  local inst2_y = 250
  local inst2_w = #"Z/X: Select" * 6
  rf.print_xy((480 - inst2_w) / 2, inst2_y, "Z/X: Select", 6)
end

function _EXIT()
  -- Cleanup when leaving menu
end

function _DONE()
  -- Shutdown (unused for this state)
end

