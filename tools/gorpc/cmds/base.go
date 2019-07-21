package cmds

import (
	"flag"
)

// Commander defines the subcmd behavior
type Commander interface {
	// subcmd usage
	UsageLine() string
	DescShort() string
	DescLong() string
	// subcmd params
	FlagSet() *flag.FlagSet
	// subcmd logic
	Run(args ...string) error
}

// Cmd defines the subcmd base behavior
type Cmd struct {
	usageLine string
	descShort string
	descLong  string
	flagSet   *flag.FlagSet
}

// UsageLine returns usage line
func (c *Cmd) UsageLine() string {
	return c.usageLine
}

// DescShort returns the short description
func (c *Cmd) DescShort() string {
	return c.descShort
}

// DescLong returns the long description
func (c *Cmd) DescLong() string {
	return c.descLong
}

// FlagSet returns the flagset
func (c *Cmd) FlagSet() *flag.FlagSet {
	return c.flagSet
}

