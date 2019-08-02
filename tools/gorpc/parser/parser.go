package parser

import (
	"github.com/hitzhangjie/go-rpc/tools/gorpc/log"
	"github.com/jhump/protoreflect/desc"
	"strings"

	"github.com/jhump/protoreflect/desc/protoparse"
)

var serverDescriptor ServerDescriptor

func ParseProtoFile(fname, protocol string, protodirs ...string) (*ServerDescriptor, error) {

	parser := protoparse.Parser{
		ImportPaths:           protodirs,
		IncludeSourceCodeInfo: true,
	}

	descriptors, err := parser.ParseFiles(fname)
	if err != nil {
		return nil, err
	}

	fd := descriptors[0]

	serverDescriptor := &ServerDescriptor{}

	// package
	log.Debug("packageName: %s", fd.GetPackage())
	serverDescriptor.PackageName = fd.GetPackage()

	// service
	log.Debug("serviceName: %s", fd.GetServices()[0].GetName())
	serverDescriptor.ServerName = fd.GetServices()[0].GetName()

	// protocol
	serverDescriptor.Protocol = protocol

	// service rpc
	service := fd.GetServices()[0]
	for _, m := range service.GetMethods() {
		rpc := ServerRPCDescriptor{
			Name:                     m.GetName(),
			RequestType:              m.GetInputType().GetFullyQualifiedName(),
			ResponseType:             m.GetOutputType().GetFullyQualifiedName(),
			//RequestTypeNameInRpcTpl:  simplify(m.GetInputType().GetFullyQualifiedName(), serverDescriptor.PackageName),
			//ResponseTypeNameInRpcTpl: simplify(m.GetOutputType().GetFullyQualifiedName(), serverDescriptor.PackageName),
		}
		if protocol == "gorpc" {
			rpc.Cmd = rpc.Name
		}
		serverDescriptor.RPC = append(serverDescriptor.RPC, rpc)
	}

	serverDescriptor.Imports = getGolangImports(fd)
	log.Debug("imports: %s", strings.Join(serverDescriptor.Imports, ","))

	return serverDescriptor, nil
}

func getGolangImports(fd *desc.FileDescriptor) []string {

	//获取跟pb中import对应的golang import
	deps := fd.GetDependencies()

	deps = append(deps, fd.GetPublicDependencies()...)
	deps = append(deps, fd.GetWeakDependencies()...)

	fnameSet := map[string]struct{}{}
	for _, dep := range deps {
		fname := dep.GetFullyQualifiedName()
		parts := strings.Split(fname, "/")
		if len(parts) > 0 {
			fname = strings.Join(parts[:len(parts)-1], "/")
		}
		fnameSet[fname] = struct{}{}
	}

	files := []string{}
	for k, _ := range fnameSet {
		files = append(files, k)
	}
	return files
}
