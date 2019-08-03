package server

import (
	"context"
	"github.com/hitzhangjie/go-rpc/router"
)

// Server represents a server instance (a server process), it can plug in ServerModules,
// including TcpServer, UdpServer, even Broker, to extend its ability.
type Server struct {
	ctx    context.Context
	opts   []*Option
	mods   []ServerModule
	closed chan (struct{})
	router *router.Router
}

// NewServer create new server with option
//
func NewServer(opts ...Option) (*Server, error) {
	s := &Server{
		ctx:    context.TODO(),
		opts:   []*Option{},
		mods:   []ServerModule{},
		router: router.NewRouter(),
	}
	return s, nil
}

func (s *Server) Start() {
	for _, m := range s.mods {
		go m.Start()
	}
	println("server started")

	<-s.closed
	println("server stopped")
}
