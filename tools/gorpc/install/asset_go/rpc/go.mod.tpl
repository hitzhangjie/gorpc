{{- with .FileOptions.go_package -}}
module {{.}}
{{- else -}}
module {{$.PackageName}}
{{- end }}

go 1.12

require github.com/hitzhangjie/go-rpc v0.0.0-20190901020304-010203040506
