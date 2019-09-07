{{- $svrName := (index .Services 0).Name -}}
{{- $pkgName := .PackageName -}}
{{- $goPkgOption := "" -}}
{{- with .FileOptions.go_package -}}
  {{- $goPkgOption = . -}}
{{- end -}}
package main

import (
	gorpc "github.com/hitzhangjie/go-rpc"

    {{ if ne $goPkgOption "" -}}
   	pb "{{$goPkgOption}}"
    {{- else -}}
    pb "{{$pkgName}}"
	{{- end }}
)

type {{$svrName|title}}ServerImpl struct {}

func main() {

	s := gorpc.NewServer()

	pb.Register{{$svrName|title}}Server(s, &{{$svrName|title}}ServerImpl{})
	s.Serve()
}