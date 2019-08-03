package balancer

type Balancer interface {
	Next() (addr string)
}
