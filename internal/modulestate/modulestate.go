package modulestate

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/AndrewDonelson/retroforge-engine/internal/gamestate"
	"github.com/AndrewDonelson/retroforge-engine/internal/statemachine"
	lua "github.com/yuin/gopher-lua"
)

// FileReader interface for reading module files
type FileReader interface {
	ReadFile(path string) ([]byte, error)
}

// ModuleLoader handles loading and registering state modules
type ModuleLoader struct {
	L             *lua.LState
	gsm           *gamestate.GameStateMachine
	fileReader    FileReader
	basePath      string                 // Base path for relative file lookups
	loadedModules map[string]*lua.LTable // Track loaded modules for persistence
}

// NewModuleLoader creates a new module loader
func NewModuleLoader(L *lua.LState, gsm *gamestate.GameStateMachine, fileReader FileReader, basePath string) *ModuleLoader {
	return &ModuleLoader{
		L:             L,
		gsm:           gsm,
		fileReader:    fileReader,
		basePath:      basePath,
		loadedModules: make(map[string]*lua.LTable),
	}
}

// ExtractStateName extracts the state name from a filename
// Pattern: {name}_state.lua or {name}.lua → state name is {name}
func ExtractStateName(filename string) string {
	// Remove path and extension
	name := filepath.Base(filename)
	name = strings.TrimSuffix(name, ".lua")

	// Remove _state suffix if present
	if strings.HasSuffix(name, "_state") {
		name = strings.TrimSuffix(name, "_state")
	}

	return name
}

// RequiredFunctions is the list of functions that must be defined in a state module
var RequiredFunctions = []string{"_INIT", "_UPDATE", "_DRAW", "_HANDLE_INPUT", "_DONE"}

// OptionalFunctions are functions that may be defined
var OptionalFunctions = []string{"_ENTER", "_EXIT"}

// ImportModule loads a Lua module file and registers it as a game state
func (ml *ModuleLoader) ImportModule(filename string) (string, error) {
	// Read the file (fileReader handles basePath internally if needed)
	content, err := ml.fileReader.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read module file '%s': %w", filename, err)
	}

	// Extract state name
	stateName := ExtractStateName(filename)
	if stateName == "" {
		return "", fmt.Errorf("invalid state name extracted from filename '%s'", filename)
	}

	// Check if already loaded
	if _, exists := ml.loadedModules[stateName]; exists {
		return stateName, nil // Already loaded, skip
	}

	// Create isolated environment for the module
	moduleEnv := ml.createModuleEnvironment(stateName)

	// Load and execute the Lua code in the module environment
	// Load the code as a function
	chunk, err := ml.L.LoadString(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to compile module '%s': %w", filename, err)
	}

	// Set the environment on the chunk
	// This ensures that when the chunk executes, any globals (including functions)
	// it defines will be stored in moduleEnv
	ml.L.SetFEnv(chunk, moduleEnv)

	// Execute the chunk (which defines functions in the environment)
	ml.L.Push(chunk)
	err = ml.L.PCall(0, lua.MultRet, nil)
	if err != nil {
		return "", fmt.Errorf("failed to execute module '%s': %w", filename, err)
	}

	// Validate required functions
	if err := ml.validateModule(moduleEnv, filename); err != nil {
		return "", err
	}

	// Wrap module functions and create callbacks
	callbacks := ml.createCallbacks(moduleEnv, stateName)

	// Register with state machine
	err = ml.gsm.RegisterState(stateName, callbacks)
	if err != nil {
		return "", fmt.Errorf("failed to register state '%s': %w", stateName, err)
	}

	// Store module for reference
	ml.loadedModules[stateName] = moduleEnv

	return stateName, nil
}

// createModuleEnvironment creates an isolated Lua environment for a module
func (ml *ModuleLoader) createModuleEnvironment(stateName string) *lua.LTable {
	env := ml.L.NewTable()

	// Set up metatable to allow access to globals and context
	mt := ml.L.NewTable()
	mt.RawSetString("__index", ml.L.NewFunction(func(L *lua.LState) int {
		key := L.CheckString(-1)

		// First check module environment
		if val := env.RawGetString(key); val != lua.LNil {
			L.Push(val)
			return 1
		}

		// Then check context
		if context := ml.L.GetGlobal("context"); context != lua.LNil {
			if tbl, ok := context.(*lua.LTable); ok {
				if val := tbl.RawGetString(key); val != lua.LNil {
					L.Push(val)
					return 1
				}
			}
		}

		// Finally check globals (rf.*, game.*, etc.)
		if val := L.GetGlobal(key); val != lua.LNil {
			L.Push(val)
			return 1
		}

		L.Push(lua.LNil)
		return 1
	}))

	env.Metatable = mt

	// Inherit important globals into module environment
	// These are directly accessible
	env.RawSetString("rf", ml.L.GetGlobal("rf"))
	env.RawSetString("game", ml.L.GetGlobal("game"))
	env.RawSetString("math", ml.L.GetGlobal("math"))
	env.RawSetString("string", ml.L.GetGlobal("string"))
	env.RawSetString("table", ml.L.GetGlobal("table"))
	env.RawSetString("type", ml.L.GetGlobal("type"))
	env.RawSetString("tostring", ml.L.GetGlobal("tostring"))
	env.RawSetString("tonumber", ml.L.GetGlobal("tonumber"))
	env.RawSetString("ipairs", ml.L.GetGlobal("ipairs"))
	env.RawSetString("pairs", ml.L.GetGlobal("pairs"))
	env.RawSetString("next", ml.L.GetGlobal("next"))
	env.RawSetString("print", ml.L.GetGlobal("print"))
	env.RawSetString("error", ml.L.GetGlobal("error"))
	env.RawSetString("assert", ml.L.GetGlobal("assert"))
	env.RawSetString("select", ml.L.GetGlobal("select"))
	env.RawSetString("unpack", ml.L.GetGlobal("unpack"))

	// Add context reference (shared across all modules)
	// Context is set dynamically, but we reference the global one
	env.RawSetString("context", ml.L.GetGlobal("context"))

	return env
}

