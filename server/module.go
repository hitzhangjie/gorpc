package server

// ServerModule
type ServerModule interface {
	Start() error
	Register(*Service)
	Closed() <-chan struct{}
	Network() string
	Address() string
	Codec() string
}
