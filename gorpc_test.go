package go_rpc

import (
	"gorpc/examples/service"
	"testing"
)

func TestService_Handle(t *testing.T) {
	g := &service.GreeterService{}

	s := NewService("gorpc/app/helloworld.GreeterService").Version("1.0.0")
	s.Handle(g)

}
