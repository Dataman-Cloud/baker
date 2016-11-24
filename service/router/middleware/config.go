package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/Dataman-Cloud/baker/config"
)

// Config is a middleware function that initializes the Configuration and
// attaches to the context of every http.Request.
func Config(cf *config.Config) gin.HandlerFunc {
	v := setupConfig(cf)
	return func(c *gin.Context) {
		c.Set("config", v)
	}
}

// helper function to create the configuration from the CLI context.
func setupConfig(cf *config.Config) *config.Config {
	return cf
}
