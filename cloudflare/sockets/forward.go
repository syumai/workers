// Deprecated: use github.com/syumai/workers-go/cloudflare/sockets instead.
package sockets

import (
	src "github.com/syumai/workers-go/cloudflare/sockets"
)

var Connect = src.Connect

type SecureTransport = src.SecureTransport

const SecureTransportOff = src.SecureTransportOff
const SecureTransportOn = src.SecureTransportOn
const SecureTransportStartTLS = src.SecureTransportStartTLS

type Socket = src.Socket
type SocketOptions = src.SocketOptions
