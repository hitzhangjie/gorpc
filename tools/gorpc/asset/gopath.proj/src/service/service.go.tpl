package service

import (
    "context"
    "../rpc/helloworld"     // template: imports
)

// template: range service
type GreeterServerImpl struct {}

// template: range rpc
func (s *GreeterServerImpl) SayHello(ctx context.Context, req *helloworld.HelloReq) (*helloworld.HelloRsp, error) {
    // your logic goes here
    return nil, nil
}

func (s *GreeterServerImpl) SayBye(ctx context.Context, req *helloworld.ByeReq) (*helloworld.ByeRsp, error) {
    // your logic goes here
    return nil, nil
}

type GreeterByeServerImpl struct {}

func (s *GreeterByeServerImpl) SayBye(ctx context.Context, req *helloworld.ByeReq) (*helloworld.ByeRsp, error) {
    // your logic goes here
    return nil, nil
}
