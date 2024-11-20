package queues

import (
	"fmt"
	"syscall/js"
)

// ConsumerMessageBatch represents a batch of messages received by the consumer. The size of the batch is determined by the
// worker configuration.
//   - https://developers.cloudflare.com/queues/configuration/configure-queues/#consumer
//   - https://developers.cloudflare.com/queues/configuration/javascript-apis/#messagebatch
type ConsumerMessageBatch struct {
	// instance - The underlying instance of the JS message object passed by the cloudflare
	instance js.Value

	// Queue - The name of the queue from which the messages were received
	Queue string

	// Messages - The messages in the batch
	Messages []*ConsumerMessage
}

func newConsumerMessageBatch(obj js.Value) (*ConsumerMessageBatch, error) {
	msgArr := obj.Get("messages")
	messages := make([]*ConsumerMessage, msgArr.Length())
	for i := 0; i < msgArr.Length(); i++ {
		m, err := newConsumerMessage(msgArr.Index(i))
		if err != nil {
			return nil, fmt.Errorf("failed to parse message %d: %v", i, err)
		}
		messages[i] = m
	}

	return &ConsumerMessageBatch{
		instance: obj,
		Queue:    obj.Get("queue").String(),
		Messages: messages,
	}, nil
}

// AckAll acknowledges all messages in the batch as successfully delivered despite the result returned from the consuming function.
//   - https://developers.cloudflare.com/queues/configuration/javascript-apis/#messagebatch
func (b *ConsumerMessageBatch) AckAll() {
	b.instance.Call("ackAll")
}

// RetryAll marks all messages in the batch to be re-delivered.
// The messages will be retried after the optional delay configured with RetryOption.
//   - https://developers.cloudflare.com/queues/configuration/javascript-apis/#messagebatch
func (b *ConsumerMessageBatch) RetryAll(opts ...RetryOption) {
	var o *retryOptions
	if len(opts) > 0 {
		o = &retryOptions{}
		for _, opt := range opts {
			opt(o)
		}
	}

	b.instance.Call("retryAll", o.toJS())
}
