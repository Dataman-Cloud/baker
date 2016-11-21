package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"

	"github.com/Dataman-Cloud/baker/config"
	"github.com/Dataman-Cloud/baker/model"
	"github.com/Dataman-Cloud/baker/store"
)

// Store is a middleware function that initializes the Store and attaches to
// the context of every http.Request.
func Store(cli *cli.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		v := setupStaticUsersStore(c)
		store.ToContext(c, v)
		c.Next()
	}
}

// helper function to create the Store from the CLI context config-path.
func setupStaticUsersStore(c *gin.Context) store.Store {
	config := c.MustGet("config").(*config.Config)
	// setup Store
	if &config.Users != nil {
		staticUsers := make(map[string]*model.StaticUser)
		for k, v := range config.Users {
			staticUsers[k] = &model.StaticUser{
				Username: k,
				Password: v.Password,
			}
		}
		return store.NewStaticUsersStore(staticUsers)
	}
	return nil
}
