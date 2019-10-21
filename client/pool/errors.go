package pool

import (
	"errors"
)

var (
	errExceedPoolLimit = errors.New("connection poolFactory limit")  // errExceedPoolLimit 连接数量超过限制错误
	errPoolClosed      = errors.New("connection poolFactory closed") // errPoolClosed 连接池关闭错误
	errConnClosed      = errors.New("conn closed")            // errConnClosed 连接关闭
)
