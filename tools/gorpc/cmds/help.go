package cmds

import (
	"flag"
	"fmt"
)

type HelpCmd struct {
	Cmd
}

func (c *HelpCmd) Run(args ...string) error {

	c.Flag.Parse(args)
	verbose := c.Flag.Lookup("v").Value.(flag.Getter).Get().(bool)

	if verbose {
		fmt.Println(c.DescLong)
	} else {
		fmt.Println(c.DescShort)
	}
	return nil
}

func NewHelpCmd() *HelpCmd {

	fs := flag.NewFlagSet("helpcmd", flag.ContinueOnError)
	fs.Bool("v", false, "verbose help info")

	u := Cmd{
		UsageLine: "goneat help",
		DescShort: `
	how to display help:
		goneat help

	how to create project:
		goneat create -protodir=. -protofile=*.proto -protocol=nrpc -httpon=false
		goneat create -protofile=*.proto -protocol=nrpc

	how to update project:
		goneat update -protodir=. -protofile=*.proto -protocol=nrpc
		goneat update -protofile=*.proto -protocol=nrpc

	how to pull/push rpc:
		goneat pull -protocol=nrpc
		goneat push -protocol=nrpc

	how to create taf server or rpc:
		goneat taf  -jceFile=xxxx.jce -jceDir=.
		goneat taf -cmd=rpc -jceFile=xxxxx.jce -jceDir=. 生成的rpc文件夹请放在goneat项目的src/model目录下
`,
		DescLong: `goneat <cmd> <options>: 

global options:
	-h display this help
	-v display verbose info

goneat create:
	-protodir, search path for protofile
	-protofile, protofile to handle
	-protocol, protocol to use, nrpc, simplesso or ilive
	-httpon, enable http mode
	-g, generate code structure conforming to global gopath

goneat update:
	-protodir, search path for protofile
	-protofile, protofile to handle
	-protocol, protocol to use, nrpc, simplesso or ilive
	
goneat pull:
	-protocol, protocol to use, nrpc, simplesso or ilive

goneat push:
	-protocol, protocol to use, nrpc, simplesso or ilive

goneat taf:
	-cmd, "server" or "rpc", default is "server" rpc命令生成的rpc文件夹请放在goneat项目的src/model目录下
	-jceFile, jceFile to handle
    -jceDir, search path for jceFile

`,
		Flag: fs,
	}

	return &HelpCmd{u}
}
