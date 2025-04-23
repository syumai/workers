//go:build !js

package fmt

import (
	"fmt"
)

// Println implements the Printer interface using fmt.Println.
func Println(args ...interface{}) {
	fmt.Println(args...)
}
