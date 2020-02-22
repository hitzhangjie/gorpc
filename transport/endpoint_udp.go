package transport

import (
	"context"
	"fmt"
	"net"

	"github.com/hitzhangjie/go-rpc/codec"
)

// UdpEndPoint udp endpoint
type UdpEndPoint struct {
	Conn  net.Conn
	ReqCh chan interface{}
	rspCh chan interface{}

	reader *UdpMessageReader
	Ctx    context.Context
	cancel context.CancelFunc

	Buf []byte
}

func (ep *UdpEndPoint) Read() {
	defer func() {
		ep.Conn.Close()
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
		ep.Conn.Close()
	}()
	for {
		// check whether server Closed
		select {
		case <-ep.Ctx.Done():
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
