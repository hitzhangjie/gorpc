package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/log"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/params"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/parser/gorpc"
	"github.com/jhump/protoreflect/desc"
	"strings"

	"github.com/jhump/protoreflect/desc/protoparse"
)

// parseProtoFile 调用jhump/protoreflect来解析pb文件，拿到proto文件描述信息
func parseProtoFile(fname string, protodirs ...string) ([]*desc.FileDescriptor, error) {

	parser := protoparse.Parser{
		ImportPaths:           protodirs,
		IncludeSourceCodeInfo: true,
	}

	return parser.ParseFiles(fname)
}

// checkRequirements 检查是否符合某些约束条件
//
// requirements:
// - 必须指定`fileoption go_package`，且go_package结尾部分必须与package diretive指定的包名一致;
// - 必须保证`packageName` === `serviceName`，如果有定义多个service只处理第一个service，其余忽略;
// - service定义数量不能为0;
func checkRequirements(fd *desc.FileDescriptor) error {

	// fixme MUST: syntax = "proto3"
	//if !fd.IsProto3() {
	//	return errors.New("syntax isn't proto3")
	//}

	// fixme MUST: option go_package = "git.code.oa.com/$group/$repo"
	// fixme MUST: option go_package trailing part, must equal to package directive
	//opts := fd.GetFileOptions()
	//if opts == nil {
	//	return errors.New(`FileOption 'go_package' missing`)
	//}
	//
	//gopkg := opts.GetGoPackage()
	//if len(gopkg) == 0 {
	//	return errors.New(`FileOption 'go_package' missing`)
	//} else {
	//	var trailing string
	//	idx := strings.LastIndex(gopkg, "/")
	//
	//	if idx < 0 {
	//		trailing = gopkg
	//	} else {
	//		trailing = gopkg[idx+1:]
	//	}
	//
	//	if trailing != fd.GetPackage() {
	//		return errors.New(`'option go_package="a/b/c"' trailing part "c" must be consistent with 'package diretive'`)
	//	}
	//}

	// MUST: service
	if len(fd.GetServices()) == 0 {
		return errors.New("service missing")
	}

	// must: packagename === services[0].name
	//if fd.GetPackage() != fd.GetServices()[0].GetName() {
	//	return errors.New(`'packageName' must be consistent with first 'serviceName'`)
	//}

	return nil
}

// ParseProtoFile 解析proto文件，返回一个构造好的可以应用于模板填充的FileDescriptor对象
//
// ParseProtoFile负责的工作包括：
// - 解析pb文件，拿到原始的描述信息
// - 检查工程约束，如是否制定了go_option、method option等自定义的一些业务开发约束
func ParseProtoFile(option *params.Option) (*FileDescriptor, error) {

	// 解析pb
	var fd *desc.FileDescriptor
	if fds, err := parseProtoFile(option.Protofile, option.Protodirs...); err != nil {
		return nil, err
	} else {
		fd = fds[0]
	}
	// 检查约束
	if err := checkRequirements(fd); err != nil {
		return nil, err
	}

	// 构造可以用于指导代码生成的FileDescriptor
	fileDescriptor := new(FileDescriptor)
	// 设置依赖(import的pb文件及其输出包名)
	fillDependencies(fd, fileDescriptor)
	// - 设置packageName
	withErrorCheck(fillPackageName(fd, fileDescriptor))
	// - 设置imports
	withErrorCheck(fillImports(fd, fileDescriptor))
	// - 设置fileOptions
	withErrorCheck(fillFileOptions(fd, fileDescriptor))
	// - 设置service
	withErrorCheck(fillServices(fd, fileDescriptor, option.AliasOn))

	return fileDescriptor, nil
}

func withErrorCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func fillDependencies(fd *desc.FileDescriptor, nfd *FileDescriptor) error {
	pbPkgMappings := map[string]string{}     // pb文件到protoc处理后package名的映射关系
	pkgPkgMappings := map[string]string{}    // pb文件package directive与protoc处理后package名的映射关系
	pkgImportMappings := map[string]string{} // pb文件package directive与导入路径的

	for _, dep := range fd.GetDependencies() {
		fname := dep.GetFullyQualifiedName()
		pkgname := dep.GetPackage()
		pkgImportMappings[pkgname] = pkgname
		validPkgName := pkgname
		if opts := dep.GetFileOptions(); opts != nil {
			if gopkgopt := opts.GetGoPackage(); len(gopkgopt) != 0 {
				pkgImportMappings[pkgname] = gopkgopt
				//idx := strings.LastIndex(gopkgopt, ".")
				idx := strings.LastIndex(gopkgopt, "/")
				if len(gopkgopt[idx+1:]) > 0 {
					validPkgName = gopkgopt[idx+1:]
				}
			}
		}
		pkgPkgMappings[pkgname] = validPkgName
		pbPkgMappings[fname] = validPkgName
	}
	nfd.Dependencies = pbPkgMappings
	nfd.pkgPkgMappings = pkgPkgMappings
	nfd.PkgImportMappings = pkgImportMappings

	return nil
}

