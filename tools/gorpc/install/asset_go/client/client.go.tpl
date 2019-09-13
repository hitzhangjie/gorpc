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

	_ "github.com/hitzhangjie/go-rpc"
	"github.com/hitzhangjie/go-rpc/client"

    {{ with .FileOptions.go_package }}
    pb "{{.}}"
    {{ else }}
    pb "{{$.PackageName}}"
    {{ end }}

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
		client.WithCodecName("whisper"),
		client.WithNetwork("tcp4"),
		client.WithTarget("ip://127.0.0.1:8000"),
	}

	cli := pb.New{{$svrName|title}}Client(opts...)

    {{ $rpcReqType := "" -}}
    {{ $rpcRspType := "" -}}
	switch *cmd {
		{{range (index .Services 0).RPC}}

        {{- /* 根据rpc请求、响应类型，及pb中fileoptions信息计算正确的类型名 */ -}}
        {{- $validReqPkg := trimright "." (gofulltype .RequestType $.FileDescriptor) -}}
        {{- $validRspPkg := trimright "." (gofulltype .ResponseType $.FileDescriptor) -}}

        {{- if or (eq $validReqPkg $.PackageName) (eq $validReqPkg $goPkgOption) -}}
        	{{- $rpcReqType = (printf "pb.%s" (splitList "." .RequestType|last|export)) -}}
        {{- else -}}
        	{{- $rpcReqType = (gofulltype .RequestType $.FileDescriptor) -}}
        {{- end -}}

        {{- if or (eq $validReqPkg $.PackageName) (eq $validReqPkg $goPkgOption) -}}
        	{{- $rpcRspType = (printf "pb.%s" (splitList "." .ResponseType|last|export)) -}}
        {{- else -}}
        	{{- $rpcRspType = (gofulltype .ResponseType $.FileDescriptor) -}}
        {{- end -}}

		case "{{.Name|title}}":
			req := &{{$rpcReqType}}{}
			rsp, err := cli.{{.Name|title}}(ctx, req)
			log.Printf("req:%v, rsp:%v, err:%v", req, rsp, err)
		{{end}}
		default:
			log.Printf("err: undefined cmd:%v", *cmd)
			os.Exit(1)
	}
}
