package server

import (
	"context"
	"github.com/hitzhangjie/go-rpc/codec"
	"net"
	"sync"
	"time"
)

// TcpServer
type TcpServer struct {
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

	opts *Options
}

const (
	tcpServerRspChanMaxLength = 1024
)

func NewTcpServer(net, addr string, codecName string, opts ...Option) (ServerModule, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	c := codec.ServerCodec(codecName)

	s := &TcpServer{
		ctx:    ctx,
		cancel: cancel,
		net:    net,
		addr:   addr,
		codec:  c,
		reader: NewTcpMessageReader(c),
		//rspChan: make(chan codec.Session, tcpServerRspChanMaxLength),
		once:   sync.Once{},
		closed: make(chan struct{}, 1),
		opts: &Options{
			router: nil,
		},
	}
	for _, o := range opts {
		o(s.opts)
	}
	return s, nil
}

func (s *TcpServer) Start() error {

	l, err := net.Listen(s.net, s.addr)
	if err != nil {
		return err
	}
	return s.serve(l)
}

func (s *TcpServer) Stop() {

	s.cancel()

	s.once.Do(func() {
		close(s.closed)
	})
}

func (s *TcpServer) Register(svr *Server) {
	s.ctx, s.cancel = context.WithCancel(svr.ctx)
	s.opts.router = svr.router
	svr.mods = append(svr.mods, s)
}

func (s *TcpServer) serve(l net.Listener) error {

	defer func() {
		s.cancel()
		l.Close()
	}()

	for {
		// check whether server closed
		select {
		case <-s.ctx.Done():
			return errServerCtxDone
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
			bufferPool.Get().([]byte),
		}
		ep.ctx, ep.cancel = context.WithCancel(s.ctx)

		go s.proc(ep.reqCh, ep.rspCh)

		go ep.Read()
		go ep.Write()
	}
}

func (s *TcpServer) proc(reqCh <-chan interface{}, rspCh chan<- interface{}) {

	builder := codec.GetSessionBuilder(s.reader.Codec.Name())

	for {
		select {
		case <-s.ctx.Done():
			s.cancel()
			return
		case req := <-reqCh:
			// build session
			session, err := builder.Build(req)
			if err != nil {
				// fixme error logging & metrics
				continue
			}
			// fixme using workerpool instead of goroutine
			r := s.opts.router

			go func() {
				// find route
				handle, err := r.Route(session.RPCName())
				if err != nil {
					session.SetErrorResponse(err)
					return
				}
				// pass session+req to handlefunc
				ctx := codec.ContextWithSession(s.ctx, session)
				rsp, err := handle(ctx, req)
				if err != nil {
					session.SetErrorResponse(err)
				} else {
					session.SetResponse(rsp)
				}
				// ready to response
				rspCh <- session
			}()
		}
	}
}

func (s *TcpServer) Closed() <-chan struct{} {
	return s.closed
}
