package transport

import (
	"context"
	"fmt"
	"github.com/hitzhangjie/go-rpc/client/transport/connpool"
	"github.com/hitzhangjie/go-rpc/codec"
	"net"
	"time"
)

type TcpTransport struct {
	ConnPool pool.ConnectionPool
	Codec    codec.Codec
}

func (t *TcpTransport) Send(ctx context.Context, network, address string, req []byte) (rsp []byte, err error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}

	conn.SetDeadline(time.Now().Add(time.Millisecond * 200))

	n, err := conn.Write(req)
	if err != nil {
		return nil, err
	}

	if len(req) != n {
		return nil, fmt.Errorf("write error, write only %d bytes, want write %d bytes", n, req)
	}

	buf := make([]byte, 64*1024)
	conn.Read(buf)

	// fixme who is reponsible for decode, absolutely we need passing decoder & encoder to transport

	return nil, nil
}
