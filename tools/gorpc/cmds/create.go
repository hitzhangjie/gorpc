package cmds

import (
	"flag"
	"gorpc/tools/gorpc/log"
	"gorpc/tools/gorpc/params"
	"gorpc/tools/gorpc/parser"
	"gorpc/tools/gorpc/spec"
	"gorpc/tools/gorpc/tpl"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"
)

var (
	protodirs params.StringArray // pb import路径
	protofile string
	protocol  string
	httpon    bool
	verbose   bool
	global    bool
	assetdir  string
)

func init() {
	alllock.Lock()
	defer alllock.Unlock()
	all["create"] = NewCreateCmd()
}

func NewCreateCmd() *CreateCmd {
	fs := flag.NewFlagSet("createcmd", flag.ContinueOnError)

	fs.Var(&protodirs, "protodir", "search path for protofile")
	fs.String("protofile", "any.proto", "protofile to handle")
	fs.String("protocol", "nrpc", "protocol to use, nrpc, simplesso or ilive")
	fs.Bool("httpon", false, "enable http mode")
	fs.Bool("g", false, "generate code structure conforming to global gopath")
	fs.Bool("v", false, "verbose help info")
	fs.String("assetdir", "", "search path for project template")

	u := Cmd{
		usageLine: `gorpc create`,
		descShort: `
how to create project:
	gorpc create -protodir=. -protofile=*.proto -protocol=nrpc -httpon=false
	gorpc create -protofile=*.proto -protocol=nrpc`,
		descLong: `
gorpc create:
	-protodir, search path for protofile
	-protofile, protofile to handle
	-protocol, protocol to use, nrpc, simplesso or ilive
	-httpon, enable http mode
	-g, generate code structure conforming to global gopath`,
		flagSet: fs,
	}

	return &CreateCmd{u}
}

type CreateCmd struct {
	Cmd
}

func (c *CreateCmd) Run(args ...string) error {

	c.flagSet.Parse(args)

	protofile = c.flagSet.Lookup("protofile").Value.(flag.Getter).Get().(string)
	protocol = c.flagSet.Lookup("protocol").Value.(flag.Getter).Get().(string)
	httpon = c.flagSet.Lookup("httpon").Value.(flag.Getter).Get().(bool)
	verbose = c.flagSet.Lookup("v").Value.(flag.Getter).Get().(bool)
	global = c.flagSet.Lookup("g").Value.(flag.Getter).Get().(bool)

	assetdir = c.flagSet.Lookup("assetdir").Value.(flag.Getter).Get().(string)
	if len(assetdir) == 0 {
		if dir, err := spec.LocateCfgPath(); err != nil {
			panic(err)
		} else {
			assetdir = filepath.Join(dir, "asset")
		}
	}

	log.InitLogging(verbose)

	return c.create()
}

func (c *CreateCmd) create() error {

	// Q:解析库jhump parser本身支持pb依赖的解析，为什么这里还要再额外的去判断导入路径呢？
	// A:nrpc页面也要支持pb依赖解析，nrpc页面上传的pb文件是按照flat layout进行组织的！
	fpaths := ImportDirs(&protodirs, protofile)

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
	server_asset, err := parser.ParseProtoFile(protofile, protocol, protodirs)
	if err != nil {
		log.Error("step 2: Parse proto file:[%s] error:[%v]", protofile, err)
		os.Exit(1)
	} else {
		log.Info("step 2: Parse proto file:[%s] succ", protofile)
		//log.Debug("[ServerDescriptor] %#v\n", server_asset)
	}
	server_asset.HttpOn = httpon
	server_asset.CreateTime = time.Now().Format("2006-01-02 15:04:05")

	// 代码生成
	fp := path.Join(fpaths[0], protofile)
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
