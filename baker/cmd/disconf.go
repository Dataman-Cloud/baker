package cmd

import (
	"github.com/urfave/cli"
)

var DisConfCmd = cli.Command{
	Name:  "config",
	Usage: "app config file management.",
	Subcommands: []cli.Command{
		DisConfPushCmd,
		DisConfPullCmd,
		DisConfListCmd,
		DisConfDelCmd,
	},
}
