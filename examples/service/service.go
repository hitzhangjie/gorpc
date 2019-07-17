package service

import (
	"context"
	pb "gorpc/examples/helloworld"
)

type GreeterService struct {
}

func (s *GreeterService) SayHello(ctx context.Context, req pb.Request) (rsp pb.Response, err error) {

	return
}

func init() {

	var e interface{}
	e = &GreeterService{}

	if _, ok := e.(Greeter); !ok {
		panic("GreeterService not implement Greeter")
	}
}
