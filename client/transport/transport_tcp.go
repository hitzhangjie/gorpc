package transport

import (
	"context"
	"fmt"
	"github.com/hitzhangjie/go-rpc/client/pool"
	"github.com/hitzhangjie/go-rpc/codec"
	"time"
)

// TcpTransport tcp transport
type TcpTransport struct {
	Pool  pool.PoolFactory
	Codec codec.Codec
}

// Send send reqHead and return rspHead, return an error if encountered
func (t *TcpTransport) Send(ctx context.Context, network, address string, reqHead interface{}) (rspHead interface{}, err error) {

	// encode
	data, err := t.Codec.Encode(reqHead)
	if err != nil {
		return nil, err
	}

	// get conn
	conn, err := t.Pool.Get(ctx, network, address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(time.Millisecond * 200))

	// conn write
	n, err := conn.Write(data)
	if err != nil {
		return nil, err
	}

	if len(data) != n {
		return nil, fmt.Errorf("write error, write only %d bytes, want write %d bytes", n, data)
	}

	// alloc buffer
	buf := bufferPool.Get().([]byte)
	defer bufferPool.Put(buf)

	// conn read
	var sz int
	for {
		n, err := conn.Read(buf[sz:])
		if err != nil {
			return nil, err
		}
		sz += n

		// decode
		rsp, _, err := t.Codec.Decode(buf[:sz])
		if err != nil {
			if err == codec.CodecReadIncomplete {
				continue
			}
			return nil, err
		}

		// fixme for now, we only support one-req-one-response transport mode
		// so, here we can return now
		return rsp, nil
	}
}
