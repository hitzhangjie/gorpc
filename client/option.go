package client

type Option func(adapter *ClientAdapter)

// ProtoType options
type ProtoType int

const (
	UDP = iota
	TCP
	UNIX
)

func ProtoTypeUDP(c *ClientAdapter) {
	c.ProtoType = UDP
}

func ProtoTypeTCP(c *ClientAdapter) {
	c.ProtoType = TCP
}

func ProtoTypeUNIX(c *ClientAdapter) {
	c.ProtoType = UNIX
}

// RpcType options
type RpcType int

const (
	SendOnly                = iota // 只发不收
	SendRecv                       // 一发一收
	SendRecvMultiplex              // 多路复用方式一发一收，发挥双工优势
	SendStreamOnly                 // 流式请求
	SendStreamAndRecv              // 流式请求
	SendAndRecvStream              // 流式请求
	SendStreamAndRecvStream        // 流式请求
)

func RpcTypeSendOnly(c *ClientAdapter) {
	c.RpcType = SendOnly
}

func RpcTypeSendRecv(c *ClientAdapter) {
	c.RpcType = SendRecv
}

func RpcTypeSendMultiplex(c *ClientAdapter) {
	c.RpcType = SendRecvMultiplex
}

func RpcTypeSendStreamOnly(c *ClientAdapter) {
	c.RpcType = SendStreamOnly
}

func RpcTypeSendStreamAndRecv(c *ClientAdapter) {
	c.RpcType = SendStreamAndRecv
}

func RpcTypeSendAndRecvStream(c *ClientAdapter) {
	c.RpcType = SendAndRecvStream
}

func RpcTypeSendStreamAndRecvStream(c *ClientAdapter) {
	c.RpcType = SendStreamAndRecvStream
}
