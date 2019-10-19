package client

import (
	"context"
	"github.com/hitzhangjie/go-rpc/client/selector"
	"github.com/hitzhangjie/go-rpc/client/transport"
	"github.com/hitzhangjie/go-rpc/codec"
	"github.com/hitzhangjie/go-rpc/codec/whisper"
)

// Client client
type Client interface {
	Invoke(ctx context.Context, req interface{}, rsp interface{}, opts ...Option) error
}

func NewClient(name string, opts ...Option) Client {

	c := &client{
		//Name:      name,
		Selector:  nil,
		TransType: TCP,
		//Addr:      addr,
		Codec:     codec.ClientCodec(whisper.Whisper),
		Transport: &transport.TcpTransport{},
		RpcType:   SendRecv,
	}

	for _, o := range opts {
		o(c)
	}
	return c
}

type client struct {
	Name      string              // 请求服务名
	Addr      string              // 必填项
	Codec     codec.Codec         // 非必填，默认为whisper
	Selector  selector.Selector   // 非必填，默认为consul
	Transport transport.Transport //
	TransType TransportType       // 非必填，默认为tcp
	RpcType   RpcType             // 非必填，默认一发一收
}

func (c *client) Invoke(ctx context.Context, req interface{}, rsp interface{}, opts ...Option) error {

	session := codec.SessionFromContext(ctx)

	//data, err := c.Codec.Encode(req)

	return nil
}

