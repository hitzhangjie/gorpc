package registry

import (
	"github.com/hitzhangjie/gorpc/server"
)

func init() {
	RegisterRegistry("noop", &NoopRegistry{})
}

// NoopRegistry noop registry implemention
type NoopRegistry struct {
}

func (n *NoopRegistry) Register(service *server.Service, opts ...Option) error {
	return nil
}

func (n *NoopRegistry) DeRegister(service *server.Service) error {
	return nil
}

func (n *NoopRegistry) GetService(name string) ([]*server.Service, error) {
	return nil, nil
}

func (n *NoopRegistry) ListServices() ([]*server.Service, error) {
	return nil, nil
}

func (n *NoopRegistry) Watcher() (Watcher, error) {
	return nil, nil
}
