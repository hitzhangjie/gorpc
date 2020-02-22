package registry

import (
	"github.com/hitzhangjie/go-rpc/server"
)

func init() {
	RegisterRegistry("noop", &NoopRegistry{})
}

type NoopRegistry struct {
}

func (n *NoopRegistry) Register(service *server.Service, opts ...Option) error {
	panic("implement me")
}

func (n *NoopRegistry) DeRegister(service *server.Service) error {
	panic("implement me")
}

func (n *NoopRegistry) GetService(name string) ([]*server.Service, error) {
	panic("implement me")
}

func (n *NoopRegistry) ListServices() ([]*server.Service, error) {
	panic("implement me")
}

func (n *NoopRegistry) Watcher() (Watcher, error) {
	panic("implement me")
}
