package server

import (
	"context"
)

// Server it represents a server instance (a server process), its ability is extensible via
// pluggable many ServerModule implementation, including TcpServer, UdpServer, even Broker.
type Server struct {
	ctx  context.Context
	opts []*Option
	mods []ServerModule
}

const (
	optionsMaxLen = 16
	modulesMaxLen = 16
)

func NewServer(opts ...Option) (*Server, error) {
	s := &Server{
		ctx:  context.TODO(),
		opts: make([]*Option, optionsMaxLen),
		mods: make([]ServerModule, modulesMaxLen),
	}
	return s, nil
}

func (s *Server) Register(m ServerModule) {
	s.mods = append(s.mods, m)
}

func (s *Server) Start() {
	for _, m := range s.mods {
		go m.Start()
	}
}

type Option struct {
}

// ServerModule
type ServerModule interface {
	Start()
	Stop()
}

// TcpServer
type TcpServer struct {
	svr *Server
}

func NewTcpServer(server *Server) ServerModule {
	s := &TcpServer{
		svr: server,
	}
	return s
}

func (s *TcpServer) Start() {}

func (s *TcpServer) Stop() {
}

// UdpServer
type UdpServer struct {
	svr *Server
}

func NewUdpServer(server *Server) ServerModule {
	s := &UdpServer{
		svr: server,
	}
	return s
}

func (s *UdpServer) Start() {}

func (s *UdpServer) Stop() {}
