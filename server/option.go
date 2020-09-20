package server

import (
	"github.com/hitzhangjie/gorpc/router"
)

type options struct {
	Router *router.Router
}

// Option specify server option
type Option func(*options)

// WithRouter specify router
func WithRouter(r *router.Router) Option {
	return func(opts *options) {
		opts.Router = r
	}
}
