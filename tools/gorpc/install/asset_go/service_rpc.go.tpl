{{- $svrName := (index .Services 0).Name -}}
{{- $pkgName := .PackageName -}}
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
{{- if or (eq (trimright "." $rpcReqType) ($pkgName)) (eq (trimright "." (gofulltype $rpcReqType $.FileDescriptor)) ($goPkgOption)) -}}
	{{- $rpcReqType = (printf "pb.%s" (splitList "." $rpcReqType|last|export)) -}}
{{- else -}}
	{{- $rpcReqType = (gofulltype $rpcReqType $.FileDescriptor) -}}
{{- end -}}

{{- if or (eq (trimright "." $rpcRspType) $pkgName) (eq (trimright "." (gofulltype $rpcRspType $.FileDescriptor)) $goPkgOption) -}}
	{{- $rpcRspType = (printf "pb.%s" (splitList "." $rpcRspType|last|export)) -}}
{{- else -}}
	{{- $rpcRspType = (gofulltype $rpcRspType $.FileDescriptor) -}}
{{- end -}}

import (
	"context"

{{/* 根据rpc请求、响应类型，判定是否需要引入当前proto对应的package */}}
{{ if or (eq (index (splitList "." $rpcReqType) 0) "pb") (eq (index (splitList "." $rpcRspType) 0) "pb") }}
{{ if ne $goPkgOption "" }}
	pb "{{ $goPkgOption }}"
{{- else }}
	pb "{{$pkgName|gopkg -}}"
{{- end }}
{{- end }}

{{/* 根据rpc请求、响应类型，确定是否需要引入对应的package */}}
{{ range .Imports }}
{{- $importPkg := . }}
{{- if or (hasprefix $importPkg $method.RequestType) (hasprefix $importPkg $method.ResponseType) }}
{{- if or (ne (index (splitList "." $rpcReqType) 0) "pb") (ne (index (splitList "." $rpcRspType) 0) "pb") }}
    "{{ $importPkg }}"
{{ end -}}
{{ end -}}
{{ end -}}
)

func (s *{{$svrName|title}}ServerImpl) {{$rpcName|title}}(ctx context.Context, req *{{$rpcReqType}}, rsp *{{$rpcRspType}}) (err error) {
	// implement business logic here ...
	// ...

	return
}

