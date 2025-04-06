package ai

import (
	"io"
	"syscall/js"

	"github.com/syumai/workers/internal/jsutil"
)

// GetOptions represents Cloudflare KV namespace get options.
//   - https://github.com/cloudflare/workers-types/blob/3012f263fb1239825e5f0061b267c8650d01b717/index.d.ts#L930
type AiOptions struct {
	Prompt string
}

func (opts *AiOptions) toJS(type_ string) js.Value {
	obj := jsutil.NewObject()

	if opts.Prompt != "" {
		obj.Set("prompt", opts.Prompt)
	}
	return obj
}

func (ns *Ai) WaitUntil(task func()) {
	exCtx := ns.instance
	exCtx.Call("run", jsutil.NewPromise(js.FuncOf(func(this js.Value, pArgs []js.Value) any {
		resolve := pArgs[0]
		go func() {
			task()
			resolve.Invoke(js.Undefined())
		}()
		return js.Undefined()
	})))
}

// GetString gets string value by the specified key.
//   - if a network error happens, returns error.
func (ns *Ai) Run(key string, opts *AiOptions) (string, error) {
	p := ns.instance.Call("run", key, opts.toJS("text"))
	v, err := jsutil.AwaitPromise(p)
	if err != nil {
		return "", err
	}
	respString := js.Global().Get("JSON").Call("stringify", v).String()
	return respString, nil
}

// GetReader gets stream value by the specified key.
//   - if a network error happens, returns error.
func (ns *Ai) GetReader(key string, opts *AiOptions) (io.Reader, error) {
	p := ns.instance.Call("run", key, opts.toJS("stream"))
	v, err := jsutil.AwaitPromise(p)
	if err != nil {
		return nil, err
	}
	return jsutil.ConvertReadableStreamToReadCloser(v), nil
}
