package config_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/hitzhangjie/gorpc/config"
)

var (
	iniCfg  *config.IniConfig
	yamlCfg *config.YamlConfig
)

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// load service.ini
	ini, err := config.NewIniConfig(filepath.Join(cwd, "testdata/service.ini"))
	if err != nil {
		panic(err)
	}
	iniCfg = ini

	// load service.yml
	yml, err := config.NewYamlConfig(filepath.Join(cwd, "testdata/service.yml"))
	if err != nil {
		panic(err)
	}
	yamlCfg = yml

}

func TestConfig_IniConfig(t *testing.T) {
	type args struct {
		key      string
		dftValue interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{"load-[]-app_mode", args{"app_mode", ""}, "development"},
		{"load-[]-not_existed", args{"not_existed", "xxx"}, "xxx"},
		{"load-[]=not_existed_again", args{"not_existed", "yyy"}, "yyy"},
		{"load-[paths]-data", args{"paths.data", ""}, "/home/git/grafana"},
		{"load-[server]-protocol", args{"server.protocol", ""}, "http"},
		{"load-[server]-http_port", args{"server.http_port", 0}, 9999},
		{"load-[server]-enforce_domain", args{"server.enforce_domain", false}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got interface{}
			switch tt.args.dftValue.(type) {
			case string:
				v := iniCfg.Read(tt.args.key, tt.args.dftValue.(string))
				got = v
			case int:
				v := iniCfg.ReadInt(tt.args.key, tt.args.dftValue.(int))
				got = v
			case bool:
				v := iniCfg.ReadBool(tt.args.key, tt.args.dftValue.(bool))
				got = v
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("case:%s, got = %v, want = %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestConfig_YamlConfig(t *testing.T) {
	type args struct {
		path     string
		dftValue interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{"load-name", args{"name", ""}, "smallfish"},
		{"load-age", args{"age", 0}, 99},
		{"load-bool", args{"bool", false}, true},
		{"load-bb.cc.dd.ee", args{"bb.cc.ee", ""}, "aaa"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got interface{}
			switch tt.args.dftValue.(type) {
			case string:
				v := yamlCfg.Read(tt.args.path, tt.args.dftValue.(string))
				got = v
			case int:
				v := yamlCfg.ReadInt(tt.args.path, tt.args.dftValue.(int))
				got = v
			case bool:
				v := yamlCfg.ReadBool(tt.args.path, tt.args.dftValue.(bool))
				got = v
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("case:%s, got = %v, want = %v", tt.name, got, tt.want)
			}
		})
	}
}
