// Package codec contains logic for encode/decode, marshal/unmarshal, compress/decompress.
package codec

import (
	"sync"
)

// Codec encode & decode
//
// Both clientside and serverside should implements their own Codec.
type Codec interface {
	// Name codec name
	Name() string

	// Encode encode 'pkg' into []byte, here 'pkg' must be []byte,
	// which is marshaled and compressed data.
	Encode(pkg interface{}) (dat []byte, err error)

	// Decode decode []byte, return decoded body data and length,
	// which is marshaled and compressed data.
	Decode(dat []byte) (req interface{}, n int, err error)
}

var (
	codecsMux = sync.RWMutex{}
	codecs    = map[string]codec{}
)

type codec struct {
	name        string
	serverCodec Codec
	clientCodec Codec
}

// RegisterCodec registers codec of protocol
func RegisterCodec(protocol string, server, client Codec) {
	codecsMux.Lock()
	defer codecsMux.Unlock()

	codecs[protocol] = codec{
		name:        protocol,
		serverCodec: server,
		clientCodec: client,
	}
}

// ServerCodec returns server side codec of protocol
func ServerCodec(protocol string) Codec {
	codecsMux.RLock()
	defer codecsMux.RUnlock()

	return codecs[protocol].serverCodec
}

// ClientCodec returns client side codec of protocol
func ClientCodec(name string) Codec {
	codecsMux.RLock()
	defer codecsMux.RUnlock()
	return codecs[name].clientCodec
}
