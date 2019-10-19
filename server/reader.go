package server

import (
	"context"
	"github.com/hitzhangjie/go-rpc/codec"
	"net"
	"sync"
	"time"
)

var mempool = &sync.Pool{
	New: func() interface{} {
		return make([]byte, 0, 64*1024)
	},
}

// MessageReader read req from `net.Conn`, if read successfully, return the req'svr session.
//
// if any error occurs, it returns nil session and error, error should be one of the following:
// - io.Timeout
// - ...
type MessageReader struct {
	Codec codec.Codec
}

func NewMessageReader(codec codec.Codec) *MessageReader {
	r := &MessageReader{Codec: codec}
	return r
}

func (r *MessageReader) Read(ctx context.Context, conn net.Conn, reqCh chan interface{}) error {

	buf := mempool.Get().([]byte)

	defer func() {
		conn.Close()
		mempool.Put(buf)
		close(reqCh)
	}()

	var (
		buflen int
		readsz int
		err    error
	)

	for {

		// check if server to be closed
		select {
		case <-ctx.Done():
			return errServerCtxDone
		default:
		}

		// fixme conn read deadline
		if readsz, err = conn.Read(buf[buflen:]); err != nil {
			// fixme check tcpconn idle & release
			if e, ok := err.(net.Error); ok && e.Temporary() {
				time.Sleep(time.Millisecond * 10)
				continue
			}
			return err
		}
		buflen += readsz

		// decode请求
		req, sz, err := r.Codec.Decode(buf[0:buflen])
		if err != nil {
			if err == codec.CodecReadIncomplete {
				continue
			}
			//return nil, err
			return err
		}

		reqCh <- req
		buf = buf[sz:]
		buflen -= sz
	}
}
