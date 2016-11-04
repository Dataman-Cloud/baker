package cmd

import (
	"log"

	"github.com/urfave/cli"
)

var CanaryCmd = cli.Command{
	Name:  "canary",
	Usage: "execute a rolling upgrade.",
	Action: func(c *cli.Context) {
		if err := canary(c); err != nil {
			log.Fatalln(err)
		}
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "r, release",
			Value: ".release.yml",
			Usage: "release file",
		},
	},
}

func canary(c *cli.Context) error {
	// execute a rolling upgrade job.

	return nil
}
