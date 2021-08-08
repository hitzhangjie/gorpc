package errors

import (
	"github.com/hitzhangjie/gorpc/internal/errors"
)

var (
	// server errors
	ErrServerCtxDone     = newFrameworkError(1000, "server ctx done")        // 服务上下文终止
	ErrServerNotInit     = newFrameworkError(1001, "server not initialized") // 服务未初始化
	ErrSessionNotExisted = newFrameworkError(1002, "session not found")      // 会话未找到
	ErrRouteNotFound     = newFrameworkError(1003, "route not found")        // 路由未找到

	// codec error
	ErrCodecDecode = newFrameworkError(2000, "decode error") // 解包失败
	ErrCodecEncode = newFrameworkError(2001, "encode error") // 组包失败

	// message reader error
	ErrCodecRead           = newFrameworkError(3000, "read error")              // 收包错误
	ErrCodecReadIncomplete = newFrameworkError(3001, "read incomplete package") // 收包不完整
	ErrCodecReadInvalid    = newFrameworkError(3002, "read invalid package")    // 收包非法
	ErrCodecReadTooBig     = newFrameworkError(3003, "read too big package")    // 收包过大

	// connection pool error
	ErrExceedPoolLimit     = newFrameworkError(4000, "connection poolFactory limit")  // 连接数量超过限制错误
	ErrPoolClosed          = newFrameworkError(4001, "connection poolFactory closed") // 连接池关闭错误
	ErrConnClosed          = newFrameworkError(4002, "conn closed")                   // 连接关闭
	ErrConnPoolExceedLimit = newFrameworkError(4003, "connpool exceed limit")         // 超出连接数限制
)

func newFrameworkError(code int, msg string) *errors.Error {
	return errors.New(code, msg, errors.Framework)
}
