package spec

import (
	"os"
	"path/filepath"
)

func LocateCfgPath() (string, error) {

	// 先检查~/.goneat是否存在
	h := os.Getenv("HOME")
	p := filepath.Join(h, ".goneat")
	_, err := os.Lstat(p)
	if err == nil {
		return p, nil
	}

	// 不存在则继续检查/etc/goneat
	if !os.IsNotExist(err) {
		return "", err
	}

	p = "/etc/goneat"
	_, err = os.Lstat(p)
	if err == nil {
		return p, nil
	}
	return "", err
}
