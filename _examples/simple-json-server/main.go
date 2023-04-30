package main

import (
	"net/http"

	"github.com/syumai/workers"
	"github.com/syumai/workers/examples/simple-json-server/app"
)

func main() {
	http.HandleFunc("/hello", app.HelloHandler)
	workers.Serve(nil) // use http.DefaultServeMux
}
