// package gorpc provides some wrappers to quickly start gorpc service.
package gorpc

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hitzhangjie/gorpc/config"
	"github.com/hitzhangjie/gorpc/registry"
	"github.com/hitzhangjie/gorpc/server"
)

// ListenAndServe quickly initialize Service and ServerModules and serve
func ListenAndServe(opts ...Option) {

	options := options{
		conf: "conf/service.ini",
	}
	for _, o := range opts {
		o(&options)
	}

	// load config
	cfg, err := loadConfig(options.conf)
	if err != nil {
		panic(err)
	}

	proc := cfg.String("service", "name", "gorpcapp")
	service := server.NewService(proc)

	// initialize transports, defined in [$codecname-service]
	for _, section := range cfg.Sections() {

		if ok := strings.HasSuffix(section.Name(), "-service"); !ok {
			continue
		}
		codec := strings.TrimSuffix(section.Name(), "-service")
		tcpport := cfg.Int(section.Name(), "tcp.port", 0)
		udpport := cfg.Int(section.Name(), "udp.port", 0)

		if err := initTransport(service, "tcp4", tcpport, codec); err != nil {
			panic(err)
		}
		if err := initTransport(service, "udp4", udpport, codec); err != nil {
			panic(err)
		}
	}

	// register to naming service
	registryName := cfg.String("service", "name", "noop")
	registry := registry.GetRegistry(registryName)
	if err := registry.Register(service); err != nil {
		panic(err)
	}
}

func loadConfig(fp string) (*config.IniConfig, error) {

	if !filepath.IsAbs(fp) {
		self, err := os.Executable()
		if err != nil {
			return nil, err
		}
		dir, _ := filepath.Split(self)
		fp = filepath.Join(dir, fp)
	}

	// load config
	cfg, err := config.NewIniConfig(fp)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func initTransport(service *server.Service, network string, port int, codec string) error {

	if !(len(network) != 0 &&
		(network == "tcp" || network == "tcp4" || network == "tcp6") ||
		(network == "udp" || network == "udp4" || network == "udp6")) {
		return fmt.Errorf("invalid network: %s", network)
	}

	if port <= 0 {
		return fmt.Errorf("invalid port: %d", port)
	}

	addr := fmt.Sprintf(":%d", port)

	err := service.ListenAndServe(context.Background(), network, addr, codec)
	if err != nil {
		return err
	}
	return nil
}
