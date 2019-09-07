package params

import (
	"github.com/hitzhangjie/go-rpc/tools/gorpc/config"
)

type Option struct {
	// pb option
	Protodirs RepeatedOption // pb import路径
	Protofile string         // protofile文件
	AliasOn   bool           // 解析MethodOption或者注释中//@alias=别名，用来代替pb文件中rpc

	// template option
	Assetdir string // 服务模板路径
	Language string // 开发语言，如go，java，cpp等
	Protocol string // 协议类型
	HttpOn   bool   // 生成http相关代码
	RpcOnly  bool   // 只生成rpc相关代码，而非完整工程
	// gorpc.json
	GoRPCConfig *config.LanguageCfg

	// gomod option
	GoMod string // 当前工程指定的gomod

	// logging option
	Verbose bool // 输出verbose日志信息
}

func NewOption() *Option {
	return &Option{
		Protodirs: RepeatedOption{},
	}
}
