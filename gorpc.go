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

// ListenAndServe quickly initialize Server and ServerModules and serve
func ListenAndServe(opts ...server.Option) {

	d, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fp := filepath.Join(d, "service.ini")
	cfg, err := config.LoadIniConfig(fp)
	if err != nil {
		panic(err)
	}

	svr, err := server.NewServer(opts...)
	if err != nil {
		panic(err)
	}

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
			mod, err := server.NewTcpServer("tcp4", fmt.Sprintf(":%d", tcpport), codec)
			if err != nil {
				panic(err)
			}
			mod.Register(svr)
		}

		// initialize udp ServerModule
		udpport := cfg.Int(section, "udp.port", 0)
		if udpport > 0 {
			mod, err := server.NewTcpServer("udp4", fmt.Sprintf(":%d", udpport), codec)
			if err != nil {
				panic(err)
			}
			mod.Register(svr)
		}
	}

	svr.Start()
}
