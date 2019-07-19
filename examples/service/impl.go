package service

import (
	"context"
	helloworld2 "github.com/hitzhangjie/gorpc/examples/helloworld"
)

type GreeterService struct {
}

func (s *GreeterService) SayHello(ctx context.Context, req helloworld2.Request) (rsp helloworld2.Response, err error) {

	return
}
