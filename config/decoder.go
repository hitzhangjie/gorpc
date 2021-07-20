package config

import (
	"encoding/json"
	"reflect"

	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v2"
)

// DecoderType 解码器类型
type DecoderType int

const (
	YAML DecoderType = iota
	JSON
	INI
)

// Decoder decoder config
type Decoder interface {
	// Decoder decode data into val, val must be pointer
	Decode(dat []byte, val interface{}) error
}

type YAMLDecoder struct {
}

func (d *YAMLDecoder) Decode(dat []byte, val interface{}) error {
	rt := reflect.TypeOf(val)
	if rt.Kind() != reflect.Ptr {
		panic("val must be pointer")
	}

	return yaml.Unmarshal(dat, val)
}

type JSONDecoder struct {
}

func (d *JSONDecoder) Decode(dat []byte, val interface{}) error {
	rt := reflect.TypeOf(val)
	if rt.Kind() != reflect.Ptr && rt.Kind() != reflect.Map {
		panic("val must be pointer or map")
	}

	return json.Unmarshal(dat, val)
}

type INIDecoder struct {
}

func (d *INIDecoder) Decode(dat []byte, val interface{}) error {
	iniCfg, ok := val.(*IniConfig)
	if !ok {
		return ini.MapTo(val, dat)
	}

	cfg, err := ini.Load(dat)
	if err != nil {
		return err
	}

	*iniCfg.cfg = *cfg
	return nil
}
