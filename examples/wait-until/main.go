package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/syumai/workers"
	"github.com/syumai/workers/cloudflare"
)

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cloudflare.WaitUntil(req.Context(), func() {
			for i := 0; i < 5; i++ {
				time.Sleep(time.Second)
			}

			fmt.Println("5-second task completed")
		})

		w.Write([]byte("http response done"))
	})
	workers.Serve(handler)
}
