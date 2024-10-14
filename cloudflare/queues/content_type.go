package queues

import (
	"fmt"
	"io"
	"syscall/js"

	"github.com/syumai/workers/internal/jsutil"
)

// QueueContentType represents the content type of a message produced to a queue.
// This information mostly affects how the message body is represented in the Cloudflare UI and is NOT
// propagated to the consumer side.
//   - https://developers.cloudflare.com/queues/configuration/javascript-apis/#queuescontenttype
type QueueContentType string

const (
	// QueueContentTypeJSON is the default content type for the produced queue message.
	// The message body is NOT being marshaled before sending and is passed to js.ValueOf directly.
	// Make sure the body is serializable to JSON.
	//  - https://pkg.go.dev/syscall/js#ValueOf
	QueueContentTypeJSON QueueContentType = "json"

	// QueueContentTypeV8 is currently treated the same as QueueContentTypeJSON.
	QueueContentTypeV8 QueueContentType = "v8"

	// QueueContentTypeText is used to send a message as a string.
	// Supported body types are string, []byte and io.Reader.
	QueueContentTypeText QueueContentType = "text"

	// QueueContentTypeBytes is used to send a message as a byte array.
	// Supported body types are string, []byte, and io.Reader.
	QueueContentTypeBytes QueueContentType = "bytes"
)

func (o QueueContentType) mapValue(val any) (js.Value, error) {
	switch o {
	case QueueContentTypeText:
		switch v := val.(type) {
		case string:
			return js.ValueOf(v), nil
		case []byte:
			return js.ValueOf(string(v)), nil
		case io.Reader:
			b, err := io.ReadAll(v)
			if err != nil {
				return js.Undefined(), fmt.Errorf("failed to read bytes from reader: %w", err)
			}
			return js.ValueOf(string(b)), nil
		default:
			return js.Undefined(), fmt.Errorf("invalid value type for text content type: %T", val)
		}

	case QueueContentTypeBytes:
		var b []byte
		switch v := val.(type) {
		case string:
			b = []byte(v)
		case []byte:
			b = v
		case io.Reader:
			var err error
			b, err = io.ReadAll(v)
			if err != nil {
				return js.Undefined(), fmt.Errorf("failed to read bytes from reader: %w", err)
			}
		default:
			return js.Undefined(), fmt.Errorf("invalid value type for bytes content type: %T", val)
		}

		ua := jsutil.NewUint8Array(len(b))
		js.CopyBytesToJS(ua, b)
		// accortind to docs, "bytes" type requires an ArrayBuffer to be sent, however practical experience shows that ArrayBufferView should
		// be used instead and with Uint8Array.buffer as a value, the send simply fails
		return ua, nil

	case QueueContentTypeJSON, QueueContentTypeV8:
		return js.ValueOf(val), nil
	}

	return js.Undefined(), fmt.Errorf("unknown content type: %s", o)
}
