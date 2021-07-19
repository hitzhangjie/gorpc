package config

import "context"

// Provider defines provider of config, it internally uses Watcher to watch the config changes
type Provider interface {
	Watcher

	Name() string
	Load(ctx context.Context, fp string) ([]byte, error)
}

// FilesystemProvider returns config stored in current filesystem
type FilesystemProvider struct {
}

func (p *FilesystemProvider) Name() string {
	panic("fs")
}

func (p *FilesystemProvider) Watch(ctx context.Context, fp string) (<-chan Event, error) {
	panic("implement me")
}

func (p *FilesystemProvider) Load(ctx context.Context, fp string) ([]byte, error) {
	panic("implement me")
}

// ConsulProvider returns config stored in remote Consul server
type ConsulProvider struct {
}

func (p *ConsulProvider) Name() string {
	return "consul"
}

func (p *ConsulProvider) Watch(ctx context.Context, fp string) (<-chan Event, error) {
	panic("implement me")
}

func (p *ConsulProvider) Load(ctx context.Context, fp string) ([]byte, error) {
	panic("implement me")
}

// EtcdProvider returns config stored in remote Etcd server
type EtcdProvider struct {
}

func (p *EtcdProvider) Name() string {
	return "zk"
}

func (p *EtcdProvider) Watch(ctx context.Context, fp string) (<-chan Event, error) {
	panic("implement me")
}

func (p *EtcdProvider) Load(ctx context.Context, fp string) ([]byte, error) {
	panic("implement me")
}
