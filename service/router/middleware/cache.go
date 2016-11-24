package middleware

import (
	"github.com/Dataman-Cloud/baker/cache"
	"github.com/Dataman-Cloud/baker/config"

	"github.com/gin-gonic/gin"
)

// Cache is a middleware function that initializes the Cache and attaches to
// the context of every http.Request.
func Cache(cf *config.Config) gin.HandlerFunc {
	v := setupCache(cf)
	return func(c *gin.Context) {
		c.Set("cache", v)
	}
}

// helper function to create the cache from the CLI context.
func setupCache(cf *config.Config) cache.Cache {
	return cache.NewCache()
}
