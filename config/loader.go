package config

import (
	"context"
	"time"
)

// Loader load the config, it may internally uses Provider to read config
type Loader interface {
	// Load returns config abount fp
	Load(ctx context.Context, fp string, opts ...Option) (Config, error)
}

type loader struct {
	opts   *options
	config *config
}

func NewLoader(ctx context.Context, opts ...Option) (Loader, error) {
	oo := options{}
	for _, o := range opts {
		o(&oo)
	}

	return &loader{opts: &oo}, nil
}

func (l *loader) Load(ctx context.Context, fp string, opts ...Option) (Config, error) {
	oo := options{
		fp:       l.opts.fp,
		reload:   l.opts.reload,
		interval: l.opts.interval,
		decoder:  l.opts.decoder,
		provider: l.opts.provider,
	}

	for _, o := range opts {
		o(&oo)
	}

	dat, err := oo.provider.Load(ctx, fp)
	if err != nil {
		return nil, err
	}

	var ldcfg interface{}

	switch v := oo.decoder.(type) {
	case *YAMLDecoder:
		cfg := YamlConfig{}
		err = v.Decode(dat, cfg.yml)
		ldcfg = cfg
	case *INIDecoder:
		cfg := IniConfig{}
		err = v.Decode(dat, cfg.cfg)
		ldcfg = cfg
	default:
		panic("not supported")
	}
	if err != nil {
		return nil, err
	}

	l.config.value.Store(ldcfg)
	return l.config, nil
}

// Option loader options
type Option func(*options)

type options struct {
	fp       string
	reload   bool
	interval time.Duration
	decoder  Decoder
	provider Provider
}

func WithReload(v bool) Option {
	return func(o *options) {
		o.reload = v
	}
}

func WithInterval(v time.Duration) Option {
	return func(o *options) {
		o.interval = v
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