// validateModule checks that all required functions are defined
func (ml *ModuleLoader) validateModule(env *lua.LTable, filename string) error {
	missing := []string{}

	for _, funcName := range RequiredFunctions {
		fn := env.RawGetString(funcName)
		if fn == lua.LNil || fn.Type() != lua.LTFunction {
			missing = append(missing, funcName+"()")
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("module '%s' is missing required functions: %s", filename, strings.Join(missing, ", "))
	}

	return nil
}

// createCallbacks wraps module functions into StateMachine callbacks
func (ml *ModuleLoader) createCallbacks(env *lua.LTable, stateName string) statemachine.LuaCallbacks {
	callbacks := statemachine.LuaCallbacks{}

	// Helper to get function from module
	getFunc := func(name string) *lua.LFunction {
		fn := env.RawGetString(name)
		if fn != lua.LNil && fn.Type() == lua.LTFunction {
			return fn.(*lua.LFunction)
		}
		return nil
	}

	// Required functions
	if fn := getFunc("_INIT"); fn != nil {
		callbacks.Initialize = func(sm *statemachine.StateMachine) error {
			return ml.callModuleFunction(env, fn, nil)
		}
	}

	if fn := getFunc("_UPDATE"); fn != nil {
		callbacks.Update = func(dt float64) {
			// Handle errors gracefully - Update errors shouldn't crash the game
			if err := ml.callModuleFunction(env, fn, lua.LNumber(dt)); err != nil {
				// Error in _UPDATE - log it but continue
				// In dev mode, this could be logged to debug output
				_ = err // Silently continue - game should keep running
			}
		}
	}

	if fn := getFunc("_DRAW"); fn != nil {
		callbacks.Draw = func() {
			ml.callModuleFunction(env, fn, nil)
		}
	}

	if fn := getFunc("_HANDLE_INPUT"); fn != nil {
		callbacks.HandleInput = func(sm *statemachine.StateMachine) {
			ml.callModuleFunction(env, fn, nil)
		}
	}

	if fn := getFunc("_DONE"); fn != nil {
		callbacks.Shutdown = func() {
			ml.callModuleFunction(env, fn, nil)
		}
	}

	// Optional functions
	if fn := getFunc("_ENTER"); fn != nil {
		callbacks.Enter = func(sm *statemachine.StateMachine) {
			// Handle errors gracefully - log but don't crash
			if err := ml.callModuleFunction(env, fn, nil); err != nil {
				// Error in _ENTER - log it but continue
				// In dev mode, this could be logged to debug output
				_ = err // Silently continue - state transition should still work
			}
		}
	}

	if fn := getFunc("_EXIT"); fn != nil {
		callbacks.Exit = func(sm *statemachine.StateMachine) {
			ml.callModuleFunction(env, fn, nil)
		}
	}

	return callbacks
}

// callModuleFunction calls a function from the module environment
func (ml *ModuleLoader) callModuleFunction(env *lua.LTable, fn *lua.LFunction, arg lua.LValue) error {
	// Safety check: if LState is nil or invalid, skip call
	if ml.L == nil {
		return fmt.Errorf("Lua state is nil (likely due to hot reload)")
	}

	ml.L.Push(fn)
	ml.L.SetFEnv(ml.L.Get(-1), env)

	if arg != nil {
		ml.L.Push(arg)
		err := ml.L.PCall(1, 0, nil)
		if err != nil {
			return err
		}
	} else {
		err := ml.L.PCall(0, 0, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetLoadedModules returns the list of loaded module names
func (ml *ModuleLoader) GetLoadedModules() []string {
	modules := make([]string, 0, len(ml.loadedModules))
	for name := range ml.loadedModules {
		modules = append(modules, name)
	}
	return modules
}

// UnloadModule removes a module (for testing/cleanup)
func (ml *ModuleLoader) UnloadModule(stateName string) {
	delete(ml.loadedModules, stateName)
	ml.gsm.UnregisterState(stateName)
}
