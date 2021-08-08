package transport

import (
	"net"
	"sync"
	"time"

	"github.com/hitzhangjie/gorpc/codec"
	"github.com/hitzhangjie/gorpc/errors"
)

var udpBufferPool = &sync.Pool{
	New: func() interface{} {
		// make sure `len` of allocated buffer is not zero,
		// otherwise conn.Read(...) returns immediately.
		return make([]byte, 64*1024)
	},
}

// UdpMessageReader read req from `net.conn`, if read successfully, return the req'svr session.
//
// if any error occurs, it returns nil session and error, error should be one of the following:
// - io.Timeout
// - ...
type UdpMessageReader struct {
	codec codec.Codec
}

func NewUdpMessageReader(codec codec.Codec) *UdpMessageReader {
	r := &UdpMessageReader{codec: codec}
	return r
}

func (r *UdpMessageReader) Read(ep *UdpEndPoint) error {

	defer func() {
		ep.conn.Close()
		udpBufferPool.Put(ep.buf)
		close(ep.reqCh)
	}()

	var (
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
		if readsz, err = ep.conn.Read(ep.buf); err != nil {
			// fixme check Udpconn idle & release
			if e, ok := err.(net.Error); ok && e.Temporary() {
				time.Sleep(time.Millisecond * 10)
				continue
			}
			return err
		}
		// decode请求
		req, _, err := r.codec.Decode(ep.buf[0:readsz])
		if err != nil {
			return err
		}

		ep.reqCh <- req
	}
}
