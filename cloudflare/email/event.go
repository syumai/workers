//go:build js && wasm

package email

import (
	"context"
	"errors"
	"syscall/js"

	"github.com/syumai/workers"
	"github.com/syumai/workers/internal/jsutil"
	"github.com/syumai/workers/internal/runtimecontext"
)

var emailHandler EmailHandler

func init() {
	emailHandler := js.FuncOf(func(this js.Value, args []js.Value) any {
		message := args[0]
		var cb js.Func
		cb = js.FuncOf(func(_ js.Value, pArgs []js.Value) any {
			defer cb.Release()
			resolve := pArgs[0]
			reject := pArgs[1]
			go func() {
				if len(args) > 1 {
					reject.Invoke(jsutil.Errorf("too many args given to handleEmail: %d", len(args)))
					return
				}
				err := processEmail(message)
				if err != nil {
					reject.Invoke(jsutil.Error(err.Error()))
					return
				}
				resolve.Invoke(js.Undefined())
			}()
			return js.Undefined()
		})
		return jsutil.NewPromise(cb)
	})
	jsutil.Binding.Set("handleEmail", emailHandler)
}

func NewForwardableEmailMessage(ctx context.Context) (*ForwardableEmailMessage, error) {
	obj := runtimecontext.MustExtractTriggerObj(ctx)
	if obj.IsUndefined() {
		return nil, errors.New("email event is null")
	}
	return &ForwardableEmailMessage{
		From: obj.Get("from").String(),
		To:   obj.Get("to").String(),
		// Headers: obj.Get("headers"),
		// Raw:     obj.Get("raw"),
		// RawSize: int(obj.Get("rawSize").Float()),
	}, nil
}

type EmailHandler func(message *ForwardableEmailMessage) error

type ForwardableEmailMessage struct {
	From string
	To   string
}

func processEmail(eventObj js.Value) error {
	ctx := runtimecontext.New(context.Background(), eventObj)
	message, err := NewForwardableEmailMessage(ctx)
	if err != nil {
		return err
	}
	return emailHandler(message)
}

func HandleEmail(handler EmailHandler) {
	emailHandler = handler
	workers.Ready()
}
