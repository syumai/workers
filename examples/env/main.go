package main

import (
	"fmt"
	"net/http"

	"github.com/syumai/workers"
	"github.com/syumai/workers/cloudflare"
)

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "MY_ENV: %s", cloudflare.Getenv("MY_ENV"))
	})
	workers.Serve(handler)
}
