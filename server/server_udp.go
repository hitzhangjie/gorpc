package server

import (
	"github.com/hitzhangjie/go-rpc/codec"
	"net"
)

// UdpServer
type UdpServer struct {
	svr    *Server
	net    string
	addr   string
	codec  codec.Codec
	reader *codec.MessageReader
	//reqChan chan codec.Session
	rspChan chan codec.Session
}

func NewUdpServer(net, addr string, codecName string, opts ...Option) (ServerModule, error) {
	s := &UdpServer{
		net:    net,
		addr:   addr,
		codec:  codec.CodecMappings[codecName],
		reader: codec.ReaderMappings[codecName],
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

func (s *UdpServer) Stop() {}

func (s *UdpServer) Register(svr *Server) {
	s.svr = svr
	svr.mods = append(svr.mods, s)
}

func (s *UdpServer) read(conn net.Conn) {
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
		session, err := s.reader.Read(conn)
		if err != nil {
			// fixme handle error
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

func (s *UdpServer) write(conn net.Conn) {
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
