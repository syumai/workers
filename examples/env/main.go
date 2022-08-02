package main

import (
	"fmt"
	"net/http"
	"syscall/js"

	"github.com/syumai/workers"
)

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		envVal := js.Global().Get("MY_ENV")
		fmt.Fprintf(w, "MY_ENV: %s", envVal)
	})
	workers.Serve(handler)
}
