package transport

import (
	"context"
	"github.com/hitzhangjie/go-rpc/client/transport/connpool"
)

type UdpTransport struct {
	ConnPool connpool.ConnectionPool
}

func (t *UdpTransport) Send(ctx context.Context, network, address string, req []byte) (rsp []byte, err error) {
	return nil, nil
}
