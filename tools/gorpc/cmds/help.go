package cmds

import (
	"flag"
	"fmt"
	"sort"
	"strings"
)

type HelpCmd struct {
	Cmd
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

func newHelpCmd() *HelpCmd {

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

	keys := c.sortedCmds()

	for _, k := range keys {
		if k == "help" {
			continue
		}
		v := cmds[k]
		b.WriteString(v.DescShort() + "\n")
	}
	return b.String()
}

func (c *HelpCmd) usageLong() string {
	b := strings.Builder{}
	b.WriteString(c.descLong + "\n")

	keys := c.sortedCmds()

	for _, k := range keys {
		if k == "help" {
			continue
		}
		v := cmds[k]
		b.WriteString(v.DescLong() + "\n")
	}
	return b.String()
}

func (c *HelpCmd) sortedCmds() []string {
	keys := make([]string, 0, len(cmds))
	for k, _ := range cmds {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
