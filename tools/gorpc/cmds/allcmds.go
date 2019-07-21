package cmds

import (
	"sync"
)

var all = map[string]Commander{}
var mux sync.Mutex

// RegisteredSubCmds return all registered subcmds.
func RegisteredSubCmds() map[string]Commander {
	return all
}
