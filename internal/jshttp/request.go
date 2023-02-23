package jshttp

import (
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"syscall/js"

	"github.com/syumai/workers/internal/jsutil"
)

// ToBody converts JavaScript sides ReadableStream (can be null) to io.ReadCloser.
//   - ReadableStream: https://developer.mozilla.org/en-US/docs/Web/API/ReadableStream
func ToBody(streamOrNull js.Value) io.ReadCloser {
	if streamOrNull.IsNull() {
		return nil
	}
	sr := streamOrNull.Call("getReader")
	return io.NopCloser(jsutil.ConvertStreamReaderToReader(sr))
}

// ToRequest converts JavaScript sides Request to *http.Request.
//   - Request: https://developer.mozilla.org/ja/docs/Web/API/Request
func ToRequest(req js.Value) (*http.Request, error) {
	reqUrl, err := url.Parse(req.Get("url").String())
	if err != nil {
		return nil, err
	}
	header := ToHeader(req.Get("headers"))

	// ignore err
	contentLength, _ := strconv.ParseInt(header.Get("Content-Length"), 10, 64)
	return &http.Request{
		Method:           req.Get("method").String(),
		URL:              reqUrl,
		Header:           header,
		Body:             ToBody(req.Get("body")),
		ContentLength:    contentLength,
		TransferEncoding: strings.Split(header.Get("Transfer-Encoding"), ","),
		Host:             header.Get("Host"),
	}, nil
}
