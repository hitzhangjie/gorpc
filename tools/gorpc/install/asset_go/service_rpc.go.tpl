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
{{ $validReqPkgName := trimright "." (gofulltype $rpcReqType $.FileDescriptor) }}
{{ $validRspPkgName := trimright "." (gofulltype $rpcRspType $.FileDescriptor) }}
{{- if or (eq $validReqPkgName $.PackageName) (eq $validReqPkgName $goPkgOption) -}}
	{{- $rpcReqType = (printf "pb.%s" (splitList "." $rpcReqType|last|export)) -}}
{{- else -}}
	{{- $rpcReqType = (gofulltype $rpcReqType $.FileDescriptor) -}}
{{- end -}}

{{- if or (eq $validRspPkgName $.PackageName) (eq $validRspPkgName $goPkgOption) -}}
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


func (s *{{$svrName|title}}ServerImpl) {{$method.Name|title}}(ctx context.Context, req *{{$rpcReqType}}) (rsp *{{$rpcRspType}}, err error) {
	// implement business logic here ...
	// ...

	return
}

