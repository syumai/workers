package sockets

import (
	"context"
	"net"
	"syscall/js"

	"github.com/syumai/workers/cloudflare/internal/cfruntimecontext"
	"github.com/syumai/workers/internal/jsutil"
)

type Dialer struct {
	connect js.Value
	opts    *SocketOptions
	ctx     context.Context
}

type SocketOptions struct {
	SecureTransport string `json:"secureTransport"`
	AllowHalfOpen   bool   `json:"allowHalfOpen"`
}

// NewDialer
func NewDialer(ctx context.Context, options *SocketOptions) (*Dialer, error) {
	connect, err := cfruntimecontext.GetRuntimeContextValue(ctx, "connect")
	if err != nil {
		return nil, err
	}
	return &Dialer{connect: connect, opts: options, ctx: ctx}, nil
}

func (d *Dialer) Dial(ctx context.Context, network, addr string) (net.Conn, error) {
	switch network {
	case "tcp":
	default:
		panic("not implemented")
	}
	optionsObj := jsutil.NewObject()
	if d.opts != nil {
		if d.opts.AllowHalfOpen {
			optionsObj.Set("allowHalfOpen", true)
		}
		if d.opts.SecureTransport != "" {
			optionsObj.Set("secureTransport", d.opts.SecureTransport)
		}
	}
	sock := &TCPSocket{}
	sock.socket = d.connect.Invoke(addr, optionsObj)
	sock.options = d.opts
	sock.init(d.ctx)
	return sock, nil
}
