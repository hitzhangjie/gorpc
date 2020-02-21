package server

import (
	"github.com/hitzhangjie/go-rpc/router"
)

type Options struct {
	Router *router.Router
}

type Option func(*Options)

func WithRouter(r *router.Router) Option {
	return func(opts *Options) {
		opts.Router = r
	}
}
