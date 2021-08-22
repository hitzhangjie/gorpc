package codec

import (
	"sync"
)

// Codec encode & decode
type Codec interface {
	// Name codec name
	Name() string

	// Encode encode pkg into []byte
	Encode(pkg interface{}) (dat []byte, err error)

	// Decode decode []byte, return decoded interface{} and number of bytes
	Decode(dat []byte) (req interface{}, n int, err error)
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

// RegisterCodec registers codec of protocol
func RegisterCodec(protocol string, server, client Codec) {
	mux.Lock()
	defer mux.Unlock()
	codecs[protocol] = codec{
		name:   protocol,
		server: server,
		client: client,
	}
}

// ServerCodec returns server side codec of protocol
func ServerCodec(protocol string) Codec {
	mux.RLock()
	defer mux.RUnlock()
	return codecs[protocol].server
}

// ClientCodec returns client side codec of protocol
func ClientCodec(name string) Codec {
	mux.RLock()
	defer mux.RUnlock()
	return codecs[name].client
}
