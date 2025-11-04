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
  game.addCredit("Game", "Animated Character Demo", "Multi-Frame Sprites & Animations")
  game.addCredit("Engine Features", "Module System", "rf.import()")
  game.addCredit("Engine Features", "State Machine", "game.* API")
  game.addCredit("Engine Features", "Sprite Pooling", "Automatic (transparent)")
  game.addCredit("Engine Features", "Physics", "Box2D Integration")
  game.addCredit("Engine Features", "Stats Display", "FPS, Memory, Objects")
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
      -- Start game
      rf.sfx("stopall") -- Stop all audio including music
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
      -- Start game
      rf.sfx("stopall") -- Stop all audio including music
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
  rf.clear_i(0)
  
  -- Draw title at topcenter
  rf.print_anchored("ANIMATED CHARACTER DEMO", "topcenter", 15)
  
  -- Draw subtitle (centered manually with offset from top)
  local subtitle_y = 20
  local subtitle_w = #"Multi-Frame Sprites & Animations" * 6
  rf.print_xy((480 - subtitle_w) / 2, subtitle_y, "Multi-Frame Sprites & Animations", 7)
  
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
  -- Stop music when leaving menu (if transitioning to play state)
  rf.sfx("stopall")
end

function _DONE()
  -- Shutdown (unused for this state)
end

