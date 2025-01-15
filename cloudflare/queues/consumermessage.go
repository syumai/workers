package queues

import (
	"fmt"
	"syscall/js"
	"time"

	"github.com/syumai/workers/internal/jsutil"
)

// ConsumerMessage represents a message of the batch received by the consumer.
//   - https://developers.cloudflare.com/queues/configuration/javascript-apis/#message
type ConsumerMessage struct {
	// instance - The underlying instance of the JS message object passed by the cloudflare
	instance js.Value

	// Id - The unique Cloudflare-generated identifier of the message
	Id string
	// Timestamp - The time when the message was enqueued
	Timestamp time.Time
	// Body - The message body. Could be accessed directly or using converting helpers as StringBody, BytesBody, IntBody, FloatBody.
	Body js.Value
	// Attempts - The number of times the message delivery has been retried.
	Attempts int
}

func newConsumerMessage(obj js.Value) (*ConsumerMessage, error) {
	timestamp, err := jsutil.DateToTime(obj.Get("timestamp"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse message timestamp: %v", err)
	}

	return &ConsumerMessage{
		instance:  obj,
		Id:        obj.Get("id").String(),
		Body:      obj.Get("body"),
		Attempts:  obj.Get("attempts").Int(),
		Timestamp: timestamp,
	}, nil
}

// Ack acknowledges the message as successfully delivered despite the result returned from the consuming function.
//   - https://developers.cloudflare.com/queues/configuration/javascript-apis/#message
func (m *ConsumerMessage) Ack() {
	m.instance.Call("ack")
}

// Retry marks the message to be re-delivered.
// The message will be retried after the optional delay configured with RetryOption.
func (m *ConsumerMessage) Retry(opts ...RetryOption) {
	var o *retryOptions
	if len(opts) > 0 {
		o = &retryOptions{}
		for _, opt := range opts {
			opt(o)
		}
	}

	m.instance.Call("retry", o.toJS())
}

func (m *ConsumerMessage) StringBody() (string, error) {
	if m.Body.Type() != js.TypeString {
		return "", fmt.Errorf("message body is not a string: %v", m.Body)
	}
	return m.Body.String(), nil
}

func (m *ConsumerMessage) BytesBody() ([]byte, error) {
	if m.Body.Type() != js.TypeObject ||
		!(m.Body.InstanceOf(jsutil.Uint8ArrayClass) || m.Body.InstanceOf(jsutil.Uint8ClampedArrayClass)) {
		return nil, fmt.Errorf("message body is not a byte array: %v", m.Body)
	}
	b := make([]byte, m.Body.Get("byteLength").Int())
	js.CopyBytesToGo(b, m.Body)
	return b, nil
}
