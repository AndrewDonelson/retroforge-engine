color_index = 1

function _INIT()
  rf.palette_set("default")
end

function _UPDATE(dt)
  color_index = color_index + dt * 2  -- slow cycle
  if color_index > 49 then color_index = 1 end
end

function _DRAW()
  rf.clear_i(0)
  rf.print_center("HELLO FROM RETROFORGE", 135, 255,255,255)
  rf.print_xy(2, 2, "SCORE:"..tostring(1234), math.floor(color_index))
end


