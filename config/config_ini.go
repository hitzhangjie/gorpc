package config

import (
	iniv1 "gopkg.in/ini.v1"
)

type IniConfig struct {
	cfg *iniv1.File
}

func (c *IniConfig) LoadConfig(fp string) error {

	cfg, err := iniv1.Load(fp)
	if err != nil {
		return err
	}
	c.cfg = cfg

	return nil
}

func (c *IniConfig) Read(section, property string, dftValue string) string {
	if c.cfg == nil {
		return dftValue
	}
	val := c.cfg.Section(section).Key(property).String()
	if len(val) == 0 {
		return dftValue
	}
	return val
}

func (c *IniConfig) ReadInt(section, property string, dftValue int) int {
	val, err := c.cfg.Section(section).Key(property).Int()
	if err != nil {
		return dftValue
	}
	return val
}

func (c *IniConfig) ReadBool(section, property string, dftValue bool) bool {
	val, err := c.cfg.Section(section).Key(property).Bool()
	if err != nil {
		return dftValue
	}
	return val
}
