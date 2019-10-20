package selector

import (
	"github.com/hitzhangjie/go-rpc/client/selector/balancer"
)

type Options struct {
	balancer balancer.Balancer
}

type Option func(*Options)
