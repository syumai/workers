package sockets

import (
	"context"
	"net"

	"github.com/syumai/workers/internal/jsutil"

	"github.com/syumai/workers/cloudflare/internal/cfruntimecontext"
)

type SecureTransport string

const (
	SecureTransportOn       SecureTransport = "on"
	SecureTransportOff      SecureTransport = "off"
	SecureTransportStartTLS SecureTransport = "starttls"
)

type SocketOptions struct {
	SecureTransport SecureTransport `json:"secureTransport"`
	AllowHalfOpen   bool            `json:"allowHalfOpen"`
}

func Connect(ctx context.Context, addr string, opts *SocketOptions) (net.Conn, error) {
	connect, err := cfruntimecontext.GetRuntimeContextValue(ctx, "connect")
	if err != nil {
		return nil, err
	}
	optionsObj := jsutil.NewObject()
	if opts != nil {
		if opts.AllowHalfOpen {
			optionsObj.Set("allowHalfOpen", true)
		}
		if opts.SecureTransport != "" {
			optionsObj.Set("secureTransport", opts.SecureTransport)
		}
	}
	sockVal := connect.Invoke(addr, optionsObj)
	return newSocket(ctx, sockVal), nil
}
