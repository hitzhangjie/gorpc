package codec

// Codec encode & decode
//
// don't use terms marshal/umarshal, why? marshal/unmarshal is used for serializing/de-serializing data,
// - when we decode a request/response binary data, data includes pkg length, magic number, req/rsp payload, etc.
// - when we encode a req/rsp payload, we also add pkg length, magic number, req/rsp payload, etc.
// Needless to say, marshal/unmarshal is not appropriate!
type Codec interface {
	// Name codec name
	Name() string
	// Encode encode pkg into []byte
	Encode(pkg interface{}) ([]byte, error)
	// Decode decode []byte into interface{}
	Decode([]byte, interface{}) error
}
