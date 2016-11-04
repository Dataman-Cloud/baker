package main

import (
	"os"

	"github.com/ianschenck/envflag"
	"github.com/urfave/cli"

	"github.com/Dataman-Cloud/baker/baker/cmd"
)

func main() {
	envflag.Parse()

	app := cli.NewApp()
	app.Name = "baker"
	app.Usage = "command line utility"

	app.Commands = []cli.Command{
		cmd.ServerCmd,
		cmd.DisConfCmd,
		cmd.BuildpackCmd,
		cmd.CanaryCmd,
		cmd.RollbackCmd,
	}

	app.Run(os.Args)
}
