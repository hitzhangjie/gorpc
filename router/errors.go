package router

import "errors"

var (
	errRouteNotFound     = errors.New("route not found")
	errSessionNotExisted = errors.New("session not found")
)
