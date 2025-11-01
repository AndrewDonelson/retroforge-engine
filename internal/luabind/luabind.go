package luabind

import (
	"fmt"
	"image/color"
	"time"

	"github.com/AndrewDonelson/retroforge-engine/internal/app"
	"github.com/AndrewDonelson/retroforge-engine/internal/audio"
	"github.com/AndrewDonelson/retroforge-engine/internal/cartio"
	"github.com/AndrewDonelson/retroforge-engine/internal/font"
	"github.com/AndrewDonelson/retroforge-engine/internal/graphics"
	"github.com/AndrewDonelson/retroforge-engine/internal/input"
	"github.com/AndrewDonelson/retroforge-engine/internal/network"
	"github.com/AndrewDonelson/retroforge-engine/internal/physics"
	"github.com/AndrewDonelson/retroforge-engine/internal/spritepool"
	lua "github.com/yuin/gopher-lua"
)

// small POD color to avoid importing image/color in caller signature
type ColorByIndex func(i int) (rgba [4]uint8)

// DevModeHandler interface for development mode operations (avoids import cycle)
type DevModeHandler interface {
	IsEnabled() bool
	AddDebugLog(msg string)
	GetStats() interface{} // Returns DevStats-compatible structure
}

// DevStats holds debugging statistics (duplicated here to avoid import cycle)
type DevStats struct {
	FPS         float64
	FrameCount  int64
	LuaMemory   int64
	LoadTime    time.Duration
	LastReload  time.Time
	ReloadCount int
}

// Register attaches rf.* drawing functions to the Lua state.
func Register(L *lua.LState, r graphics.Renderer, colorByIndex ColorByIndex, setPalette func(string), sfxMap cartio.SFXMap, musicMap cartio.MusicMap, spritesMap cartio.SpriteMap, physWorld *physics.World, netMgr *network.NetworkManager) {
	state := NewState()
	RegisterWithState(L, r, colorByIndex, setPalette, sfxMap, musicMap, spritesMap, physWorld, state, netMgr)
}

// RegisterWithDev attaches rf.* drawing functions with dev mode support
func RegisterWithDev(L *lua.LState, r graphics.Renderer, colorByIndex ColorByIndex, setPalette func(string), sfxMap cartio.SFXMap, musicMap cartio.MusicMap, spritesMap cartio.SpriteMap, physWorld *physics.World, devMode DevModeHandler, netMgr *network.NetworkManager) {
	state := NewState()
	RegisterWithDevMode(L, r, colorByIndex, setPalette, sfxMap, musicMap, spritesMap, physWorld, state, devMode, netMgr)
}

// RegisterWithState attaches rf.* drawing functions with state management
func RegisterWithState(L *lua.LState, r graphics.Renderer, colorByIndex ColorByIndex, setPalette func(string), sfxMap cartio.SFXMap, musicMap cartio.MusicMap, spritesMap cartio.SpriteMap, physWorld *physics.World, state *State, netMgr *network.NetworkManager) {
	RegisterWithDevMode(L, r, colorByIndex, setPalette, sfxMap, musicMap, spritesMap, physWorld, state, nil, netMgr)
}

