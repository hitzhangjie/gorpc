package tpl

import (
	"fmt"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/log"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/params"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/parser"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

func GenerateFiles(fd *parser.FileDescriptor, protofilePath, outputdir string, option *params.Option) (err error) {

	serviceIdx := 0
	if len(fd.Services) > 1 {
		// todo 忽略多余的service定义
		log.Info("You have defined more than one service which will be ignored")
	}

	// 准备输出目录
	if err := prepareOutputdir(outputdir); err != nil {
		return fmt.Errorf("GenerateFiles prepareOutputdir:%v", err)
	}

	// 遍历模板文件进行处理
	f := func(path string, info os.FileInfo, err error) error {
		return fileEntryHandler(path, info, err, &mixedOptions{fd, serviceIdx, outputdir, option})
	}

	err = filepath.Walk(option.Assetdir, f)
	if err != nil {
		return fmt.Errorf("GenerateFiles filepath.Walk:%v", err)
	}

	return nil
}

func prepareOutputdir(outputdir string) (err error) {

	_, err = os.Lstat(outputdir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(outputdir, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func GenerateFile(fd *parser.FileDescriptor, infile, outfile string, option *params.Option, rpcIndex ...int) (err error) {

	assetdir := option.Assetdir
	if !path.IsAbs(assetdir) {
		return errors.New("assetdir must be specified an absolute path")
	}

	// stat template
	tplFilePath := infile
	if _, err = os.Stat(tplFilePath); err != nil {
		log.Error("%v", err)
		return err
	}

	// create output file
	fout, err := os.Create(outfile)
	if err != nil {
		log.Error("%v", err)
		return err
	}
	defer fout.Close()

	// template execute and populate the output file
	var tplInstance *template.Template

	baseName := path.Base(infile)
	if funcMap == nil {
		tplInstance, err = template.New(baseName).ParseFiles(tplFilePath)
	} else {
		tplInstance, err = template.New(baseName).Funcs(funcMap).ParseFiles(tplFilePath)
	}
	if err != nil {
		log.Error("%v", err)
		return err
	}

	// 将需要的descriptor信息、命令行控制参数信息、其他分文件需要的rpcindex信息传入
	err = tplInstance.Execute(fout, struct {
		*parser.FileDescriptor
		*params.Option
		RPCIndex int
	}{
		fd,
		option,
		func() int {
			if len(rpcIndex) != 0 {
				return rpcIndex[0]
			}
			return 999999
		}(),
	})
	if err != nil {
		log.Error("%v", err)
		return err
	}

	return nil
}

// fileEntryHandler 处理模板文件
func fileEntryHandler(entry string, info os.FileInfo, err error, options *mixedOptions) error {

	fd := options.FileDescriptor
	option := options.Option
	outputdir := options.OutputDir
	serviceIdx := options.ServiceIdx

	// if incoming error encounter, return at once
	if err != nil {
		return err
	}

	// ignore ${asset_dir}
	var relPath string
	if relPath = strings.TrimPrefix(entry, option.Assetdir); len(relPath) == 0 {
		return nil
	}
	relPath = strings.TrimPrefix(relPath, "/")

	log.Debug("entry srcPath:%s", entry)

	// 如果server stub需要分文件，则指定rpc_server_stub模板文件名
	sd := fd.Services[serviceIdx]
	if relPath == option.GoRPCConfig.RPCServerStub {
		for idx, rpc := range sd.RPC {
			outPath := filepath.Join(outputdir, relPath)
			dir := filepath.Dir(outPath)
			base := sd.Name + "_" + rpc.Name + "." + option.GoRPCConfig.Language
			outPath = filepath.Join(dir, base)
			if err := GenerateFile(fd, entry, outPath, option, idx); err != nil {
				return err
			}
		}
		return nil
	}

	outPath := filepath.Join(outputdir, relPath)
	log.Debug("entry destPath: %s", outPath)

	// if `entry` is directory, create the same entry in `outputdir`
	if info.IsDir() {
		return os.MkdirAll(outPath, os.ModePerm)
	}
	outPath = strings.TrimSuffix(outPath, option.GoRPCConfig.TplFileExt)

	return GenerateFile(fd, entry, outPath, option)
}

// mixedOptions 将众多选项聚集在一起，简化方法签名
type mixedOptions struct {
	*parser.FileDescriptor
	ServiceIdx int
	OutputDir  string
	*params.Option
}
