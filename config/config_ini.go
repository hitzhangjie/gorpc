package config

import (
	"strings"

	"gopkg.in/ini.v1"
)

// IniConfig ini config
type IniConfig struct {
	cfg *ini.File
}

// NewIniConfig create a new config from ini configfile `fp`
func NewIniConfig(fp string) (*IniConfig, error) {
	cfg, err := ini.Load(fp)
	if err != nil {
		return nil, err
	}
	return &IniConfig{cfg}, nil
}

func (c *IniConfig) Sections() []*ini.Section {
	if c.cfg == nil {
		return nil
	}
	return c.cfg.Sections()
}

func (c *IniConfig) String(section, property string, dftValue string) string {
	if c.cfg == nil {
		return dftValue
	}
	val := c.cfg.Section(section).Key(property).String()
	if len(val) == 0 {
		return dftValue
	}
	return val
}

func (c *IniConfig) Int(section, property string, dftValue int) int {
	val, err := c.cfg.Section(section).Key(property).Int()
	if err != nil {
		return dftValue
	}
	return val
}

func (c *IniConfig) Bool(section, property string, dftValue bool) bool {
	val, err := c.cfg.Section(section).Key(property).Bool()
	if err != nil {
		return dftValue
	}
	return val
}

func (c *IniConfig) Read(key string, dftValue string) string {
	s, p := c.split(key)
	return c.String(s, p, dftValue)
}

func (c *IniConfig) ReadInt(key string, dftValue int) int {
	s, p := c.split(key)
	return c.Int(s, p, dftValue)
}

func (c *IniConfig) ReadBool(key string, dftValue bool) bool {
	s, p := c.split(key)
	return c.Bool(s, p, dftValue)
}

func (c *IniConfig) ToStruct(cfg interface{}) error {
	return c.cfg.StrictMapTo(cfg)
}

func (c *IniConfig) split(key string) (string, string) {
	v := strings.SplitN(key, ".", 2)
	switch len(v) {
	case 0:
		return "", ""
	case 1:
		return "", v[0]
	case 2:
		return v[0], v[1]
	default:
		return "", ""
	}
}
