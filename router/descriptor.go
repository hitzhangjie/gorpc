package router

import (
	"context"
)

type ServiceDesc struct {
	Name        string                 // "helloworld.greeter"
	ServiceType interface{}            // (*helloworld.GreeterServer)(nil)
	Method      map[string]*MethodDesc // "SayHello"
	Stream      map[string]*StreamDesc
}

type HandleFunc = func(svr interface{}, ctx context.Context, req interface{}) (rsp interface{}, err error)

type MethodDesc struct {
	Name   string
	Method HandleFunc
}

type StreamDesc struct {
}
