package log

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var (
	mux     = sync.Mutex{}
	loggers = map[string]Logger{}
)

func NewLogger(fname string, level Level, opts ...Option) (Logger, error) {
	// test if this logger created before
	fp, err := filepath.Abs(fname)
	if err != nil {
		return nil, err
	}

	// try to create a new logger
	mux.Lock()
	defer mux.Unlock()

	l, ok := loggers[fp]
	if ok && l != nil {
		return l, nil
	}

	_ = os.MkdirAll(filepath.Dir(fp), os.ModePerm)

	oo := options{
		level:      Info,
		fpath:      fp,
		writerType: FileWriter,
		rollType:   RollNONE,
		async:      true,
	}
	for _, o := range opts {
		o(&oo)
	}

	b, err := GetWriterBuilder(oo.writerType)
	if err != nil {
		return nil, err
	}
	writer, err := b(&oo)
	if err != nil {
		return nil, err
	}

	return &logger{
		name:   fp,
		mux:    new(sync.Mutex),
		writer: writer,
		opts:   options{},
	}, nil
}

type logger struct {
	name string

	mux    *sync.Mutex
	writer Writer
	opts   options
	prefix string

	parent *logger
}

func (l *logger) Trace(s string, v ...interface{}) {
	l.tryWrite(Trace, s, v...)
}

func (l *logger) Debug(s string, v ...interface{}) {
	l.tryWrite(Debug, s, v...)
}

func (l *logger) Info(s string, v ...interface{}) {
	l.tryWrite(Info, s, v...)
}

func (l *logger) Warn(s string, v ...interface{}) {
	l.tryWrite(Info, s, v...)
}

func (l *logger) Error(s string, v ...interface{}) {
	l.tryWrite(Error, s, v...)
}

func (l *logger) Fatal(s string, v ...interface{}) {
	l.tryWrite(Fatal, s, v...)
}

func (l *logger) WithPrefix(s string, v ...interface{}) Logger {
	s1 := fmt.Sprintf("%s %s", l.prefix, s)
	s1 = fmt.Sprintf(s1, v...)

	p := l
	if l.parent != nil {
		p = l.parent
	}

	return &logger{
		name:   l.name,
		opts:   l.opts,
		prefix: s1,
		parent: p,
	}
}

func (l *logger) tryWrite(level Level, s string, args ...interface{}) {
	if l.opts.level > level {
		return
	}
	s0 := fmt.Sprintf("[%s] %s %s", level, l.prefix, s)
	s1 := fmt.Sprintf(s0, args...)

	var n int
	var err error

	if l.parent != nil {
		n, err = l.parent.writer.Write([]byte(s1))
	} else {
		n, err = l.writer.Write([]byte(s1))
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "file write err: %v\n", err)
	}
	if n != len(s1) {
		fmt.Fprintf(os.Stderr, "file write: only %d/%d bytes written", n, len(s1))
	}

	if level == Fatal {
		os.Exit(1)
	}
}

type options struct {
	fpath      string
	level      Level
	rollType   RollType
	writerType WriterType
	maxFileSZ  int
	async      bool
}

// Option options to to create a logger
type Option func(*options)

// WithRollType specifies logfile rolltype
func WithRollType(typ RollType) Option {
	return func(opts *options) {
		opts.rollType = typ
	}
}

// WithAsyncWrite enable async write
func WithAsyncWrite(async bool) Option {
	return func(opts *options) {
		opts.async = async
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
