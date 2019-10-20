package transport

import (
	"context"
	"github.com/hitzhangjie/go-rpc/client/transport/connpool"
	"github.com/hitzhangjie/go-rpc/codec"
)

type UdpTransport struct {
	ConnPool pool.ConnectionPool
	Codec    codec.Codec
}

func (t *UdpTransport) Send(ctx context.Context, network, address string, req []byte) (rsp []byte, err error) {
	return nil, nil
}
