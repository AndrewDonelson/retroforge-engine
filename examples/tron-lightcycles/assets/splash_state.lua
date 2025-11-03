-- Game splash screen - shows Light Bikes splash image

local splash_time = 0
local splash_duration = 3.0 -- 3 seconds

-- State module functions (required by module system)

function _INIT()
  -- Initialization (if needed)
end

function _ENTER()
  splash_time = 0
  -- Start menu music
  rf.music("menu_music")
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
  
  -- Any button (0-10) skips splash
  -- Button 0: SELECT
  if rf.btnp(0) then
    if game then
      game.changeState("menu")
    end
    return
  end
  
  -- Button 1: START
  if rf.btnp(1) then
    if game then
      game.changeState("menu")
    end
    return
  end
  
  -- Button 2: UP
  if rf.btnp(2) then
    if game then
      game.changeState("menu")
    end
    return
  end
  
  -- Button 3: DOWN
  if rf.btnp(3) then
    if game then
      game.changeState("menu")
    end
    return
  end
  
  -- Button 4: LEFT
  if rf.btnp(4) then
    if game then
      game.changeState("menu")
    end
    return
  end
  
  -- Button 5: RIGHT
  if rf.btnp(5) then
    if game then
      game.changeState("menu")
    end
    return
  end
  
  -- Button 6: A
  if rf.btnp(6) then
    if game then
      game.changeState("menu")
    end
    return
  end
  
  -- Button 7: B
  if rf.btnp(7) then
    if game then
      game.changeState("menu")
    end
    return
  end
  
  -- Button 8: X
  if rf.btnp(8) then
    if game then
      game.changeState("menu")
    end
    return
  end
  
  -- Button 9: Y
  if rf.btnp(9) then
    if game then
      game.changeState("menu")
    end
    return
  end
  
  -- Button 10: TURBO
  if rf.btnp(10) then
    if game then
      game.changeState("menu")
    end
    return
  end
end

function _UPDATE(dt)
  splash_time = splash_time + dt
  
  -- Auto-transition after duration
  if splash_time >= splash_duration then
    -- Transition to menu
    if game then
      game.changeState("menu")
    end
  end
end

function _DRAW()
  -- Clear to charcoal gray background (very dark gray-brown, index 43 in Super Mario 50 palette)
  -- Index 43: RGB[79,9,0] - dark red-brown that appears as charcoal gray
  rf.clear_i(43)
  
  -- Draw splash sprite as full background (480x270)
  local sprite = rf.sprite("splash")
  if sprite then
    rf.spr("splash", 0, 0)
  else
    -- Fallback: draw text if sprite not loaded
    rf.print_anchored("TRON", "topcenter", 48)
    local cycles_y = 70
    local cycles_x = 240 - string.len("LIGHT CYCLES") * 3
    rf.print_xy(cycles_x, cycles_y, "LIGHT CYCLES", 48)
  end
  
  -- Show "Press any key" message (pulsing)
  if math.floor(splash_time * 2) % 2 == 0 then
    local msg = "Press any key to continue"
    local msg_x = 240 - string.len(msg) * 3
    rf.print_xy(msg_x, 250, msg, 1)
  end
end

function _EXIT()
  -- Stop music before exiting splash
  rf.sfx("stopall")
end

function _DONE()
  -- Shutdown (unused)
end
