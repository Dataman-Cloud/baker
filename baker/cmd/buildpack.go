package cmd

import (
	"github.com/urfave/cli"
)

var BuildpackCmd = cli.Command{
	Name:  "buildpack",
	Usage: "execute buildpack job.",
	Subcommands: []cli.Command{
		BuildpackImportCmd,     // import app files in baker fileserver.
		BuildpackListCmd,       // list app files in baker fileserver.
		BuildpackDelCmd,        // delete app files in baker fileserver.
		BuildpackDockerfileCmd, // pull or push dockerfile in baker fileserver.
		BuildpackImageCmd,      // pull or push image in docker registry.
	},
}
