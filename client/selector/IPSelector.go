package selector

import (
	"github.com/hitzhangjie/gorpc/client/selector/balancer"
)

// IPSelector selector based on []IP
type IPSelector struct {
	network  string
	addrs    []string
	balancer balancer.Balancer
}

// NewIPSelector create a new IPSelector
func NewIPSelector(network string, addrs []string) *IPSelector {
	return &IPSelector{
		network:  network,
		addrs:    addrs,
		balancer: &balancer.RandomBalancer{addrs},
	}
}

// Select return next node
func (s *IPSelector) Select(service string) (*Node, error) {
	addr := s.balancer.Next()
	node := Node{
		Network: s.network,
		Address: addr,
	}
	return &node, nil
}

// Update do nothing in IPSelector
func (s *IPSelector) Update(node *Node, err error) error {
	return nil
}
