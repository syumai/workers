package fetch

import (
	"errors"
	"net/http"
	"strconv"
	"syscall/js"

	"github.com/syumai/workers/internal/jshttp"
	"github.com/syumai/workers/internal/jsutil"
)

// fetch is a function that reproduces cloudflare fetch.
// Docs: https://developers.cloudflare.com/workers/runtime-apis/fetch/
func fetch(namespace js.Value, req *http.Request, init *RequestInit) (*http.Response, error) {
	if namespace.IsUndefined() {
		return nil, errors.New("fetch function not found")
	}
	fetchFunc := namespace.Get("fetch")
	promise := fetchFunc.Invoke(
		// The Request object to fetch.
		// Docs: https://developers.cloudflare.com/workers/runtime-apis/request
		jshttp.ToJSRequest(req),
		// The content of the request.
		// Docs: https://developers.cloudflare.com/workers/runtime-apis/request#requestinit
		init.ToJS(),
	)

	jsRes, err := jsutil.AwaitPromise(promise)
	if err != nil {
		return nil, err
	}

	// Create TransformStream
	ts := js.Global().Get("IdentityTransformStream").New()
	readable := ts.Get("readable")
	writable := ts.Get("writable")
	jsRes.Get("body").Call("pipeTo", writable)

	// Create response
	res := new(http.Response)
	res.StatusCode = jsRes.Get("status").Int()
	res.Status = strconv.Itoa(res.StatusCode) + " " + jsRes.Get("statusText").String()
	res.Header = jshttp.ToHeader(jsRes.Get("headers"))
	res.Body = jsutil.ConvertReadableStreamToReadCloser(readable)

	return res, nil
}
