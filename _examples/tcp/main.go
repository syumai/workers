package main

import (
	"bufio"
	"net/http"
	"time"

	"github.com/syumai/workers"
	"github.com/syumai/workers/cloudflare/socket"
)

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		dialer, err := socket.NewDialer(req.Context(), nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		conn, err := dialer.Dial(req.Context(), "tcp", "tcpbin.com:4242")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer conn.Close()
		conn.SetDeadline(time.Now().Add(1 * time.Hour))
		_, err = conn.Write([]byte("hello.\n"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rd := bufio.NewReader(conn)
		bts, err := rd.ReadBytes('.')
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(bts)
	})
	workers.Serve(handler)
}
