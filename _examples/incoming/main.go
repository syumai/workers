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

		encoder := json.NewEncoder(w)
		w.Header().Set("Content-Type", "application/json")
		if err := encoder.Encode(p); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err), http.StatusInternalServerError)
			return
		}
	})
	workers.Serve(handler)
}
