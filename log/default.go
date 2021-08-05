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
		level:      level,
		fpath:      fp,
		writerType: FileWriter,
		rollType:   RollNONE,
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
		opts:   oo,
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
	l.tryWrite(Warn, s, v...)
}

func (l *logger) Error(s string, v ...interface{}) {
	l.tryWrite(Error, s, v...)
}

func (l *logger) Fatal(s string, v ...interface{}) {
	l.tryWrite(Fatal, s, v...)
}

func (l *logger) Flush() error {
	return l.writer.Flush()
}

func (l *logger) WithPrefix(s string, v ...interface{}) Logger {
	l.mux.Lock()
	defer l.mux.Unlock()

	p := l
	if l.parent != nil {
		p = l.parent
	}

	return &logger{
		name:   l.name,
		opts:   l.opts,
		prefix: fmt.Sprintf(s, v...),
		parent: p,
	}
}

func (l *logger) tryWrite(level Level, s string, args ...interface{}) {
	if l.opts.level > level {
		return
	}
	s = fmt.Sprintf("[%s] %s\n", level, s)
	if len(l.prefix) != 0 {
		s = fmt.Sprintf("[%s] %s %s\n", level, l.prefix, s)
	}
	s = fmt.Sprintf(s, args...)

	var err error

	if l.parent != nil {
		_, err = l.parent.writer.Write([]byte(s))
	} else {
		_, err = l.writer.Write([]byte(s))
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "file write err: %v\n", err)
	}

	if level == Fatal {
		os.Exit(1)
	}
}