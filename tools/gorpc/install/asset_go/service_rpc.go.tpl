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

import (
	"context"
{{ if or (eq (index (splitList "." $rpcReqType) 0) "pb") (eq (index (splitList "." $rpcRspType) 0) "pb") }}
{{ if ne $goPkgOption "" }}
	pb "{{ $goPkgOption }}"
{{- else }}
	pb "{{$pkgName|gopkg -}}"
{{- end }}
{{- end }}
{{ range .Imports }}
{{- $importPkg := . }}
{{- if or (hasprefix $importPkg $rpcReqType) (hasprefix $importPkg $rpcRspType) }}
{{- if or (ne (index (splitList "." $rpcReqType) 0) "pb") (ne (index (splitList "." $rpcRspType) 0) "pb") }}
    "{{ $importPkg }}"
{{ end -}}
{{ end -}}
{{ end -}}
)

func (s *{{$svrName|title}}ServerImpl) {{$rpcName|title}}(ctx context.Context, req *{{$rpcReqType}},rsp *{{$rpcRspType}}) (err error) {
	// implement business logic here ...
	// ...

	return
}

