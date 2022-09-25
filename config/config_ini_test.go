package config_test

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hitzhangjie/gorpc/config"
)

var (
	iniCfg *config.IniConfig
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

/*
# possible values : production, development
app_mode = development

[paths]
# Path to where grafana can store temp files, sessions, and the sqlite3 db (if that is used)
data = /home/git/grafana

[server]
# Protocol (http or https)
protocol = http

# The http port  to use
http_port = 9999

# Redirect to correct domain if host header does not match domain
# Prevents DNS rebinding attacks
enforce_domain = true
*/
type iniServiceConfig struct {
	app_mode string
	paths    struct {
		data string
	}
	server struct {
		protocol       string
		http_port      int
		enforce_domain bool
	}
}

func TestConfig_IniConfig_ToStruct(t *testing.T) {
	vv := iniServiceConfig{}
	err := iniCfg.ToStruct(&vv)
	assert.Nil(t, err)
	fmt.Println(vv)
}
