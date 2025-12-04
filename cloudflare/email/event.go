//go:build js && wasm

package email

import (
	"context"
	"errors"
	"fmt"
	"io"
	"syscall/js"
	"time"

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
				err := invokeEmailHandler(message)
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

func newForwardableEmailMessage(ctx context.Context) (*forwardableEmailMessage, error) {
	obj := runtimecontext.MustExtractTriggerObj(ctx)
	if obj.IsUndefined() {
		return nil, errors.New("email event is null")
	}

	return &forwardableEmailMessage{
		from: obj.Get("from").String(),
		to:   obj.Get("to").String(),
		raw:  obj.Get("raw"),
		// rawSize: obj.Get("rawSize").Int(),
	}, nil
}

type Handler func(m ForwardableEmailMessage) error

type Email interface {
	From() string
	To() string
	Raw() io.ReadCloser
}

// Emails that originate from inbound handler, can forward it onward or drop, etc
type ForwardableEmailMessage interface {
	From() string
	To() string
	Raw() io.ReadCloser
}

type forwardableEmailMessage struct {
	from string
	to   string
	raw  js.Value
}

func (f *forwardableEmailMessage) From() string {
	return f.from
}
func (f *forwardableEmailMessage) To() string {
	return f.to
}
func (f *forwardableEmailMessage) Raw() io.ReadCloser {
	return jsutil.ConvertReadableStreamToReadCloser(f.raw)
}

// Emails that we're sending outbound
type EmailSendable interface {
	From() string
	To() string
}

type EmailMessage struct {
	from string
	to   string
}

func (e *EmailMessage) From() string {
	return e.from
}
func (e *EmailMessage) To() string {
	return e.to
}
func NewEmailMessage(from string, to string) *EmailMessage {
	return &EmailMessage{
		from: from,
		to:   to,
	}
}

type EmailClient struct {
	bind js.Value
}

func NewClient(bind js.Value) *EmailClient {
	return &EmailClient{
		bind: bind,
	}
}
func (c *EmailClient) Send(m EmailSendable) error {
	// Check if the EMAIL binding is available
	if c.bind.IsUndefined() {
		return errors.New("EMAIL binding not found. Make sure you have [[send_email]] configured in your wrangler.toml or wrangler.jsonc")
	}

	// Check if the binding has a send method
	sendMethod := c.bind.Get("send")
	if sendMethod.IsUndefined() {
		return errors.New("EMAIL binding does not have a send method. Make sure you have [[send_email]] configured in your wrangler.toml or wrangler.jsonc")
	}

	// Get EmailMessage from the global RuntimeContext (jsutil.RuntimeContext)
	runtimeCtx := jsutil.RuntimeContext
	emailMessageCtor := runtimeCtx.Get("EmailMessage")

	if emailMessageCtor.IsUndefined() {
		return errors.New("EmailMessage not found in runtime context")
	}

	// Create a proper MIME message with Message-ID
	messageId := fmt.Sprintf("<%d@example.com>", time.Now().UnixNano())
	mimeMessage := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: Test Email\r\nMessage-ID: %s\r\nDate: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\nThis is the email body",
		m.From(), m.To(), messageId, time.Now().Format(time.RFC1123),
	)

	// Create EmailMessage instance using the constructor
	emailMsg := emailMessageCtor.New(m.From(), m.To(), mimeMessage)

	// Send using the EMAIL binding
	_, err := jsutil.AwaitPromise(c.bind.Call("send", emailMsg))
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

func invokeEmailHandler(eventObj js.Value) error {
	ctx := runtimecontext.New(context.Background(), eventObj)
	message, err := newForwardableEmailMessage(ctx)
	if err != nil {
		return err
	}
	return emailHandler(message)
}

func Handle(handler Handler) {
	emailHandler = handler
	workers.Ready()
}
