package server

import (
	"context"
	"fmt"
	"github.com/hitzhangjie/go-rpc/codec"
	"net"
	"sync"
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

func (s *TcpServer) Start() {
	l, err := net.Listen(s.net, s.addr)
	if err != nil {
		panic(err)
	}
	go s.serve(l)
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

func (s *TcpServer) serve(l net.Listener) {

	defer func() {
		l.Close()
	}()

	for {
		// check whether server closed
		select {
		case <-s.ctx.Done():
			return
		default:
		}
		// accept request

		conn, err := l.Accept()
		if err != nil {
			// fixme handle error
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
		router := s.opts.router

		go func() {
			service, handle, err := router.Route(session)
			if err != nil {
				session.SetErrorResponse(err)
				return
			}
			err = handle(service, s.ctx, session)
			if err != nil {
				session.SetErrorResponse(err)
			}
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
