package transport

import (
	"context"
	"github.com/hitzhangjie/go-rpc/client/transport/connpool"
)

type UdpTransport struct {
	ConnPool connpool.ConnectionPool
}

func (t *UdpTransport) Send(ctx context.Context, req interface{}) (rsp interface{}, err error) {
	return nil, nil
}