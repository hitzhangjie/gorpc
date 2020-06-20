package transport

import (
	"context"
	"fmt"
	"net"

	"github.com/hitzhangjie/gorpc/codec"
)

// UdpEndPoint udp endpoint
type UdpEndPoint struct {
	conn  net.Conn
	reqCh chan interface{}
	rspCh chan interface{}

	reader *UdpMessageReader
	ctx    context.Context
	cancel context.CancelFunc

	buf []byte
}

func (ep *UdpEndPoint) Read() {
	defer func() {
		ep.conn.Close()
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
		ep.conn.Close()
	}()
	for {
		// check whether server Closed
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
			data, err := ep.reader.codec.Encode(rsp)
			if err != nil {
				// fixme handle error
			}
			// fixme set write deadline
			ep.conn.Write(data)
		}
	}
}
