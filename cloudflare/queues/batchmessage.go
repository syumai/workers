package queues

import (
	"errors"
	"syscall/js"

	"github.com/syumai/workers/internal/jsutil"
)

type BatchMessage struct {
	body    any
	options *sendOptions
}

// NewBatchMessage creates a single message to be batched before sending to a queue.
func NewBatchMessage(body any, opts ...SendOption) *BatchMessage {
	options := defaultSendOptions()
	for _, opt := range opts {
		opt(options)
	}
	return &BatchMessage{body: body, options: options}
}

func (m *BatchMessage) toJS() (js.Value, error) {
	if m == nil {
		return js.Undefined(), errors.New("message is nil")
	}

	jsValue, err := m.options.ContentType.mapValue(m.body)
	if err != nil {
		return js.Undefined(), err
	}

	obj := jsutil.NewObject()
	obj.Set("body", jsValue)
	obj.Set("options", m.options.toJS())

	return obj, nil
}
