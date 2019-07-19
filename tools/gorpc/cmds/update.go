package cmds

import (
	"flag"
	"github.com/hitzhangjie/gorpc/tools/gorpc/log"
	"github.com/hitzhangjie/gorpc/tools/gorpc/parser"
	"github.com/hitzhangjie/gorpc/tools/gorpc/spec"
	"github.com/hitzhangjie/gorpc/tools/gorpc/tpl"
	"os"
	"path"
	"path/filepath"
	"time"
)

func init() {
	alllock.Lock()
	defer alllock.Unlock()
	all["update"] = NewUpdateCmd()
}

func NewUpdateCmd() *UpdateCmd {
	fs := flag.NewFlagSet("updatecmd", flag.ContinueOnError)

	fs.Var(&protodirs, "protodir", "search path for protofile")
	fs.String("protofile", "any.proto", "protofile to handle")
	fs.String("protocol", "nrpc", "protocol to use, nrpc, simplesso or ilive")
	fs.Bool("g", false, "generate code structure conforming to global gopath")
	fs.Bool("v", false, "verbose help info")
	fs.String("assetdir", "", "search path for project template")

	u := Cmd{
		usageLine: `gorpc update`,
		descShort: `
how to update project:
	gorpc update -protodir=. -protofile=*.proto -protocol=nrpc
	gorpc update -protofile=*.proto -protocol=nrpc`,

		descLong: `
gorpc update:
	-protodir, search path for protofile
	-protofile, protofile to handle
	-protocol, protocol to use, nrpc, simplesso or ilive`,
		flagSet: fs,
	}

	return &UpdateCmd{u}
}

type UpdateCmd struct {
	Cmd
}

func (c *UpdateCmd) Run(args ...string) error {

	c.flagSet.Parse(args)

	protofile = c.flagSet.Lookup("protofile").Value.(flag.Getter).Get().(string)
	protocol = c.flagSet.Lookup("protocol").Value.(flag.Getter).Get().(string)
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

	return c.update()
}

func (c *UpdateCmd) update() error {
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
		log.Error("step 2: Parse proto file:[%s] error:[%v]", err)
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

	err = tpl.GenerateFiles(server_asset, fp, false, options)

	if err != nil {
		log.Error("step 3: Generate service files error:[%v]", err)
		os.Exit(1)
	} else {
		log.Info("step 3: Generate service files success.")
	}

	return nil
}
