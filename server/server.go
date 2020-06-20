package server

import (
	"context"
	"sync"

	"github.com/hitzhangjie/gorpc-framework/router"
	"github.com/hitzhangjie/gorpc-framework/transport"
)

// Service represents a server instance (a server process),
//
// Any server can include more than one service, i.e, any server can be
// plugged into multile modules, like TcpServerTransport, UdpServerTransport, Broker, etc.
// By this way, we can implement more modules to extend server's abilities.
type Service struct {
	name   string
	ctx    context.Context
	cancel context.CancelFunc
	opts   *options

	trans    []transport.Transport
	transLck *sync.Mutex

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
		trans:     []transport.Transport{},
		transLck:  &sync.Mutex{},
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

func (s *Service) ListenAndServe(ctx context.Context, net, addr, codec string, opts ...Option) error {

	var (
		trans transport.Transport
		err   error
	)

	options := options{}
	for _, o := range opts {
		o(&options)
	}

	// transport options
	toptions := []transport.Option{}
	if options.Router != nil {
		toptions = append(toptions, transport.WithRouter(options.Router))
	}

	if net == "tcp" || net == "tcp4" || net == "tcp6" {
		trans, err = transport.NewTcpServerTransport(ctx, net, addr, codec, toptions...)
		if err != nil {
			return err
		}

	}

	if net == "udp" || net == "udp4" || net == "udp6" {
		trans, err = transport.NewUdpServerTransport(ctx, net, addr, codec, toptions...)
		if err != nil {
			return err
		}
	}

	s.transLck.Lock()
	s.trans = append(s.trans, trans)
	s.transLck.Unlock()

	go trans.ListenAndServe()

	return nil
}

func (s *Service) ServerModules() []transport.Transport {
	return s.trans
}

// ListenAndServe starts every Transport, after this, Service may be registered to remote naming service
//func (s *Service) Start() error {
//
//	if len(s.trans) == 0 {
//		return errors.New("server: no modules registered")
//	}
//
//	var err error
//
//	s.startOnce.Do(func() {
//		cherr := make(chan error, len(s.trans))
//		chok := make(chan struct{}, len(s.trans))
//
//		wg := sync.WaitGroup{}
//
//		// start all server trans
//		for _, m := range s.trans {
//			wg.Add(1)
//			go func() {
//				defer wg.Done()
//				if err := m.ListenAndServe(); err != nil && err != errServerCtxDone {
//					cherr <- err
//				} else {
//					chok <- struct{}{}
//				}
//			}()
//		}
//
//		// process stop signals to exit
//		go func() {
//			chexit := make(chan os.Signal)
//			signal.Notify(chexit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
//			<-chexit
//
//			println("server: got stop signal")
//			s.stop()
//			println("server: server stopped")
//		}()
//
//		// wait all server trans exit
//		wg.Wait()
//
//		select {
//		case err = <-cherr:
//			err = fmt.Errorf("server: inner module error: %v", err)
//		default:
//			println("server: ready to stop")
//		}
//		println("server: stopped")
//
//	})
//
//	return err
//}

// stop stop all server modules, server exit.
func (s *Service) stop() {

	s.stopOnce.Do(func() {

		s.cancel()

		wg := sync.WaitGroup{}
		for _, m := range s.trans {
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
