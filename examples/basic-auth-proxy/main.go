package main

import (
	"io"
	"log"
	"net/http"

	"github.com/syumai/workers"
)

const (
	userName     = "user"
	userPassword = "password"
)

func authenticate(req *http.Request) bool {
	username, password, ok := req.BasicAuth()
	return ok && username == userName && password == userPassword
}

func handleError(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	w.Write([]byte(msg + "\n"))
}

func handleRequest(w http.ResponseWriter, req *http.Request) {
	if !authenticate(req) {
		w.Header().Add("WWW-Authenticate", `Basic realm="login is required"`)
		handleError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	u := *req.URL
	u.Scheme = "https"
	u.Host = "syum.ai"
	proxyReq, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Internal Error")
		log.Printf("failed to create proxy request: %v\n", err)
		return
	}
	resp, err := (*Transport).RoundTrip(nil, proxyReq)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Internal Error")
		log.Printf("failed to execute proxy request: %v\n", err)
		return
	}
	for k, values := range resp.Header {
		for _, v := range values {
			w.Header().Add(k, v)
		}
	}
	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}

func main() {
	workers.Serve(http.HandlerFunc(handleRequest))
}
