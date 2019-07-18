package main

import (
	"git.code.oa.com/go-neat/tools/codegen/cmds"
	"git.code.oa.com/go-neat/tools/codegen/i18n"
	"git.code.oa.com/go-neat/tools/codegen/params"
	"os"
)

var usage map[string]string = i18n.UsagesEn

var (
	protodirs params.StringArray // pb import路径
	protofile *string            // pb 文件
	protocol  *string            // 协议类型，nrpc、simplesso、ilive
	httpon    *bool              // 是否开启http
	assetdir  *string            // 模板路径
	global    *bool              // 生成代码时使用全局GOPATH
	verbose   *bool              // 打印详细日志
)

func init() {
}

var subcmds = map[string]cmds.Commander{
	"help":   cmds.NewHelpCmd(),
	"create": cmds.NewCreateCmd(),
	"update": cmds.NewUpdateCmd(),
	"taf": cmds.NewTafCmd(),
}

func main() {
	if l := len(os.Args); l == 1 {
		subcmds["help"].Run()
		return
	}
	switch os.Args[1] {
	case "help":
		if len(os.Args) > 2 {
			subcmds["help"].Run(os.Args[2:]...)
		} else {
			subcmds["help"].Run()
		}
	case "create", "update", "pull", "push":
		if cmd, ok := subcmds[os.Args[1]]; !ok {
			panic("subcmd invalid")
		} else {
			if len(os.Args) == 2 {
				panic("subcmd options invalid")
			}
			cmd.Run(os.Args[2:]...)
		}
	case "taf":
		if cmd, ok := subcmds[os.Args[1]]; !ok {
			panic("subcmd invalid")
		} else {
			if len(os.Args) == 2 {
				panic("subcmd options invalid")
			}
			cmd.Run(os.Args[2:]...)
		}

	default:
		subcmds["help"].Run()
	}
}

