package server

import (
	"context"
	"github.com/hitzhangjie/go-rpc/codec"
	"net"
	"sync"
	"time"
)

// UdpServer
type UdpServer struct {
	ctx    context.Context
	cancel context.CancelFunc

	net  string
	addr string

	codec  codec.Codec
	reader *UdpMessageReader

	//reqChan chan codec.Session
	rspChan chan codec.Session

	once   sync.Once
	closed chan struct{}

	opts *Options
}

func NewUdpServer(net, addr string, codecName string, opts ...Option) (ServerModule, error) {
	c := codec.ServerCodec(codecName)
	s := &UdpServer{
		net:    net,
		addr:   addr,
		codec:  c,
		reader: NewUdpMessageReader(c),
		once:   sync.Once{},
		closed: make(chan struct{}, 1),
		opts:   &Options{},
	}
	for _, o := range opts {
		o(s.opts)
	}
	return s, nil
}

func (s *UdpServer) Start() error {

	var (
		udpconn *net.UDPConn
		err     error
	)

	addr, err := net.ResolveUDPAddr(s.net, s.addr)
	if err != nil {
		return err
	}

	for {
		udpconn, err = net.ListenUDP(s.net, addr)
		if err != nil {
			if e, ok := err.(net.Error); ok && e.Temporary() {
				time.Sleep(time.Millisecond * 10)
				continue
			}
			return err
		}
		break
	}

	ep := UdpEndPoint{
		udpconn,
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

	return nil
}

func (s *UdpServer) Stop() {
	s.once.Do(func() {
		close(s.closed)
	})
}

func (s *UdpServer) Register(svr *Server) {
	s.ctx, s.cancel = context.WithCancel(svr.ctx)
	s.opts.router = svr.router
	svr.mods = append(svr.mods, s)
}

func (s *UdpServer) Closed() <-chan struct{} {
	return s.closed
}

// fixme this method `proc` appears in TcpServer, too. That's unnessary, refactor this
func (s *UdpServer) proc(reqCh <-chan interface{}, rspCh chan<- interface{}) {

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
