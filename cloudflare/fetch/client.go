package fetch

import (
	"syscall/js"

	"github.com/syumai/workers/internal/jsutil"
)

// Client is an HTTP client.
type Client struct {
	// namespace - Objects that Fetch API belongs to. Default is Global
	namespace js.Value
}

// NewClient returns new Client
func NewClient() *Client {
	return &Client{
		namespace: jsutil.Global,
	}
}
