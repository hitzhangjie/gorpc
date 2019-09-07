package server

// ServerModule
type ServerModule interface {
	Start() error
	Stop()
	Register(*Server)
	Closed() <-chan struct{}
}
