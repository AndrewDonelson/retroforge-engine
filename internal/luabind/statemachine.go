package luabind

import (
	"github.com/AndrewDonelson/retroforge-engine/internal/gamestate"
	"github.com/AndrewDonelson/retroforge-engine/internal/statemachine"
	lua "github.com/yuin/gopher-lua"
)

// RegisterStateMachine attaches game.* state machine functions to the Lua state
func RegisterStateMachine(L *lua.LState, gsm *gamestate.GameStateMachine) {
	game := L.NewTable()
	L.SetGlobal("game", game)

	// Store state machine in userdata
	ud := L.NewUserData()
	ud.Value = gsm
	L.SetMetatable(ud, L.GetTypeMetatable("GameStateMachine"))

	// game.registerState(name, stateTable)
	L.SetField(game, "registerState", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		stateTable := L.CheckTable(2)

		// Extract callbacks from the table
		callbacks := statemachine.LuaCallbacks{}

		// Initialize callback
		if fn := L.GetField(stateTable, "initialize"); fn != lua.LNil {
			if lfn, ok := fn.(*lua.LFunction); ok {
				callbacks.Initialize = func(sm *statemachine.StateMachine) error {
					L.Push(lfn)
					// Pass nil for sm parameter - Lua code uses game.* functions instead
					L.Push(lua.LNil)
					err := L.PCall(1, 0, nil)
					if err != nil {
						return err
					}
					return nil
				}
			}
		}

		// Enter callback
		if fn := L.GetField(stateTable, "enter"); fn != lua.LNil {
			if lfn, ok := fn.(*lua.LFunction); ok {
				callbacks.Enter = func(sm *statemachine.StateMachine) {
					L.Push(lfn)
					// Pass nil for sm parameter - Lua code uses game.* functions instead
					L.Push(lua.LNil)
					L.PCall(1, 0, nil)
				}
			}
		}

		// HandleInput callback
		if fn := L.GetField(stateTable, "handleInput"); fn != lua.LNil {
			if lfn, ok := fn.(*lua.LFunction); ok {
				callbacks.HandleInput = func(sm *statemachine.StateMachine) {
					L.Push(lfn)
					// Pass nil for sm parameter - Lua code uses game.* functions instead
					L.Push(lua.LNil)
					L.PCall(1, 0, nil)
				}
			}
		}

		// Update callback
		if fn := L.GetField(stateTable, "update"); fn != lua.LNil {
			if lfn, ok := fn.(*lua.LFunction); ok {
				callbacks.Update = func(dt float64) {
					L.Push(lfn)
					L.Push(lua.LNumber(dt))
					L.PCall(1, 0, nil)
				}
			}
		}

		// Draw callback
		if fn := L.GetField(stateTable, "draw"); fn != lua.LNil {
			if lfn, ok := fn.(*lua.LFunction); ok {
				callbacks.Draw = func() {
					L.Push(lfn)
					L.PCall(0, 0, nil)
				}
			}
		}

		// Exit callback
		if fn := L.GetField(stateTable, "exit"); fn != lua.LNil {
			if lfn, ok := fn.(*lua.LFunction); ok {
				callbacks.Exit = func(sm *statemachine.StateMachine) {
					L.Push(lfn)
					// Pass nil for sm parameter - Lua code uses game.* functions instead
					L.Push(lua.LNil)
					L.PCall(1, 0, nil)
				}
			}
		}

		// Shutdown callback
		if fn := L.GetField(stateTable, "shutdown"); fn != lua.LNil {
			if lfn, ok := fn.(*lua.LFunction); ok {
				callbacks.Shutdown = func() {
					L.Push(lfn)
					L.PCall(0, 0, nil)
				}
			}
		}

		err := gsm.RegisterState(name, callbacks)
		if err != nil {
			L.RaiseError("failed to register state: %v", err)
			return 0
		}

		return 0
	}))

	// game.unregisterState(name)
	L.SetField(game, "unregisterState", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		err := gsm.UnregisterState(name)
		if err != nil {
			L.RaiseError("failed to unregister state: %v", err)
			return 0
		}
		return 0
	}))

	// game.changeState(name)
	L.SetField(game, "changeState", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		err := gsm.ChangeState(name)
		if err != nil {
			L.RaiseError("failed to change state: %v", err)
			return 0
		}
		return 0
	}))

	// game.pushState(name)
	L.SetField(game, "pushState", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		err := gsm.PushState(name)
		if err != nil {
			L.RaiseError("failed to push state: %v", err)
			return 0
		}
		return 0
	}))

	// game.popState()
	L.SetField(game, "popState", L.NewFunction(func(L *lua.LState) int {
		err := gsm.PopState()
		if err != nil {
			L.RaiseError("failed to pop state: %v", err)
			return 0
		}
		return 0
	}))

	// game.popAllStates()
	L.SetField(game, "popAllStates", L.NewFunction(func(L *lua.LState) int {
		gsm.PopAllStates()
		return 0
	}))

	// Context Management

	// game.setContext(key, value)
	L.SetField(game, "setContext", L.NewFunction(func(L *lua.LState) int {
		key := L.CheckString(1)
		value := L.CheckAny(2)

		// Convert Lua value to Go interface{}
		var goValue interface{}
		switch v := value.(type) {
		case lua.LString:
			goValue = string(v)
		case lua.LNumber:
			goValue = float64(v)
		case lua.LBool:
			goValue = bool(v)
		case *lua.LTable:
			// Convert table to map
			goValue = tableToMap(L, v)
		default:
			goValue = nil
		}

		gsm.SetContext(key, goValue)
		return 0
	}))

	// game.getContext(key) -> value or nil
	L.SetField(game, "getContext", L.NewFunction(func(L *lua.LState) int {
		key := L.CheckString(1)
		value, exists := gsm.GetContext(key)
		if !exists {
			L.Push(lua.LNil)
			return 1
		}

		// Convert Go value to Lua value
		L.Push(goValueToLua(L, value))
		return 1
	}))

	// game.hasContext(key) -> bool
	L.SetField(game, "hasContext", L.NewFunction(func(L *lua.LState) int {
		key := L.CheckString(1)
		exists := gsm.HasContext(key)
		L.Push(lua.LBool(exists))
		return 1
	}))

	// game.clearContext(key)
	L.SetField(game, "clearContext", L.NewFunction(func(L *lua.LState) int {
		key := L.CheckString(1)
		gsm.ClearContext(key)
		return 0
	}))

	// game.clearAllContext()
	L.SetField(game, "clearAllContext", L.NewFunction(func(L *lua.LState) int {
		gsm.ClearAllContext()
		return 0
	}))

	// Credits API

	// game.addCredit(category, name, role)
	L.SetField(game, "addCredit", L.NewFunction(func(L *lua.LState) int {
		category := L.CheckString(1)
		name := L.CheckString(2)
		role := L.CheckString(3)
		gsm.AddCreditEntry(category, name, role)
		return 0
	}))

	// Control

	// game.exit() - Transition to credits then exit
	L.SetField(game, "exit", L.NewFunction(func(L *lua.LState) int {
		err := gsm.Exit()
		if err != nil {
			L.RaiseError("failed to exit: %v", err)
			return 0
		}
		return 0
	}))

	// Utility

	// game.drawPreviousState() - Draw state underneath in stack
	L.SetField(game, "drawPreviousState", L.NewFunction(func(L *lua.LState) int {
		gsm.DrawPreviousState()
		return 0
	}))

	// game.getStackDepth() -> number
	L.SetField(game, "getStackDepth", L.NewFunction(func(L *lua.LState) int {
		depth := gsm.GetStackDepth()
		L.Push(lua.LNumber(depth))
		return 1
	}))
}

