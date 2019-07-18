{{- $pkgName := .PackageName -}}
{{- $svrName := .ServerName -}}
package exec

import (
	"context"
	"{{.ProtoSpec.RepoPrefix}}/{{$svrName}}"
	"git.code.oa.com/go-neat/core/nserver"
	"git.code.oa.com/go-neat/core/nserver/nsession"
	"git.code.oa.com/go-neat/tencent/attr"
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
