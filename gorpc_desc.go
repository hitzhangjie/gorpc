package gorpc

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

type MethodDesc struct {
	Name   string
	Method func(service interface{}, ctx context.Context, session codec.Session)
}

type StreamDesc struct {
}
