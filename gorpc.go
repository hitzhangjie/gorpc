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

// ListenAndServe quickly initialize Service and ServerModules and serve
func ListenAndServe(opts ...Option) {

	options := options{
		configfile: "conf/service.ini",
	}
	for _, o := range opts {
		o(&options)
	}

	// parse config
	fp := options.configfile
	if !filepath.IsAbs(options.configfile) {

		self, err := os.Executable()
		if err != nil {
			panic(err)
		}
		dir, _ := filepath.Split(self)
		fp = filepath.Join(dir, options.configfile)
	}

	// load config
	cfg, err := config.LoadIniConfig(fp)
	if err != nil {
		panic(err)
	}

	self, _ := os.Executable()

	service := server.NewService(self)

	for _, section := range cfg.Sections() {

		// enable support for protocols
		ok := strings.HasSuffix(section, "-service")
		if !ok {
			continue
		}
		codec := strings.TrimSuffix(section, "-service")

		// initialize tcp Transport
		tcpport := cfg.Int(section, "tcp.port", 0)
		if tcpport > 0 {
			err := service.ListenAndServe("tcp4", fmt.Sprintf(":%s", tcpport), codec)
			if err != nil {
				panic(err)
			}
		}

		// initialize udp Transport
		udpport := cfg.Int(section, "udp.port", 0)
		if udpport > 0 {
			err := service.ListenAndServe("udp4", fmt.Sprintf(":%self", udpport), codec)
			if err != nil {
				panic(err)
			}
		}
	}

	// register to naming service
	for _, mod := range service.ServerModules() {
		section := mod.Codec() + "-service"
		if name := cfg.String(section, "name", ""); len(name) != 0 {
			// fixme nameing service register this mod.Net+mod.Address+mod.Codec
		}
	}

	//service.Start()
}
