//go:build js && wasm

package ai

import (
	"fmt"
	"io"
	"syscall/js"

	"github.com/syumai/workers/cloudflare/internal/cfruntimecontext"
	"github.com/syumai/workers/internal/jsutil"
)

type AIInterface interface {
	WaitUntil(task func())
	Run(key string, opts map[string]interface{}) (string, error)
	RunReader(key string, opts map[string]interface{}) (io.Reader, error)
}

type AI struct {
	instance js.Value
}

// NewNamespace returns Namespace for given variable name.
//   - variable name must be defined in wrangler.toml as kv_namespace's binding.
//   - if the given variable name doesn't exist on runtime context, returns error.
//   - This function panics when a runtime context is not found.
func NewNamespace(varName string) (*AI, error) {
	inst := cfruntimecontext.MustGetRuntimeContextEnv().Get(varName)
	if inst.IsUndefined() {
		return nil, fmt.Errorf("%s is undefined", varName)
	}
	return &AI{instance: inst}, nil
}

func (ns *AI) WaitUntil(task func()) {
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

func mapToJS(opts map[string]interface{}, type_ string) js.Value {
	obj := jsutil.NewObject()
	for k, v := range opts {
		obj.Set(k, v)
	}
	obj.Set("type", type_)
	return obj
}

func (ns *AI) Run(key string, opts map[string]interface{}) (string, error) {
	p := ns.instance.Call("run", key, mapToJS(opts, "text"))
	v, err := jsutil.AwaitPromise(p)
	if err != nil {
		return "", err
	}
	respString := js.Global().Get("JSON").Call("stringify", v).String()
	return respString, nil
}

// GetReader gets stream value by the specified key.
//   - if a network error happens, returns error.
func (ns *AI) RunReader(key string, opts map[string]interface{}) (io.Reader, error) {
	p := ns.instance.Call("run", key, mapToJS(opts, "stream"))
	v, err := jsutil.AwaitPromise(p)
	if err != nil {
		return nil, err
	}
	return jsutil.ConvertReadableStreamToReadCloser(v), nil
}
