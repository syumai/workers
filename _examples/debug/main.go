package main

import (
	"fmt"
	"net/http"

	"github.com/syumai/workers"
	fmtMock "github.com/syumai/workers/cloudflare/fmt/mock"
)

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("name")
		if name == "" {
			name = "world"
		}
		fmt.Fprintf(w, "Hello, %s!", name)
		fmtMock.Println("Request received:", req.Method, req.URL)
	})
	workers.Serve(handler)
}
