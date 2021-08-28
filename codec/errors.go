package codec

import "errors"

var (
	ErrMarshalInvalidPB = errors.New("pkg isn't proto.Message")
)
