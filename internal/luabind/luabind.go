package luabind

import (
	"fmt"
	"image/color"
	"sort"
	"strings"
	"time"

	"github.com/AndrewDonelson/retroforge-engine/internal/animation"
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

// AnimationUpdater provides a way to update animation states for direct sprites
type AnimationUpdater struct {
	animationStates map[string]*animation.AnimationState
	spritesMap      *cartio.SpriteMap
}

// UpdateAnimations updates all animation states for direct sprites
func (au *AnimationUpdater) UpdateAnimations(deltaTime time.Duration) {
	deltaTimeMs := int64(deltaTime / time.Millisecond)
	for spriteName, animState := range au.animationStates {
		if animState != nil && animState.Playing && !animState.Paused {
			if sprite, ok := (*au.spritesMap)[spriteName]; ok {
				if sprite.Type == cartio.SpriteTypeAnimation {
					animation.UpdateAnimationState(animState, &sprite, deltaTimeMs)
				}
			}
		}
	}
}

// Register attaches rf.* drawing functions to the Lua state.
func Register(L *lua.LState, r graphics.Renderer, colorByIndex ColorByIndex, setPalette func(string), sfxMap cartio.SFXMap, musicMap cartio.MusicMap, spritesMap cartio.SpriteMap, physWorld *physics.World, netMgr *network.NetworkManager) {
	state := NewState()
	RegisterWithState(L, r, colorByIndex, setPalette, sfxMap, musicMap, spritesMap, nil, physWorld, state, netMgr)
}

// RegisterWithDev attaches rf.* drawing functions with dev mode support
// Note: This function doesn't receive tilemapsMap, so it passes nil
// For full tilemap support, use RegisterWithDevMode directly
func RegisterWithDev(L *lua.LState, r graphics.Renderer, colorByIndex ColorByIndex, setPalette func(string), sfxMap cartio.SFXMap, musicMap cartio.MusicMap, spritesMap cartio.SpriteMap, physWorld *physics.World, devMode DevModeHandler, netMgr *network.NetworkManager) {
	state := NewState()
	RegisterWithDevMode(L, r, colorByIndex, setPalette, sfxMap, musicMap, spritesMap, nil, physWorld, state, devMode, netMgr)
}

// RegisterWithState attaches rf.* drawing functions with state management
func RegisterWithState(L *lua.LState, r graphics.Renderer, colorByIndex ColorByIndex, setPalette func(string), sfxMap cartio.SFXMap, musicMap cartio.MusicMap, spritesMap cartio.SpriteMap, tilemapsMap map[string]*cartio.TileMapData, physWorld *physics.World, state *State, netMgr *network.NetworkManager) {
	RegisterWithDevMode(L, r, colorByIndex, setPalette, sfxMap, musicMap, spritesMap, tilemapsMap, physWorld, state, nil, netMgr)
}

// RegisterWithDevMode attaches rf.* drawing functions with dev mode support
func RegisterWithDevMode(L *lua.LState, r graphics.Renderer, colorByIndex ColorByIndex, setPalette func(string), sfxMap cartio.SFXMap, musicMap cartio.MusicMap, spritesMap cartio.SpriteMap, tilemapsMap map[string]*cartio.TileMapData, physWorld *physics.World, state *State, devMode DevModeHandler, netMgr *network.NetworkManager) {
	rf := L.NewTable()
	L.SetGlobal("rf", rf)

	// Store devMode in closure for debug functions
	devModePtr := devMode
	
	// Helper function for debug logging (only logs when dev mode is enabled)
	debugLog := func(format string, args ...interface{}) {
		if devModePtr != nil && devModePtr.IsEnabled() {
			msg := fmt.Sprintf(format, args...)
			devModePtr.AddDebugLog(msg)
		}
	}

	// Create pool manager for automatic sprite pooling
	poolManager := spritepool.NewPoolManager()
	
	// Track animation states for non-pooled sprite drawing (for direct rf.spr() calls)
	// Key: sprite name, Value: animation state
	animationStates := make(map[string]*animation.AnimationState)
	
	// Note: Animation states for direct sprites are updated lazily in rf.spr() calls
	// For pooled sprites, animations are updated in pool.Update() which is called in engine loop
	// For better accuracy, we could expose animUpdater to engine, but lazy update works for now

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
		result := input.Btnp(i)
		L.Push(lua.LBool(result))
		return 1
	}))
	L.SetField(rf, "btnr", L.NewFunction(func(L *lua.LState) int {
		_ = L.CheckInt(1) // Button ID - release detection not yet implemented
		// Button release detection (not yet implemented in input package)
		// For now, return false
		L.Push(lua.LBool(false))
		return 1
	}))
	
	// Modifier keys
	L.SetField(rf, "shift", L.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LBool(input.Shift()))
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

		// Get tokens for this SFX (handles both new and legacy formats)
		tokens := cartio.GetSFXTokens(sfxMap, name)
		if len(tokens) > 0 {
			// Check if this is a single loop token (L code) - handle immediately
			if len(tokens) == 1 {
				token := tokens[0]
				waveType, freq, _, gain, err := cartio.ParseToken(token)
				if err == nil && waveType == "thrust" {
					// Single loop token - handle on/off directly
					if action != "off" {
						audio.PlayThrust(true, freq, gain)
					} else {
						audio.PlayThrust(false, freq, gain)
					}
					return 0
				}
			}

			// Handle special stopall token
			if len(tokens) == 1 && tokens[0] == "STOPALL" {
				audio.StopAll()
				return 0
			}

			// Play tokens sequentially (for multi-token sequences)
			go func() {
				for _, token := range tokens {
					// Handle special stopall token
					if token == "STOPALL" {
						audio.StopAll()
						continue
					}

					waveType, freq, duration, gain, err := cartio.ParseToken(token)
					if err != nil {
						continue // Skip invalid tokens
					}

					switch waveType {
					case "sine":
						audio.PlaySine(freq, duration, gain)
					case "noise":
						audio.PlayNoise(duration, gain)
					case "thrust":
						// For thrust tokens in sequences, still play them but don't wait
						// (This shouldn't happen in normal sequences, but handle it)
						if action != "off" {
							audio.PlayThrust(true, freq, gain)
						} else {
							audio.PlayThrust(false, freq, gain)
						}
						// For loops, we don't wait - they play until stopped
						continue
					case "rest":
						// Rest: just wait (duration already calculated)
						time.Sleep(time.Duration(duration * float64(time.Second)))
						continue
					}

					// Wait for non-looped sounds to finish
					time.Sleep(time.Duration(duration * float64(time.Second)))
				}
			}()
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
				// Play music with loop support
				loop := music.Loop
				if loop {
					// Play in a loop until stopped
					// StopAll() will set musicStopFlag, which PlayNotes checks internally
					go func() {
						for {
							audio.PlayNotes(music.Tokens, bpm, gain)
							// Calculate total duration to wait before looping
							beat := 60.0 / bpm
							totalBeats := 0.0
							for _, t := range music.Tokens {
								s := strings.ToUpper(strings.TrimSpace(t))
								if s == "" {
									continue
								}
								length := 1
								if last := s[len(s)-1]; last >= '0' && last <= '9' {
									length = int(last - '0')
								}
								totalBeats += float64(length)
							}
							totalDuration := totalBeats * beat
							// Sleep in chunks to allow responsive stopping
							chunkDuration := 100 * time.Millisecond
							elapsed := time.Duration(0)
							for elapsed < time.Duration(totalDuration*float64(time.Second)) {
								time.Sleep(chunkDuration)
								elapsed += chunkDuration
								// PlayNotes will check musicStopFlag internally and exit early if set
							}
						}
					}()
				} else {
					// Play once
					audio.PlayNotes(music.Tokens, bpm, gain)
				}
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
		tbl.RawSetString("type", lua.LString(string(sprite.Type)))

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

	// Sprite drawing: rf.spr(name, x, y, [frameNameOrAnimation, flip_x, flip_y])
	// Draws a sprite by name at position (x, y)
	// For static sprites: rf.spr(name, x, y, [flip_x, flip_y])
	// For frames sprites: rf.spr(name, x, y, frameName, [flip_x, flip_y])
	// For animation sprites: rf.spr(name, x, y, [animationName, flip_x, flip_y]) - uses active animation state
	L.SetField(rf, "spr", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		x := L.CheckInt(2)
		y := L.CheckInt(3)
		
		// Parse optional arguments - could be frameName/animation, flipX, or just flipX
		var frameNameOrAnim string
		var flipX, flipY bool
		
		argCount := L.GetTop()
		if argCount >= 4 {
			// Check if 4th arg is bool (flipX) or string (frameName/animation)
			arg4 := L.Get(4)
			if arg4.Type() == lua.LTString {
				frameNameOrAnim = L.CheckString(4)
				flipX = L.OptBool(5, false)
				flipY = L.OptBool(6, false)
			} else {
				flipX = L.OptBool(4, false)
				flipY = L.OptBool(5, false)
			}
		} else {
			flipX = L.OptBool(4, false)
			flipY = L.OptBool(5, false)
		}

		sprite, ok := (*spriteMapPtr)[name]
		if !ok {
			return 0 // Sprite not found, do nothing
		}

		var pixels [][]int
		var err error

		// Handle different sprite types
		switch sprite.Type {
		case cartio.SpriteTypeStatic:
			// Static sprite - use pixels directly
			pixels = sprite.Pixels
			
		case cartio.SpriteTypeFrames:
			// Frames sprite - require frame name
			if frameNameOrAnim == "" {
				L.RaiseError("frames sprite '%s' requires frame name parameter", name)
				return 0
			}
			pixels, err = sprite.GetFramePixels(frameNameOrAnim)
			if err != nil {
				L.RaiseError("frame '%s' not found in sprite '%s': %v", frameNameOrAnim, name, err)
				return 0
			}
			
		case cartio.SpriteTypeAnimation:
			// Animation sprite - use active animation state or specified animation
			if frameNameOrAnim != "" {
				// Specified animation name - use it
				anim, err := sprite.GetAnimation(frameNameOrAnim)
				if err != nil {
					L.RaiseError("animation '%s' not found in sprite '%s': %v", frameNameOrAnim, name, err)
					return 0
				}
				// Get or create animation state for this sprite
				animState, exists := animationStates[name]
				if !exists || animState.AnimationName != frameNameOrAnim {
					animState = animation.NewAnimationState(name)
					if err := animation.PlayAnimation(animState, &sprite, frameNameOrAnim); err != nil {
						L.RaiseError("failed to play animation '%s': %v", frameNameOrAnim, err)
						return 0
					}
					animationStates[name] = animState
				}
				// Update animation state with delta time (fallback - proper update should happen in engine loop)
				// This ensures frames advance even if engine loop doesn't update
				if animState.Playing && !animState.Paused {
					animation.UpdateAnimationState(animState, &sprite, 16) // ~60fps = ~16ms per frame
				}
				// Get current frame from animation state
				currentFrameName := animation.GetCurrentFrameName(animState)
				if currentFrameName == "" {
					// Fallback to first frame if no current frame
					if len(anim.FrameRefs) > 0 {
						currentFrameName = anim.FrameRefs[0]
					}
				}
				if currentFrameName != "" {
					pixels, err = sprite.GetFramePixels(currentFrameName)
					if err != nil {
						L.RaiseError("frame '%s' not found: %v", currentFrameName, err)
						return 0
					}
				} else {
					L.RaiseError("animation '%s' has no frames", frameNameOrAnim)
					return 0
				}
			} else {
				// No animation specified - check if there's an active animation state
				animState, exists := animationStates[name]
				if exists && animState.Playing && animState.Sequence != nil {
					// Update animation state (approximate delta time - this is called during draw)
					// For more accurate timing, animations should be updated in engine update loop
					// For now, we'll update with a small delta to ensure frames advance
					// Note: This is a fallback - proper update should happen in engine loop
					animation.UpdateAnimationState(animState, &sprite, 16) // ~60fps = ~16ms per frame
					
					currentFrameName := animation.GetCurrentFrameName(animState)
					if currentFrameName != "" {
						pixels, err = sprite.GetFramePixels(currentFrameName)
						if err != nil {
							// Fallback to first frame if current frame not found
							if len(animState.Sequence.FrameRefs) > 0 {
								pixels, err = sprite.GetFramePixels(animState.Sequence.FrameRefs[0])
							}
						}
					}
				}
				// If no active animation or frame not found, use first frame from first animation
				if pixels == nil {
					if len(sprite.Animations) > 0 && len(sprite.Animations[0].FrameRefs) > 0 {
						pixels, err = sprite.GetFramePixels(sprite.Animations[0].FrameRefs[0])
					}
				}
				if pixels == nil || err != nil {
					L.RaiseError("sprite '%s' has no active animation or default frame", name)
					return 0
				}
			}
			
		default:
			// Default to static behavior for backward compatibility
			if sprite.Pixels != nil && len(sprite.Pixels) > 0 {
				pixels = sprite.Pixels
			} else {
				return 0 // No pixels to draw
			}
		}

		if pixels == nil {
			return 0 // No pixels to draw
		}

		// Draw sprite pixels
		for sy := 0; sy < len(pixels); sy++ {
			if sy >= sprite.Height {
				break
			}
			for sx := 0; sx < len(pixels[sy]); sx++ {
				if sx >= sprite.Width {
					break
				}
				// Calculate source coordinates with flipping
				srcX := sx
				srcY := sy
				if flipX {
					srcX = sprite.Width - 1 - sx
				}
				if flipY {
					srcY = sprite.Height - 1 - sy
				}

				// Ensure source coordinates are within bounds
				if srcX < 0 || srcY < 0 || srcX >= len(pixels[srcY]) || srcY >= len(pixels) {
					continue
				}

				colorIdx := pixels[srcY][srcX]
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
		isUI := L.OptBool(4, true) // Default to true for UI sprite

		// Validate sprite size using new rules
		if err := cartio.ValidateSpriteSize(width, height, isUI); err != nil {
			L.RaiseError("%s", err.Error())
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
			Type:         cartio.SpriteTypeStatic, // Default to static sprite
			UseCollision: false,
			MountPoints:  []cartio.MountPoint{},
			IsUI:         isUI, // Use provided value or default true
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
		tbl.RawSetString("isUI", lua.LBool(isUI))
		tbl.RawSetString("lifetime", lua.LNumber(0))
		tbl.RawSetString("maxSpawn", lua.LNumber(0))
		tbl.RawSetString("mountPoints", L.NewTable())

		L.Push(tbl)
		return 1
	}))

	// rf.newSpriteFrames(name, width, height, [isUI]) -> creates a multi-frame sprite
	L.SetField(rf, "newSpriteFrames", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		width := L.CheckInt(2)
		height := L.CheckInt(3)
		isUI := L.OptBool(4, true) // Default to true for UI sprite

		// Validate sprite size
		if err := cartio.ValidateSpriteSize(width, height, isUI); err != nil {
			L.RaiseError("%s", err.Error())
			return 0
		}

		// Initialize pixels array (for frame reference dimensions)
		pixels := make([][]int, height)
		for y := range pixels {
			pixels[y] = make([]int, width)
			for x := range pixels[y] {
				pixels[y][x] = -1 // Transparent
			}
		}

		// Create new frames sprite
		newSprite := cartio.SpriteData{
			Width:        width,
			Height:       height,
			Type:         cartio.SpriteTypeFrames,
			Frames:       []cartio.SpriteFrame{},
			UseCollision: false,
			MountPoints:  []cartio.MountPoint{},
			IsUI:         isUI,
			Lifetime:     0,
			MaxSpawn:     0,
		}

		// Add to sprite map
		(*spriteMapPtr)[name] = newSprite

		// Return sprite table
		tbl := L.NewTable()
		tbl.RawSetString("width", lua.LNumber(width))
		tbl.RawSetString("height", lua.LNumber(height))
		tbl.RawSetString("type", lua.LString("frames"))
		tbl.RawSetString("isUI", lua.LBool(isUI))

		L.Push(tbl)
		return 1
	}))

	// rf.addSpriteFrame(spriteName, frameName, pixels) -> adds a frame to a frames sprite
	L.SetField(rf, "addSpriteFrame", L.NewFunction(func(L *lua.LState) int {
		spriteName := L.CheckString(1)
		frameName := L.CheckString(2)
		pixelsTbl := L.CheckTable(3)

		// Validate frame name
		if err := cartio.ValidateFrameName(frameName); err != nil {
			L.RaiseError("invalid frame name: %v", err)
			return 0
		}

		sprite, ok := (*spriteMapPtr)[spriteName]
		if !ok {
			L.RaiseError("sprite '%s' not found", spriteName)
			return 0
		}

		// Ensure sprite is frames or animation type
		if sprite.Type != cartio.SpriteTypeFrames && sprite.Type != cartio.SpriteTypeAnimation {
			L.RaiseError("sprite '%s' is not a frames or animation sprite (type: %s)", spriteName, sprite.Type)
			return 0
		}

		// Check for duplicate frame name
		for _, frame := range sprite.Frames {
			if frame.Name == frameName {
				L.RaiseError("frame '%s' already exists in sprite '%s'", frameName, spriteName)
				return 0
			}
		}

		// Convert Lua table to pixels array
		pixels := make([][]int, sprite.Height)
		pixelsTbl.ForEach(func(key, value lua.LValue) {
			rowIdx := int(lua.LVAsNumber(key)) - 1 // Lua is 1-indexed
			if rowIdx >= 0 && rowIdx < sprite.Height {
				rowTbl, ok := value.(*lua.LTable)
				if !ok {
					return
				}
				pixels[rowIdx] = make([]int, sprite.Width)
				rowTbl.ForEach(func(colKey, colValue lua.LValue) {
					colIdx := int(lua.LVAsNumber(colKey)) - 1
					if colIdx >= 0 && colIdx < sprite.Width {
						pixels[rowIdx][colIdx] = int(lua.LVAsNumber(colValue))
					}
				})
			}
		})

		// Validate pixels dimensions
		if len(pixels) != sprite.Height {
			L.RaiseError("pixels array height (%d) does not match sprite height (%d)", len(pixels), sprite.Height)
			return 0
		}
		for i, row := range pixels {
			if len(row) != sprite.Width {
				L.RaiseError("row %d width (%d) does not match sprite width (%d)", i, len(row), sprite.Width)
				return 0
			}
		}

		// Add frame
		newFrame := cartio.SpriteFrame{
			Name:     frameName,
			Pixels:   pixels,
			Duration: 100, // Default 100ms
		}
		sprite.Frames = append(sprite.Frames, newFrame)
		(*spriteMapPtr)[spriteName] = sprite

		return 0
	}))

	// rf.addSpriteAnimation(spriteName, animName, frameRefs, [speed], [loop], [loopType]) -> adds animation to sprite
	L.SetField(rf, "addSpriteAnimation", L.NewFunction(func(L *lua.LState) int {
		spriteName := L.CheckString(1)
		animName := L.CheckString(2)
		frameRefsTbl := L.CheckTable(3)
		speed := L.OptNumber(4, 1.0)
		loop := L.OptBool(5, true)
		loopType := L.OptString(6, "forward")

		// Validate animation name
		if err := cartio.ValidateFrameName(animName); err != nil {
			L.RaiseError("invalid animation name: %v", err)
			return 0
		}

		// Validate loop type
		if loopType != "forward" && loopType != "reverse" && loopType != "pingpong" {
			L.RaiseError("invalid loopType '%s', must be 'forward', 'reverse', or 'pingpong'", loopType)
			return 0
		}

		sprite, ok := (*spriteMapPtr)[spriteName]
		if !ok {
			L.RaiseError("sprite '%s' not found", spriteName)
			return 0
		}

		// Ensure sprite is animation type
		if sprite.Type != cartio.SpriteTypeAnimation {
			L.RaiseError("sprite '%s' is not an animation sprite (type: %s)", spriteName, sprite.Type)
			return 0
		}

		// Check for duplicate animation name
		for _, anim := range sprite.Animations {
			if anim.Name == animName {
				L.RaiseError("animation '%s' already exists in sprite '%s'", animName, spriteName)
				return 0
			}
		}

		// Convert frame refs table to slice
		var frameRefs []string
		frameRefsTbl.ForEach(func(key, value lua.LValue) {
			if str, ok := value.(lua.LString); ok {
				frameRefs = append(frameRefs, string(str))
			}
		})

		if len(frameRefs) == 0 {
			L.RaiseError("animation must have at least one frame reference")
			return 0
		}

		// Validate frame references exist
		frameNames := make(map[string]bool)
		for _, frame := range sprite.Frames {
			frameNames[frame.Name] = true
		}
		for _, frameRef := range frameRefs {
			if !frameNames[frameRef] {
				L.RaiseError("frame reference '%s' not found in sprite '%s'", frameRef, spriteName)
				return 0
			}
		}

		// Validate speed
		speedFloat := float64(speed)
		if speedFloat <= 0 {
			speedFloat = 1.0
		}

		// Add animation
		newAnim := cartio.AnimationSequence{
			Name:      animName,
			FrameRefs: frameRefs,
			Speed:     speedFloat,
			Loop:      loop,
			LoopType:  loopType,
		}
		sprite.Animations = append(sprite.Animations, newAnim)
		(*spriteMapPtr)[spriteName] = sprite

		return 0
	}))

	// Helper function to update animation states for all instances of a sprite
	updateSpriteAnimations := func(spriteName string, updateFunc func(*animation.AnimationState)) {
		// Update direct animation state (for non-pooled sprites)
		if animState, exists := animationStates[spriteName]; exists {
			updateFunc(animState)
		}

		// Update all active pool instances
		poolManager.ForEachActiveInstance(spriteName, func(instance *spritepool.SpriteInstance) {
			if instance.AnimationState != nil {
				updateFunc(instance.AnimationState)
			}
		})
	}

	// rf.playAnimation(spriteName, animationName) -> starts animation for all instances
	L.SetField(rf, "playAnimation", L.NewFunction(func(L *lua.LState) int {
		spriteName := L.CheckString(1)
		animationName := L.CheckString(2)

		sprite, ok := (*spriteMapPtr)[spriteName]
		if !ok {
			L.RaiseError("sprite '%s' not found", spriteName)
			return 0
		}

		if sprite.Type != cartio.SpriteTypeAnimation {
			L.RaiseError("sprite '%s' is not an animation sprite", spriteName)
			return 0
		}

		updateSpriteAnimations(spriteName, func(animState *animation.AnimationState) {
			if err := animation.PlayAnimation(animState, &sprite, animationName); err != nil {
				// Log error but don't fail - individual instances might fail
				if devModePtr != nil {
					devModePtr.AddDebugLog(fmt.Sprintf("Failed to play animation '%s' on sprite '%s': %v", animationName, spriteName, err))
				}
			}
		})

		// Also ensure direct state exists for rf.spr() calls
		animState, exists := animationStates[spriteName]
		if !exists {
			animState = animation.NewAnimationState(spriteName)
			animationStates[spriteName] = animState
		}
		if err := animation.PlayAnimation(animState, &sprite, animationName); err != nil {
			L.RaiseError("failed to play animation '%s': %v", animationName, err)
			return 0
		}

		return 0
	}))

	// rf.pauseAnimation(spriteName) -> pauses animation for all instances
	L.SetField(rf, "pauseAnimation", L.NewFunction(func(L *lua.LState) int {
		spriteName := L.CheckString(1)

		updateSpriteAnimations(spriteName, func(animState *animation.AnimationState) {
			animation.PauseAnimation(animState)
		})

		return 0
	}))

	// rf.resumeAnimation(spriteName) -> resumes paused animation for all instances
	L.SetField(rf, "resumeAnimation", L.NewFunction(func(L *lua.LState) int {
		spriteName := L.CheckString(1)

		updateSpriteAnimations(spriteName, func(animState *animation.AnimationState) {
			animation.ResumeAnimation(animState)
		})

		return 0
	}))

	// rf.stopAnimation(spriteName) -> stops animation for all instances
	L.SetField(rf, "stopAnimation", L.NewFunction(func(L *lua.LState) int {
		spriteName := L.CheckString(1)

		updateSpriteAnimations(spriteName, func(animState *animation.AnimationState) {
			animation.StopAnimation(animState)
		})

		return 0
	}))

	// rf.setAnimationSpeed(spriteName, speed) -> sets speed multiplier for all instances
	L.SetField(rf, "setAnimationSpeed", L.NewFunction(func(L *lua.LState) int {
		spriteName := L.CheckString(1)
		speed := L.CheckNumber(2)

		updateSpriteAnimations(spriteName, func(animState *animation.AnimationState) {
			animation.SetAnimationSpeed(animState, float64(speed))
		})

		return 0
	}))

	// rf.setAnimationFrame(spriteName, frameIndex) -> sets frame index for all instances
	L.SetField(rf, "setAnimationFrame", L.NewFunction(func(L *lua.LState) int {
		spriteName := L.CheckString(1)
		frameIndex := L.CheckInt(2)

		updateSpriteAnimations(spriteName, func(animState *animation.AnimationState) {
			if err := animation.SetAnimationFrame(animState, frameIndex); err != nil {
				// Log error but don't fail - individual instances might fail
				if devModePtr != nil {
					devModePtr.AddDebugLog(fmt.Sprintf("Failed to set frame %d on sprite '%s': %v", frameIndex, spriteName, err))
				}
			}
		})

		return 0
	}))

	// rf.getAnimationFrame(spriteName) -> returns current frame index (from first available instance)
	L.SetField(rf, "getAnimationFrame", L.NewFunction(func(L *lua.LState) int {
		spriteName := L.CheckString(1)

		// Check direct animation state first
		if animState, exists := animationStates[spriteName]; exists {
			frame := animation.GetAnimationFrame(animState)
			L.Push(lua.LNumber(frame))
			return 1
		}

		// Check pool instances
		var foundFrame int = -1
		poolManager.ForEachActiveInstance(spriteName, func(instance *spritepool.SpriteInstance) {
			if instance.AnimationState != nil && foundFrame == -1 {
				foundFrame = animation.GetAnimationFrame(instance.AnimationState)
			}
		})
		if foundFrame != -1 {
			L.Push(lua.LNumber(foundFrame))
			return 1
		}

		// No active animation
		L.Push(lua.LNumber(-1))
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
				c := colorByIndexRemapped(tileIndex % 64) // Use tile index as color (0-63)
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

	// Tilemap rendering: rf.drawTilemap(mapName, offsetX, offsetY)
	// Draws a tilemap at the specified offset
	tilemapsMapPtr := &tilemapsMap
	L.SetField(rf, "drawTilemap", L.NewFunction(func(L *lua.LState) int {
		mapName := L.CheckString(1)
		offsetX := L.OptInt(2, 0)
		offsetY := L.OptInt(3, 0)

		// Get tilemap
		tilemapData, exists := (*tilemapsMapPtr)[mapName]
		if !exists {
			debugLog("drawTilemap: Tilemap '%s' not found", mapName)
			// List available tilemaps for debugging
			availableMaps := make([]string, 0, len(*tilemapsMapPtr))
			for name := range *tilemapsMapPtr {
				availableMaps = append(availableMaps, name)
			}
			debugLog("drawTilemap: Available tilemaps: %v", availableMaps)
			return 0
		}
		
		// Only log tilemap info once (smart logging will deduplicate)
		debugLog("drawTilemap: Found tilemap '%s', IsISO=%v", mapName, tilemapData.IsISO)
		
		// Check if tileset is empty
		if len(tilemapData.Tileset) == 0 {
			debugLog("drawTilemap: ERROR - Tileset is empty for tilemap '%s'", mapName)
			return 0
		}
		
		// Check if tilemap has any tiles
		if len(tilemapData.Tiles) == 0 || len(tilemapData.Tiles[0]) == 0 {
			debugLog("drawTilemap: ERROR - Tilemap '%s' is empty", mapName)
			return 0
		}

		// Get first tile to determine tile dimensions (for isometric calculations)
		var tileWidth, tileHeight int
		for _, row := range tilemapData.Tiles {
			for _, tileName := range row {
				if tileName != "" {
					if tile, exists := tilemapData.Tileset[tileName]; exists {
						tileWidth = tile.Width
						tileHeight = tile.Height
						break
					}
				}
			}
			if tileWidth > 0 {
				break
			}
		}
		
		// If we couldn't find a tile, can't render
		if tileWidth == 0 || tileHeight == 0 {
			debugLog("drawTilemap: ERROR - Could not determine tile dimensions (width=%d, height=%d)", tileWidth, tileHeight)
			return 0
		}

		// Draw each tile in the map (sorted back-to-front for isometric)
		type tileRenderInfo struct {
			gridX    int
			gridY    int
			tileName string
			screenX  int
			screenY  int
		}

		renderQueue := make([]tileRenderInfo, 0)

		// Build render queue
		for mapY, row := range tilemapData.Tiles {
			for mapX, tileName := range row {
				if tileName == "" {
					continue // Empty tile
				}

				// Get tile from tileset
				tile, exists := tilemapData.Tileset[tileName]
				if !exists {
					continue // Tile not found in tileset
				}

				// Calculate screen position
				var screenX, screenY int
				if tilemapData.IsISO {
					// Isometric positioning using basis vectors
					// For 32×24 isometric tiles (2:1 ratio):
					// - TileWidth: 32 (full tile width)
					// - DiamondHeight: 16 (height of top diamond face, not full tile height)
					// - Full tile height is 24 (16 diamond + 8 side faces)
					// Basis vectors:
					// i_hat: (tileWidth/2, diamondHeight/2) = (16, 8) for 32×24 tiles
					// j_hat: (-tileWidth/2, diamondHeight/2) = (-16, 8) for 32×24 tiles
					// screen = (gridX * i_hat) + (gridY * j_hat) + offset
					
					// Calculate diamond height (for 2:1 isometric ratio: width/2)
					diamondHeight := tileWidth / 2  // 32 / 2 = 16 for 32×24 tiles
					
					// Basis vectors
					iHatX := float64(tileWidth) / 2.0      // 16
					iHatY := float64(diamondHeight) / 2.0  // 8
					jHatX := -float64(tileWidth) / 2.0     // -16
					jHatY := float64(diamondHeight) / 2.0  // 8

					// Matrix transformation: screen = (gridX * i_hat) + (gridY * j_hat)
					screenX = int(float64(mapX)*iHatX + float64(mapY)*jHatX)
					screenY = int(float64(mapX)*iHatY + float64(mapY)*jHatY)

					// Apply origin offset
					// For isometric, anchor is at the top center of the diamond
					screenX += offsetX - (tileWidth / 2)
					screenY += offsetY

					// Add to render queue for depth sorting
					renderQueue = append(renderQueue, tileRenderInfo{
						gridX:    mapX,
						gridY:    mapY,
						tileName: tileName,
						screenX:  screenX,
						screenY:  screenY,
					})
				} else {
					// Normal orthogonal positioning: simple grid
					screenX = offsetX + mapX*tile.Width
					screenY = offsetY + mapY*tile.Height

					// Add to render queue (no sorting needed for orthogonal)
					renderQueue = append(renderQueue, tileRenderInfo{
						gridX:    mapX,
						gridY:    mapY,
						tileName: tileName,
						screenX:  screenX,
						screenY:  screenY,
					})
				}
			}
		}

		// Sort isometric tiles back-to-front (gridX + gridY)
		if tilemapData.IsISO {
			sort.Slice(renderQueue, func(i, j int) bool {
				return renderQueue[i].gridX+renderQueue[i].gridY < renderQueue[j].gridX+renderQueue[j].gridY
			})
		}

		// Render all tiles
		tilesDrawn := 0
		pixelsDrawn := 0
		debugLog("drawTilemap: Rendering %d tiles (IsISO=%v)", len(renderQueue), tilemapData.IsISO)
		for _, tileInfo := range renderQueue {
			tile, exists := tilemapData.Tileset[tileInfo.tileName]
			if !exists {
				continue
			}

			// Draw tile using rf.spr (treat tiles as sprites)
			// For now, we'll draw the tile directly pixel by pixel
			// Get tile pixels (handle static/frames/animation)
			var pixels [][]int
			switch tile.Type {
			case cartio.SpriteTypeStatic:
				pixels = tile.Pixels
			case cartio.SpriteTypeFrames, cartio.SpriteTypeAnimation:
				// For frames/animation, use first frame
				if len(tile.Frames) > 0 {
					pixels = tile.Frames[0].Pixels
				} else {
					continue
				}
			default:
				continue
			}

			// Apply tile variation (rotation/flipping) for normal tiles only
			// Isometric tiles cannot be rotated/flipped (would break 3D appearance)
			// For seamless tiles, we only apply variations when adjacent tiles are different types
			var drawWidth, drawHeight int = tile.Width, tile.Height
			if !tilemapData.IsISO {
				// Check if adjacent tiles are the same type to maintain seamlessness
				shouldVary := true
				adjacentTiles := []string{
					"", // top (will check if exists)
					"", // right
					"", // bottom
					"", // left
				}
				
				// Get adjacent tile names (if they exist)
				if tileInfo.gridY > 0 && tileInfo.gridY-1 < len(tilemapData.Tiles) && tileInfo.gridX < len(tilemapData.Tiles[tileInfo.gridY-1]) {
					adjacentTiles[0] = tilemapData.Tiles[tileInfo.gridY-1][tileInfo.gridX] // top
				}
				if tileInfo.gridX+1 < len(tilemapData.Tiles[tileInfo.gridY]) {
					adjacentTiles[1] = tilemapData.Tiles[tileInfo.gridY][tileInfo.gridX+1] // right
				}
				if tileInfo.gridY+1 < len(tilemapData.Tiles) && tileInfo.gridX < len(tilemapData.Tiles[tileInfo.gridY+1]) {
					adjacentTiles[2] = tilemapData.Tiles[tileInfo.gridY+1][tileInfo.gridX] // bottom
				}
				if tileInfo.gridX > 0 {
					adjacentTiles[3] = tilemapData.Tiles[tileInfo.gridY][tileInfo.gridX-1] // left
				}
				
				// If any adjacent tile is the same type, don't apply variation to maintain seamlessness
				for _, adjTileName := range adjacentTiles {
					if adjTileName == tileInfo.tileName {
						shouldVary = false
						break
					}
				}
				
				if shouldVary {
					variation := cartio.GetTileVariation(tilemapData.Seed, tileInfo.gridX, tileInfo.gridY, false)
					var transformedPixels [][]int
					transformedPixels, drawWidth, drawHeight = cartio.ApplyVariation(pixels, tile.Width, tile.Height, variation)
					pixels = transformedPixels
				}
			}

			// For isometric tiles: enforce DIRT sides (bottom left and bottom right)
			// The side faces are always DIRT, cannot be any other type
			if tilemapData.IsISO {
				// Isometric tiles have structure:
				// - Top diamond: Y=0 to Y=tileHeight (16px for 32x24 tiles)
				// - Side faces: Y=tileHeight to Y=tileHeight+sideHeight (16 to 24px)
				//   - Left side: X=0 to X=tileWidth/2 (0 to 16px)
				//   - Right side: X=tileWidth/2 to X=tileWidth (16 to 32px)
				// Replace side face pixels with DIRT tile pixels if available
				dirtTile, hasDirt := tilemapData.Tileset["dirt"]
				if hasDirt && len(pixels) > tile.Height {
					// Get DIRT tile pixels (use first frame if multi-frame)
					var dirtPixels [][]int
					switch dirtTile.Type {
					case cartio.SpriteTypeStatic:
						dirtPixels = dirtTile.Pixels
					case cartio.SpriteTypeFrames, cartio.SpriteTypeAnimation:
						if len(dirtTile.Frames) > 0 {
							dirtPixels = dirtTile.Frames[0].Pixels
						}
					}
					
					if dirtPixels != nil && len(dirtPixels) > 0 {
						// Calculate side face dimensions
						// For 32x24 isometric tiles: tileHeight=16, sideHeight=8
						// Side faces start at Y=tileHeight (16) and go to Y=tileHeight+sideHeight (24)
						sideStartY := tile.Height // 16 for 32x24 tiles
						sideHeight := len(pixels) - tile.Height // 8 for 32x24 tiles
						sideWidth := tile.Width / 2 // 16 for 32x24 tiles
						
						// Replace left side face (X=0 to X=sideWidth, Y=sideStartY to Y=sideStartY+sideHeight)
						for y := sideStartY; y < len(pixels) && y < sideStartY+sideHeight; y++ {
							dirtY := y - sideStartY
							if dirtY < len(dirtPixels) {
								for x := 0; x < sideWidth && x < len(pixels[y]) && x < len(dirtPixels[dirtY]); x++ {
									// Replace with DIRT pixel (use modulo to handle size differences)
									dirtX := x % len(dirtPixels[dirtY])
									pixels[y][x] = dirtPixels[dirtY][dirtX]
								}
							}
						}
						
						// Replace right side face (X=sideWidth to X=tileWidth, Y=sideStartY to Y=sideStartY+sideHeight)
						for y := sideStartY; y < len(pixels) && y < sideStartY+sideHeight; y++ {
							dirtY := y - sideStartY
							if dirtY < len(dirtPixels) {
								for x := sideWidth; x < tile.Width && x < len(pixels[y]); x++ {
									// Replace with DIRT pixel (use modulo to handle size differences)
									dirtX := (x - sideWidth) % len(dirtPixels[dirtY])
									pixels[y][x] = dirtPixels[dirtY][dirtX]
								}
							}
						}
					}
				}
			}

			// Draw tile pixels
			// Validate pixels array dimensions
			if len(pixels) != drawHeight {
				debugLog("drawTilemap: Tile '%s' has invalid pixel height: %d (expected %d)", tileInfo.tileName, len(pixels), drawHeight)
				continue // Skip invalid tile
			}
			tilePixelsDrawn := 0
			for ty, pixelRow := range pixels {
				if len(pixelRow) != drawWidth {
					debugLog("drawTilemap: Tile '%s' row %d has invalid width: %d (expected %d)", tileInfo.tileName, ty, len(pixelRow), drawWidth)
					continue // Skip invalid row
				}
				for tx, colorIdx := range pixelRow {
					// Skip transparent pixels (colorIdx < 0)
					if colorIdx < 0 {
						continue
					}
					// Valid color index is 0-63 (16 built-in + 48 game palette)
					if colorIdx >= 0 && colorIdx < 64 {
						px := tileInfo.screenX + tx
						py := tileInfo.screenY + ty
						// Only draw if within screen bounds
						if px >= 0 && px < 480 && py >= 0 && py < 270 {
							c := colorByIndexRemapped(colorIdx)
							r.PSet(px, py, color.RGBA{c[0], c[1], c[2], c[3]})
							tilePixelsDrawn++
							pixelsDrawn++
						}
					}
				}
			}
			if tilePixelsDrawn > 0 {
				tilesDrawn++
			}
		}
		
		// Only log if there's an issue (no tiles drawn) or very rarely (smart logging will throttle)
		if tilesDrawn == 0 {
			debugLog("drawTilemap: WARNING - No tiles drawn for tilemap '%s'", mapName)
		} else {
			debugLog("drawTilemap: Successfully rendered %d tiles (%d pixels) for '%s'", tilesDrawn, pixelsDrawn, mapName)
		}

		return 0
	}))
}
