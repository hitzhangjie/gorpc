package helloworld;

import (
    "context"
)

// template: range service

// ------------------------
// server
type GreeterServer interface{
    // template: range rpc
    SayHello(ctx context.Context, req *helloworld.HelloReq) returns(*helloworld.HelloRsp, error)
    SayByte(ctx context.Context, req *helloworld.ByteReq) returns(*helloworld.ByteRsp, error)
}

type gorpc.ServiceDesc struct {

}

func RegisterService()

// client
type GreeterClient interface{
    // template: range rpc
    SayHello(ctx context.Context, req *helloworld.HelloReq) returns(*helloworld.HelloRsp, error)
    SayByte(ctx context.Context, req *helloworld.ByteReq) returns(*helloworld.ByteRsp, error)
}




// -------------------------
// server
type GreeterByeServer interface{
    // template: range rpc
    SayByte(ctx context.Context, req *helloworld.ByteReq) returns(*helloworld.ByteRsp, error)
}

// client
type GreeterByeClient interface{
    // template: range rpc
    SayByte(ctx context.Context, req *helloworld.ByteReq) returns(*helloworld.ByteRsp, error)
}



