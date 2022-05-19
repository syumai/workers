package main

import (
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

func handleRequest(w http.ResponseWriter, req *http.Request) {
	if !authenticate(req) {
		w.Header().Add("WWW-Authenticate", `Basic realm="login is required"`)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized\n"))
		return
	}
	w.Write([]byte("Authorized!\n"))
}

func main() {
	workers.Serve(http.HandlerFunc(handleRequest))
}
