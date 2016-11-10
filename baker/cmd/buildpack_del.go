package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var BuildpackDelCmd = cli.Command{
	Name:  "del",
	Usage: "delete app files in fileserver.",
	Action: func(c *cli.Context) {
		if err := buildpackDel(c); err != nil {
			logrus.Fatal(err)
		}
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "path",
			Usage: "delete files in the path",
		},
	},
}

func buildpackDel(c *cli.Context) error {
	return nil
}
