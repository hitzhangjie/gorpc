package server

// Transport
type Transport interface {
	ListenAndServe() error
	Register(*Service)
	Closed() <-chan struct{}
	Network() string
	Address() string
	Codec() string
}
