package cloudflare

import (
	"context"
	"syscall/js"

	"github.com/syumai/workers/cloudflare/internal/cfruntimecontext"
	"github.com/syumai/workers/internal/jsutil"
)

// WaitUntil extends the lifetime of the "fetch" event.
// It accepts an asynchronous task which the Workers runtime will execute before the handler terminates but without blocking the response.
// see: https://developers.cloudflare.com/workers/runtime-apis/fetch-event/#waituntil
func WaitUntil(ctx context.Context, task func()) {
	executionContext := cfruntimecontext.GetExecutionContext(ctx)

	executionContext.Call("waitUntil", jsutil.NewPromise(js.FuncOf(func(this js.Value, pArgs []js.Value) any {
		resolve := pArgs[0]
		go func() {
			task()
			resolve.Invoke(js.Undefined())
		}()
		return js.Undefined()
	})))
}
