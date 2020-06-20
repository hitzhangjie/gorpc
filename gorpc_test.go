package gorpc_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	gorpc "github.com/hitzhangjie/gorpc"
	_ "github.com/hitzhangjie/gorpc/codec/whisper"
)

// Test `gorpc.ListenAndServe`
func TestListenAndServe(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fp := filepath.Join(dir, "testcase/service.ini")

	gorpc.ListenAndServe(gorpc.WithConfig(fp))

	// Linux: run `fuser port/tcp` or `fuser port/udp` to check whether server working
	// macOS: run `lsof -i tcp:port` or `lsof -i udp:port` to check whether server working
	time.Sleep(time.Second * 10)
}
