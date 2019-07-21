package parser

import (
	"fmt"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/log"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/spec"
	"github.com/jhump/protoreflect/desc"
	"strconv"
	"strings"

	"github.com/jhump/protoreflect/desc/protoparse"
)

var (
	serverDescriptor ServerDescriptor

	// swan
	swanBigCmd = 0
	swanSubCmd = []int{}

	// chick
	chickCmd = []int{}
)

func GetNameWithPackageCheck(fullTypeName string, goPackageName string) string {
	//根据go文件的package来判断是使用全限定的类型名(如package_a.TypeA)，还是直接使用简单类型名(如TypeA)
	eles := strings.Split(fullTypeName, ".")
	if eles != nil && len(eles) > 1 {
		//type所在package名
		typePackageName := strings.Join(eles[:len(eles)-1], ".")
		//type简单名
		typeSimpleName := eles[len(eles)-1]

		if typePackageName == goPackageName {
			//如果type就在当前go文件所在package中，则使用简单类型名
			return typeSimpleName
		}
	}

	return fullTypeName
}

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

	// spec
	serverDescriptor.ProtoSpec = *spec.GetTypeSpec(protocol)

	// serviceCmd: swan bigCmd+subCmd, chick serviceCmd, gorpc don't need this.
	log.Debug("enums: %v", fd.GetEnumTypes())
	for _, e := range fd.GetEnumTypes() {
		name := e.GetName()
		if protocol == "swan" {
			if name == "BIG_CMD" {
				swanBigCmd = int(e.GetValues()[0].GetNumber())
			}
			if name == "SUB_CMD" {
				for _, v := range e.GetValues() {
					swanSubCmd = append(swanSubCmd, int(v.GetNumber()))
				}
			}
		} else if protocol == "chick" {
			if name == "SERVICE_CMD" {
				for _, v := range e.GetValues() {
					chickCmd = append(chickCmd, int(v.GetNumber()))
				}
			}
		}
	}

	// service rpc
	service := fd.GetServices()[0]
	for idx, m := range service.GetMethods() {
		rpc := ServerRPCDescriptor{
			Name:                     m.GetName(),
			RequestType:              m.GetInputType().GetFullyQualifiedName(),
			ResponseType:             m.GetOutputType().GetFullyQualifiedName(),
			RequestTypeNameInRpcTpl:  GetNameWithPackageCheck(m.GetInputType().GetFullyQualifiedName(), serverDescriptor.PackageName),
			ResponseTypeNameInRpcTpl: GetNameWithPackageCheck(m.GetOutputType().GetFullyQualifiedName(), serverDescriptor.PackageName),
		}
		if protocol == "swan" {
			rpc.Cmd = fmt.Sprintf("%#x_%#x", swanBigCmd, swanSubCmd[idx])
		}
		if protocol == "chick" {
			rpc.Cmd = strconv.FormatInt(int64(chickCmd[idx]), 10)
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
