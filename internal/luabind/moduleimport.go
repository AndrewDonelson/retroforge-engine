package luabind

import (
	"github.com/AndrewDonelson/retroforge-engine/internal/gamestate"
	"github.com/AndrewDonelson/retroforge-engine/internal/modulestate"
	lua "github.com/yuin/gopher-lua"
)

var moduleLoaders = make(map[*lua.LState]*modulestate.ModuleLoader)

// RegisterModuleImport attaches rf.import() function to the Lua state
// This requires both the GameStateMachine and a file reader
func RegisterModuleImport(L *lua.LState, gsm *gamestate.GameStateMachine, fileReader modulestate.FileReader, basePath string) {
	loader := modulestate.NewModuleLoader(L, gsm, fileReader, basePath)

	// Store loader in map keyed by Lua state
	moduleLoaders[L] = loader

	// rf.import(filename) -> stateName
	rf := L.GetGlobal("rf").(*lua.LTable)
	L.SetField(rf, "import", L.NewFunction(func(L *lua.LState) int {
		filename := L.CheckString(1)

		// Get loader from map
		loader, exists := moduleLoaders[L]
		if !exists || loader == nil {
			L.RaiseError("module loader not initialized")
			return 0
		}

		stateName, err := loader.ImportModule(filename)
		if err != nil {
			L.RaiseError("failed to import module '%s': %v", filename, err)
			return 0
		}

		L.Push(lua.LString(stateName))
		return 1
	}))
}

// RegisterModuleImportWithMap attaches rf.import() using a file map (for cart mode)
func RegisterModuleImportWithMap(L *lua.LState, gsm *gamestate.GameStateMachine, files map[string][]byte) {
	fileReader := modulestate.NewMapFileReader(files)
	basePath := "" // Not needed for map reader
	RegisterModuleImport(L, gsm, fileReader, basePath)
}

// RegisterModuleImportWithFilesystem attaches rf.import() using filesystem (for dev mode)
func RegisterModuleImportWithFilesystem(L *lua.LState, gsm *gamestate.GameStateMachine, basePath string) {
	fileReader := modulestate.NewFileSystemReader(basePath)
	RegisterModuleImport(L, gsm, fileReader, basePath)
}

// GetModuleLoader retrieves the module loader for a Lua state
func GetModuleLoader(L *lua.LState) *modulestate.ModuleLoader {
	return moduleLoaders[L]
}

// UnregisterModuleLoader removes the module loader when Lua state is closed
func UnregisterModuleLoader(L *lua.LState) {
	delete(moduleLoaders, L)
}
