package transport

import (
	"context"
	"github.com/hitzhangjie/go-rpc/client/transport/connpool"
)

type UnixTransport struct {
	ConnPool connpool.ConnectionPool
}

func (t *UnixTransport) Send(ctx context.Context, network, address string, req []byte) (rsp []byte, err error) {
	return nil, nil
}
