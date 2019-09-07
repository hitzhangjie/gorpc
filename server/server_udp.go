package server

import (
	"context"
	"fmt"
	"github.com/hitzhangjie/go-rpc/codec"
	"github.com/hitzhangjie/go-rpc/router"
	"net"
	"sync"
)

// UdpServer
type UdpServer struct {
	ctx    context.Context
	cancel context.CancelFunc

	net  string
	addr string

	codec  codec.Codec
	reader *codec.MessageReader

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
		reader: codec.NewMessageReader(c),
		once:   sync.Once{},
		closed: make(chan struct{}, 1),
		opts:   &Options{},
	}
	for _, o := range opts {
		o(s.opts)
	}
	return s, nil
}

func (s *UdpServer) Start() {
	addr, err := net.ResolveUDPAddr(s.net, s.addr)
	if err != nil {
		panic(err)
	}
	udpconn, err := net.ListenUDP(s.net, addr)
	if err != nil {
		panic(err)
	}
	go s.read(udpconn)
	go s.write(udpconn)
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

func (s *UdpServer) read(conn net.Conn) {
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
			continue
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
			ctx := context.WithValue(s.ctx, router.SessionKey(), session)
			rsp, err := handle(ctx, session)
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

func (s *UdpServer) write(conn net.Conn) {
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

func (s *UdpServer) Closed() <-chan struct{} {
	return s.closed
}
