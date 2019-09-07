package parser

import (
	"fmt"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/params"
	"os"
	"path"
	"path/filepath"
)

// pb导入路径解析，返回$protodir/$protofile下存在的文件目录列表
func ImportDirs(protodirs *params.RepeatedOption, protofile string) ([]string, error) {

	// `-protodir not specified` or `-protodir=.`
	if len(*protodirs) == 0 || (len(*protodirs) == 1 && (*protodirs)[0] == ".") {
		abs, _ := filepath.Abs(".")


		//protodirs.Set(".")

		return []string{abs}, nil
	}

	// $protodir/$protofile
	dirs := Uniq(*protodirs)
	protodirs.Replace(&dirs)

	//查找protofile的绝对路径
	fpaths := []string{}
	for _, dir := range *protodirs {
		p := path.Join(dir, protofile)

		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			fpaths = append(fpaths, dir)
		}
	}
	if len(fpaths) == 0 {
		return nil, fmt.Errorf("protofile:%s not found in dirs:%v", protofile, protodirs.String())
	} else if len(fpaths) > 1 {
		return nil, fmt.Errorf("protofile:%s found duplicate ones under dirs:%v", protofile, fpaths)
	}
	return fpaths, nil
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
