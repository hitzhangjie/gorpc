package server

// In gorpc framework, you have more than one method to quickly start a server.
//
// 	Method1: use `gorpc.ListenAndServe()`
//	This method read conf/service.ini and start `ServerModules` defined in `service.ini`.
//  Also, it'll register to remote naming service if information of naming service provided.
//
// 	method2: explicity create a `server` and
//	```go
// 	service := gorpc.NewService(name)
// 	service.RegisterServer(&server)
//	```
//
// method3:
//		```go
//		tcpSvr := NewTcpServerModule(...)
//		udpSvr := NewUdpServerModule(...)
//		service := gorpc.NewService(name)
//		service.RegisterModule(tcpSvr)
//		```
type Service struct {
	name    string
	version string
	server  *Server
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

func (s *Service) Start() error {

	if s.server == nil {
		return errServerNotInit
	}

	// start
	if err := s.server.Start(); err != nil {
		return err
	}
	// register
	if err := s.RegisterNaming(); err != nil {
		return err
	}
	// de-register
	<-s.server.Closed()
	if err := s.DeRegisterNaming(); err != nil {
		return err
	}

	return nil
}

func (s *Service) RegisterServer(svr *Server) {
	panic("implement me")
}

func (s *Service) RegisterNaming() error {
	panic("implement me")
}

func (s *Service) DeRegisterNaming() error {
	panic("implement me")
}
