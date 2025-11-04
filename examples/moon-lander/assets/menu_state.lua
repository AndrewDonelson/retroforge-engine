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
      -- Start game - reset score and level
      rf.sfx("stopall") -- Stop all audio including music
      score = 0
      level = 1
      set_level(level)
      rf.sfx("select")
      rf.music("start_melody")
      game.changeState("play")
    elseif selected_option == 2 then
      -- Exit - transition to credits (which will then exit)
      rf.sfx("select")
      rf.sfx("stopall") -- Stop all audio before exit
      game.exit()
    end
  end
  
  -- Button 1: START
  if rf.btnp(1) then
    -- START bypasses menu and goes straight to PLAY
    rf.sfx("stopall") -- Stop all audio including music
    score = 0
    level = 1
    set_level(level)
    rf.sfx("select")
    rf.music("start_melody")
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
  
  -- Button 4: LEFT
  -- Not used in menu
  
  -- Button 5: RIGHT
  -- Not used in menu
  
  -- Button 6: A
  if rf.btnp(6) then
    -- A button also confirms menu selection (same as SELECT)
    if selected_option == 1 then
      -- Start game - reset score and level
      rf.sfx("stopall") -- Stop all audio including music
      score = 0
      level = 1
      set_level(level)
      rf.sfx("select")
      rf.music("start_melody")
      game.changeState("play")
    elseif selected_option == 2 then
      -- Exit - transition to credits (which will then exit)
      rf.sfx("select")
      rf.sfx("stopall") -- Stop all audio before exit
      game.exit()
    end
  end
  
  -- Button 7: B
  -- Not used in menu
  
  -- Button 8: X
  -- Not used in menu
  
  -- Button 9: Y
  -- Not used in menu
  
  -- Button 10: TURBO
  -- Not used in menu
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
  -- Stop music when leaving menu (if transitioning to play state)
  rf.sfx("stopall")
end

function _DONE()
  -- Shutdown (unused for this state)
end

