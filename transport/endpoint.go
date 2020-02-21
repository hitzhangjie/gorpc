package transport

// EndPoint endpoint represents one side of net.Conn
//
// Read read data from net.Conn
// Write write data to net.Conn
type EndPoint interface {
	Read()
	Write()
}
