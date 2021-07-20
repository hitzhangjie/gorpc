package config

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
)

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

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	ch := make(chan Event)
	go func() {
		defer func() {
			watcher.Close()
			close(ch)
		}()

		for {
			select {
			case ev, ok := <-watcher.Events:
				if !ok {
					fmt.Println("watch error: wather.Events closed")
					return
				}
				if ev.Op&fsnotify.Write == 0 {
					continue
				}
				//fmt.Println("watch file modified:", ev.Name)
				ch <- Event{
					typ:  Update,
					meta: ev.String(),
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("watch error:", err)
				return
			case <-ctx.Done():
				fmt.Println("watch cancelled")
				return
			}
		}
	}()

	err = watcher.Add(fp)
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (p *FilesystemProvider) Load(ctx context.Context, fp string) ([]byte, error) {
	fin, err := os.Lstat(fp)
	if err != nil {
		return nil, err
	}

	if fin.IsDir() {
		return nil, errors.New("not normal file")
	}

	return os.ReadFile(fp)
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
