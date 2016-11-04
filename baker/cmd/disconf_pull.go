package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var DisConfPullCmd = cli.Command{
	Name:  "pull",
	Usage: "pull config files from config management",
	Action: func(c *cli.Context) {
		if err := disConfPull(c); err != nil {
			logrus.Fatal(err)
		}
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "path",
			Usage: "pull config files in path",
		},
	},
}

func disConfPull(c *cli.Context) error {
	logrus.Infof("download config files in the path: %s", c.String("path"))
	return nil
}
