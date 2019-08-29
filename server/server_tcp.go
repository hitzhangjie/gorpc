package server

import (
	"fmt"
	"github.com/hitzhangjie/go-rpc/codec"
	"net"
)

// TcpServer
type TcpServer struct {
	svr    *Server
	net    string
	addr   string
	codec  codec.Codec
	reader *codec.MessageReader
	//reqChan chan codec.Session
	rspChan chan codec.Session
}

const (
	tcpServerRspChanMaxLength = 1024
)

func NewTcpServer(net, addr string, codecName string, opts ...Option) (ServerModule, error) {
	c := codec.ServerCodec(codecName)
	s := &TcpServer{
		net:     net,
		addr:    addr,
		codec:   c,
		reader:  codec.NewMessageReader(c),
		rspChan: make(chan codec.Session, tcpServerRspChanMaxLength),
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
}

func (s *TcpServer) Register(svr *Server) {
	s.svr = svr
	svr.mods = append(svr.mods, s)
}

func (s *TcpServer) serve(l net.Listener) {

	defer func() {
		l.Close()
	}()

	for {
		// check whether server closed
		select {
		case <-s.svr.ctx.Done():
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
		case <-s.svr.ctx.Done():
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
		go func() {
			service, handle, err := s.svr.router.Route(session)
			if err != nil {
				session.SetErrorResponse(err)
				return
			}
			err = handle(service, s.svr.ctx, session)
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
		case <-s.svr.ctx.Done():
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
