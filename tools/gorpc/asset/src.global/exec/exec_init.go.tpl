package exec

import (
	"git.code.oa.com/go-neat/core/nserver/default_nserver"
)

func init() {
    {{- range .RPC}}
    default_nserver.AddExec("{{.Cmd}}", {{.Name}})
    {{- end}}
}
