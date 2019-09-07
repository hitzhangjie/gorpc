package server

import (
	"context"
	"github.com/hitzhangjie/go-rpc/router"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Server represents a server instance (a server process), it can plug in ServerModules,
// including TcpServer, UdpServer, even Broker, to extend its ability.
type Server struct {
	ctx    context.Context
	cancel context.CancelFunc
	opts   []*Option
	mods   []ServerModule
	router *router.Router
	once   sync.Once
	closed chan (struct{})
}

// NewServer create new server with option
//
func NewServer(opts ...Option) (*Server, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	s := &Server{
		ctx:    ctx,
		cancel: cancel,
		opts:   []*Option{},
		mods:   []ServerModule{},
		router: router.NewRouter(),
		once:   sync.Once{},
		closed: make(chan struct{}, 1),
	}
	return s, nil
}

func (s *Server) Start() {
	for _, m := range s.mods {
		go m.Start()
	}
	println("server started")

	// process recv following signals to exit
	chexit := make(chan os.Signal)
	signal.Notify(chexit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	<-chexit

	s.Stop()

	println("server stopped")
}

func (s *Server) Stop() {

	s.cancel()

	for _, m := range s.mods {
		if !m.Closed() {
			return false
		}
	}

	s.once.Do(func() {
		close(s.closed)
	})
}

func (s *Server) Closed() chan struct{} {
	return s.closed
}
