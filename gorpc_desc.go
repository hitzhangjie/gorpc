package gorpc

import (
	"context"
)

type ServiceDesc struct {
	Name        string                // "helloworld.greeter"
	ServiceType interface{}           // (*helloworld.GreeterServer)(nil)
	Method      map[string]MethodDesc // "SayHello"
	Stream      map[string]StreamDesc
}

type MethodDesc struct {
	Name   string
	Method func(service interface{}, ctx context.Context, session Session)
}

type StreamDesc struct {
}


