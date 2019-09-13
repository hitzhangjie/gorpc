{{- $svrName := (index .Services 0).Name -}}
{{- $goPkgOption := "" -}}
{{- with .FileOptions.go_package -}}
  {{- $goPkgOption = . -}}
{{- end -}}
package main

{{ $service0 := index .Services 0 -}}
{{ $method := (index $service0.RPC .RPCIndex) -}}
{{- $rpcName := $method.Name -}}
{{- $rpcReqType := $method.RequestType -}}
{{- $rpcRspType := $method.ResponseType -}}

{{/* 计算rpc中请求类型、响应类型的正确类型名 */}}
{{- if or (eq (trimright "." $rpcReqType) ($.PackageName)) (eq (trimright "." (gofulltype $rpcReqType $.FileDescriptor)) ($goPkgOption)) -}}
	{{- $rpcReqType = (printf "pb.%s" (splitList "." $rpcReqType|last|export)) -}}
{{- else -}}
	{{- $rpcReqType = (gofulltype $rpcReqType $.FileDescriptor) -}}
{{- end -}}

{{- if or (eq (trimright "." $rpcRspType) $.PackageName) (eq (trimright "." (gofulltype $rpcRspType $.FileDescriptor)) $goPkgOption) -}}
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
{{ if ne $goPkgOption "" }}
	pb "{{ $goPkgOption }}"
{{- else }}
	pb "{{$.PackageName|gopkg -}}"
{{- end }}
{{- end }}

{{/* 根据rpc请求、响应类型，确定是否需要引入对应的package */}}
{{ $reqTypeImportPkg := (index $.PkgImportMappings $reqTypePackage) }}
{{ $rspTypeImportPkg := (index $.PkgImportMappings $rspTypePackage) }}

{{ range .Imports }}
{{- if or (hasprefix . $reqTypeImportPkg) (hasprefix . $rspTypeImportPkg) }}
    "{{.}}"
{{ end -}}
{{ end -}}
)

func (s *{{$svrName|title}}ServerImpl) {{$rpcName|title}}(ctx context.Context, req *{{$rpcReqType}}, rsp *{{$rpcRspType}}) (err error) {
	// implement business logic here ...
	// ...

	return
}

