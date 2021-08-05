package log

type options struct {
	fpath      string
	level      Level
	rollType   RollType
	writerType WriterType
	maxFileSZ  int
}

// Option options to to create a logger
type Option func(*options)

// WithRollType specifies logfile rolltype
func WithRollType(typ RollType) Option {
	return func(opts *options) {
		opts.rollType = typ
	}
}

// WithWriteType specifies the writer type
func WithWriteType(typ WriterType) Option {
	return func(opts *options) {
		opts.writerType = typ
	}
}

// WithMaxFileSZ specifies the max file size
func WithMaxFileSZ(sz int) Option {
	return func(opts *options) {
		opts.maxFileSZ = sz
	}
}
