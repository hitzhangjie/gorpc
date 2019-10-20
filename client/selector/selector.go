package selector

type Selector interface {
	Select(service string) (*Node, error)
	Update(node *Node, err error) error
}

type Node struct {
	Network string // tcp, tcp4, tcp6 or udp, udp4, udp6
	Address string // ip:port
}
