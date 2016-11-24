package middleware

import (
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"

	"github.com/Dataman-Cloud/baker/config"
	"github.com/Dataman-Cloud/baker/executor"
)

// BakeWorkPool is a middleware function that initializes the BakeWorkPool and attaches to
// the context of every http.Request.
func BakeWorkPool(cli *cli.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		v := setupBakeWorkPool(c)
		c.Set("bakeworkpool", v)
	}
}

// helper function to create the bakeWorkPool from the CLI context.
func setupBakeWorkPool(c *gin.Context) *executor.WorkPool {
	cf := c.MustGet("config").(*config.Config)
	workPool := cf.WorkPool
	imagePushWorkPoolOptions := workPool["imagepush"]
	imagePushWorkPoolSize, _ := strconv.Atoi(imagePushWorkPoolOptions.MaxWorkers)
	if imagePushWorkPoolSize < 1 {
		logrus.Errorf("must provide positive maxWorkers; provided %d", imagePushWorkPoolSize)
		return nil
	}
	pool, err := executor.NewWorkPool(imagePushWorkPoolSize)
	if err != nil {
		logrus.Errorf("error create bakeworkpool: %s", err)
		return nil
	}
	return pool
}
