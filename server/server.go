package server

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/hitzhangjie/go-rpc/router"
)

// Server represents a server instance (a server process),
//
// Any server can include more than one service, i.e, any server can be
// plugged into multile modules, like TcpServerModule, UdpServerModule, Broker, etc.
// By this way, we can implement more modules to extend server's abilities.
type Server struct {
	ctx       context.Context
	cancel    context.CancelFunc
	opts      *options
	mods      []ServerModule
	router    *router.Router
	startOnce sync.Once
	stopOnce  sync.Once
	closed    chan (struct{})
}

// NewServer create new server with option
func NewServer(opts ...Option) (*Server, error) {

	ctx, cancel := context.WithCancel(context.TODO())

	s := &Server{
		ctx:       ctx,
		cancel:    cancel,
		opts:      &options{},
		mods:      []ServerModule{},
		router:    router.NewRouter(),
		startOnce: sync.Once{},
		stopOnce:  sync.Once{},
		closed:    make(chan struct{}, 1),
	}
	for _, o := range opts {
		o(s.opts)
	}
	return s, nil
}

// Start starts every ServerModule, after this, Service may be registered to remote naming service
func (s *Server) Start() error {

	if len(s.mods) == 0 {
		return errors.New("server: no modules registered")
	}

	var err error

	s.startOnce.Do(func() {
		cherr := make(chan error, len(s.mods))
		chok := make(chan struct{}, len(s.mods))

		wg := sync.WaitGroup{}

		// start all server mods
		for _, m := range s.mods {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := m.Start(); err != nil && err != errServerCtxDone {
					cherr <- err
				} else {
					chok <- struct{}{}
				}
			}()
		}

		// process stop signals to exit
		go func() {
			chexit := make(chan os.Signal)
			signal.Notify(chexit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
			<-chexit

			println("server: got stop signal")
			s.stop()
			println("server: server stopped")
		}()

		// wait all server mods exit
		wg.Wait()

		select {
		case err = <-cherr:
			err = fmt.Errorf("server: inner module error: %v", err)
		default:
			println("server: ready to stop")
		}
		println("server: stopped")

	})

	return err
}

// stop stop all server modules, server exit.
func (s *Server) stop() {

	s.stopOnce.Do(func() {

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

		close(s.closed)
	})
}

// Closed return whether server is Closed or not
func (s *Server) Closed() chan struct{} {
	return s.closed
}
