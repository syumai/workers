package workers

import (
	"io"
	"net/http"
	"sync"
)

type responseWriterBuffer struct {
	header     http.Header
	statusCode int
	reader     *io.PipeReader
	writer     *io.PipeWriter
	readyCh    chan struct{}
	once       sync.Once
}

var _ http.ResponseWriter = &responseWriterBuffer{}

// ready indicates that responseWriterBuffer is ready to be converted to Response.
func (w *responseWriterBuffer) ready() {
	w.once.Do(func() {
		close(w.readyCh)
	})
}

func (w *responseWriterBuffer) Write(data []byte) (n int, err error) {
	w.ready()
	return w.writer.Write(data)
}

func (w *responseWriterBuffer) Header() http.Header {
	return w.header
}

func (w *responseWriterBuffer) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}
