//go:build !(js && wasm)

package statemachine

type nilJSValue struct{}

func (n nilJSValue) Truthy() bool { return false }
func (n nilJSValue) Call(method string, args ...interface{}) interface{} { return nil }

type jsValue interface {
	Truthy() bool
	Call(method string, args ...interface{}) interface{}
}

func getJSConsole() jsValue {
	return nilJSValue{}
}

