package transport

import (
	"context"
)

type Transport interface {
	// fixme add some ...options
	Send(ctx context.Context, network, addr string, req []byte) (rsp []byte, err error)
}
