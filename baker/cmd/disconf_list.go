package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var DisConfListCmd = cli.Command{
	Name:  "list",
	Usage: "list config files from config management",
	Action: func(c *cli.Context) {
		if err := disConfList(c); err != nil {
			logrus.Fatal(err)
		}
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "path",
			Usage: "list config files in path",
		},
	},
}

func disConfList(c *cli.Context) error {
	logrus.Infof("list config files in the path: %s", c.String("path"))
	return nil
}
