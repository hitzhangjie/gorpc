package main

import (
	"github.com/hitzhangjie/go-rpc/tools/gorpc/cmds"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/params"
	"os"
)

var (
	protodirs params.List // pb import路径
	protofile *string     // pb 文件
	protocol  *string     // 协议类型，如gorpc等等
	httpon    *bool       // 是否开启http
	assetdir  *string     // 模板路径
	global    *bool       // 生成代码时使用全局GOPATH
	verbose   *bool       // 打印详细日志
)

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
