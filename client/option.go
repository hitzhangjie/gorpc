package client

import (
	"github.com/hitzhangjie/gorpc-framework/client/selector"
	"github.com/hitzhangjie/gorpc-framework/client/transport"
	"github.com/hitzhangjie/gorpc-framework/codec"
)

type Option func(*client)

// WithAddress specify the address that client requests
func WithAddress(addr string) Option {
	return func(c *client) {
		c.Addr = addr
	}
}

// TransportType options
type TransportType int

const (
	UDP = iota
	UDP4
	UDP6
	TCP
	TCP4
	TCP6
	UNIX
)

func (t TransportType) String() string {
	switch t {
	case UDP:
		return "udp"
	case UDP4:
		return "udp4"
	case UDP6:
		return "udp6"
	case TCP:
		return "tcp"
	case TCP4:
		return "tcp4"
	case TCP6:
		return "tcp6"
	case UNIX:
		return "unix"
	default:
		return ""
	}
}

func (t TransportType) Valid() bool {
	if t == UDP || t == UDP4 || t == UDP6 ||
		t == TCP || t == TCP4 || t == TCP6 ||
		t == UNIX {
		return true
	}
	return false
}

// WithTransportType specify the transport type, support UDP, TCP, Unix
func WithTransportType(typ TransportType) Option {
	return func(c *client) {
		c.TransType = typ
		switch typ {
		case TCP, TCP4, TCP6:
			c.Transport = &transport.TcpTransport{
				Pool:  defaultPoolFactory,
				Codec: codec.ClientCodec("whisper"),
			}
		case UDP, UDP4, UDP6:
			c.Transport = &transport.UdpTransport{} //fixme
		case UNIX:
			c.Transport = &transport.UnixTransport{} //fixme
		}
	}
}

// RpcType options
type RpcType int

const (
	SendOnly                = iota // 只发不收
	SendRecv                       // 一发一收
	SendRecvMultiplex              // 多路复用方式一发一收，发挥双工优势
	SendStreamOnly                 // 流式请求
	SendStreamAndRecv              // 流式请求
	SendAndRecvStream              // 流式请求
	SendStreamAndRecvStream        // 流式请求
)

// WithRpcType specify the rpc type, support SendOnly, SendRecv, SendRecvWithMultiplex, etc.
func WithRpcType(typ RpcType) Option {
	return func(c *client) {
		c.RpcType = typ
	}
}

// WithSelector specify the selector
func WithSelector(selector selector.Selector) Option {
	return func(c *client) {
		c.Selector = selector
	}
}

// WithCodec specify the codec
func WithCodec(name string) Option {
	return func(c *client) {
		c.Codec = codec.ClientCodec(name)
	}
}
