package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/syumai/workers"
	"github.com/syumai/workers/cloudflare/ai"
)

func main() {
	http.HandleFunc("/ia", func(w http.ResponseWriter, req *http.Request) {

		// initialize KV namespace instance
		aiCaller, err := ai.NewNamespace("AI")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to init KV: %v", err)
			os.Exit(1)
		}

		countStr, err := aiCaller.Run("@cf/meta/llama-3.1-8b-instruct", &ai.AiOptions{
			Prompt: "What is the origin of the phrase Hello, World",
		})

		if err != nil {
			fmt.Println(w, "failed to get current count\n", err)
			return
		}

		fmt.Println(countStr)

		// bind := cloudflare.GetBinding("AI")
		// fc := fetch.NewClient(fetch.WithBinding(bind))

		// hc := fc.HTTPClient(fetch.RedirectModeFollow)

		// payload := map[string]string{
		// 	"prompt": "What is the origin of the phrase Hello, World",
		// }

		// // Convertimos el mapa a JSON
		// jsonData, err := json.Marshal(payload)
		// if err != nil {
		// 	fmt.Println("Error al convertir a JSON:", err)
		// 	return
		// }

		// reqAI, _ := http.NewRequest("POST", "@cf/meta/llama-3.1-8b-instruct", bytes.NewBuffer(jsonData))
		// res, err := hc.Do(reqAI)
		// if err != nil {
		// 	fmt.Println(err)
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	return
		// }

		io.Copy(w, strings.NewReader(countStr))

		// io.Copy(w, strings.NewReader(countStr))

	})
	http.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		msg := "Hello!"
		w.Write([]byte(msg))
	})
	http.HandleFunc("/echo", func(w http.ResponseWriter, req *http.Request) {
		b, err := io.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}
		io.Copy(w, bytes.NewReader(b))
	})
	workers.Serve(nil) // use http.DefaultServeMux
}
