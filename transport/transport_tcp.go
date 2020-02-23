package transport

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"github.com/hitzhangjie/go-rpc/codec"
	"github.com/hitzhangjie/go-rpc/errs"
	"github.com/hitzhangjie/go-rpc/router"
)

// TcpServerTransport
type TcpServerTransport struct {
	ctx    context.Context
	cancel context.CancelFunc

	net  string
	addr string

	codec  codec.Codec
	reader *TcpMessageReader

	//reqChan chan interface{}
	//reqChan chan codec.Session
	//rspChan chan codec.Session

	wg sync.WaitGroup

	once   sync.Once
	closed chan struct{}

	opts *options
}

const (
	tcpServerRspChanMaxLength = 1024
)

type options struct {
	Router *router.Router
}

type Option func(*options)

func WithRouter(router *router.Router) Option {
	return func(opts *options) {
		opts.Router = router
	}
}

func NewTcpServerTransport(ctx context.Context, net, addr, codecName string, opts ...Option) (Transport, error) {

	ctx, cancel := context.WithCancel(ctx)
	c := codec.ServerCodec(codecName)

	s := &TcpServerTransport{
		ctx:    ctx,
		cancel: cancel,
		net:    net,
		addr:   addr,
		codec:  c,
		reader: NewTcpMessageReader(c),
		//rspChan: make(chan codec.Session, tcpServerRspChanMaxLength),
		once:   sync.Once{},
		closed: make(chan struct{}, 1),
		opts:   &options{},
	}
	for _, o := range opts {
		o(s.opts)
	}
	return s, nil
}

func (s *TcpServerTransport) ListenAndServe() error {

	l, err := net.Listen(s.net, s.addr)
	if err != nil {
		return err
	}
	err = s.serve(l)

	s.once.Do(func() {
		close(s.closed)
	})

	return err
}

//func (s *TcpServerTransport) Register(svr *server.Service) {
//	s.ctx, s.cancel = context.WithCancel(svr.ctx)
//	s.opts.router = svr.router
//	svr.trans = append(svr.trans, s)
//}

func (s *TcpServerTransport) serve(l net.Listener) error {

	defer func() {
		s.cancel()
		l.Close()
	}()

	for {
		// check whether server Closed
		select {
		case <-s.ctx.Done():
			return errs.ErrServerCtxDone
		default:
		}
		// accept tcpconn
		conn, err := l.Accept()
		if err != nil {
			if e, ok := err.(net.Error); ok && e.Temporary() {
				time.Sleep(time.Millisecond * 10)
				continue
			}
			return err
		}

		ep := TcpEndPoint{
			conn,
			make(chan interface{}, 1024),
			make(chan interface{}, 1024),
			s.reader,
			nil,
			nil,
			tcpBufferPool.Get().([]byte),
		}
		ep.ctx, ep.cancel = context.WithCancel(s.ctx)

		go s.proc(ep.reqCh, ep.rspCh)

		go ep.Read()
		go ep.Write()
	}
}

func (s *TcpServerTransport) proc(reqCh <-chan interface{}, rspCh chan<- interface{}) {

	builder := codec.GetSessionBuilder(s.reader.codec.Name())

	for {
		select {
		case <-s.ctx.Done():
			s.cancel()
			return
		case req, ok := <-reqCh:
			if !ok {
				log.Printf("one endpoint.reqCh is Closed")
				return
			}
			// build session
			session, err := builder.Build(req)
			if err != nil {
				log.Fatalf("tcp build session error:%v", err)
				continue
			}
			// fixme using workerpool instead of goroutine
			r := s.opts.Router
			if r == nil {
				log.Fatalf("tcp router not initialized")
			}

			go func() {
				// find route
				handle, err := r.Route(session.RPCName())
				if err != nil {
					session.SetErrorResponse(err)
					log.Fatalf("tcp router route error:%v", err)
					return
				}
				// pass session+req to handlefunc
				ctx := codec.ContextWithSession(s.ctx, session)
				rsp, err := handle(ctx, req)
				if err != nil {
					session.SetErrorResponse(err)
					log.Fatalf("tcp handle func error:%v, rsp:%v", err, rsp)
				} else {
					session.SetResponse(rsp)
				}
				// ready to response
				rspCh <- session
			}()
		}
	}
}

func (s *TcpServerTransport) Closed() <-chan struct{} {
	return s.closed
}

func (s *TcpServerTransport) Network() string {
	return s.net
}

func (s *TcpServerTransport) Address() string {
	return s.addr
}

func (s *TcpServerTransport) Codec() string {
	return s.codec.Name()
}
