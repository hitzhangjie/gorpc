package client

import (
	"github.com/hitzhangjie/go-rpc/client/pool"
	"github.com/hitzhangjie/go-rpc/client/selector"
	"github.com/hitzhangjie/go-rpc/client/transport"
	"github.com/hitzhangjie/go-rpc/codec"
	"time"
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
	TCP
	UNIX
)

func (t TransportType) String() string {
	switch t {
	case UDP:
		return "udp"
	case TCP:
		return "tcp"
	case UNIX:
		return "unix"
	default:
		return ""
	}
}

func (t TransportType) Valid() bool {
	if t == UDP || t == TCP || t == UNIX {
		return true
	}
	return false
}

// WithTransportType specify the transport type, support UDP, TCP, Unix
func WithTransportType(typ TransportType) Option {
	return func(c *client) {
		c.TransType = typ
		switch typ {
		case TCP:
			c.Transport = &transport.TcpTransport{
				ConnPool: pool.ConnPool{
					MinIdle:         2,
					MaxIdle:         4,
					MaxActive:       8,
					Wait:            true,
					IdleTimeout:     time.Minute * 5,
					MaxConnLifetime: time.Hour * 1,
				},
				Codec: nil,
			}
		case UDP:
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
