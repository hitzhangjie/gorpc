package log

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

func init() {
	RegisterWriterBuilder(FileWriter, NewFileWriter)
}

type fileWriter struct {
	opts *options
	fout atomic.Value
	ch   chan []byte
	done chan struct{}
}

func (w *fileWriter) Write(b []byte) (n int, err error) {
	return w.fout.Load().(*os.File).Write(b)
}

func (w *fileWriter) AsyncWrite(b []byte) {
	w.ch <- b
}

func (w *fileWriter) asyncWrite() {
	for {
		select {
		case m := <-w.ch:
			w.fout.Load().(*os.File).Write(m)
		case <-w.done:
			w.fout.Load().(*os.File).Close()
			return
		}
	}
}

func (w *fileWriter) roll() {
	tick := time.NewTicker(time.Second)

	for range tick.C {
		fp := w.fout.Load().(*os.File).Name()
		inf, err := os.Lstat(fp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "lstat err: %v", err)
			continue
		}

		sz := inf.Size()
		if sz < int64(w.opts.maxFileSZ) {
			continue
		}

		// rename the files
		os.Rename(fp, fmt.Sprintf("%s.%d", fp, sz))
		f, err := os.OpenFile(fp, os.O_APPEND|os.O_CREATE, 0755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "open file err: %v\n", err)
			continue
		}
		w.fout.Store(f)

		time.AfterFunc(time.Second, func() {
			w.fout.Load().(*os.File).Close()
		})
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
	fw.fout.Store(fout)

	if fw.opts.async {
		go fw.asyncWrite()
	}

	if fw.opts.rollType != RollNONE {
		go fw.roll()
	}

	return fw, nil
}
