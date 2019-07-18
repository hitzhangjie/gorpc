{{$pkgname := .PackageName}}
package {{$pkgname}}

import (
	"context"
    "git.code.oa.com/go-neat/core/nclient"
    "git.code.oa.com/go-neat/core/nlog"
    "git.code.oa.com/go-neat/core/nserver/default_nserver"
    "git.code.oa.com/go-neat/core/nserver/nsession"
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
		//cltRpc = clt.NewNRPCRpcClient("test_nrpc", addr, proto, logger)
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

	{{- /*nrpc*/}}
	{{- if eq $protocol "nrpc"}}
	err = cltRpc.SendWithContext(ctx, session, "{{$serverName}}", "{{.Name}}", req, rsp)
	{{- end}}

    {{- /*ilive*/}}
	{{- if eq $protocol "ilive"}}
	err = cltRpc.SendWithContext(ctx, session, {{splitIliveCmd .Cmd}}, req, rsp)
	{{- end}}

    {{- /*simplesso*/}}
	{{- if eq $protocol "simplesso"}}
	err = cltRpc.SendWithContext(ctx, session, {{.Cmd}}, req, rsp)
	{{- end}}

	return rsp, err
}
{{end}}

