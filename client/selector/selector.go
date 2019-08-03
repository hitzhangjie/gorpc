package selector

type Selector interface {
	Select(service string) (*Node, error)
	Update(node *Node, err error) error
}

type Node struct {
	Net  string // tcp, tcp4, tcp6 or udp, udp4, udp6
	Addr string // ip:port
	RPC  string // /$pkgname.$service/$method
}
