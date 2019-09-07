{{- $svrName := (index .Services 0).Name -}}
{{- $pkgName := .PackageName -}}
{{- $goPkgOption := "" -}}
{{- with .FileOptions.go_package -}}
  {{- $goPkgOption = . -}}
{{- end -}}
{{- if ne $goPkgOption "" -}}
package {{ (splitList "/" $goPkgOption)|last|gopkg}}
{{- else -}}
package {{ $pkgName|gopkg }}
{{- end }}

import (
	"context"

   _ "github.com/hitzhangjie/go-rpc"
   _ "github.com/hitzhangjie/go-rpc/http"

    "github.com/hitzhangjie/go-rpc/server"
    "github.com/hitzhangjie/go-rpc/client"
    "github.com/hitzhangjie/go-rpc/codec"
    //gorpc "github.com/hitzhangjie/go-rpc"
    {{ range .Imports }}
    "{{- . -}}"
    {{- end }}
)

/* ************************************ Service Definition ************************************ */

// {{$svrName|title}} defines service
type {{$svrName|title}}Server interface {

	{{ range (index .Services 0).RPC }}
	{{- $rpcReqType := (simplify (gofulltype .RequestType $.FileDescriptor) $pkgName)|export }}
	{{- $rpcRspType := (simplify (gofulltype .ResponseType $.FileDescriptor) $pkgName)|export }}
	// {{.Name|title}} {{.LeadingComments}}
	{{.Name|title }}(ctx context.Context, req *{{$rpcReqType}},rsp *{{$rpcRspType}}) (err error) // {{.TrailingComments}}
{{ end -}}
}

{{range (index .Services 0).RPC -}}
func {{$svrName|title}}Server_{{.Name|title}}_Handler(svr interface{}, ctx context.Context, f server.FilterFunc) (rspbody interface{}, err error) {
    {{- $rpcReqType := (simplify (gofulltype .RequestType $.FileDescriptor) $pkgName)|export }}
    {{- $rpcRspType := (simplify (gofulltype .ResponseType $.FileDescriptor) $pkgName)|export }}

    req := &{{$rpcReqType}}{}
	rsp := &{{$rpcRspType}}{}
	filters, err := f(req)
    if err != nil {
    	return nil, err
    }
	handleFunc := func(ctx context.Context, reqbody interface{}, rspbody interface{}) error {
	    return svr.({{$svrName|title}}Server).{{.Name|title}}(ctx, reqbody.(*{{$rpcReqType}}), rspbody.(*{{$rpcRspType}}))
	}

	err = filters.Handle(ctx, req, rsp, handleFunc)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

{{end -}}

// {{$svrName|title}}Server_ServiceDesc descriptor for server.RegisterService
var {{$svrName|title}}Server_ServiceDesc = server.ServiceDesc {
    ServiceName: "{{$pkgName}}.{{$svrName}}",
    HandlerType: ((*{{$svrName|title}}Server)(nil)),
    Methods: []server.Method{
        {{- range (index .Services 0).RPC}}
	        {Name: "{{.FullyQualifiedCmd}}", Func: {{$svrName|title}}Server_{{.Name|title}}_Handler},
        {{- end}}
    },
}

// Register{{$svrName|title}}Server register service
func Register{{$svrName|title}}Server(s server.Service, svr {{$svrName|title}}Server) {
	s.Register(&{{$svrName|title}}Server_ServiceDesc, svr)
}

/* ************************************ Client Definition ************************************ */

// {{$svrName|title}}ClientProxy defines service client proxy
type {{$svrName|title}}ClientProxy interface {
	{{ range (index .Services 0).RPC}}
	{{- $rpcReqType := (simplify (gofulltype .RequestType $.FileDescriptor) $pkgName)|export }}
   	{{- $rpcRspType := (simplify (gofulltype .ResponseType $.FileDescriptor) $pkgName)|export }}
   	// {{.Name|title}} {{.LeadingComments}}
	{{.Name|title}}(ctx context.Context, req *{{$rpcReqType}}, opts ...client.Option) (rsp *{{$rpcRspType}}, err error) // {{.TrailingComments}}
{{ end -}}
}

type {{$svrName|title}}ClientProxyImpl struct{
	client client.Client
	opts []client.Option
}

func New{{$svrName|title}}ClientProxy(opts...client.Option) {{$svrName|title}}ClientProxy {
	return &{{$svrName|title}}ClientProxyImpl {client: client.DefaultClient, opts: opts}
}

{{range $idx, $rpc := (index .Services 0).RPC}}
{{- $rpcReqType := (simplify (gofulltype .RequestType $.FileDescriptor) $pkgName)|export }}
{{- $rpcRspType := (simplify (gofulltype .ResponseType $.FileDescriptor) $pkgName)|export }}
func (c *{{$svrName|title}}ClientProxyImpl) {{.Name|title}}(ctx context.Context, req *{{$rpcReqType}}, opts ...client.Option) (rsp *{{$rpcRspType}}, err error) {

	ctx, msg := codec.WithCloneMessage(ctx)

	msg.WithClientRpcName({{$svrName|title}}Server_ServiceDesc.Methods[{{$idx}}].Name)
	msg.WithCalleeServiceName({{$svrName|title}}Server_ServiceDesc.ServiceName)

	callopts := make([]client.Option, 0, len(c.opts)+len(opts))
	callopts = append(callopts, c.opts...)
	callopts = append(callopts, opts...)
	rsp = &{{$rpcRspType}}{}

	err = c.client.Invoke(ctx, req, rsp, callopts...)
	if err != nil {
	    return nil, err
	}

	return rsp, nil
}
{{end}}
