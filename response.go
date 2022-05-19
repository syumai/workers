package workers

import (
	"net/http"
	"syscall/js"
)

func toJSHeader(header http.Header) js.Value {
	h := headersClass.New()
	for key, values := range header {
		for _, value := range values {
			h.Call("append", key, value)
		}
	}
	return h
}

func toJSResponse(w *responseWriterBuffer) (js.Value, error) {
	<-w.readyCh // wait until ready
	status := w.statusCode
	if status == 0 {
		status = http.StatusOK
	}
	respInit := newObject()
	respInit.Set("status", status)
	respInit.Set("statusText", http.StatusText(status))
	respInit.Set("headers", toJSHeader(w.Header()))
	readableStream := convertReaderToReadableStream(w.reader)
	return responseClass.New(readableStream, respInit), nil
}
