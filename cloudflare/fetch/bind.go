package fetch

import (
	"net/http"
	"syscall/js"

	"github.com/syumai/workers/internal/jshttp"
	"github.com/syumai/workers/internal/jsutil"
)

var globalFetchFunc = js.Global().Get("fetch")

// fetch is a function that reproduces cloudflare fetch.
// Docs: https://developers.cloudflare.com/workers/runtime-apis/fetch/
func fetch(namespace js.Value, req *http.Request, init *RequestInit) (*http.Response, error) {
	// The Request object to fetch.
	// Docs: https://developers.cloudflare.com/workers/runtime-apis/request
	reqObj := jshttp.ToJSRequest(req)
	// The content of the request.
	// Docs: https://developers.cloudflare.com/workers/runtime-apis/request#requestinit
	reqInit := init.ToJS()

	promise := func() js.Value {
		if namespace.IsUndefined() {
			return globalFetchFunc.Invoke(reqObj, reqInit)
		}
		return namespace.Call("fetch", reqObj, reqInit)
	}()

	jsRes, err := jsutil.AwaitPromise(promise)
	if err != nil {
		return nil, err
	}

	return jshttp.ToResponse(jsRes)
}
