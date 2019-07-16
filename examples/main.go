package main

import (
	"gorpc"
	"gorpc/examples/service"
)

func main() {
	serv := gorpc.NewService("grpc/helloworld.greeter").Version("1.0")
	serv.Handle(service.GreeterService{})
	serv.Start()
}
