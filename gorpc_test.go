package go_rpc

import (
	"github.com/hitzhangjie/go-rpc/examples/service"
	"testing"
)

func TestService_Handle(t *testing.T) {
	g := &service.GreeterService{}

	s := Version("1.0.0")
	Handle(g)
}
