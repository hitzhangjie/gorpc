package cmds

import (
	"flag"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/log"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/params"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/parser"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/tpl"
	"os"
	"os/exec"
	"path"
)

var (
	protodirs params.List // pb import路径
	protofile string
	protocol  string
	httpon    bool
	verbose   bool
	global    bool
	assetdir  string
)

func init() {
	mux.Lock()
	defer mux.Unlock()
	all["create"] = NewCreateCmd()
}

func NewCreateCmd() *CreateCmd {
	fs := flag.NewFlagSet("createcmd", flag.ContinueOnError)

	fs.Var(&protodirs, "protodir", "search path for protofile")
	fs.String("protofile", "any.proto", "protofile to handle")
	fs.String("protocol", "gorpc", "protocol to use, gorpc, chick or swan")
	fs.Bool("httpon", false, "enable http mode")
	fs.Bool("g", false, "generate code structure conforming to global gopath")
	fs.Bool("v", false, "verbose help info")
	fs.String("assetdir", "", "search path for project template")

	u := Cmd{
		usageLine: `go-rpc create`,
		descShort: `
how to create project:
	go-rpc create -protodir=. -protofile=*.proto -protocol=gorpc -httpon=false
	go-rpc create -protofile=*.proto -protocol=gorpc`,
		descLong: `
go-rpc create:
	-protodir, search path for protofile
	-protofile, protofile to handle
	-protocol, protocol to use, gorpc, chick or swan
	-httpon, enable http mode
	-g, generate code structure conforming to global gopath`,
		flagSet: fs,
	}

	return &CreateCmd{u}
}

type CreateCmd struct {
	Cmd
}

func (c *CreateCmd) Run(args ...string) (err error) {

	c.flagSet.Parse(args)

	protofile = c.flagSet.Lookup("protofile").Value.(flag.Getter).Get().(string)
	protocol = c.flagSet.Lookup("protocol").Value.(flag.Getter).Get().(string)
	httpon = c.flagSet.Lookup("httpon").Value.(flag.Getter).Get().(bool)
	verbose = c.flagSet.Lookup("v").Value.(flag.Getter).Get().(bool)
	global = c.flagSet.Lookup("g").Value.(flag.Getter).Get().(bool)

	assetdir = c.flagSet.Lookup("assetdir").Value.(flag.Getter).Get().(string)
	if len(assetdir) == 0 {
		if assetdir, err = defaultAssetDir(); err != nil {
			return err
		}
	}
	log.InitLogging(verbose)
	return c.create()
}

func (c *CreateCmd) create() error {

	// Q:解析库jhump parser本身支持pb依赖的解析，为什么这里还要再额外的去判断导入路径呢？
	// A:protoc是支持这种操作的，import声明不一定指明导入路径，可通过-I选项指定导入路径.
	fpaths := parser.ImportDirs(&protodirs, protofile)

	if len(fpaths) == 0 {
		log.Error("step 1: proto file:[%s] not found in dirs:%v", protofile, protodirs.String())
		os.Exit(1)
	} else if len(fpaths) > 1 {
		log.Error("step 1: proto file:[%s] found in multiple dirs:[%v], cannot determine which one to use", protofile, fpaths)
		os.Exit(1)
	} else {
		log.Info("step 1: found proto file:[%s] in following dirs:[%v]", protofile, fpaths)
	}

	// 解析pb
	server_asset, err := parser.ParseProtoFile(protofile, protocol, protodirs...)
	if err != nil {
		log.Error("step 2: Parse proto file:[%s] error:[%v]", protofile, err)
		os.Exit(1)
	} else {
		log.Info("step 2: Parse proto file:[%s] succ", protofile)
		//log.Debug("[ServerDescriptor] %#v\n", server_asset)
	}
	server_asset.HttpOn = httpon

	// 代码生成
	fp := path.Join(fpaths[0], protofile)
	log.Info(fp)
	options := map[string]interface{}{
		"protodir":  protodirs,
		"protofile": protofile,
		"assetdir":  assetdir,
		"g":         global,
		"v":         verbose,
	}

	err = tpl.GenerateFiles(server_asset, fp, true, options)

	if err != nil {
		log.Error("step 3: Generate service files error:[%v]", err)
		err = exec.Command("rm", "-r", server_asset.ServerName).Run()
		if err != nil {
			log.Error("remove file fail with error:[%v], you should remove %s by yourself", err, server_asset.ServerName)
		}
		os.Exit(1)
	} else {
		log.Info("step 3: Generate service files success.")
	}

	return nil
}
