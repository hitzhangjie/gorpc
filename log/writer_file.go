package log

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

func init() {
	RegisterWriterBuilder(FileWriter, NewFileWriter)
}

type fileWriter struct {
	sync.Mutex
	opts *options
	fout *os.File
	ch   chan []byte
	num  int64
}

func NewFileWriter(opts *options) (Writer, error) {
	fw := &fileWriter{
		opts: opts,
		num:  1,
		ch:   make(chan []byte, 1024),
	}

	fout, err := os.OpenFile(opts.fpath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	fw.fout = fout

	go fw.write()

	if fw.opts.rollType != RollNONE {
		go fw.roll()
	}

	return fw, nil
}

func (w *fileWriter) Write(b []byte) (n int, err error) {
	w.ch <- b
	return
}

func (w *fileWriter) write() {
	for {
		m, ok := <-w.ch
		if !ok {
			return
		}

		w.Lock()
		_, err := w.fout.Write(m)
		if err != nil {
			fmt.Fprintf(os.Stderr, "write error: %v\n", err)
		}
		w.Unlock()
	}
}

func (w *fileWriter) roll() {
	tick := time.NewTicker(time.Millisecond * 500)

	for range tick.C {
		w.Lock()
		fp := w.fout.Name()
		inf, err := os.Lstat(fp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "lstat err: %v", err)
			w.Unlock()
			continue
		}

		sz := inf.Size()
		if sz < int64(w.opts.maxFileSZ) {
			w.Unlock()
			continue
		}

		// rename the files
		os.Rename(fp, fmt.Sprintf("%s.%d", fp, atomic.LoadInt64(&w.num)))
		atomic.AddInt64(&w.num, 1)

		f, err := os.OpenFile(fp, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "open file err: %v\n", err)
			w.Unlock()
			continue
		}
		w.fout.Sync()
		w.fout.Close()
		w.fout = f
		w.Unlock()
	}
}

func (w *fileWriter) Flush() error {
	w.Lock()
	w.Unlock()
	return w.fout.Sync()
}

func (w *fileWriter) Close() error {
	w.Lock()
	w.Unlock()
	return w.fout.Close()
}
