package queues

import (
	"errors"
	"fmt"
	"syscall/js"

	"github.com/syumai/workers/cloudflare/internal/cfruntimecontext"
	"github.com/syumai/workers/internal/jsutil"
)

type Producer struct {
	// queue - Objects that Queue API belongs to. Default is Global
	queue js.Value
}

// NewProducer creates a new Producer object to send messages to a queue.
// queueName is the name of the queue environment var to send messages to.
// In Cloudflare API documentation, this object represents the Queue.
//   - https://developers.cloudflare.com/queues/configuration/javascript-apis/#producer
func NewProducer(queueName string) (*Producer, error) {
	inst := cfruntimecontext.MustGetRuntimeContextEnv().Get(queueName)
	if inst.IsUndefined() {
		return nil, fmt.Errorf("%s is undefined", queueName)
	}
	return &Producer{queue: inst}, nil
}

// Send sends a single message to a queue. This function allows setting send options for the message.
// If no options are provided, the default options are used (QueueContentTypeJSON and no delay).
//
//   - https://developers.cloudflare.com/queues/configuration/javascript-apis/#producer
//   - https://developers.cloudflare.com/queues/configuration/javascript-apis/#queuesendoptions
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

// SendBatch sends multiple messages to a queue. This function allows setting options for	each message.
func (p *Producer) SendBatch(messages []*BatchMessage, opts ...BatchSendOption) error {
	if p.queue.IsUndefined() {
		return errors.New("queue object not found")
	}

	if len(messages) == 0 {
		return nil
	}

	var options *batchSendOptions
	if len(opts) > 0 {
		options = &batchSendOptions{}
		for _, opt := range opts {
			opt(options)
		}
	}

	jsArray := jsutil.NewArray(len(messages))
	for i, message := range messages {
		jsValue, err := message.toJS()
		if err != nil {
			return fmt.Errorf("failed to convert message %d to JS: %w", i, err)
		}
		jsArray.SetIndex(i, jsValue)
	}

	prom := p.queue.Call("sendBatch", jsArray, options.toJS())
	_, err := jsutil.AwaitPromise(prom)
	return err
}
