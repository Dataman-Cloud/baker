package middleware

import (
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"

	"github.com/Dataman-Cloud/baker/config"
)

// Config is a middleware function that initializes the Configuration and
// attaches to the context of every http.Request.
func Config(cli *cli.Context) gin.HandlerFunc {
	v := setupConfig(cli)
	return func(c *gin.Context) {
		c.Set("config", v)
	}
}

// helper function to create the configuration from the CLI context.
func setupConfig(c *cli.Context) *config.Config {
	// read the configuration
	logrus.Infof("Configuration path: %s", c.String("config-path"))
	cf, err := config.Decode(c.String("config-path"))
	if err != nil {
		logrus.Infof("Configuration error: %s", err)
		return nil
	}
	return cf
}
