package registry

import (
	"sync"

	"github.com/hitzhangjie/go-rpc/server"
)

var (
	registries = map[string]Registry{}
	lock       = sync.RWMutex{}
)

// Registry registry interacts with the remote Nameing Service
//
// Register, register service
// UnRegister, unregister service
// GetService, get services by name, which may have more than one version
// ListServices, list all registered services
// Watcher, returns a watcher, which watches events on NamingService backend
type Registry interface {
	Register(service *server.Service, opts ...Option) error
	DeRegister(service *server.Service) error

	GetService(name string) ([]*server.Service, error)
	ListServices() ([]*server.Service, error)

	Watcher() (Watcher, error)
}

// Option registry option
type Option func(options *options)

type options struct{}

// Watcher watch event from remote Naming Service
type Watcher interface {
	Next() (*Result, error)
	Stop()
}

// Result watch result of event
type Result struct {
	Action ActionType
}

type ActionType = int

const (
	ActionTypeCreate = iota
	ActionTypeUpdate
	ActionTypeDelete
)

func RegisterRegistry(name string, registry Registry) {
	lock.Lock()
	registries[name] = registry
	lock.Unlock()
}

func GetRegistry(name string) Registry {
	lock.RLock()
	defer lock.RUnlock()

	v, ok := registries[name]
	if !ok {
		return nil
	}
	return v
}
