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
	// fixme service应该生成桩代码，里面定义好各个rpc名字与对应handler的映射关系
	// 类似于完成goneat中AddExec的操作！
	// 考虑不同业务协议的问题，可能有些业务协议使用的是int类型的cmd来区分接口，因此在协议之上还要抽象一个层，通过req体到rpc名字的映射，
	// func RpcName(req interface{}) string
}

func (s *Service) Start() {
}
