package mock

import (
	"fmt"
)

// Println implements the Printer interface using fmt.Println.
func Println(args ...interface{}) {
	fmt.Println(args...)
}
