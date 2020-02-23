package gorpc

type options struct {
	conf string
}

// Option gorpc.ListenAndServe optionns
type Option func(*options)

// WithConfigfile specify config path
func WithConfig(fpath string) Option {

	return func(opts *options) {
		opts.conf = fpath
	}
}
