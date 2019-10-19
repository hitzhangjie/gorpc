package server

import "errors"

var (
	errServerNotInit = errors.New("server not initialized")
	errServerCtxDone = errors.New("server ctx done")
)
