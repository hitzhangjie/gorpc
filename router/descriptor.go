package router

import (
	"context"
	"github.com/hitzhangjie/go-rpc/codec"
)

type ServiceDesc struct {
	Name        string                 // "helloworld.greeter"
	ServiceType interface{}            // (*helloworld.GreeterServer)(nil)
	Method      map[string]*MethodDesc // "SayHello"
	Stream      map[string]*StreamDesc
}

type HandleFunc = func(svr interface{}, ctx context.Context, session codec.Session) error

type MethodDesc struct {
	Name   string
	Method HandleFunc
}

type StreamDesc struct {
}
