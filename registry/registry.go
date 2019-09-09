package registry

import (
	"git.code.oa.com/trpc-go/trpc-go/server"
)

type Registry interface {
	Register(service *server.Service, opts ...RegisterOption) error
	DeRegister(service *server.Service) error
	GetService(name string) ([]*server.Service, error)
	ListServices() ([]*server.Service, error)
	Watcher() (Watcher, error)
}

type RegisterOption func()

type RegisterOptions struct {
}

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
