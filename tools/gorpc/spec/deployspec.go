package spec

import (
    "encoding/json"
    "fmt"
    "git.code.oa.com/go-neat/tools/codegen/log"
    "io/ioutil"
    "os"
    "path/filepath"
)

// DeploySpec 定义了具体某个用户类型对应的deploy配置信息
type DeploySpec struct {
    Ips          string `json:"ips"`          //ip列表信息
    Author       string `json:"author"`       //作者rtx
}

// name => DeploySpec
var deploySpecs DeploySpec

func init() {

    dir, err := LocateCfgPath()
    if err != nil {
        panic(err)
    }
    conf := filepath.Join(dir, "deployspec.json")

    fin, err := os.Open(conf)
    if err != nil {
        msg := fmt.Sprintf("Read config file error, file:%s, err:%v", conf, err)
        log.Error(msg)
        // 文件没读成功，程序不结束，使用默认值
        return
    }

    data, err := ioutil.ReadAll(fin)
    if err != nil {
        msg := fmt.Sprintf("Read config file error, file:%s, err:%v", conf, err)
        log.Error(msg)
        // 文件没读成功，程序不结束，使用默认值
        return
    }

    err = json.Unmarshal(data, &deploySpecs)
    if err != nil {
        msg := fmt.Sprintf("Unmarshal config file error, file:%s, err:%v", conf, err)
        log.Error(msg)
        // 文件没读成功，程序不结束，使用默认值
        return
    }

    log.Debug("Load config file succ, file:%s", conf)
}

func GetDeploySpec() *DeploySpec {
    return &deploySpecs
}
