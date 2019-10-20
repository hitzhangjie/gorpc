package transport

import (
	"context"
	"github.com/hitzhangjie/go-rpc/client/pool"
)

type UnixTransport struct {
	ConnPool pool.ConnPool
}

func (t *UnixTransport) Send(ctx context.Context, network, address string, req []byte) (rsp []byte, err error) {
	return nil, nil
}
