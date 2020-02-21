package server

import (
	"github.com/hitzhangjie/go-rpc/router"
)

type options struct {
	router *router.Router
}

type Option func(*options)

func WithRouter(r *router.Router) Option {
	return func(opts *options) {
		opts.router = r
	}
}
