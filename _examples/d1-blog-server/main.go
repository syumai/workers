package main

import (
	"net/http"

	"github.com/syumai/workers"
	"github.com/syumai/workers/_examples/d1-blog-server/app"
)

func main() {
	http.Handle("/articles", app.NewArticleHandler())
	workers.Serve(nil) // use http.DefaultServeMux
}
