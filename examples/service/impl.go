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
