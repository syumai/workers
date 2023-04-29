package jshttp

import (
	"io"
	"net/http"
	"sync"
	"syscall/js"
)

type ResponseWriterBuffer struct {
	HeaderValue http.Header
	StatusCode  int
	Reader      *io.PipeReader
	Writer      *io.PipeWriter
	ReadyCh     chan struct{}
	Once        sync.Once
}

var _ http.ResponseWriter = &ResponseWriterBuffer{}

// Ready indicates that ResponseWriterBuffer is ready to be converted to Response.
func (w *ResponseWriterBuffer) Ready() {
	w.Once.Do(func() {
		close(w.ReadyCh)
	})
}

func (w *ResponseWriterBuffer) Write(data []byte) (n int, err error) {
	w.Ready()
	return w.Writer.Write(data)
}

func (w *ResponseWriterBuffer) Header() http.Header {
	return w.HeaderValue
}

func (w *ResponseWriterBuffer) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
}

// ToJSResponse converts *ResponseWriterBuffer to JavaScript sides Response.
//   - Response: https://developer.mozilla.org/docs/Web/API/Response
func (w *ResponseWriterBuffer) ToJSResponse() js.Value {
	return newJSResponse(w.StatusCode, w.HeaderValue, w.Reader)
}
