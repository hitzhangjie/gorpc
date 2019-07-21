package tpl

import (
	"fmt"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/fs"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/log"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/params"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/parser"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/spec"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"text/template"
)

func GenerateFiles(asset *parser.ServerDescriptor, fAbsPath string, create bool, options map[string]interface{}) error {

	var (
		outputdir string // 输出目录
		err       error
	)

	defer func() {
		if !create {
			os.RemoveAll(outputdir)
		}
	}()

	// 准备输出目录
	if create {
		outputdir, err = getOutputdir(asset)
		failfast(err)
	} else {
		outputdir = path.Join(os.TempDir(), asset.ServerName)
	}

	err = PrepareOutputDir(outputdir)
	failfast(err)

	// 处理模板文件
	tpllist := map[string]string{}

	// 新建模式go-rpc create需要生成服务代码，gorpc update只更新rpc就可以了
	if create {
		tpllist = filelist(asset, options)
		for fin, fout := range tpllist {
			generateFile(asset, fin, fout, nil, outputdir, options)
		}

		// global GOPATH
		if global, ok := options["g"].(bool); ok && global {
			fs.Move(path.Join(outputdir, "src", asset.ServerName+".go"), outputdir)
			fs.Move(path.Join(outputdir, "src/exec"), outputdir)
			os.Remove(path.Join(outputdir, "src"))
		}
	}

	// clientRpcStub + copy proto + generate *.pb.go
	funcMap := template.FuncMap{"splitIliveCmd": splitIliveCmd}
	generateFile(asset, "rpc/rpc.go.tpl", "rpc/"+asset.ServerName+"_rpc.go", funcMap, outputdir, options)

	// move rpcStub/pb/pb.go to /data/home/go-rpc/src/github.com/hitzhangjie/go-rpc-protos/${server}
	protofile := options["protofile"].(string)
	protodirs := options["protodir"].(params.List)

	// - copy pb to /rpc + /proto
	src := fAbsPath
	dest := path.Join(outputdir, "proto", protofile)
	fs.Copy(src, dest)
	dest = path.Join(outputdir, "rpc", protofile)
	fs.Copy(src, dest)

	// generate *.pb.go in /rpc
	err = runProtocGoOut(protodirs, protofile, outputdir)
	if err != nil {
		return err
	}

	// move outputdir/rpc to public/servername
	src = path.Join(outputdir, "rpc")
	dest = path.Join(spec.GetTypeSpec(asset.Protocol).LocalPrefix, asset.ServerName)
	if err = os.RemoveAll(dest); err != nil && os.IsNotExist(err) {
		log.Error("remove file error:%v, file:%s", err, dest)
		return err
	}

	// cannot handle invalid cross-device link, try copy and delete, or use `mv` instead.
	/*
		if err = fs.Move(src, dest); err != nil {
			log.Error("move file error:%v, src:%s to dest:%s", err, src, dest)
			return err
		}
	*/
	if err = exec.Command("mv", src, dest).Run(); err != nil {
		log.Error("move file error:%v, src:%s to dest:%s", err, src, dest)
		return err
	}
	log.Debug("move file success, src:%s to dest:%s", src, dest)

	// 生成log目录
	os.Mkdir(path.Join(outputdir, "log"), os.ModePerm)

	return nil
}

func PrepareOutputDir(outputdir string) error {

	_, err := os.Stat(outputdir)
	if err == nil {
		return fmt.Errorf("Stat Output directory existed: %s", outputdir)
	} else {
		if os.IsNotExist(err) {
			os.MkdirAll(outputdir, os.ModePerm)
			log.Debug("Create Output directory ok: %s", outputdir)
		} else {
			return err
		}
	}

	dirs := []string{"/src/", "/src/exec/", "/rpc/", "/client/", "/conf/", "/bin/", "/DEVOPS/"}
	for _, dir := range dirs {
		abs := path.Join(outputdir, dir)
		err := os.MkdirAll(abs, os.ModePerm)
		if err != nil {
			log.Debug("Create directory error: %v", err)
			return fmt.Errorf("Create directory error: %v", err)
		}
	}

	return nil
}

