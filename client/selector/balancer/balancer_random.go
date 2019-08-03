package balancer

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// RandomBalancer randomly select next node
type RandomBalancer struct {
	Addrs []string
}

func (r *RandomBalancer) Next() (addr string) {
	idx := rand.Int() % len(r.Addrs)
	return r.Addrs[idx]
}
