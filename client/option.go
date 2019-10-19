package client

type Option func(adapter *client)

// ProtoType options
type ProtoType int

const (
	UDP = iota
	TCP
	UNIX
)

func ProtoTypeUDP(c *client) {
	c.ProtoType = UDP
}

func ProtoTypeTCP(c *client) {
	c.ProtoType = TCP
}

func ProtoTypeUNIX(c *client) {
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

func RpcTypeSendOnly(c *client) {
	c.RpcType = SendOnly
}

func RpcTypeSendRecv(c *client) {
	c.RpcType = SendRecv
}

func RpcTypeSendMultiplex(c *client) {
	c.RpcType = SendRecvMultiplex
}

func RpcTypeSendStreamOnly(c *client) {
	c.RpcType = SendStreamOnly
}

func RpcTypeSendStreamAndRecv(c *client) {
	c.RpcType = SendStreamAndRecv
}

func RpcTypeSendAndRecvStream(c *client) {
	c.RpcType = SendAndRecvStream
}

func RpcTypeSendStreamAndRecvStream(c *client) {
	c.RpcType = SendStreamAndRecvStream
}
