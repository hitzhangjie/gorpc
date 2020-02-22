package gorpc_test

import (
	"os"
	"path/filepath"
	"testing"

	gorpc "github.com/hitzhangjie/go-rpc"
	_ "github.com/hitzhangjie/go-rpc/codec/whisper"
)

func TestListenAndServe(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fp := filepath.Join(dir, "testcase/service.ini")

	gorpc.ListenAndServe(gorpc.WithConfig(fp))
}
