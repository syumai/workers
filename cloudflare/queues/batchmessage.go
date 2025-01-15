package queues

import (
	"syscall/js"

	"github.com/syumai/workers/internal/jsutil"
)

// FIXME: rename to MessageSendRequest
// see: https://developers.cloudflare.com/queues/configuration/javascript-apis/#messagesendrequest
type BatchMessage struct {
	body    js.Value
	options *sendOptions
}

// NewTextBatchMessage creates a single text message to be batched before sending to a queue.
func NewTextBatchMessage(content string, opts ...SendOption) *BatchMessage {
	return newBatchMessage(js.ValueOf(content), contentTypeText, opts...)
}

// NewBytesBatchMessage creates a single byte array message to be batched before sending to a queue.
func NewBytesBatchMessage(content []byte, opts ...SendOption) *BatchMessage {
	return newBatchMessage(js.ValueOf(content), contentTypeBytes, opts...)
}

// NewJSONBatchMessage creates a single JSON message to be batched before sending to a queue.
func NewJSONBatchMessage(content any, opts ...SendOption) *BatchMessage {
	return newBatchMessage(js.ValueOf(content), contentTypeJSON, opts...)
}

// NewV8BatchMessage creates a single raw JS value message to be batched before sending to a queue.
func NewV8BatchMessage(content js.Value, opts ...SendOption) *BatchMessage {
	return newBatchMessage(content, contentTypeV8, opts...)
}

// newBatchMessage creates a single message to be batched before sending to a queue.
func newBatchMessage(body js.Value, contentType contentType, opts ...SendOption) *BatchMessage {
	options := sendOptions{
		ContentType: contentType,
	}
	for _, opt := range opts {
		opt(&options)
	}
	return &BatchMessage{body: body, options: &options}
}

func (m *BatchMessage) toJS() js.Value {
	obj := jsutil.NewObject()
	obj.Set("body", m.body)
	obj.Set("options", m.options.toJS())
	return obj
}
