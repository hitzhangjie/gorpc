package log

import (
	"errors"
	"io"
)

// WriterType writer type
type WriterType int

const (
	FileWriter WriterType = iota // log to files
)

var writers = map[WriterType]WriterBuilder{}

// Writer defines where the logging messages are written
type Writer interface {
	io.Writer
	io.Closer
	Flush() error
}

// WriterBuilder builder of writer
type WriterBuilder func(opts *options) (Writer, error)

// RegisterWriterBuilder register new writer 'w' for type 'typ'
func RegisterWriterBuilder(typ WriterType, w WriterBuilder) {
	if w == nil {
		return
	}
	writers[typ] = w
}

func GetWriterBuilder(typ WriterType) (WriterBuilder, error) {
	if w, ok := writers[typ]; ok {
		return w, nil
	}
	return nil, errors.New("writer not found")
}
