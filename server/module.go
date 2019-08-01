package server

// ServerModule
type ServerModule interface {
	Start()
	Stop()
	Register(*Server)
}
