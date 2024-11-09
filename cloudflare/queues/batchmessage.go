package queues

import (
	"syscall/js"

	"github.com/syumai/workers/internal/jsutil"
)

type BatchMessage struct {
	body    js.Value
	options *sendOptions
}

func NewTextBatchMessage(content string, opts ...SendOption) *BatchMessage {
	return newBatchMessage(js.ValueOf(content), contentTypeText, opts...)
}

func NewBytesBatchMessage(content []byte, opts ...SendOption) *BatchMessage {
	return newBatchMessage(js.ValueOf(content), contentTypeBytes, opts...)
}

func NewJSONBatchMessage(content any, opts ...SendOption) *BatchMessage {
	return newBatchMessage(js.ValueOf(content), contentTypeJSON, opts...)
}

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
