{{- $svrName := (index .Services 0).Name -}}
{{- $pkgName := .PackageName -}}
{{- $goPkgOption := "" -}}

{{- with .FileOptions.go_package -}}
  {{- $goPkgOption = . -}}
{{- end -}}
package main

import (
	"context"
	"flag"
	"os"
	"time"
	"log"

	_ "git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/client"

    {{ if ne $goPkgOption "" -}}
	pb "{{$goPkgOption}}"
    {{- else -}}
	pb "{{$pkgName|gopkg}}"
    {{- end }}
    {{ range .Imports }}
	"{{- . -}}"
    {{- end }}
)

// common options
var addr = flag.String("addr", "ip://127.0.0.1:8000", "addr, supporting ip://<ip>:<port>, l5://mid:cid, cmlb://appid[:sysid]")
var cmd = flag.String("cmd", "{{(index (index .Services 0).RPC 0).Name|title}}", "{{range (index .Services 0).RPC}}, {{.Name|title}}{{end}}")
var timeout = flag.Int("timeout", 2000, "timeout in milliseconds")

func init() {
	flag.Parse()
	log.SetFlags(log.LstdFlags|log.Lshortfile)
}

func main() {

	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*2000)
	defer cancel()

	opts := []client.Option{
		client.WithCodecName("trpc"),
		client.WithCheckerName("trpc"),
		client.WithNetwork("tcp4"),
		client.WithTarget("ip://127.0.0.1:8000"),
	}

	clientProxy := pb.New{{$svrName|title}}ClientProxy(opts...)

	switch *cmd {
		{{range (index .Services 0).RPC}}
        {{- $rpcReqType := .RequestType -}}
        {{- $rpcRspType := .ResponseType -}}

        {{- if or (eq (trimright "." $rpcReqType|gopkg) ($pkgName|gopkg)) (eq (trimright "." (gofulltype $rpcReqType $.FileDescriptor)|gopkg) ($goPkgOption|gopkg)) -}}
        	{{- $rpcReqType = (printf "pb.%s" (splitList "." $rpcReqType|last|export)) -}}
        {{- else -}}
        	{{- $rpcReqType = (gofulltype $rpcReqType $.FileDescriptor) -}}
        {{- end -}}

        {{- if or (eq (trimright "." $rpcRspType|gopkg) ($pkgName|gopkg)) (eq (trimright "." (gofulltype $rpcRspType $.FileDescriptor)|gopkg) ($goPkgOption|gopkg)) -}}
        	{{- $rpcRspType = (printf "pb.%s" (splitList "." $rpcRspType|last|export)) -}}
        {{- else -}}
        	{{- $rpcRspType = (gofulltype $rpcRspType $.FileDescriptor) -}}
        {{- end -}}
		case "{{.Name|title}}":
			req := &{{$rpcReqType}}{}
			rsp, err := clientProxy.{{.Name|title}}(ctx, req)
			log.Printf("req:%v, rsp:%v, err:%v", req, rsp, err)
		{{end}}
		default:
			log.Printf("err: undefined cmd:%v", *cmd)
			os.Exit(1)
	}
}
