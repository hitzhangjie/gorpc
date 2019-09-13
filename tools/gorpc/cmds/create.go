package cmds

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/config"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/params"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/parser"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/parser/gomod"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/tpl"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/util/fs"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/util/log"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/util/pb"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// CreateCmd 创建服务工程，或者创建rpcstub (-rpconly)
//
// 1. create project:
// 		gorpc create -protodir=<dir1> -protodir=<dir2> -protofile=greeter.proto -protocol=gorpc -v
//		or
//		gorpc create -protofile=greeter.proto
// 2. create rpcstub:
// 		gorpc create -protofile=greeter.proto -rpconly
type CreateCmd struct {
	Cmd
	*params.Option
}

// newCreateCmd 创建一个CreateCmd
func newCreateCmd() *CreateCmd {

	cmd := Cmd{
		usageLine: `gorpc create`,
		descShort: `
	how to create project:
		gorpc create -protodir=. -protofile=*.proto -protocol=gorpc -alias
		gorpc create -protofile=*.proto -protocol=gorpc`,
		descLong: `
	gorpc create:
		-protodir, search path for protofile, default: "."
		-protofile, protofile to handle
		-protocol, protocol to use, including: gorpc, nrpc, ilive, sso, default: gorpc 
		-lang, language including: go, java, cpp, default: go
		-alias, enable alias mode, //@alias=${rpcName}, default: false
		-rpconly, generate rpc stub only, default: false"`,
		flagSet: newCreateFlagSet(),
	}

	return &CreateCmd{cmd, params.NewOption()}
}

// newCreateFlagSet 为CreateCmd创建专有的参数
func newCreateFlagSet() *flag.FlagSet {

	fs := flag.NewFlagSet("createcmd", flag.ContinueOnError)

	fs.Var(&params.RepeatedOption{}, "protodir", "search path of protofile")
	fs.String("protofile", "any.proto", "protofile to handle")
	fs.String("protocol", "gorpc", "protocol to use, gorpc, chick or swan")
	fs.Bool("v", false, "verbose mode")
	fs.String("assetdir", "", "search path of project template")
	fs.Bool("alias", false, "rpcname alias mode")
	fs.Bool("rpconly", false, "generate rpc stub only")
	fs.String("lang", "go", "language, including go, java, cpp, etc")

	return fs
}

// Run 执行CreateCmd创建逻辑
func (c *CreateCmd) Run(args ...string) (err error) {

	c.flagSet.Parse(args)

	var protofile string

	params.LookupFlag(c.flagSet, "protodir", &c.Protodirs)
	params.LookupFlag(c.flagSet, "protofile", &protofile)
	params.LookupFlag(c.flagSet, "lang", &c.Language)
	params.LookupFlag(c.flagSet, "protocol", &c.Protocol)
	params.LookupFlag(c.flagSet, "alias", &c.AliasOn)
	params.LookupFlag(c.flagSet, "rpconly", &c.RpcOnly)
	params.LookupFlag(c.flagSet, "assetdir", &c.Assetdir)
	params.LookupFlag(c.flagSet, "v", &c.Verbose)

	// `-protofile=abc/d.proto`, works like `-protodir=abc -protofile=d.proto`ma
	p, err := filepath.Abs(protofile)
	if err != nil {
		panic(err)
	}
	c.Protofile = filepath.Base(p)
	c.Protodirs = append(c.Protodirs, filepath.Dir(p))

	// load language config in gorpc.json
	c.GoRPCConfig, err = config.GetLanguageCfg(c.Language)
	if err != nil {
		return err
	}

	// using assetdir in gorpc.json
	if len(c.Assetdir) == 0 {
		c.Assetdir = c.GoRPCConfig.AssetDir
	}

	// init logging level
	log.InitLogging(c.Verbose)

	if !c.RpcOnly {
		return c.create()
	}
	return c.generateRPCStub()
}

func (c *CreateCmd) create() error {

	// 检查pb中的导入路径
	fpaths, err := parser.ImportDirs(&c.Protodirs, c.Protofile)
	if err != nil {
		return err
	}
	log.Info("Found protofile:%s in following dir:%v", c.Protofile, fpaths)

	// 解析pb
	fd, err := parser.ParseProtoFile(c.Option)
	if err != nil {
		return fmt.Errorf("Parse protofile:%s error:%v", c.Protofile, err)
	}

	// 解析gomod
	mod, err := gomod.LoadGoMod()
	if err == nil && len(mod) != 0 {
		c.GoMod = mod
	}
	dump(fd)

	// 代码生成
	// - 准备输出目录
	outputdir, err := getOutputDir(fd, c.Option)
	if err != nil {
		return err
	}
	// - 生成代码
	protofileAbsPath := path.Join(fpaths[0], c.Protofile)

	err = tpl.GenerateFiles(fd, protofileAbsPath, outputdir, c.Option)
	if err != nil {
		os.RemoveAll(outputdir)
		return err
	}
	// - generate *.pb.go or *.java or *.pb.h/*.pb.cc under outputdir/rpc/
	pbOutDir := path.Join(outputdir, "rpc")
	if err = pb.Protoc(c.Option.Protodirs, c.Option.Protofile, c.Option.Language, pbOutDir, fd.Dependencies); err != nil {
		return fmt.Errorf("GenerateFiles: %v", err)
	}
	// - copy *.proto to outpoutdir/rpc/
	basename := path.Base(protofileAbsPath)
	if err := fs.Copy(protofileAbsPath, path.Join(pbOutDir, basename)); err != nil {
		return err
	}
	// - move outputdir/rpc to outputdir/dir($gopkgdir)
	src := path.Join(outputdir, "rpc")
	fileOption := fmt.Sprintf("%s_package", c.Option.GoRPCConfig.Language)
	gopkgdir := fd.PackageName
	if fo := fd.FileOptions[fileOption]; fo != nil {
		if v := fd.FileOptions[fileOption].(string); len(v) != 0 {
			gopkgdir = v
		}
	}

	// - 将outputdir/rpc移动到outputdir/$gopkgdir/
	dest := path.Join(outputdir, gopkgdir)
	if err := os.MkdirAll(path.Dir(dest), os.ModePerm); err != nil {
		return err
	}
	if err := fs.Move(src, dest); err != nil {
		return err
	}
	// - 将stub文件gorpc.go重命名，
	// fixme handle .gorpc.go
	sd := fd.Services[0]
	err = filepath.Walk(dest, func(fpath string, info os.FileInfo, err error) error {
		if strings.HasSuffix(dest, ".go") {
			err := gofmt(dest)
			if err != nil {
				log.Error("Warn: gofmt file:%s error:%v", dest, err)
			}
		}

		if fname := path.Base(fpath); fname == "gorpc.go" {
			fs.Move(fpath, path.Join(path.Dir(fpath), sd.Name+".gorpc.go"))
		}
		return nil
	})
	if err != nil {
		return err
	}

	log.Info("Generate project %s```%s```%s success", log.COLOR_RED, sd.Name, log.COLOR_GREEN)
	return nil
}

