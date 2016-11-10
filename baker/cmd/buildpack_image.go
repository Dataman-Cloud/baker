package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var BuildpackImageCmd = cli.Command{
	Name:  "image",
	Usage: "execute dockerimage command",
	Subcommands: []cli.Command{
		BuildpackImagePushCmd, // push docker image to docker registry.
	},
}
