package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/syumai/workers"
	"github.com/syumai/workers/cloudflare/ai"
	// ai "github.com/syumai/workers/cloudflare/ai/mock"
)

func main() {
	http.HandleFunc("/ai", func(w http.ResponseWriter, req *http.Request) {

		// initialize KV namespace instance
		aiCaller, err := ai.NewNamespace("AI")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to init KV: %v", err)
			os.Exit(1)
		}

		countStr, err := aiCaller.Run("@cf/meta/llama-3.1-8b-instruct", map[string]interface{}{
			"prompt": "What is the origin of the phrase Hello, World",
		})

		if err != nil {
			fmt.Println(w, "failed to get current count\n", err)
			return
		}

		fmt.Println(countStr)

		io.Copy(w, strings.NewReader(countStr))
	})
	workers.Serve(nil) // use http.DefaultServeMux
}
