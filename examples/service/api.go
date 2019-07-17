package service

import (
	"context"
	"fmt"
	pb "gorpc/examples/helloworld"
)

type Greeter interface {
	SayHello(ctx context.Context, req pb.Request) (rsp pb.Response, err error)
}

var mapping = map[string]interface{}{}

func init() {
	mapping["SayHello"] = GreeterService.SayHello

	for n, m := range mapping {
		fmt.Println("call rpc:", n, ", result:", n, m.(*GreeterService).SayHello)
	}
}
