package codec

import (
	"sync"

	"github.com/hitzhangjie/go-rpc/transport"
)

// MessageReader message reader for any codec
type MessageReader interface {
	Codec() Codec
	Read(ep transport.EndPoint) Session
}

var (
	readers    = make(map[string]MessageReader)
	readersLck = sync.RWMutex{}
)

// RegisterReader register message reader
func RegisterReader(reader MessageReader) {

	readersLck.Lock()
	defer readersLck.Unlock()

	codec := reader.Codec()
	readers[codec.Name()] = reader
	readersLck.Unlock()
}

// Reader return message reader for `codec`
func Reader(ccodec string) MessageReader {

	readersLck.RLock()
	defer readersLck.RUnlock()

	c, ok := readers[ccodec]
	if !ok {
		return nil
	}
	return c
}
