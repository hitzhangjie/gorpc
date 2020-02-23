package config

import "sync"

const (
	ConfTypeYaml = "yaml"
	ConfTypeIni  = "ini"
)

type Loader interface {
	Load(fp string) (Config, error)
}

var (
	loaders    = map[string]Loader{}
	loadersLck = sync.RWMutex{}
)

func Register(typ string, loader Loader) {
	loadersLck.Lock()
	loaders[typ] = loader
	loadersLck.Unlock()
}

func GetLoader(typ string) Loader {
	loadersLck.RLock()
	defer loadersLck.RUnlock()

	v, ok := loaders[typ]
	if !ok {
		return nil
	}
	return v
}
