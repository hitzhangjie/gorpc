package server

import (
	"context"
	"fmt"
	"github.com/hitzhangjie/go-rpc/codec"
	"net"
	"time"
)

// EndPoint
type EndPoint interface {
	Read()
	Write()
}

// TcpEndPoint tcp endpoint
type TcpEndPoint struct {
	net.Conn
	reqCh chan interface{}
	rspCh chan interface{}

	reader *TcpMessageReader
	ctx    context.Context
	cancel context.CancelFunc

	buf []byte
}

func (ep *TcpEndPoint) Read() {

	defer func() {
		ep.Close()
		ep.cancel()
	}()

	// keep reading message until we encounter some non-temporary errors
	err := ep.reader.Read(ep)
	if err != nil {
		// fixme handle error
		fmt.Println("read error:", err)
	}
}

func (ep *TcpEndPoint) Write() {

	defer func() {
		ep.Close()
	}()

	for {
		// check whether server closed
		select {
		case <-ep.ctx.Done():
			return
		default:
		}

		// write response
		select {
		case v := <-ep.rspCh:
			session := v.(codec.Session)
			rsp := session.Response()
			data, err := ep.reader.Codec.Encode(rsp)
			if err != nil {
				// fixme handle error
				continue
			}

			// fixme set write deadline, make the value configurable
			ep.SetWriteDeadline(time.Now().Add(time.Millisecond * 2000))

			n, err := ep.Conn.Write(data)
			if err != nil || len(data) != n {
				// fixme handle error
				continue
			}
		}
	}
}

// UdpEndPoint udp endpoint
type UdpEndPoint struct {
	net.Conn
	reqCh chan interface{}
	rspCh chan interface{}

	reader *UdpMessageReader
	ctx    context.Context
	cancel context.CancelFunc

	buf []byte
}

func (ep *UdpEndPoint) Read() {
	defer func() {
		ep.Close()
	}()

	// keep reading message, until when we encounter any non-temporary errors
	err := ep.reader.Read(ep)
	if err != nil {
		// fixme handle error
		fmt.Println("read error:", err)
	}
}

func (ep *UdpEndPoint) Write() {
	defer func() {
		ep.Close()
	}()
	for {
		// check whether server closed
		select {
		case <-ep.ctx.Done():
			ep.cancel()
			return
		default:
		}
		// write response
		select {
		case v := <-ep.rspCh:
			session := v.(codec.Session)
			rsp := session.Response()
			data, err := ep.reader.Codec.Encode(rsp)
			if err != nil {
				// fixme handle error
			}
			// fixme set write deadline
			ep.Conn.Write(data)
		}
	}
}
