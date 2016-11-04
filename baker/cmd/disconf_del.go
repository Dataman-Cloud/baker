package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var DisConfDelCmd = cli.Command{
	Name:  "del",
	Usage: "delele config files from config management",
	Action: func(c *cli.Context) {
		if err := disConfDel(c); err != nil {
			logrus.Fatal(err)
		}
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "path",
			Usage: "delete config files in path",
		},
	},
}

func disConfDel(c *cli.Context) error {
	logrus.Infof("delete config files in the path: %s", c.String("path"))
	return nil
}
