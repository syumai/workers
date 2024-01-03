package socket

import (
	"context"
	"io"
	"net"
	"os"
	"syscall/js"
	"time"

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

func (sock *TCPSocket) init(ctx context.Context) {
	sock.SetDeadline(time.Now().Add(999999 * time.Hour))
	sock.writer = sock.socket.Get("writable").Call("getWriter")
	sock.reader = sock.socket.Get("readable").Call("getReader")
	sock.rd = jsutil.ConvertStreamReaderToReader(sock.reader)
	sock.ctx, sock.cancel = context.WithCancel(ctx)
	return
}

type TCPSocket struct {
	socket js.Value
	writer js.Value
	reader js.Value

	rd io.Reader

	options *SocketOptions

	readDeadline  time.Time
	writeDeadline time.Time

	ctx    context.Context
	cancel context.CancelFunc
}

func (t *TCPSocket) Socket() js.Value {
	return t.socket
}

var _ net.Conn = (*TCPSocket)(nil)

// Read reads data from the connection.
// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
func (t *TCPSocket) Read(b []byte) (n int, err error) {
	ctx, cancel := context.WithDeadline(t.ctx, t.readDeadline)
	defer cancel()
	done := make(chan struct{})
	go func() {
		n, err = t.rd.Read(b)
		close(done)
	}()
	select {
	case <-done:
		return
	case <-ctx.Done():
		return 0, os.ErrDeadlineExceeded
	}
}

// Write writes data to the connection.
// Write can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetWriteDeadline.
func (t *TCPSocket) Write(b []byte) (n int, err error) {
	ctx, cancel := context.WithDeadline(t.ctx, t.writeDeadline)
	defer cancel()
	done := make(chan struct{})
	go func() {
		arr := jsutil.NewUint8Array(len(b))
		js.CopyBytesToJS(arr, b)
		_, err = jsutil.AwaitPromise(t.writer.Call("write", arr))
		if err == nil {
			n = len(b)
		}
		close(done)
	}()
	select {
	case <-done:
		return
	case <-ctx.Done():
		return 0, os.ErrDeadlineExceeded
	}
}

// StartTls will call startTls on the socket
func (t *TCPSocket) StartTls() *TCPSocket {
	sock := &TCPSocket{}
	sock.socket = t.socket.Call("startTls")
	sock.options = t.options
	sock.init(t.ctx)
	return sock
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (t *TCPSocket) Close() error {
	t.cancel()
	t.socket.Call("close")
	return nil
}

// CloseRead closes the read side of the connection.
func (t *TCPSocket) CloseRead() error {
	t.reader.Call("close")
	return nil
}

// CloseWrite closes the write side of the connection.
func (t *TCPSocket) CloseWrite() error {
	t.writer.Call("close")
	return nil
}

// LocalAddr returns the local network address, if known.
func (t *TCPSocket) LocalAddr() net.Addr {
	return nil
}

// RemoteAddr returns the remote network address, if known.
func (t *TCPSocket) RemoteAddr() net.Addr {
	return nil
}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
//
// A deadline is an absolute time after which I/O operations
// fail instead of blocking. The deadline applies to all future
// and pending I/O, not just the immediately following call to
// Read or Write. After a deadline has been exceeded, the
// connection can be refreshed by setting a deadline in the future.
//
// If the deadline is exceeded a call to Read or Write or to other
// I/O methods will return an error that wraps os.ErrDeadlineExceeded.
// This can be tested using errors.Is(err, os.ErrDeadlineExceeded).
// The error's Timeout method will return true, but note that there
// are other possible errors for which the Timeout method will
// return true even if the deadline has not been exceeded.
//
// An idle timeout can be implemented by repeatedly extending
// the deadline after successful Read or Write calls.
//
// A zero value for t means I/O operations will not time out.
func (t *TCPSocket) SetDeadline(deadline time.Time) error {
	t.SetReadDeadline(deadline)
	t.SetWriteDeadline(deadline)
	return nil
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (t *TCPSocket) SetReadDeadline(deadline time.Time) error {
	t.readDeadline = deadline
	return nil
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (t *TCPSocket) SetWriteDeadline(deadline time.Time) error {
	t.writeDeadline = deadline
	return nil
}
