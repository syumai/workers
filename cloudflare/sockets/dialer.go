package sockets

import (
	"context"
	"net"

	"github.com/syumai/workers/cloudflare/internal/cfruntimecontext"
	"github.com/syumai/workers/internal/jsutil"
)

type SocketOptions struct {
	SecureTransport string `json:"secureTransport"`
	AllowHalfOpen   bool   `json:"allowHalfOpen"`
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
	sock := &TCPSocket{}
	sock.socket = connect.Invoke(addr, optionsObj)
	sock.options = opts
	sock.init(ctx)
	return sock, nil
}