func fillPackageName(fd *desc.FileDescriptor, nfd *FileDescriptor) error {
	nfd.PackageName = fd.GetPackage()
	return nil
}

func fillImports(fd *desc.FileDescriptor, nfd *FileDescriptor) error {
	nfd.Imports = getImports(fd, nfd)
	return nil
}

func fillFileOptions(fd *desc.FileDescriptor, nfd *FileDescriptor) error {

	opts := fd.GetFileOptions()
	if opts == nil {
		return nil
	}

	v, err := json.Marshal(opts)
	if err != nil {
		return err
	}
	m := make(map[string]interface{})
	if err := json.Unmarshal(v, &m); err != nil {
		return err
	}

	if nfd.FileOptions == nil {
		nfd.FileOptions = make(map[string]interface{})
	}

	for k, v := range m {
		nfd.FileOptions[k] = v
	}
	return nil
}

func fillServices(fd *desc.FileDescriptor, nfd *FileDescriptor, aliasMode bool) error {

	for _, sd := range fd.GetServices() {

		nsd := new(ServiceDescriptor)
		nfd.Services = append(nfd.Services, nsd)

		// service name
		nsd.Name = sd.GetName()

		// service methods
		for _, m := range sd.GetMethods() {

			rpc := &RPCDescriptor{
				Name: m.GetName(),
				Cmd:  m.GetName(),
				// fixme 这里写死了rpc的拼接规则为/$package.$service/$method
				FullyQualifiedCmd: fmt.Sprintf("/%s.%s/%s", fd.GetPackage(), sd.GetName(), m.GetName()),
				RequestType:       m.GetInputType().GetFullyQualifiedName(),
				ResponseType:      m.GetOutputType().GetFullyQualifiedName(),
				LeadingComments:   m.GetSourceInfo().GetLeadingComments(),
				TrailingComments:  m.GetSourceInfo().GetTrailingComments(),
			}
			nsd.RPC = append(nsd.RPC, rpc)

			// check method option, if gorpc.alias exists, use it as rpc.Cmd
			hasMethodOptions := false
			if v, err := proto.GetExtension(m.GetMethodOptions(), gorpc.E_Alias); err == nil {
				s := v.(*string)
				if s == nil {
					log.Debug("method:%s.%s parse methodOptions option gorpc.alias not specified", sd.GetName(), m.GetName())
				} else {
					log.Debug("method:%s.%s parse methodOptions, name:%s = %s ", sd.GetName(), m.GetName(), gorpc.E_Alias, *(v.(*string)))
					if s != nil {
						if cmd := strings.TrimSpace(*s); len(cmd) != 0 {
							rpc.FullyQualifiedCmd = cmd
							hasMethodOptions = true
						}
					}
				}
			}

			if !hasMethodOptions && aliasMode {
				// check comment //@alias=${rpcName}
				annotation := "@alias="
				hasLeadingAlias := strings.Contains(rpc.LeadingComments, annotation)
				hasTrailingAlias := strings.Contains(rpc.TrailingComments, annotation)

				if hasLeadingAlias && hasTrailingAlias {
					return fmt.Errorf("service:%s, method:%s, leading and trailing aliases conflict", sd.GetName(), m.GetName())
				}

				if hasLeadingAlias {
					s := strings.Split(rpc.LeadingComments, annotation)
					if len(s) != 2 {
						panic(fmt.Sprintf("invalid alias annotation:%s", rpc.LeadingComments))
					}
					cmd := s[1]
					if len(cmd) == 0 {
						panic(fmt.Sprintf("invalid alias annotation:%s", rpc.LeadingComments))
					}
					rpc.FullyQualifiedCmd = cmd
				}

				if hasTrailingAlias {
					s := strings.Split(rpc.TrailingComments, annotation)
					if len(s) != 2 {
						panic(fmt.Sprintf("invalid alias annotation:%s", rpc.TrailingComments))
					}
					cmd := s[1]
					if len(cmd) == 0 {
						panic(fmt.Sprintf("invalid alias annotation:%s", rpc.TrailingComments))
					}
					rpc.FullyQualifiedCmd = cmd
				}
			}
		}
	}

	return nil
}

func getImports(fd *desc.FileDescriptor, nfd *FileDescriptor) []string {

	pkgs := []string{}

	// 遍历rpc，检查是否有req\rsp出现在对应的pkg中，是则允许添加到pkgs，否则从中剔除
	m := map[string]struct{}{}
	for _, rpc := range fd.GetServices()[0].GetMethods() {
		p1 := TrimRight(".", rpc.GetInputType().GetFullyQualifiedName())
		p2 := TrimRight(".", rpc.GetOutputType().GetFullyQualifiedName())
		m[p1] = struct{}{}
		m[p2] = struct{}{}
	}
	for k, _ := range m {
		//if v, ok := nfd.pkgPkgMappings[k]; ok && len(v) != 0 {
		//	pkgs = append(pkgs, v)
		//}
		if v, ok := nfd.PkgImportMappings[k]; ok && len(v) != 0 {
			pkgs = append(pkgs, v)
		}
	}

	return pkgs
}
