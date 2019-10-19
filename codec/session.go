package codec

import (
	"context"
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

type SessionBuilder interface {
	Build(reqHead interface{}) (Session, error)
}

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

var sessionKey = "session"

func SessionKey() string {
	return sessionKey
}

func SessionFromContext(ctx context.Context) Session {
	v := ctx.Value(sessionKey)
	session ,ok := v.(Session)
	if !ok {
		return nil
	}
	return session
}

func ContextWithSession(ctx context.Context, session Session) context.Context {
	return context.WithValue(ctx, sessionKey, session)
}
