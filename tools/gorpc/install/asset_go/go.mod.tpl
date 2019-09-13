{{- if eq .GoMod "" -}}
module {{ (index .Services 0).Name }}
{{- else -}}
module {{.GoMod}}
{{- end }}

go 1.12

{{ with .FileOptions.go_package }}
replace {{.}} => ./{{.}}
{{ else }}
replace {{.PackageName}} => ./{{.PackageName}}
{{ end }}
