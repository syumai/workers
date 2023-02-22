package jsutil

import (
	"io"
	"net/http"
	"sync"
)

type ResponseWriterBuffer struct {
	header     http.Header
	statusCode int
	reader     *io.PipeReader
	writer     *io.PipeWriter
	readyCh    chan struct{}
	once       sync.Once
}

var _ http.ResponseWriter = &ResponseWriterBuffer{}

// ready indicates that ResponseWriterBuffer is ready to be converted to Response.
func (w *ResponseWriterBuffer) ready() {
	w.once.Do(func() {
		close(w.readyCh)
	})
}

func (w *ResponseWriterBuffer) Write(data []byte) (n int, err error) {
	w.ready()
	return w.writer.Write(data)
}

func (w *ResponseWriterBuffer) Header() http.Header {
	return w.header
}

func (w *ResponseWriterBuffer) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}
