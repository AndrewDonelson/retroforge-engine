//go:build !(js && wasm)

package gamestate

type nilJSValue struct{}

func (n nilJSValue) Truthy() bool { return false }
func (n nilJSValue) Call(method string, args ...interface{}) interface{} { return nil }

type jsValue interface {
	Truthy() bool
	Call(method string, args ...interface{}) interface{}
}

// getJSConsole returns a no-op value for non-WASM builds
func getJSConsole() jsValue {
	return nilJSValue{}
}

