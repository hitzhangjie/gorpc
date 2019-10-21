package selector

import (
	"testing"
)

func TestIPSelector(t *testing.T) {
	var (
		network = "tcp"
		address = "127.0.0.1:8888"
	)
	s := NewIPSelector(network, []string{address})
	if s == nil {
		t.Fatalf("ipselector create failed")
	}

	n, err := s.Select("")
	if err != nil {
		t.Fatalf("ipselector select error:%v", err)
	}

	if n.Network != "tcp" && n.Address != "127.0.0.1:8888" {
		t.Fatalf("ipselector select error, got %s/%s, want %s/%s", n.Address, n.Network, address, network)
	}
}
