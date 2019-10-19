package codec

import (
	"net"
)

// MessageReader read req from `net.Conn`, if read successfully, return the req's session.
//
// if any error occurs, it returns nil session and error, error should be one of the following:
// - io.Timeout
// - ...
type MessageReader struct {
	Codec Codec
}

func NewMessageReader(codec Codec) *MessageReader {
	r := &MessageReader{Codec: codec}
	return r
}

func (r *MessageReader) Read(conn net.Conn) (interface{}, error) {

	var (
		req interface{}
		err error
		n   int
		m   int
	)

	for {
		if m, err = conn.Read(buf[n:]); err != nil {
			return nil, err
		}
		n += m

		dat := buf[:n]

		// decode请求
		req, err = r.Codec.Decode(dat)
		if err != nil {
			if err == CodecReadIncomplete {
				continue
			}
			return nil, err
		}
	}

	return req, nil
}
