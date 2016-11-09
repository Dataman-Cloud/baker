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
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "s, server",
			Usage:  "server location",
			EnvVar: "BAKER_SERVER",
			Value:  ":8000",
		},
		cli.StringFlag{
			Name:   "u, username",
			Usage:  "username",
			EnvVar: "BAKER_USERNAME",
			Value:  "admin",
		},
		cli.StringFlag{
			Name:   "p, password",
			Usage:  "password",
			EnvVar: "BAKER_PASSWORD",
			Value:  "badmin",
		},
	}

	app.Commands = []cli.Command{
		cmd.ServerCmd,
		cmd.DisConfCmd,
		cmd.BuildpackCmd,
		cmd.CanaryCmd,
		cmd.RollbackCmd,
	}

	app.Run(os.Args)
}
