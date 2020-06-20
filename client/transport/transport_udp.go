package transport

import (
	"context"
	"github.com/hitzhangjie/gorpc-framework/client/pool"
	"github.com/hitzhangjie/gorpc-framework/codec"
)

type UdpTransport struct {
	ConnPool pool.ConnPool
	Codec    codec.Codec
}

func (t *UdpTransport) Send(ctx context.Context, network, address string, reqHead interface{}) (rspHead interface{}, err error) {
	return nil, nil
}
