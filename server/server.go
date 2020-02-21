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

// Service represents a server instance (a server process),
//
// Any server can include more than one service, i.e, any server can be
// plugged into multile modules, like TcpServerModule, UdpServerModule, Broker, etc.
// By this way, we can implement more modules to extend server's abilities.
type Service struct {
	name   string
	ctx    context.Context
	cancel context.CancelFunc
	opts   *options

	mods []ServerModule
	lock *sync.Mutex

	router    *router.Router
	startOnce sync.Once
	stopOnce  sync.Once
	closed    chan (struct{})
}

// NewService create new server with option
func NewService(name string, opts ...Option) *Service {

	s := &Service{
		name:      name,
		opts:      &options{},
		mods:      []ServerModule{},
		lock:      &sync.Mutex{},
		router:    router.NewRouter(),
		startOnce: sync.Once{},
		stopOnce:  sync.Once{},
		closed:    make(chan struct{}, 1),
	}
	s.ctx, s.cancel = context.WithCancel(context.TODO())

	for _, o := range opts {
		o(s.opts)
	}
	return s
}

func (s *Service) AddServerModule(net, addr, codec string, opts ...Option) error {

	var (
		mod ServerModule
		err error
	)

	if net == "tcp" || net == "tcp4" || net == "tcp6" {
		mod, err = NewTcpServerModule(net, addr, codec, opts...)
		if err != nil {
			return err
		}

	}

	if net == "udp" || net == "udp4" || net == "udp6" {
		mod, err = NewUdpServerModule(net, addr, codec, opts...)
		if err != nil {
			return err
		}
	}

	s.lock.Lock()
	s.mods = append(s.mods, mod)
	s.lock.Unlock()

	return nil
}

func (s *Service) ServerModules() []ServerModule {
	return s.mods
}

// Start starts every ServerModule, after this, Service may be registered to remote naming service
func (s *Service) Start() error {

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
func (s *Service) stop() {

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
