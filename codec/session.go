package codec

import (
	"context"
	"sync"
)

// Session defines the rpc context for a request
//
// Usually, session can be stored in client or server side, we adopt the term `session`
// instead of `RpcContext` for simplicity.
type Session interface {

	// RPC return rpc name, i.e., the method name defined in pb service.rpc.name
	RPC() string

	// Request return the request header
	Request() interface{}
	// SetRequest set the request header
	SetRequest(req interface{})

	// Response return the response header
	Response() interface{}
	// SetResponse set the response header
	SetResponse(rsp interface{})
	// SetErrorResponse set the error response
	SetErrorResponse(error)

	// TraceContext return the tracing context
	TraceContext() interface{}
}

// BaseSession implements some basic methods defined in `Session`
type BaseSession struct {
	Request  interface{}
	Response interface{}
}

var (
	lock     sync.RWMutex
	builders = map[string]SessionBuilder{}
)

// SessionBuilder when extending protocols, SessionBuilder should be
// implemented and registered to help build the `Session`.
type SessionBuilder interface {
	Build(reqHead interface{}) (Session, error)
}

// RegisterSessionBuilder register extended SessionBuilder for protocol `proto`
func RegisterSessionBuilder(proto string, builder SessionBuilder) {
	lock.Lock()
	defer lock.Unlock()
	builders[proto] = builder
}

// GetSessionBuilder return SessionBuilder for protocol `proto`
func GetSessionBuilder(proto string) SessionBuilder {
	lock.RLock()
	defer lock.RUnlock()
	return builders[proto]
}

var sessionKey = "session"

// SessionFromContext return Session carried by `ctx`
func SessionFromContext(ctx context.Context) Session {
	v := ctx.Value(sessionKey)
	session ,ok := v.(Session)
	if !ok {
		return nil
	}
	return session
}

// ContextWithSession return new context carrying value `session`
func ContextWithSession(ctx context.Context, session Session) context.Context {
	return context.WithValue(ctx, sessionKey, session)
}
