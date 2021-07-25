package log_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hitzhangjie/gorpc/log"
)

func TestNewFileWriter_SyncWriter(t *testing.T) {
	fp := filepath.Join(os.TempDir(), "gorpc_test.log")
	fmt.Println(fp)
	_ = os.RemoveAll(fp)

	opts := []log.Option{
		log.WithWriteType(log.FileWriter),
		log.WithRollType(log.RollNONE),
		log.WithAsyncWrite(false),
	}

	l, err := log.NewLogger(fp, log.Debug, opts...)
	assert.Nil(t, err)
	assert.NotNil(t, l)
	assert.FileExists(t, fp)

	type row struct {
		level  log.Level
		text   string
		repeat int
	}
	rows := []row{
		{level: log.Trace, text: "this is one trace log message", repeat: 1},
		{level: log.Debug, text: "this is one debug log message", repeat: 2},
		{level: log.Info, text: "this is one info message", repeat: 1},
		{level: log.Warn, text: "this is one warn message", repeat: 1},
		{level: log.Error, text: "this is one error message", repeat: 1},
		//{level: log.Fatal, text: "this is one fatal message", repeat: 1}, // skip this
	}

	var sb strings.Builder

	for _, r := range rows {
		for i := 0; i < r.repeat; i++ {
			if r.level >= log.Debug {
				sb.WriteString(fmt.Sprintf("[%s] %s\n", r.level, r.text))
			}
			switch r.level {
			case log.Trace:
				l.Trace(r.text)
			case log.Debug:
				l.Debug(r.text)
			case log.Info:
				l.Info(r.text)
			case log.Warn:
				l.Warn(r.text)
			case log.Error:
				l.Error(r.text)
			default:
			}
		}
	}
	l.Flush()

	b, err := ioutil.ReadFile(fp)
	assert.Nil(t, err)
	assert.Equal(t, sb.String(), string(b))

	_ = os.RemoveAll(fp)
}

func TestNewFileWriter_AsyncWrite(t *testing.T) {

}

func TestFileWriter_RollBy_FixSize(t *testing.T) {

}

func TestFileWriter_RollBy_Day(t *testing.T) {

}
