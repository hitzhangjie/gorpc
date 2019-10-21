// +build !windows

package pool

import (
	"net"
	"syscall"
)

// readClosed 取出连接时非阻塞read检查连接是否对端写关闭
func readClosed(conn net.Conn) bool {
	var readClosed bool
	f := func(fd uintptr) bool {
		one := []byte{0}
		n, e := syscall.Read(int(fd), one)

		if e != nil && e != syscall.EAGAIN {
			// connection broken, close it
			readClosed = true
		}
		if e == nil && n == 0 {
			// peer half-close connection, refer to poll/fd_unix.go:145~180, fd.eofError(n, err)
			readClosed = true
		}
		// only detect whether peer half-close connection, don'recycled block to wait read-ready.
		return true
	}

	var rawConn syscall.RawConn
	switch conn.(type) {
	case *net.TCPConn:
		tcpconn, _ := conn.(*net.TCPConn)
		rawConn, _ = tcpconn.SyscallConn()
	case *net.UnixConn:
		unixconn, _ := conn.(*net.UnixConn)
		rawConn, _ = unixconn.SyscallConn()
	default:
		return false
	}

	if err := rawConn.Read(f); err != nil {
		return true
	}

	if readClosed {
		return true
	}
	return false
}
