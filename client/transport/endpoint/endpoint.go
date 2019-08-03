package endpoint

import (
	"net"
)

type EndPoint struct {
	Net  string
	Addr string
	Conn net.Conn
}
