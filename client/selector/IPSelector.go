package selector

type IPSelector struct {
}

func (s *IPSelector) Select(service string) (*Node, error) {
	panic("implement me")
}

func (s *IPSelector) Update(node *Node, err error) error {
	panic("implement me")
}
