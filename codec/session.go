package codec

import (
	"context"
	"sync"
)

// Session defines the rpc context for a request.
//
// Usually, session can be stored in client and server side:
// - unary rpc:
// 	- when client starts an rpc, new session is initialized locally
// 	- when server receives the request, new session is initialized remotely
// - streaming rpc:
// 	- when client starts an stream, new session is initialized locally
// 	- when server receives the openning stream request, new session is initialized remotely
//
// we uses term `session` instead of `RpcContext` so as to avoid being
// misunderstood btw 'context.Context' and 'codec.RpcContext', though 'RpcContext'
// is a good name.
type Session interface {
	// RPCName returns rpc name, i.e., the method name defined in pb service.rpc.name
	RPCName() string

	// ReqHead returns the request header of unary rpc, or next request header of streaming rpc
	Request() interface{}

	// RspHead returns the response header
	Response() interface{}

	// SetResponse sets the response header
	SetResponse(rsp interface{})

	// SetError sets the error status of response
	SetError(error)
}

// BaseSession implements some basic methods defined in `Session`
type BaseSession struct {
	ReqHead interface{}
	RspHead interface{}
}

// Request returns request header
func (s *BaseSession) Request() interface{} {
	if s != nil {
		return s.ReqHead
	}
	return nil
}

// Response returns response header
func (s *BaseSession) Response() interface{} {
	if s != nil {
		return s.RspHead
	}
	return nil
}

// SetRequest sets response header
func (s *BaseSession) SetResponse(rsp interface{}) {
	if s != nil {
		s.RspHead = rsp
	}
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
	session, ok := v.(Session)
	if !ok {
		return nil
	}
	return session
}

// ContextWithSession return new context carrying value `session`
func ContextWithSession(ctx context.Context, session Session) context.Context {
	return context.WithValue(ctx, sessionKey, session)
}
