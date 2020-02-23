package transport

import (
	"net"
	"sync"
	"time"

	"github.com/hitzhangjie/go-rpc/codec"
	"github.com/hitzhangjie/go-rpc/errs"
)

var tcpBufferPool = &sync.Pool{
	New: func() interface{} {
		// make sure `len` of allocated buffer is not zero,
		// otherwise conn.Read(...) returns immediately.
		return make([]byte, 64*1024)
	},
}

// TcpMessageReader read req from `net.Conn`, if read successfully, return the req'svr session.
//
// if any error occurs, it returns nil session and error, error should be one of the following:
// - io.Timeout
// - ...
type TcpMessageReader struct {
	Codec codec.Codec
}

func NewTcpMessageReader(codec codec.Codec) *TcpMessageReader {
	r := &TcpMessageReader{Codec: codec}
	return r
}

func (r *TcpMessageReader) Read(ep *TcpEndPoint) error {

	defer func() {
		ep.Conn.Close()
		tcpBufferPool.Put(ep.Buf)
		close(ep.ReqCh)
	}()

	var (
		buflen int
		readsz int
		err    error
	)

	for {
		// check if server to be Closed
		select {
		case <-ep.Ctx.Done():
			return errs.ErrServerCtxDone
		default:
		}

		// fixme conn read deadline
		ep.Conn.SetReadDeadline(time.Now().Add(time.Second * 30))
		if readsz, err = ep.Conn.Read(ep.Buf[buflen:]); err != nil {
			// fixme check tcpconn idle & release
			if e, ok := err.(net.Error); ok && e.Temporary() {
				time.Sleep(time.Millisecond * 10)
				continue
			}
			return err
		}
		buflen += readsz

		// decode请求
		req, sz, err := r.Codec.Decode(ep.Buf[0:buflen])
		if err != nil {
			if err == errs.CodecReadIncomplete {
				continue
			}
			//return nil, err
			return err
		}

		ep.ReqCh <- req
		ep.Buf = ep.Buf[sz:]
		buflen -= sz
	}
}
