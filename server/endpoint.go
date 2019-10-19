package server

import (
	"context"
	"fmt"
	"github.com/hitzhangjie/go-rpc/codec"
	"net"
	"time"
)

type TcpEndPoint struct {
	net.Conn
	reqCh chan interface{}
	rspCh chan interface{}

	reader *MessageReader
	ctx    context.Context
	cancel context.CancelFunc
}

func (ep *TcpEndPoint) Read() {

	defer func() {
		ep.Close()
		ep.cancel()
	}()

	// keep reading message until we encounter some non-temporary errors
	err := ep.reader.Read(ep.ctx, ep.Conn, ep.reqCh)
	if err != nil {
		// fixme handle error
		fmt.Println("read error:", err)
		return
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
