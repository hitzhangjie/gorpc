package transport

// Transport
type Transport interface {
	ListenAndServe() error
	//Register(*server.Service)
	Closed() <-chan struct{}
	Network() string
	Address() string
	Codec() string
}
