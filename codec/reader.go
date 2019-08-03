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

func (r *MessageReader) Read(conn net.Conn) (Session, error) {
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

	// 根据请求构建session，这里就意味着Codec要做更多事情，NewSession
	builder := GetSessionBuilder(r.Codec.Name())
	session, err := builder(req)
	if err != nil {
		return nil, err
	}
	return session, nil
}
