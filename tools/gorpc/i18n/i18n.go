package i18n

var UsagesCh = map[string]string{
	"protofile": "protobuff 文件名",
	"protocol":  "协议名：nrpc, simplesso, ilive",
	"assetdir":  "模板文件目录",
	"v":         "详细日志",
	"g":         "使用全局GOPATH",
	"httpon":    "开启http模式",
	"gopath":    "是否在$GOPATH/src中搜索profile和其imports文件，即将$GOPATH/src加入protodir参数，如果$GOPATH有多条路径，则都会搜索",
	"protodir":  "指定protofile及其imports文件的搜索目录，可以多次指定，如果未指定，则默认当前目录",
}

var UsagesEn = map[string]string{
	"protofile": "protobuf filename",
	"protocol":  "protocol：nrpc, simplesso, ilive",
	"assetdir":  "template asset dir",
	"v":         "verbose logging info",
	"g":         "using global GOPATH",
	"httpon":    "enable http mode",
	"gopath":    "enable search proto file and its imports in $GOPATH/src.",
	"protodir":  "where should we search for protofile and dependencies, it can be specified multiple times. if not given, the current working directory is used.",
}
