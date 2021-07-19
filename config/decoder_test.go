package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestYAMLDecoder_Decode(t *testing.T) {
	d := &YAMLDecoder{}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fp := filepath.Join(cwd, "testdata/service.yml")
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		panic(err)
	}

	v := yamlStruct{}
	err = d.Decode(b, &v)
	if err != nil {
		t.Fatalf("yaml decode error: %v", err)
	}
	t.Logf("yaml decode ok, data: %+v", v)
}

type yamlStruct struct {
	Name   string   `yaml:"name"`
	Age    int      `yaml:"age"`
	Float  float64  `yaml:"float"`
	Bool   bool     `yaml:"bool"`
	Emails []string `yaml:"emails"`
	Bb     struct {
		Cc struct {
			Dd []int  `yaml:"dd"`
			Ee string `yaml:"ee"`
		} `yaml:"cc"`
	} `yaml:"bb"`
}

func TestINIDecoder_Decode(t *testing.T) {
	d := &INIDecoder{}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fp := filepath.Join(cwd, "testdata/service.ini")
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		panic(err)
	}

	v := iniStruct{}
	err = d.Decode(b, &v)
	if err != nil {
		t.Fatalf("ini decode error: %v", err)
	}
	t.Logf("ini decode ok, data: %+v", v)
}

type iniStruct struct {
	AppMode string `ini:"app_mode"`

	Paths struct {
		Data string `ini:"data"`
	} `ini:"[paths]"`
}
