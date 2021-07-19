package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

var fsProvider = &FilesystemProvider{}

func TestFilesystemProvider_Load(t *testing.T) {
	d, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fp := filepath.Join(d, "testdata/service.ini")

	b, err := fsProvider.Load(context.TODO(), fp)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	t.Logf("load ok: %s", string(b))
}

func TestFilesystemProvider_Watch(t *testing.T) {
	d, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fp := filepath.Join(d, "testdata/change.ini")

	ch, err := fsProvider.Watch(context.TODO(), fp)
	if err != nil {
		t.Fatalf("watch error: %v", err)
	}

	for {
		select {
		case ev, ok := <-ch:
			if !ok {
				break
			}
			t.Logf("load ok: %s", ev.meta)
		default:
		}
	}
}
