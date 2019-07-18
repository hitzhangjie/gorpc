package cmds

import (
	"flag"
	"fmt"
	"git.code.oa.com/go-neat/tools/codegen/params"
	"os"
	"path"
	"path/filepath"
)

type Commander interface {
	Run(args ...string) error
}

type Cmd struct {
	UsageLine string
	DescShort string
	DescLong  string
	Flag      *flag.FlagSet
}

func (c *Cmd) Run(args ...string) error {
	c.Flag.Parse(args)
	fmt.Println("command:", c)
	return nil
}

func ImportDirs(fileDirs *params.StringArray, fileName string) []string {

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