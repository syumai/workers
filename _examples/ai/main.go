package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/syumai/workers"
	"github.com/syumai/workers/cloudflare/ai"
)

func main() {

	http.HandleFunc("/ai", func(w http.ResponseWriter, req *http.Request) {

		// initialize AI namespace instance
		aiCaller, err := ai.NewNamespace("AI")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to init AI instance: %v", err)
			return
		}

		aiJsonResultStr, err := aiCaller.Run("@cf/meta/llama-3.1-8b-instruct", map[string]interface{}{
			"prompt": "What is the origin of the phrase Hello, World",
		})

		if err != nil {
			fmt.Println(w, "failed to get result from AI\n", err)
			return
		}

		fmt.Println(aiJsonResultStr)

		io.Copy(w, strings.NewReader(aiJsonResultStr))
	})

	workers.Serve(nil) // use http.DefaultServeMux
}