func (c *CreateCmd) generateRPCStub() error {

	// 检查pb中的导入路径
	fpaths, err := parser.ImportDirs(&c.Protodirs, c.Protofile)
	if err != nil {
		return err
	}
	log.Info("Found protofile:%s in following dir:%v", c.Protofile, fpaths)

	// 解析pb
	fd, err := parser.ParseProtoFile(c.Option)
	if err != nil {
		return fmt.Errorf("Parse protofile:%s error:%v", c.Protofile, err)
	}

	// 解析gomod
	mod, err := gomod.LoadGoMod()
	if err == nil && len(mod) != 0 {
		c.GoMod = mod
	}
	dump(fd)

	// 代码生成
	// - 准备输出目录
	outputdir, err := os.Getwd()
	if err != nil {
		return err
	}
	// - 生成代码，只处理clientstub
	for _, f := range c.Option.GoRPCConfig.RPCClientStub {
		in := path.Join(c.Assetdir, f)
		log.Debug("handle:%s", in)
		out := path.Join(outputdir, strings.TrimSuffix(path.Base(in), c.GoRPCConfig.TplFileExt))
		if err := tpl.GenerateFile(fd, in, out, c.Option); err != nil {
			return err
		}
	}
	// 将stub文件gorpc.go重命名
	// fixme, handle .gorpc.go
	sd := fd.Services[0]
	err = filepath.Walk(outputdir, func(fpath string, info os.FileInfo, err error) error {
		if fname := path.Base(fpath); fname == "gorpc.go" {
			fs.Move(fpath, path.Join(path.Dir(fpath), sd.Name+".gorpc.go"))
		}
		return nil
	})
	if err != nil {
		return err
	}
	// - generate *.pb.go or *.java or *.pb.h/*.pb.cc under outputdir/rpc/
	if err = pb.Protoc(c.Option.Protodirs, c.Option.Protofile, c.Option.Language, outputdir, fd.Dependencies); err != nil {
		return fmt.Errorf("GenerateFiles: %v", err)
	}
	log.Info("Generate rpc stub success")
	return nil
}

func dump(fd *parser.FileDescriptor) {
	log.Debug("************************** FileDescriptor ***********************")
	buf, _ := json.MarshalIndent(fd, "", "  ")
	log.Debug("\n%s", string(buf))
	log.Debug("*****************************************************************")
}

// purgeNonRpcStub 清理非rpcstub文件
//
// todo 不同语言将${lang}_package转换成文件系统路径时，处理方式不同，先硬编码几个
func purgeNonRpcStub(fd *parser.FileDescriptor, outputdir string, option *params.Option) error {

	// 先简单粗暴搞一发
	pkg, err := pkgFileOption(fd, option.Language)
	if err != nil {
		return err
	}

	// copy outputdir/rpc/ to tmpdir/$pkg/
	src := path.Join(outputdir, pkg)
	tmpdir := path.Join(os.TempDir(), fd.PackageName)
	if err := fs.Copy(src, tmpdir); err != nil {
		return err
	}

	// delete any file or directory under outputdir
	err = fs.DeleteFilesUnderDir(outputdir)
	if err != nil {
		return err
	}

	// copy any file under tmpdir/$pkg/ into outputdir
	err = fs.CopyFileUnderDir(tmpdir, outputdir)
	if err != nil {
		return err
	}

	// remove tmpdir
	if err := os.RemoveAll(tmpdir); err != nil {
		return err
	}

	return nil
}

// pkgFileOption 获取pb文件中的FileOption:${lang}_package
func pkgFileOption(fd *parser.FileDescriptor, lang string) (string, error) {
	pkg := fd.PackageName
	if strings.ToLower(lang) != "cpp" {
		fo := fmt.Sprintf("%s_package", lang)
		opt, ok := fd.FileOptions[fo]
		if !ok || opt == nil {
			return "", fmt.Errorf("invalid FileOption:%s", fo)
		}
		if v, ok := opt.(string); !ok || len(pkg) == 0 {
			return "", fmt.Errorf("invalid FileOption:%s", fo)
		} else {
			pkg = v
		}
	}
	return pkg, nil
}
