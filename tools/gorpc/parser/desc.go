package parser

// FileDescriptor 文件作用域相关的描述信息
type FileDescriptor struct {
	PackageName    string                 // pb包名称
	Imports        []string               // 跟pb文件中import对应的golang import路径
	FileOptions    map[string]interface{} // fileoptions
	Services       []*ServiceDescriptor   // 支持多service
	Dependencies   map[string]string      // 依赖pb文件对应的输出包名
	pkgPkgMappings map[string]string      // pkg到pkg的映射关系
}

// ServiceDescriptor service作用域相关的描述信息
type ServiceDescriptor struct {
	Name string           // 服务名称
	RPC  []*RPCDescriptor // rpc接口定义
}

// RPCDescriptor rpc作用域相关的描述信息
//
// RequestType由于涉及到
type RPCDescriptor struct {
	Name              string // RPC方法名
	Cmd               string // RPC命令字
	FullyQualifiedCmd string // 完整的RPC命令字，用于ServiceDesc、client请求时命令字
	RequestType       string // RPC请求消息类型，包含package，比如package_a.TypeA
	ResponseType      string // RPC响应消息类型，包含package，比如package_b.TypeB
	LeadingComments   string // RPC前置注释信息
	TrailingComments  string // RPC后置注释信息
}
