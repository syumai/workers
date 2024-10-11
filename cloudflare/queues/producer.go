package queues

import (
	"errors"
	"fmt"
	"syscall/js"

	"github.com/syumai/workers/cloudflare/internal/cfruntimecontext"
	"github.com/syumai/workers/internal/jsutil"
)

type BatchMessage struct {
	body    any
	options *sendOptions
}

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

type Producer struct {
	// queue - Objects that Queue API belongs to. Default is Global
	queue js.Value
}

func NewProducer(queueName string) (*Producer, error) {
	inst := cfruntimecontext.MustGetRuntimeContextEnv().Get(queueName)
	if inst.IsUndefined() {
		return nil, fmt.Errorf("%s is undefined", queueName)
	}
	return &Producer{queue: inst}, nil
}

func (p *Producer) Send(content any, opts ...SendOption) error {
	if p.queue.IsUndefined() {
		return errors.New("queue object not found")
	}

	options := defaultSendOptions()
	for _, opt := range opts {
		opt(options)
	}

	jsValue, err := options.ContentType.mapValue(content)
	if err != nil {
		return err
	}

	prom := p.queue.Call("send", jsValue, options.toJS())
	_, err = jsutil.AwaitPromise(prom)
	return err
}

func (p *Producer) SendBatch(messages []*BatchMessage) error {
	if p.queue.IsUndefined() {
		return errors.New("queue object not found")
	}

	if len(messages) == 0 {
		return nil
	}

	jsArray := jsutil.NewArray(len(messages))
	for i, message := range messages {
		jsValue, err := message.toJS()
		if err != nil {
			return fmt.Errorf("failed to convert message %d to JS: %w", i, err)
		}
		jsArray.SetIndex(i, jsValue)
	}

	prom := p.queue.Call("sendBatch", jsArray)
	_, err := jsutil.AwaitPromise(prom)
	return err
}

func (p *Producer) SendJsonBatch(messages ...any) error {
	batch := make([]*BatchMessage, len(messages))
	for i, message := range messages {
		batch[i] = NewBatchMessage(message)
	}

	return p.SendBatch(batch)
}
