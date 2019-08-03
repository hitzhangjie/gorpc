package cmds

import (
	"flag"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/log"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/parser"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/tpl"
	"os"
	"path"
	"time"
)

func init() {
	mux.Lock()
	defer mux.Unlock()
	all["update"] = NewUpdateCmd()
}

func NewUpdateCmd() *UpdateCmd {
	fs := flag.NewFlagSet("updatecmd", flag.ContinueOnError)

	fs.Var(&protodirs, "protodir", "search path for protofile")
	fs.String("protofile", "any.proto", "protofile to handle")
	fs.String("protocol", "gorpc", "protocol to use, gorpc, chick or swan")
	fs.Bool("g", false, "generate code structure conforming to global gopath")
	fs.Bool("v", false, "verbose help info")
	fs.String("assetdir", "", "search path for project template")

	u := Cmd{
		usageLine: `go-rpc update`,
		descShort: `
how to update project:
	go-rpc update -protodir=. -protofile=*.proto -protocol=gorpc
	go-rpc update -protofile=*.proto -protocol=gorpc`,

		descLong: `
go-rpc update:
	-protodir, search path for protofile
	-protofile, protofile to handle
	-protocol, protocol to use, gorpc, chick or swan`,
		flagSet: fs,
	}

	return &UpdateCmd{u}
}

type UpdateCmd struct {
	Cmd
}

func (c *UpdateCmd) Run(args ...string) (err error) {

	c.flagSet.Parse(args)

	protofile = c.flagSet.Lookup("protofile").Value.(flag.Getter).Get().(string)
	protocol = c.flagSet.Lookup("protocol").Value.(flag.Getter).Get().(string)
	verbose = c.flagSet.Lookup("v").Value.(flag.Getter).Get().(bool)
	global = c.flagSet.Lookup("g").Value.(flag.Getter).Get().(bool)

	assetdir = c.flagSet.Lookup("assetdir").Value.(flag.Getter).Get().(string)
	if len(assetdir) == 0 {
		if assetdir, err = defaultAssetDir(); err != nil {
			return err
		}
	}
	log.InitLogging(verbose)
	return c.update()
}

func (c *UpdateCmd) update() error {
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
		log.Error("step 2: Parse proto file:[%s] error:[%v]", err)
		os.Exit(1)
	} else {
		log.Info("step 2: Parse proto file:[%s] succ", protofile)
		//log.Debug("[ServerDescriptor] %#v\n", server_asset)
	}
	server_asset.HttpOn = httpon

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
