package luabind

import (
	"image"

	"github.com/AndrewDonelson/retroforge-engine/internal/imgtool"
	lua "github.com/yuin/gopher-lua"
)

// RegisterImgToolAPI registers imgtool functions with Lua state
func RegisterImgToolAPI(L *lua.LState) {
	// Get or create rf table
	rf := L.GetGlobal("rf")
	var rfTable *lua.LTable
	if rf == lua.LNil {
		rfTable = L.NewTable()
		L.SetGlobal("rf", rfTable)
	} else {
		rfTable = rf.(*lua.LTable)
	}

	// Create imgtool table
	imgToolTable := L.NewTable()

	// Register functions
	L.SetField(imgToolTable, "quantize", L.NewFunction(luaQuantize))
	L.SetField(imgToolTable, "map_palette", L.NewFunction(luaMapPalette))
	L.SetField(imgToolTable, "scale", L.NewFunction(luaScale))
	L.SetField(imgToolTable, "to_sprite", L.NewFunction(luaToSprite))
	L.SetField(imgToolTable, "load_png", L.NewFunction(luaLoadPNG))
	L.SetField(imgToolTable, "load_palette", L.NewFunction(luaLoadPalette))

	// Register under rf.imgtool
	L.SetField(rfTable, "imgtool", imgToolTable)
}

// luaQuantize: rf.imgtool.quantize(png_data, options)
func luaQuantize(L *lua.LState) int {
	// Get PNG data (userdata or string)
	pngData := L.CheckAny(1)

	// Parse options table (optional)
	opts := imgtool.DefaultQuantizeOptions()
	if L.GetTop() >= 2 {
		optsTable := L.CheckTable(2)
		opts = parseQuantizeOptions(L, optsTable)
	}

	// Load image
	img, err := loadImageFromLua(L, pngData)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	// Quantize
	palette, err := imgtool.Quantize(img, opts)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	// Convert to Lua table (array of hex strings)
	paletteTable := paletteToLuaTable(L, palette)
	L.Push(paletteTable)
	return 1
}

// luaMapPalette: rf.imgtool.map_palette(png_data, palette_json, options)
func luaMapPalette(L *lua.LState) int {
	// Get PNG data
	pngData := L.CheckAny(1)

	// Get palette (table or JSON string)
	paletteData := L.CheckAny(2)

	// Parse options
	opts := imgtool.DefaultMapPaletteOptions()
	if L.GetTop() >= 3 {
		optsTable := L.CheckTable(3)
		opts = parseMapPaletteOptions(L, optsTable)
	}

	// Load image and palette
	img, err := loadImageFromLua(L, pngData)
	if err != nil {
		return luaError(L, err)
	}

	palette, err := loadPaletteFromLua(L, paletteData)
	if err != nil {
		return luaError(L, err)
	}

	// Map palette
	indices, err := imgtool.MapPalette(img, palette, opts)
	if err != nil {
		return luaError(L, err)
	}

	// Convert to Lua 2D table
	indicesTable := indices2DToLuaTable(L, indices)
	L.Push(indicesTable)
	return 1
}

// luaScale: rf.imgtool.scale(png_data, options)
func luaScale(L *lua.LState) int {
	pngData := L.CheckAny(1)

	opts := imgtool.DefaultScaleOptions()
	if L.GetTop() >= 2 {
		optsTable := L.CheckTable(2)
		opts = parseScaleOptions(L, optsTable)
	}

	img, err := loadImageFromLua(L, pngData)
	if err != nil {
		return luaError(L, err)
	}

	rgbData, err := imgtool.Scale(img, opts)
	if err != nil {
		return luaError(L, err)
	}

	// Convert to Lua 3D table [row][col][R,G,B,A]
	rgbTable := rgb3DToLuaTable(L, rgbData)
	L.Push(rgbTable)
	return 1
}

// luaToSprite: rf.imgtool.to_sprite(png_data, palette_json, options)
func luaToSprite(L *lua.LState) int {
	pngData := L.CheckAny(1)
	paletteData := L.CheckAny(2)
	optsTable := L.CheckTable(3)

	opts := parseToSpriteOptions(L, optsTable)

	img, err := loadImageFromLua(L, pngData)
	if err != nil {
		return luaError(L, err)
	}

	palette, err := loadPaletteFromLua(L, paletteData)
	if err != nil {
		return luaError(L, err)
	}

	sprite, err := imgtool.ToSprite(img, palette, opts)
	if err != nil {
		return luaError(L, err)
	}

	// Convert sprite to Lua table
	spriteTable := spriteToLuaTable(L, sprite)
	L.Push(spriteTable)
	return 1
}

