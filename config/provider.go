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

func (f FilesystemProvider) Name() string {
	panic("fs")
}

func (f FilesystemProvider) Watch(ctx context.Context, fp string) (<-chan Event, error) {
	panic("implement me")
}

func (f FilesystemProvider) Load(ctx context.Context, fp string) ([]byte, error) {
	panic("implement me")
}

// ConsulProvider returns config stored in remote Consul server
type ConsulProvider struct {
}

func (c ConsulProvider) Name() string {
	return "consul"
}

func (c *ConsulProvider) Watch(ctx context.Context, fp string) (<-chan Event, error) {
	panic("implement me")
}

func (c ConsulProvider) Load(ctx context.Context, fp string) ([]byte, error) {
	panic("implement me")
}

// EtcdProvider returns config stored in remote Etcd server
type EtcdProvider struct {
}

func (Z EtcdProvider) Name() string {
	return "zk"
}

func (Z *EtcdProvider) Watch(ctx context.Context, fp string) (<-chan Event, error) {
	panic("implement me")
}

func (Z EtcdProvider) Load(ctx context.Context, fp string) ([]byte, error) {
	panic("implement me")
}
