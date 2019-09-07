package pb

import (
	"fmt"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/log"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/params"
	"os"
	"os/exec"
	"strings"
)

func Protoc(protodirs params.RepeatedOption, protofile, language, outputdir string, pbpkgMapping map[string]string) error {

	// protoc bug:
	//
	// File does not reside within any path specified using --proto_path (or -I).
	// You must specify a --proto_path which encompasses this file.
	// Note that the proto_path must be an exact prefix of the .proto file names
	// -- protoc is too dumb to figure out when two paths (e.g. absolute and relative)
	// are equivalent (it's harder than you think).
	//
	// 当指定如下选项时，protoc仍然无法处理，这其实是个明显的bug，protoc v3.6.0+及以上版本都可以正常执行：
	// protoc --proto_path=/root/test --go_out=paths=source_relative:/root/test/greeter/rpc greeter.proto
	// or
	// protoc --proto_path=. --go_out=paths=source_relative:/root/test/greeter/rpc greeter.proto
	//
	// 下面存在一些"排除--proto_path为当前路径的操作"，纯粹是为了兼容老的protoc处理相对路径、绝对路径的bug
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	args := []string{}
	for _, protodir := range protodirs {
		if protodir == wd {
			continue
		}
		args = append(args, fmt.Sprintf("--proto_path=%s", protodir))
	}

	var out string
	if len(pbpkgMapping) == 0 {
		out = fmt.Sprintf("--%s_out=paths=source_relative:%s", language, outputdir)
	} else {
		pbpkg := ""
		for k, v := range pbpkgMapping {
			pbpkg += "M" + k + "=" + v + ","
		}
		pbpkg = pbpkg[0 : len(pbpkg)-1]
		out = fmt.Sprintf("--%s_out=paths=source_relative,%s:%s", language, pbpkg, outputdir)
	}
	log.Debug("protoc %s", out)

	args = append(args, out, protofile)
	cmd := exec.Command("protoc", args...)

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("Run command: `%s`, error: %s", strings.Join(cmd.Args, " "), string(output))
	}

	return nil
}
