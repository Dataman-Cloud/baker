package cmd

import (
	"log"

	"github.com/urfave/cli"
)

var BuildpackCmd = cli.Command{
	Name:  "buildpack",
	Usage: "execute buildpack job from source code.",
	Action: func(c *cli.Context) {
		if err := build(c); err != nil {
			log.Fatalln(err)
		}
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "n, name",
			Usage: "image name",
		},
		cli.StringFlag{
			Name:  "r, release",
			Value: ".release.yml",
			Usage: "release file",
		},
	},
}

func build(c *cli.Context) error {
	// build image from source code.

	return nil
}
