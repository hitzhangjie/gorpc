package errs

import (
	"github.com/hitzhangjie/gorpc/internal/errors"
)

var (
	// server errors
	ErrServerCtxDone     = newFrameworkError(1000, "server ctx done")
	ErrServerNotInit     = newFrameworkError(1001, "server not initialized")
	ErrSessionNotExisted = newFrameworkError(1002, "session not found")
	ErrRouteNotFound     = newFrameworkError(1003, "route not found")

	// codec error
	CodecDecodeError = newFrameworkError(2000, "decode error")
	CodecEncodeError = newFrameworkError(2001, "encode error")

	// message reader error
	CodecReadError      = newFrameworkError(3000, "read error")
	CodecReadIncomplete = newFrameworkError(3001, "read incomplete package")
	CodecReadInvalid    = newFrameworkError(3002, "read invalid package")
	CodecReadTooBig     = newFrameworkError(3003, "read too big package")

	// connection pool error
	ErrExceedPoolLimit     = newFrameworkError(4000, "connection poolFactory limit")  // ErrExceedPoolLimit 连接数量超过限制错误
	ErrPoolClosed          = newFrameworkError(4001, "connection poolFactory closed") // ErrPoolClosed 连接池关闭错误
	ErrConnClosed          = newFrameworkError(4002, "conn closed")                   // ErrConnClosed 连接关闭
	ErrConnPoolExceedLimit = newFrameworkError(4003, "connpool exceed limit")         // ErrConnPoolExceedLimit 超出连接数限制
)

func newFrameworkError(code int, msg string) *errors.Error {
	return errors.New(code, msg, errors.Framework)
}
