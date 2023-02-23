package jshttp

import (
	"net/http"
	"syscall/js"

	"github.com/syumai/workers/internal/jsutil"
)

func ToJSHeader(header http.Header) js.Value {
	h := jsutil.HeadersClass.New()
	for key, values := range header {
		for _, value := range values {
			h.Call("append", key, value)
		}
	}
	return h
}

func ToJSResponse(w *ResponseWriterBuffer) (js.Value, error) {
	<-w.ReadyCh // wait until ready
	status := w.StatusCode
	if status == 0 {
		status = http.StatusOK
	}
	respInit := jsutil.NewObject()
	respInit.Set("status", status)
	respInit.Set("statusText", http.StatusText(status))
	respInit.Set("headers", ToJSHeader(w.Header()))
	readableStream := jsutil.ConvertReaderToReadableStream(w.Reader)
	return jsutil.ResponseClass.New(readableStream, respInit), nil
}
