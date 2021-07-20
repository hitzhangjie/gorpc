package config_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hitzhangjie/gorpc/config"
)

func TestLoader_Load(t *testing.T) {
	opts := []config.Option{
		config.WithProvider(&config.FilesystemProvider{}),
		config.WithDecoder(&config.INIDecoder{}),
	}
	ld, err := config.NewLoader(context.TODO(), opts...)
	if err != nil {
		t.Fatalf("new loader error: %v", err)
	}

	d, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	c, err := ld.Load(context.TODO(), filepath.Join(d, "testdata/service.ini"))
	if err != nil {
		t.Fatalf("load service.ini error: %v", err)
	}

	assert.Equal(t, "development", c.Read("app_mode", ""))
	assert.Equal(t, "/home/git/grafana", c.Read("paths.data", ""))
}
