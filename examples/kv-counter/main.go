package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/syumai/workers"
)

// counterNamespace is a bounded KV namespace for storing counter.
const counterNamespace = "COUNTER"

// countKey is a key to store current count value to the KV namespace.
const countKey = "count"

func handleErr(w http.ResponseWriter, msg string, err error) {
	log.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(msg))
}

func main() {
	// initialize KV namespace instance
	kv, err := workers.NewKVNamespace(counterNamespace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init KV: %v", err)
		os.Exit(1)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		countStr, err := kv.GetString(countKey, nil)
		if err != nil {
			handleErr(w, "failed to get current count\n", err)
			return
		}

		/*
			countReader, err := kv.GetReader(countKey, nil)
			if err != nil {
				handleErr(w, "failed to get current count\n", err)
				return
			}
			b, _ := io.ReadAll(countReader)
			countStr := string(b)
		*/

		// ignore err and treat count value as 0
		count, _ := strconv.Atoi(countStr)

		nextCountStr := strconv.Itoa(count + 1)

		err = kv.PutString(countKey, nextCountStr, nil)
		if err != nil {
			handleErr(w, "failed to put next count\n", err)
			return
		}

		/*
			err = kv.PutReader(countKey, strings.NewReader(nextCountStr), nil)
			if err != nil {
				handleErr(w, "failed to put next count\n", err)
				return
			}
		*/

		w.Header().Set("Content-Type", "text/plain")

		/*
			// List returns only `count` as the keys in this namespace.
			v, err := kv.List(nil)
			if err != nil {
				handleErr(w, "failed to list\n", err)
				return
			}
			for i, key := range v.Keys {
				fmt.Fprintf(w, "%d: %s\n", i, key.Name)
			}
		*/

		w.Write([]byte(nextCountStr))
	})

	workers.Serve(nil)
}
