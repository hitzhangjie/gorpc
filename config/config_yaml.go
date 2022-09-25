package config

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/smallfish/simpleyaml"
)

// YamlConfig yaml config
type YamlConfig struct {
	yml *simpleyaml.Yaml
}

// NewYamlConfig create a new config from yaml configfile `fp`
func NewYamlConfig(fp string) (*YamlConfig, error) {
	fin, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(fin)
	if err != nil {
		return nil, err
	}

	yml, err := simpleyaml.NewYaml(data)
	if err != nil {
		return nil, err
	}

	return &YamlConfig{yml}, nil
}

func (c *YamlConfig) Read(key string, dftValue string) string {
	if c.yml == nil {
		return dftValue
	}

	paths := c.path(key)

	v, err := c.yml.GetPath(paths...).String()
	if err != nil || len(v) == 0 {
		return dftValue
	}
	return v
}

func (c *YamlConfig) ReadInt(key string, dftValue int) int {
	if c.yml == nil {
		return dftValue
	}

	path := c.path(key)

	v, err := c.yml.GetPath(path...).Int()
	if err != nil || v == 0 {
		return dftValue
	}
	return v
}

func (c *YamlConfig) ReadBool(key string, dftValue bool) bool {
	if c.yml == nil {
		return dftValue
	}

	path := c.path(key)

	v, err := c.yml.GetPath(path...).Bool()
	if err != nil || !v {
		return dftValue
	}
	return v
}

func (c *YamlConfig) path(key string) []interface{} {
	path := strings.Split(key, ".")
	paths := make([]interface{}, len(path), len(path))
	for idx, p := range path {
		paths[idx] = p
	}
	return paths
}

// YamlConfigLoader yaml config loader
type YamlConfigLoader struct {
}
