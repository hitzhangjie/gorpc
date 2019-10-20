package client

import (
	"context"
	"github.com/hitzhangjie/go-rpc/client/pool"
	"github.com/hitzhangjie/go-rpc/client/selector"
	"github.com/hitzhangjie/go-rpc/client/transport"
	"github.com/hitzhangjie/go-rpc/codec"
	"github.com/hitzhangjie/go-rpc/codec/whisper"
	"strings"
	"time"
)

var defaultPoolFactory = pool.NewConnPoolFactory(
	pool.WithMinIdle(2),
	pool.WithMaxIdle(4),
	pool.WithMaxActive(8),
	pool.WithDialTimeout(time.Second*2),
	pool.WithIdleTimeout(time.Minute*5),
	pool.WithMaxConnLifetime(time.Minute*30),
	pool.WithWait(true),
)

// Client client
type Client interface {
	Invoke(ctx context.Context, reqHead interface{}, opts ...Option) (rspHead interface{}, err error)
}

func NewClient(name string, opts ...Option) Client {

	c := &client{
		//Name:      name,
		Selector:  nil,
		TransType: TCP,
		//Transport: &transport.TcpTransport{},
		//Address:      addr,
		Codec:   codec.ClientCodec(whisper.Whisper),
		RpcType: SendRecv,
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

func (c *client) Invoke(ctx context.Context, reqHead interface{}, opts ...Option) (rspHead interface{}, err error) {

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
			return nil, err
		}
		network = node.Network
		address = node.Address
	}

	rsp, err := c.Transport.Send(ctx, network, address, reqHead)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}
