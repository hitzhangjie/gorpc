package codec

import (
	"errors"
)

var (
	// Codec error
	CodecDecodeError = errors.New("decode error")
	CodecEncodeError = errors.New("encode error")

	// MsgReader error
	CodecReadError      = errors.New("read error")
	CodecReadIncomplete = errors.New("read incomplete package")
	CodecReadInvalid    = errors.New("read invalid package")
	CodecReadTooBig     = errors.New("read too big package")
)
