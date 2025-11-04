-- Splash State - Shows game title and description

local splash_time = 0
local splash_duration = 3.0 -- 3 seconds

function _INIT()
  -- Initialization
end

function _ENTER()
  splash_time = 0
end

function _HANDLE_INPUT()
  -- Any button (0-10) skips splash and goes to menu
  for i = 0, 10 do
    if rf.btnp(i) then
      rf.sfx("stopall") -- Stop all audio before transition
      if game then
        game.changeState("menu")
      end
      return
    end
  end
end

function _UPDATE(dt)
  splash_time = splash_time + dt
  
  -- Auto-transition to menu after duration
  if splash_time >= splash_duration then
    if game then
      game.changeState("menu")
    end
  end
end

function _DRAW()
  -- Clear screen with light background
  rf.cls(14)
  
  -- Title at top center
  rf.print_anchored("Animated Character Demo", "topcenter", 15)
  
  -- Description in middle center
  rf.print_anchored("A cute demonstration of multi-frame sprites and animations. Features a walking character with directional states and animated walk cycles.", "middlecenter", 0)
  
  -- Instructions at bottom center
  rf.print_anchored("Press any key to continue", "bottomcenter", 6)
end

function _EXIT()
  -- Cleanup
end

function _DONE()
  -- Shutdown
end

