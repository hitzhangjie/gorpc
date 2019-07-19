package main

import (
	"gorpc/tools/gorpc/i18n"
	"gorpc/tools/gorpc/params"
	"gorpc/tools/gorpc/cmds"
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

var subcmds = cmds.RegisteredSubCmds()

func main() {

	if l := len(os.Args); l == 1 {
		subcmds["help"].Run()
		return
	}

	f, ok := subcmds[os.Args[1]]
	if !ok || f == nil {
		subcmds["help"].Run()
		return
	}

	f.Run(os.Args[2:]...)
}
