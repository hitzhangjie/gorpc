package transport

// EndPoint endpoint represents one side of net.conn
//
// Read read data from net.conn
// Write write data to net.conn
type EndPoint interface {
	Read()
	Write()
}
