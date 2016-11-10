package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var BuildpackDockerfilePullCmd = cli.Command{
	Name:  "pull",
	Usage: "execute pull dockerfile in fileserver",
	Action: func(c *cli.Context) {
		if err := buildpackDockerfilePull(c); err != nil {
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

func buildpackDockerfilePull(c *cli.Context) error {
	return nil
}
