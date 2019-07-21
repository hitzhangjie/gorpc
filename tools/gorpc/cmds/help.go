package cmds

import (
	"flag"
	"fmt"
	"strings"
)

type HelpCmd struct {
	Cmd
}

func init() {
	mux.Lock()
	defer mux.Unlock()
	all["help"] = NewHelpCmd()
}

func (c *HelpCmd) Run(args ...string) error {

	c.FlagSet().Parse(args)
	verbose := c.FlagSet().Lookup("v").Value.(flag.Getter).Get().(bool)

	var tip string
	if verbose {
		tip = c.usageLong()
	} else {
		tip = c.usageShort()
	}
	fmt.Println(tip)
	return nil
}

func NewHelpCmd() *HelpCmd {

	fs := flag.NewFlagSet("helpcmd", flag.ContinueOnError)
	fs.Bool("v", false, "verbose help info")

	u := Cmd{
		usageLine: "gorpc help",
		descShort: `
how to display help:
	gorpc help`,
		descLong: `
gorpc <cmd> <options>: 

global options:
	-h display this help
	-v display verbose info`,
		flagSet: fs,
	}

	return &HelpCmd{u}
}

func (c *HelpCmd) usageShort() string {
	b := strings.Builder{}
	b.WriteString(c.descShort + "\n")

	for k, v := range all {
		if k == "help" {
			continue
		}
		b.WriteString(v.DescShort() + "\n")
	}
	return b.String()
}

func (c *HelpCmd) usageLong() string {
	b := strings.Builder{}
	b.WriteString(c.descLong + "\n")

	for k, v := range all {
		if k == "help" {
			continue
		}
		b.WriteString(v.DescLong() + "\n")
	}
	return b.String()
}
