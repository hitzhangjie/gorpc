package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

// LanguageCfg 开发语言相关的配置信息，如对应的模板工程目录、模板工程中的serverstub文件、clientstub文件
type LanguageCfg struct {
	Language      string   `json:"language"`        // required: 语言名称，如go、java
	AssetDir      string   `json:"asset_dir"`       // required: 语言对应的工程目录
	TplFileExt    string   `json:"tpl_file_ext"`    // required: 工程中模板文件的后缀名，如.tpl
	RPCServerStub string   `json:"rpc_server_stub"` // optional: 工程中对应的rpc server stub文件名（按service.method分文件生成时有用)
	RPCClientStub []string `json:"rpc_client_stub"` // required: 工程中对应的rpc client stub文件列表
}

// configs 所有语言的配置信息，汇总在此
var configs = map[string]*LanguageCfg{}

func init() {

	// 加载gorpc安装目录下的配置文件gorpc.json
	dir, err := gorpcInstallPrefix()
	fin, err := os.Open(path.Join(dir, "gorpc.json"))
	if err != nil {
		panic(err)
	}

	dat, err := ioutil.ReadAll(fin)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(dat, &configs)
	if err != nil {
		panic(err)
	}

	for k, v := range configs {
		if err := validate(k, v); err != nil {
			panic(err)
		}
	}
}

// GetLanguageCfg 加载开发语言对应的配置信息
func GetLanguageCfg(lang string) (*LanguageCfg, error) {
	cfg, ok := configs[lang]
	if !ok {
		return nil, fmt.Errorf("language:%s not supported, check config 'gorpc.json'", lang)
	}
	return cfg, nil
}

// gorpcInstallPrefix 获取gorpc安装路径，root安装到/etc/gorpc，非root用户安装到$HOME/.gorpc
func gorpcInstallPrefix() (dir string, err error) {
	u, err := user.Current()
	if err != nil {
		return
	}
	if u.Username != "root" {
		dir = filepath.Join(u.HomeDir, ".gorpc")
	} else {
		dir = "/etc/gorpc"
	}
	return
}

func validate(lang string, cfg *LanguageCfg) error {

	dir, err := gorpcInstallPrefix()
	if err != nil {
		return err
	}

	if len(lang) == 0 {
		return errors.New("invalid language, check config 'gorpc.json'")
	}
	cfg.Language = lang
	// asset dir
	if len(cfg.AssetDir) == 0 {
		return errors.New("invalid asset_dir, check config 'gorpc.json'")
	}
	if !path.IsAbs(cfg.AssetDir) {
		cfg.AssetDir = path.Join(dir, cfg.AssetDir)
	}
	// tpl_file_ext
	if len(cfg.TplFileExt) == 0 {
		return errors.New("invalid tpl_file_ext, check config 'gorpc.json'")
	}
	// rpc_server_stub, 分文件用，不设置也ok

	// rpc_client_stub，-rpconly用
	if len(cfg.RPCClientStub) == 0 {
		return errors.New("invalid rpc_client_stub, check config 'gorpc.json'")
	}
	return nil
}
