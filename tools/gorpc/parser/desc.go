package parser

import (
)

// ServerDescriptor 将pb文件中用于描述服务信息的内容提取出来，包括所属的包名、服务名、支持的协议、启用HTTP、RPC请求等信息，
// 后续结合go template很方便地实现自动化代码生成，基于模板的方式也更容易后期调整、维护、扩展。
type ServerDescriptor struct {
	PackageName string                // pb包名称
	Imports     []string              //跟pb文件中import对应的golang import路径
	ServerName  string                // 服务名称
	Protocol    string                // 协议类型，如gorpc等等
	HttpOn      bool                  // http开关
	RPC         []ServerRPCDescriptor // rpc接口定义
	CreateTime  string                // 服务创建时间
	Author      string                // 作者
}

// ServerRPCDescriptor 将pb文件中用于描述RPC的内容提取出来，包括RPC的名字、命令字、请求类型、响应类型等信息。
type ServerRPCDescriptor struct {
	Name                     string // RPC方法名
	Cmd                      string // RPC命令字
	RequestType              string // RPC请求消息类型，包含package，比如package_a.TypeA
	ResponseType             string // RPC响应消息类型，包含package，比如package_b.TypeB
	RequestTypeNameInRpcTpl  string //在rpc.go.tpl中request type的name
	ResponseTypeNameInRpcTpl string //在rpc.go.tpl中response type的name
}
