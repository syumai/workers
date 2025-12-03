//go:build js && wasm

package email

import (
	"context"
	"errors"
	"io"
	"net/mail"
	"syscall/js"

	"github.com/syumai/workers"
	"github.com/syumai/workers/internal/jsutil"
	"github.com/syumai/workers/internal/runtimecontext"
)

var emailHandler Handler

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
		raw: obj.Get("raw"),
	}, nil
}

type Handler func(msg *mail.Message) error

type ForwardableEmailMessage struct {
	// 'from', 'to', 'headers' are also available here,
	// but we'll hand off to golang's mail pkg for parsing
	raw js.Value
}

func (f *ForwardableEmailMessage) RawReader() io.Reader {
	return jsutil.ConvertReadableStreamToReadCloser(f.raw)
}

func processEmail(eventObj js.Value) error {
	ctx := runtimecontext.New(context.Background(), eventObj)
	message, err := NewForwardableEmailMessage(ctx)
	if err != nil {
		return err
	}
	msg, err := mail.ReadMessage(message.RawReader())
	if err != nil {
		return err
	}
	return emailHandler(msg)
}

func Handle(handler Handler) {
	emailHandler = handler
	workers.Ready()
}
