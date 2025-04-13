package main

import (
	"fmt"
	"net/http"

	"github.com/syumai/workers"
	cloudflare "github.com/syumai/workers/cloudflare/env"
	// mock "github.com/syumai/workers/cloudflare/env/mock"
)

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "MY_ENV: %s", cloudflare.Getenv("MY_ENV"))
	})
	workers.Serve(handler)
}
