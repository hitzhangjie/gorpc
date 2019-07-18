package cmds

import (
	"flag"
	"git.code.oa.com/go-neat/tools/codegen/log"
	"git.code.oa.com/go-neat/tools/codegen/params"
	"git.code.oa.com/go-neat/tools/codegen/parser"
	"git.code.oa.com/go-neat/tools/codegen/spec"
	"git.code.oa.com/go-neat/tools/codegen/tpl"
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
		UsageLine: "goneat create",
		Flag:      fs,
	}

	return &CreateCmd{u}
}

type CreateCmd struct {
	Cmd
}

func (c *CreateCmd) Run(args ...string) error {

	c.Flag.Parse(args)

	protofile = c.Flag.Lookup("protofile").Value.(flag.Getter).Get().(string)
	protocol = c.Flag.Lookup("protocol").Value.(flag.Getter).Get().(string)
	httpon = c.Flag.Lookup("httpon").Value.(flag.Getter).Get().(bool)
	verbose = c.Flag.Lookup("v").Value.(flag.Getter).Get().(bool)
	global = c.Flag.Lookup("g").Value.(flag.Getter).Get().(bool)

	assetdir = c.Flag.Lookup("assetdir").Value.(flag.Getter).Get().(string)
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
