package config_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

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

func TestLoader_ReLoad(t *testing.T) {
	opts := []config.Option{
		config.WithProvider(&config.FilesystemProvider{}),
		config.WithDecoder(&config.INIDecoder{}),
		config.WithReload(true),
	}
	ld, err := config.NewLoader(context.TODO(), opts...)
	if err != nil {
		t.Fatalf("new loader error: %v", err)
	}

	d, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fp := filepath.Join(d, "testdata/service.ini")
	c, err := ld.Load(context.TODO(), fp)
	if err != nil {
		t.Fatalf("load service.ini error: %v", err)
	}

	assert.Equal(t, "development", c.Read("app_mode", ""))

	// change testdata/service.ini
	b, err := os.ReadFile(fp)
	if err != nil {
		panic(err)
	}
	defer os.WriteFile(fp, b, 0666)

	n := strings.ReplaceAll(string(b), "app_mode = development", "app_mode = production")
	err = os.WriteFile(fp, []byte(n), 0666)
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second*3)
	assert.Equal(t, "production", c.Read("app_mode", ""))
}
