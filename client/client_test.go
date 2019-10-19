package client

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/hitzhangjie/go-rpc/client/selector"
	"github.com/hitzhangjie/go-rpc/codec"
	"github.com/hitzhangjie/go-rpc/codec/whisper"
	"sync/atomic"
	"testing"
)

func newClient(name, address, codec string, selector selector.Selector) Client {

	opts := []Option{
		WithAddress(address),
		WithCodec(codec),
		WithSelector(selector),
	}
	client := NewClient(name, opts...)
	return client
}

func TestNewClient(t *testing.T) {

	client := newClient("greeter", "ip://127.0.0.1:8888", "whisper", &selector.IPSelector{})
	t.Logf("client:%+v", client)
}

func TestClientInvoke(t *testing.T) {
	client := newClient("greeter", "ip://127.0.0.1:8888", "whisper", &selector.IPSelector{})

	ctx := context.Background()

	err := client.Invoke(ctx, "Hello", req, rsp, opts...)
}

type XXXXClient struct {
	client
}

var seqno uint64

func (c *XXXXClient) Hello(ctx context.Context, req, rsp interface{}) error {

	rpcName := "Hello"

	reqHead := &whisper.Request{
		Seqno:   proto.Uint64(atomic.AddUint64(&seqno, 1)),
		Rpcname: proto.String(rpcName),
	}
	rspHead := &whisper.Response{}

	data, err := proto.Marshal(req.(proto.Message))
	if err != nil {
		return err
	}

	reqHead.Body = data

	err := c.Invoke(ctx, reqHead, rspHead)
	if err != nil {
		return err
	}


	return nil
}
