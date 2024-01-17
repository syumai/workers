package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/syumai/workers"
	"github.com/syumai/workers/cloudflare/incoming"
)

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		p := incoming.NewProperties(req.Context())

		buf, _ := json.Marshal(p)
		fmt.Fprintf(w, "%s", string(buf))
	})
	workers.Serve(handler)
}
