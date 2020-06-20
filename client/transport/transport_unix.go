package transport

import (
	"context"
	"github.com/hitzhangjie/gorpc-framework/client/pool"
)

type UnixTransport struct {
	ConnPool pool.ConnPool
}

func (t *UnixTransport) Send(ctx context.Context, network, address string, reqHead interface{}) (rsp interface{}, err error) {
	return nil, nil
}
