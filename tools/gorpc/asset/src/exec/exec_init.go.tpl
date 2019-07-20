package exec

import (
	"gorpc/nserver/default_nserver"
)

func init() {
	//注册服务接口
    {{- range .RPC}}
    default_nserver.AddExec("{{.Cmd}}", {{.Name}})
    {{- end}}

    {{- if .HttpOn}}
    //注册http接口
    {{- range .RPC}}
    default_nserver.AddExec("/{{.Cmd}}", {{.Name}})
    {{- end}}
    {{- end}}
}
