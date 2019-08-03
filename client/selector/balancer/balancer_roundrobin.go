package balancer

import (
	"sync/atomic"
)

var rrIdx int64

// RoundRobinBalancer select next node by roundrobin
type RoundRobinBalancer struct {
	Addrs []string
}

func (r *RoundRobinBalancer) Next() string {
	idx := atomic.AddInt64(&rrIdx, 1)
	idx = idx % int64(len(r.Addrs))
	return r.Addrs[idx]
}
