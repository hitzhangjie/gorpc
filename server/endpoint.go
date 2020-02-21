package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/hitzhangjie/go-rpc/codec"
)

// EndPoint endpoint represents one side of net.Conn
//
// Read read data from net.Conn
// Write write data to net.Conn
type EndPoint interface {
	Read()
	Write()
}

// TcpEndPoint endpoint of tcp connection
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
		if err == io.EOF {
			log.Printf("peer connection closed now, local:%s->remote:%s", ep.Conn.LocalAddr().String(), ep.Conn.RemoteAddr().String())
			return
		}
		log.Fatalf("tcp read request error:%v", err)
	}
}

func (ep *TcpEndPoint) Write() {

	defer func() {
		ep.Close()
	}()

	for {

		select {
		// check whether server closed
		case <-ep.ctx.Done():
			return
		// write response
		case v := <-ep.rspCh:
			fmt.Println("handle response")
			session := v.(codec.Session)
			rsp := session.Response()
			data, err := ep.reader.Codec.Encode(rsp)
			if err != nil {
				log.Fatalf("tcp encode respone error:%v", err)
				continue
			}

			// fixme set write deadline, make the value configurable
			ep.SetWriteDeadline(time.Now().Add(time.Millisecond * 2000))

			n, err := ep.Conn.Write(data)
			if err != nil || len(data) != n {
				// fixme handle error
				log.Fatalf("tcp send response error:%v, bytes written got:%d, want:%d", err, n, len(data))
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
