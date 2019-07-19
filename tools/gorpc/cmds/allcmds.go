package cmds

import (
	"sync"
)

var all map[string]Commander = map[string]Commander{}
var alllock sync.Mutex

// RegisteredSubCmds return all registered subcmds.
func RegisteredSubCmds() map[string]Commander {
	alllock.Lock()
	defer alllock.Unlock()
	return all
}