// RegisterWithDevMode attaches rf.* drawing functions with dev mode support
func RegisterWithDevMode(L *lua.LState, r graphics.Renderer, colorByIndex ColorByIndex, setPalette func(string), sfxMap cartio.SFXMap, musicMap cartio.MusicMap, spritesMap cartio.SpriteMap, physWorld *physics.World, state *State, devMode DevModeHandler, netMgr *network.NetworkManager) {
	rf := L.NewTable()
	L.SetGlobal("rf", rf)

	// Store devMode in closure for debug functions
	devModePtr := devMode

	// Create pool manager for automatic sprite pooling
	poolManager := spritepool.NewPoolManager()

	// Register existing sprites that meet pooling criteria
	for spriteName, spriteData := range spritesMap {
		if spritepool.ShouldPool(spriteData) {
			if err := poolManager.RegisterSprite(spriteName, spriteData); err != nil {
				// Log error but don't fail - pooling is optional
				if devMode != nil {
					devMode.AddDebugLog(fmt.Sprintf("Failed to create pool for sprite '%s': %v", spriteName, err))
				}
			}
		}
	}

	// Create wrapper for colorByIndex that applies remapping
	colorByIndexRemapped := func(i int) (c [4]uint8) {
		remapped := state.GetPalRemap(i)
		return colorByIndex(remapped)
	}

	// rf.print_anchored(text, anchor, index)
	// anchor: "topleft", "topcenter", "topright", "middleleft", "middlecenter", "middleright", "bottomleft", "bottomcenter", "bottomright"
	L.SetField(rf, "print_anchored", L.NewFunction(func(L *lua.LState) int {
		txt := L.CheckString(1)
		anchor := L.CheckString(2)
		idx := L.CheckInt(3)
		c := colorByIndexRemapped(idx)
		r.PrintAnchored(txt, anchor, color.RGBA{R: c[0], G: c[1], B: c[2], A: c[3]})
		return 0
	}))

	// rf.clear_i(idx)
	L.SetField(rf, "clear_i", L.NewFunction(func(L *lua.LState) int {
		idx := L.CheckInt(1)
		c := colorByIndexRemapped(idx)
		r.Clear(color.RGBA{R: c[0], G: c[1], B: c[2], A: c[3]})
		return 0
	}))

	// rf.print_xy(x,y,text, [idx]) - If idx omitted, use cursor/color state
	L.SetField(rf, "print_xy", L.NewFunction(func(L *lua.LState) int {
		x := L.CheckInt(1)
		y := L.CheckInt(2)
		txt := L.CheckString(3)

		var idx int
		if L.GetTop() >= 4 {
			idx = L.CheckInt(4)
		} else {
			// Use cursor/color state if available
			if colorIdx, hasColor := state.GetTextColor(); hasColor {
				idx = colorIdx
			} else {
				idx = 15 // Default white
			}
			if cx, cy, hasCursor := state.GetCursor(); hasCursor {
				x = cx
				y = cy
			}
		}

		c := colorByIndexRemapped(idx)
		r.Print(txt, x, y, color.RGBA{R: c[0], G: c[1], B: c[2], A: c[3]})

		// Update cursor position after printing (handle newlines)
		// Note: This matches PICO-8 behavior
		finalX := x
		finalY := y
		runes := []rune(txt)
		for _, r := range runes {
			if r == '\n' {
				finalX = x                // Return to start of line
				finalY += font.Height + 1 // Advance to next line
			} else {
				finalX += font.Advance
			}
		}
		state.SetCursor(finalX, finalY)
		return 0
	}))

	// rf.cursor([x, y]) - Set text cursor position. No args resets cursor.
	L.SetField(rf, "cursor", L.NewFunction(func(L *lua.LState) int {
		if L.GetTop() == 0 {
			// No args = reset cursor
			state.ResetCursor()
		} else {
			x := L.CheckInt(1)
			y := L.CheckInt(2)
			state.SetCursor(x, y)
		}
		return 0
	}))

	// rf.color([index]) - Set text color. No args resets color.
	L.SetField(rf, "color", L.NewFunction(func(L *lua.LState) int {
		if L.GetTop() == 0 {
			// No args = reset color
			state.ResetColor()
		} else {
			idx := L.CheckInt(1)
			state.SetTextColor(idx)
		}
		return 0
	}))

	// rf.print(text, [x, y, index]) - PICO-8-like print with optional cursor/color state
	L.SetField(rf, "print", L.NewFunction(func(L *lua.LState) int {
		txt := L.CheckString(1)

		var x, y, idx int
		useState := false

		if L.GetTop() >= 3 {
			// Explicit x, y provided
			x = L.CheckInt(2)
			y = L.CheckInt(3)
			if L.GetTop() >= 4 {
				idx = L.CheckInt(4)
			} else {
				// Use color state if available
				if colorIdx, hasColor := state.GetTextColor(); hasColor {
					idx = colorIdx
				} else {
					idx = 15 // Default white
				}
			}
		} else {
			// Use cursor/color state
			useState = true
			if cx, cy, hasCursor := state.GetCursor(); hasCursor {
				x = cx
				y = cy
			} else {
				x = 0
				y = 0
			}
			if colorIdx, hasColor := state.GetTextColor(); hasColor {
				idx = colorIdx
			} else {
				idx = 15 // Default white
			}
		}

		c := colorByIndexRemapped(idx)
		r.Print(txt, x, y, color.RGBA{R: c[0], G: c[1], B: c[2], A: c[3]})

		// Update cursor position after printing if using state (handle newlines)
		if useState {
			finalX := x
			finalY := y
			runes := []rune(txt)
			for _, r := range runes {
				if r == '\n' {
					finalX = x                // Return to start of line
					finalY += font.Height + 1 // Advance to next line
				} else {
					finalX += font.Advance
				}
			}
			state.SetCursor(finalX, finalY)
		}
		return 0
	}))

	// palette.set(name)
	L.SetField(rf, "palette_set", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		if setPalette != nil {
			setPalette(name)
		}
		return 0
	}))

	// Drawing primitives (index-colored)
	L.SetField(rf, "pset", L.NewFunction(func(L *lua.LState) int {
		x := L.CheckInt(1)
		y := L.CheckInt(2)
		idx := L.CheckInt(3)
		c := colorByIndexRemapped(idx)
		r.PSet(x, y, color.RGBA{c[0], c[1], c[2], c[3]})
		return 0
	}))
	L.SetField(rf, "line", L.NewFunction(func(L *lua.LState) int {
		x0 := L.CheckInt(1)
		y0 := L.CheckInt(2)
		x1 := L.CheckInt(3)
		y1 := L.CheckInt(4)
		idx := L.CheckInt(5)
		c := colorByIndexRemapped(idx)
		r.Line(x0, y0, x1, y1, color.RGBA{c[0], c[1], c[2], c[3]})
		return 0
	}))
	L.SetField(rf, "rect", L.NewFunction(func(L *lua.LState) int {
		x0 := L.CheckInt(1)
		y0 := L.CheckInt(2)
		x1 := L.CheckInt(3)
		y1 := L.CheckInt(4)
		idx := L.CheckInt(5)
		c := colorByIndexRemapped(idx)
		r.Rect(x0, y0, x1, y1, color.RGBA{c[0], c[1], c[2], c[3]})
		return 0
	}))
	L.SetField(rf, "rectfill", L.NewFunction(func(L *lua.LState) int {
		x0 := L.CheckInt(1)
		y0 := L.CheckInt(2)
		x1 := L.CheckInt(3)
		y1 := L.CheckInt(4)
		idx := L.CheckInt(5)
		c := colorByIndexRemapped(idx)
		r.RectFill(x0, y0, x1, y1, color.RGBA{c[0], c[1], c[2], c[3]})
		return 0
	}))
	L.SetField(rf, "circ", L.NewFunction(func(L *lua.LState) int {
		x := L.CheckInt(1)
		y := L.CheckInt(2)
		rad := L.CheckInt(3)
		idx := L.CheckInt(4)
		c := colorByIndexRemapped(idx)
		r.Circ(x, y, rad, color.RGBA{c[0], c[1], c[2], c[3]})
		return 0
	}))
	L.SetField(rf, "circfill", L.NewFunction(func(L *lua.LState) int {
		x := L.CheckInt(1)
		y := L.CheckInt(2)
		rad := L.CheckInt(3)
		idx := L.CheckInt(4)
		c := colorByIndexRemapped(idx)
		r.CircFill(x, y, rad, color.RGBA{c[0], c[1], c[2], c[3]})
		return 0
	}))

	// Shape primitives
	L.SetField(rf, "triangle", L.NewFunction(func(L *lua.LState) int {
		x := L.CheckInt(1)
		y := L.CheckInt(2)
		radius := L.CheckInt(3)
		filled := L.OptBool(4, false)
		idx := L.CheckInt(5)
		c := colorByIndexRemapped(idx)
		r.Triangle(x, y, radius, filled, color.RGBA{c[0], c[1], c[2], c[3]})
		return 0
	}))
	L.SetField(rf, "diamond", L.NewFunction(func(L *lua.LState) int {
		x := L.CheckInt(1)
		y := L.CheckInt(2)
		radius := L.CheckInt(3)
		filled := L.OptBool(4, false)
		idx := L.CheckInt(5)
		c := colorByIndexRemapped(idx)
		r.Diamond(x, y, radius, filled, color.RGBA{c[0], c[1], c[2], c[3]})
		return 0
	}))
	L.SetField(rf, "square", L.NewFunction(func(L *lua.LState) int {
		x := L.CheckInt(1)
		y := L.CheckInt(2)
		radius := L.CheckInt(3)
		filled := L.OptBool(4, false)
		idx := L.CheckInt(5)
		c := colorByIndexRemapped(idx)
		r.Square(x, y, radius, filled, color.RGBA{c[0], c[1], c[2], c[3]})
		return 0
	}))
	L.SetField(rf, "pentagon", L.NewFunction(func(L *lua.LState) int {
		x := L.CheckInt(1)
		y := L.CheckInt(2)
		radius := L.CheckInt(3)
		filled := L.OptBool(4, false)
		idx := L.CheckInt(5)
		c := colorByIndexRemapped(idx)
		r.Pentagon(x, y, radius, filled, color.RGBA{c[0], c[1], c[2], c[3]})
		return 0
	}))
	L.SetField(rf, "hexagon", L.NewFunction(func(L *lua.LState) int {
		x := L.CheckInt(1)
		y := L.CheckInt(2)
		radius := L.CheckInt(3)
		filled := L.OptBool(4, false)
		idx := L.CheckInt(5)
		c := colorByIndexRemapped(idx)
		r.Hexagon(x, y, radius, filled, color.RGBA{c[0], c[1], c[2], c[3]})
		return 0
	}))
	L.SetField(rf, "star", L.NewFunction(func(L *lua.LState) int {
		x := L.CheckInt(1)
		y := L.CheckInt(2)
		radius := L.CheckInt(3)
		filled := L.OptBool(4, false)
		idx := L.CheckInt(5)
		c := colorByIndexRemapped(idx)
		r.Star(x, y, radius, filled, color.RGBA{c[0], c[1], c[2], c[3]})
		return 0
	}))

	// Pixel reading
	L.SetField(rf, "pget", L.NewFunction(func(L *lua.LState) int {
		x := L.CheckInt(1)
		y := L.CheckInt(2)
		col := r.PGet(x, y)
		// Find closest palette index (simplified - just return RGB components)
		// For now, return table with r, g, b, a
		tbl := L.NewTable()
		tbl.RawSetString("r", lua.LNumber(col.R))
		tbl.RawSetString("g", lua.LNumber(col.G))
		tbl.RawSetString("b", lua.LNumber(col.B))
		tbl.RawSetString("a", lua.LNumber(col.A))
		L.Push(tbl)
		return 1
	}))

	// Clipping
	L.SetField(rf, "clip", L.NewFunction(func(L *lua.LState) int {
		if L.GetTop() == 0 {
			// No args = disable clipping
			r.SetClip(0, 0, 0, 0)
		} else {
			x := L.CheckInt(1)
			y := L.CheckInt(2)
			w := L.CheckInt(3)
			h := L.CheckInt(4)
			r.SetClip(x, y, w, h)
		}
		return 0
	}))

	// Camera
	L.SetField(rf, "camera", L.NewFunction(func(L *lua.LState) int {
		if L.GetTop() == 0 {
			// No args = reset camera
			r.SetCamera(0, 0)
		} else {
			x := L.CheckInt(1)
			y := L.CheckInt(2)
			r.SetCamera(x, y)
		}
		return 0
	}))

	// Ellipse drawing
	L.SetField(rf, "elli", L.NewFunction(func(L *lua.LState) int {
		x := L.CheckInt(1)
		y := L.CheckInt(2)
		rx := L.CheckInt(3)
		ry := L.CheckInt(4)
		idx := L.CheckInt(5)
		c := colorByIndexRemapped(idx)
		r.Ellipse(x, y, rx, ry, color.RGBA{c[0], c[1], c[2], c[3]})
		return 0
	}))
	L.SetField(rf, "ellifill", L.NewFunction(func(L *lua.LState) int {
		x := L.CheckInt(1)
		y := L.CheckInt(2)
		rx := L.CheckInt(3)
		ry := L.CheckInt(4)
		idx := L.CheckInt(5)
		c := colorByIndexRemapped(idx)
		r.EllipseFill(x, y, rx, ry, color.RGBA{c[0], c[1], c[2], c[3]})
		return 0
	}))

	// Input
	L.SetField(rf, "btn", L.NewFunction(func(L *lua.LState) int {
		i := L.CheckInt(1)

		// Check if this is a multiplayer call: rf.btn(player_id, button)
		if L.GetTop() >= 2 {
			playerID := i
			buttonID := L.CheckInt(2)
			if netMgr != nil && netMgr.IsHost() {
				// Host can check other players' inputs
				L.Push(lua.LBool(netMgr.GetPlayerInput(playerID, buttonID)))
				return 1
			}
			// Non-host trying to check other player - not allowed
			L.Push(lua.LBool(false))
			return 1
		}

		// Normal single-player or local player input
		L.Push(lua.LBool(input.Btn(i)))
		return 1
	}))
	L.SetField(rf, "btnp", L.NewFunction(func(L *lua.LState) int {
		i := L.CheckInt(1)
		L.Push(lua.LBool(input.Btnp(i)))
		return 1
	}))
	L.SetField(rf, "btnr", L.NewFunction(func(L *lua.LState) int {
		_ = L.CheckInt(1) // Button ID - release detection not yet implemented
		// Button release detection (not yet implemented in input package)
		// For now, return false
		L.Push(lua.LBool(false))
		return 1
	}))

	// Multiplayer API
	if netMgr != nil {
		// rf.is_multiplayer() → boolean
		L.SetField(rf, "is_multiplayer", L.NewFunction(func(L *lua.LState) int {
			L.Push(lua.LBool(netMgr.IsMultiplayer()))
			return 1
		}))

		// rf.player_count() → number (1-6)
		L.SetField(rf, "player_count", L.NewFunction(func(L *lua.LState) int {
			L.Push(lua.LNumber(netMgr.PlayerCount()))
			return 1
		}))

		// rf.my_player_id() → number (1-6)
		L.SetField(rf, "my_player_id", L.NewFunction(func(L *lua.LState) int {
			L.Push(lua.LNumber(netMgr.PlayerID()))
			return 1
		}))

		// rf.is_host() → boolean
		L.SetField(rf, "is_host", L.NewFunction(func(L *lua.LState) int {
			L.Push(lua.LBool(netMgr.IsHost()))
			return 1
		}))

		// rf.network_sync(table, tier)
		L.SetField(rf, "network_sync", L.NewFunction(func(L *lua.LState) int {
			tbl := L.CheckTable(1)
			tierStr := L.CheckString(2)

			// Convert Lua table to map
			stateMap := make(map[string]interface{})
			tbl.ForEach(func(k, v lua.LValue) {
				key := k.String()
				var val interface{}
				switch lv := v.(type) {
				case lua.LNumber:
					val = float64(lv)
				case lua.LString:
					val = string(lv)
				case lua.LBool:
					val = bool(lv)
				case *lua.LTable:
					// Nested table - convert recursively (simplified)
					nested := make(map[string]interface{})
					lv.ForEach(func(k2, v2 lua.LValue) {
						nested[k2.String()] = v2.String()
					})
					val = nested
				default:
					val = lv.String()
				}
				stateMap[key] = val
			})

			// Register table for sync (using table path from Lua)
			var tier network.SyncTier
			switch tierStr {
			case "fast":
				tier = network.SyncTierFast
			case "moderate":
				tier = network.SyncTierModerate
			case "slow":
				tier = network.SyncTierSlow
			default:
				msg := lua.LString("invalid tier: must be 'fast', 'moderate', or 'slow'")
				L.Push(msg)
				return 1 // Return error message
			}

			// Use a simple path (we'll need to track table references properly)
			tablePath := "lua_table" // Simplified - in real implementation, track by Lua reference
			err := netMgr.RegisterSyncedTable(tablePath, tier, stateMap)
			if err != nil {
				msg := lua.LString(err.Error())
				L.Push(msg)
				return 1 // Return error message
			}

			return 0
		}))

		// rf.network_unsync(table)
		L.SetField(rf, "network_unsync", L.NewFunction(func(L *lua.LState) int {
			tbl := L.CheckTable(1)
			_ = tbl                  // For now, simplified - would need to track table references
			tablePath := "lua_table" // Simplified
			netMgr.UnregisterSyncedTable(tablePath)
			return 0
		}))
	}

	// Sound effects
	L.SetField(rf, "sfx", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		action := L.OptString(2, "")
		_ = audio.Init()

		// Check if SFX exists in loaded map
		if sfx, ok := sfxMap[name]; ok {
			switch sfx.Type {
			case "sine":
				audio.PlaySine(sfx.Freq, sfx.Duration, sfx.Gain)
			case "noise":
				audio.PlayNoise(sfx.Duration, sfx.Gain)
			case "thrust":
				audio.Thrust(action != "off")
			case "stopall":
				audio.StopAll()
			}
			return 0
		}

		// Fallback to hardcoded defaults for backward compatibility
		switch name {
		case "thrust":
			audio.Thrust(action != "off")
		case "land":
			audio.PlaySine(880, 0.12, 0.3)
		case "crash":
			audio.PlayNoise(0.25, 0.4)
		case "move":
			audio.PlaySine(520, 0.05, 0.25)
		case "select":
			audio.PlaySine(700, 0.08, 0.3)
		case "stopall":
			audio.StopAll()
		}
		return 0
	}))

	// Raw tone/noise
	L.SetField(rf, "tone", L.NewFunction(func(L *lua.LState) int {
		_ = audio.Init()
		f := L.CheckNumber(1)
		d := L.CheckNumber(2)
		g := L.OptNumber(3, 0.3)
		audio.PlaySine(float64(f), float64(d), float64(g))
		return 0
	}))
	L.SetField(rf, "noise", L.NewFunction(func(L *lua.LState) int {
		_ = audio.Init()
		d := L.CheckNumber(1)
		g := L.OptNumber(2, 0.3)
		audio.PlayNoise(float64(d), float64(g))
		return 0
	}))

	// Music: rf.music("trackname", bpm, gain) or rf.music({"1G#2","R1","A3"}, bpm, gain)
	L.SetField(rf, "music", L.NewFunction(func(L *lua.LState) int {
		_ = audio.Init()

		// Check if first arg is a string (track name) or table (inline notes)
		firstArg := L.Get(1)
		if str, ok := firstArg.(lua.LString); ok {
			// Try to find in music map
			trackName := string(str)
			if music, ok := musicMap[trackName]; ok {
				var bpm float64
				if L.GetTop() >= 2 {
					bpm = float64(L.OptNumber(2, lua.LNumber(music.BPM)))
				} else {
					bpm = music.BPM
				}
				if bpm == 0 {
					bpm = 120 // default
				}
				var gain float64
				if L.GetTop() >= 3 {
					gain = float64(L.OptNumber(3, lua.LNumber(music.Gain)))
				} else {
					gain = music.Gain
				}
				if gain == 0 {
					gain = 0.3 // default
				}
				audio.PlayNotes(music.Tokens, bpm, gain)
				return 0
			}
		}

		// Fallback to table (inline notes) for backward compatibility
		tbl := L.CheckTable(1)
		bpm := L.OptNumber(2, 120)
		gain := L.OptNumber(3, 0.3)
		var toks []string
		tbl.ForEach(func(k, v lua.LValue) {
			if s, ok := v.(lua.LString); ok {
				toks = append(toks, string(s))
			}
		})
		audio.PlayNotes(toks, float64(bpm), float64(gain))
		return 0
	}))

	// Store spritesMap pointer for modification
	// Maps are reference types in Go, so we can store the map directly
	spriteMapPtr := &spritesMap

	// Store pool manager for automatic pool registration
	// This allows pools to be created/updated when sprite properties change

	// Sprites: rf.sprite(name) returns table with width, height, pixels, useCollision, mountPoints, isUI, lifetime, maxSpawn
	L.SetField(rf, "sprite", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		sprite, ok := (*spriteMapPtr)[name]
		if !ok {
			L.Push(lua.LNil)
			return 1
		}

		// Return table: {width=w, height=h, pixels={{row1}, {row2}, ...}, useCollision=bool, mountPoints={{x,y}, ...}, isUI=bool, lifetime=int, maxSpawn=int}
		tbl := L.NewTable()
		tbl.RawSetString("width", lua.LNumber(sprite.Width))
		tbl.RawSetString("height", lua.LNumber(sprite.Height))

		pixelsTbl := L.NewTable()
		for y, row := range sprite.Pixels {
			rowTbl := L.NewTable()
			for x, colorIdx := range row {
				rowTbl.RawSetInt(x+1, lua.LNumber(colorIdx))
			}
			pixelsTbl.RawSetInt(y+1, rowTbl)
		}
		tbl.RawSetString("pixels", pixelsTbl)
		tbl.RawSetString("useCollision", lua.LBool(sprite.UseCollision))
		tbl.RawSetString("isUI", lua.LBool(sprite.IsUI))
		tbl.RawSetString("lifetime", lua.LNumber(sprite.Lifetime))
		tbl.RawSetString("maxSpawn", lua.LNumber(sprite.MaxSpawn))

		mountPointsTbl := L.NewTable()
		for i, mp := range sprite.MountPoints {
			mpTbl := L.NewTable()
			mpTbl.RawSetString("x", lua.LNumber(mp.X))
			mpTbl.RawSetString("y", lua.LNumber(mp.Y))
			if mp.Name != "" {
				mpTbl.RawSetString("name", lua.LString(mp.Name))
				// Also set by name for direct access: mountPoints["thrust"] -> mount point
				mountPointsTbl.RawSetString(mp.Name, mpTbl)
			}
			// Set by index (1-based): mountPoints[1], mountPoints[2], etc.
			mountPointsTbl.RawSetInt(i+1, mpTbl)
		}
		tbl.RawSetString("mountPoints", mountPointsTbl)

		L.Push(tbl)
		return 1
	}))

	// Sprite drawing: rf.spr(name, x, y, [flip_x, flip_y])
	// Draws a sprite by name at position (x, y) with optional flipping
	L.SetField(rf, "spr", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		x := L.CheckInt(2)
		y := L.CheckInt(3)
		flipX := L.OptBool(4, false)
		flipY := L.OptBool(5, false)

		sprite, ok := (*spriteMapPtr)[name]
		if !ok {
			return 0 // Sprite not found, do nothing
		}

		// Draw sprite pixels
		for sy := 0; sy < sprite.Height; sy++ {
			for sx := 0; sx < sprite.Width; sx++ {
				// Calculate source coordinates with flipping
				srcX := sx
				srcY := sy
				if flipX {
					srcX = sprite.Width - 1 - sx
				}
				if flipY {
					srcY = sprite.Height - 1 - sy
				}

				colorIdx := sprite.Pixels[srcY][srcX]
				if colorIdx >= 0 { // -1 is transparent
					c := colorByIndexRemapped(colorIdx)
					r.PSet(x+sx, y+sy, color.RGBA{c[0], c[1], c[2], c[3]})
				}
			}
		}
		return 0
	}))

	// Sprite region: rf.sspr(sx, sy, sw, sh, dx, dy, [dw, dh, flip_x, flip_y])
	// Draws a region of a sprite. For RetroForge, we'll use sprite name and draw sub-region
	// Note: PICO-8's sspr works differently (sprite sheet), but we'll adapt it
	L.SetField(rf, "sspr", L.NewFunction(func(L *lua.LState) int {
		// For now, simplified version - draw sprite region
		// Full implementation would need sprite sheet support
		name := L.OptString(1, "")
		if name == "" {
			return 0
		}
		sx := L.CheckInt(2)
		sy := L.CheckInt(3)
		sw := L.CheckInt(4)
		sh := L.CheckInt(5)
		dx := L.CheckInt(6)
		dy := L.CheckInt(7)
		dw := L.OptInt(8, sw)
		dh := L.OptInt(9, sh)
		flipX := L.OptBool(10, false)
		flipY := L.OptBool(11, false)

		sprite, ok := (*spriteMapPtr)[name]
		if !ok {
			return 0
		}

		// Draw scaled/flipped region
		xScale := float64(dw) / float64(sw)
		yScale := float64(dh) / float64(sh)

		for dyi := 0; dyi < dh; dyi++ {
			for dxi := 0; dxi < dw; dxi++ {
				// Map destination to source
				srcX := int(float64(dxi)/xScale) + sx
				srcY := int(float64(dyi)/yScale) + sy

				if srcX < 0 || srcY < 0 || srcX >= sprite.Width || srcY >= sprite.Height {
					continue
				}

				// Apply flipping
				drawX := srcX
				drawY := srcY
				if flipX {
					drawX = sprite.Width - 1 - drawX
				}
				if flipY {
					drawY = sprite.Height - 1 - drawY
				}

				colorIdx := sprite.Pixels[drawY][drawX]
				if colorIdx >= 0 {
					c := colorByIndexRemapped(colorIdx)
					r.PSet(dx+dxi, dy+dyi, color.RGBA{c[0], c[1], c[2], c[3]})
				}
			}
		}
		return 0
	}))

	// Sprite creation: rf.newSprite(name, width, height) -> sprite table
	// Creates a new empty sprite (all pixels transparent, defaults: isUI=true, lifetime=0, maxSpawn=0)
	L.SetField(rf, "newSprite", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		width := L.CheckInt(2)
		height := L.CheckInt(3)

		if width <= 0 || height <= 0 {
			L.RaiseError("sprite width and height must be positive")
			return 0
		}
		if width > 256 || height > 256 {
			L.RaiseError("sprite dimensions cannot exceed 256x256")
			return 0
		}

		// Initialize pixels with all transparent (-1)
		pixels := make([][]int, height)
		for y := range pixels {
			pixels[y] = make([]int, width)
			for x := range pixels[y] {
				pixels[y][x] = -1 // Transparent
			}
		}

		// Create new sprite with defaults
		newSprite := cartio.SpriteData{
			Width:        width,
			Height:       height,
			Pixels:       pixels,
			UseCollision: false,
			MountPoints:  []cartio.MountPoint{},
			IsUI:         true, // Default true
			Lifetime:     0,    // 0 = no lifetime limit
			MaxSpawn:     0,    // 0 = no spawn limit
		}

		// Add to sprite map
		(*spriteMapPtr)[name] = newSprite

		// Automatically register pool if sprite meets criteria (isUI=false, maxSpawn>10)
		// Note: Default sprite has isUI=true and maxSpawn=0, so it won't be pooled by default
		// Pool will be created automatically when properties are changed via setSpriteProperty

		// Return sprite table
		tbl := L.NewTable()
		tbl.RawSetString("width", lua.LNumber(width))
		tbl.RawSetString("height", lua.LNumber(height))

		pixelsTbl := L.NewTable()
		for y, row := range pixels {
			rowTbl := L.NewTable()
			for x, colorIdx := range row {
				rowTbl.RawSetInt(x+1, lua.LNumber(colorIdx))
			}
			pixelsTbl.RawSetInt(y+1, rowTbl)
		}
		tbl.RawSetString("pixels", pixelsTbl)
		tbl.RawSetString("useCollision", lua.LBool(false))
		tbl.RawSetString("isUI", lua.LBool(true))
		tbl.RawSetString("lifetime", lua.LNumber(0))
		tbl.RawSetString("maxSpawn", lua.LNumber(0))
		tbl.RawSetString("mountPoints", L.NewTable())

		L.Push(tbl)
		return 1
	}))

	// Helper functions for sprite drawing (defined before use)
	abs := func(n int) int {
		if n < 0 {
			return -n
		}
		return n
	}

	// Helper to draw line in sprite
	spriteLine := func(name string, x0, y0, x1, y1, idx int, sprite cartio.SpriteData, spriteMapPtr *cartio.SpriteMap) {
		dx := abs(x1 - x0)
		dy := abs(y1 - y0)
		sx := 1
		if x0 > x1 {
			sx = -1
		}
		sy := 1
		if y0 > y1 {
			sy = -1
		}
		err := dx - dy

		x, y := x0, y0
		for {
			if x >= 0 && x < sprite.Width && y >= 0 && y < sprite.Height {
				sprite.Pixels[y][x] = idx
			}

			if x == x1 && y == y1 {
				break
			}

			e2 := 2 * err
			if e2 > -dy {
				err -= dy
				x += sx
			}
			if e2 < dx {
				err += dx
				y += sy
			}
		}
		(*spriteMapPtr)[name] = sprite
	}

	// Helper to set pixel in sprite with bounds checking
	setSpritePixel := func(name string, x, y, idx int, sprite cartio.SpriteData, spriteMapPtr *cartio.SpriteMap) {
		if x >= 0 && x < sprite.Width && y >= 0 && y < sprite.Height {
			sprite.Pixels[y][x] = idx
			(*spriteMapPtr)[name] = sprite
		}
	}

	// Sprite primitive drawing functions
	// These draw to sprite pixels instead of the screen

	// rf.sprite_pset(sprite_name, x, y, index) - Set pixel in sprite
	L.SetField(rf, "sprite_pset", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		x := L.CheckInt(2)
		y := L.CheckInt(3)
		idx := L.CheckInt(4)

		sprite, ok := (*spriteMapPtr)[name]
		if !ok {
			L.RaiseError("sprite '%s' not found", name)
			return 0
		}

		if x < 0 || y < 0 || x >= sprite.Width || y >= sprite.Height {
			return 0 // Out of bounds, ignore
		}

		sprite.Pixels[y][x] = idx
		(*spriteMapPtr)[name] = sprite
		return 0
	}))

	// rf.sprite_line(sprite_name, x0, y0, x1, y1, index) - Draw line in sprite
	L.SetField(rf, "sprite_line", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		x0 := L.CheckInt(2)
		y0 := L.CheckInt(3)
		x1 := L.CheckInt(4)
		y1 := L.CheckInt(5)
		idx := L.CheckInt(6)

		sprite, ok := (*spriteMapPtr)[name]
		if !ok {
			L.RaiseError("sprite '%s' not found", name)
			return 0
		}

		// Bresenham's line algorithm
		dx := abs(x1 - x0)
		dy := abs(y1 - y0)
		sx := 1
		if x0 > x1 {
			sx = -1
		}
		sy := 1
		if y0 > y1 {
			sy = -1
		}
		err := dx - dy

		x, y := x0, y0
		for {
			if x >= 0 && x < sprite.Width && y >= 0 && y < sprite.Height {
				sprite.Pixels[y][x] = idx
			}

			if x == x1 && y == y1 {
				break
			}

			e2 := 2 * err
			if e2 > -dy {
				err -= dy
				x += sx
			}
			if e2 < dx {
				err += dx
				y += sy
			}
		}

		(*spriteMapPtr)[name] = sprite
		return 0
	}))

	// rf.sprite_rect(sprite_name, x0, y0, x1, y1, index) - Draw rectangle outline in sprite
	L.SetField(rf, "sprite_rect", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		x0 := L.CheckInt(2)
		y0 := L.CheckInt(3)
		x1 := L.CheckInt(4)
		y1 := L.CheckInt(5)
		idx := L.CheckInt(6)

		sprite, ok := (*spriteMapPtr)[name]
		if !ok {
			L.RaiseError("sprite '%s' not found", name)
			return 0
		}

		// Draw four lines
		spriteLine(name, x0, y0, x1, y0, idx, sprite, spriteMapPtr) // Top
		spriteLine(name, x1, y0, x1, y1, idx, sprite, spriteMapPtr) // Right
		spriteLine(name, x1, y1, x0, y1, idx, sprite, spriteMapPtr) // Bottom
		spriteLine(name, x0, y1, x0, y0, idx, sprite, spriteMapPtr) // Left

		return 0
	}))

	// rf.sprite_rectfill(sprite_name, x0, y0, x1, y1, index) - Draw filled rectangle in sprite
	L.SetField(rf, "sprite_rectfill", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		x0 := L.CheckInt(2)
		y0 := L.CheckInt(3)
		x1 := L.CheckInt(4)
		y1 := L.CheckInt(5)
		idx := L.CheckInt(6)

		sprite, ok := (*spriteMapPtr)[name]
		if !ok {
			L.RaiseError("sprite '%s' not found", name)
			return 0
		}

		// Ensure x0 < x1 and y0 < y1
		if x0 > x1 {
			x0, x1 = x1, x0
		}
		if y0 > y1 {
			y0, y1 = y1, y0
		}

		// Fill rectangle
		for y := y0; y <= y1; y++ {
			if y >= 0 && y < sprite.Height {
				for x := x0; x <= x1; x++ {
					if x >= 0 && x < sprite.Width {
						sprite.Pixels[y][x] = idx
					}
				}
			}
		}

		(*spriteMapPtr)[name] = sprite
		return 0
	}))

	// rf.sprite_circ(sprite_name, x, y, radius, index) - Draw circle outline in sprite
	L.SetField(rf, "sprite_circ", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		x := L.CheckInt(2)
		y := L.CheckInt(3)
		radius := L.CheckInt(4)
		idx := L.CheckInt(5)

		sprite, ok := (*spriteMapPtr)[name]
		if !ok {
			L.RaiseError("sprite '%s' not found", name)
			return 0
		}

		// Midpoint circle algorithm
		xx := radius
		yy := 0
		err := 0

		for xx >= yy {
			// Draw 8 points of symmetry
			setSpritePixel(name, x+xx, y+yy, idx, sprite, spriteMapPtr)
			setSpritePixel(name, x-xx, y+yy, idx, sprite, spriteMapPtr)
			setSpritePixel(name, x+xx, y-yy, idx, sprite, spriteMapPtr)
			setSpritePixel(name, x-xx, y-yy, idx, sprite, spriteMapPtr)
			setSpritePixel(name, x+yy, y+xx, idx, sprite, spriteMapPtr)
			setSpritePixel(name, x-yy, y+xx, idx, sprite, spriteMapPtr)
			setSpritePixel(name, x+yy, y-xx, idx, sprite, spriteMapPtr)
			setSpritePixel(name, x-yy, y-xx, idx, sprite, spriteMapPtr)

			if err <= 0 {
				yy++
				err += 2*yy + 1
			}
			if err > 0 {
				xx--
				err -= 2*xx + 1
			}
		}

		return 0
	}))

	// rf.sprite_circfill(sprite_name, x, y, radius, index) - Draw filled circle in sprite
	L.SetField(rf, "sprite_circfill", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		x := L.CheckInt(2)
		y := L.CheckInt(3)
		radius := L.CheckInt(4)
		idx := L.CheckInt(5)

		sprite, ok := (*spriteMapPtr)[name]
		if !ok {
			L.RaiseError("sprite '%s' not found", name)
			return 0
		}

		// Fill circle by scanning and checking distance
		for sy := -radius; sy <= radius; sy++ {
			if y+sy >= 0 && y+sy < sprite.Height {
				for sx := -radius; sx <= radius; sx++ {
					if x+sx >= 0 && x+sx < sprite.Width {
						if sx*sx+sy*sy <= radius*radius {
							sprite.Pixels[y+sy][x+sx] = idx
						}
					}
				}
			}
		}

		(*spriteMapPtr)[name] = sprite
		return 0
	}))

	// rf.setSpriteProperty(sprite_name, property, value) - Set sprite property (useCollision, isUI, lifetime, maxSpawn)
	L.SetField(rf, "setSpriteProperty", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		property := L.CheckString(2)
		value := L.CheckAny(3)

		sprite, ok := (*spriteMapPtr)[name]
		if !ok {
			L.RaiseError("sprite '%s' not found", name)
			return 0
		}

		switch property {
		case "useCollision":
			if b, ok := value.(lua.LBool); ok {
				sprite.UseCollision = bool(b)
			} else {
				L.RaiseError("useCollision must be boolean")
				return 0
			}
		case "isUI":
			if b, ok := value.(lua.LBool); ok {
				sprite.IsUI = bool(b)
			} else {
				L.RaiseError("isUI must be boolean")
				return 0
			}
		case "lifetime":
			if n, ok := value.(lua.LNumber); ok {
				sprite.Lifetime = int(n)
			} else {
				L.RaiseError("lifetime must be number")
				return 0
			}
		case "maxSpawn":
			if n, ok := value.(lua.LNumber); ok {
				sprite.MaxSpawn = int(n)
			} else {
				L.RaiseError("maxSpawn must be number")
				return 0
			}
		default:
			L.RaiseError("unknown property: %s (use: useCollision, isUI, lifetime, maxSpawn)", property)
			return 0
		}

		(*spriteMapPtr)[name] = sprite

		// Automatically register/update pool if sprite now meets criteria
		if spritepool.ShouldPool(sprite) {
			// Register pool if it doesn't exist
			if err := poolManager.RegisterSprite(name, sprite); err != nil {
				// Log error but don't fail - pooling is optional
				if devMode != nil {
					devMode.AddDebugLog(fmt.Sprintf("Failed to create pool for sprite '%s': %v", name, err))
				}
			}
		} else {
			// Sprite no longer meets criteria - remove pool if it exists
			if poolManager.HasPool(name) {
				poolManager.RemovePool(name)
			}
		}

		return 0
	}))

	// Tilemap functions: mget, mset, map
	L.SetField(rf, "mget", L.NewFunction(func(L *lua.LState) int {
		x := L.CheckInt(1)
		y := L.CheckInt(2)
		val := state.GetTileMap().Get(x, y)
		L.Push(lua.LNumber(val))
		return 1
	}))
	L.SetField(rf, "mset", L.NewFunction(func(L *lua.LState) int {
		x := L.CheckInt(1)
		y := L.CheckInt(2)
		v := L.CheckInt(3)
		state.GetTileMap().Set(x, y, v)
		return 0
	}))
	L.SetField(rf, "map", L.NewFunction(func(L *lua.LState) int {
		celX := L.CheckInt(1)
		celY := L.CheckInt(2)
		sx := L.CheckInt(3)
		sy := L.CheckInt(4)
		celW := L.CheckInt(5)
		celH := L.CheckInt(6)

		// Draw tilemap region using sprites
		tm := state.GetTileMap()
		tm.Draw(celX, celY, sx, sy, celW, celH, func(x, y, tileIndex int) {
			// Convert tile index to sprite name (simplified: assume tile index is sprite name index)
			// For now, draw an 8x8 rectangle representing the tile
			// Full implementation would look up sprite by index
			if tileIndex > 0 {
				c := colorByIndexRemapped(tileIndex % 50) // Use tile index as color
				r.RectFill(x, y, x+7, y+7, color.RGBA{c[0], c[1], c[2], c[3]})
			}
		})
		return 0
	}))

	// Color remapping: pal(c0, c1, [p])
	L.SetField(rf, "pal", L.NewFunction(func(L *lua.LState) int {
		if L.GetTop() == 0 {
			// No args = reset all remapping
			state.ResetPalRemap()
		} else {
			c0 := L.CheckInt(1)
			c1 := L.OptInt(2, c0) // Default to same color if not provided
			p := L.OptBool(3, true)
			state.SetPalRemap(c0, c1, p)
		}
		return 0
	}))

	// Memory functions: poke, peek
	L.SetField(rf, "poke", L.NewFunction(func(L *lua.LState) int {
		addr := L.CheckInt(1)
		val := L.CheckInt(2)
		mem := state.GetMemory()
		if addr >= 0 && addr < len(mem) {
			mem[addr] = byte(val & 0xFF)
		}
		return 0
	}))
	L.SetField(rf, "peek", L.NewFunction(func(L *lua.LState) int {
		addr := L.CheckInt(1)
		mem := state.GetMemory()
		if addr >= 0 && addr < len(mem) {
			L.Push(lua.LNumber(mem[addr]))
		} else {
			L.Push(lua.LNumber(0))
		}
		return 1
	}))
	L.SetField(rf, "poke2", L.NewFunction(func(L *lua.LState) int {
		addr := L.CheckInt(1)
		val := L.CheckInt(2)
		mem := state.GetMemory()
		if addr >= 0 && addr+1 < len(mem) {
			mem[addr] = byte(val & 0xFF)
			mem[addr+1] = byte((val >> 8) & 0xFF)
		}
		return 0
	}))
	L.SetField(rf, "peek2", L.NewFunction(func(L *lua.LState) int {
		addr := L.CheckInt(1)
		mem := state.GetMemory()
		if addr >= 0 && addr+1 < len(mem) {
			val := int(mem[addr]) | (int(mem[addr+1]) << 8)
			L.Push(lua.LNumber(val))
		} else {
			L.Push(lua.LNumber(0))
		}
		return 1
	}))
	L.SetField(rf, "poke4", L.NewFunction(func(L *lua.LState) int {
		addr := L.CheckInt(1)
		val := int64(L.CheckNumber(2))
		mem := state.GetMemory()
		if addr >= 0 && addr+3 < len(mem) {
			mem[addr] = byte(val & 0xFF)
			mem[addr+1] = byte((val >> 8) & 0xFF)
			mem[addr+2] = byte((val >> 16) & 0xFF)
			mem[addr+3] = byte((val >> 24) & 0xFF)
		}
		return 0
	}))
	L.SetField(rf, "peek4", L.NewFunction(func(L *lua.LState) int {
		addr := L.CheckInt(1)
		mem := state.GetMemory()
		if addr >= 0 && addr+3 < len(mem) {
			val := int64(mem[addr]) | (int64(mem[addr+1]) << 8) | (int64(mem[addr+2]) << 16) | (int64(mem[addr+3]) << 24)
			L.Push(lua.LNumber(val))
		} else {
			L.Push(lua.LNumber(0))
		}
		return 1
	}))

	// Cart persistence: cstore(dest_addr, src_addr, len) - Copy from runtime memory to cart storage
	L.SetField(rf, "cstore", L.NewFunction(func(L *lua.LState) int {
		destAddr := L.CheckInt(1)
		srcAddr := L.CheckInt(2)
		length := L.CheckInt(3)

		runtimeMem := state.GetMemory()
		cartStore := state.GetCartStore()

		// Validate addresses and length
		if srcAddr < 0 || destAddr < 0 || length < 0 {
			return 0 // Invalid parameters, do nothing
		}
		if srcAddr >= len(runtimeMem) || destAddr >= len(cartStore) {
			return 0 // Out of bounds
		}
		if srcAddr+length > len(runtimeMem) {
			length = len(runtimeMem) - srcAddr // Clamp to available memory
		}
		if destAddr+length > len(cartStore) {
			length = len(cartStore) - destAddr // Clamp to available cart storage
		}

		// Copy bytes from runtime memory to cart storage
		copy(cartStore[destAddr:destAddr+length], runtimeMem[srcAddr:srcAddr+length])
		return 0
	}))

	// Cart persistence: reload(dest_addr, src_addr, len) - Copy from cart storage to runtime memory
	L.SetField(rf, "reload", L.NewFunction(func(L *lua.LState) int {
		destAddr := L.CheckInt(1)
		srcAddr := L.CheckInt(2)
		length := L.CheckInt(3)

		runtimeMem := state.GetMemory()
		cartStore := state.GetCartStore()

		// Validate addresses and length
		if srcAddr < 0 || destAddr < 0 || length < 0 {
			return 0 // Invalid parameters, do nothing
		}
		if srcAddr >= len(cartStore) || destAddr >= len(runtimeMem) {
			return 0 // Out of bounds
		}
		if srcAddr+length > len(cartStore) {
			length = len(cartStore) - srcAddr // Clamp to available cart storage
		}
		if destAddr+length > len(runtimeMem) {
			length = len(runtimeMem) - destAddr // Clamp to available memory
		}

		// Copy bytes from cart storage to runtime memory
		copy(runtimeMem[destAddr:destAddr+length], cartStore[srcAddr:srcAddr+length])
		return 0
	}))

	// Physics functions (only if physics world is provided)
	if physWorld != nil {
		// Store physics bodies in Lua userdata
		physicsBodies := make(map[int]*physics.Body)
		nextBodyId := 1

		L.SetField(rf, "physics_create_body", L.NewFunction(func(L *lua.LState) int {
			bodyTypeStr := L.CheckString(1)
			x := float64(L.CheckNumber(2))
			y := float64(L.CheckNumber(3))

			var body *physics.Body
			switch bodyTypeStr {
			case "static":
				body = physWorld.CreateStaticBody(x, y)
			case "dynamic":
				body = physWorld.CreateDynamicBody(x, y)
			case "kinematic":
				body = physWorld.CreateKinematicBody(x, y)
			default:
				L.Push(lua.LNil)
				return 1
			}

			bodyId := nextBodyId
			nextBodyId++
			physicsBodies[bodyId] = body
			L.Push(lua.LNumber(bodyId))
			return 1
		}))

		L.SetField(rf, "physics_body_add_box", L.NewFunction(func(L *lua.LState) int {
			bodyId := L.CheckInt(1)
			width := float64(L.CheckNumber(2))
			height := float64(L.CheckNumber(3))
			density := float64(L.OptNumber(4, 1.0))
			restitution := float64(L.OptNumber(5, 0.0))
			friction := float64(L.OptNumber(6, 0.2))

			body, ok := physicsBodies[bodyId]
			if !ok {
				return 0
			}
			// Use properties version if restitution or friction provided
			if restitution > 0 || friction != 0.2 {
				body.CreateBoxFixtureWithProps(width, height, density, restitution, friction)
			} else {
				body.CreateBoxFixture(width, height, density)
			}
			return 0
		}))

		L.SetField(rf, "physics_body_add_circle", L.NewFunction(func(L *lua.LState) int {
			bodyId := L.CheckInt(1)
			radius := float64(L.CheckNumber(2))
			density := float64(L.OptNumber(3, 1.0))
			restitution := float64(L.OptNumber(4, 0.0))
			friction := float64(L.OptNumber(5, 0.2))

			body, ok := physicsBodies[bodyId]
			if !ok {
				return 0
			}
			// Use properties version if restitution or friction provided
			if restitution > 0 || friction != 0.2 {
				body.CreateCircleFixtureWithProps(radius, density, restitution, friction)
			} else {
				body.CreateCircleFixture(radius, density)
			}
			return 0
		}))

		L.SetField(rf, "physics_body_set_position", L.NewFunction(func(L *lua.LState) int {
			bodyId := L.CheckInt(1)
			x := float64(L.CheckNumber(2))
			y := float64(L.CheckNumber(3))

			body, ok := physicsBodies[bodyId]
			if !ok {
				return 0
			}
			body.SetPosition(x, y)
			return 0
		}))

		L.SetField(rf, "physics_body_get_position", L.NewFunction(func(L *lua.LState) int {
			bodyId := L.CheckInt(1)

			body, ok := physicsBodies[bodyId]
			if !ok {
				L.Push(lua.LNumber(0))
				L.Push(lua.LNumber(0))
				return 2
			}
			x, y := body.GetPosition()
			L.Push(lua.LNumber(x))
			L.Push(lua.LNumber(y))
			return 2
		}))

		L.SetField(rf, "physics_body_set_velocity", L.NewFunction(func(L *lua.LState) int {
			bodyId := L.CheckInt(1)
			vx := float64(L.CheckNumber(2))
			vy := float64(L.CheckNumber(3))

			body, ok := physicsBodies[bodyId]
			if !ok {
				return 0
			}
			body.SetVelocity(vx, vy)
			return 0
		}))

		// rf.physics_body_set_gravity_scale(bodyId, scale) - Set gravity scale (0 = no gravity, 1 = normal)
		L.SetField(rf, "physics_body_set_gravity_scale", L.NewFunction(func(L *lua.LState) int {
			bodyId := L.CheckInt(1)
			scale := float64(L.CheckNumber(2))

			body, ok := physicsBodies[bodyId]
			if !ok {
				return 0
			}
			body.SetGravityScale(scale)
			return 0
		}))

		L.SetField(rf, "physics_body_get_velocity", L.NewFunction(func(L *lua.LState) int {
			bodyId := L.CheckInt(1)

			body, ok := physicsBodies[bodyId]
			if !ok {
				L.Push(lua.LNumber(0))
				L.Push(lua.LNumber(0))
				return 2
			}
			vx, vy := body.GetVelocity()
			L.Push(lua.LNumber(vx))
			L.Push(lua.LNumber(vy))
			return 2
		}))

		L.SetField(rf, "physics_body_apply_force", L.NewFunction(func(L *lua.LState) int {
			bodyId := L.CheckInt(1)
			fx := float64(L.CheckNumber(2))
			fy := float64(L.CheckNumber(3))
			px := float64(L.CheckNumber(4))
			py := float64(L.CheckNumber(5))

			body, ok := physicsBodies[bodyId]
			if !ok {
				return 0
			}
			body.ApplyForce(fx, fy, px, py)
			return 0
		}))

		L.SetField(rf, "physics_body_destroy", L.NewFunction(func(L *lua.LState) int {
			bodyId := L.CheckInt(1)

			body, ok := physicsBodies[bodyId]
			if !ok {
				return 0
			}
			body.Destroy()
			delete(physicsBodies, bodyId)
			return 0
		}))
	}

	// Quit request
	L.SetField(rf, "quit", L.NewFunction(func(L *lua.LState) int {
		app.RequestQuit()
		return 0
	}))

	// Debug functions (only available in development mode)
	if devModePtr != nil {
		// printh(str) - Print to debug log
		L.SetField(rf, "printh", L.NewFunction(func(L *lua.LState) int {
			if devModePtr.IsEnabled() {
				str := L.OptString(1, "")
				devModePtr.AddDebugLog(str)
			}
			return 0
		}))

		// stat(n) - Get system statistics
		L.SetField(rf, "stat", L.NewFunction(func(L *lua.LState) int {
			if !devModePtr.IsEnabled() {
				L.Push(lua.LNumber(0))
				return 1
			}

			statId := L.OptInt(1, 0)
			statsIface := devModePtr.GetStats()

			// Type assert to DevStats structure
			type StatsStruct struct {
				FPS         float64
				FrameCount  int64
				LuaMemory   int64
				LoadTime    time.Duration
				LastReload  time.Time
				ReloadCount int
			}
			stats, ok := statsIface.(StatsStruct)
			if !ok {
				// Try map-based access as fallback
				L.Push(lua.LNumber(0))
				return 1
			}

			switch statId {
			case 0: // FPS
				L.Push(lua.LNumber(stats.FPS))
				return 1
			case 1: // Frame count
				L.Push(lua.LNumber(stats.FrameCount))
				return 1
			case 2: // Lua memory (in bytes)
				L.Push(lua.LNumber(stats.LuaMemory))
				return 1
			case 3: // Load time (in milliseconds)
				L.Push(lua.LNumber(stats.LoadTime.Milliseconds()))
				return 1
			case 4: // Last reload time (Unix timestamp)
				if stats.LastReload.IsZero() {
					L.Push(lua.LNumber(0))
				} else {
					L.Push(lua.LNumber(stats.LastReload.Unix()))
				}
				return 1
			case 5: // Reload count
				L.Push(lua.LNumber(stats.ReloadCount))
				return 1
			default:
				L.Push(lua.LNumber(0))
				return 1
			}
		}))

		// time() - Get current time in seconds (Unix timestamp)
		L.SetField(rf, "time", L.NewFunction(func(L *lua.LState) int {
			if devModePtr.IsEnabled() {
				now := time.Now()
				L.Push(lua.LNumber(float64(now.UnixNano()) / 1e9))
			} else {
				L.Push(lua.LNumber(0))
			}
			return 1
		}))
	}

	// Bitwise operations (PICO-8-compatible)
	// rf.shl(x, y) - Shift left: x << y
	L.SetField(rf, "shl", L.NewFunction(func(L *lua.LState) int {
		x := int64(L.CheckNumber(1))
		y := int64(L.CheckNumber(2))
		if y < 0 {
			// Negative shift = shift right
			result := x >> (-y)
			L.Push(lua.LNumber(result))
		} else if y > 63 {
			// Shift by more than 63 bits = 0 (for 64-bit integers)
			L.Push(lua.LNumber(0))
		} else {
			result := x << y
			L.Push(lua.LNumber(result))
		}
		return 1
	}))

	// rf.shr(x, y) - Shift right (arithmetic): x >> y
	L.SetField(rf, "shr", L.NewFunction(func(L *lua.LState) int {
		x := int64(L.CheckNumber(1))
		y := int64(L.CheckNumber(2))
		if y < 0 {
			// Negative shift = shift left
			result := x << (-y)
			L.Push(lua.LNumber(result))
		} else if y > 63 {
			// Shift by more than 63 bits = 0 or sign-extended
			if x < 0 {
				L.Push(lua.LNumber(-1)) // Sign-extended for negative numbers
			} else {
				L.Push(lua.LNumber(0))
			}
		} else {
			result := x >> y
			L.Push(lua.LNumber(result))
		}
		return 1
	}))

	// rf.band(x, y) - Bitwise AND: x & y
	L.SetField(rf, "band", L.NewFunction(func(L *lua.LState) int {
		x := int64(L.CheckNumber(1))
		y := int64(L.CheckNumber(2))
		result := x & y
		L.Push(lua.LNumber(result))
		return 1
	}))

	// rf.bor(x, y) - Bitwise OR: x | y
	L.SetField(rf, "bor", L.NewFunction(func(L *lua.LState) int {
		x := int64(L.CheckNumber(1))
		y := int64(L.CheckNumber(2))
		result := x | y
		L.Push(lua.LNumber(result))
		return 1
	}))

	// rf.bxor(x, y) - Bitwise XOR: x ~ y (Lua uses ~ for XOR)
	L.SetField(rf, "bxor", L.NewFunction(func(L *lua.LState) int {
		x := int64(L.CheckNumber(1))
		y := int64(L.CheckNumber(2))
		result := x ^ y
		L.Push(lua.LNumber(result))
		return 1
	}))

	// rf.bnot(x) - Bitwise NOT: ~x
	L.SetField(rf, "bnot", L.NewFunction(func(L *lua.LState) int {
		x := int64(L.CheckNumber(1))
		result := ^x
		L.Push(lua.LNumber(result))
		return 1
	}))

	// PICO-8-style helper functions

	// rf.flr(x) - Floor function: math.floor(x)
	L.SetField(rf, "flr", L.NewFunction(func(L *lua.LState) int {
		x := L.CheckNumber(1)
		L.Push(lua.LNumber(float64(int64(x)))) // Truncate towards zero (Lua's floor)
		return 1
	}))

	// rf.ceil(x) - Ceiling function: math.ceil(x)
	L.SetField(rf, "ceil", L.NewFunction(func(L *lua.LState) int {
		x := float64(L.CheckNumber(1))
		if x == float64(int64(x)) {
			L.Push(lua.LNumber(x))
		} else if x > 0 {
			L.Push(lua.LNumber(float64(int64(x) + 1)))
		} else {
			L.Push(lua.LNumber(float64(int64(x))))
		}
		return 1
	}))

	// rf.rnd([x]) - Random number: 0-1 if no arg, 0-x if arg provided
	L.SetField(rf, "rnd", L.NewFunction(func(L *lua.LState) int {
		if L.GetTop() == 0 {
			// No arguments: return 0.0 to 1.0 (exclusive of 1.0)
			val := state.NextRandom()
			L.Push(lua.LNumber(val))
		} else {
			// Argument provided: return 0.0 to x (exclusive of x)
			x := float64(L.CheckNumber(1))
			if x < 0 {
				// Negative range: return x to 0.0
				val := state.NextRandom() * x // Will be negative
				L.Push(lua.LNumber(val))
			} else if x == 0 {
				L.Push(lua.LNumber(0))
			} else {
				val := state.NextRandom() * x
				L.Push(lua.LNumber(val))
			}
		}
		return 1
	}))

	// rf.mid(x, y, z) - Clamp value x between y and z
	L.SetField(rf, "mid", L.NewFunction(func(L *lua.LState) int {
		x := float64(L.CheckNumber(1))
		y := float64(L.CheckNumber(2))
		z := float64(L.CheckNumber(3))
		// Ensure y <= z (swap if needed)
		if y > z {
			y, z = z, y
		}
		result := x
		if result < y {
			result = y
		}
		if result > z {
			result = z
		}
		L.Push(lua.LNumber(result))
		return 1
	}))

	// rf.sgn(x) - Sign function: -1 if x < 0, 0 if x == 0, 1 if x > 0
	L.SetField(rf, "sgn", L.NewFunction(func(L *lua.LState) int {
		x := float64(L.CheckNumber(1))
		var result float64
		if x > 0 {
			result = 1
		} else if x < 0 {
			result = -1
		} else {
			result = 0
		}
		L.Push(lua.LNumber(result))
		return 1
	}))

	// rf.chr(n) - Convert number to character
	L.SetField(rf, "chr", L.NewFunction(func(L *lua.LState) int {
		n := int(L.CheckNumber(1))
		// Clamp to valid byte range (0-255)
		if n < 0 {
			n = 0
		} else if n > 255 {
			n = 255
		}
		// Create a string with a single byte (not UTF-8 encoded rune)
		L.Push(lua.LString(string([]byte{byte(n)})))
		return 1
	}))

	// rf.ord(c) - Convert character to number (first character of string)
	L.SetField(rf, "ord", L.NewFunction(func(L *lua.LState) int {
		str := L.CheckString(1)
		if len(str) == 0 {
			L.Push(lua.LNumber(0))
		} else {
			L.Push(lua.LNumber(int(str[0])))
		}
		return 1
	}))
}
