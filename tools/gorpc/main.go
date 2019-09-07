package main

import (
	"github.com/hitzhangjie/go-rpc/tools/gorpc/cmds"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/log"
	"os"
)

func main() {
	subcmds := cmds.SubCmds()

	if l := len(os.Args); l == 1 {
		subcmds["help"].Run()
		return
	}

	cmd, ok := subcmds[os.Args[1]]
	if !ok || cmd == nil {
		subcmds["help"].Run()
		return
	}

	if err := cmd.Run(os.Args[2:]...); err != nil {
		log.Error("Run command:%v error:\n\t\t%v", os.Args, err)
	}
}
