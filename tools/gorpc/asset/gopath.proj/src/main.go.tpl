{{- $pkgName := .PackageName -}}
{{- $svrName := .ServerName -}}
{{- $protocol := .Protocol -}}
package main

import (
    "github.com/hitzhangjie/go-rpc"
)

func main() {
    gorpc.NewService("{{$pkgName}}.{{$svrName|Title}}").Version("1.0.0").Start()
}
