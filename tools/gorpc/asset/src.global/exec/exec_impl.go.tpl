{{- $pkgName := .PackageName -}}
{{- $svrName := .ServerName -}}
package exec

import (
    "context"
	"{{.ProtoSpec.RepoPrefix}}/{{$svrName}}"
	"git.code.oa.com/go-neat/core/nserver/nsession"
	{{range .Imports}}
    _ "{{.}}"
    {{end}}
)

{{range .RPC -}}
func {{.Name}}Impl(ctx context.Context, session nsession.NSession, req *{{.RequestType}}, rsp *{{.ResponseType}}) error {
	// business logic
	return nil
}

{{end -}}

