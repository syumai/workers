package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/syumai/workers"
	"github.com/syumai/workers/cloudflare"
	"github.com/syumai/workers/cloudflare/fetch"
)

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		obj := cloudflare.GetBinding(ctx, "hello")
		cli := fetch.NewClient(fetch.WithBinding(obj))
		r, err := fetch.NewRequest(ctx, http.MethodGet, req.URL.String(), nil)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		res, err := cli.Do(r)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		io.Copy(w, res.Body)
	})
	workers.Serve(handler)
}