// Helper function to convert Lua table to Go map
func tableToMap(L *lua.LState, tbl *lua.LTable) map[string]interface{} {
	result := make(map[string]interface{})
	tbl.ForEach(func(key lua.LValue, value lua.LValue) {
		keyStr := key.String()
		var goValue interface{}
		switch v := value.(type) {
		case lua.LString:
			goValue = string(v)
		case lua.LNumber:
			goValue = float64(v)
		case lua.LBool:
			goValue = bool(v)
		case *lua.LTable:
			goValue = tableToMap(L, v)
		default:
			goValue = nil
		}
		result[keyStr] = goValue
	})
	return result
}

// Helper function to convert Go value to Lua value
func goValueToLua(L *lua.LState, value interface{}) lua.LValue {
	switch v := value.(type) {
	case string:
		return lua.LString(v)
	case float64:
		return lua.LNumber(v)
	case int:
		return lua.LNumber(v)
	case int64:
		return lua.LNumber(v)
	case bool:
		return lua.LBool(v)
	case map[string]interface{}:
		tbl := L.NewTable()
		for k, val := range v {
			L.SetField(tbl, k, goValueToLua(L, val))
		}
		return tbl
	case []interface{}:
		tbl := L.NewTable()
		for i, val := range v {
			L.RawSetInt(tbl, i+1, goValueToLua(L, val))
		}
		return tbl
	default:
		return lua.LNil
	}
}
