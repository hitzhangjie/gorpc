package client

import (
	"context"
	"github.com/hitzhangjie/go-rpc/client/selector"
	"github.com/hitzhangjie/go-rpc/client/transport"
	"github.com/hitzhangjie/go-rpc/codec"
	"github.com/hitzhangjie/go-rpc/codec/whisper"
)

type Client interface {
	Invoke(ctx context.Context, req interface{}, rsp interface{}, opts ...Option) error
}

type ClientAdapter struct {
	Addr      string              // 必填项
	Codec     codec.Codec         // 非必填，默认为whisper
	Selector  selector.Selector   // 非必填，默认为consul
	Transport transport.Transport //
	ProtoType ProtoType           // 非必填，默认为tcp
	RpcType   RpcType             // 非必填，默认一发一收
}

func (c *ClientAdapter) Invoke(ctx context.Context, req interface{}, rsp interface{}, opts ...Option) error {
	return nil
}

func NewClientAdapter(protoType ProtoType, addr, codecName string, rpcType RpcType, opts ...Option) (*ClientAdapter, error) {

	c := &ClientAdapter{
		Selector:  nil,
		ProtoType: TCP,
		Addr:      addr,
		Codec:     codec.ClientCodec(whisper.Whisper),
		Transport: &transport.TcpTransport{},
		RpcType:   SendRecv,
	}

	for _, o := range opts {
		o(c)
	}
	return c, nil
}
