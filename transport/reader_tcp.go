package transport

import (
	"net"
	"sync"
	"time"

	"github.com/hitzhangjie/gorpc/codec"
	"github.com/hitzhangjie/gorpc/errors"
)

var tcpBufferPool = &sync.Pool{
	New: func() interface{} {
		// make sure `len` of allocated buffer is not zero,
		// otherwise conn.Read(...) returns immediately.
		return make([]byte, 64*1024)
	},
}

// TcpMessageReader read req from `net.conn`, if read successfully, return the req'svr session.
//
// if any error occurs, it returns nil session and error, error should be one of the following:
// - io.Timeout
// - ...
type TcpMessageReader struct {
	codec codec.Codec
}

func NewTcpMessageReader(codec codec.Codec) *TcpMessageReader {
	r := &TcpMessageReader{codec: codec}
	return r
}

func (r *TcpMessageReader) Read(ep *TcpEndPoint) error {

	defer func() {
		ep.conn.Close()
		tcpBufferPool.Put(ep.buf)
		close(ep.reqCh)
	}()

	var (
		buflen int
		readsz int
		err    error
	)

	for {
		// check if server to be Closed
		select {
		case <-ep.ctx.Done():
			return errors.ErrServerCtxDone
		default:
		}

		// fixme conn read deadline
		ep.conn.SetReadDeadline(time.Now().Add(time.Second * 30))
		if readsz, err = ep.conn.Read(ep.buf[buflen:]); err != nil {
			// fixme check tcpconn idle & release
			if e, ok := err.(net.Error); ok && e.Temporary() {
				time.Sleep(time.Millisecond * 10)
				continue
			}
			return err
		}
		buflen += readsz

		// decode请求
		req, sz, err := r.codec.Decode(ep.buf[0:buflen])
		if err != nil {
			if err == errors.ErrCodecReadIncomplete {
				continue
			}
			//return nil, err
			return err
		}

		ep.reqCh <- req
		ep.buf = ep.buf[sz:]
		buflen -= sz
	}
}
