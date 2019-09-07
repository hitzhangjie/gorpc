{{- $svrName := (index .Services 0).Name -}}
{{- $pkgName := .PackageName -}}

{{- $goPkgOption := "" -}}
{{- with .FileOptions.go_package -}}
  {{- $goPkgOption = . -}}
{{- end -}}

{{- if eq .GoMod "" -}}
module {{$svrName}}
{{- else -}}
module {{.GoMod}}
{{- end }}

go 1.12

{{ if ne $goPkgOption "" -}}
replace {{$goPkgOption}} => ./{{$goPkgOption}}
{{- else -}}
replace {{$pkgName}} => ./{{$pkgName}}
{{- end -}}