func generateFile(asset *parser.ServerDescriptor, infile, outfile string, funcMap template.FuncMap, outputdir string, options map[string]interface{}) (err error) {

	defer func() {
		if err != nil {
			log.Error("generate file:[%s] error:[%v]", outfile, err)
		} else {
			log.Debug("generate file:[%s] succ", outfile)
		}
	}()

	assetdir := options["assetdir"].(string)
	if !path.IsAbs(assetdir) {
		return errors.New("assetdir must be specified an absolute path")
	}

	// stat template
	tplFilePath := path.Join(assetdir, infile)
	_, err = os.Stat(tplFilePath)
	failfast(err)

	// create output file
	dest := path.Join(outputdir, outfile)
	fout, err := os.Create(dest)
	failfast(err)
	defer fout.Close()

	// template execute and populate the output file
	var tplInstance *template.Template

	baseName := infile[strings.LastIndex(infile, "/")+1:]
	if funcMap == nil {
		tplInstance, err = template.New(baseName).ParseFiles(tplFilePath)
	} else {
		tplInstance, err = template.New(baseName).Funcs(funcMap).ParseFiles(tplFilePath)
	}

	failfast(err)

	err = tplInstance.Execute(fout, *asset)
	failfast(err)

	return nil
}

func splitIliveCmd(cmdStr string) string {
	cmds := strings.Split(cmdStr, "_")
	bigCmd, _ := strconv.ParseInt(cmds[0], 0, 32)
	subCmd, _ := strconv.ParseInt(cmds[1], 0, 32)

	return strconv.Itoa(int(bigCmd)) + ", " + strconv.Itoa(int(subCmd))
}

func getOutputdir(asset *parser.ServerDescriptor) (string, error) {

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	log.Debug("Current working directory: %s", wd)

	// 准备输出目录
	return path.Join(wd, asset.ServerName), nil
}

func runProtocGoOut(protodirs params.List, protofile, outputdir string) error {

	args := make([]string, 0)

	//wd, _ := os.Getwd()

	for _, protodir := range protodirs {
		arg_proto_path := fmt.Sprintf("--proto_path=%s", protodir)
		args = append(args, arg_proto_path)
	}
	arg_go_out := fmt.Sprintf("--go_out=%s", path.Join(outputdir, "rpc"))
	arg_proto_file := protofile

	args = append(args, arg_go_out)
	args = append(args, arg_proto_file)

	log.Debug("run: protoc %s", strings.Join(args, " "))
	cmd := exec.Command("protoc", args...)
	output, err := cmd.CombinedOutput()

	log.Debug(string(output))
	if err != nil {
		return fmt.Errorf("Run error: %v, errmsg: %s", err, string(output))
	}

	return nil
}

func failfast(err error) {
	if err != nil {
		log.Error("Error: %v", err)
		os.Exit(1)
	}
}

func filelist(asset *parser.ServerDescriptor, options map[string]interface{}) map[string]string {

	filelist := map[string]string{
		"src/svr_main.go.tpl":          path.Join("src", asset.ServerName+".go"),
		"src/exec/exec.go.tpl":         path.Join("src/exec", "exec_"+asset.ServerName+".go"),
		"src/exec/exec_impl.go.tpl":    path.Join("src/exec", "exec_"+asset.ServerName+"_impl.go"),
		"src/exec/exec_init.go.tpl":    path.Join("src/exec", "exec_"+asset.ServerName+"_init.go"),
		"client/client.go.tpl":         path.Join("client", asset.ServerName+"_client.go"),
		"conf/service.ini.tpl":         path.Join("conf", "service.ini"),
		"conf/metric.ini.tpl":         path.Join("conf", "metric.ini"),
		"conf/trace.ini.tpl":           path.Join("conf", "trace.ini"),
		"conf/log.ini.tpl":             path.Join("conf", "log.ini"),
		"bin/run.sh.tpl":               path.Join("bin", "run.sh"),
		"README.md.tpl":                "README.md",
		"Makefile.tpl":                 "Makefile",
	}

	filelistG := map[string]string{

		"src.global/svr_main.go.tpl":       path.Join("src", asset.ServerName+".go"),
		"src.global/exec/exec.go.tpl":      path.Join("src/exec", "exec_"+asset.ServerName+".go"),
		"src.global/exec/exec_impl.go.tpl": path.Join("src/exec", "exec_"+asset.ServerName+"_impl.go"),
		"src.global/exec/exec_init.go.tpl": path.Join("src/exec", "exec_"+asset.ServerName+"_init.go"),
		"client/client.go.tpl":             path.Join("client", asset.ServerName+"_client.go"),
		"conf/service.ini.tpl":             path.Join("conf", "service.ini"),
		"conf/metric.ini.tpl":             path.Join("conf", "metric.ini"),
		"conf/trace.ini.tpl":               path.Join("conf", "trace.ini"),
		"conf/log.ini.tpl":                 path.Join("conf", "log.ini"),
		"bin/run.sh.tpl":                   path.Join("bin", "run.sh"),
		"README.md.tpl":                    "README.md",
		"Makefile.tpl.global":              "Makefile",
	}

	global := options["g"].(bool)

	if global {
		return filelistG
	} else {
		return filelist
	}
}
