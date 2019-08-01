package codec

import (
	"sync"
)

var (
	Mux            = sync.Mutex{}
	CodecMappings  = map[string]Codec{}
	ReaderMappings = map[string]*MessageReader{}
)
