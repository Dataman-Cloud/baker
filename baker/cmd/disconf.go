package cmd

import (
	"github.com/urfave/cli"
)

var DisConfCmd = cli.Command{
	Name:  "disconf",
	Usage: "app property file management.",
	Subcommands: []cli.Command{
		DisConfPushCmd,
		DisConfPullCmd,
		DisConfListCmd,
		DisConfDelCmd,
		DisConfUnzipCmd,
	},
}
