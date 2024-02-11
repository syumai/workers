package hono

import (
	"context"
	"io"
	"net/http"
	"sync"
	"syscall/js"

	"github.com/syumai/workers/internal/jshttp"
	"github.com/syumai/workers/internal/jsutil"
	"github.com/syumai/workers/internal/runtimecontext"
)

type Context struct {
	ctxObj  js.Value
	reqFunc func() *http.Request
}

func newContext(ctxObj js.Value) *Context {
	return &Context{
		ctxObj: ctxObj,
		reqFunc: sync.OnceValue(func() *http.Request {
			reqObj := ctxObj.Get("req").Get("raw")
			req, err := jshttp.ToRequest(reqObj)
			if err != nil {
				panic(err)
			}
			ctx := runtimecontext.New(context.Background(), reqObj, jsutil.RuntimeContext)
			req = req.WithContext(ctx)
			return req
		}),
	}
}

func (c *Context) Request() *http.Request {
	return c.reqFunc()
}

func (c *Context) Header() Header {
	return &header{
		headerObj: c.ctxObj.Get("req").Get("headers"),
	}
}

func (c *Context) SetStatus(statusCode int) {
	c.ctxObj.Call("status", statusCode)
}

func (c *Context) ResponseBody() io.ReadCloser {
	return jsutil.ConvertReadableStreamToReadCloser(c.ctxObj.Get("res").Get("body"))
}

func (c *Context) SetResponseBody(body io.ReadCloser) {
	var res js.Value
	if sr, ok := body.(jsutil.RawJSBodyGetter); ok {
		res = jsutil.ResponseClass.New(sr, c.ctxObj.Get("res"))
	} else {
		bodyObj := jsutil.ConvertReaderToReadableStream(body)
		res = jsutil.ResponseClass.New(bodyObj, c.ctxObj.Get("res"))
	}
	c.ctxObj.Set("res", res)
}
