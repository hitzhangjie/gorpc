package errs

import "errors"

var (
	// server errors
	ErrServerCtxDone     = errors.New("server Ctx done")
	ErrServerNotInit     = errors.New("server not initialized")
	ErrSessionNotExisted = errors.New("session not found")
	ErrRouteNotFound     = errors.New("route not found")

	// Codec error
	CodecDecodeError = errors.New("decode error")
	CodecEncodeError = errors.New("encode error")

	// MsgReader error
	CodecReadError      = errors.New("read error")
	CodecReadIncomplete = errors.New("read incomplete package")
	CodecReadInvalid    = errors.New("read invalid package")
	CodecReadTooBig     = errors.New("read too big package")

	// conn pool error
	ErrExceedPoolLimit     = errors.New("connection poolFactory limit")  // ErrExceedPoolLimit 连接数量超过限制错误
	ErrPoolClosed          = errors.New("connection poolFactory closed") // ErrPoolClosed 连接池关闭错误
	ErrConnClosed          = errors.New("conn closed")                   // ErrConnClosed 连接关闭
	ErrConnPoolExceedLimit = errors.New("connpool exceed limit")         // ErrConnPoolExceedLimit 超出连接数限制
)
