package server

import (
	"context"
	"fmt"
	"github.com/hitzhangjie/go-rpc/codec"
	"github.com/hitzhangjie/go-rpc/router"
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
	reader *codec.MessageReader

	//reqChan chan codec.Session
	rspChan chan codec.Session

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
		ctx:     ctx,
		cancel:  cancel,
		net:     net,
		addr:    addr,
		codec:   c,
		reader:  codec.NewMessageReader(c),
		rspChan: make(chan codec.Session, tcpServerRspChanMaxLength),
		once:    sync.Once{},
		closed:  make(chan struct{}, 1),
		opts:    &Options{},
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
				time.Sleep(time.Millisecond*10)
				continue
			}
			return nil
		}

		go s.read(conn)
		go s.write(conn)
	}
}

func (s *TcpServer) read(conn net.Conn) {
	defer func() {
		conn.Close()
	}()

	for {
		// check whether server closed
		select {
		case <-s.ctx.Done():
			return
		default:
		}
		// fixme set read deadline
		// read message
		req, err := s.reader.Read(conn)
		if err != nil {
			// fixme handle error
			fmt.Println("read error:", err)
			return
		}

		// fixme build session
		builder := codec.GetSessionBuilder(s.reader.Codec.Name())
		session, err := builder.Build(req)
		if err != nil {
			return
		}

		// fixme using workerpool instead of goroutine
		r := s.opts.router

		go func() {
			// find route
			handle, err := r.Route(session.RPC())
			if err != nil {
				session.SetErrorResponse(err)
				return
			}
			// pass session+req to handlefunc
			ctx := context.WithValue(s.ctx, router.SessionKey(), req)
			rsp, err := handle(ctx, req)
			if err != nil {
				session.SetErrorResponse(err)
			} else {
				session.SetResponse(rsp)
			}
			// ready to response
			s.rspChan <- session
		}()
	}

}

func (s *TcpServer) write(conn net.Conn) {
	defer func() {
		conn.Close()
	}()
	for {
		// check whether server closed
		select {
		case <-s.ctx.Done():
			return
		default:
		}
		// write response
		select {
		case session := <-s.rspChan:
			rsp := session.Response()
			data, err := s.codec.Encode(rsp)
			if err != nil {
				// fixme handle error
			}
			// fixme set write deadline
			conn.Write(data)
		}
	}
}

func (s *TcpServer) Closed() <-chan struct{} {
	return s.closed
}
