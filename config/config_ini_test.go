package config_test

import (
	"github.com/hitzhangjie/go-rpc/config"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var (
	cwd    string
	err    error
	iniCfg *config.IniConfig
)

func TestMain(m *testing.M) {

	cwd, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	iniCfg, err = config.LoadIniConfig(filepath.Join(cwd, "testdata/service.ini"))
	if err != nil {
		panic(err)
	}
	m.Run()
}

func TestIniConfig(t *testing.T) {
	type args struct {
		section  string
		property string
		dftValue interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"load-[]-app_mode", args{"", "app_mode", ""}, "development"},
		{"load-[]-not_existed", args{"", "not_existed", "xxx"}, "xxx"},
		{"load-[]=not_existed_again", args{"", "not_existed", "yyy"}, "yyy"},
		{"load-[paths]-data", args{"paths", "data", ""}, "/home/git/grafana"},
		{"load-[server]-protocol", args{"server", "protocol", ""}, "http"},
		{"load-[server]-http_port", args{"server", "http_port", ""}, "9999"},
		{"load-[server]-enforce_domain", args{"server", "enforce_domain", ""}, "true"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got interface{}
			switch tt.args.dftValue.(type) {
			case string:
				v := iniCfg.String(tt.args.section, tt.args.property, tt.args.dftValue.(string))
				got = v
			case int8:
				v := iniCfg.Int(tt.args.section, tt.args.property, tt.args.dftValue.(int))
				got = v
			case bool:
				v := iniCfg.Bool(tt.args.section, tt.args.property, tt.args.dftValue.(bool))
				got = v
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("case:%s, got = %v, want = %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestIniConfig_ImplConfig(t *testing.T) {
	type args struct {
		key      string
		dftValue interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"load-[]-app_mode", args{"app_mode", ""}, "development"},
		{"load-[]-not_existed", args{"not_existed", "xxx"}, "xxx"},
		{"load-[]=not_existed_again", args{"not_existed", "yyy"}, "yyy"},
		{"load-[paths]-data", args{"paths.data", ""}, "/home/git/grafana"},
		{"load-[server]-protocol", args{"server.protocol", ""}, "http"},
		{"load-[server]-http_port", args{"server.http_port", ""}, "9999"},
		{"load-[server]-enforce_domain", args{"server.enforce_domain", ""}, "true"},
	}
	var cfg config.Config = iniCfg
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got interface{}
			switch tt.args.dftValue.(type) {
			case string:
				v := cfg.Read(tt.args.key, tt.args.dftValue.(string))
				got = v
			case int8:
				v := cfg.ReadInt(tt.args.key, tt.args.dftValue.(int))
				got = v
			case bool:
				v := cfg.ReadBool(tt.args.key, tt.args.dftValue.(bool))
				got = v
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("case:%s, got = %v, want = %v", tt.name, got, tt.want)
			}
		})
	}
}
