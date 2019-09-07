package cmds

import (
	"flag"
	"os/user"
	"path/filepath"
)

// Commander defines the subcmd behavior
type Commander interface {
	// UsageLine cmd example
	UsageLine() string

	// DescShort cmd brief description
	DescShort() string

	// DescLong cmd detailed description
	DescLong() string

	// FlagSet cmd flagset
	FlagSet() *flag.FlagSet

	// Run cmd run
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

func defaultAssetDir() (dir string, err error) {
	u, err := user.Current()
	if err != nil {
		return
	}
	if u.Username != "root" {
		dir = filepath.Join(u.HomeDir, ".gorpc/asset")
	} else {
		dir = "/etc/gorpc/assetdir"
	}
	return
}
