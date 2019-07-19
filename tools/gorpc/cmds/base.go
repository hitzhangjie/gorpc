package cmds

import (
	"flag"
	"gorpc/tools/gorpc/params"
	"os"
	"path"
	"path/filepath"
)

type Commander interface {
	Run(args ...string) error
	UsageLine() string
	DescShort() string
	DescLong() string
	FlagSet() *flag.FlagSet
}

type Cmd struct {
	usageLine string
	descShort string
	descLong  string
	flagSet   *flag.FlagSet
}

//func (c *Cmd) Run(args ...string) error {
//	c.flagSet.Parse(args)
//	fmt.Println("command:", c)
//	return nil
//}

func (c *Cmd) UsageLine() string {
	return c.usageLine
}

func (c *Cmd) DescShort() string {
	return c.descShort
}

func (c *Cmd) DescLong() string {
	return c.descLong
}

func (c *Cmd) FlagSet() *flag.FlagSet {
	return c.flagSet
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
