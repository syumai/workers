//go:build js && wasm

package fmt

import (
	"syscall/js"
)

// Printer is an interface that abstracts the Println method.
type Printer interface {
	Println(args ...interface{})
}

// Println uses JavaScript's console.log to print the arguments.
func Println(args ...interface{}) {
	js.Global().Get("console").Call("log", args...)
}
