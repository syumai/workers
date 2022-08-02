package main

import (
	"fmt"
	"net/http"

	"github.com/syumai/workers"
)

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "MY_ENV: %s", workers.Getenv("MY_ENV"))
	})
	workers.Serve(handler)
}
