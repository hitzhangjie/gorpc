package log_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hitzhangjie/gorpc/log"
)

type row struct {
	level  log.Level
	text   string
	repeat int
}

var rows = []row{
	{level: log.Trace, text: "this is one trace log message", repeat: 1},
	{level: log.Debug, text: "this is one debug log message", repeat: 2},
	{level: log.Info, text: "this is one info message", repeat: 1},
	{level: log.Warn, text: "this is one warn message", repeat: 1},
	{level: log.Error, text: "this is one error message", repeat: 10},
	//{level: log.Fatal, text: "this is one fatal message", repeat: 1}, // skip this
}

func expect(rows []row) string {
	var sb strings.Builder
	for _, r := range rows {
		for i := 0; i < r.repeat; i++ {
			if r.level >= log.Debug {
				sb.WriteString(fmt.Sprintf("[%s] %s\n", r.level, r.text))
			}
		}
	}
	return sb.String()
}

func TestNewFileWriter_Write(t *testing.T) {
	fp := filepath.Join(os.TempDir(), "gorpc_logwrite.log")
	testWrite(t, fp, log.RollNONE, 0, 0)

	b, err := ioutil.ReadFile(fp)
	assert.Nil(t, err)
	assert.Equal(t, expect(rows), string(b))

	_ = os.Remove(fp)
}

func TestFileWriter_Roll(t *testing.T) {
	t.Run("roll by filesz", func(t *testing.T) {
		fp := filepath.Join(os.TempDir(), "gorpc_logwrite_rollby_filesz.log")
		testWrite(t, fp, log.RollByFileSZ, 16, time.Second)
		files, err := filepath.Glob(fp + "*")
		assert.Nil(t, err)
		assert.NotEmpty(t, files)
		assert.Greater(t, len(files), 1)
		for _, f := range files {
			_ = os.Remove(f)
		}
	})

	t.Run("roll by day", func(t *testing.T) {
		fp := filepath.Join(os.TempDir(), "gorpc_logwrite_rollby_day.log")
		testWrite(t, fp, log.RollByDay, 16, time.Second)
		_ = os.RemoveAll(filepath.Dir(fp))
	})
}

func testWrite(t *testing.T, fp string, rolltyp log.RollType, rollsz int, d time.Duration) {
	opts := []log.Option{
		log.WithWriteType(log.FileWriter),
		log.WithRollType(rolltyp),
		log.WithMaxFileSZ(rollsz),
	}

	l, err := log.NewLogger(fp, log.Debug, opts...)
	assert.Nil(t, err)
	assert.NotNil(t, l)
	assert.FileExists(t, fp)

	for _, r := range rows {
		for i := 0; i < r.repeat; i++ {
			time.Sleep(d)

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
}
