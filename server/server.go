package server

import (
	"context"
	"fmt"
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
func NewServer(opts ...Option) (*Server, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	s := &Server{
		ctx:    ctx,
		cancel: cancel,

		opts: []*Option{},

		mods:   []ServerModule{},
		router: router.NewRouter(),

		once:   sync.Once{},
		closed: make(chan struct{}, 1),
	}
	return s, nil
}

// Start starts every ServerModule, after this, Service may be registered to remote naming service
func (s *Server) Start() error {
	cherr := make(chan error, len(s.mods))
	chok := make(chan struct{}, len(s.mods))

	wg := sync.WaitGroup{}
	for _, m := range s.mods {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := m.Start(); err != nil {
				cherr <- err
			} else {
				chok <- struct{}{}
			}
		}()
	}
	wg.Wait()

	select {
	case err := <-cherr:
		return fmt.Errorf("server start error: %v", err)
	default:
		println("server started")
	}

	// process recv following signals to exit
	go func() {
		chexit := make(chan os.Signal)
		signal.Notify(chexit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
		<-chexit

		s.Stop()
		println("server stopped")
	}()

	return nil
}

func (s *Server) Stop() {

	s.cancel()

	wg := sync.WaitGroup{}
	for _, m := range s.mods {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-m.Closed()
		}()
	}
	wg.Wait()

	s.once.Do(func() {
		close(s.closed)
	})
}

func (s *Server) Closed() chan struct{} {
	return s.closed
}
