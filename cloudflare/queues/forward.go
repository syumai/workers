// Deprecated: use github.com/syumai/workers-go/cloudflare/queues instead.
package queues

import (
	src "github.com/syumai/workers-go/cloudflare/queues"
)

type BatchSendOption = src.BatchSendOption

var Consume = src.Consume
var ConsumeNonBlock = src.ConsumeNonBlock

type Consumer = src.Consumer
type Message = src.Message
type MessageBatch = src.MessageBatch
type MessageSendRequest = src.MessageSendRequest

var NewBytesMessageSendRequest = src.NewBytesMessageSendRequest
var NewJSONMessageSendRequest = src.NewJSONMessageSendRequest
var NewProducer = src.NewProducer
var NewTextMessageSendRequest = src.NewTextMessageSendRequest
var NewV8MessageSendRequest = src.NewV8MessageSendRequest

type Producer = src.Producer
type RetryOption = src.RetryOption
type SendOption = src.SendOption

var WithBatchDelaySeconds = src.WithBatchDelaySeconds
var WithDelaySeconds = src.WithDelaySeconds
var WithRetryDelay = src.WithRetryDelay
