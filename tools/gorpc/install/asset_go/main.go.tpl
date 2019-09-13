package main

import (
	gorpc "github.com/hitzhangjie/go-rpc"

	{{ with .FileOptions.go_package }}
	pb "{{.}}"
	{{ else }}
	pb "{{.PackageName}}"
	{{ end }}
)

{{- $svrName := (index .Services 0).Name }}
type {{$svrName|title}}ServerImpl struct {}

func main() {

	s := gorpc.NewServer()

	pb.Register{{$svrName|title}}Server(s, &{{$svrName|title}}ServerImpl{})
	s.Serve()
}