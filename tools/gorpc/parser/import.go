package parser

import (
	"github.com/hitzhangjie/go-rpc/tools/gorpc/params"
	"os"
	"path"
	"path/filepath"
)

// pb导入路径解析
func ImportDirs(fileDirs *params.List, fileName string) []string {

	// ./$protofile
	if len(*fileDirs) == 0 {
		abs, _ := filepath.Abs(".")
		fileDirs.Set(".")
		return []string{abs}
	}

	// $protodir/$protofile
	dirs := Uniq(*fileDirs)
	fileDirs.Replace(&dirs)

	//查找protofile的绝对路径
	fpath := []string{}
	for _, dir := range *fileDirs {
		p := path.Join(dir, fileName)

		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			fpath = append(fpath, p)
		}
	}

	return fpath
}

// 字符串去重
func Uniq(dirs []string) []string {

	set := map[string]struct{}{}
	for _, p := range dirs {
		abs, _ := filepath.Abs(p)
		set[abs] = struct{}{}
	}

	uniq := []string{}
	for dir, _ := range set {
		uniq = append(uniq, dir)
	}

	return uniq
}