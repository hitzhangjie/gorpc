package gorpc

import (
	"gorpc/server"
)

// Service represents a service running somewhere, maybe multiple hosts or Cloud.
//
// Service defines a service running in multiple nodes, while server is an instance
// running in one node.
type Service struct {
	name    string
	version string
	server  *server.Server
}

// NewService create a new service
func NewService(name string) *Service {
	s := &Service{
		name: name,
	}
	return s
}

// Version set service version, each service can have serveral versions' api.
func (s *Service) Version(v string) *Service {
	s.version = v
	return s
}

func (s *Service) Handle(service interface{}) {
	svr := server.NewServer()
}

func (s *Service) Start() {
}
