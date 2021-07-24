package log

import (
	"fmt"
	"os"
	"time"
)

func init() {
	RegisterWriterBuilder(FileWriter, NewFileWriter)
}

type fileWriter struct {
	opts *options
	fout *os.File
	ch   chan []byte
	done chan struct{}
}

func (w *fileWriter) Write(b []byte) (n int, err error) {
	return w.fout.Write(b)
}

func (w *fileWriter) AsyncWrite(b []byte) {
	w.ch <- b
}

func (w *fileWriter) asyncWrite() {
	for {
		select {
		case m := <-w.ch:
			w.fout.Write(m)
		case <-w.done:
			w.fout.Close()
			return
		}
	}
}

func (w *fileWriter) roll() {
	tick := time.NewTicker(time.Second)

	for range tick.C {
		fp := w.fout.Name()
		inf, err := os.Lstat(fp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "lstat err: %v", err)
			continue
		}

		sz := inf.Size()
		if sz < int64(w.opts.maxFileSZ) {
			continue
		}

		// TODO rename the files
	}
}

func (w *fileWriter) Close() error {
	panic("implement me")
}

func NewFileWriter(opts *options) (Writer, error) {
	fw := &fileWriter{opts: opts}

	fout, err := os.OpenFile(opts.fpath, os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	fw.fout = fout

	if fw.opts.async {
		go fw.asyncWrite()
	}

	if fw.opts.rollType != RollNONE {
		go fw.roll()
	}

	return fw, nil
}
