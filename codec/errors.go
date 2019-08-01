package codec

import (
	"errors"
)

var (
	// Codec error
	codecDecodeError = errors.New("decode error")
	codecEncodeError = errors.New("encode error")

	// MsgReader error
	codecReadError = errors.New("read error")
)
