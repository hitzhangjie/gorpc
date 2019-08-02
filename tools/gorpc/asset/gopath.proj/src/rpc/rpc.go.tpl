{{- $pkgName := .PackageName -}}
{{- $svrName := .ServerName -}}
{{- $protocol := .Protocol -}}
package helloworld;

import (
    "context"
)

// {{$svrName|Title}}Server service definition
type {{$svrName|Title}}Server interface{
    // template: range rpc
    {{- range .RPC}}
    {{.Name}}(ctx context.Context, req *{{Simplify .RequestType $pkgName}}) returns(*{{Simplify .ResponseType $pkgName}}, error)
    {{- end}}
}

type gorpc.ServiceDesc struct {

}

func RegisterService()

// client
type GreeterClient interface{
    // template: range rpc
    SayHello(ctx context.Context, req *helloworld.HelloReq) returns(*helloworld.HelloRsp, error)
    SayByte(ctx context.Context, req *helloworld.ByteReq) returns(*helloworld.ByteRsp, error)
}

// -------------------------
// server
type GreeterByeServer interface{
    // template: range rpc
    SayByte(ctx context.Context, req *helloworld.ByteReq) returns(*helloworld.ByteRsp, error)
}

// client
type GreeterByeClient interface{
    // template: range rpc
    SayByte(ctx context.Context, req *helloworld.ByteReq) returns(*helloworld.ByteRsp, error)
}



