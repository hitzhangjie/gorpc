package transport

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/hitzhangjie/go-rpc/codec"
)

// TcpEndPoint endpoint of tcp connection
type TcpEndPoint struct {
	Conn  net.Conn
	ReqCh chan interface{}
	rspCh chan interface{}

	reader *TcpMessageReader
	Ctx    context.Context
	cancel context.CancelFunc

	Buf []byte
}

func (ep *TcpEndPoint) Read() {

	defer func() {
		ep.Conn.Close()
		ep.cancel()
	}()

	// keep reading message until we encounter some non-temporary errors
	err := ep.reader.Read(ep)
	if err != nil {
		// fixme handle error
		if err == io.EOF {
			log.Printf("peer connection Closed now, local:%s->remote:%s", ep.Conn.LocalAddr().String(), ep.Conn.RemoteAddr().String())
			return
		}
		log.Fatalf("tcp read request error:%v", err)
	}
}

func (ep *TcpEndPoint) Write() {

	defer func() {
		ep.Conn.Close()
	}()

	for {

		select {
		// check whether server Closed
		case <-ep.Ctx.Done():
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
			ep.Conn.SetWriteDeadline(time.Now().Add(time.Millisecond * 2000))

			n, err := ep.Conn.Write(data)
			if err != nil || len(data) != n {
				// fixme handle error
				log.Fatalf("tcp send response error:%v, bytes written got:%d, want:%d", err, n, len(data))
				continue
			}
		}
	}
}
