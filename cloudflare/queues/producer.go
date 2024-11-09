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

func (p *Producer) SendText(content string, opts ...SendOption) error {
	return p.send(js.ValueOf(content), contentTypeText, opts...)
}

func (p *Producer) SendBytes(content []byte, opts ...SendOption) error {
	ua := jsutil.NewUint8Array(len(content))
	js.CopyBytesToJS(ua, content)
	// accortind to docs, "bytes" type requires an ArrayBuffer to be sent, however practical experience shows that ArrayBufferView should
	// be used instead and with Uint8Array.buffer as a value, the send simply fails
	return p.send(ua, contentTypeBytes, opts...)
}

func (p *Producer) SendJSON(content any, opts ...SendOption) error {
	return p.send(js.ValueOf(content), contentTypeJSON, opts...)
}

func (p *Producer) SendV8(content js.Value, opts ...SendOption) error {
	return p.send(content, contentTypeV8, opts...)
}

// send sends a single message to a queue. This function allows setting send options for the message.
// If no options are provided, the default options are used (QueueContentTypeJSON and no delay).
//
//   - https://developers.cloudflare.com/queues/configuration/javascript-apis/#producer
//   - https://developers.cloudflare.com/queues/configuration/javascript-apis/#queuesendoptions
func (p *Producer) send(body js.Value, contentType contentType, opts ...SendOption) error {
	options := &sendOptions{
		ContentType: contentType,
	}
	for _, opt := range opts {
		opt(options)
	}

	prom := p.queue.Call("send", body, options.toJS())
	_, err := jsutil.AwaitPromise(prom)
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
