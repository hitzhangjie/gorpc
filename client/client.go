package client

import (
	"context"
	"github.com/hitzhangjie/go-rpc/client/selector"
	"github.com/hitzhangjie/go-rpc/client/transport"
	"github.com/hitzhangjie/go-rpc/codec"
	"github.com/hitzhangjie/go-rpc/codec/whisper"
	"strings"
)

// Client client
type Client interface {
	Invoke(ctx context.Context, reqHead interface{}, rspHead interface{}, opts ...Option) error
}

func NewClient(name string, opts ...Option) Client {

	c := &client{
		//Name:      name,
		Selector:  nil,
		TransType: TCP,
		//Transport: &transport.TcpTransport{},
		//Address:      addr,
		Codec:     codec.ClientCodec(whisper.Whisper),
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

func (c *client) Invoke(ctx context.Context, reqHead interface{}, rspHead interface{}, opts ...Option) error {

	data, err := c.Codec.Encode(reqHead)
	if err != nil {
		return err
	}

	var (
		network string
		address string
	)

	if c.Addr != "" && c.TransType.Valid() {
		network = c.TransType.String()
		address = strings.TrimPrefix(c.Addr, "ip://")
	} else if c.Name != "" && c.Selector != nil {
		node, err := c.Selector.Select(c.Name)
		if err != nil {
			return err
		}
		network = node.Network
		address = node.Address
	}

	rsp, err := c.Transport.Send(ctx, network, address, data)
	if err != nil {
		return err
	}

	c.Codec.Decode(rsp)

	return nil
}
