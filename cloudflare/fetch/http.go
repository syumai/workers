package fetch

import (
	"net/http"
)

// Do sends an HTTP request and returns an HTTP response
func (c *Client) Do(req *Request) (*http.Response, error) {
	return fetch(c.namespace, req.Request)
}
