package connpool

import (
	"net"
)

type ConnectionPool interface {
	GetConn(addr string) (net.Conn, error)
	FreeConn(conn net.Conn) error
}
