package service

import (
	"context"
	pb "gorpc/examples/helloworld"
)

type Greeter interface {
	SayHello(ctx context.Context, req pb.Request) (rsp pb.Response, err error)
}
