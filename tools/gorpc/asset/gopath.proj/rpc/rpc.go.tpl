{{$pkgname := .PackageName}}
package {{$pkgname}}

import (
	"context"
    "gorpc/nclient"
    "gorpc/nlog"
    "gorpc/nserver/default_nserver"
    "gorpc/nserver/nsession"
    clt "{{.ProtoSpec.ClientPkg}}"
    "sync"
)

var cltRpc *clt.{{.ProtoSpec.ClientType}}
var once = new(sync.Once)

func init() {
    logger := default_nserver.GetSvrLog()
    if logger == nil {
    	return
    }
    Init("l5://1:1", nclient.ConnType(0), logger)
}

func Init(addr string, proto nclient.ConnType, logger *nlog.NLog) {
	once.Do(func() {
		cltRpc = clt.{{.ProtoSpec.ClientFactory}}("{{.ServerName}}", addr, proto, logger)
	})
}

{{$protocol := .Protocol}}
{{$serverName := .ServerName}}

{{range .RPC}}
func {{.Name}}(ctx context.Context, session nsession.NSession, req *{{.RequestTypeNameInRpcTpl}}) (rsp *{{.ResponseTypeNameInRpcTpl}}, err error) {
	if rsp == nil {
		rsp = &{{.ResponseTypeNameInRpcTpl}}{}
	}
	if cltRpc == nil {
		return nil, nclient.CreateErrorWithMsg(20, "client not init", false)
	}

	{{- /*gorpc*/}}
	{{- if eq $protocol "gorpc"}}
	err = cltRpc.SendWithContext(ctx, session, "{{$serverName}}", "{{.Name}}", req, rsp)
	{{- end}}

    {{- /*swan*/}}
	{{- if eq $protocol "swan"}}
	err = cltRpc.SendWithContext(ctx, session, {{splitSwanCmd .Cmd}}, req, rsp)
	{{- end}}

    {{- /*chick*/}}
	{{- if eq $protocol "chick"}}
	err = cltRpc.SendWithContext(ctx, session, {{.Cmd}}, req, rsp)
	{{- end}}

	return rsp, err
}
{{end}}

