// Package config provides support of config loading and reading for
// different types of config, including ini, yaml, etc.
package config

import (
	"errors"
	"fmt"
	"os"
	"sync/atomic"
)

// Config config
type Config interface {
	// Read, read string value by `key`, if not found, return dftValue
	Read(key string, dftValue string) string

	// ReadInt, read int value by `key`, if not found, return dftValue
	ReadInt(key string, dftValue int) int

	// ReadBool, read bool value by `key`, if not found, return dftValue
	ReadBool(key string, dftValue bool) bool

	// ToStruct convert config data to specified struct
	ToStruct(cfg interface{}) error
}

type config struct {
	// reload config if needed
	value atomic.Value
	opts  options
}

func (c *config) Read(key string, dftValue string) string {
	cfg := c.value.Load()

	switch v := cfg.(type) {
	case *YamlConfig:
		return v.Read(key, dftValue)
	case *IniConfig:
		return v.Read(key, dftValue)
	default:
		fmt.Fprintf(os.Stderr, "not supported config: %T\n", v)
		return dftValue
	}
}

func (c *config) ReadInt(key string, dftValue int) int {
	cfg := c.value.Load()

	switch v := cfg.(type) {
	case *YamlConfig:
		return v.ReadInt(key, dftValue)
	case *IniConfig:
		return v.ReadInt(key, dftValue)
	default:
		fmt.Fprintln(os.Stderr, "not supported config: ", v)
		return dftValue
	}
}

func (c *config) ReadBool(key string, dftValue bool) bool {
	cfg := c.value.Load()

	switch v := cfg.(type) {
	case *YamlConfig:
		return v.ReadBool(key, dftValue)
	case *IniConfig:
		return v.ReadBool(key, dftValue)
	default:
		fmt.Fprintln(os.Stderr, "not supported config: ", v)
		return dftValue
	}
}

func (c *config) ToStruct(v interface{}) error {
	cfg := c.value.Load()

	switch v := cfg.(type) {
	case *YamlConfig:
		return v.ToStruct(v)
	case *IniConfig:
		return v.ToStruct(v)
	default:
		return errors.New("invalid Config type")
	}
}
