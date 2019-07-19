package main

import (
	gorpc2 "github.com/hitzhangjie/gorpc"
	"github.com/hitzhangjie/gorpc/examples/service"
)

func main() {
	serv := gorpc2.NewService("grpc/helloworld.greeter").Version("1.0")
	serv.Handle(service.GreeterService{})
	serv.Start()
}
