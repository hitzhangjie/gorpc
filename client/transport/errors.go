package transport

import (
	"errors"
)

var (
	errConnPoolExceedLimit = errors.New("connpool exceed limit")
)
