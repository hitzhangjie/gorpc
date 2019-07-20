{{- $pkgName := .PackageName -}}
{{- $svrName := .ServerName -}}
package exec

import (
	"context"
	"{{.ProtoSpec.RepoPrefix}}/{{$svrName}}"
	"gorpc/nserver"
	"gorpc/nserver/nsession"
	"gorpc/attr"
	{{range .Imports}}
    _ "{{.}}"
    {{end}}
)

{{range .RPC -}}
func {{.Name}}(ctx context.Context, session nsession.NSession) (interface{}, error) {
	req := &{{.RequestType}}{}
	err := session.ParseRequestBody(req)
	
	if err != nil {
		attr.Monitor(0, 1) //{{.RequestType}}解析失败
		session.Logger().Error("parse req err %v", err)
		return nil, err
	}
	
	rsp := &{{.ResponseType}}{}
	err = {{.Name}}Impl(ctx, session, req, rsp)
	if err != nil {
		attr.Monitor(0, 1) //{{.RequestType}}处理异常
		session.Logger().Error("handle req err %v", err)

		if _, ok := err.(nserver.Error); ok {
        	return nil, err
        }
        return nil, nserver.CreateError(nserver.EXEC_HANDLE_ERROR, err)
	}
	
	return rsp, nil
}

{{end -}}
