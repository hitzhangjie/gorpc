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
	// fixme using sync.Pool instead of []byte
	buf := make([]byte, 1024, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	dat := buf[:n]

	// decode请求
	req, err := r.Codec.Decode(dat)
	if err != nil {
		return nil, err
	}

	return req, nil
}
