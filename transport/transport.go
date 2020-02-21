package transport

import "github.com/hitzhangjie/go-rpc/server"

// Transport
type Transport interface {
	ListenAndServe() error
	Register(*server.Service)
	Closed() <-chan struct{}
	Network() string
	Address() string
	Codec() string
}
