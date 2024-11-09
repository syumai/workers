package queues

import (
	"errors"
	"fmt"
	"syscall/js"
	"testing"
	"time"

	"github.com/syumai/workers/internal/jsutil"
)

func validatingProducer(t *testing.T, validateFn func(message js.Value, options js.Value) error) *Producer {
	sendFn := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		sendArg := args[0] // this should be batch (in case of SendBatch) or a single message (in case of Send)
		var options js.Value
		if len(args) > 1 {
			options = args[1]
		}
		return jsutil.NewPromise(js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			resolve := args[0]
			go func() {
				if err := validateFn(sendArg, options); err != nil {
					// must be non-fatal to avoid a deadlock
					t.Errorf("validation failed: %v", err)
				}
				resolve.Invoke(js.Undefined())
			}()
			return js.Undefined()
		}))
	})

	queue := jsutil.NewObject()
	queue.Set("send", sendFn)
	queue.Set("sendBatch", sendFn)

	return &Producer{queue: queue}
}

func TestSend(t *testing.T) {
	t.Run("text content type", func(t *testing.T) {
		validation := func(message js.Value, options js.Value) error {
			if message.Type() != js.TypeString {
				return errors.New("message body must be a string")
			}
			if message.String() != "hello" {
				return errors.New("message body must be 'hello'")
			}
			if options.Get("contentType").String() != "text" {
				return errors.New("content type must be text")
			}
			return nil
		}

		producer := validatingProducer(t, validation)
		err := producer.Send("hello", WithContentType(QueueContentTypeText))
		if err != nil {
			t.Fatalf("Send failed: %v", err)
		}
	})

	t.Run("json content type", func(t *testing.T) {
		validation := func(message js.Value, options js.Value) error {
			if message.Type() != js.TypeString {
				return errors.New("message body must be a string")
			}
			if message.String() != "hello" {
				return errors.New("message body must be 'hello'")
			}
			if options.Get("contentType").String() != "json" {
				return errors.New("content type must be json")
			}
			return nil
		}

		producer := validatingProducer(t, validation)
		err := producer.Send("hello", WithContentType(QueueContentTypeJSON))
		if err != nil {
			t.Fatalf("Send failed: %v", err)
		}
	})
}

func TestSend_ContentTypeOption(t *testing.T) {
	tests := []struct {
		name                string
		options             []SendOption
		expectedContentType string
		expectedDelaySec    int
		wantErr             bool
	}{
		{
			name:                "text",
			options:             []SendOption{WithContentType(QueueContentTypeText)},
			expectedContentType: "text",
		},
		{
			name:                "json",
			options:             []SendOption{WithContentType(QueueContentTypeJSON)},
			expectedContentType: "json",
		},
		{
			name:                "default",
			options:             nil,
			expectedContentType: "json",
		},
		{
			name:                "v8",
			options:             []SendOption{WithContentType(QueueContentTypeV8)},
			expectedContentType: "v8",
		},
		{
			name:                "bytes",
			options:             []SendOption{WithContentType(QueueContentTypeBytes)},
			expectedContentType: "bytes",
		},

		{
			name:                "delay",
			options:             []SendOption{WithDelaySeconds(5 * time.Second)},
			expectedDelaySec:    5,
			expectedContentType: "json",
		},

		{
			name:    "invalid content type",
			options: []SendOption{WithContentType("invalid")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validation := func(message js.Value, options js.Value) error {
				gotCT := options.Get("contentType").String()
				if gotCT != string(tt.expectedContentType) {
					return fmt.Errorf("expected content type %q, got %q", tt.expectedContentType, gotCT)
				}
				gotDelaySec := jsutil.MaybeInt(options.Get("delaySeconds"))
				if gotDelaySec != tt.expectedDelaySec {
					return fmt.Errorf("expected delay %d, got %d", tt.expectedDelaySec, gotDelaySec)
				}
				return nil
			}

			producer := validatingProducer(t, validation)
			err := producer.Send("hello", tt.options...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error: %t, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestSendBatch_Defaults(t *testing.T) {
	validation := func(batch js.Value, options js.Value) error {
		if batch.Type() != js.TypeObject {
			return errors.New("message batch must be an object (array)")
		}
		if batch.Length() != 2 {
			return fmt.Errorf("expected 2 messages, got %d", batch.Length())
		}
		first := batch.Index(0)
		if first.Get("body").String() != "hello" {
			return fmt.Errorf("first message body must be 'hello', was %s", first.Get("body"))
		}
		if first.Get("options").Get("contentType").String() != "json" {
			return fmt.Errorf("first message content type must be json, was %s", first.Get("options").Get("contentType"))
		}

		second := batch.Index(1)
		if second.Get("body").String() != "world" {
			return fmt.Errorf("second message body must be 'world', was %s", second.Get("body"))
		}
		if second.Get("options").Get("contentType").String() != "text" {
			return fmt.Errorf("second message content type must be text, was %s", second.Get("options").Get("contentType"))
		}

		return nil
	}

	var batch []*BatchMessage = []*BatchMessage{
		NewBatchMessage("hello"),
		NewBatchMessage("world", WithContentType(QueueContentTypeText)),
	}

	producer := validatingProducer(t, validation)
	err := producer.SendBatch(batch)
	if err != nil {
		t.Fatalf("SendBatch failed: %v", err)
	}
}

func TestSendBatch_Options(t *testing.T) {
	validation := func(_ js.Value, options js.Value) error {
		if options.Get("delaySeconds").Int() != 5 {
			return fmt.Errorf("expected delay 5, got %d", options.Get("delaySeconds").Int())
		}
		return nil
	}

	var batch []*BatchMessage = []*BatchMessage{
		NewTextBatchMessage("hello"),
	}

	producer := validatingProducer(t, validation)
	err := producer.SendBatch(batch, WithBatchDelaySeconds(5*time.Second))
	if err != nil {
		t.Fatalf("SendBatch failed: %v", err)
	}
}
