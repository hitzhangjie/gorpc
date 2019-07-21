package spec

import (
	"encoding/json"
	"fmt"
	"github.com/hitzhangjie/go-rpc/tools/gorpc/log"
	"io/ioutil"
	"os"
	"path/filepath"
)

// ProtoSpec 定义了具体某个协议类型对应的一些配置信息
type ProtoSpec struct {
	Name          string `json:"name"`          //协议名称
	Handler       string `json:"handler"`       //协议handler
	Repo          string `json:"repo"`          //协议repo
	RepoPrefix    string `json:"repoPrefix"`    //协议repo导入前缀
	LocalPrefix   string `json:"localPrefix"`   //协议local路径前缀
	Asset         string `json:"asset"`         //协议资源pkg
	Spec          string `json:"spec"`          //协议定义pkg
	ClientPkg     string `json:"clientPkg"`     //协议client pkg
	ClientType    string `json:"clientType"`    //协议client type
	ClientFactory string `json:"clientFactory"` //协议client factory
	SessionType   string `json:"sessionType"`   //会话类型
}

var protoSpecs map[string]*ProtoSpec

func init() {

	log.Debug("step 0: Load config ~/.gorpc/protospec.json or /etc/gorpc/protospec.json")

	dir, err := LocateCfgPath()
	if err != nil {
		panic(err)
	}
	conf := filepath.Join(dir, "protospec.json")

	fin, err := os.Open(conf)
	if err != nil {
		msg := fmt.Sprintf("step 0: Load config file error, file:%s, err:%v", conf, err)
		log.Error(msg)
		os.Exit(1)
	}

	data, err := ioutil.ReadAll(fin)
	if err != nil {
		msg := fmt.Sprintf("step 0: Load config file error, file:%s, err:%v", conf, err)
		log.Error(msg)
		os.Exit(1)
	}

	err = json.Unmarshal(data, &protoSpecs)
	if err != nil {
		msg := fmt.Sprintf("step 0: Unmarshal config file error, file:%s, err:%v", conf, err)
		log.Error(msg)
		os.Exit(1)
	}

	log.Debug("step 0: Load config file succ, file:%s", conf)
}

func GetTypeSpec(protocol string) *ProtoSpec {

	log.Debug("param: ", protocol)
	log.Debug("specs: ", protoSpecs)

	typespec, ok := protoSpecs[protocol]
	if !ok {
		return nil
	}
	return typespec
}
