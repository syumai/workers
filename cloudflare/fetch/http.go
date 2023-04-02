package fetch

import (
	"net/http"

	"github.com/syumai/workers/internal/jshttp"
	"github.com/syumai/workers/internal/jsutil"
)

// Do sends an HTTP request and returns an HTTP response
func (c *Client) Do(req *Request) (*http.Response, error) {
	jsReq := jshttp.ToJSRequest(req.Request)

	init := jsutil.NewObject()
	promise := c.namespace.Call("fetch2", jsReq, init)
	jsRes, err := jsutil.AwaitPromise(promise)
	if err != nil {
		return nil, err
	}

	return jshttp.ToResponse(jsRes)
}
