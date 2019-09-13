{{- $service := index .Services 0 -}}
{{- $svrName := $service.Name -}}
{{- $method := (index $service.RPC .RPCIndex) -}}
{{- $goPkgOption := "" -}}
{{- with .FileOptions.go_package -}}
  {{- $goPkgOption = . -}}
{{- end -}}
package main

{{ $rpcReqType := $method.RequestType -}}
{{- $rpcRspType := $method.ResponseType -}}

{{/* 计算rpc中请求类型、响应类型的正确类型名 */}}
{{ $validReqPkg := trimright "." (gofulltype $rpcReqType $.FileDescriptor) }}
{{ $validRspPkg := trimright "." (gofulltype $rpcRspType $.FileDescriptor) }}
{{- if or (eq $validReqPkg $.PackageName) (eq $validReqPkg $goPkgOption) -}}
	{{- $rpcReqType = (printf "pb.%s" (splitList "." $rpcReqType|last|export)) -}}
{{- else -}}
	{{- $rpcReqType = (gofulltype $rpcReqType $.FileDescriptor) -}}
{{- end -}}

{{- if or (eq $validRspPkg $.PackageName) (eq $validRspPkg $goPkgOption) -}}
	{{- $rpcRspType = (printf "pb.%s" (splitList "." $rpcRspType|last|export)) -}}
{{- else -}}
	{{- $rpcRspType = (gofulltype $rpcRspType $.FileDescriptor) -}}
{{- end -}}

import (
	"context"

{{/* 根据rpc请求、响应类型，判定是否需要引入当前proto对应的package */}}
{{ $reqTypePackage := (trimright "." $method.RequestType) }}
{{ $rspTypePackage := (trimright "." $method.ResponseType) }}

{{ if or (eq $reqTypePackage $.PackageName) (eq $rspTypePackage $.PackageName) }}
    {{ with .FileOptions.go_package }}
	pb "{{.}}"
    {{ else }}
	pb "{{$.PackageName}}"
    {{ end }}
{{- end }}

{{/* 根据rpc请求、响应类型，确定是否需要引入对应的package */}}
{{ $reqTypeImportPath := (index $.ImportPathMappings $reqTypePackage) }}
{{ $rspTypeImportPath := (index $.ImportPathMappings $rspTypePackage) }}

{{ range .Imports }}
{{- if or (eq . $reqTypeImportPath) (eq . $rspTypeImportPath) }}
    "{{.}}"
{{ end -}}
{{ end -}}
)

func (s *{{$svrName|title}}ServerImpl) {{$method.Name|title}}(ctx context.Context, req *{{$rpcReqType}}) (rsp *{{$rpcRspType}}, err error) {
	// implement business logic here ...
	// ...

	return
}

