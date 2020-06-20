package selector

import (
	"github.com/hitzhangjie/gorpc-framework/client/selector/balancer"
)

type Options struct {
	balancer balancer.Balancer
}

type Option func(*Options)
