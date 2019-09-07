package cmds

import (
	"sync"
)

var (
	cmds = map[string]Commander{}
	mux  sync.Mutex
)

// SubCmds return cmds registered subcmds.
func SubCmds() map[string]Commander {
	return cmds
}

func RegisterSubCmd(name string, commander Commander) {
	mux.Lock()
	cmds[name] = commander
	mux.Unlock()
}

func init() {
	RegisterSubCmd("create", newCreateCmd())
	RegisterSubCmd("update", newUpdateCmd())
	RegisterSubCmd("help", newHelpCmd())
}
