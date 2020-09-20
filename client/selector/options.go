package selector

import (
	"github.com/hitzhangjie/gorpc/client/selector/balancer"
)

type Options struct {
	balancer balancer.Balancer
}

type Option func(*Options)
