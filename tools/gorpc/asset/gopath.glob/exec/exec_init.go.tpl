package exec

import (
	"gorpc/nserver/default_nserver"
)

func init() {
    {{- range .RPC}}
    default_nserver.AddExec("{{.Cmd}}", {{.Name}})
    {{- end}}
}
