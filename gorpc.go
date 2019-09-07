// package gorpc defines the Service and provides some wrappers to quickly start gorpc service.
package gorpc

import (
	"github.com/hitzhangjie/go-rpc/server"
)

// Service represents a service running somewhere, maybe deployed in multi-hosts or in Cloud.
//
// Service vs Server, these terms, as I see, they're two different views of running service.
// - Server, it's a running process or instance.
// - Service, it's deployed in public environment and provides service via `naming mechanism`.
//
// In go-rpc, you can start a `Server` via a `server.NewTcpServer()` or `server.NewUdpServer()`,
// If you want to register this service to remote naming service, you can use:
//
// 	method1:
//		```go
// 		gorpc.NewService()
//		```
// 	method2:
//		```go
// 		service := gorpc.NewService(name)
// 		service.RegisterServer(&server)
//		```
//
// method3:
//		```go
//		tcpSvr := NewTcpServer(...)
//		udpSvr := NewUdpServer(...)
//		service := gorpc.NewService(name)
//		service.RegisterModule(tcpSvr)
//		```
type Service struct {
	name    string
	version string
	server  *server.Server
}

// NewService create a new service
func NewService(name string) *Service {
	s := &Service{
		name:    name,
		version: "0.0.1",
		server:  nil,
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
	if s.server == nil {
		panic(errServerNotInit)
	}
	s.server.Start()

	<- s.server.Closed()
}

func (s *Service) RegisterServer(svr *server.Server) {
	panic("implement me")
}

func (s *Service) RegisterSModule(mod *server.ServerModule) {
	panic("implement me")
}