// luaLoadPNG: rf.imgtool.load_png(filename)
func luaLoadPNG(L *lua.LState) int {
	_ = L.CheckString(1) // filename - reserved for future use
	// In a real implementation, this would load from cart assets
	// For now, return a placeholder
	L.Push(lua.LNil)
	L.Push(lua.LString("load_png not yet implemented - requires asset loader"))
	return 2
}

// luaLoadPalette: rf.imgtool.load_palette(filename)
func luaLoadPalette(L *lua.LState) int {
	_ = L.CheckString(1) // filename - reserved for future use
	// In a real implementation, this would load from cart assets
	// For now, return a placeholder
	L.Push(lua.LNil)
	L.Push(lua.LString("load_palette not yet implemented - requires asset loader"))
	return 2
}

// Helper functions for type conversion

func loadImageFromLua(L *lua.LState, data lua.LValue) (image.Image, error) {
	switch v := data.(type) {
	case *lua.LUserData:
		// Assume it's image.Image wrapped in userdata
		if img, ok := v.Value.(image.Image); ok {
			return img, nil
		}
		return nil, &imgtool.ImgToolError{
			Code:    imgtool.ErrCodeInvalidImage,
			Message: "invalid image userdata",
		}
	case lua.LString:
		// Assume it's a file path
		img, err := imgtool.LoadPNGFile(string(v))
		if err != nil {
			return nil, err
		}
		return img, nil
	default:
		return nil, &imgtool.ImgToolError{
			Code:    imgtool.ErrCodeInvalidImage,
			Message: "invalid image data type",
		}
	}
}

func loadPaletteFromLua(L *lua.LState, data lua.LValue) (*imgtool.Palette, error) {
	switch v := data.(type) {
	case *lua.LTable:
		// Convert Lua table to Palette
		return tableToPalette(L, v)
	case lua.LString:
		// Assume it's JSON string or file path
		palette, err := imgtool.LoadPaletteFile(string(v))
		if err != nil {
			// Try as JSON
			palette, err = imgtool.LoadPalette([]byte(string(v)))
		}
		return palette, err
	default:
		return nil, &imgtool.ImgToolError{
			Code:    imgtool.ErrCodeInvalidPalette,
			Message: "invalid palette data type",
		}
	}
}

func paletteToLuaTable(L *lua.LState, palette *imgtool.Palette) *lua.LTable {
	tbl := L.NewTable()
	for i, color := range palette.Colors {
		L.RawSetInt(tbl, i+1, lua.LString(color))
	}
	return tbl
}

func tableToPalette(L *lua.LState, tbl *lua.LTable) (*imgtool.Palette, error) {
	palette := &imgtool.Palette{Colors: make([]string, 0, 48)} // Game palette is 48 colors
	tbl.ForEach(func(key lua.LValue, value lua.LValue) {
		if str, ok := value.(lua.LString); ok {
			palette.Colors = append(palette.Colors, string(str))
		}
	})
	if err := palette.Validate(); err != nil {
		return nil, err
	}
	return palette, nil
}

func indices2DToLuaTable(L *lua.LState, indices [][]int) *lua.LTable {
	tbl := L.NewTable()
	for i, row := range indices {
		rowTbl := L.NewTable()
		for j, idx := range row {
			L.RawSetInt(rowTbl, j+1, lua.LNumber(idx))
		}
		L.RawSetInt(tbl, i+1, rowTbl)
	}
	return tbl
}

func rgb3DToLuaTable(L *lua.LState, rgb [][][]uint8) *lua.LTable {
	tbl := L.NewTable()
	for i, row := range rgb {
		rowTbl := L.NewTable()
		for j, pixel := range row {
			pixelTbl := L.NewTable()
			for k, component := range pixel {
				L.RawSetInt(pixelTbl, k+1, lua.LNumber(component))
			}
			L.RawSetInt(rowTbl, j+1, pixelTbl)
		}
		L.RawSetInt(tbl, i+1, rowTbl)
	}
	return tbl
}

func spriteToLuaTable(L *lua.LState, sprite *imgtool.Sprite) *lua.LTable {
	tbl := L.NewTable()
	L.SetField(tbl, "width", lua.LNumber(sprite.Width))
	L.SetField(tbl, "height", lua.LNumber(sprite.Height))
	L.SetField(tbl, "pixels", indices2DToLuaTable(L, sprite.Pixels))
	L.SetField(tbl, "useCollision", lua.LBool(sprite.UseCollision))
	L.SetField(tbl, "isUI", lua.LBool(sprite.IsUI))
	L.SetField(tbl, "lifetime", lua.LNumber(sprite.Lifetime))
	L.SetField(tbl, "maxSpawn", lua.LNumber(sprite.MaxSpawn))
	mountPointsTbl := L.NewTable()
	for i, point := range sprite.MountPoints {
		pointTbl := L.NewTable()
		L.SetField(pointTbl, "x", lua.LNumber(point.X))
		L.SetField(pointTbl, "y", lua.LNumber(point.Y))
		L.RawSetInt(mountPointsTbl, i+1, pointTbl)
	}
	L.SetField(tbl, "mountPoints", mountPointsTbl)
	return tbl
}

