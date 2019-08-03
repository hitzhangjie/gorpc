package codec

import (
	"sync"
)

type Session interface {
	RPC() string

	Request() interface{}
	SetRequest(req interface{})

	Response() interface{}
	SetResponse(rsp interface{})
	SetErrorResponse(error)

	TraceContext() interface{}
}

type BaseSession struct {
	Request  interface{}
	Response interface{}
}

var (
	lock     sync.RWMutex
	builders = map[string]SessionBuilder{}
)

type SessionBuilder func(reqHead []byte) (Session, error)

func RegisterSessionBuilder(name string, builder SessionBuilder) {
	lock.Lock()
	defer lock.Unlock()
	builders[name] = builder
}

func GetSessionBuilder(name string) SessionBuilder {
	lock.RLock()
	defer lock.RUnlock()
	return builders[name]
}
