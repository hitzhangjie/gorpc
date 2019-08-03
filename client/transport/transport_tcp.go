package transport

import (
	"context"
	"github.com/hitzhangjie/go-rpc/client/transport/connpool"
)

type TcpTransport struct {
	ConnPool connpool.ConnectionPool
}

func (t *TcpTransport) Send(ctx context.Context, req interface{}) (rsp interface{}, err error) {
	return nil, nil
}
