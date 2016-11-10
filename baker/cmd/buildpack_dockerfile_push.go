package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var BuildpackDockerfilePushCmd = cli.Command{
	Name:  "push",
	Usage: "execute push dockerfile in fileserver",
	Action: func(c *cli.Context) {
		if err := buildpackDockerfilePush(c); err != nil {
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

func buildpackDockerfilePush(c *cli.Context) error {
	return nil
}
