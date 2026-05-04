package fetch

import (
	"errors"
	"net/http"
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
	var promise js.Value
	// The Request object to fetch.
	// Docs: https://developers.cloudflare.com/workers/runtime-apis/request
	var reqValue = jshttp.ToJSRequest(req)
	// The content of the request.
	// Docs: https://developers.cloudflare.com/workers/runtime-apis/request#requestinit
	var initValue = init.ToJS()
	if namespace.Equal(js.Global()) {
		// global fetch
		fetchFunc := namespace.Get("fetch")
		promise = fetchFunc.Invoke(reqValue, initValue)
	} else {
		// service-binding
		promise = namespace.Call("fetch", reqValue, initValue)
	}

	jsRes, err := jsutil.AwaitPromise(promise)
	if err != nil {
		return nil, err
	}

	return jshttp.ToResponse(jsRes)
}
