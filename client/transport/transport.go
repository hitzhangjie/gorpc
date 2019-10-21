package transport

import (
	"context"
	"sync"
)

type Transport interface {
	// fixme add some ...options
	Send(ctx context.Context, network, addr string, reqHead interface{}) (rsp interface{}, err error)
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 64*1024)
	},
}
