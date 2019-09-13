package cmds

import (
	"flag"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/log"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/params"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/parser"
	"go/format"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

// Commander defines the subcmd behavior
type Commander interface {
	// UsageLine cmd example
	UsageLine() string

	// DescShort cmd brief description
	DescShort() string

	// DescLong cmd detailed description
	DescLong() string

	// FlagSet cmd flagset
	FlagSet() *flag.FlagSet

	// Run cmd run
	Run(args ...string) error
}

// Cmd defines the subcmd base behavior
type Cmd struct {
	usageLine string
	descShort string
	descLong  string
	flagSet   *flag.FlagSet
}

// UsageLine returns usage line
func (c *Cmd) UsageLine() string {
	return c.usageLine
}

// DescShort returns the short description
func (c *Cmd) DescShort() string {
	return c.descShort
}

// DescLong returns the long description
func (c *Cmd) DescLong() string {
	return c.descLong
}

// FlagSet returns the flagset
func (c *Cmd) FlagSet() *flag.FlagSet {
	return c.flagSet
}

func defaultAssetDir() (dir string, err error) {
	u, err := user.Current()
	if err != nil {
		return
	}
	if u.Username != "root" {
		dir = filepath.Join(u.HomeDir, ".gorpc/asset")
	} else {
		dir = "/etc/gorpc/assetdir"
	}
	return
}

func getOutputDir(fd *parser.FileDescriptor, options *params.Option) (string, error) {

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if len(options.GoMod) != 0 {
		return wd, nil
	}

	// 准备输出目录
	//pkgName := fd.PackageName
	//switch options.Language {
	//case "go", "java":
	//	fo := fmt.Sprintf("%s_package", options.Language)
	//	if v, ok := fd.FileOptions[fo]; ok && len(v.(string)) != 0 {
	//		pkgName = v.(string)
	//	}
	//}
	//return path.Join(wd, pkgName), nil

	return path.Join(wd, fd.Services[0].Name), nil
}

func gofmt(fpath string) error {
	in, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}

	out, err := format.Source(in)
	if err != nil {
		log.Error("%v", err)
		return err
	}

	err = ioutil.WriteFile(fpath, out, 0644)
	if err != nil {
		return err
	}

	return nil
}
