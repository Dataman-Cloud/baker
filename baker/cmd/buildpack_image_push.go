package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var BuildpackImagePushCmd = cli.Command{
	Name:  "push",
	Usage: "execute push image to docker registry",
	Action: func(c *cli.Context) {
		if err := buildpackImagePush(c); err != nil {
			logrus.Fatal(err)
		}
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "app name",
		},
	},
}

func buildpackImagePush(c *cli.Context) error {
	return nil
}
