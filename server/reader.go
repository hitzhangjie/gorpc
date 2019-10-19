package server

import (
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
		mempool.Put(ep.buf)
		close(ep.reqCh)
	}()

	var (
		buflen int
		readsz int
		err    error
	)

	for {
		// check if server to be closed
		select {
		case <-ep.ctx.Done():
			return errServerCtxDone
		default:
		}

		// fixme conn read deadline
		ep.Conn.SetReadDeadline(time.Now().Add(time.Second*30))
		if readsz, err = ep.Conn.Read(ep.buf[buflen:]); err != nil {
			// fixme check tcpconn idle & release
			if e, ok := err.(net.Error); ok && e.Temporary() {
				time.Sleep(time.Millisecond * 10)
				continue
			}
			return err
		}
		buflen += readsz

		// decode请求
		req, sz, err := r.Codec.Decode(ep.buf[0:buflen])
		if err != nil {
			if err == codec.CodecReadIncomplete {
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


// UdpMessageReader read req from `net.Conn`, if read successfully, return the req'svr session.
//
// if any error occurs, it returns nil session and error, error should be one of the following:
// - io.Timeout
// - ...
type UdpMessageReader struct {
	Codec codec.Codec
}

func NewUdpMessageReader(codec codec.Codec) *UdpMessageReader {
	r := &UdpMessageReader{Codec: codec}
	return r
}

func (r *UdpMessageReader) Read(ep *UdpEndPoint) error {

	defer func() {
		ep.Conn.Close()
		mempool.Put(ep.buf)
		close(ep.reqCh)
	}()

	var (
		readsz int
		err    error
	)

	for {
		// check if server to be closed
		select {
		case <-ep.ctx.Done():
			return errServerCtxDone
		default:
		}

		// fixme conn read deadline
		ep.Conn.SetReadDeadline(time.Now().Add(time.Second*30))
		if readsz, err = ep.Conn.Read(ep.buf); err != nil {
			// fixme check Udpconn idle & release
			if e, ok := err.(net.Error); ok && e.Temporary() {
				time.Sleep(time.Millisecond * 10)
				continue
			}
			return err
		}
		// decode请求
		req, _, err := r.Codec.Decode(ep.buf[0:readsz])
		if err != nil {
			return err
		}

		ep.reqCh <- req
	}
}