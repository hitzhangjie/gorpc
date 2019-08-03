package codec

import (
	"sync"
)

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
	Encode(pkg interface{}) (dat []byte, err error)
	// Decode decode []byte into interface{}
	Decode(dat []byte) (req interface{}, err error)
}

var (
	mux    = sync.RWMutex{}
	codecs = map[string]codec{}
)

type codec struct {
	name   string
	server Codec
	client Codec
}

func RegisterCodec(name string, server, client Codec) {
	mux.Lock()
	defer mux.Unlock()
	codecs[name] = codec{
		name:   name,
		server: server,
		client: client,
	}
}

func ServerCodec(name string) Codec {
	mux.RLock()
	defer mux.RUnlock()
	return codecs[name].server
}

func ClientCodec(name string) Codec {
	mux.RLock()
	defer mux.RUnlock()
	return codecs[name].client
}
