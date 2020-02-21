// package gorpc provides some wrappers to quickly start gorpc service.
package gorpc

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hitzhangjie/go-rpc/config"
	"github.com/hitzhangjie/go-rpc/server"
)

const (
	unKnownServiceName = "unknown"
)

// ListenAndServe quickly initialize Server and ServerModules and serve
func ListenAndServe(opts ...Option) {

	// 处理选项
	options := options{
		configfile: "conf/service.ini",
	}
	for _, o := range opts {
		o(&options)
	}

	// 解析配置路径
	fp := options.configfile
	if !filepath.IsAbs(options.configfile) {

		self, err := os.Executable()
		if err != nil {
			panic(err)
		}
		dir, _ := filepath.Split(self)
		fp = filepath.Join(dir, options.configfile)
	}

	// 加载配置
	cfg, err := config.LoadIniConfig(fp)
	if err != nil {
		panic(err)
	}

	svr := server.NewServer()

	for _, section := range cfg.Sections() {
		// enable support for protocols
		ok := strings.HasSuffix(section, "-service")
		if !ok {
			continue
		}
		codec := strings.TrimSuffix(section, "-service")

		// initialize tcp ServerModule
		tcpport := cfg.Int(section, "tcp.port", 0)
		if tcpport > 0 {
			mod, err := server.NewTcpServerModule("tcp4", fmt.Sprintf(":%self", tcpport), codec)
			if err != nil {
				panic(err)
			}
			mod.Register(svr)

			if name := cfg.String(section, "name", ""); len(name) != 0 {
				server.NewService(name).RegisterModule(&mod)
			}
		}

		// initialize udp ServerModule
		udpport := cfg.Int(section, "udp.port", 0)
		if udpport > 0 {
			mod, err := server.NewTcpServerModule("udp4", fmt.Sprintf(":%self", udpport), codec)
			if err != nil {
				panic(err)
			}
			mod.Register(svr)

			if name := cfg.String(section, "name", ""); len(name) != 0 {
				server.NewService(name).RegisterModule(&mod)
			}
		}
	}

	svr.Start()
}
