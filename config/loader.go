package config

import (
	"context"

	"gopkg.in/ini.v1"
)

// Loader load the config, it may internally uses Provider to read config
type Loader interface {
	// Load returns config abount fp
	Load(ctx context.Context, fp string, opts ...Option) (Config, error)
}

type loader struct {
	opts   options
	config config
}

func NewLoader(ctx context.Context, opts ...Option) (Loader, error) {
	oo := options{}
	for _, o := range opts {
		o(&oo)
	}

	return &loader{opts: oo}, nil
}

func (l *loader) Load(ctx context.Context, fp string, opts ...Option) (Config, error) {
	oo := options{
		fp:       fp,
		reload:   l.opts.reload,
		decoder:  l.opts.decoder,
		provider: l.opts.provider,
	}

	for _, o := range opts {
		o(&oo)
	}

	if oo.reload {
		l.reload(ctx, fp, oo)
	}

	return l.load(ctx, fp, oo)
}

func (l *loader) load(ctx context.Context, fp string, opts options) (Config, error) {
	dat, err := opts.provider.Load(ctx, fp)
	if err != nil {
		return nil, err
	}

	cfg, err := l.decode(ctx, dat, opts)
	if err != nil {
		return nil, err
	}

	l.config.value.Store(cfg)
	return &l.config, nil
}

func (l *loader) reload(ctx context.Context, fp string, opts options) error {
	ch, err := opts.provider.Watch(ctx, fp)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case v := <-ch:
				if v.typ != Update {
					continue
				}
				l.load(ctx, fp, opts)
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (l *loader) decode(ctx context.Context, dat []byte, opts options) (interface{}, error) {
	var cfg interface{}
	var err error

	switch v := opts.decoder.(type) {
	case *YAMLDecoder:
		c := YamlConfig{}
		err = v.Decode(dat, c.yml)
		cfg = &c
	case *INIDecoder:
		c := IniConfig{
			cfg: &ini.File{},
		}
		err = v.Decode(dat, &c)
		cfg = &c
	default:
		panic("not supported decoder type")
	}

	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// Option loader options
type Option func(*options)

type options struct {
	fp       string
	reload   bool
	decoder  Decoder
	provider Provider
}

func WithReload(v bool) Option {
	return func(o *options) {
		o.reload = v
	}
}

func WithDecoder(v Decoder) Option {
	return func(o *options) {
		o.decoder = v
	}
}

func WithProvider(v Provider) Option {
	return func(o *options) {
		o.provider = v
	}
}
