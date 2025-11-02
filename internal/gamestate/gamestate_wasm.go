//go:build js && wasm

package gamestate

import "syscall/js"

type jsValue interface {
	Truthy() bool
	Call(method string, args ...interface{}) js.Value
}

// getJSConsole returns the JavaScript console object for WASM builds
func getJSConsole() jsValue {
	return js.Global().Get("console")
}

