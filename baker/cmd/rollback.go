package cmd

import (
	"log"

	"github.com/urfave/cli"
)

var RollbackCmd = cli.Command{
	Name:  "rollback",
	Usage: "execute a rollback job",
	Action: func(c *cli.Context) {
		if err := rollback(c); err != nil {
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

func rollback(c *cli.Context) error {
	// execute a rollback job.

	return nil
}
