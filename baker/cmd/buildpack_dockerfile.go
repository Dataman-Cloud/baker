package cmd

import (
	"github.com/urfave/cli"
)

var BuildpackDockerfileCmd = cli.Command{
	Name:  "dockerfile",
	Usage: "execute dockerfile command.",
	Subcommands: []cli.Command{
		BuildpackDockerfilePullCmd, // pull dockerfile from baker fileserver.
		BuildpackDockerfilePushCmd, // push dockerfile into baker fileserver.
	},
}
