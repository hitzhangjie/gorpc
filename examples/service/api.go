package service

import (
	"context"
	"fmt"
	helloworld2 "github.com/hitzhangjie/gorpc/examples/helloworld"
)

type Greeter interface {
	SayHello(ctx context.Context, req helloworld2.Request) (rsp helloworld2.Response, err error)
}

var mapping = map[string]interface{}{}

func init() {
	mapping["SayHello"] = SayHello

	for n, m := range mapping {
		fmt.Println("call rpc:", n, ", result:", n, SayHello)
	}
}
