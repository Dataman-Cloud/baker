package cmd

import (
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/Dataman-Cloud/baker/util"
)

var DisConfUnzipCmd = cli.Command{
	Name:  "unzip",
	Usage: "unzip config files in local",
	Action: func(c *cli.Context) {
		if err := disConfUnzip(c); err != nil {
			logrus.Fatal(err)
		}
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "file",
			Usage: "file name",
		},
		cli.StringFlag{
			Name:  "path",
			Usage: "path name",
		},
	},
}

func disConfUnzip(c *cli.Context) error {
	file := c.String("file")
	if file == "" {
		logrus.Fatal("no file input")
		return errors.New("no file input")
	}
	path := c.String("path")
	if path == "" {
		logrus.Fatal("no path input")
		return errors.New("no path input")
	}
	err := util.Unzip(file, path)
	if err != nil {
		logrus.Error("error file unzip: %s", err)
		return err
	}
	return nil
}
