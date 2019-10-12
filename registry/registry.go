package registry

import (
	"github.com/hitzhangjie/go-rpc/server"
)

type Registry interface {
	// Register register service
	Register(service *server.Service, opts ...RegisterOption) error
	// UnRegister unregister service
	UnRegister(service *server.Service) error
	// GetService get services by name, which may have more than one version
	GetService(name string) ([]*server.Service, error)
	// ListServices list all registered services
	ListServices() ([]*server.Service, error)
	// Watcher returns a watcher, which watches events on NamingService backend
	Watcher() (Watcher, error)
}

// RegisterOption
type RegisterOption func(options *RegisterOptions)

type RegisterOptions struct{}

type Watcher interface {
	Next() (*Result, error)
	Stop()
}

type Result struct {
	Action ActionType
}

type ActionType = int

const (
	actionTypeCreate = iota
	actionTypeUpdate
	actionTypeDelete
)
