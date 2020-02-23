package transport

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/hitzhangjie/go-rpc/codec"
)

// UdpServerTransport
type UdpServerTransport struct {
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

	opts *options
}

func NewUdpServerTransport(ctx context.Context, net, addr string, codecName string, opts ...Option) (Transport, error) {
	c := codec.ServerCodec(codecName)
	s := &UdpServerTransport{
		net:    net,
		addr:   addr,
		codec:  c,
		reader: NewUdpMessageReader(c),
		once:   sync.Once{},
		closed: make(chan struct{}, 1),
		opts:   &options{},
	}
	s.ctx, s.cancel = context.WithCancel(ctx)
	for _, o := range opts {
		o(s.opts)
	}
	return s, nil
}

func (s *UdpServerTransport) ListenAndServe() error {

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
		udpBufferPool.Get().([]byte),
	}
	ep.ctx, ep.cancel = context.WithCancel(s.ctx)

	go s.proc(ep.reqCh, ep.rspCh)
	go ep.Read()
	go ep.Write()

	s.once.Do(func() {
		close(s.closed)
	})
	return nil
}

//func (s *UdpServerTransport) Register(svr *server.Service) {
//	s.ctx, s.cancel = context.WithCancel(svr.ctx)
//	s.opts.router = svr.router
//	svr.trans = append(svr.trans, s)
//}

func (s *UdpServerTransport) Closed() <-chan struct{} {
	return s.closed
}

// fixme this method `proc` appears in TcpServerTransport, too. That's unnessary, refactor this
func (s *UdpServerTransport) proc(reqCh <-chan interface{}, rspCh chan<- interface{}) {

	builder := codec.GetSessionBuilder(s.reader.codec.Name())

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
			r := s.opts.Router

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

func (s *UdpServerTransport) Network() string {
	return s.net
}

func (s *UdpServerTransport) Address() string {
	return s.addr
}

func (s *UdpServerTransport) Codec() string {
	return s.codec.Name()
}
