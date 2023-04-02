package fetch

import (
	"net/http"
	"syscall/js"

	"github.com/syumai/workers/internal/jsutil"
)

// Client is an HTTP client.
type Client struct {
	*http.Client

	// namespace - Objects that Fetch API belongs to. Default is Global
	namespace js.Value
}

// NewClient returns new Client
func NewClient() *Client {
	return &Client{
		Client:    &http.Client{},
		namespace: jsutil.Global,
	}
}
