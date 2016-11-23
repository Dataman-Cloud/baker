package cmd

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/contrib/ginrus"
	"github.com/urfave/cli"

	"github.com/Dataman-Cloud/baker/service/router"
	"github.com/Dataman-Cloud/baker/service/router/middleware"
)

var ServerCmd = cli.Command{
	Name:  "server",
	Usage: "starts the baker server daemon",
	Action: func(c *cli.Context) {
		if err := server(c); err != nil {
			logrus.Fatal(err)
		}
	},
	Flags: []cli.Flag{
		cli.BoolFlag{
			EnvVar: "BAKER_DEBUG",
			Name:   "debug",
			Usage:  "start the server in debug mode",
		},
		cli.StringFlag{
			EnvVar: "BAKER_SERVER_ADDR",
			Name:   "server-addr",
			Usage:  "server address",
			Value:  ":8000",
		},
		cli.StringFlag{
			EnvVar: "BAKER_SERVER_CERT",
			Name:   "server-cert",
			Usage:  "server ssl cert",
		},
		cli.StringFlag{
			EnvVar: "BAKER_SERVER_KEY",
			Name:   "server-key",
			Usage:  "server ssl key",
		},
		cli.StringFlag{
			EnvVar: "BAKER_CONFIG_PATH",
			Name:   "config-path",
			Value:  "config.yml",
		},
	},
}

func server(c *cli.Context) error {
	// debug level if requested by user
	if c.Bool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.WarnLevel)
	}

	// setup the server and start the listener
	handler := router.Load(
		ginrus.Ginrus(logrus.StandardLogger(), time.RFC3339, true),
		middleware.Config(c),
		middleware.Store(c),
		middleware.Cache(c),
		middleware.BakeWorkPool(c),
	)

	// start the server with tls enabled
	if c.String("server-cert") != "" {
		return http.ListenAndServeTLS(
			c.String("server-addr"),
			c.String("server-cert"),
			c.String("server-key"),
			handler,
		)
	}

	// start the server without tls enabled
	return http.ListenAndServe(
		c.String("server-addr"),
		handler,
	)
}