func parseQuantizeOptions(L *lua.LState, tbl *lua.LTable) imgtool.QuantizeOptions {
	opts := imgtool.DefaultQuantizeOptions()
	tbl.ForEach(func(key lua.LValue, value lua.LValue) {
		keyStr := key.String()
		switch keyStr {
		case "dither":
			if str, ok := value.(lua.LString); ok {
				opts.DitherAlgorithm = string(str)
			}
		case "alpha_threshold":
			if num, ok := value.(lua.LNumber); ok {
				opts.AlphaThreshold = uint8(num)
			}
		case "enforce_black_white":
			if b, ok := value.(lua.LBool); ok {
				opts.EnforceBlackWhite = bool(b)
			}
		}
	})
	return opts
}

func parseMapPaletteOptions(L *lua.LState, tbl *lua.LTable) imgtool.MapPaletteOptions {
	opts := imgtool.DefaultMapPaletteOptions()
	tbl.ForEach(func(key lua.LValue, value lua.LValue) {
		keyStr := key.String()
		switch keyStr {
		case "dither":
			if str, ok := value.(lua.LString); ok {
				opts.DitherAlgorithm = string(str)
			}
		case "alpha_threshold":
			if num, ok := value.(lua.LNumber); ok {
				opts.AlphaThreshold = uint8(num)
			}
		case "transparent_index":
			if num, ok := value.(lua.LNumber); ok {
				opts.TransparentIndex = int(num)
			}
		}
	})
	return opts
}

func parseScaleOptions(L *lua.LState, tbl *lua.LTable) imgtool.ScaleOptions {
	opts := imgtool.DefaultScaleOptions()
	tbl.ForEach(func(key lua.LValue, value lua.LValue) {
		keyStr := key.String()
		switch keyStr {
		case "width":
			if num, ok := value.(lua.LNumber); ok {
				opts.Width = int(num)
			}
		case "height":
			if num, ok := value.(lua.LNumber); ok {
				opts.Height = int(num)
			}
		case "algorithm":
			if str, ok := value.(lua.LString); ok {
				opts.Algorithm = string(str)
			}
		case "ensure_divisible":
			if b, ok := value.(lua.LBool); ok {
				opts.EnsureDivisible = bool(b)
			}
		case "preserve_aspect":
			if b, ok := value.(lua.LBool); ok {
				opts.PreserveAspect = bool(b)
			}
		}
	})
	return opts
}

func parseToSpriteOptions(L *lua.LState, tbl *lua.LTable) imgtool.ToSpriteOptions {
	opts := imgtool.DefaultToSpriteOptions("sprite")
	tbl.ForEach(func(key lua.LValue, value lua.LValue) {
		keyStr := key.String()
		switch keyStr {
		case "name":
			if str, ok := value.(lua.LString); ok {
				opts.Name = string(str)
			}
		case "width", "target_width":
			if num, ok := value.(lua.LNumber); ok {
				opts.TargetWidth = int(num)
			}
		case "height", "target_height":
			if num, ok := value.(lua.LNumber); ok {
				opts.TargetHeight = int(num)
			}
		case "use_collision":
			if b, ok := value.(lua.LBool); ok {
				opts.UseCollision = bool(b)
			}
		case "is_ui":
			if b, ok := value.(lua.LBool); ok {
				opts.IsUI = bool(b)
			}
		case "lifetime":
			if num, ok := value.(lua.LNumber); ok {
				opts.Lifetime = int(num)
			}
		case "max_spawn":
			if num, ok := value.(lua.LNumber); ok {
				opts.MaxSpawn = int(num)
			}
		case "dither":
			if str, ok := value.(lua.LString); ok {
				opts.DitherAlgorithm = string(str)
			}
		case "alpha_threshold":
			if num, ok := value.(lua.LNumber); ok {
				opts.AlphaThreshold = uint8(num)
			}
		}
	})
	return opts
}

func luaError(L *lua.LState, err error) int {
	L.Push(lua.LNil)
	L.Push(lua.LString(err.Error()))
	return 2
}
