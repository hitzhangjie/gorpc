package metrics

var options = &Options{}

type Options struct {
	Meta          map[string]interface{}
	HistogramMeta map[string]HistogramMeta
}

type HistogramMeta struct {
	BucketRange
	value interface{}
}

func GetOptions() Options {
	return *options
}

type Option func(opts *Options)

func WithMeta(meta map[string]interface{}) Option {
	return func(opts *Options) {
		if opts == nil {
			return
		}
		if opts.Meta == nil {
			opts.Meta = meta
			return
		}
		for k, v := range meta {
			opts.Meta[k] = v
		}
	}
}

func WithHistogramMeta(meta map[string]HistogramMeta) Option {
	return func(opts *Options) {
		if opts == nil {
			return
		}
		if opts.HistogramMeta == nil {
			opts.HistogramMeta = meta
			return
		}
		for k, v := range meta {
			opts.HistogramMeta[k] = v
		}
	}
}
