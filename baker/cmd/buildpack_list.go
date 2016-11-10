package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var BuildpackListCmd = cli.Command{
	Name:  "list",
	Usage: "list app files in fileserver.",
	Action: func(c *cli.Context) {
		if err := buildpackList(c); err != nil {
			logrus.Fatal(err)
		}
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "path",
			Usage: "list app files in the path.",
		},
	},
}

func buildpackList(c *cli.Context) error {
	return nil
}
