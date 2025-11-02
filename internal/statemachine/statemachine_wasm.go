//go:build js && wasm

package statemachine

import (
	"runtime"
	"syscall/js"
)

type jsValue interface {
	Truthy() bool
	Call(method string, args ...interface{}) js.Value
}

func getJSConsole() jsValue {
	if runtime.GOOS == "js" {
		return js.Global().Get("console")
	}
	return nil
}

