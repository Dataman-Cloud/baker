package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var BuildpackImportCmd = cli.Command{
	Name:  "import",
	Usage: "import app files into baker fileserver.",
	Action: func(c *cli.Context) {
		if err := buildpackImport(c); err != nil {
			logrus.Fatal(err)
		}
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "app name",
		},
		cli.StringFlag{
			Name:  "from",
			Usage: "base image",
		},
		cli.StringFlag{
			Name:  "binaryFile",
			Usage: "binary file",
		},
		cli.StringFlag{
			Name:  "binaryPath",
			Usage: "container path of binary file.",
		},
		cli.StringFlag{
			Name:  "propsFile",
			Usage: "zip file of property files.",
		},
		cli.StringFlag{
			Name:  "startupFile",
			Usage: "startup script file",
		},
		cli.StringFlag{
			Name:  "startCmd",
			Usage: "startup command",
		},
		cli.StringFlag{
			Name:  "reloadCmd",
			Usage: "reload command",
		},
	},
}

func buildpackImport(c *cli.Context) error {
	return nil
}
